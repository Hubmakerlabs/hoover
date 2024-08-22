package bluesky

import (
	"fmt"
	"time"

	"github.com/Hubmakerlabs/hoover/pkg"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/repo"
	typegen "github.com/whyrusleeping/cbor-gen"
)

// {
//   "lexicon": 1,
//   "id": "app.bsky.feed.post",
//   "defs": {
//     "main": {
//       "type": "record",
//       "description": "Record containing a Bluesky post.",
//       "key": "tid",
//       "record": {
//         "type": "object",
//         "required": ["text", "createdAt"],
//         "properties": {
//           "text": {
//             "type": "string",
//             "maxLength": 3000,
//             "maxGraphemes": 300,
//             "description": "The primary post content. May be an empty string, if there are embeds."
//           },
//           "entities": {
//             "type": "array",
//             "description": "DEPRECATED: replaced by app.bsky.richtext.facet.",
//             "items": { "type": "ref", "ref": "#entity" }
//           },
//           "facets": {
//             "type": "array",
//             "description": "Annotations of text (mentions, URLs, hashtags, etc)",
//             "items": { "type": "ref", "ref": "app.bsky.richtext.facet" }
//           },
//           "reply": { "type": "ref", "ref": "#replyRef" },
//           "embed": {
//             "type": "union",
//             "refs": [
//               "app.bsky.embed.images",
//               "app.bsky.embed.external",
//               "app.bsky.embed.record",
//               "app.bsky.embed.recordWithMedia"
//             ]
//           },
//           "langs": {
//             "type": "array",
//             "description": "Indicates human language of post primary text content.",
//             "maxLength": 3,
//             "items": { "type": "string", "format": "language" }
//           },
//           "labels": {
//             "type": "union",
//             "description": "Self-label values for this post. Effectively content warnings.",
//             "refs": ["com.atproto.label.defs#selfLabels"]
//           },
//           "tags": {
//             "type": "array",
//             "description": "Additional hashtags, in addition to any included in post text and facets.",
//             "maxLength": 8,
//             "items": { "type": "string", "maxLength": 640, "maxGraphemes": 64 }
//           },
//           "createdAt": {
//             "type": "string",
//             "format": "datetime",
//             "description": "Client-declared timestamp when this post was originally created."
//           }
//         }
//       }
//     },
//     "replyRef": {
//       "type": "object",
//       "required": ["root", "parent"],
//       "properties": {
//         "root": { "type": "ref", "ref": "com.atproto.repo.strongRef" },
//         "parent": { "type": "ref", "ref": "com.atproto.repo.strongRef" }
//       }
//     },
//     "entity": {
//       "type": "object",
//       "description": "Deprecated: use facets instead.",
//       "required": ["index", "type", "value"],
//       "properties": {
//         "index": { "type": "ref", "ref": "#textSlice" },
//         "type": {
//           "type": "string",
//           "description": "Expected values are 'mention' and 'link'."
//         },
//         "value": { "type": "string" }
//       }
//     },
//     "textSlice": {
//       "type": "object",
//       "description": "Deprecated. Use app.bsky.richtext instead -- A text segment. Start is inclusive, end is exclusive. Indices are for utf16-encoded strings.",
//       "required": ["start", "end"],
//       "properties": {
//         "start": { "type": "integer", "minimum": 0 },
//         "end": { "type": "integer", "minimum": 0 }
//       }
//     }
//   }
// }

