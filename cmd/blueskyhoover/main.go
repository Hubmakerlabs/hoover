package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/events"
	"github.com/bluesky-social/indigo/events/schedulers/sequential"
	"github.com/gorilla/websocket"
)

const subReposURL = "wss://bsky.network/xrpc/com.atproto.sync.subscribeRepos"

func main() {
	var conn *websocket.Conn
	var err error
	if conn, _, err = websocket.DefaultDialer.Dial(subReposURL,
		http.Header{}); chk.E(err) {

		os.Exit(1)
	}
	rsc := &events.RepoStreamCallbacks{
		RepoCommit: func(evt *atproto.SyncSubscribeRepos_Commit) (err error) {
			fmt.Println("Event from ", evt.Repo)
			for _, op := range evt.Ops {
				fmt.Printf(" - %s record %s\n", op.Action, op.Path)
			}
			return
		},
	}
	sched := sequential.NewScheduler("myfirehose", rsc.EventHandler)
	events.HandleRepoStream(context.Background(), conn, sched)
}
