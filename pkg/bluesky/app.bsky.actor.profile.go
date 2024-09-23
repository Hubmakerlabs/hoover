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
//   "id": "app.bsky.actor.profile",
//   "defs": {
//     "main": {
//       "type": "record",
//       "description": "A declaration of a Bluesky account profile.",
//       "key": "literal:self",
//       "record": {
//         "type": "object",
//         "properties": {
//           "displayName": {
//             "type": "string",
//             "maxGraphemes": 64,
//             "maxLength": 640
//           },
//           "description": {
//             "type": "string",
//             "description": "Free-form profile description text.",
//             "maxGraphemes": 256,
//             "maxLength": 2560
//           },
//           "avatar": {
//             "type": "blob",
//             "description": "Small image to be displayed next to posts from account. AKA, 'profile picture'",
//             "accept": ["image/png", "image/jpeg"],
//             "maxSize": 1000000
//           },
//           "banner": {
//             "type": "blob",
//             "description": "Larger horizontal image to display behind profile view.",
//             "accept": ["image/png", "image/jpeg"],
//             "maxSize": 1000000
//           },
//           "labels": {
//             "type": "union",
//             "description": "Self-label values, specific to the Bluesky application, on the overall account.",
//             "refs": ["com.atproto.label.defs#selfLabels"]
//           },
//           "joinedViaStarterPack": {
//             "type": "ref",
//             "ref": "com.atproto.repo.strongRef"
//           },
//           "createdAt": { "type": "string", "format": "datetime" }
//         }
//       }
//     }
//   }
// }

// FromBskyActorProfile is
func FromBskyActorProfile(evt Ev, op Op, rr Repo, rec Rec,
	data *ao.EventData) (bundle BundleItem, err error) {
	var createdAt time.Time
	var to any
	if to, createdAt, err = UnmarshalEvent(evt, rec, &bsky.ActorProfile{}); chk.E(err) {
		return
	}
	if to == nil {
		err = errorf.E("failed to unmarshal post")
		return
	}
	profile, ok := to.(*bsky.ActorProfile)
	if !ok {
		err = errorf.E("did not get", Kinds(Profile))
		return
	}
	bundle = &types.BundleItem{}
	if err = GetCommon(bundle, rr, createdAt, op, evt); chk.E(err) {
		return
	}
	if profile.DisplayName != nil && *profile.DisplayName != "" {
		data.Append(J(Display, Name), *profile.DisplayName)
	}
	if profile.Description != nil {
		data.Append(Bio, *profile.Description)
	}
	if profile.Avatar != nil {
		AppendLexBlobTags(data, J(Avatar, Image), profile.Avatar)
	}
	if profile.Banner != nil {
		AppendLexBlobTags(data, J(Banner, Image), profile.Banner)
	}
	return
}

// ToBskyActorProfile is
func ToBskyActorProfile() {

}