// FromBskyFeedPost is
func FromBskyFeedPost(
	evt *atproto.SyncSubscribeRepos_Commit,
	op *atproto.SyncSubscribeRepos_RepoOp,
	rr *repo.Repo,
	rec typegen.CBORMarshaler,
) (bundle *types.BundleItem, err error) {

	var createdAt time.Time
	var to any
	if to, createdAt, err = UnmarshalEvent(evt, rec, &bsky.FeedPost{}); chk.E(err) {
		return
	}
	if to == nil {
		err = errorf.E("failed to unmarshal post")
		return
	}
	pst, ok := to.(*bsky.FeedPost)
	if !ok {
		err = errorf.E("did not get app.bsky.feed.post")
		return
	}
	bundle = &types.BundleItem{}
	bundle.Tags = GetCommon(rr, createdAt, op, evt)
	if pst.Embed != nil {
		if pst.Embed.EmbedRecord != nil {
			AppendTags(bundle, "#embedrecord",
				[]S{pst.Embed.EmbedRecord.Record.Uri, pst.Embed.EmbedRecord.Record.Cid})
		}
		if pst.Embed.EmbedImages != nil {
			for _, img := range pst.Embed.EmbedImages.Images {
				if img != nil {
					EmbedImages(bundle, pkg.EmbedImage, img)
				}
			}
		}
		if pst.Embed.EmbedRecordWithMedia != nil {
			EmbedExternalRecordWithMedia(bundle, pkg.Embed, pst.Embed.EmbedRecordWithMedia)
		}
		if pst.Embed.EmbedExternal != nil {
			EmbedExternal(bundle, pkg.EmbedExternal, pst.Embed.EmbedExternal)
		}
	}
	if pst.Entities != nil {
		for i, entity := range pst.Entities {
			var index string
			if entity.Index != nil {
				index = fmt.Sprintf("%d-%d", entity.Index.Start, entity.Index.End)
			}
			AppendTags(bundle, fmt.Sprintf("%s-%03d", pkg.Entities, i), []string{index,
				entity.Type, entity.Value})
		}
	}
	if pst.Facets != nil {
		for i := range pst.Facets {
			if pst.Facets[i].Features != nil {
				for _, feats := range pst.Facets[i].Features {
					if feats.RichtextFacet_Mention != nil {
						if feats.RichtextFacet_Mention.Did != "" {
							AppendTag(bundle,
								fmt.Sprintf("%s-Mention", pkg.Richtext),
								// "#facet_features_richtext_mention",
								feats.RichtextFacet_Mention.Did)
						}
					}
					if feats.RichtextFacet_Link != nil {
						if feats.RichtextFacet_Link.Uri != "" {
							AppendTag(bundle,
								fmt.Sprintf("%s-Link", pkg.Richtext),
								feats.RichtextFacet_Link.Uri)
						}
					}
					if feats.RichtextFacet_Tag != nil {
						if feats.RichtextFacet_Tag.Tag != "" {
							AppendTag(bundle,
								fmt.Sprintf("%s-Tag", pkg.Richtext),
								// "#facet_features_richtext_tag",
								feats.RichtextFacet_Tag.Tag)
						}
					}
				}
			}
		}
	}
	if pst.Labels != nil {
		if pst.Labels.LabelDefs_SelfLabels != nil {
			if pst.Labels.LabelDefs_SelfLabels.Values != nil {
				if pst.Labels.LabelDefs_SelfLabels.Values != nil {
					var labels []string
					for _, label := range pst.Labels.LabelDefs_SelfLabels.Values {
						if label != nil {
							labels = append(labels)
						}
					}
					if len(labels) > 0 {
						for i := range labels {
							AppendTag(bundle, fmt.Sprintf("%s-%03d", pkg.Label, i), labels[i])
						}
						// AppendTags(bundle, "#labels", labels)
					}
				}
			}
		}
	}
	if pst.Langs != nil && len(pst.Langs) > 0 {
		if len(pst.Langs) == 1 {
			AppendTag(bundle, fmt.Sprintf("%s", pkg.Language), pst.Langs[0])
		} else {
			for i := range pst.Langs {
				AppendTag(bundle, fmt.Sprintf("%s-%03d", pkg.Language, i), pst.Langs[i])
			}
		}
		// AppendTags(bundle, "#langs", pst.Langs)
	}
	if pst.Reply != nil {
		if pst.Reply.Root != nil && pst.Reply.Root.Uri != "" {
			AppendTag(bundle, fmt.Sprintf("%s-%s-%s", pkg.Reply, pkg.Root, pkg.Id),
				pst.Reply.Parent.Cid)
			AppendTag(bundle, fmt.Sprintf("%s-%s-%s", pkg.Reply, pkg.Root, pkg.URI),
				pst.Reply.Parent.Uri)
			// AppendTags(bundle, fmt.Sprintf("%s-%s", pkg.Reply, pkg.Root),
			// 	[]string{pst.Reply.Root.Cid, pst.Reply.Root.Uri})
		}
		if pst.Reply.Parent != nil && pst.Reply.Parent.Uri != "" {
			AppendTag(bundle, fmt.Sprintf("%s-%s-%s", pkg.Reply, pkg.Parent, pkg.Id),
				pst.Reply.Parent.Cid)
			AppendTag(bundle, fmt.Sprintf("%s-%s-%s", pkg.Reply, pkg.Parent, pkg.URI),
				pst.Reply.Parent.Uri)
		}
	}
	if pst.Tags != nil && len(pst.Tags) > 0 {
		for i := range pst.Tags {
			AppendTag(bundle, fmt.Sprintf("%s-%03d", pkg.Tag, i), pst.Tags[i])
		}
		// AppendTags(bundle, "#tags", pst.Tags)
	}
	return
}
