package bluesky

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Hubmakerlabs/hoover/pkg"
	ao "github.com/Hubmakerlabs/hoover/pkg/arweave"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/repo"
	"github.com/bluesky-social/indigo/repomgr"
	"github.com/ipfs/go-cid"
	"github.com/whyrusleeping/cbor-gen"
)

func RepoCommit(ctx context.Context,
	cancel context.CancelFunc, fn func(bundle *types.BundleItem) (err error)) func(
	evt *atproto.SyncSubscribeRepos_Commit) (err error) {
	return func(evt *atproto.SyncSubscribeRepos_Commit) (err error) {
		var rr *repo.Repo
		if rr, err = repo.ReadRepoFromCar(ctx, bytes.NewReader(evt.Blocks)); chk.E(err) {
			return
		}
		var bundle *types.BundleItem
		data := &ao.EventData{}
		for _, op := range evt.Ops {
			ek := repomgr.EventKind(op.Action)
			switch ek {
			case repomgr.EvtKindCreateRecord, repomgr.EvtKindUpdateRecord:
				var rc cid.Cid
				var rec typegen.CBORMarshaler
				if rc, rec, err = rr.GetRecord(ctx, op.Path); err != nil {
					err = fmt.Errorf("getting record %s (%s) within seq %d for %s: %w",
						op.Path, *op.Cid, evt.Seq, evt.Repo, err)
					return nil
				}
				if util.LexLink(rc) != *op.Cid {
					err = errorf.E("mismatch in record and op cid: %s != %s", rc, *op.Cid)
					return
				}
				switch {
				case strings.HasPrefix(op.Path, Kinds(pkg.Post)):
					if bundle, err = FromBskyFeedPost(evt, op, rr, rec, data); chk.E(err) {
						err = nil
						continue
					}
				case strings.HasPrefix(op.Path, Kinds(pkg.Like)):
					if bundle, err = FromBskyFeedLike(evt, op, rr, rec, data); err != nil {
						err = nil
						continue
					}
				case strings.HasPrefix(op.Path, Kinds(pkg.Follow)):
					if bundle, err = FromBskyGraphFollow(evt, op, rr, rec); chk.E(err) {
						err = nil
						continue
					}
				case strings.HasPrefix(op.Path, Kinds(pkg.Repost)):
					if bundle, err = FromBskyFeedRepost(evt, op, rr, rec, data); chk.E(err) {
						err = nil
						continue
					}
				case strings.HasPrefix(op.Path, Kinds(pkg.Profile)):
					if bundle, err = FromBskyActorProfile(evt, op, rr, rec, data); chk.E(err) {
						err = nil
						continue
					}
				}
			default:
				// log.I.Ln(ek)
			}
		}
		if bundle == nil {
			return
		}
		if data.Content != "" || data.EventTags != nil {
			// put the ao.EventData into JSON form and place in the bundle.Data field
			var b []byte
			b, err = json.Marshal(data)
			if err != nil {
				return
			}
			bundle.Data = string(b)
			if err = fn(bundle); err != nil {
				cancel()
				return
			}
		}
		return
	}
}
