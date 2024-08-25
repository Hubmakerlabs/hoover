package bluesky

import (
	"fmt"
	"time"

	. "github.com/Hubmakerlabs/hoover/pkg"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/bluesky-social/indigo/api/bsky"
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
func FromBskyFeedPost(evt Ev, op Op, rr Repo, rec Rec) (bundle BundleItem, err error) {
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
	if err = GetCommon(bundle, rr, createdAt, op, evt); chk.E(err) {
		return
	}
	if pst.Text != "" {
		bundle.Data = pst.Text
	}
	if pst.Reply != nil {
		if pst.Reply.Root != nil && pst.Reply.Root.Uri != "" {
			AppendTag(bundle, J(Reply, Root, Id), pst.Reply.Parent.Cid)
			AppendTag(bundle, J(Reply, Root, Uri), pst.Reply.Parent.Uri)
		}
		if pst.Reply.Parent != nil && pst.Reply.Parent.Uri != "" {
			AppendTag(bundle, J(Reply, Parent, Id), pst.Reply.Parent.Cid)
			AppendTag(bundle, J(Reply, Parent, Uri), pst.Reply.Parent.Uri)
		}
	}
	if pst.Embed != nil {
		if pst.Embed.EmbedRecord != nil {
			AppendTag(bundle, J(Embed, Record, Uri), pst.Embed.EmbedRecord.Record.Uri)
			AppendTag(bundle, J(Embed, Record, Id), pst.Embed.EmbedRecord.Record.Cid)
		}
		if pst.Embed.EmbedImages != nil {
			// var count int
			// for _, embeds := range pst.Embed.EmbedImages.Images {
			//
			// }
		}
		if pst.Embed.EmbedRecordWithMedia != nil {
			EmbedExternalRecordWithMedia(bundle, Embed, pst.Embed.EmbedRecordWithMedia)
		}
		if pst.Embed.EmbedExternal != nil {
			EmbedExternal(bundle, J(Embed, External), pst.Embed.EmbedExternal)
		}
	}
	if pst.Entities != nil {
		for i, entity := range pst.Entities {
			var index string
			if entity.Index != nil {
				index = fmt.Sprintf("%d-%d", entity.Index.Start, entity.Index.End)
			}
			AppendTag(bundle, J(Entities, i, Index), index)
			AppendTag(bundle, J(Entities, i, Type), entity.Type)
			AppendTag(bundle, J(Entities, i, Value), entity.Value)
		}
	}
	if pst.Facets != nil {
		for i := range pst.Facets {
			if pst.Facets[i].Features != nil {
				for _, feats := range pst.Facets[i].Features {
					if feats.RichtextFacet_Mention != nil {
						if feats.RichtextFacet_Mention.Did != "" {
							AppendTag(bundle, J(Richtext, Mention),
								feats.RichtextFacet_Mention.Did)
						}
					}
					if feats.RichtextFacet_Link != nil {
						if feats.RichtextFacet_Link.Uri != "" {
							AppendTag(bundle, J(Richtext, Link),
								feats.RichtextFacet_Link.Uri)
						}
					}
					if feats.RichtextFacet_Tag != nil {
						if feats.RichtextFacet_Tag.Tag != "" {
							AppendTag(bundle, J(Richtext, Tag),
								feats.RichtextFacet_Tag.Tag)
						}
					}
				}
			}
			if pst.Facets[i].Index != nil {
				AppendTag(bundle, J(Richtext, Tag, Start),
					fmt.Sprint(pst.Facets[i].Index.ByteStart))
				AppendTag(bundle, J(Richtext, Tag, End),
					fmt.Sprint(pst.Facets[i].Index.ByteEnd))

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
							AppendTag(bundle, fmt.Sprintf("%s-%03d", Label, i), labels[i])
						}
						// AppendTags(bundle, "#labels", labels)
					}
				}
			}
		}
	}
	if pst.Langs != nil && len(pst.Langs) > 0 {
		if len(pst.Langs) == 1 {
			AppendTag(bundle, fmt.Sprintf("%s", Language), pst.Langs[0])
		} else {
			for i := range pst.Langs {
				AppendTag(bundle, fmt.Sprintf("%s-%03d", Language, i), pst.Langs[i])
			}
		}
		// AppendTags(bundle, "#langs", pst.Langs)
	}
	if pst.Tags != nil && len(pst.Tags) > 0 {
		for i := range pst.Tags {
			AppendTag(bundle, J(Tag, i), pst.Tags[i])
		}
	}
	return
}
