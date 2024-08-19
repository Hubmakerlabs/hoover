package bluesky

import (
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/repo"
	typegen "github.com/whyrusleeping/cbor-gen"
)

// {
//   "lexicon": 1,
//   "id": "app.bsky.graph.listitem",
//   "defs": {
//     "main": {
//       "type": "record",
//       "description": "Record representing an account's inclusion on a specific list. The AppView will ignore duplicate listitem records.",
//       "key": "tid",
//       "record": {
//         "type": "object",
//         "required": ["subject", "list", "createdAt"],
//         "properties": {
//           "subject": {
//             "type": "string",
//             "format": "did",
//             "description": "The account which is included on the list."
//           },
//           "list": {
//             "type": "string",
//             "format": "at-uri",
//             "description": "Reference (AT-URI) to the list record (app.bsky.graph.list)."
//           },
//           "createdAt": { "type": "string", "format": "datetime" }
//         }
//       }
//     }
//   }
// }

// FromBskyGraphListitem is
func FromBskyGraphListitem(
	evt *atproto.SyncSubscribeRepos_Commit,
	op *atproto.SyncSubscribeRepos_RepoOp,
	rr *repo.Repo,
	rec typegen.CBORMarshaler,
) (bundle *types.BundleItem, err error) {

	return
}

// ToBskyGraphListitem is
func ToBskyGraphListitem() {

}
