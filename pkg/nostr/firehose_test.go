package nostr

import (
	"context"
	"fmt"
	"testing"

	ao "github.com/Hubmakerlabs/hoover/pkg/arweave"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/Hubmakerlabs/replicatr/pkg/interrupt"
)

var relays = []string{
	"wss://purplepag.es",
	"wss://njump.me",
	"wss://relay.snort.social",
	"wss://relay.damus.io",
	"wss://relay.primal.net",
}

func TestFirehose(t *testing.T) {
	c, cancel := context.WithCancel(context.Background())
	interrupt.AddHandler(cancel)
	Firehose(c, cancel, relays, func(bundle *types.BundleItem) (err error) {
		fmt.Println()
		ao.PrintBundleItem(bundle)
		return
	})
}
