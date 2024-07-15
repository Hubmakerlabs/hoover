package schema

import (
	cacheSchema "github.com/Hubmakerlabs/hoover/pkg/arweave/everpay-go/cache/schema"
	confSchema "github.com/Hubmakerlabs/hoover/pkg/arweave/everpay-go/config/schema"
	tokSchema "github.com/Hubmakerlabs/hoover/pkg/arweave/everpay-go/token/schema"
)

type RespErr struct {
	Err string `json:"error"`
}

func (r RespErr) Error() string {
	return r.Err
}

type WithdrawTxResponse struct {
	EverHash    string
	Token       string
	Status      string
	WithdrawFee string
	WithdrawTx  string
	RefundTx    string
	Error       string
}

type TokenInfo struct {
	Tag                string                           `json:"tag"`
	ID                 string                           `json:"id"`
	Symbol             string                           `json:"symbol"`
	Decimals           int                              `json:"decimals"`
	TotalSupply        string                           `json:"totalSupply"`
	ChainType          string                           `json:"chainType"`
	ChainID            string                           `json:"chainID"`
	BurnFees           map[string]string                `json:"burnFees"` // key: targetChainType, val: fee
	TransferFee        string                           `json:"transferFee"`
	BundleFee          string                           `json:"bundleFee"`
	HolderNum          int                              `json:"holderNum"`
	CrossChainInfoList map[string]tokSchema.TargetChain `json:"crossChainInfoList"` // key: targetChainType
}

type Info struct {
	IsSynced        bool              `json:"isSynced"`
	IsClosed        bool              `json:"isClosed"`
	BalanceRootHash string            `json:"balanceRootHash"`
	RootHash        string            `json:"rootHash"`
	EverRootHash    string            `json:"everRootHash"`
	Owner           string            `json:"owner"`
	EthChainID      string            `json:"ethChainID"`
	FeeRecipient    string            `json:"feeRecipient"`
	EthLocker       string            `json:"ethLocker"`
	ArLocker        string            `json:"arLocker"`
	Lockers         map[string]string `json:"lockers"`
	TokenList       []TokenInfo       `json:"tokenList"`
}

type TuringInfo struct {
	IsSynced                    bool   `json:"isSynced"`
	RollupWatcherArId           string `json:"rollupWatcherArId"`
	RollupLastArId              string `json:"rollupLastArId"`
	RollupLastOnChainEverTxHash string `json:"rollupLastOnChainEverTxHash"`
	TrackerLastArId             string `json:"trackerLastArId"`
	PendingTxNum                int    `json:"pendingTxNum"`
	CurRollupTxNum              int    `json:"curRollupTxNum"`
	RollupAddr                  string `json:"rollupAddr"`
	TrackerAddr                 string `json:"trackerAddr"`
}

type LimitIp struct {
	Limit bool `json:"limit"`
}

type Balance struct {
	Tag      string `json:"tag"`
	Amount   string `json:"amount"`
	Decimals int    `json:"decimals"`
}

type AccBalance struct {
	AccId   string  `json:"accid"`
	Balance Balance `json:"balance"`
}

type AccBalances struct {
	AccId    string    `json:"accid"`
	Balances []Balance `json:"balances"`
}

type Txs struct {
	Txs         []*cacheSchema.TxResponse `json:"txs"`
	CurrentPage int                       `json:"currentPage"`
	TotalPages  int                       `json:"totalPages"`
}

type AccTxs struct {
	AccId string `json:"accid"`
	Txs
}

type Tx struct {
	Tx *cacheSchema.TxResponse `json:"tx"`
}

type RespStatus struct {
	Status string `json:"status"`
}

type PendingTxs struct {
	Total       int                       `json:"total"`
	HasNextPage bool                      `json:"hasNextPage"` // true means can get more
	Txs         []*cacheSchema.TxResponse `json:"txs"`
}

type Fee struct {
	Fee confSchema.TokenFee `json:"fee"`
}

type Fees struct {
	Fees []confSchema.TokenFee `json:"fees"`
}
