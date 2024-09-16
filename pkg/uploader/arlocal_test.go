package uploader

import (
	_ "embed"
	"fmt"
	"net/http"
	"testing"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar"
)

//go:embed keyfile.json
var key []byte

// this requires the use of https://github.com/textury/arlocal todo: maybe we can spawn it in a container?

const arlocal = "http://localhost:1984"

// Mine triggers arlocal to commit the pending transactions into a block
func Mine(endpoint string, t *testing.T) (err error) {
	_, err = http.Get(
		fmt.Sprintf("%s/mine", endpoint))
	if err != nil {
		return
	}
	return
}

func Mint(endpoint, account string, amount int64, t *testing.T) (err error) {
	var res *http.Response
	var body string
	res, err = http.Get(fmt.Sprintf("%s/mint/%s/%d",
		endpoint, account, amount))
	if err != nil {
		t.Fatal(err)
	}
	if body, err = GetBody(res); err != nil {
		t.Fatal(err)
	}
	t.Log(body)
	return
}

// BumpBalance checks that the balance is not low and adds winstons to it to bring it to a
// decent amount.
func BumpBalance(endpoint, address string, t *testing.T) {
	var err error
	var bal int64
	if bal, err = GetBalance(endpoint, address); err != nil {
		t.Fatal(err)
	}
	// do we need to mint some more?
	if bal < 1000000 {
		// load the wallet to a decent starting amount of winston
		if err = Mint(endpoint, address, 1000000000-bal, t); err != nil {
			t.Fatal(err)
		}
		// mine it
		if err = Mine(endpoint, t); err != nil {
			t.Fatal(err)
		}
	}
}

func GetTestWallet(endpoint string, t *testing.T) (address string, wallet *goar.Wallet, err error) {
	if wallet, err = goar.NewWallet(key, endpoint); err != nil {
		t.Fatal(err)
	}
	address = wallet.Signer.Address
	return
}
