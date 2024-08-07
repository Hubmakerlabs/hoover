package bluesky

import (
	"encoding/hex"
	"encoding/json"
	"strconv"
	"time"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/repo"
	"github.com/bluesky-social/indigo/util"
	"github.com/whyrusleeping/cbor-gen"
)

// {
//   "lexicon": 1,
//   "id": "app.bsky.feed.like",
//   "defs": {
//     "main": {
//       "type": "record",
//       "description": "Record declaring a 'like' of a piece of subject content.",
//       "key": "tid",
//       "record": {
//         "type": "object",
//         "required": ["subject", "createdAt"],
//         "properties": {
//           "subject": { "type": "ref", "ref": "com.atproto.repo.strongRef" },
//           "createdAt": { "type": "string", "format": "datetime" }
//         }
//       }
//     }
//   }
// }

// FromBskyFeedLike is
func FromBskyFeedLike(evt *atproto.SyncSubscribeRepos_Commit, op *atproto.SyncSubscribeRepos_RepoOp, rr *repo.Repo,
	rec typegen.CBORMarshaler) (bundle *types.BundleItem, err error) {
	banana := lexutil.LexiconTypeDecoder{
		Val: rec,
	}
	pst := bsky.FeedLike{}
	var b B
	if b, err = banana.MarshalJSON(); chk.E(err) {
		return
	}
	if err = json.Unmarshal(b, &pst); chk.E(err) {
		return
	}
	bundle = &types.BundleItem{}
	var createdAt time.Time
	if createdAt, err = time.Parse(util.ISO8601, evt.Time); chk.E(err) {
		return
	}
	bundle.Tags = []types.Tag{
		{Name: "source", Value: "bsky"},
		{Name: "id", Value: op.Cid.String()},
		{Name: "pubkey", Value: rr.SignedCommit().Did},
		{Name: "created_at", Value: strconv.FormatInt(createdAt.Unix(), 10)},
		{Name: "kind", Value: "app.bsky.feed.like"},
		{Name: "repo", Value: evt.Repo},
		{Name: "path", Value: op.Path},
		{Name: "sig", Value: hex.EncodeToString(rr.SignedCommit().Sig)},
	}
	return
}

// todo: no way we are trying to reconstruct this lol... just for it to be authenticated, ok, enough.
// // ToBskyFeedLike is
// func ToBskyFeedLike() {
//
// }
