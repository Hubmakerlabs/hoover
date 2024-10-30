package bluesky

import (
	"time"

	. "github.com/Hubmakerlabs/hoover/pkg"
	ao "github.com/Hubmakerlabs/hoover/pkg/arweave"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/bluesky-social/indigo/api/bsky"
)

// {
//   "lexicon": 1,
//   "id": "app.bsky.graph.follow",
//   "defs": {
//     "main": {
//       "type": "record",
//       "description": "Record declaring a social 'follow' relationship of another account. Duplicate follows will be ignored by the AppView.",
//       "key": "tid",
//       "record": {
//         "type": "object",
//         "required": ["subject", "createdAt"],
//         "properties": {
//           "subject": { "type": "string", "format": "did" },
//           "createdAt": { "type": "string", "format": "datetime" }
//         }
//       }
//     }
//   }
// }

// FromBskyGraphFollow is for a follow.
//
// In bluesky protocol, the reverse operation, unfollow, is actually from a delete operation.
//
// Todo: for now, the reverse operations will not be handled but it should be done for MS2
func FromBskyGraphFollow(evt Ev, op Op, rr Repo, rec Rec) (bundle BundleItem, err error) {
	var createdAt time.Time
	var to any
	if to, createdAt, err = UnmarshalEvent(evt, rec, &bsky.GraphFollow{}); chk.E(err) {
		return
	}
	if to == nil {
		err = errorf.E("failed to unmarshal post")
		return
	}
	fol, ok := to.(*bsky.GraphFollow)
	if !ok {
		err = errorf.E("did not get", Kinds(Follow))
		return
	}
	bundle = &types.BundleItem{}
	var userID, protocol, timestamp string
	if userID, protocol, timestamp, err = GetCommon(bundle, rr, createdAt, op, evt); chk.E(err) {
		return
	}
	follow_id := fol.Subject
	ao.AppendTag(bundle, J(Follow, User, Id), follow_id)
	title := userID + " followed another user on " + protocol + " at " + timestamp
	ao.AppendTag(bundle, Title, title)

	description := userID + " followed " + follow_id + " on " + protocol + " at " + timestamp
	ao.AppendTag(bundle, Description, description)
	return
}

// ToBskyGraphFollow is
func ToBskyGraphFollow() {

}
