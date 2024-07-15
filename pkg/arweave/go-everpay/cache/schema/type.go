package schema

import (
	"encoding/json"

	paySchema "github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/pay/schema"
)

type TxInfo struct {
	EverHash  string
	ChainType string
	Action    string
	Data      string
	Timestamp int64
}

type TxResponse struct {
	RawId          int64  `json:"rawId"`
	ID             string `json:"id"` // AR tx id
	TokenSymbol    string `json:"tokenSymbol"`
	Action         string `json:"action"`
	From           string `json:"from"`
	To             string `json:"to"`
	Amount         string `json:"amount"`
	Fee            string `json:"fee"`
	FeeRecipient   string `json:"feeRecipient"`
	Nonce          int64  `json:"nonce"`
	TokenID        string `json:"tokenID"`
	ChainType      string `json:"chainType"`
	ChainID        string `json:"chainID"`
	Data           string `json:"data"`
	Version        string `json:"version"`
	Sig            string `json:"sig"`
	EverHash       string `json:"everHash"`
	Status         string `json:"status"`
	InternalStatus string `json:"internalStatus"` // if internal tx (bundle tx) execute success, return "success" then return err info
	Timestamp      int64  `json:"timestamp"`      // arTx timestamp

	// for cross chain
	TargetChainTxHash string `json:"targetChainTxHash"`
	// for express
	Express Express `json:"express"`
}

type Express struct {
	ChainTxHash    string `json:"chainTxHash"`
	WithdrawFee    string `json:"withdrawFee"`
	RefundEverHash string `json:"refundEverHash"`
	Err            string `json:"err"`
}

type InternalStatus struct {
	Status string `json:"status"` // "success" or "failed"
	*paySchema.InternalErr
}

func (i InternalStatus) Marshal() string {
	js, _ := json.Marshal(i)
	return string(js)
}
