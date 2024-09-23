// Package bluesky is a set of functions for pulling events from bluesky and sending them to
// Arweave AO.
package bluesky

import (
	"sync"

	. "github.com/Hubmakerlabs/hoover/pkg"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/repo"
	typegen "github.com/whyrusleeping/cbor-gen"
)

// bluesky names are crazy ugly stutter parties

type (
	Repo       = *repo.Repo
	Op         = *atproto.SyncSubscribeRepos_RepoOp
	Ev         = *atproto.SyncSubscribeRepos_Commit
	Rec        = typegen.CBORMarshaler
	BundleItem = *types.BundleItem
)

var mx sync.Mutex

func Kinds(k string) (s string) {
	mx.Lock()
	defer mx.Unlock()
	var ok bool
	if s, ok = kinds[k]; ok {
		return
	}
	return ""
}
func BskyKinds(k string) (s string) {
	mx.Lock()
	defer mx.Unlock()
	var ok bool
	if s, ok = bskyKinds[k]; ok {
		return
	}
	return ""
}

var kinds = map[string]string{
	Like:    "app.bsky.feed.like",
	Post:    "app.bsky.feed.post",
	Follow:  "app.bsky.graph.follow",
	Repost:  "app.bsky.feed.repost",
	Block:   "app.bsky.graph.block",
	Profile: "app.bsky.actor.profile",
}
var bskyKinds = map[string]string{
	"app.bsky.feed.like":     Like,
	"app.bsky.feed.post":     Post,
	"app.bsky.graph.follow":  Follow,
	"app.bsky.feed.repost":   Repost,
	"app.bsky.graph.block":   Block,
	"app.bsky.actor.profile": Profile,
}

func IsRelevant(kind string) (is bool) {
	mx.Lock()
	defer mx.Unlock()
	for i := range kinds {
		if kind == kinds[i] {
			is = true
			break
		}
	}
	return
}
