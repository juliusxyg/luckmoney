package main 

import (
	"luckymoneysrv/luckymoney"
	"luckymoneysrv/misc"
	"net"
	"os"
	"io"
	_"time"
	//profiling
	"github.com/davecheney/profile"
	//codec
	"bufio"
	"bytes"
	"encoding/binary"
	//json-api
	"strconv"
	"encoding/json"
	//misc
	_"fmt"
	//mongodb client
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//mutex
	"sync"
)

const serverAddr string = "127.0.0.1:9001"
const mongodbAddr string = "mongodb://127.0.0.1:27017"

func main() {
	defer profile.Start(profile.CPUProfile).Stop()
	/**
	 start uuid service
	 */
	err := misc.StartUUId()
	if err!=nil {
		misc.ERR("Fatal error: ", err.Error())
    os.Exit(0)
	}
	/** 
	 check mongo
	 */
	misc.INFO("checking mongodb connection")
	session, err := mgo.Dial(mongodbAddr)
  if err != nil {
    misc.ERR("Fatal error: ", err.Error())
    os.Exit(0)
  }
  err = session.Ping()
  if err != nil {
    misc.ERR("Fatal error: ", err.Error())
    os.Exit(0)
  }
  session.Close()
  misc.INFO("mongodb checked")

	/**
	 start listening
	 */
	misc.INFO("establishing server")
	srvAddr := serverAddr
	tcpAddr, err := net.ResolveTCPAddr("tcp4", srvAddr)
	if err != nil {
		misc.ERR("Fatal error: ", err.Error())
		os.Exit(0)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		misc.ERR("Fatal error: ", err.Error())
		os.Exit(0)
	}
	misc.INFO("server started")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue;
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {		
		data, readErr := Decode(reader)

		if readErr == io.EOF {
			misc.INFO("connection closed\r\n")
			return
		}else if readErr != nil {
			misc.INFO("Read Parse Error: %s\r\n", readErr.Error())
			return
		}

		cmd := &Command{}
		//json unmarshal big number will covert to float64 and lose precison
		d := json.NewDecoder(bytes.NewBuffer(data))
    d.UseNumber()
    if err := d.Decode(cmd); err != nil {
      misc.WARN("failed to decode\r\n")
			return
    }

		switch cmd.Cmd {
		case "SET":
			doSet(conn, cmd.Args)
		case "GET":
			doGet(conn, cmd.Args)
		}

		break;
	}
}

func sendResponse(respCode *RespCode, conn net.Conn) {
	respJSON, respJSONErr := json.Marshal(respCode);
	if respJSONErr != nil {
		misc.INFO("resp json encode error")
		return
	}

	respJSONData, respJSONDataErr := Encode(respJSON);
	if respJSONDataErr != nil {
		misc.INFO("resp json codec error")
		return
	}

	_, respErr := conn.Write(respJSONData)
	if respErr != nil {
		return
	}
}

//parogram mark - SET command {cmd:"SET", args:[moneyInCents int, numberOfEnvelops int]}
func doSet(conn net.Conn, args []interface{}) {
	var ok bool
	var n json.Number
	respCode := &RespCode{}

	defer func() {
		if err := recover(); err != nil {
			respCode.Code = 1
			respCode.Message, _ = err.(string)
			respCode.Data = nil
			sendResponse(respCode, conn)
		}
	}()

	if n, ok = args[0].(json.Number); !ok {
		panic("FAILED: bad money")
	}
	moneyInCent, _ := strconv.ParseInt(string(n), 10, 64)

	if n, ok = args[1].(json.Number); !ok {
		panic("FAILED: bad number")
	}
	numberArg, _ := strconv.ParseInt(string(n), 10, 64)

	number := int(numberArg)
  money := float64(moneyInCent/100)
  id := misc.UUId()
	ok = luckymoney.Distribute(id, money, number)

	if !ok {
		misc.WARN("Distribute failed. money: ", money, " number: ", number)
		panic("FAILED: create error")
	}

	misc.DEBUG("Distribute success. money:  ", money, " number: ", number)
	misc.DEBUG("envelop id: ", id)

	respCode.Code = 0
	respCode.Message = "SUCCESS"
	respCode.Data = id
	

	sendResponse(respCode, conn)

	go storeInMgo(id)
}

