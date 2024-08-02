package bluesky

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/bluesky-social/indigo/events"
	"github.com/bluesky-social/indigo/events/schedulers/sequential"
	"github.com/bluesky-social/indigo/repo"
	"github.com/gorilla/websocket"
)

const subReposURL = "wss://bsky.network/xrpc/com.atproto.sync.subscribeRepos"

func TestFirehose(t *testing.T) {
	var conn *websocket.Conn
	var err error
	if conn, _, err = websocket.DefaultDialer.Dial(subReposURL,
		http.Header{}); chk.E(err) {

		os.Exit(1)
	}
	c, cancel := context.WithCancel(context.Background())
	const limit = 24
	// now := time.Now()
	knownEvents := map[string]struct{}{}
	rsc := &events.RepoStreamCallbacks{
		RepoCommit: func(evt *atproto.SyncSubscribeRepos_Commit) (err error) {

			// var did syntax.DID
			// if did, err = syntax.ParseDID(evt.Repo); chk.E(err) {
			// 	return nil
			// }
			// log.I.F("\nEvent from %s\n", did)
			// var b B
			// b, err = evt.Blocks.MarshalJSON()
			// log.I.F("%s", b)
			// log.I.S(evt)
			var rr *repo.Repo
			if rr, err = repo.ReadRepoFromCar(c, bytes.NewReader(evt.Blocks)); chk.E(err) {
				return nil
			}
			_ = rr
			// log.I.S(evt.Ops)
			for _, op := range evt.Ops {
				var collection syntax.NSID
				var rkey syntax.RecordKey
				if collection, rkey, err = splitRepoPath(op.Path); chk.E(err) {
					return
				}
				_, _ = collection, rkey
				if _, ok := knownEvents[collection.String()]; !ok {
					log.I.Ln(collection.String())
				}
				knownEvents[collection.String()] = struct{}{}
				// log.I.S(knownEvents)
			}
			if len(knownEvents) > limit {
				cancel()
			}
			// log.I.S(rr)
			// for i := range evt.Ops {
			// 	log.I.S(evt.Ops[i])
			// }
			return
		},
	}
	sched := sequential.NewScheduler("myfirehose", rsc.EventHandler)
	err = events.HandleRepoStream(c, conn, sched)
	chk.E(err)
}

// TODO: move this to a "ParsePath" helper in syntax package?
func splitRepoPath(path string) (syntax.NSID, syntax.RecordKey, error) {
	parts := strings.SplitN(path, "/", 3)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid record path: %s", path)
	}
	collection, err := syntax.ParseNSID(parts[0])
	if err != nil {
		return "", "", err
	}
	rkey, err := syntax.ParseRecordKey(parts[1])
	if err != nil {
		return "", "", err
	}
	return collection, rkey, nil
}
