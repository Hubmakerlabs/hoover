package nostr

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	ao "github.com/Hubmakerlabs/hoover/pkg/arweave"
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
		if bundle, err = EventToBundleItem(ev, "archive"); chk.E(err) {
			t.Fatal(err)
		}
		if bundle == nil {
			continue
		}
		// fmt.Println()
		// fmt.Println(S(b))
		fmt.Println()
		ao.PrintBundleItem(bundle)
		b = b[:0]
		b = nil
	}
}
