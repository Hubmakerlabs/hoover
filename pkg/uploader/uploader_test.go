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
	BumpBalance(arlocal, address,10000000000, t)
	c, cancel := context.WithCancel(context.Background())
	interrupt.AddHandler(cancel)
	go func() {
		time.Sleep(time.Second * 10)
		cancel()
	}()
	var wg sync.WaitGroup
	fmt.Println()
	var speedFactor int64 = 1
	var count byte
	multi.Firehose(c, cancel, &wg, nostr.Relays, func(bundle *types.BundleItem) (err error) {
		var reward int64
		reward, err = wallet.Client.GetTransactionPrice(len(bundle.Data), nil)
		if err != nil {
			// if he dies, he dies
			return nil
		}
		rew := reward * (100 + speedFactor) / 100
		if rew == 0 {
			rew = 1000
		}
		tx := &types.Transaction{
			Format:   2,
			Target:   "",
			Quantity: "0",
			Tags:     utils.TagsEncode(bundle.Tags),
			Data:     utils.Base64Encode([]byte(bundle.Data)),
			DataSize: fmt.Sprintf("%d", len(bundle.Data)),
			Reward:   fmt.Sprintf("%d", rew),
		}
		if _, err = wallet.SendTransaction(tx); err != nil {
			// if he dies, he dies
			t.Log(err)
			return nil
		}
		count++
		if count==0{
			Mine(arlocal, t)
		}
		return
	})
}

func TestUpload(t *testing.T) {
	address, wallet, err := GetTestWallet(arlocal, t)
	BumpBalance(arlocal, address,100000000000, t)
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
