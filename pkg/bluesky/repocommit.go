package bluesky

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/Hubmakerlabs/hoover/pkg"
	"github.com/Hubmakerlabs/hoover/pkg/arweave"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/repo"
	"github.com/bluesky-social/indigo/repomgr"
	"github.com/ipfs/go-cid"
	"github.com/whyrusleeping/cbor-gen"
)

func RepoCommit(ctx context.Context,
	cancel context.CancelFunc) func(evt *atproto.SyncSubscribeRepos_Commit) (err error) {
	return func(evt *atproto.SyncSubscribeRepos_Commit) (err error) {
		var rr *repo.Repo
		if rr, err = repo.ReadRepoFromCar(ctx, bytes.NewReader(evt.Blocks)); chk.E(err) {
			return
		}
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
					var post *types.BundleItem
					if post, err = FromBskyFeedPost(evt, op, rr, rec); chk.E(err) {
						// normally would return but this shuts down the firehose processing
						err = nil
						continue
						// return
					}
					_ = post
					arweave.PrintBundleItem(post)
					fmt.Println()
					case strings.HasPrefix(op.Path, Kinds(pkg.Like)):
						var like *types.BundleItem
						if like, err = FromBskyFeedLike(evt, op, rr, rec); err != nil {
							// normally would return but this shuts down the firehose processing
							err = nil
							continue
							// return
						}
						_ = like
						arweave.PrintBundleItem(like)
						fmt.Println()
					case strings.HasPrefix(op.Path, Kinds(pkg.Follow)):
						var follow *types.BundleItem
						if follow, err = FromBskyGraphFollow(evt, op, rr, rec); chk.E(err) {
							// normally would return but this shuts down the firehose processing
							// err = nil
							// continue
							return
						}
						_ = follow
						arweave.PrintBundleItem(follow)
						fmt.Println()
					case strings.HasPrefix(op.Path, Kinds(pkg.Repost)):
						var repost *types.BundleItem
						if repost, err = FromBskyFeedRepost(evt, op, rr, rec); chk.E(err) {
							// normally would return but this shuts down the firehose processing
							// err = nil
							// continue
							return
						}
						_ = repost
						arweave.PrintBundleItem(repost)
						fmt.Println()
					case strings.HasPrefix(op.Path, Kinds(pkg.Profile)):
						var profile *types.BundleItem
						if profile, err = FromBskyActorProfile(evt, op, rr, rec); chk.E(err) {
							// normally would return but this shuts down the firehose processing
							// err = nil
							// continue
							return
						}
						_ = profile
						arweave.PrintBundleItem(profile)
						fmt.Println()
				}
			default:
				// log.I.Ln(ek)
			}
		}
		// cancel()
		_ = rr
		return
	}
}
