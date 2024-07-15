package schema

import (
	"encoding/json"

	"github.com/Hubmakerlabs/hoover/pkg/eth/accounts"
	"github.com/Hubmakerlabs/hoover/pkg/eth/common/hexutil"
)

const (
	BundleTxVersionV1 = "v1"
	TxActionBundle    = "bundle"
)

type BundleItem struct {
	Tag     string `json:"tag"` // token tag
	ChainID string `json:"chainID"`
	From    string `json:"from"`
	To      string `json:"to"`
	Amount  string `json:"amount"`
}

type Bundle struct {
	Items      []BundleItem `json:"items"`
	Expiration int64        `json:"expiration"` // second
	Salt       string       `json:"salt"`       // uuid
	Version    string       `json:"version"`
}

type BundleWithSigs struct {
	Bundle
	Sigs map[string]string `json:"sigs"` // accid -> sig
}

type BundleData struct {
	Bundle BundleWithSigs `json:"bundle"`
}

func (s *Bundle) String() string {
	by, _ := json.Marshal(s)
	return string(by)
}

func (s *Bundle) Hash() []byte {
	return accounts.TextHash([]byte(s.String()))
}

func (s *Bundle) HashHex() string {
	return hexutil.Encode(s.Hash())
}
