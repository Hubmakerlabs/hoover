package multi

import (
	"sync"

	"github.com/Hubmakerlabs/replicatr/pkg/nostr/context"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/Hubmakerlabs/hoover/pkg/bluesky"
	"github.com/Hubmakerlabs/hoover/pkg/farcaster"
	"github.com/Hubmakerlabs/hoover/pkg/nostr"
)

// Firehose runs concurrent connections capturing new events appearing on each
// of the social network protocols, and calls a closure that provides access to
// the bundle formatted from a protocol event.
func Firehose(
	c context.T,
	cancel context.F,
	wait *sync.WaitGroup,
	nostrRelays, blueskyEndpoints, farcasterHubs []string,
	fn func(bundle *types.BundleItem) (err error),
	resolverPath string,
) {
	go bluesky.Firehose(c, cancel, wait, blueskyEndpoints, fn, resolverPath)
	go nostr.Firehose(c, cancel, wait, nostrRelays, fn)
	go farcaster.Firehose(c, cancel, wait, farcasterHubs, fn)
	select {
	case <-c.Done():
		return
	}
}
