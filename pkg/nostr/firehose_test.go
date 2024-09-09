package nostr

import (
	"context"
	"fmt"
	"sync"
	"testing"

	ao "github.com/Hubmakerlabs/hoover/pkg/arweave"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/Hubmakerlabs/replicatr/pkg/interrupt"
)

func TestFirehose(t *testing.T) {
	c, cancel := context.WithCancel(context.Background())
	interrupt.AddHandler(cancel)
	var wg sync.WaitGroup
	Firehose(c, cancel, &wg, Relays,
		func(bundle *types.BundleItem) (err error) {
			ao.PrintBundleItem(bundle)
			fmt.Println()
			return
		})
}
