package bluesky

import (
	"encoding/json"
	"fmt"

	"github.com/bluesky-social/indigo/api/atproto"
)

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
