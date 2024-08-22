// Package bluesky is a set of functions for pulling events from bluesky and
// sending them to Arweave AO.
//
// <insert more notes here>
package bluesky

import (
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/Hubmakerlabs/hoover/pkg"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/repo"
	typegen "github.com/whyrusleeping/cbor-gen"
)

// bluesky names are crazy ugly stutter parties

type (
	Repo       = *repo.Repo
	Time       = time.Time
	Op         = *atproto.SyncSubscribeRepos_RepoOp
	Ev         = *atproto.SyncSubscribeRepos_Commit
	Rec        = typegen.CBORMarshaler
	BundleItem = *types.BundleItem
)

var Kinds = map[string]string{
	pkg.Like:    "app.bsky.feed.like",
	pkg.Post:    "app.bsky.feed.post",
	pkg.Follow:  "app.bsky.graph.follow",
	pkg.Repost:  "app.bsky.feed.repost",
	pkg.Block:   "app.bsky.graph.block",
	pkg.Profile: "app.bsky.actor.profile",
}
var BskyKinds = map[string]string{
	"app.bsky.feed.like":     pkg.Like,
	"app.bsky.feed.post":     pkg.Post,
	"app.bsky.graph.follow":  pkg.Follow,
	"app.bsky.feed.repost":   pkg.Repost,
	"app.bsky.graph.block":   pkg.Block,
	"app.bsky.actor.profile": pkg.Profile,
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
	split := strings.Split(op.Path, "/")
	return []types.Tag{
		{Name: pkg.Protocol, Value: pkg.Bsky},
		{Name: pkg.Kind, Value: BskyKinds[split[0]]},
		{Name: pkg.EventId, Value: op.Cid.String()},
		{Name: pkg.UserId, Value: rr.SignedCommit().Did},
		{Name: pkg.Timestamp, Value: strconv.FormatInt(createdAt.Unix(), 10)},
		{Name: pkg.Repository, Value: evt.Repo},
		{Name: pkg.Path, Value: op.Path},
		{Name: pkg.Signature, Value: hex.EncodeToString(rr.SignedCommit().Sig)},
	}

}
