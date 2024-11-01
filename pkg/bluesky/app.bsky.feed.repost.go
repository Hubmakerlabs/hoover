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
func FromBskyFeedRepost(evt Ev, op Op, rr Repo, rec Rec, data *ao.EventData) (bundle BundleItem,
	err error) {
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
		err = errorf.E("did not get", Kinds(Repost))
		return
	}
	bundle = &types.BundleItem{}
	var userID, protocol, timestamp string
	if userID, protocol, timestamp, err = GetCommon(bundle, rr, createdAt, op, evt); chk.E(err) {
		return
	}
	postId := repost.Subject.Cid
	ao.AppendTag(bundle, J(Repost, Event, Id), postId)
	data.Append(J(Repost, Event, Uri), repost.Subject.Uri)

	title := userID + " reposted on " + protocol + " at " + timestamp
	ao.AppendTag(bundle, Title, title)

	description := userID + " reposted on " + protocol + " at " + timestamp + ". Id of original post: " + postId
	ao.AppendTag(bundle, Description, description[:min(len(description), 300)])
	return
}

// ToBskyFeedRepost is
func ToBskyFeedRepost() {

}
