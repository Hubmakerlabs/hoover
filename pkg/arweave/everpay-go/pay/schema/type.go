package schema

import (
	"github.com/Hubmakerlabs/hoover/pkg/arweave/everpay-go/token/utils"
	"github.com/Hubmakerlabs/hoover/pkg/eth/accounts"
	"github.com/Hubmakerlabs/hoover/pkg/eth/common/hexutil"
)

type Transaction struct {
	TokenSymbol  string `json:"tokenSymbol"`
	Action       string `json:"action"`
	From         string `json:"from"`
	To           string `json:"to"`
	Amount       string `json:"amount"`
	Fee          string `json:"fee"`
	FeeRecipient string `json:"feeRecipient"`
	Nonce        string `json:"nonce"`
	TokenID      string `json:"tokenID"`
	ChainType    string `json:"chainType"`
	ChainID      string `json:"chainID"`
	Data         string `json:"data"`
	Version      string `json:"version"`
	Sig          string `json:"sig"`

	ArOwner     string `json:"-"`
	ArTxID      string `json:"-"`
	ArTimestamp int64  `json:"-"`
}

func (t *Transaction) String() string {
	return "tokenSymbol:" + t.TokenSymbol + "\n" +
		"action:" + t.Action + "\n" +
		"from:" + t.From + "\n" +
		"to:" + t.To + "\n" +
		"amount:" + t.Amount + "\n" +
		"fee:" + t.Fee + "\n" +
		"feeRecipient:" + t.FeeRecipient + "\n" +
		"nonce:" + t.Nonce + "\n" +
		"tokenID:" + t.TokenID + "\n" +
		"chainType:" + t.ChainType + "\n" +
		"chainID:" + t.ChainID + "\n" +
		"data:" + t.Data + "\n" +
		"version:" + t.Version
}

// Tag is the unique identifier of token
func (t *Transaction) Tag() string {
	return utils.Tag(t.ChainType, t.TokenSymbol, t.TokenID)
}

func (t *Transaction) Hash() []byte {
	return accounts.TextHash([]byte(t.String()))
}

func (t *Transaction) HexHash() string {
	return hexutil.Encode(t.Hash())
}
