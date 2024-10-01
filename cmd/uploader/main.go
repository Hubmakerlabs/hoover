package main

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/utils"
	"github.com/Hubmakerlabs/hoover/pkg/bluesky"
	"github.com/Hubmakerlabs/hoover/pkg/farcaster"
	"github.com/Hubmakerlabs/hoover/pkg/multi"
	"github.com/Hubmakerlabs/hoover/pkg/nostr"
	"github.com/Hubmakerlabs/replicatr/pkg/interrupt"
)

func main() {
	address, wallet, err := GetTestWallet(arlocal)
	if err != nil {
		os.Exit(1)
	}
	balanceTarget := utils.ARToWinston(big.NewFloat(100000))
	BumpBalance(arlocal, address, balanceTarget)
	c, cancel := context.WithCancel(context.Background())
	interrupt.AddHandler(cancel)
	// go func() {
	// 	time.Sleep(time.Second * 30)
	// 	cancel()
	// }()
	var wg sync.WaitGroup
	fmt.Println()
	var speedFactor int64 = 1
	go func() {
		tick := time.NewTicker(time.Second)
		for {
			select {
			case <-tick.C:
				Mine(arlocal)
			case <-c.Done():
				return
			}
		}
	}()
	multi.Firehose(c, cancel, &wg, nostr.Relays, bluesky.Urls, farcaster.Urls,

		func(bundle *types.BundleItem) (err error) {
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
			if _, err = wallet.SendTransaction(tx); err != nil {
				// we need to add more winstons to pay for this probably
				BumpBalance(arlocal, address, balanceTarget)
				if _, err = wallet.SendTransaction(tx); err != nil {
					// if he dies, he dies
					return nil
				}
			}
			return
		})
}

//go:embed keyfile.json
var key []byte

// this requires the use of https://github.com/textury/arlocal todo: maybe we can spawn it in a container?

const arlocal = "http://localhost:1984"

// Mine triggers arlocal to commit the pending transactions into a block
func Mine(endpoint string) (err error) {
	_, err = http.Get(
		fmt.Sprintf("%s/mine", endpoint))
	if err != nil {
		return
	}
	return
}

func Mint(endpoint, account string, amount *big.Int) (err error) {
	var res *http.Response
	var body string
	res, err = http.Get(fmt.Sprintf("%s/mint/%s/%d",
		endpoint, account, amount))
	if err != nil {
		return
	}
	if body, err = GetBody(res); err != nil {
		return
	}
	_ = body
	return
}

// BumpBalance checks that the balance is not low and adds winstons to it to bring it to a
// decent amount.
func BumpBalance(endpoint, address string, amount *big.Int) {
	var err error
	// if bal, err = GetBalance(endpoint, address); err != nil {
	// 	t.Fatal(err)
	// }
	// load the wallet to a decent starting amount of winston
	if err = Mint(endpoint, address, amount); err != nil {
		return
	}
	// mine it
	if err = Mine(endpoint); err != nil {
		return
	}
}

func GetTestWallet(endpoint string) (address string, wallet *goar.Wallet,
	err error) {
	if wallet, err = goar.NewWallet(key, endpoint); err != nil {
		return
	}
	address = wallet.Signer.Address
	return
}

// GetBody reads out the body from a http.Response
func GetBody(res *http.Response) (s string, err error) {
	if res == nil {
		return
	}
	if res.Body == nil {
		return
	}
	defer res.Body.Close()
	var body []byte
	body, err = io.ReadAll(res.Body)
	if err != nil {
		return
	}
	s = string(body)
	return
}

func GetBalance(endpoint, account string) (bal int64, err error) {
	var res *http.Response
	address := fmt.Sprintf("%s/wallet/%s/balance", endpoint, account)
	res, err = http.Get(address)
	if err != nil {
		return
	}
	var body string
	if body, err = GetBody(res); err != nil {
		return
	}
	return strconv.ParseInt(body, 10, 64)
}

func Upload() (code int, err error) {
	return
}
