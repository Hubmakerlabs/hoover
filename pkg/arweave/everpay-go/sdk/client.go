package sdk

import (
	"encoding/json"
	"errors"
	"fmt"

	cacheSchema "github.com/Hubmakerlabs/hoover/pkg/arweave/everpay-go/cache/schema"
	paySchema "github.com/Hubmakerlabs/hoover/pkg/arweave/everpay-go/pay/schema"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/everpay-go/sdk/schema"
	serverSchema "github.com/Hubmakerlabs/hoover/pkg/arweave/everpay-go/server/schema"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/body"
)

type Client struct {
	cli *gentleman.Client
}

func NewClient(payURL string) *Client {
	return &Client{
		cli: gentleman.New().URL(payURL),
	}
}

func (c *Client) SetHeader(key, val string) {
	c.cli.SetHeader(key, val)
}

func (c *Client) GetInfo() (info serverSchema.Info, err error) {
	req := c.cli.Request()
	req.Path("/info")

	res, err := req.Send()
	if err != nil {
		return
	}
	defer res.Close()
	if !res.Ok {
		err = decodeRespErr(res.Bytes())
		return
	}

	info = serverSchema.Info{}
	err = json.Unmarshal(res.Bytes(), &info)
	return
}

func (c *Client) LimitIp() (isLimit bool, err error) {
	req := c.cli.Request()
	req.Path("/limit_ip")
	res, err := req.Send()
	if err != nil {
		return
	}
	defer res.Close()
	if !res.Ok {
		err = decodeRespErr(res.Bytes())
		return
	}
	result := serverSchema.LimitIp{}
	if err = json.Unmarshal(res.Bytes(), &result); err != nil {
		return
	}
	return result.Limit, nil
}

func (c *Client) Balance(tokenTag, accid string) (balance serverSchema.AccBalance,
	err error) {
	req := c.cli.Request()
	req.Path(fmt.Sprintf("/balance/%s/%s", tokenTag, accid))

	res, err := req.Send()
	if err != nil {
		return
	}
	defer res.Close()
	if !res.Ok {
		err = decodeRespErr(res.Bytes())
		return
	}

	balance = serverSchema.AccBalance{}
	err = json.Unmarshal(res.Bytes(), &balance)
	return
}

func (c *Client) Balances(accid string) (balances serverSchema.AccBalances,
	err error) {
	req := c.cli.Request()
	req.Path(fmt.Sprintf("/balances/%s", accid))

	res, err := req.Send()
	if err != nil {
		return
	}
	defer res.Close()
	if !res.Ok {
		err = decodeRespErr(res.Bytes())
		return
	}

	balances = serverSchema.AccBalances{}
	err = json.Unmarshal(res.Bytes(), &balances)
	return
}

// Txs
// option args: tokenId, action, withoutAction
// default value: page(1), orderBy(desc)
func (c *Client) Txs(page int, orderBy, tokenSymbol string,
	action, withoutAction string) (txs serverSchema.Txs, err error) {
	req := c.cli.Request()
	req.Path("/txs")
	req.AddQuery("page", fmt.Sprintf("%v", page))
	req.AddQuery("order", orderBy)
	req.AddQuery("symbol", tokenSymbol)
	req.AddQuery("action", action)
	req.AddQuery("withoutAction", withoutAction)

	res, err := req.Send()
	if err != nil {
		return
	}
	defer res.Close()
	if !res.Ok {
		err = decodeRespErr(res.Bytes())
		return
	}

	txs = serverSchema.Txs{}
	err = json.Unmarshal(res.Bytes(), &txs)
	return
}

// TxsByAcc
// option args: tokenId, action, withoutAction
// default value: page(1), orderBy(desc)
func (c *Client) TxsByAcc(accid string, page int, orderBy string,
	tokenSymbol, action, withoutAction string) (txs serverSchema.AccTxs,
	err error) {
	req := c.cli.Request()
	req.Path(fmt.Sprintf("/txs/%s", accid))
	req.AddQuery("page", fmt.Sprintf("%v", page))
	req.AddQuery("order", orderBy)
	req.AddQuery("symbol", tokenSymbol)
	req.AddQuery("action", action)
	req.AddQuery("withoutAction", withoutAction)

	res, err := req.Send()
	if err != nil {
		return
	}
	defer res.Close()
	if !res.Ok {
		err = decodeRespErr(res.Bytes())
		return
	}

	txs = serverSchema.AccTxs{}
	err = json.Unmarshal(res.Bytes(), &txs)
	return
}

func (c *Client) CursorTxs(startCursor uint64,
	tokenSymbol, action, withoutAction string) (txs serverSchema.Txs,
	err error) {
	req := c.cli.Request()
	req.Path("/txs")
	req.AddQuery("symbol", tokenSymbol)
	req.AddQuery("action", action)
	req.AddQuery("withoutAction", withoutAction)
	if startCursor == 0 {
		startCursor = 1
	}
	req.AddQuery("cursor", fmt.Sprintf("%d", startCursor))

	res, err := req.Send()
	if err != nil {
		return
	}
	defer res.Close()
	if !res.Ok {
		err = decodeRespErr(res.Bytes())
		return
	}

	txs = serverSchema.Txs{}
	err = json.Unmarshal(res.Bytes(), &txs)
	return
}

