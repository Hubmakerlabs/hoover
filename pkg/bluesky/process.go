package bluesky

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"strings"

	comatproto "github.com/bluesky-social/indigo/api/atproto"
	appbsky "github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/repo"
	"github.com/bluesky-social/indigo/repomgr"
	"github.com/ipfs/go-cid"
	typegen "github.com/whyrusleeping/cbor-gen"
)

func RepoCommit(ctx context.Context,
	cancel context.CancelFunc) func(evt *comatproto.SyncSubscribeRepos_Commit) (err error) {
	return func(evt *comatproto.SyncSubscribeRepos_Commit) (err error) {
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
				if rc, rec, err = rr.GetRecord(ctx, op.Path); chk.E(err) {
					err = errorf.E("getting record %s (%s) within seq %d for %s: %w",
						op.Path, *op.Cid, evt.Seq, evt.Repo, err)
					return nil
				}
				if lexutil.LexLink(rc) != *op.Cid {
					err = errorf.E("mismatch in record and op cid: %s != %s", rc, *op.Cid)
					return
				}
				if strings.HasPrefix(op.Path, "app.bsky.feed.like") {
					log.I.S(FromBskyFeedLike(evt, op, rr, rec))
					return
				}
				if strings.HasPrefix(op.Path, "app.bsky.feed.post") {
					banana := lexutil.LexiconTypeDecoder{
						Val: rec,
					}
					pst := appbsky.FeedPost{}
					var b B
					if b, err = banana.MarshalJSON(); chk.E(err) {
						return
					}
					if err = json.Unmarshal(b, &pst); chk.E(err) {
						return
					}
					// var xrpcc *xrpc.Client
					// var userProfile *appbsky.ActorDefs_ProfileViewDetailed
					// var replyUserProfile *appbsky.ActorDefs_ProfileViewDetailed
					// log.I.S(rr)
					// log.I.S(pst)
					// Handle if its a post
					switch pst.LexiconTypeID {
					case "app.bsky.feed.post":
						var reply string
						if len(pst.Text) == 0 {
							return
						}
						if pst.Reply != nil {
							reply = "Parent: " + pst.Reply.Parent.Uri + "\nRoot: " + pst.Reply.Root.Uri + "\n"
						}
						if pst.Facets != nil {
							for i := range pst.Facets {
								if pst.Facets[i].Features != nil {
									for j := range pst.Facets[i].Features {
										if pst.Facets[i].Features[j].RichtextFacet_Mention != nil {

										}
										if pst.Facets[i].Features[j].RichtextFacet_Link != nil {

										}
										if pst.Facets[i].Features[j].RichtextFacet_Tag != nil {

										}
									}
								}

							}
						}
						if pst.Langs != nil {

						}
						if pst.Tags != nil {

						}
						if pst.Embed != nil {
							if pst.Embed.EmbedRecord != nil {
								log.I.S(pst.Embed.EmbedRecord)
							}
							if pst.Embed.EmbedImages != nil {
								for _, img := range pst.Embed.EmbedImages.Images {
									log.I.S(img)
								}
							}
							if pst.Embed.EmbedRecordWithMedia != nil {
								if pst.Embed.EmbedRecordWithMedia.Record != nil {
									log.I.S(pst.Embed.EmbedRecordWithMedia.Record)
								}
								if pst.Embed.EmbedRecordWithMedia.Media != nil {
									if pst.Embed.EmbedRecordWithMedia.Media.EmbedImages != nil {
										for _, img := range pst.Embed.EmbedRecordWithMedia.Media.EmbedImages.Images {
											log.I.S(img)
										}
										if pst.Embed.EmbedRecordWithMedia.Media.EmbedExternal != nil {
											log.I.S(pst.Embed.EmbedRecordWithMedia.Media.EmbedExternal.External)
										}
									}
								}
								log.I.S(pst.Embed.EmbedRecordWithMedia.Media.EmbedImages.Images)

							}
							if pst.Embed.EmbedExternal != nil {
								log.I.S(pst.Embed.EmbedExternal.External)
							}
						}
						log.I.F("Post\nEvent ID: %s\nCreated At: %s\nUser: %s\nContent: [%d]\n`%s`\n%s",
							op.Cid, pst.CreatedAt, rr.SignedCommit().Did, len(pst.Text), pst.Text, reply)
						// log.I.S(pst)
					}
					return
				}
			default:
				log.I.Ln(ek)
			}
		}
		// cancel()
		_ = rr
		return
	}
}

