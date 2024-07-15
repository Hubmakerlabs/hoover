package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	cacheSchema "github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/cache/schema"
	paySchema "github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/pay/schema"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/sdk/schema"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/token"
	tokSchema "github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/token/schema"
	"github.com/tidwall/sjson"
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

func (c *Client) GetInfo() (info schema.Info, err error) {
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

	info = schema.Info{}
	err = json.Unmarshal(res.Bytes(), &info)
	return
}

func (c *Client) GetTokens() (tokens map[string]*token.Token, err error) {
	tokens = map[string]*token.Token{}
	info, err := c.GetInfo()
	if err != nil {
		return
	}
	for _, tokenInfo := range info.TokenList {
		targetChains := make([]tokSchema.TargetChain, 0)
		for _, v := range tokenInfo.CrossChainInfoList {
			targetChains = append(targetChains, tokSchema.TargetChain{
				ChainId:   v.ChainId,
				ChainType: v.ChainType,
				Decimals:  v.Decimals,
				TokenID:   v.TokenID,
			})
		}

		tok := token.New(
			tokenInfo.ID, tokenInfo.Symbol, tokenInfo.ChainType,
			tokenInfo.ChainID,
			tokenInfo.Decimals, targetChains,
		)
		tokens[tok.Tag()] = tok
	}
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
	result := schema.LimitIp{}
	if err = json.Unmarshal(res.Bytes(), &result); err != nil {
		return
	}
	return result.Limit, nil
}

func (c *Client) Balance(tokenTag, accid string) (balance schema.AccBalance,
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

	balance = schema.AccBalance{}
	err = json.Unmarshal(res.Bytes(), &balance)
	return
}

func (c *Client) Balances(accid string) (balances schema.AccBalances,
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

	balances = schema.AccBalances{}
	err = json.Unmarshal(res.Bytes(), &balances)
	return
}

func (c *Client) BlackList(tokenTag string) ([]string, error) {
	req := c.cli.Request()
	req.Path(fmt.Sprintf("/black_list/%s", tokenTag))

	res, err := req.Send()
	if err != nil {
		return nil, err
	}
	defer res.Close()
	if !res.Ok {
		err = decodeRespErr(res.Bytes())
		return nil, err
	}

	result := make([]string, 0)
	err = json.Unmarshal(res.Bytes(), &result)
	return result, err
}

func (c *Client) WhiteList(tokenTag string) ([]string, error) {
	req := c.cli.Request()
	req.Path(fmt.Sprintf("/white_list/%s", tokenTag))

	res, err := req.Send()
	if err != nil {
		return nil, err
	}
	defer res.Close()
	if !res.Ok {
		err = decodeRespErr(res.Bytes())
		return nil, err
	}

	result := make([]string, 0)
	err = json.Unmarshal(res.Bytes(), &result)
	return result, err
}

func (c *Client) AccInfo(accid string) (resp schema.RespAcc, err error) {
	req := c.cli.Request()
	req.Path(fmt.Sprintf("/account/%s", accid))
	res, err := req.Send()
	if err != nil {
		return
	}
	defer res.Close()
	if !res.Ok {
		err = decodeRespErr(res.Bytes())
		return
	}
	err = json.Unmarshal(res.Bytes(), &resp)
	return
}

