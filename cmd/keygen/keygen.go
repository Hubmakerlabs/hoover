package main

import (
	_ "embed"
	"fmt"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/wallet"
	"github.com/everFinance/gojwk"
)

func main() {
	var err error
	w := wallet.GenerateWallet()
	var k B
	if k, err = gojwk.Marshal(w.Key); err != nil {
		panic(err)
	}
	fmt.Printf("%s", k)
}
