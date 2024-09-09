package bluesky

import (
	"context"
	"fmt"
	"testing"

	arweave "github.com/Hubmakerlabs/hoover/pkg/arweave"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/Hubmakerlabs/replicatr/pkg/interrupt"
)

func TestFirehose(t *testing.T) {
	c, cancel := context.WithCancel(context.Background())
	interrupt.AddHandler(cancel)
	Firehose(c, cancel, func(bundle *types.BundleItem) (err error) {
		fmt.Println()
		arweave.PrintBundleItem(bundle)
		return
	})
}