func (c *Client) CursorTxsByAcc(accid string, startCursor uint64,
	tokenSymbol, action, withoutAction string) (txs serverSchema.Txs,
	err error) {
	req := c.cli.Request()
	req.Path(fmt.Sprintf("/txs/%s", accid))
	req.AddQuery("symbol", tokenSymbol)
	req.AddQuery("action", action)
	req.AddQuery("withoutAction", withoutAction)
	if startCursor == 0 {
		startCursor = 1
	}
	req.AddQuery("cursor", fmt.Sprintf("%d", startCursor))

	res, err := req.Send()
	if err != nil {
		return
	}
	defer res.Close()
	if !res.Ok {
		err = decodeRespErr(res.Bytes())
		return
	}

	accTxs := serverSchema.AccTxs{}
	err = json.Unmarshal(res.Bytes(), &accTxs)
	txs = accTxs.Txs
	return
}

// SubscribeTxs
// fq.StartCursor: option
// fq.Address: option
// fq.TokenSymbol: option
// fq.Action: option
// fq.WithoutAction: option
func (c *Client) SubscribeTxs(fq schema.FilterQuery) *SubscribeTx {
	sub := newSubscribeTx(c, fq)
	go sub.run()
	return sub
}

func (c *Client) TxByHash(everHash string) (tx serverSchema.Tx, err error) {
	req := c.cli.Request()
	req.Path(fmt.Sprintf("/tx/%s", everHash))

	res, err := req.Send()
	if err != nil {
		return
	}
	defer res.Close()
	if !res.Ok {
		err = decodeRespErr(res.Bytes())
		return
	}

	tx = serverSchema.Tx{}
	err = json.Unmarshal(res.Bytes(), &tx)
	return
}

func (c *Client) BundleByHash(everHash string) (
	tx cacheSchema.TxResponse,
	bundle paySchema.BundleWithSigs,
	internalStatus cacheSchema.InternalStatus,
	err error) {

	txRes, err := c.TxByHash(everHash)
	if err != nil {
		return
	}
	tx = *txRes.Tx

	if tx.Action != paySchema.TxActionBundle {
		err = ErrNotBundleTx
		return
	}

	bundleData := paySchema.BundleData{}
	if err = json.Unmarshal([]byte(tx.Data), &bundleData); err != nil {
		return
	}
	bundle = bundleData.Bundle

	err = json.Unmarshal([]byte(tx.InternalStatus), &internalStatus)
	return
}

// MintTx get minted everTx by onChain mint txHash
func (c *Client) MintTx(chainHash string) (tx serverSchema.Tx, err error) {
	req := c.cli.Request()
	req.Path(fmt.Sprintf("/minted/%s", chainHash))

	res, err := req.Send()
	if err != nil {
		return
	}
	defer res.Close()
	if !res.Ok {
		err = decodeRespErr(res.Bytes())
		return
	}

	tx = serverSchema.Tx{}
	err = json.Unmarshal(res.Bytes(), &tx)
	return
}

// PendingTxs get pending Txs
// everHash: means get from the everTx
func (c *Client) PendingTxs(everHash string) (txs serverSchema.PendingTxs,
	err error) {
	req := c.cli.Request()
	req.Path("/tx/pending")
	req.AddQuery("everHash", everHash)
	res, err := req.Send()
	if err != nil {
		return
	}

	defer res.Close()
	if !res.Ok {
		err = decodeRespErr(res.Bytes())
		return
	}
	txs = serverSchema.PendingTxs{}
	err = res.JSON(&txs)
	return
}

func (c *Client) Fee(tokenTag string) (fee serverSchema.Fee, err error) {
	req := c.cli.Request()
	req.Path(fmt.Sprintf("/fee/%s", tokenTag))

	res, err := req.Send()
	if err != nil {
		return
	}
	defer res.Close()
	if !res.Ok {
		err = decodeRespErr(res.Bytes())
		return
	}

	fee = serverSchema.Fee{}
	err = json.Unmarshal(res.Bytes(), &fee)
	return
}

func (c *Client) Fees() (fees serverSchema.Fees, err error) {
	req := c.cli.Request()
	req.Path("/fees")
	res, err := req.Send()
	if err != nil {
		return
	}
	defer res.Close()
	if !res.Ok {
		err = decodeRespErr(res.Bytes())
		return
	}

	fees = serverSchema.Fees{}
	err = json.Unmarshal(res.Bytes(), &fees)
	return
}

func (c *Client) SubmitTx(tx paySchema.Transaction) (err error) {
	req := c.cli.Request()
	req.Path(fmt.Sprintf("/tx"))
	req.Method("POST")
	req.Use(body.JSON(tx))

	res, err := req.Send()
	if err != nil {
		return
	}
	defer res.Close()
	if !res.Ok {
		err = decodeRespErr(res.Bytes())
		return
	}

	// check status is "ok"
	respStatus := serverSchema.RespStatus{}
	if err = json.Unmarshal(res.Bytes(), &respStatus); err != nil {
		return
	}
	if respStatus.Status != "ok" {
		err = decodeRespErr(res.Bytes())
	}

	return
}

func decodeRespErr(errMsg []byte) error {
	resErr := serverSchema.RespErr{}
	if err := json.Unmarshal(errMsg, &resErr); err != nil {
		return errors.New(string(errMsg))
	}
	return resErr
}
