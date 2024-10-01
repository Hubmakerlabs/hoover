package farcaster

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Hubmakerlabs/hoover/pkg/arweave"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
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
	Firehose(c, cancel, &wg, Urls, func(bundle *types.BundleItem) (err error) {
		ao.PrintBundleItem(bundle)
		return
	})

}
