package bluesky

import (
	"encoding/json"

	"github.com/bluesky-social/indigo/api/atproto"
)

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
