package luckydraw

import (
	"math"
	"math/rand"
	"time"
)

const MIN_AMOUNT float64 = 0.01
const TWO_PI float64 = 2 * math.Pi
const RAND_MAX float64 = 2147483647
const SIGMA_FACTOR_MAX int = 10
const SIGMA_FACTOR_MIN int = 2
const DECIMAL int = 2

/*
data structure:
{ID, money, size, created_at, opened [piece1, piece2, ..., pieceN]}
pieceN: {i, money, grabber, grabtime}
*/
//row
type M_envelop_piece struct {
	I        int
	Money    float64
	Grabber  string
	GrabTime int64
}

type M_envelop struct {
	Id        uint64
	Money     float64
	Size      int
	CreatedAt int64
	Opened    int
	Pieces    []M_envelop_piece
}

//table
var TableEnvelopes = make(map[uint64]*M_envelop)

func init() {
	//need to set seed when init() called, otherwise rand out value will be the same
	rand.Seed(time.Now().UTC().UnixNano())
}

func (remain *M_envelop) OpenLastFind(name string) *M_envelop_piece {
	if remain.Opened == 0 {
		return nil
	}

	if len(name) == 0 {
		return nil
	}

	for _, cursor := range remain.Pieces {
		if cursor.Grabber == name {
			return &cursor
		}
	}

	return nil
}

//随机取出
func (remain *M_envelop) OpenRandom(name string) *M_envelop_piece {
	if remain.Size == remain.Opened {
		return nil
	}

	if len(name) == 0 {
		return nil
	}

	//从slice里随机取出一个，原子性
	randIdx := rand.Intn(remain.Size - remain.Opened)
	itr := 0

	for index, cursor := range remain.Pieces {
		if cursor.GrabTime > 0 {
			continue
		}
		if itr == randIdx {
			itr = index
			break
		}
		itr++
	}

	remain.Pieces[itr].Grabber = name
	remain.Pieces[itr].GrabTime = time.Now().Unix()

	remain.Opened++

	return &remain.Pieces[itr]
}

//初始化红包堆
func Distribute(id uint64, money float64, number int) bool {
	if money <= 0 || number <= 0 || id <= 0 {
		return false
	}
	moneyLeft := money - float64(number)*MIN_AMOUNT
	var mu float64
	var sigma float64
	var noise float64
	//rand a sigma factor
	sigma_factor := rand.Intn(SIGMA_FACTOR_MAX-SIGMA_FACTOR_MIN) + SIGMA_FACTOR_MIN

	envelops := new(M_envelop)
	envelops.Size = number
	envelops.Id = id
	envelops.CreatedAt = time.Now().Unix()
	envelops.Opened = 0
	envelops.Money = money
	envelops.Pieces = make([]M_envelop_piece, envelops.Size)

	TableEnvelopes[envelops.Id] = envelops

	for i := 0; i < number; i++ {
		mu = moneyLeft / float64(number-i)
		sigma = mu / float64(sigma_factor)
		noise = generateNoise(mu, sigma)

		if noise < 0 {
			noise = 0
		}

		if noise > moneyLeft {
			noise = moneyLeft
		}
		//构建单元，并从头开始往后依次插入remain链表中，第一个插入的是尾
		envelop := M_envelop_piece{}
		envelop.I = i
		envelop.Money = toFixed(noise+MIN_AMOUNT, DECIMAL)
		envelop.Grabber = ""
		envelop.GrabTime = 0

		envelops.Pieces[i] = envelop

		moneyLeft -= noise
	}

	return true
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
		return (sigma*math.Sqrt(rand1)*math.Sin(rand2) + mu)
	}

	haveSpare = true
	rand1 = rand.Float64() / RAND_MAX
	if rand1 < math.Pow10(-100) {
		rand1 = math.Pow10(-100)
	}
	rand1 = -2 * math.Log(rand1)
	rand2 = (rand.Float64() / RAND_MAX) * TWO_PI

	return toFixed((sigma*math.Sqrt(rand1)*math.Cos(rand2) + mu), DECIMAL)
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
