package bluesky

import (
	"errors"

	. "github.com/Hubmakerlabs/hoover/pkg"
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

// FromBskyFeedLike is for a like.
//
// In bluesky protocol, the reverse operation, unlike, is actually from a delete operation.
//
// Todo: for now, the reverse operations will not be handled but it should be done for MS2
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
		err = errorf.E("did not get %", BskyKinds(Like))
		return
	}
	if like.Subject == nil {
		err = errors.New("like has no subject, data of no use, refers to nothing")
		return
	}
	bundle = new(types.BundleItem)
	if err = GetCommon(bundle, rr, createdAt, op, evt); chk.E(err) {
		return
	}
	AppendTag(bundle, J(Like, Event, Id), like.Subject.Cid)
	AppendTag(bundle, J(Like, Uri), like.Subject.Uri)
	return
}
