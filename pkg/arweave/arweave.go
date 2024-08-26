package ao

import "github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"

const Gateway = "https://arweave.net"

func AppendTag(bundle *types.BundleItem, name, value string) {
	bundle.Tags = append(bundle.Tags, types.Tag{Name: name, Value: value})
}
