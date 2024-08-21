package bluesky

import (
	"encoding/json"
	"fmt"

	"github.com/bluesky-social/indigo/api/atproto"
)

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
