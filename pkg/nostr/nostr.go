// Package nostr is a set of functions for pulling events from nostr and sending
// them to Arweave AO.
//
// <insert more notes here>
package nostr

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Hubmakerlabs/hoover/pkg"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/event"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/kind"
)

// GetNostrKindToBundle checks the event kind of a nostr event and returns the form used in
// Hoover bundles or an empty string if it is not a relevant kind.
func GetNostrKindToBundle(k kind.T) (s string) {
	switch k {
	case kind.TextNote, kind.LongFormContent:
		s = pkg.Post
	case kind.Repost, kind.GenericRepost:
		s = pkg.Repost
	case kind.Follows:
		s = pkg.Follow
	case kind.Reaction:
		s = pkg.Like
	case kind.MuteList:
		s = pkg.Block
	case kind.ProfileMetadata:
		s = pkg.Profile
	}
	return
}

// EventToBundleItem constructs a bundle from relevant nostr events.
//
// If there is no error but the bundle returns nil the provided event is not relevant as per
// GetNostrKindToBundle.
func EventToBundleItem(ev *event.T) (bundle *types.BundleItem, err error) {
	var k string
	if k = GetNostrKindToBundle(ev.Kind); k == "" {
		return
	}
	bundle = &types.BundleItem{}
	bundle.Data = ev.Content
	bundle.Tags = []types.Tag{
		{Name: pkg.Protocol, Value: pkg.Nostr},
		{Name: pkg.Kind, Value: k},
		{Name: pkg.EventId, Value: ev.ID.String()},
		{Name: pkg.UserId, Value: ev.PubKey},
		{Name: pkg.Timestamp, Value: strconv.FormatInt(ev.CreatedAt.I64(), 10)},
		{Name: pkg.Signature, Value: ev.Sig},
	}
	switch k {
	case pkg.Post:
		for i, t := range ev.Tags {
			switch ev.Tags[i][0] {
			case "e":
				if len(t) > 1 {
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.ReplyTo,
						Value: t[1],
					})
				}
			case "p":
				if len(t) > 1 {
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.Mention,
						Value: t[1],
					})
				}
			case "t":
				if len(t) > 1 {
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.Hashtag,
						Value: t[1],
					})
				}
			case "proxy":
				if len(t) > 1 {
					var sauce, uri S
					uri = t[1]
					if len(t) > 2 {
						sauce = t[2]
					}
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.Source,
						Value: fmt.Sprintf("%s,%s", sauce, uri),
					})
				}
			case "emoji":
				if len(t) > 1 {
					var sauce, uri S
					uri = t[1]
					if len(t) > 2 {
						sauce = t[2]
					}
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.Emoji,
						Value: fmt.Sprintf("%s,%s", uri, sauce),
					})
				}
			case "content-warning":
				var desc S
				if len(t) > 1 {
					desc = t[1]
				}
				bundle.Tags = append(bundle.Tags, types.Tag{
					Name:  pkg.ContentWarning,
					Value: desc,
				})
			case "l":
				if len(t) > 1 {
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.Label,
						Value: t[1],
					})
				}
			case "L":
				if len(t) > 1 {
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.LabelNamespace,
						Value: t[1],
					})
				}
			}
		}
	case pkg.Repost:
		for i, t := range ev.Tags {
			switch ev.Tags[i][0] {
			case "e":
				if len(t) > 1 {
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.RepostEventId,
						Value: t[1],
					})
				}
			case "proxy":
				if len(t) > 1 {
					var sauce, uri S
					uri = t[1]
					if len(t) > 2 {
						sauce = t[2]
					}
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.Source,
						Value: fmt.Sprintf("%s,%s", sauce, uri),
					})
				}
			case "p":
				if len(t) > 1 {
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.Mention,
						Value: t[1],
					})
				}
			case "l":
				if len(t) > 1 {
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.Label,
						Value: t[1],
					})
				}
			case "L":
				if len(t) > 1 {
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.LabelNamespace,
						Value: t[1],
					})
				}
			}
		}
	case pkg.Like:
		for i, t := range ev.Tags {
			switch ev.Tags[i][0] {
			case "e":
				if len(t) > 1 {
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.LikeEventId,
						Value: t[1],
					})
				}
			case "p":
				if len(t) > 1 {
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.Mention,
						Value: t[1],
					})
				}
			case "proxy":
				if len(t) > 1 {
					var sauce, uri S
					uri = t[1]
					if len(t) > 2 {
						sauce = t[2]
					}
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.Source,
						Value: fmt.Sprintf("%s,%s", sauce, uri),
					})
				}
			}
		}
	case pkg.Follow:
		for i, t := range ev.Tags {
			switch ev.Tags[i][0] {
			case "p":
				if len(t) > 1 {
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.FollowUserId,
						Value: t[1],
					})
				}
			case "t":
				if len(t) > 1 {
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.FollowTag,
						Value: t[1],
					})
				}
			}
		}
	case pkg.Block:
		for i, t := range ev.Tags {
			switch ev.Tags[i][0] {
			case "p":
				if len(t) > 1 {
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.BlockUserId,
						Value: t[1],
					})
				}
			case "t":
				if len(t) > 1 {
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.BlockTag,
						Value: t[1],
					})
				}
			}
		}

	case pkg.Profile:
		// remove data field in case it's empty
		bundle.Data = ""
		var prf ProfileMetadata
		if err = json.Unmarshal(B(ev.Content), &prf); chk.E(err) {
			return
		}
		var hasContent bool
		if prf.Name != "" {
			hasContent = true
			bundle.Tags = append(bundle.Tags, types.Tag{
				Name:  pkg.UserName,
				Value: prf.Name,
			})
		}
		if prf.DisplayName != "" {
			hasContent = true
			bundle.Tags = append(bundle.Tags, types.Tag{
				Name:  pkg.DisplayName,
				Value: prf.DisplayName,
			})
		}
		if prf.About != "" {
			hasContent = true
			bundle.Data = prf.About
		}
		if prf.Picture != "" {
			hasContent = true
			bundle.Tags = append(bundle.Tags, types.Tag{
				Name:  pkg.AvatarImage,
				Value: prf.Picture,
			})
		}
		if prf.Banner != "" {
			hasContent = true
			bundle.Tags = append(bundle.Tags, types.Tag{
				Name:  pkg.BannerImage,
				Value: prf.Banner,
			})
		}
		if prf.Website != "" {
			hasContent = true
			bundle.Tags = append(bundle.Tags, types.Tag{
				Name:  pkg.Website,
				Value: prf.Website,
			})
		}
		if prf.NIP05 != "" {
			hasContent = true
			bundle.Tags = append(bundle.Tags, types.Tag{
				Name:  pkg.Verification,
				Value: prf.NIP05,
			})
		}
		if prf.LUD16 != "" {
			hasContent = true
			bundle.Tags = append(bundle.Tags, types.Tag{
				Name:  pkg.PaymentAddress,
				Value: prf.LUD16,
			})
		}
		for i, t := range ev.Tags {
			hasContent = true
			switch ev.Tags[i][0] {
			case "e":
				if len(t) > 1 {
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.ReplyTo,
						Value: t[1],
					})
				}
			case "p":
				if len(t) > 1 {
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.Mention,
						Value: t[1],
					})
				}
			case "t":
				if len(t) > 1 {
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.Hashtag,
						Value: t[1],
					})
				}
			case "proxy":
				if len(t) > 1 {
					var sauce, uri S
					uri = t[1]
					if len(t) > 2 {
						sauce = t[2]
					}
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.Source,
						Value: fmt.Sprintf("%s,%s", sauce, uri),
					})
				}
			case "emoji":
				if len(t) > 1 {
					var sauce, uri S
					uri = t[1]
					if len(t) > 2 {
						sauce = t[2]
					}
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.Emoji,
						Value: fmt.Sprintf("%s,%s", uri, sauce),
					})
				}
			case "content-warning":
				var desc S
				if len(t) > 1 {
					desc = t[1]
				}
				bundle.Tags = append(bundle.Tags, types.Tag{
					Name:  pkg.ContentWarning,
					Value: desc,
				})
			case "l":
				if len(t) > 1 {
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.Label,
						Value: t[1],
					})
				}
			case "L":
				if len(t) > 1 {
					bundle.Tags = append(bundle.Tags, types.Tag{
						Name:  pkg.LabelNamespace,
						Value: t[1],
					})
				}
			}
		}
		if !hasContent {
			bundle = nil
			return
		}
	}
	// for _, tt := range ev.Tags {
	// 	// tags are prefixed by a hash symbol so they don't conflict with the
	// 	// above standard event names
	// 	name := "#" + tt[0]
	// 	var b B
	// 	if b, err = json.Marshal(tt[1:]); chk.E(err) {
	// 		return
	// 	}
	// 	bundle.Tags = append(bundle.Tags,
	// 		types.Tag{Name: name, Value: string(b)})
	// }
	return
}

type ProfileMetadata struct {
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	About       string `json:"about,omitempty"`
	Website     string `json:"website,omitempty"`
	Picture     string `json:"picture,omitempty"`
	Banner      string `json:"banner,omitempty"`
	NIP05       string `json:"nip05,omitempty"`
	LUD16       string `json:"lud16,omitempty"`
}
