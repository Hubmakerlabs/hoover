package farcaster

import (
	"context"
	"sync"
	"testing"
	"time"

	arweave "github.com/Hubmakerlabs/hoover/pkg/arweave"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/Hubmakerlabs/replicatr/pkg/interrupt"
)

func TestFirehose(t *testing.T) {
	c, cancel := context.WithCancel(context.Background())
	interrupt.AddHandler(func() {
		cancel()
	})
	go func() {
		select {
		case <-c.Done():
		case <-time.After(time.Second * 10):
		}
		cancel()
	}()
	var wg sync.WaitGroup
	Firehose(c, cancel, &wg, func(bundle *types.BundleItem) (err error) {
		arweave.PrintBundleItem(bundle)
		return
	})

}
