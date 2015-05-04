package main 

import (
	"fmt"
	"flag"
	"luckymoney"
)

func main() {
	// args := os.Args[1:]

	// totalMoney, _ := strconv.ParseFloat(args[0], 64)
	// numberOfEnvelopes, _ := strconv.Atoi(args[1])

	totalMoneyPtr := flag.Int("money", 1, "an integer calculated in cents")
	numberOfEnvelopesPtr := flag.Int("number", 1, "an integer")

	flag.Parse()

	totalMoney := float64(*totalMoneyPtr/100)
	numberOfEnvelopes := *numberOfEnvelopesPtr

	id := luckymoney.Distribute(totalMoney, numberOfEnvelopes)

	moneyEnvelopes := luckymoney.RemainEnvelopes[id]
	
	fmt.Println("all envelops rolled out:")
	cursor := new(luckymoney.M_envelop)
	cursor = moneyEnvelopes.Head

	for cursor != nil {
		fmt.Print(cursor.Money, " ")
		cursor = cursor.Next
	}
	fmt.Println("")

	fmt.Println("remain index:")
	fmt.Printf("%v\n",luckymoney.RemainEnvelopes)
	fmt.Println("opened index:")
	fmt.Printf("%v\n",luckymoney.OpenedEnvelops)

	fmt.Println("6 users going to grab a random:")
	fmt.Println(moneyEnvelopes.OpenRandom("Julius"))
	fmt.Println(moneyEnvelopes.OpenRandom("Peter"))
	fmt.Println(moneyEnvelopes.OpenRandom("Rock"))
	fmt.Println(moneyEnvelopes.OpenRandom("Tim"))
	fmt.Println(moneyEnvelopes.OpenRandom("John"))
	fmt.Println(moneyEnvelopes.OpenRandom("Tiffa"))

	openedEnvelopes := luckymoney.OpenedEnvelops[id]
	fmt.Println("result list:")
	cursor = openedEnvelopes.Head

	for cursor != nil {
		fmt.Println(cursor.Money, cursor.Grabber, cursor.GrabTime)
		cursor = cursor.Next
	}

	fmt.Println("remain index:")
	fmt.Printf("%v\n",luckymoney.RemainEnvelopes)
	fmt.Println("opened index:")
	fmt.Printf("%v\n",luckymoney.OpenedEnvelops)

}



