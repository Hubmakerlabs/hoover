package bluesky

import (
	"encoding/json"
	"fmt"

	"github.com/bluesky-social/indigo/api/atproto"
)

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
