// Package nostr is a set of functions for pulling events from nostr and sending
// them to Arweave AO.
//
// <insert more notes here>
package nostr

import (
	"encoding/json"
	"strconv"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/Hubmakerlabs/replicatr/pkg/ec/hex"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/event"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/eventid"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/kind"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/tag"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/timestamp"
	"github.com/mleku/nodl/pkg/text"
)

// EventToBundleItem constructs the data parts of an Arweave Bundle prior to
// adding the cryptographic authentication parts.
//
// The "content" field of the event goes in the "data" field, and all of the
// event fields are named the same as the event JSON. The tags are the same
// except the first field of each tag is prefixed with a `#` to designate it as
// a tag, and the Value is the JSON escaped string of the JSON encoding of a
// slice of strings, like:
//
//	"[\"string1\",\"string2\"]"
//
// To distinguish between protocols, the first tag is "source" and its value is
// the protocol name, in this case, "nostr"
func EventToBundleItem(ev *event.T) (bundle *types.BundleItem, err error) {
	bundle = &types.BundleItem{}
	bundle.Data = string(text.NostrEscape(nil, text.B(ev.Content)))
	bundle.Tags = []types.Tag{
		{Name: "source", Value: "nostr"},
		{Name: "id", Value: ev.ID.String()},
		{Name: "pubkey", Value: ev.PubKey},
		{Name: "created_at", Value: strconv.FormatInt(ev.CreatedAt.I64(), 10)},
		{Name: "kind", Value: strconv.Itoa(int(ev.Kind.ToUint16()))},
		{Name: "sig", Value: ev.Sig},
	}
	for _, tt := range ev.Tags {
		// tags are prefixed by a hash symbol so they don't conflict with the
		// above standard event names
		name := "#" + tt[0]
		var b B
		if b, err = json.Marshal(tt[1:]); chk.E(err) {
			return
		}
		bundle.Tags = append(bundle.Tags,
			types.Tag{Name: name, Value: string(b)})
	}
	return
}

func BundleItemToEvent(bundle *types.BundleItem) (ev *event.T, err error) {
	// first check that the first tag is nostr
	if bundle.Tags[0].Name != "source" && bundle.Tags[0].Value == "nostr" {
		err = errorf.E("first tag of bundle is not \"source\" and value is not \"nostr\"")
		return
	}
	ev = &event.T{}
	ev.Content = string(text.NostrUnescape([]byte(bundle.Data)))
	for _, tt := range bundle.Tags[1:] {
		switch tt.Name {
		case "id":
			ev.ID = eventid.T(tt.Value)
			b := make(B, len(ev.ID)/2)
			if _, err = hex.Decode(b, B(ev.ID)); chk.E(err) {
				return
			}
		case "pubkey":
			ev.PubKey = tt.Value
			b := make(B, len(ev.PubKey)/2)
			if _, err = hex.Decode(b, B(ev.PubKey)); chk.E(err) {
				return
			}
		case "created_at":
			var ca int
			if ca, err = strconv.Atoi(tt.Value); chk.E(err) {
				return
			}
			ev.CreatedAt = timestamp.T(ca)
		case "kind":
			var k int
			if k, err = strconv.Atoi(tt.Value); chk.E(err) {
				return
			}
			ev.Kind = kind.T(k)
		case "sig":
			ev.Sig = tt.Value
			b := make(B, len(ev.Sig)/2)
			if _, err = hex.Decode(b, B(ev.Sig)); chk.E(err) {
				return
			}
		default:
			if tt.Name[0] != '#' {
				err = errorf.E("tags must have a # prefix '%s'", tt.Name)
				return
			}
			var val tag.T
			if err = json.Unmarshal([]byte(tt.Value), &val); chk.E(err) {
				return
			}
			t := tag.T{tt.Name[1:]}
			t = append(t, val...)
			ev.Tags = append(ev.Tags, t)
		}
	}
	return
}
