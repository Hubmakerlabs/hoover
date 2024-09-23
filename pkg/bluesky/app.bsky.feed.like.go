package bluesky

import (
	"errors"
	"strings"
	"time"

	. "github.com/Hubmakerlabs/hoover/pkg"
	ao "github.com/Hubmakerlabs/hoover/pkg/arweave"
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
func FromBskyFeedLike(evt Ev, op Op, rr Repo, rec Rec, data *ao.EventData) (bundle BundleItem,
	err error) {
	var createdAt time.Time
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
	ao.AppendTag(bundle, J(Like, Event, Id), like.Subject.Cid)
	// so there is a mention with the poster's ID to search on
	// data.Append(J(Like, Uri), like.Subject.Uri)
	s1 := strings.Split(like.Subject.Uri, "://")
	if len(s1) > 1 {
		s2 := strings.Split(s1[1], "/")
		if len(s2) > 1 {
			data.Append(J(Like, Path), strings.Join(s2[1:], "/"))
		}
		if len(s2) > 0 {
			data.Append(Mention, s2[0])
		}
	}
	return
}
