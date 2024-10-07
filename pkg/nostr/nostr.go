// Package nostr is a set of functions for pulling events from nostr and sending
// them to Arweave AO.
//
// <insert more notes here>
package nostr

import (
	"encoding/json"
	"fmt"
	"strconv"

	. "github.com/Hubmakerlabs/hoover/pkg"
	ao "github.com/Hubmakerlabs/hoover/pkg/arweave"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/event"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/kind"
)

// GetNostrKindToBundle checks the event kind of a nostr event and returns the form used in
// Hoover bundles or an empty string if it is not a relevant kind.
func GetNostrKindToBundle(k kind.T) (s string) {
	switch k {
	case kind.TextNote, kind.LongFormContent:
		s = Post
	case kind.Repost, kind.GenericRepost:
		s = Repost
	case kind.Follows:
		s = Follow
	case kind.Reaction:
		s = Like
	case kind.ProfileMetadata:
		s = Profile
	}
	return
}

var RelevantKinds = []int{
	kind.TextNote.ToInt(), kind.LongFormContent.ToInt(),
	kind.Reaction.ToInt(),
	kind.Follows.ToInt(),
	kind.Repost.ToInt(), kind.GenericRepost.ToInt(),
	kind.ProfileMetadata.ToInt(),
}

// EventToBundleItem constructs a bundle from relevant nostr events.
//
// If there is no error but the bundle returns nil the provided event is not relevant as per
// GetNostrKindToBundle.
func EventToBundleItem(ev *event.T, relay string) (bundle *types.BundleItem, err error) {
	var k string
	if k = GetNostrKindToBundle(ev.Kind); k == "" {
		return
	}
	bundle = &types.BundleItem{}
	// bundle.Data = ev.Content
	data := ao.NewEventData(ev.Content)
	bundle.Tags = []types.Tag{
		{J(App, Name), AppNameValue},
		{J(App, Version), AppVersion},
		{Protocol, Nostr},
		{Repository, relay},
		{Kind, k},
		{J(Event, Id), ev.ID.String()},
		{J(User, Id), ev.PubKey},
		{J(Unix, Time), strconv.FormatInt(ev.CreatedAt.I64(), 10)},
		{Signature, ev.Sig},
	}
out:
	switch k {
	case Post:
		for i, t := range ev.Tags {
			switch ev.Tags[i][0] {
			case "e":
				if len(t) > 1 {
					// reply parent/roots need to be in the tags
					if len(t) > 3 {
						// it probably has a reply relation specifier
						switch t[3] {
						case "root":
							ao.AppendTag(bundle, J(Reply, Root, Id), t[1])
						case "reply":
							ao.AppendTag(bundle, J(Reply, Parent, Id), t[1])
						}
					} else {
						ao.AppendTag(bundle, J(Reply, Parent, Id), t[1])
					}
				}
			case "p":
				if len(t) > 1 {
					data.Append(Mention, t[1])
				}
			case "t":
				if len(t) > 1 {
					data.Append(Hashtag, t[1])
				}
			case "proxy":
				if len(t) > 1 {
					var sauce, uri S
					uri = t[1]
					if len(t) > 2 {
						sauce = t[2]
					}
					data.Append(J(Source), sauce)
					data.Append(J(Source, Uri), uri)
				}
			case "emoji":
				if len(t) > 1 {
					var sauce, uri S
					uri = t[1]
					if len(t) > 2 {
						sauce = t[2]
					}
					data.Append(Emoji, fmt.Sprintf("%s,%s", uri, sauce))
				}
			case "content-warning":
				var desc S
				if len(t) > 1 {
					desc = t[1]
				}
				// these also need to be indexable i think
				ao.AppendTag(bundle, J(Content, Warning), desc)
			case "l":
				if len(t) > 1 {
					data.Append(Label, t[1])
				}
			case "L":
				if len(t) > 1 {
					data.Append(J(Label, Namespace), t[1])
				}
			}
		}
	case Repost:
		for i, t := range ev.Tags {
			switch ev.Tags[i][0] {
			case "e":
				if len(t) > 1 {
					ao.AppendTag(bundle, J(Repost, Event, Id), t[1])
				}
			case "proxy":
				if len(t) > 1 {
					var sauce, uri S
					uri = t[1]
					if len(t) > 2 {
						sauce = t[2]
					}
					data.Append(J(Source), sauce)
					data.Append(J(Source, Uri), uri)
				}
			case "p":
				if len(t) > 1 {
					data.Append(Mention, t[1])
				}
			case "l":
				if len(t) > 1 {
					data.Append(Label, t[1])
				}
			case "L":
				if len(t) > 1 {
					data.Append(J(Label, Namespace), t[1])
				}
			}
		}
	case Like:
		for i, t := range ev.Tags {
			switch ev.Tags[i][0] {
			case "e":
				if len(t) > 1 {
					// likes need to also be in the tags
					ao.AppendTag(bundle, J(Like, Event, Id), t[1])
				}
			case "p":
				if len(t) > 1 {
					data.Append(Mention, t[1])
				}
			case "proxy":
				if len(t) > 1 {
					var sauce, uri S
					uri = t[1]
					if len(t) > 2 {
						sauce = t[2]
					}
					data.Append(J(Source), sauce)
					data.Append(J(Source, Uri), uri)
				}
			}
		}
	case Follow:
		// we don't need the content field of follow events
		data.Content = ""
		for i, t := range ev.Tags {
			switch ev.Tags[i][0] {
			case "p":
				if len(t) > 1 {
					data.Append(J(Follow, User, Id), t[1])
				}
			case "t":
				if len(t) > 1 {
					data.Append(J(Follow, Tag), t[1])
				}
			}
		}
	case Profile:
		// remove data field in case it's empty
		bundle.Data = ""
		var prf ProfileMetadata
		if err = json.Unmarshal(B(ev.Content), &prf); chk.E(err) {
			log.I.F("%s", ev.Content)
			break out
		}
		if prf.Name != "" {
			ao.AppendTag(bundle, J(User, Name),
				prf.Name)
		}
		if prf.DisplayName != "" {
			ao.AppendTag(bundle, J(Display, Name), prf.DisplayName)
		}
		if prf.About != "" {
			data.Append(Bio, prf.About)
		}
		if prf.Picture != "" {
			ao.AppendTag(bundle, J(Avatar, Image), prf.Picture)
		}
		if prf.Banner != "" {
			ao.AppendTag(bundle, J(Banner, Image), prf.Banner)
		}
		if prf.Website != "" {
			data.Append(Website, prf.Website)
		}
		if nip05, ok := prf.NIP05.(string); ok {
			if nip05 != "" {
				data.Append(Verification, nip05)
			}
		}
		if prf.LUD16 != "" {
			data.Append(J(Payment, Address), prf.LUD16)
		}
		for i, t := range ev.Tags {
			switch ev.Tags[i][0] {
			case "e":
				if len(t) > 1 {
					data.Append(J(Mention, Event, Id), t[1])
				}
			case "p":
				if len(t) > 1 {
					data.Append(Mention, t[1])
				}
			case "t":
				if len(t) > 1 {
					data.Append(Hashtag, t[1])
				}
			case "proxy":
				if len(t) > 1 {
					var sauce, uri S
					uri = t[1]
					if len(t) > 2 {
						sauce = t[2]
					}
					data.Append(Source, sauce, uri)
					data.Append(J(Source, Uri), uri)
				}
			case "emoji":
				if len(t) > 1 {
					var sauce, uri S
					uri = t[1]
					if len(t) > 2 {
						sauce = t[2]
					}
					data.Append(Emoji, fmt.Sprintf("%s,%s", uri, sauce))
				}
			case "content-warning":
				var desc S
				if len(t) > 1 {
					desc = t[1]
				}
				data.Append(J(Content, Warning), desc)
			case "l":
				if len(t) > 1 {
					data.Append(Label, t[1])
				}
			case "L":
				if len(t) > 1 {
					data.Append(J(Label, Namespace), t[1])
				}
			}
		}
	}
	// put the ao.EventData into JSON form and place in the bundle.Data field
	var b []byte
	b, err = json.Marshal(data)
	if err != nil {
		return
	}
	bundle.Data = string(b)
	return
}

type ProfileMetadata struct {
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	About       string `json:"about,omitempty"`
	Website     string `json:"website,omitempty"`
	Picture     string `json:"picture,omitempty"`
	Banner      string `json:"banner,omitempty"`
	NIP05       any    `json:"nip05,omitempty"`
	LUD16       string `json:"lud16,omitempty"`
}