func PrintPost(pst appbsky.FeedPost,
	userProfile, replyUserProfile, likingUserProfile *appbsky.ActorDefs_ProfileViewDetailed, postPath string) {
	if userProfile != nil && userProfile.FollowersCount != nil {
		// Try to use the display name and follower count if we can get it
		var rply, likedTxt string
		if pst.Reply != nil && replyUserProfile != nil && replyUserProfile.FollowersCount != nil {
			rply = " ➡️ " + replyUserProfile.Handle + ":" + strconv.Itoa(int(*userProfile.FollowersCount)) + "\n" // + "https://staging.bsky.app/profile/" + strings.Split(pst.Reply.Parent.Uri, "/")[2] + "/post/" + path.Base(pst.Reply.Parent.Uri) + "\n"
		} else if likingUserProfile != nil {
			likedTxt = likingUserProfile.Handle + ":" + strconv.Itoa(int(*likingUserProfile.FollowersCount)) + " ❤️ "
			rply = ":\n"
		} else {
			rply = ":\n"
		}

		url := "https://bsky.app/profile/" + userProfile.Handle + "/post/" + path.Base(postPath)
		fmtdstring := likedTxt + userProfile.Handle + ":" + strconv.Itoa(int(*userProfile.FollowersCount)) + rply + pst.Text + "\n" + url + "\n"
		log.I.Ln(fmtdstring)
	}
}

func RepoHandle() func(handle *comatproto.SyncSubscribeRepos_Handle) error {
	return func(handle *comatproto.SyncSubscribeRepos_Handle) error {
		b, err := json.MarshalIndent(handle, "", "  ")
		if err != nil {
			return err
		}
		log.I.F("RepoHandle:\n%s", b)
		return nil
	}
}

func RepoInfo() func(info *comatproto.SyncSubscribeRepos_Info) error {
	return func(info *comatproto.SyncSubscribeRepos_Info) error {
		b, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return err
		}
		log.I.F("RepoInfo\r%s", b)
		log.I.F("INFO: %s: %v\n", info.Name, info.Message)
		return nil
	}
}

func RepoMigrate() func(mig *comatproto.SyncSubscribeRepos_Migrate) error {
	return func(mig *comatproto.SyncSubscribeRepos_Migrate) error {
		b, err := json.Marshal(mig)
		if err != nil {
			return err
		}
		fmt.Println("RepoMigrate")
		fmt.Println(string(b))
		return nil
	}
}

func RepoTombstone() func(tomb *comatproto.SyncSubscribeRepos_Tombstone) error {
	return func(tomb *comatproto.SyncSubscribeRepos_Tombstone) error {
		b, err := json.Marshal(tomb)
		if err != nil {
			return err
		}
		fmt.Println("RepoTombstone")
		fmt.Println(string(b))
		return nil
	}
}

func LabelLabels() func(labels *comatproto.LabelSubscribeLabels_Labels) error {
	return func(labels *comatproto.LabelSubscribeLabels_Labels) error {
		b, err := json.Marshal(labels)
		if err != nil {
			return err
		}
		fmt.Println("LabelLabels")
		fmt.Println(string(b))
		return nil
	}
}

func LabelInfo() func(info *comatproto.LabelSubscribeLabels_Info) error {
	return func(info *comatproto.LabelSubscribeLabels_Info) error {
		b, err := json.Marshal(info)
		if err != nil {
			return err
		}
		fmt.Println("LabelInfo")
		fmt.Println(string(b))
		return nil
	}
}
