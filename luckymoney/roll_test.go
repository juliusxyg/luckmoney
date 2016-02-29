package luckymoney

import (
	"testing"
	"math"
)

const EPSILON float64 = 0.00000001

//program mark - unittest: go test -v luckymoneysrv/luckymoney
func TestDistributeNormal(t *testing.T) {
	money := 99.99
	number := 10
	var id uint64 = 1
	Distribute(id, money, number)

	if _, ok := TableEnvelopes[id]; !ok {
		t.Error("Distribute Normal Remain list failed")
	}

	var moneyPiece float64
	var total int

	for _, cursor := range TableEnvelopes[id].Pieces {
		moneyPiece += cursor.Money
		total++
	}

	if total != number {
		t.Error("Distribute Normal total number incorrect")
	}

	if math.Abs(moneyPiece-money) > EPSILON {
		t.Errorf("Distribute Normal total money incorrect, %0.2f != %0.2f", moneyPiece, money)
	}
}

func TestDistributeZeroMoney(t *testing.T) {
	money := 0.00
	number := 10
	var id uint64 = 1

	res := Distribute(id, money, number)

	if res != false {
		t.Error("Distribute Zero Money failed")
	}
}

func TestDistributeNegtiveMoney(t *testing.T) {
	money := -100.00
	number := 10
	var id uint64 = 1

	res := Distribute(id, money, number)

	if res != false {
		t.Error("Distribute Negtive Money failed")
	}
}

func TestDistributeZeroNumber(t *testing.T) {
	money := 100.00
	number := 0
	var id uint64 = 1

	res := Distribute(id, money, number)

	if res != false {
		t.Error("Distribute Zero Number failed")
	}
}

func TestDistributeNegtiveNumber(t *testing.T) {
	money := 100.00
	number := -10
	var id uint64 = 1

	res := Distribute(id, money, number)

	if res != false {
		t.Error("Distribute Negtive Number failed")
	}
}

func TestOpenRandomNormal(t *testing.T) {
	money := 99.99
	number := 5
	var id uint64 = 1
	Distribute(id, money, number)

	envelop := TableEnvelopes[id].OpenRandom("Julius")

	if 0 == TableEnvelopes[id].Opened {
		t.Error("Open Random Opened list failed")
	}

	if envelop.Money < MIN_AMOUNT {
		t.Error("Open Random rolled out money incorrect")
	}

	if envelop.Grabber != "Julius" {
		t.Error("Open Random rolled out grabber incorrect")
	}

	if envelop.GrabTime == 0 {
		t.Error("Open Random rolled out grabber time incorrect")
	}

	envelop2 := TableEnvelopes[id].OpenRandom("Julius2")
	envelop3 := TableEnvelopes[id].OpenRandom("Julius3")
	envelop4 := TableEnvelopes[id].OpenRandom("Julius4")
	envelop5 := TableEnvelopes[id].OpenRandom("Julius5")

	if TableEnvelopes[id].Opened != TableEnvelopes[id].Size {
		t.Error("Open Random Remain list exists after all rolled out")
	}

	moneyTotal := envelop.Money + envelop2.Money + envelop3.Money + envelop4.Money + envelop5.Money
	if math.Abs(moneyTotal-money) > EPSILON {
		t.Errorf("Open Random total money incorrect, %0.2f != %0.2f", moneyTotal, money)
	}
}

func TestOpenRandomNoName(t *testing.T) {
	money := 99.99
	number := 5
	var id uint64 = 1
	Distribute(id, money, number)

	envelop := TableEnvelopes[id].OpenRandom("")

	if envelop != nil {
		t.Error("Open Random no name failed")
	}
}

//program mark - benchmark: go test -bench=. luckymoneysrv/luckymoney
var bench_result bool
func benchmarkDistribute(number int, b *testing.B) {
	money := 999.99
	var id uint64 = 1
	var r bool
	
	for n := 0; n < b.N; n++ {
		// always record the result of Fib to prevent
    // the compiler eliminating the function call.
    r = Distribute(id, money, number)
  }
  // always store the result to a package level variable
  // so the compiler cannot eliminate the Benchmark itself.
  bench_result = r
}

func BenchmarkDistribute10(b *testing.B) {
	benchmarkDistribute(10, b)
}

func BenchmarkDistribute50(b *testing.B) {
	benchmarkDistribute(50, b)
}

func BenchmarkDistribute100(b *testing.B) {
	benchmarkDistribute(100, b)
}

func BenchmarkDistribute150(b *testing.B) {
	benchmarkDistribute(150, b)
}

func BenchmarkDistribute200(b *testing.B) {
	benchmarkDistribute(200, b)
}




