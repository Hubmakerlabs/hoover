package bluesky

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Hubmakerlabs/hoover/pkg/arweave"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/bluesky-social/indigo/api/atproto"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/repo"
	"github.com/bluesky-social/indigo/repomgr"
	"github.com/ipfs/go-cid"
	typegen "github.com/whyrusleeping/cbor-gen"
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
				if rc, rec, err = rr.GetRecord(ctx, op.Path); chk.D(err) {
					err = errorf.E("getting record %s (%s) within seq %d for %s: %w",
						op.Path, *op.Cid, evt.Seq, evt.Repo, err)
					return nil
				}
				if lexutil.LexLink(rc) != *op.Cid {
					err = errorf.E("mismatch in record and op cid: %s != %s", rc, *op.Cid)
					return
				}
				switch {
				case strings.HasPrefix(op.Path, Kinds["like"]):
					// var like *types.BundleItem
					// if like, err = FromBskyFeedLike(evt, op, rr, rec); chk.E(err) {
					// 	// normally would return but this shuts down the firehose processing
					// 	err = nil
					// 	continue
					// 	// return
					// }
					// _ = like
					// arweave.PrintBundleItem(like)
					// fmt.Println()
				case strings.HasPrefix(op.Path, Kinds["post"]):
					// var post *types.BundleItem
					// if post, err = FromBskyFeedPost(evt, op, rr, rec); chk.E(err) {
					// 	// normally would return but this shuts down the firehose processing
					// 	err = nil
					// 	continue
					// 	// return
					// }
					// _ = post
					// arweave.PrintBundleItem(post)
					// fmt.Println()
				case strings.HasPrefix(op.Path, Kinds["follow"]):
					// var follow *types.BundleItem
					// if follow, err = FromBskyGraphFollow(evt, op, rr, rec); chk.E(err) {
					// 	// normally would return but this shuts down the firehose processing
					// 	// err = nil
					// 	// continue
					// 	return
					// }
					// _ = follow
					// arweave.PrintBundleItem(follow)
					// fmt.Println()
				case strings.HasPrefix(op.Path, Kinds["repost"]):
					// var repost *types.BundleItem
					// if repost, err = FromBskyFeedRepost(evt, op, rr, rec); chk.E(err) {
					// 	// normally would return but this shuts down the firehose processing
					// 	// err = nil
					// 	// continue
					// 	return
					// }
					// _ = repost
					// arweave.PrintBundleItem(repost)
					// fmt.Println()
				case strings.HasPrefix(op.Path, Kinds["block"]):
					var block *types.BundleItem
					if block, err = FromBskyGraphBlock(evt, op, rr, rec); chk.E(err) {
						// normally would return but this shuts down the firehose processing
						// err = nil
						// continue
						return
					}
					_ = block
					arweave.PrintBundleItem(block)
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

func RepoHandle() func(handle *atproto.SyncSubscribeRepos_Handle) error {
	return func(handle *atproto.SyncSubscribeRepos_Handle) error {
		b, err := json.MarshalIndent(handle, "", "  ")
		if err != nil {
			return err
		}
		log.I.F("RepoHandle:\n%s", b)
		return nil
	}
}

func RepoInfo() func(info *atproto.SyncSubscribeRepos_Info) error {
	return func(info *atproto.SyncSubscribeRepos_Info) error {
		b, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return err
		}
		log.I.F("RepoInfo\r%s", b)
		log.I.F("INFO: %s: %v\n", info.Name, info.Message)
		return nil
	}
}

func RepoMigrate() func(mig *atproto.SyncSubscribeRepos_Migrate) error {
	return func(mig *atproto.SyncSubscribeRepos_Migrate) error {
		b, err := json.Marshal(mig)
		if err != nil {
			return err
		}
		fmt.Println("RepoMigrate")
		fmt.Println(string(b))
		return nil
	}
}

func RepoTombstone() func(tomb *atproto.SyncSubscribeRepos_Tombstone) error {
	return func(tomb *atproto.SyncSubscribeRepos_Tombstone) error {
		b, err := json.Marshal(tomb)
		if err != nil {
			return err
		}
		fmt.Println("RepoTombstone")
		fmt.Println(string(b))
		return nil
	}
}

func LabelLabels() func(labels *atproto.LabelSubscribeLabels_Labels) error {
	return func(labels *atproto.LabelSubscribeLabels_Labels) error {
		b, err := json.Marshal(labels)
		if err != nil {
			return err
		}
		fmt.Println("LabelLabels")
		fmt.Println(string(b))
		return nil
	}
}

func LabelInfo() func(info *atproto.LabelSubscribeLabels_Info) error {
	return func(info *atproto.LabelSubscribeLabels_Info) error {
		b, err := json.Marshal(info)
		if err != nil {
			return err
		}
		fmt.Println("LabelInfo")
		fmt.Println(string(b))
		return nil
	}
}
