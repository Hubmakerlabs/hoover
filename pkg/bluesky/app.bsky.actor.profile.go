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
	var userID, protocol, timestamp string
	if userID, protocol, timestamp, err = GetCommon(bundle, rr, createdAt, op, evt); chk.E(err) {
		return
	}
	changes := []string{}
	if profile.DisplayName != nil && *profile.DisplayName != "" {
		ao.AppendTag(bundle, J(Display, Name), *profile.DisplayName)
		changes = append(changes, "display name")
	}
	if profile.Description != nil {
		data.Append(Bio, *profile.Description)
		changes = append(changes, "bio")
	}
	if profile.Avatar != nil {
		ao.AppendTag(bundle, J(Avatar, Image), profile.Avatar.Ref.String())
		AppendLexBlobTags(data, J(Avatar, Image), profile.Avatar)
		changes = append(changes, "avatar")
	}
	if profile.Banner != nil {
		ao.AppendTag(bundle, J(Banner, Image), profile.Banner.Ref.String())
		AppendLexBlobTags(data, J(Banner, Image), profile.Banner)
		changes = append(changes, "banner")
	}
	var change string
	var noChange bool
	if len(changes) == 0 {
		noChange = true
	} else if len(changes) == 1 {
		change = changes[0]
	} else if len(changes) == 2 {
		change = changes[0] + " and " + changes[1]
	} else {
		for i, s := range changes {

			if i == len(changes)-1 {
				change += "and a " + s
			} else {
				change += "a " + s + ", "
			}
		}
	}
	if !noChange {
		change = ". New profile includes " + change
	}
	title := "Profile Update:" + userID + " updated their profile on " + protocol + " at " + timestamp
	ao.AppendTag(bundle, Title, title)
	description := "Profile Update:" + userID + " updated their profile on " + protocol + " at " + timestamp + change
	ao.AppendTag(bundle, Description, description[:min(300, len(description))])
	return
}

// ToBskyActorProfile is
func ToBskyActorProfile() {

}
