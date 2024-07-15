package schema

import "math/big"

type Transaction struct {
	Action          string
	From            string // notice: Case Sensitive !!!
	To              string // notice: Case Sensitive !!!
	Amount          *big.Int
	Fee             *big.Int
	FeeRecipient    string // notice: Case Sensitive !!!
	Data            string
	TargetChainType string
}

type TargetChain struct {
	ChainId   string  `json:"targetChainId"`
	ChainType string  `json:"targetChainType"` // e.g: "avalanche" "arweave" "ethereum","moon"
	Decimals  int     `json:"targetDecimals"`  // e.g: 18
	TokenID   string  `json:"targetTokenId"`   // target chain token address
	Locker    *Locker `json:"-"`
	// for oracle verify
	Rpc   string `json:"-"` // chain node url
	PstGw string `json:"-"` // pst token verify gateway
}
