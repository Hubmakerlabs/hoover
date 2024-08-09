package nostr

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Hubmakerlabs/hoover/pkg/arweave"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/event"
	"github.com/mleku/nodl/pkg/event/examples"
)

var evs = examples.Cache

func TestEventToBundleItem(t *testing.T) {
	scanner := bufio.NewScanner(bytes.NewBuffer(examples.Cache))
	buf := make(B, 1_000_000)
	scanner.Buffer(buf, len(buf))
	var err error
	for scanner.Scan() {
		b := scanner.Bytes()
		bc := make(B, len(b))
		copy(bc, b)
		ev := &event.T{}
		if err = json.Unmarshal(b, ev); chk.E(err) {
			t.Fatal(err)
		}
		var bundle *types.BundleItem
		if bundle, err = EventToBundleItem(ev); chk.E(err) {
			t.Fatal(err)
		}
		fmt.Println()
		fmt.Println()
		fmt.Println(S(b))
		fmt.Println()
		arweave.PrintBundleItem(bundle)
		var ev2 *event.T
		if ev2, err = BundleItemToEvent(bundle); chk.E(err) {
			t.Fatal(err)
		}
		var b2 B
		if b2, err = ev2.MarshalJSON(); chk.E(err) {
			t.Fatal(err)
		}
		if !equals(bc, b2) {
			t.Errorf("DIFFERENT\n%s\n\n%s\n\n", bc, b2)
		}
		b = b[:0]
		b = nil
	}
}
