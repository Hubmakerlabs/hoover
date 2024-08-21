package arweave

import (
	"encoding/json"
	"fmt"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
)

func PrintBundleItem(bi *types.BundleItem) {
	if bi == nil {
		return
	}
	var err error
	for _, t := range bi.Tags {
		if len(t.Name) < 1 {
			continue
		}
		// tag lists have hash symbol prefixes
		if t.Name[0] == '#' {
			// present as a list of strings in a line like nostr tag
			fmt.Printf("\"%s\"", t.Name)
			var elements []S
			if err = json.Unmarshal(B(t.Value), &elements); err == nil {
				for i := range elements {
					fmt.Printf(",\"%s\"", elements[i])
					// if i < len(elements)-1 {
					// 	fmt.Print(",")
					// }
				}
			} else {
				fmt.Printf(",\"%s\"", t.Value)
			}
			fmt.Println()
		} else {
			// usually these are first anyway
			fmt.Printf("\"%s\",\"%s\"\n", t.Name, t.Value)
		}
	}
	// fmt.Printf("Data: %d\n", len(bi.Data))
	if len(bi.Data)>0{
		fmt.Printf("Data: %d\n```\n%s\n```\n", len(bi.Data), bi.Data)
	}
}
