package bluesky

import (
	"errors"

	"github.com/Hubmakerlabs/hoover/pkg"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/bluesky-social/indigo/api/bsky"
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
func FromBskyFeedLike(evt Ev, op Op, rr Repo, rec Rec) (bundle BundleItem, err error) {
	var createdAt Time
	var to any
	if to, createdAt, err = UnmarshalEvent(evt, rec, &bsky.FeedLike{}); err != nil {
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
		err = errors.New("like has no subject, data of no use, refers to nothing")
		return
	}
	bundle = &types.BundleItem{}
	bundle.Tags = GetCommon(rr, createdAt, op, evt)
	AppendTag(bundle, pkg.LikeEventId, like.Subject.Cid)
	AppendTag(bundle, pkg.URI, like.Subject.Uri)
	return
}
