package sdk

import (
	"testing"
	"time"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/sdk/schema"
	"github.com/stretchr/testify/assert"
)

var testClient *Client

func init() {
	testClient = NewClient("https://api-dev.everpay.io")
}

// func TestGetInfo(t *testing.T) {
// 	info, err := testClient.GetInfo()
// 	assert.NoError(t, err)
// 	assert.Equal(t, "42", info.EthChainID)
// }

// func Test_IpWhite(t *testing.T) {
// testClient.SetHeader("origin","https://auction.everpay.io")
//
// for {
// 	go func() {
// 		info, err := testClient.GetInfo()
// 		assert.NoError(t, err)
// 		t.Log(info.EthLocker)
// 	}()
// 	time.Sleep(10 * time.Millisecond)
// }
// }

func TestBalance(t *testing.T) {
	bal, err := testClient.Balance("ethereum-eth-0x0000000000000000000000000000000000000000",
		"0x2ca81e1253f9426c62Df68b39a22A377164eeC92")
	assert.NoError(t, err)
	assert.Equal(t, "0x2ca81e1253f9426c62Df68b39a22A377164eeC92", bal.AccId)
}

func TestBalances(t *testing.T) {
	bal, err := testClient.Balances("0x2ca81e1253f9426c62Df68b39a22A377164eeC92")
	assert.NoError(t, err)
	assert.Equal(t, "0x2ca81e1253f9426c62Df68b39a22A377164eeC92", bal.AccId)
}

// func TestTxs(t *testing.T) {
// 	txs, err := testClient.Txs(1, "asc", 0, schema.TxOpts{})
// 	assert.NoError(t, err)
// 	assert.Equal(t, "HOKX5WM3bZw4mTXBnVj6llXgNsHcnK-Pfz4s4f5n1LA",
// 		txs.Txs[0].ID)
// }

// func TestTxByHash(t *testing.T) {
// 	tx, err := testClient.TxByHash("0x13acdb2097ba66a0466ba93cb350259cce240df90d8573506284e8e311a6ef1a")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "O3sPMMqDU9duy7xc7P6iVTKpjDAtrtZjXuUY48mzqPs", tx.Tx.ID)
// }

// func TestClient_BundleByHash(t *testing.T) {
// 	tx, bundle, status, err := testClient.BundleByHash("0xb375c4b3b1955d993781d5363693824710a7c65cd7d07dac5e3a53d016ec03af")
// 	assert.NoError(t, err)
// 	assert.Equal(t,
// 		"0xb375c4b3b1955d993781d5363693824710a7c65cd7d07dac5e3a53d016ec03af",
// 		tx.EverHash)
// 	assert.Equal(t,
// 		"0x65093420bd5dbcd4164387a4d8f0781c376bff1c5d41bd151cc91b672e689426",
// 		bundle.HashHex())
// 	assert.Equal(t, status.Index, 0)
// 	assert.Equal(t, status.Msg, "err_insufficient_balance")
// }

// func TestClient_PendingTxs(t *testing.T) {
// 	txs, err := testClient.PendingTxs("")
// 	assert.NoError(t, err)
// 	t.Log(txs)
// }

func TestClient_TokenFee(t *testing.T) {
	tag := "ethereum-usdt-0xd85476c906b5301e8e9eb58d174a6f96b9dfc5ee"
	fee, err := testClient.Fee(tag)
	assert.NoError(t, err)
	t.Log(fee)
}

func TestClient_SubscribeTxs_FilterToken(t *testing.T) {
	testClient = NewClient("https://api-dev.everpay.io")
	accid := "0x4002ED1a1410aF1b4930cF6c479ae373dEbD6223"
	sub := testClient.SubscribeTxs(schema.FilterQuery{
		Address:  accid,
		TokenTag: "arweave,ethereum-ar-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA,0xcc9141efa8c20c7df0778748255b1487957811be",
	})
	go func() {
		// for test
		time.Sleep(10 * time.Second)
		sub.Unsubscribe()
	}()

	for {
		select {
		case tx := <-sub.Subscribe():
			t.Log(tx.RawId, tx.EverHash)
		case <-sub.quit:
			return
		}
	}
}

func TestClient_SubscribeTxs_Cursor(t *testing.T) {
	testClient = NewClient("https://api-dev.everpay.io")
	accid := "0x4002ED1a1410aF1b4930cF6c479ae373dEbD6223"
	sub := testClient.SubscribeTxs(schema.FilterQuery{
		StartCursor: 155457,
		Address:     accid,
	})
	go func() {
		// for test
		time.Sleep(10 * time.Second)
		sub.Unsubscribe()
	}()

	for {
		select {
		case tx := <-sub.Subscribe():
			t.Log(tx.RawId, tx.EverHash)
		case <-sub.quit:
			return
		}
	}
}

func TestGetTokens(t *testing.T) {
	tokens, err := testClient.GetTokens()
	assert.NoError(t, err)
	t.Log(len(tokens), tokens)
}
