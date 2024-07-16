package listener

import (
	"github.com/Hubmakerlabs/hoover/pkg/arweave/utils/arweave"
)

type Payload struct {
	BlockHash      arweave.Base64String
	BlockHeight    int64
	BlockTimestamp int64
	Transactions   []*arweave.Transaction
}
