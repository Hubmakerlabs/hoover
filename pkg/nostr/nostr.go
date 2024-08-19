// Package nostr is a set of functions for pulling events from nostr and sending
// them to Arweave AO.
//
// <insert more notes here>
package nostr

import (
	"encoding/json"
	"strconv"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/event"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/kind"
)

func GetNostrKindToBundle(k kind.T) (s string) {
	switch k {
	case kind.Repost, kind.GenericRepost:
		s = "Repost"
	}
	return
}

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
	var b B
	if b, err = json.Marshal(ev.Content); chk.E(err) {
		return
	}
	bundle.Data = S(b)
	bundle.Tags = []types.Tag{
		{Name: "Protocol", Value: "nostr"},
		{Name: "Event-ID", Value: ev.ID.String()},
		{Name: "User-ID", Value: ev.PubKey},
		{Name: "Timestamp", Value: strconv.FormatInt(ev.CreatedAt.I64(), 10)},
		{Name: "Kind", Value: strconv.Itoa(int(ev.Kind.ToUint16()))},
		{Name: "Signature", Value: ev.Sig},
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
