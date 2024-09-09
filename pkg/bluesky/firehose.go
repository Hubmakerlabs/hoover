package bluesky

import (
	"sync"
	"sync/atomic"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/context"
	"github.com/bluesky-social/indigo/events"
	"github.com/bluesky-social/indigo/events/schedulers/sequential"
	"github.com/gorilla/websocket"
)

const SubReposURL = "wss://bsky.network/xrpc/com.atproto.sync.subscribeRepos"

func Firehose(c context.T, cancel context.F, wg *sync.WaitGroup,
	fn func(bundle *types.BundleItem) (err error)) {

	wg.Add(1)
	var ready atomic.Bool
	ready.Store(false)
	var conn *websocket.Conn
	var err error
	if conn, err = Connect(c); chk.E(err) {
		return
	}
	rscb := &events.RepoStreamCallbacks{
		RepoCommit: RepoCommit(c, cancel, func(bundle *types.BundleItem) (err error) {
			if !ready.Load() {
				ready.Store(true)
				wg.Done()
			}
			if err = fn(bundle); err != nil {
				return
			}
			return
		}),
		RepoHandle:    RepoHandle(),
		RepoInfo:      RepoInfo(),
		RepoMigrate:   RepoMigrate(),
		RepoTombstone: RepoTombstone(),
		LabelLabels:   LabelLabels(),
		LabelInfo:     LabelInfo(),
	}
	seqScheduler := sequential.NewScheduler(conn.RemoteAddr().String(), rscb.EventHandler)
	if err = events.HandleRepoStream(c, conn, seqScheduler); chk.E(err) {
		return
	}
}
