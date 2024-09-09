package bluesky

import (
	"fmt"

	arweave "github.com/Hubmakerlabs/hoover/pkg/arweave"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/context"
	"github.com/bluesky-social/indigo/events"
	"github.com/bluesky-social/indigo/events/schedulers/sequential"
	"github.com/gorilla/websocket"
)

const SubReposURL = "wss://bsky.network/xrpc/com.atproto.sync.subscribeRepos"

func Firehose(c context.T, cancel context.F, fn func(bundle *types.BundleItem) (err error)) {
	var conn *websocket.Conn
	var err error
	if conn, err = Connect(c); chk.E(err) {
		return
	}
	rscb := &events.RepoStreamCallbacks{
		RepoCommit: RepoCommit(c, cancel, func(bundle *types.BundleItem) (err error) {
			fmt.Println()
			arweave.PrintBundleItem(bundle)
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
