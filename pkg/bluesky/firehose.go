package bluesky

import (
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/Hubmakerlabs/replicatr/pkg/nostr/context"
	"github.com/bluesky-social/indigo/events"
	"github.com/bluesky-social/indigo/events/schedulers/sequential"
	"github.com/gorilla/websocket"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
)

const SubReposURL = "wss://bsky.network/xrpc/com.atproto.sync.subscribeRepos"

var Urls = []string{
	SubReposURL,
}

func Firehose(c context.T, cancel context.F, wg *sync.WaitGroup, endpoints []string,
	fn func(bundle *types.BundleItem) (err error), resolverPath string) {

	wg.Add(1)
	var ready atomic.Bool
	ready.Store(false)
	// set up resolver service
	resolv, err := NewResolver(filepath.Join(resolverPath, "bluesky"))
	if err != nil {
		// really this should not happen unless there is a misconfiguration or filesystem
		// problem.
		panic(err)
	}
	defer resolv.Close()
	for {
		select {
		case <-c.Done():
			return
		default:
			var conn *websocket.Conn
			for _, endpoint := range endpoints {
				if conn, err = Connect(c, endpoint); chk.E(err) {
					continue
				}
				// if it worked, continue
				break
				// todo: we only actually use the main bluesky inc. one here and never see
				//  problems with it but what if it actually ever did stop being VC funded???
			}
			if err != nil {
				continue
			}
			rscb := &events.RepoStreamCallbacks{
				RepoCommit: RepoCommit(c, cancel, func(bundle *types.BundleItem) (err error) {
					if !ready.Load() {
						ready.Store(true)
						wg.Done()
					}
					if err = fn(bundle); chk.E(err) {
						return nil
					}
					return
				}, resolv),
				RepoHandle:    RepoHandle(),
				RepoInfo:      RepoInfo(),
				RepoMigrate:   RepoMigrate(),
				RepoTombstone: RepoTombstone(),
				LabelLabels:   LabelLabels(),
				LabelInfo:     LabelInfo(),
			}
			seqScheduler := sequential.NewScheduler(conn.RemoteAddr().String(),
				rscb.EventHandler)
			if err = events.HandleRepoStream(c, conn, seqScheduler); chk.E(err) {
				continue
			}
		}
	}
}