//parogram mark - GET command {cmd:"GET", args:[envelopId-integer, nickname-string]}
var mutex = &sync.Mutex{}
func doGet(conn net.Conn, args []interface{}) {
	respCode := &RespCode{}

	defer func() {
		if err := recover(); err != nil {
			respCode.Code = 1
			respCode.Message, _ = err.(string)
			respCode.Data = nil
			sendResponse(respCode, conn)
		}
	}()

	n, ok := args[0].(json.Number);
	if !ok {
		panic("FAILED: bad id")
	}
	id, _ := strconv.ParseUint(string(n), 10, 64)

	name, ok := args[1].(string);
	if !ok {
		panic("FAILED: bad name")
	}

	envelop, ok := luckymoney.TableEnvelopes[id]
	if !ok {
		misc.DEBUG("to read from mongo")
		envelop = readFromMgo(id)
		if envelop != nil {
			luckymoney.TableEnvelopes[id] = envelop
			ok = true
		}
	}

	if ok {
		opened := envelop.OpenLastFind(name)
		if opened == nil {
			mutex.Lock()
			opened = envelop.OpenRandom(name)
			mutex.Unlock()
		}

		if opened == nil {
			respCode.Code = 0
			respCode.Message = "SUCCESS"
			respCode.Data = 0
		}else{
			misc.DEBUG("envelop id: ",id, " grabber: ", opened.Grabber, " money: ", opened.Money, " timestamp: ", opened.GrabTime)
			misc.DEBUG("total envelopes: ", len(luckymoney.TableEnvelopes))
			respCode.Code = 0
			respCode.Message = "SUCCESS"
			respCode.Data = opened.Money
		}

		sendResponse(respCode, conn)

		go storeInMgo(id)

		return
	}
}

//program mark - read from mongo
func readFromMgo(id uint64) *luckymoney.M_envelop {
	session, err := mgo.Dial(mongodbAddr)
  if err != nil {
    misc.ERR("[mongodb]", err)
    return nil
  }
  defer session.Close()

  session.SetMode(mgo.Monotonic, true)
  c := session.DB("luckymoney").C("envelops")

  envelop := new(luckymoney.M_envelop)

  err = c.Find(bson.M{"id": id}).One(envelop)
  if err != nil {
    misc.ERR("[mongodb]", err)
    return nil
  }

  return envelop
}

//program mark - store data in mongo aysnced, channed used limit connections
var bufferedMongoChannel = make(chan bool, 100)
func storeInMgo(id uint64) {
	if envelop, ok := luckymoney.TableEnvelopes[id]; ok {
		bufferedMongoChannel <- true
		go func() {
			defer func(){ <-bufferedMongoChannel }()
			session, err := mgo.Dial(mongodbAddr)
	    if err != nil {
	      misc.ERR("[mongodb]", err)
	      return
	    }
	    defer session.Close()

	    session.SetMode(mgo.Monotonic, true)
	    c := session.DB("luckymoney").C("envelops")
	    _, err = c.Upsert(bson.M{"id": id}, envelop)
	    if err != nil {
	      misc.ERR("[mongodb]", err)
	      return
		  }

		  misc.DEBUG("save in mongo, id: ", id)

		  if envelop.Opened == envelop.Size {
		  	misc.DEBUG("all grabbed, delete from memory")
		  	delete(luckymoney.TableEnvelopes, id)
		  }
		}()
	}
}

//Command - json structure
type Command struct {
	Cmd string `json:"cmd"`
	Args []interface{} `json:"args"`
}

//RespCode - json structure
type RespCode struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

//Encode - codec
func Encode(message []byte) ([]byte, error) {
  length := int32(len(message))
  pkg := new(bytes.Buffer)
  // write body-length
  err := binary.Write(pkg, binary.LittleEndian, length)
  if err != nil {
    return nil, err
  }
  // write body-content
  err = binary.Write(pkg, binary.LittleEndian, message)
  if err != nil {
     return nil, err
  }

  return pkg.Bytes(), nil
}

//Decode - codec
func Decode(reader *bufio.Reader) ([]byte, error) {
  // get body-length binary from input and covert to fix-sized variable
  lengthByte, _ := reader.Peek(4)
  lengthBuff := bytes.NewBuffer(lengthByte)
  var length int32
  err := binary.Read(lengthBuff, binary.LittleEndian, &length)
  if err != nil {
     return nil, err
  }

  if int32(reader.Buffered()) < length+4 {
     return nil, err
  }
  // get body-content
  pkg := make([]byte, int(4+length))
  _, err = reader.Read(pkg)
  if err != nil {
     return nil, err
  }
  return pkg[4:], nil
}

