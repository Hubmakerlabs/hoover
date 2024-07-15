package schema

import "github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"

type OptionItem struct {
	Target string
	Anchor string
	Tags   []types.Tag
}
