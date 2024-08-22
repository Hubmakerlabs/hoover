package bluesky

import (
	"time"

	"github.com/Hubmakerlabs/hoover/pkg"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/repo"
	typegen "github.com/whyrusleeping/cbor-gen"
)

// {
//   "lexicon": 1,
//   "id": "app.bsky.graph.follow",
//   "defs": {
//     "main": {
//       "type": "record",
//       "description": "Record declaring a social 'follow' relationship of another account. Duplicate follows will be ignored by the AppView.",
//       "key": "tid",
//       "record": {
//         "type": "object",
//         "required": ["subject", "createdAt"],
//         "properties": {
//           "subject": { "type": "string", "format": "did" },
//           "createdAt": { "type": "string", "format": "datetime" }
//         }
//       }
//     }
//   }
// }

// FromBskyGraphFollow is
func FromBskyGraphFollow(
	evt *atproto.SyncSubscribeRepos_Commit,
	op *atproto.SyncSubscribeRepos_RepoOp,
	rr *repo.Repo,
	rec typegen.CBORMarshaler,
) (bundle *types.BundleItem, err error) {

	var createdAt time.Time
	var to any
	if to, createdAt, err = UnmarshalEvent(evt, rec, &bsky.GraphFollow{}); chk.E(err) {
		return
	}
	if to == nil {
		err = errorf.E("failed to unmarshal post")
		return
	}
	fol, ok := to.(*bsky.GraphFollow)
	if !ok {
		err = errorf.E("did not get", Kinds["follow"])
		return
	}
	bundle = &types.BundleItem{}
	bundle.Tags = GetCommon(rr, createdAt, op, evt)
	AppendTag(bundle, pkg.LikeEventId, fol.Subject)
	return
}

// ToBskyGraphFollow is
func ToBskyGraphFollow() {

}
