package bluesky

import (
	"encoding/hex"
	"strconv"
	"time"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/repo"
	"github.com/bluesky-social/indigo/util"
	typegen "github.com/whyrusleeping/cbor-gen"
)

/*
{
  "lexicon": 1,
  "id": "app.bsky.graph.block",
  "defs": {
    "main": {
      "type": "record",
      "description": "Record declaring a 'block' relationship against another account. NOTE: blocks are public in Bluesky; see blog posts for details.",
      "key": "tid",
      "record": {
        "type": "object",
        "required": ["subject", "createdAt"],
        "properties": {
          "subject": {
            "type": "string",
            "format": "did",
            "description": "DID of the account to be blocked."
          },
          "createdAt": { "type": "string", "format": "datetime" }
        }
      }
    }
  }
}
*/

// FromBskyGraphBlock is
func FromBskyGraphBlock(
	evt *atproto.SyncSubscribeRepos_Commit,
	op *atproto.SyncSubscribeRepos_RepoOp,
	rr *repo.Repo,
	rec typegen.CBORMarshaler,
) (bundle *types.BundleItem, err error) {

	var createdAt time.Time
	var to any
	if to, createdAt, err = UnmarshalEvent(evt, rec, &bsky.GraphBlock{}); chk.E(err) {
		return
	}
	if to == nil {
		err = errorf.E("failed to unmarshal post")
		return
	}
	repost, ok := to.(*bsky.GraphBlock)
	if !ok {
		err = errorf.E("did not get", Kinds["block"])
		return
	}
	bundle = &types.BundleItem{}
	bundle.Tags = []types.Tag{
		{Name: "protocol", Value: "bsky"},
		{Name: "kind", Value: Kinds["block"]},
		{Name: "id", Value: op.Cid.String()},
		{Name: "pubkey", Value: rr.SignedCommit().Did},
		{Name: "created_at", Value: strconv.FormatInt(createdAt.Unix(), 10)},
		{Name: "repo", Value: evt.Repo},
		{Name: "path", Value: op.Path},
		{Name: "sig", Value: hex.EncodeToString(rr.SignedCommit().Sig)},
	}
	if createdAt, err = time.Parse(util.ISO8601, repost.CreatedAt); chk.E(err) {
		return
	}
	AppendTag(bundle, "#sent_timestamp", strconv.FormatInt(createdAt.Unix(), 10))
	AppendTag(bundle, "#subject", repost.Subject)
	return
}

// ToBskyGraphBlock is
func ToBskyGraphBlock() {

}
