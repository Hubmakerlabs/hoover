package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/utils"
	"github.com/Hubmakerlabs/replicatr/pkg/interrupt"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "insufficient arguments for testharness:\n\n"+
			"usage: testharness <arlocal endpoint address> <wallet address> <balance target in AR>\n")
		os.Exit(1)
	}
	endpoint := os.Args[1]
	address := os.Args[2]
	balanceTarget, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse balance target: '%s' %s\n", os.Args[3],
			err.Error)
		os.Exit(1)
	}
	c, cancel := context.WithCancel(context.Background())
	interrupt.AddHandler(cancel)
	TestHarness(c, cancel, endpoint, address, balanceTarget)
}

// TestHarness is a helper for running a test/demonstration for hoover social network protocol
// event data bundler/uploader.
//
// In order to use this, you need to install arlocal https://github.com/textury/arlocal - this
// can be done with nodejs 20+ installed and npx:
//
//	npx arlocal
//
// This will install and run arlocal, which provides a simple virtual arweave gateway dev
// testnet which has a number of features that enable minting tokens and triggering the mining
// of blocks using web service endpoints.
//
// This test harness triggers a mine of a block once a second as well as maintaining the balance
// of a provided wallet address to be amply sufficient for most test cases
func TestHarness(
	c context.Context,
	cancel context.CancelFunc,
	endpoint, address string,
	balanceTarget int,
) {
	var err error
	bt := big.NewFloat(float64(balanceTarget))
	var bal *big.Int
	if bal, err = GetBalance(endpoint, address); err != nil {
		fmt.Fprintf(os.Stderr, "failed to get balance: '%s'\n", err)
		// probably arlocal not running
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "balance %d winstons\n", bal)
	wt := utils.ARToWinston(bt)
	diff := wt.Sub(wt, bal)
	if diff.Cmp(big.NewInt(1000000)) == 1 {
		if err = Mint(endpoint, address, diff); err != nil {
			fmt.Fprintf(os.Stderr, "failed to mint: '%s'\n", err.Error())
			// probably arlocal not running
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "balance supplemented to %d AR with %d winstons\n",
			balanceTarget, diff)
	}
	ticker := time.NewTicker(time.Second)
out:
	for {
		select {
		case <-c.Done():
			break out
		case <-ticker.C:
			// check balance and bump
			if bal, err = GetBalance(endpoint, address); err != nil {
				fmt.Fprintf(os.Stderr, "failed to get balance: '%s'\n", err)
				// probably arlocal not running
				cancel()
				continue // will now select on Done
			}
			diff = wt.Sub(wt, bal)
			if diff.Cmp(big.NewInt(10000)) == 1 {
				if err = Mint(endpoint, address, diff); err != nil {
					fmt.Fprintf(os.Stderr, "failed to mint: '%s'\n", err.Error())
					// probably arlocal not running
					cancel()
					continue // will now select on Done
				}
				fmt.Fprintf(os.Stderr, "balance supplemented to %d with %d winstons\n",
					balanceTarget, diff)
			}
			// mine
			var height int
			if height, err = Mine(endpoint); err != nil {
				fmt.Fprintf(os.Stderr, "failed to mine: '%s'\n", err.Error())
				// probably arlocal not running
				cancel()
				continue // will now select on Done
			}
			fmt.Fprintf(os.Stderr, "block height: %d balance: %d winstons\n", height, bal)
		}
	}
}

// Mine triggers arlocal to commit the pending transactions into a block
func Mine(endpoint string) (height int, err error) {
	var resp *http.Response
	if resp, err = http.Get(fmt.Sprintf("%s/mine", endpoint)); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err.Error())
		return
	}
	var body []byte
	if body, err = GetBody(resp); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err.Error())
		return
	}
	var res map[string]interface{}
	if err = json.Unmarshal(body, &res); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err.Error())
		return
	}
	if block, ok := res["blocks"]; ok {
		if h, ok := block.(float64); ok {
			height = int(h)
		}
	}
	resp.Body.Close()
	return
}

func Mint(endpoint, account string, amount *big.Int) (err error) {
	// var res *http.Response
	// var body []byte
	_, err = http.Get(fmt.Sprintf("%s/mint/%s/%d",
		endpoint, account, amount))
	if err != nil {
		return
	}
	// if body, err = GetBody(res); err != nil {
	// 	return
	// }
	return
}

// GetBody reads out the body from a http.Response
func GetBody(res *http.Response) (body []byte, err error) {
	if res == nil {
		return
	}
	if res.Body == nil {
		return
	}
	defer res.Body.Close()
	body, err = io.ReadAll(res.Body)
	if err != nil {
		return
	}
	return
}

func GetBalance(endpoint, address string) (bal *big.Int, err error) {
	var res *http.Response
	uri := fmt.Sprintf("%s/wallet/%s/balance", endpoint, address)
	res, err = http.Get(uri)
	if err != nil {
		return
	}
	var body []byte
	if body, err = GetBody(res); err != nil {
		return
	}
	amt := big.NewInt(0)
	err = amt.UnmarshalText(body)
	return amt, nil
}
