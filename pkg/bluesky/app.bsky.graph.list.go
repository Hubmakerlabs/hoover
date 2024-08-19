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
//   "id": "app.bsky.graph.list",
//   "defs": {
//     "main": {
//       "type": "record",
//       "description": "Record representing a list of accounts (actors). Scope includes both moderation-oriented lists and curration-oriented lists.",
//       "key": "tid",
//       "record": {
//         "type": "object",
//         "required": ["name", "purpose", "createdAt"],
//         "properties": {
//           "purpose": {
//             "type": "ref",
//             "description": "Defines the purpose of the list (aka, moderation-oriented or curration-oriented)",
//             "ref": "app.bsky.graph.defs#listPurpose"
//           },
//           "name": {
//             "type": "string",
//             "maxLength": 64,
//             "minLength": 1,
//             "description": "Display name for list; can not be empty."
//           },
//           "description": {
//             "type": "string",
//             "maxGraphemes": 300,
//             "maxLength": 3000
//           },
//           "descriptionFacets": {
//             "type": "array",
//             "items": { "type": "ref", "ref": "app.bsky.richtext.facet" }
//           },
//           "avatar": {
//             "type": "blob",
//             "accept": ["image/png", "image/jpeg"],
//             "maxSize": 1000000
//           },
//           "labels": {
//             "type": "union",
//             "refs": ["com.atproto.label.defs#selfLabels"]
//           },
//           "createdAt": { "type": "string", "format": "datetime" }
//         }
//       }
//     }
//   }
// }

// FromBskyGraphList is
func FromBskyGraphList(
	evt *atproto.SyncSubscribeRepos_Commit,
	op *atproto.SyncSubscribeRepos_RepoOp,
	rr *repo.Repo,
	rec typegen.CBORMarshaler,
) (bundle *types.BundleItem, err error) {

	var createdAt time.Time
	var to any
	if to, createdAt, err = UnmarshalEvent(evt, rec, &bsky.GraphList{}); chk.E(err) {
		return
	}
	if to == nil {
		err = errorf.E("failed to unmarshal post")
		return
	}
	list, ok := to.(*bsky.GraphList)
	if !ok {
		err = errorf.E("did not get", Kinds["list"])
		return
	}
	bundle = &types.BundleItem{}
	bundle.Tags = []types.Tag{
		{Name: "protocol", Value: "bsky"},
		{Name: "kind", Value: Kinds["list"]},
		{Name: "id", Value: op.Cid.String()},
		{Name: "pubkey", Value: rr.SignedCommit().Did},
		{Name: "created_at", Value: strconv.FormatInt(createdAt.Unix(), 10)},
		{Name: "repo", Value: evt.Repo},
		{Name: "path", Value: op.Path},
		{Name: "sig", Value: hex.EncodeToString(rr.SignedCommit().Sig)},
	}
	if createdAt, err = time.Parse(util.ISO8601, list.CreatedAt); chk.E(err) {
		return
	}
	AppendTag(bundle, "#updated_at", strconv.FormatInt(createdAt.Unix(), 10))
	if list.Name != "" {
		AppendTag(bundle, "#name", list.Name)
	}
	if list.Purpose != nil {
		AppendTag(bundle, "#purpose", *list.Purpose)
	}
	if list.Avatar != nil {
		AppendTags(bundle, "#avatar", GetLexBlobTags(list.Avatar))
	}
	if list.Description != nil {
		AppendTag(bundle, "#description", *list.Description)
	}
	if list.DescriptionFacets != nil {
		for i := range list.DescriptionFacets {
			if list.DescriptionFacets[i].Features != nil {
				for _, feats := range list.DescriptionFacets[i].Features {
					if feats.RichtextFacet_Mention != nil {
						if feats.RichtextFacet_Mention.Did != "" {
							AppendTag(bundle,
								"#facet_features_richtext_mention",
								feats.RichtextFacet_Mention.Did)
						}
					}
					if feats.RichtextFacet_Link != nil {
						if feats.RichtextFacet_Link.Uri != "" {
							AppendTag(bundle,
								"#facet_features_richtext_link", feats.RichtextFacet_Link.Uri)
						}
					}
					if feats.RichtextFacet_Tag != nil {
						if feats.RichtextFacet_Tag.Tag != "" {
							AppendTag(bundle,
								"#facet_features_richtext_tag", feats.RichtextFacet_Tag.Tag)
						}
					}
				}
			}
		}
	}
	if list.Labels != nil {
		if list.Labels.LabelDefs_SelfLabels != nil {
			if list.Labels.LabelDefs_SelfLabels.Values != nil {
				if list.Labels.LabelDefs_SelfLabels.Values != nil {
					var labels []string
					for _, label := range list.Labels.LabelDefs_SelfLabels.Values {
						if label != nil {
							labels = append(labels)
						}
					}
					if len(labels) > 0 {
						AppendTags(bundle, "#labels", labels)
					}
				}
			}
		}
	}

	return
}

// ToBskyGraphList is
func ToBskyGraphList() {

}
