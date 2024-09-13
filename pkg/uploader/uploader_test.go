package uploader

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"
)

//go:embed keyfile.json
var key []byte

// this requires the use of https://github.com/textury/arlocal todo: maybe we can spawn it in a container?

const arlocal = "http://localhost:1984"

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

func Mine(t *testing.T) (err error) {
	var res *http.Response
	res, err = http.Get(
		"http://localhost:1984/mine")
	if err != nil {
		return
	}
	var body string
	if body, err = GetBody(res); err != nil {
		return
	}
	t.Log(body)
	return
}

func TestUpload(t *testing.T) {
	var err error
	var res *http.Response
	res, err = http.Get(
		"http://localhost:1984/wallet/27xHJ0MNsBUKFIdOiQ3OlrZdDzSNfBPGnp6YVmWKKxU/balance")
	if err != nil {
		t.Fatal(err)
	}
	var body string
	if body, err = GetBody(res); err != nil {
		t.Fatal(err)
	}
	var bal int64
	t.Log(body)
	bal, err = strconv.ParseInt(body, 10, 64)
	// do we need to mint some more?
	if bal < 10000 {
		// load the wallet to a decent starting amount of winston
		res, err = http.Get(
			fmt.Sprintf("http://localhost:1984/mint"+
				"/27xHJ0MNsBUKFIdOiQ3OlrZdDzSNfBPGnp6YVmWKKxU"+
				"/%d", 1000000000-bal))
		if err != nil {
			t.Fatal(err)
		}
		if body, err = GetBody(res); err != nil {
			t.Fatal(err)
		}
		t.Log(body)
		// mine it
		if err = Mine(t); err != nil {
			t.Fatal(err)
		}
	}
	// spew.Dump(res)
	// t.Log(res)
	// b := make([]byte, res.ContentLength)
	// var n int
	// if n, err = res.Body.Read(b);err!=nil {
	// 	t.Fatal(err)
	// }
	// t.Log(string(b[:n]))
	// var wallet *goar.Wallet
	// if wallet, err = goar.NewWallet(key, arlocal); err != nil {
	// 	t.Fatal(err)
	// }
	// data := utils.Base64Encode([]byte("this is test data"))
	// spew.Dump(data)
	// var reward int64
	// if reward, err = wallet.Client.GetTransactionPrice(len(data), nil); err != nil {
	// 	t.Fatal(err)
	// }
	// var speedFactor int64
	// tx := &types.Transaction{
	// 	Format: 2,
	// 	ID:     "",
	// 	LastTx: "",
	// 	Owner:  "",
	// 	Tags: utils.TagsEncode([]types.Tag{{Name: "Name-Tag",
	// 		Value: "this is a test tag"}}),
	// 	Target:     "",
	// 	Quantity:   "0",
	// 	Data:       data,
	// 	DataReader: nil,
	// 	DataSize:   fmt.Sprintf("%d", len(data)),
	// 	DataRoot:   "",
	// 	Reward:     fmt.Sprintf("%d", reward*(100+speedFactor)/100),
	// 	Signature:  "",
	// 	Chunks:     nil,
	// }
	// spew.Dump(tx)
	// if *tx, err = wallet.SendTransaction(tx); err != nil {
	// 	t.Fatal(err)
	// }
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
