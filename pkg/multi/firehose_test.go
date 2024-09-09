package multi

import (
	"context"
	"fmt"
	"sync"
	"testing"

	arweave "github.com/Hubmakerlabs/hoover/pkg/arweave"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/Hubmakerlabs/hoover/pkg/nostr"
	"github.com/Hubmakerlabs/replicatr/pkg/interrupt"
)

func TestFirehose(t *testing.T) {
	c, cancel := context.WithCancel(context.Background())
	interrupt.AddHandler(cancel)
	var wg sync.WaitGroup
	fmt.Println()
	Firehose(c, cancel, &wg, nostr.Relays,func(bundle *types.BundleItem) (err error) {
		arweave.PrintBundleItem(bundle)
		return
	})
}
