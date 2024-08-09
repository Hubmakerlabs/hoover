package bluesky

import (
	"encoding/hex"
	"strconv"
	"time"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
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
func FromBskyFeedLike(
	evt *atproto.SyncSubscribeRepos_Commit,
	op *atproto.SyncSubscribeRepos_RepoOp,
	rr *repo.Repo,
	rec typegen.CBORMarshaler,
) (bundle *types.BundleItem, err error) {

	var createdAt time.Time
	var to any
	if to, createdAt, err = UnmarshalEvent(evt, rec, &bsky.FeedLike{}); chk.E(err) {
		return
	}
	if to == nil {
		err = errorf.E("failed to unmarshal post")
		return
	}
	like, ok := to.(*bsky.FeedLike)
	if !ok {
		err = errorf.E("did not get app.bsky.feed.post")
		return
	}
	// banana := lexutil.LexiconTypeDecoder{
	// 	Val: rec,
	// }
	// pst := bsky.FeedLike{}
	// var b B
	// if b, err = banana.MarshalJSON(); chk.E(err) {
	// 	return
	// }
	// if err = json.Unmarshal(b, &pst); chk.E(err) {
	// 	return
	// }
	// var createdAt time.Time
	// if createdAt, err = time.Parse(util.ISO8601, evt.Time); chk.E(err) {
	// 	return
	// }
	if like.Subject == nil {
		err = errorf.E("like has no subject, data of no use, refers to nothing")
		return
	}
	bundle = &types.BundleItem{}
	bundle.Tags = []types.Tag{
		{Name: "protocol", Value: "bsky"},
		{Name: "id", Value: op.Cid.String()},
		{Name: "pubkey", Value: rr.SignedCommit().Did},
		{Name: "created_at", Value: strconv.FormatInt(createdAt.Unix(), 10)},
		{Name: "kind", Value: "app.bsky.feed.like"},
		{Name: "repo", Value: evt.Repo},
		{Name: "path", Value: op.Path},
		{Name: "sig", Value: hex.EncodeToString(rr.SignedCommit().Sig)},
	}
	if createdAt, err = time.Parse(util.ISO8601, like.CreatedAt); chk.E(err) {
		return
	}
	AppendTag(bundle, "#sent_timestamp", strconv.FormatInt(createdAt.Unix(), 10))
	AppendTags(bundle, "#subject", []string{like.Subject.Cid, like.Subject.Uri})
	return
}

// todo: no way we are trying to reconstruct this lol... just for it to be authenticated, ok, enough.

// // ToBskyFeedLike is
// func ToBskyFeedLike() {
//
// }
