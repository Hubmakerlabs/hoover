package bluesky

import (
	"strconv"
	"time"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/repo"
	"github.com/bluesky-social/indigo/util"
	typegen "github.com/whyrusleeping/cbor-gen"
)

// {
//   "lexicon": 1,
//   "id": "app.bsky.feed.repost",
//   "defs": {
//     "main": {
//       "description": "Record representing a 'repost' of an existing Bluesky post.",
//       "type": "record",
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

// FromBskyFeedRepost is
func FromBskyFeedRepost(
	evt *atproto.SyncSubscribeRepos_Commit,
	op *atproto.SyncSubscribeRepos_RepoOp,
	rr *repo.Repo,
	rec typegen.CBORMarshaler,
) (bundle *types.BundleItem, err error) {

	var createdAt time.Time
	var to any
	if to, createdAt, err = UnmarshalEvent(evt, rec, &bsky.FeedRepost{}); chk.E(err) {
		return
	}
	if to == nil {
		err = errorf.E("failed to unmarshal post")
		return
	}
	repost, ok := to.(*bsky.FeedRepost)
	if !ok {
		err = errorf.E("did not get", Kinds["repost"])
		return
	}
	bundle = &types.BundleItem{}
	bundle.Tags = GetCommon(rr, createdAt, op, evt)
	// bundle.Tags = []types.Tag{
	// 	{Name: "protocol", Value: "bsky"},
	// 	{Name: "kind", Value: Kinds["repost"]},
	// 	{Name: "id", Value: op.Cid.String()},
	// 	{Name: "pubkey", Value: rr.SignedCommit().Did},
	// 	{Name: "created_at", Value: strconv.FormatInt(createdAt.Unix(), 10)},
	// 	{Name: "repo", Value: evt.Repo},
	// 	{Name: "path", Value: op.Path},
	// 	{Name: "sig", Value: hex.EncodeToString(rr.SignedCommit().Sig)},
	// }
	if createdAt, err = time.Parse(util.ISO8601, repost.CreatedAt); chk.E(err) {
		return
	}
	AppendTag(bundle, "#updated_at", strconv.FormatInt(createdAt.Unix(), 10))
	AppendTags(bundle, "#subject", []S{repost.Subject.Cid, repost.Subject.Uri})
	return
}

// ToBskyFeedRepost is
func ToBskyFeedRepost() {

}