func (c *Client) Txs(startCursor int64, orderBy string, limit int,
	opts schema.TxOpts) (txs schema.Txs, err error) {
	req := c.cli.Request()
	req.Path("/txs")
	if startCursor > 0 {
		req.AddQuery("cursor", fmt.Sprintf("%d", startCursor))
	}
	if len(orderBy) > 0 {
		req.AddQuery("order", orderBy)
	}
	if limit > 0 {
		req.AddQuery("count", fmt.Sprintf("%d", limit))
	}

	if len(opts.Address) > 0 {
		req.AddQuery("address", opts.Address)
	}
	if len(opts.TokenTag) > 0 {
		req.AddQuery("tokenTag", opts.TokenTag)
	}
	if len(opts.Action) > 0 {
		req.AddQuery("action", opts.Action)
	}
	if len(opts.WithoutAction) > 0 {
		req.AddQuery("withoutAction", opts.WithoutAction)
	}

	res, err := req.Send()
	if err != nil {
		return
	}
	defer res.Close()
	if !res.Ok {
		err = decodeRespErr(res.Bytes())
		return
	}

	txs = schema.Txs{}
	err = json.Unmarshal(res.Bytes(), &txs)
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

func (c *Client) TxByHash(everHash string) (tx schema.Tx, err error) {
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

	tx = schema.Tx{}
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
func (c *Client) MintTx(chainHash string) (tx schema.Tx, err error) {
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

	tx = schema.Tx{}
	err = json.Unmarshal(res.Bytes(), &tx)
	return
}

// PendingTxs get pending Txs
// everHash: means get from the everTx
func (c *Client) PendingTxs(everHash string) (txs schema.PendingTxs,
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
	txs = schema.PendingTxs{}
	err = res.JSON(&txs)
	return
}

func (c *Client) Fee(tokenTag string) (fee schema.Fee, err error) {
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

	fee = schema.Fee{}
	err = json.Unmarshal(res.Bytes(), &fee)
	return
}

func (c *Client) Fees() (fees schema.Fees, err error) {
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

	fees = schema.Fees{}
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
	respStatus := schema.RespStatus{}
	if err = json.Unmarshal(res.Bytes(), &respStatus); err != nil {
		return
	}
	if respStatus.Status != "ok" {
		err = decodeRespErr(res.Bytes())
	}

	return
}

func (c *Client) Mint102WithoutSig(tokenTag, from, to, amount string) (everTx paySchema.Transaction,
	err error) {
	return c.AssembleTxWithoutSig(tokenTag, from, to, amount, "0",
		tokSchema.TxActionMint, "")
}

func (c *Client) TransferWithoutSig(tokenTag, from, to, amount string) (everTx paySchema.Transaction,
	err error) {
	return c.AssembleTxWithoutSig(tokenTag, from, to, amount, "0",
		tokSchema.TxActionTransfer, "")
}

func (c *Client) AddWhiteListWithoutSig(tokenTag, from string,
	whiteList []string) (everTx paySchema.Transaction, err error) {
	data, err := sjson.Set("", "whiteList", whiteList)
	if err != nil {
		return
	}
	return c.AssembleTxWithoutSig(tokenTag, from, from, "0", "0",
		tokSchema.TxActionAddWhiteList, data)
}

func (c *Client) AddBlackListWithoutSig(tokenTag, from string,
	blackList []string) (everTx paySchema.Transaction, err error) {
	data, err := sjson.Set("", "blackList", blackList)
	if err != nil {
		return
	}
	return c.AssembleTxWithoutSig(tokenTag, from, from, "0", "0",
		tokSchema.TxActionAddBlackList, data)
}

func (c *Client) Burn102WithoutSig(tokenTag, from, amount string) (everTx paySchema.Transaction,
	err error) {
	data, err := sjson.Set("", "targetChainType", tokSchema.ChainTypeEverpay)
	if err != nil {
		return
	}
	return c.AssembleTxWithoutSig(tokenTag, from, tokSchema.ZeroAddress, amount,
		"0", tokSchema.TxActionBurn, data)
}

func (c *Client) AssembleTxWithoutSig(tokenTag, from, to, amount, fee, action, data string) (everTx paySchema.Transaction,
	err error) {
	info, err := c.GetInfo()
	if err != nil {
		return
	}
	tokenInfo := schema.TokenInfo{}
	for _, t := range info.TokenList {
		if t.Tag == tokenTag {
			tokenInfo = t
		}
		break
	}

	// assemble tx
	everTx = paySchema.Transaction{
		TokenSymbol:  tokenInfo.Symbol,
		Action:       action,
		From:         from,
		To:           to,
		Amount:       amount,
		Fee:          fee,
		FeeRecipient: info.FeeRecipient,
		Nonce:        fmt.Sprintf("%d", time.Now().UnixNano()/1000000),
		TokenID:      tokenInfo.ID,
		ChainType:    tokenInfo.ChainType,
		ChainID:      tokenInfo.ChainID,
		Data:         data,
		Version:      tokSchema.TxVersionV1,
		Sig:          "",
	}
	return
}

func decodeRespErr(errMsg []byte) error {
	resErr := schema.RespErr{}
	if err := json.Unmarshal(errMsg, &resErr); err != nil {
		return errors.New(string(errMsg))
	}
	return resErr
}
