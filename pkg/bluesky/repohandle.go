package bluesky

import (
	"encoding/json"

	"github.com/bluesky-social/indigo/api/atproto"
)

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
