package main

import (
	"github.com/hhong0326/hhongcoin/cli"
	"github.com/hhong0326/hhongcoin/db"
)

// Concept
// B1
// 	b1Hash. = data + "x"(ifChange)
// B2
// 	b2Hash.. = data + b1Hash.(=prevHash)
// B3
// 	b3Hash... = data + b2Hash..(=prevHash)

func main() {

	// Phase I
	// hong := person.Person{}
	// hong.SetDeatails("sunil", 28)

	// fmt.Println("Main hong", hong)
	// hong.SayHello()

	// fmt.Println(hong.Name())

	// chain := blockchain.GetBlockChain()

	// chain.AddBlock("Second Block")
	// chain.AddBlock("Third Block")
	// chain.AddBlock("Fourth Block")
	// for _, b := range chain.AllBlocks() { // another way : blocks -> Blocks
	// 	fmt.Printf("Data: %s\n", b.Data)
	// 	fmt.Printf("Hash: %s\n", b.Hash)
	// 	fmt.Printf("PrevHash: %s\n", b.PrevHash)
	// }

	defer db.Close()
	cli.Start()

}
