package bluesky

import (
	"fmt"
	"strconv"
	"time"

	. "github.com/Hubmakerlabs/hoover/pkg"
	ao "github.com/Hubmakerlabs/hoover/pkg/arweave"
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
func FromBskyFeedPost(evt Ev, op Op, rr Repo, rec Rec, data *ao.EventData) (bundle BundleItem,
	err error) {
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
		data.Content = pst.Text
	}
	if pst.Reply != nil {
		if pst.Reply.Root != nil && pst.Reply.Root.Uri != "" {
			ao.AppendTag(bundle, J(Reply, Root, Id), pst.Reply.Parent.Cid)
			data.Append(J(Reply, Root, Uri), pst.Reply.Parent.Uri)
		}
		if pst.Reply.Parent != nil && pst.Reply.Parent.Uri != "" {
			ao.AppendTag(bundle, J(Reply, Parent, Id), pst.Reply.Parent.Cid)
			data.Append(J(Reply, Parent, Uri), pst.Reply.Parent.Uri)
		}
	}
	if pst.Embed != nil {
		if pst.Embed.EmbedRecord != nil {
			data.Append(J(Embed, Record, Uri), pst.Embed.EmbedRecord.Record.Uri)
			data.Append(J(Embed, Record, Id), pst.Embed.EmbedRecord.Record.Cid)
		}
		if pst.Embed.EmbedImages != nil {
			var count int
			if len(pst.Embed.EmbedImages.Images) == 1 {
				AppendImageTags(data, J(Embed, Image), pst.Embed.EmbedImages.Images[0])
			} else {
				for _, embed := range pst.Embed.EmbedImages.Images {
					if embed != nil {
						AppendImageTags(data, J(Embed, Image, count), embed)
						count++
					}
				}
			}
		}
		if pst.Embed.EmbedRecordWithMedia != nil {
			EmbedExternalRecordWithMedia(data, Embed, pst.Embed.EmbedRecordWithMedia)
		}
		if pst.Embed.EmbedExternal != nil {
			EmbedExternal(data, J(Embed, External), pst.Embed.EmbedExternal)
		}
	}
	if pst.Entities != nil {
		for i, entity := range pst.Entities {
			data.Append(J(Entities, i, Index, Start),
				strconv.FormatInt(entity.Index.Start, 10))
			data.Append(J(Entities, i, Index), strconv.FormatInt(entity.Index.End, 10))
			data.Append(J(Entities, i, Type), entity.Type)
			data.Append(J(Entities, i, Value), entity.Value)
		}
	}
	if pst.Facets != nil {
		for i := range pst.Facets {
			if pst.Facets[i].Features != nil {
				var prefix string
				if len(pst.Facets) == 1 {
					prefix = J(Richtext)
				} else {
					prefix = J(Richtext, i)
				}
				for j, feats := range pst.Facets[i].Features {
					if len(pst.Facets[i].Features) == 1 {
						prefix = J(prefix)
					} else {
						prefix = J(prefix, j)
					}
					if feats.RichtextFacet_Mention != nil {
						if feats.RichtextFacet_Mention.Did != "" {
							data.Append(J(prefix, Mention),
								feats.RichtextFacet_Mention.Did)
						}
					}
					if feats.RichtextFacet_Link != nil {
						if feats.RichtextFacet_Link.Uri != "" {
							data.Append(J(prefix, Uri),
								feats.RichtextFacet_Link.Uri)
						}
					}
					if feats.RichtextFacet_Tag != nil {
						if feats.RichtextFacet_Tag.Tag != "" {
							data.Append(J(Hashtag),
								feats.RichtextFacet_Tag.Tag)
						}
					}
				}
			}
			if pst.Facets[i].Index != nil {
				if len(pst.Facets) == 1 {
					data.Append(J(Richtext, Tag, Start),
						fmt.Sprint(pst.Facets[i].Index.ByteStart))
					data.Append(J(Richtext, Tag, End),
						fmt.Sprint(pst.Facets[i].Index.ByteEnd))
				} else {
					data.Append(J(Richtext, i, Tag, Start),
						fmt.Sprint(pst.Facets[i].Index.ByteStart))
					data.Append(J(Richtext, i, Tag, End),
						fmt.Sprint(pst.Facets[i].Index.ByteEnd))
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
						if len(labels) == 1 {
							data.Append(J(Label), labels[0])
						} else {
							for i := range labels {
								data.Append(J(Label, i), labels[i])
							}
						}
					}
				}
			}
		}
	}
	if pst.Langs != nil && len(pst.Langs) > 0 {
		if len(pst.Langs) == 1 {
			data.Append(Language, pst.Langs[0])
		} else {
			for i := range pst.Langs {
				data.Append(J(Language, i), pst.Langs[i])
			}
		}
		// AppendTags(bundle, "#langs", pst.Langs)
	}
	if pst.Tags != nil && len(pst.Tags) > 0 {
		if len(pst.Tags) == 1 {
			data.Append(Tag, pst.Tags[0])
		} else {
			for i := range pst.Tags {
				data.Append(J(Tag, i), pst.Tags[i])
			}
		}
	}
	return
}
