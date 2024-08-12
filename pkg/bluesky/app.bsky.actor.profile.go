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
	typegen "github.com/whyrusleeping/cbor-gen"
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
func FromBskyActorProfile(
	evt *atproto.SyncSubscribeRepos_Commit,
	op *atproto.SyncSubscribeRepos_RepoOp,
	rr *repo.Repo,
	rec typegen.CBORMarshaler,
) (bundle *types.BundleItem, err error) {

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
		err = errorf.E("did not get", Kinds["profile"])
		return
	}
	bundle = &types.BundleItem{}
	bundle.Tags = []types.Tag{
		{Name: "protocol", Value: "bsky"},
		{Name: "kind", Value: Kinds["profile"]},
		{Name: "id", Value: op.Cid.String()},
		{Name: "pubkey", Value: rr.SignedCommit().Did},
		{Name: "created_at", Value: strconv.FormatInt(createdAt.Unix(), 10)},
		{Name: "repo", Value: evt.Repo},
		{Name: "path", Value: op.Path},
		{Name: "sig", Value: hex.EncodeToString(rr.SignedCommit().Sig)},
	}
	if profile.CreatedAt!=nil {
		if createdAt, err = time.Parse(util.ISO8601, *profile.CreatedAt); chk.E(err) {
			return
		}
	}
	AppendTag(bundle, "#updated_at", strconv.FormatInt(createdAt.Unix(), 10))
	if profile.DisplayName!=nil {
		AppendTag(bundle, "#displayname", *profile.DisplayName)
	}
	if profile.Description!=nil {
		AppendTag(bundle, "#description", *profile.Description)
	}
	if profile.Avatar != nil {
		AppendTags(bundle, "#avatar", GetLexBlobTags(profile.Avatar))
	}
	if profile.Banner != nil {
		AppendTags(bundle, "#banner", GetLexBlobTags(profile.Banner))
	}
	if profile.JoinedViaStarterPack != nil {
		AppendTags(bundle, "#starterpack", []S{profile.JoinedViaStarterPack.Cid, profile.JoinedViaStarterPack.Uri})
	}
	return
}

// ToBskyActorProfile is
func ToBskyActorProfile() {

}
