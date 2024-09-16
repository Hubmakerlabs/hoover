package uploader

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/utils"
	"github.com/Hubmakerlabs/hoover/pkg/multi"
	"github.com/Hubmakerlabs/hoover/pkg/nostr"
	"github.com/Hubmakerlabs/replicatr/pkg/interrupt"
	"github.com/davecgh/go-spew/spew"
	"lukechampine.com/frand"
)

func TestMultiFirehose(t *testing.T) {
	address, wallet, err := GetTestWallet(arlocal, t)
	if err != nil {
		t.Fatal(err)
	}
	BumpBalance(arlocal, address, t)
	c, cancel := context.WithCancel(context.Background())
	interrupt.AddHandler(cancel)
	go func() {
		time.Sleep(time.Second * 10)
		cancel()
	}()
	var wg sync.WaitGroup
	fmt.Println()
	multi.Firehose(c, cancel, &wg, nostr.Relays, func(bundle *types.BundleItem) (err error) {
		// var tx types.Transaction
		if _, err = wallet.SendData([]byte(bundle.Data), bundle.Tags); err != nil {
			// just continue because this is a test
			return nil
		}
		Mine(arlocal, t)
		return
	})
}

func TestUpload(t *testing.T) {
	address, wallet, err := GetTestWallet(arlocal, t)
	BumpBalance(arlocal, address, t)
	b := make([]byte, 32)
	if _, err = frand.Read(b); err != nil {
		t.Fatal(err)
	}
	data := utils.Base64Encode(b)
	spew.Dump(b, data)
	var reward int64
	if reward, err = wallet.Client.GetTransactionPrice(len(data), nil); err != nil {
		t.Fatal(err)
	}
	var speedFactor int64
	tx := &types.Transaction{
		Format: 2,
		ID:     "",
		LastTx: "",
		Owner:  "",
		Tags: utils.TagsEncode([]types.Tag{{Name: "Name-Tag",
			Value: "this is a test tag"}}),
		Target:     "",
		Quantity:   "0",
		Data:       data,
		DataReader: nil,
		DataSize:   fmt.Sprintf("%d", len(data)),
		DataRoot:   "",
		Reward:     fmt.Sprintf("%d", reward*(100+speedFactor)/100),
		Signature:  "",
		Chunks:     nil,
	}
	spew.Dump(tx)
	if *tx, err = wallet.SendTransaction(tx); err != nil {
		t.Fatal(err)
	}
	if err = Mine(arlocal, t); err != nil {
		t.Fatal(err)
	}
	//
	// res, err = http.Get(
	// 	"http://localhost:1984/mine")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// // spew.Dump(res)
	//
	// // spew.Dump(tx)
	// if *tx, err = wallet.SendData([]byte("aoeu"), []types.Tag{{Name: "Name",
	// 	Value: "testing testing 1 2 3"}}); err != nil {
	// 	t.Fatal(err)
	// }
	// spew.Dump(tx)
	//
	// res, err = http.Get(
	// 	"http://localhost:1984/mine")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// spew.Dump(res)

}
