package multi

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	arweave "github.com/Hubmakerlabs/hoover/pkg/arweave"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/Hubmakerlabs/hoover/pkg/bluesky"
	"github.com/Hubmakerlabs/hoover/pkg/farcaster"
	"github.com/Hubmakerlabs/hoover/pkg/nostr"
	"github.com/Hubmakerlabs/replicatr/pkg/interrupt"
)

func TestFirehose(t *testing.T) {
	c, cancel := context.WithCancel(context.Background())
	interrupt.AddHandler(cancel)
	go func() {
		time.Sleep(time.Second * 10)
		cancel()
	}()
	var wg sync.WaitGroup
	fmt.Println()
	Firehose(c, cancel, &wg, nostr.Relays, bluesky.Urls, farcaster.Urls,
		func(bundle *types.BundleItem) (err error) {
			arweave.PrintBundleItem(bundle)
			return
		})
}
