package bluesky

import (
	"encoding/json"
	"fmt"

	"github.com/bluesky-social/indigo/api/atproto"
)

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
