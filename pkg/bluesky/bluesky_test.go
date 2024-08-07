package bluesky

import (
	"context"
	"os/signal"
	"syscall"
	"testing"

	"github.com/bluesky-social/indigo/events"
	"github.com/bluesky-social/indigo/events/schedulers/sequential"
	"github.com/gorilla/websocket"
)

const subReposURL = "wss://bsky.network/xrpc/com.atproto.sync.subscribeRepos"

func TestFirehose(t *testing.T) {
	var err error
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()
	var conn *websocket.Conn
	if conn, err = Connect(ctx); chk.E(err) {
		t.Fatal(err)
	}
	rscb := &events.RepoStreamCallbacks{
		RepoCommit:    RepoCommit(ctx, cancel),
		RepoHandle:    RepoHandle(),
		RepoInfo:      RepoInfo(),
		RepoMigrate:   RepoMigrate(),
		RepoTombstone: RepoTombstone(),
		LabelLabels:   LabelLabels(),
		LabelInfo:     LabelInfo(),
	}
	seqScheduler := sequential.NewScheduler(conn.RemoteAddr().String(), rscb.EventHandler)
	if err = events.HandleRepoStream(ctx, conn, seqScheduler); chk.E(err) {
		t.Fatal(err)
	}
}
