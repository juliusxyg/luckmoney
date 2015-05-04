package main 

import (
	"fmt"
	"luckymoney"
	"net"
	"os"
	"regexp"
	"strconv"
	"io"
	"log"
)

func main() {
	//socket srv
	srvAddr := "127.0.0.1:9001"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", srvAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

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

	buf := make([]byte, 128)
	for {
		n, err := conn.Read(buf[0:])
		if err == io.EOF {
			log.Printf("connection closed\r\n")
			return
		}else if err != nil {
			log.Printf("Read Error: %s\r\n", err.Error())
			return
		}

		log.Printf("Read %d bytes\r\n", n)

		s := string(buf[0:n])
		if len(s) < 3 { //protect slice boundary
			continue
		}
		//SET 10000<money, cents, integer> 10<number, integer>
		if s[0:3] == "SET" {
			doSet(conn, s[4:])
		}else if s[0:3] == "GET" {
			doGet(conn, s[4:])
		}

	}
}

func doSet(conn net.Conn, data string) {
	log.Printf("SET %s\r\n", data)

	reg := regexp.MustCompile("[0-9]+")
	dataArr := reg.FindAllString(data, -1)
	if len(dataArr) != 2 {
		_, err := conn.Write([]byte("bad SET parameter 1\r\n"))
		if err != nil {
			return
		}
		return
	}

	if isMatched, _ := regexp.MatchString("^[0-9]+$", dataArr[0]); !isMatched {
		_, err := conn.Write([]byte("bad SET parameter 2\r\n"))
		if err != nil {
			return
		}
		return
	}

	if isMatched, _ := regexp.MatchString("^[0-9]+$", dataArr[1]); !isMatched {
		_, err := conn.Write([]byte("bad SET parameter 3\r\n"))
		if err != nil {
			return
		}
		return
	}

	log.Printf("set pool for money: %s, number: %s.\r\n", dataArr[0], dataArr[1])

	money, _ := strconv.ParseFloat(dataArr[0], 64)
	number, _ := strconv.Atoi(dataArr[1])
	id := luckymoney.Distribute(money, number)

	_, err := conn.Write([]byte(fmt.Sprintf("id: %d\r\n",id)))
	if err != nil {
		return
	}
}

func doGet(conn net.Conn, data string) {
	log.Printf("GET %s\r\n", data)

	reg := regexp.MustCompile("[0-9]+")
	idStr := reg.FindString(data)

	if isMatched, _ := regexp.MatchString("^[0-9]+$", idStr); !isMatched {
		_, err := conn.Write([]byte("bad id\r\n"))
		if err != nil {
			return
		}
		return
	}

	reg = regexp.MustCompile("[a-zA-Z]+")
	nameStr := reg.FindString(data)

	if isMatched, _ := regexp.MatchString("^[a-zA-Z]+$", nameStr); !isMatched {
		_, err := conn.Write([]byte("bad name\r\n"))
		if err != nil {
			return
		}
		return
	}

	id, _ := strconv.ParseInt(idStr, 10, 64)

	if envelop, ok := luckymoney.RemainEnvelopes[id]; ok {
		opened := envelop.OpenRandom(nameStr)

		log.Printf("grabber : %s, money: %.2f, timestamp: %d\r\n", opened.Grabber, opened.Money, opened.GrabTime)
		log.Printf("remain envelop: %d\r\n", len(luckymoney.RemainEnvelopes))
		log.Printf("opened envelop: %d\r\n", len(luckymoney.OpenedEnvelops))

		_, err := conn.Write([]byte(fmt.Sprintf("You get: %.2f\r\n",opened.Money)))
		if err != nil {
			return
		}
		return
	}else{
		_, err := conn.Write([]byte("no such envelop or envelop are all taken\r\n"))
		if err != nil {
			return
		}
		return
	}
}