package uploader

import (
	"context"
	"fmt"
	"math/big"
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
	balanceTarget := utils.ARToWinston(big.NewFloat(100000))
	BumpBalance(arlocal, address, balanceTarget, t)
	c, cancel := context.WithCancel(context.Background())
	interrupt.AddHandler(cancel)
	go func() {
		time.Sleep(time.Second * 30)
		cancel()
	}()
	var wg sync.WaitGroup
	fmt.Println()
	var speedFactor int64 = 1
	go func() {
		tick := time.NewTicker(time.Second)
		for {
			select {
			case <-tick.C:
				Mine(arlocal, t)
			case <-c.Done():
				return
			}
		}
	}()
	multi.Firehose(c, cancel, &wg, nostr.Relays, func(bundle *types.BundleItem) (err error) {
		t.Log(bundle.Tags)
		t.Log(bundle.Data)
		tx := &types.Transaction{
			Format:   2,
			Target:   "",
			Quantity: "0",
			Tags:     utils.TagsEncode(bundle.Tags),
			Data:     utils.Base64Encode([]byte(bundle.Data)),
			DataSize: fmt.Sprintf("%d", len(bundle.Data)),
		}
		var sum int
		for i := range tx.Tags {
			sum += len(tx.Tags[i].Name) + len(tx.Tags[i].Value)
		}
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
		tx.Reward = fmt.Sprintf("%d", rew)
		t.Log("tags data", sum, "data", len(bundle.Data), "reward", tx.Reward)
		if _, err = wallet.SendTransaction(tx); err != nil {
			// we need to add more winstons to pay for this probably
			BumpBalance(arlocal, address, balanceTarget, t)
			if _, err = wallet.SendTransaction(tx); err != nil {
				t.Log("failed to fund for upload")
				// if he dies, he dies
				t.Log(err)
				return nil
			}
		}
		return
	})
}

func TestUpload(t *testing.T) {
	address, wallet, err := GetTestWallet(arlocal, t)
	BumpBalance(arlocal, address, big.NewInt(100000000000000), t)
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
}
