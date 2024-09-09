package nostr

import (
	"sort"
	"sync"
	"time"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/context"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/event"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/eventid"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/kind"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/tag"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/tags"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/timestamp"
	"github.com/nbd-wtf/go-nostr"
)

var Relays = []string{
	"wss://purplepag.es",
	"wss://njump.me",
	"wss://relay.snort.social",
	"wss://relay.damus.io",
	"wss://relay.primal.net",
	"wss://relay.nostr.band",
}

type sortId struct {
	id string
	ts int64
}

type sortIds []sortId

func (s sortIds) Len() int           { return len(s) }
func (s sortIds) Less(i, j int) bool { return s[i].ts < s[j].ts }
func (s sortIds) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Firehose connects to a list of relays and pulls recent events of relevant kinds, bundles
// them in arweave bundle items and runs a provided closure on them.
//
// If the closure returns an error, the Firehose halts and returns to the caller.
func Firehose(c context.T, cancel context.F, wg *sync.WaitGroup, relays []S,
	fn func(bundle *types.BundleItem) (err error)) {

	var ready bool
	wg.Add(1)
	pool := nostr.NewSimplePool(c)
	stop := func() {
		pool.Relays.Range(func(_ string, relay *nostr.Relay) bool {
			relay.Close()
			return true
		})
		cancel()
	}
	defer stop()
	var err error
	ff := nostr.Filters{{Kinds: RelevantKinds}}
	idMap := make(map[string]int64)
	for evt := range pool.SubMany(c, relays, ff) {
		var bundle *types.BundleItem
		var ev *event.T
		if ev, err = ToEvent(evt.Event); err != nil {
			continue
		}
		if _, ok := idMap[evt.ID]; ok {
			// skip it if we've seen it
			continue
		}
		idMap[evt.ID] = time.Now().UnixMilli()
		if bundle, err = EventToBundleItem(ev); chk.E(err) {
			continue
		}
		if bundle == nil {
			continue
		}
		if !ready {
			ready = true
			wg.Done()
		}
		wg.Wait()
		if err = fn(bundle); err != nil {
			return
		}
		// prune the idMap as we probably won't see the same event again from any of the relays
		// after (half) this many
		if len(idMap) > 4096 {
			var ids sortIds
			for i := range idMap {
				ids = append(ids, sortId{id: i, ts: idMap[i]})
			}
			sort.Sort(ids)
			// trim the top half
			idMap = make(map[string]int64)
			for _, id := range ids[:2048] {
				idMap[id.id] = id.ts
			}
		}
	}
}

func ToEvent(evt *nostr.Event) (ev *event.T, err error) {
	var id eventid.T
	id, err = eventid.New(evt.ID)
	if err != nil {
		return
	}
	var tgs tags.T
	for _, tt := range evt.Tags {
		var t tag.T
		for _, ttt := range tt {
			t = append(t, ttt)
		}
		tgs = append(tgs, t)
	}
	ev = &event.T{
		ID:        id,
		PubKey:    evt.PubKey,
		CreatedAt: timestamp.T(evt.CreatedAt),
		Kind:      kind.T(evt.Kind),
		Tags:      tgs,
		Content:   evt.Content,
		Sig:       evt.Sig,
	}
	return
}
