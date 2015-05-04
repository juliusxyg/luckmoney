package luckymoney

import (
	"math"
	"math/rand"
	"time"
)

const MIN_AMOUNT float64 = 0.01
const TWO_PI float64 = 2 * math.Pi
const RAND_MAX float64 = 2147483647
const SIGMA_FACTOR float64 = 8
const DECIMAL int = 2

type M_envelop struct {
	Money float64
	Grabber string
	GrabTime int64 

	Next *M_envelop
	Prev *M_envelop
}

type M_envelopes_remain struct {
	Head *M_envelop
	Size int
	Id int64
}

type M_envelops_opened struct {
	Head *M_envelop
	Size int
	Id int64
}

var RemainEnvelopes = make(map[int64]*M_envelopes_remain)
var OpenedEnvelops  = make(map[int64]*M_envelops_opened)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func (remain *M_envelopes_remain) OpenRandom(name string) *M_envelop {
	if remain.Head == nil {
		return nil
	}

	if len(name) == 0 {
		return nil
	}
	//从链表里随机取出一个
	randIdx := rand.Intn(remain.Size)
	itr := 0
	cursor := new(M_envelop)
	cursor = remain.Head

	for cursor != nil {
		if itr == randIdx {
			break
		}
		cursor = cursor.Next
		itr++
	} 

	cursor.Grabber = name
	cursor.GrabTime = time.Now().Unix()

	//把单元从remain链表结构中分离出
	if cursor.Next != nil {
		cursor.Next.Prev = cursor.Prev
	}
	if cursor.Prev != nil {
		cursor.Prev.Next = cursor.Next
	}
	if cursor == remain.Head {
		remain.Head = cursor.Next
	}
	remain.Size--
	//把取出的单元从头往后一次插入到opened的链表中
	opened, ok := OpenedEnvelops[remain.Id]
	if !ok {
		opened = new(M_envelops_opened)
	  opened.Head = nil
	  opened.Size = 0
	  opened.Id = remain.Id
		OpenedEnvelops[opened.Id] = opened
	}

	cursor.Next = opened.Head
	cursor.Prev = nil
	if opened.Head != nil {
		opened.Head.Prev = cursor
	}
	opened.Head = cursor
	opened.Size++
	//维护索引列表
	if remain.Head == nil {
		delete(RemainEnvelopes, remain.Id)
	}

	return cursor
}

func Distribute(money float64, number int) int64 {
	if money <= 0 || number <= 0 {
		return -1
	}
	moneyLeft := money - float64(number) * MIN_AMOUNT
	var mu float64
	var sigma float64
	var noise float64
	//初始化两个链表
	id_of_envelops := time.Now().UTC().UnixNano();

  envelops := new(M_envelopes_remain)
  envelops.Head = nil
  envelops.Size = number
  envelops.Id = id_of_envelops

  RemainEnvelopes[envelops.Id] = envelops

	for i := 0; i<number; i++ {
		mu = moneyLeft / float64(number - i)
		sigma = mu / SIGMA_FACTOR
		noise = generateNoise(mu, sigma)

		if noise < 0 {
			noise = 0
		}

		if noise > moneyLeft {
			noise = moneyLeft
		}
		//构建单元，并从头开始往后依次插入remain链表中，第一个插入的是尾
		envelop := new(M_envelop)
		envelop.Money = noise + MIN_AMOUNT
		envelop.Grabber = ""
		envelop.GrabTime = 0
		envelop.Next = envelops.Head
		envelop.Prev = nil

		if envelops.Head != nil {
			envelops.Head.Prev = envelop
		}
		envelops.Head = envelop

		moneyLeft -= noise
	}

	return id_of_envelops
}

//期望值mu 和 标准差sigma， mu 为当前红包的均值，
//当到第i个红包是所剩金额是totalMoneyLeftAt_i ＝ totalMoney-0.01*numberOfEnvelopes - moneyInEnvelope[0] - ... - moneyInEnvelope[i-1]，
// mu = totalMoneyLeftAt_i / (numberOfEnvelopes - i)
//截尾正态分布 红包金额范围是[0, totalMoneyLeftAt_i]
func generateNoise(mu float64, sigma float64) float64 {
	haveSpare := false
	var rand1 float64
	var rand2 float64

	if haveSpare {
		haveSpare = false
		return (sigma * math.Sqrt(rand1) * math.Sin(rand2) + mu)
	}

	haveSpare = true
	rand1 = rand.Float64() / RAND_MAX
	if rand1 < math.Pow10(-100) {
		rand1 = math.Pow10(-100)
	}
	rand1 = -2 * math.Log(rand1)
	rand2 = ( rand.Float64() / RAND_MAX ) * TWO_PI

	return toFixed(( sigma * math.Sqrt(rand1) * math.Cos(rand2) + mu ), DECIMAL)
}

func round(num float64) int {
    return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
    output := math.Pow(10, float64(precision))
    return float64(round(num * output)) / output
}