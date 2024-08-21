package bluesky

import (
	"github.com/Hubmakerlabs/hoover/pkg"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/bluesky-social/indigo/api/bsky"
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
func FromBskyFeedLike(evt Ev, op Op, rr Repo, rec typegen.CBORMarshaler,
) (bundle *types.BundleItem, err error) {

	var createdAt Time
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
		err = errorf.E("did not get app.bsky.feed.like")
		return
	}
	if like.Subject == nil {
		err = errorf.E("like has no subject, data of no use, refers to nothing")
		return
	}
	bundle = &types.BundleItem{}
	bundle.Tags = GetCommon(rr, createdAt, op, evt)
	// bundle.Tags = []types.Tag{
	// 	{Name: pkg.Protocol, Value: pkg.Bsky},
	// 	{Name: pkg.Kind, Value: pkg.Like},
	// 	{Name: pkg.EventId, Value: op.Cid.String()},
	// 	{Name: pkg.UserId, Value: rr.SignedCommit().Did},
	// 	{Name: pkg.Timestamp, Value: strconv.FormatInt(createdAt.Unix(), 10)},
	// 	{Name: pkg.Repository, Value: evt.Repo},
	// 	{Name: pkg.Path, Value: op.Path},
	// 	{Name: pkg.Signature, Value: hex.EncodeToString(rr.SignedCommit().Sig)},
	// }
	AppendTag(bundle, pkg.LikeEventId, like.Subject.Cid)
	AppendTag(bundle, pkg.URI, like.Subject.Uri)
	return
}
