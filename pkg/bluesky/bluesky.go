// Package bluesky is a set of functions for pulling events from bluesky and
// sending them to Arweave AO.
//
// <insert more notes here>
package bluesky

import (
	"encoding/hex"
	"strconv"
	"time"

	"github.com/Hubmakerlabs/hoover/pkg"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/repo"
	typegen "github.com/whyrusleeping/cbor-gen"
)

// bluesky names are crazy ugly stutter parties

type (
	Repo = *repo.Repo
	Time = time.Time
	Op = *atproto.SyncSubscribeRepos_RepoOp
	Ev = *atproto.SyncSubscribeRepos_Commit
	Rec = typegen.CBORMarshaler
)

var Kinds = map[string]string{
	"like":      "app.bsky.feed.like",
	"post":      "app.bsky.feed.post",
	"follow":    "app.bsky.graph.follow",
	"repost":    "app.bsky.feed.repost",
	"block":     "app.bsky.graph.block",
	"profile":   "app.bsky.actor.profile",
	"list":      "app.bsky.graph.list",
	"listitem":  "app.bsky.graph.listitem",
	"listblock": "app.bsky.graph.listblock",
}

func IsRelevant(kind S) (is bool) {
	for i := range Kinds {
		if kind == Kinds[i] {
			is = true
			break
		}
	}
	return
}

func GetCommon(rr *repo.Repo, createdAt Time, op Op, evt Ev) []types.Tag {
	return []types.Tag{
		{Name: pkg.Protocol, Value: pkg.Bsky},
		{Name: pkg.Kind, Value: pkg.Like},
		{Name: pkg.EventId, Value: op.Cid.String()},
		{Name: pkg.UserId, Value: rr.SignedCommit().Did},
		{Name: pkg.Timestamp, Value: strconv.FormatInt(createdAt.Unix(), 10)},
		{Name: pkg.Repository, Value: evt.Repo},
		{Name: pkg.Path, Value: op.Path},
		{Name: pkg.Signature, Value: hex.EncodeToString(rr.SignedCommit().Sig)},
	}

}
