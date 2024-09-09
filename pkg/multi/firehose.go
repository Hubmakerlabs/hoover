package multi

import (
	"sync"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/Hubmakerlabs/hoover/pkg/bluesky"
	"github.com/Hubmakerlabs/hoover/pkg/nostr"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/context"
)

func Firehose(c context.T, cancel context.F, wait *sync.WaitGroup, nostrRelays []string,
	fn func(bundle *types.BundleItem) (err error)) {
	go bluesky.Firehose(c, cancel, wait, fn)
	go nostr.Firehose(c, cancel, wait, nostrRelays, fn)
	select {
	case <-c.Done():
		return
	}
}
