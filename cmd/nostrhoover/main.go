package main

import (
	"os"

	ao "github.com/Hubmakerlabs/hoover/pkg/arweave"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
)

func main() {
	var err error
	cl := goar.NewClient(ao.Gateway)
	var nodeInfo *types.NetworkInfo
	if nodeInfo, err = cl.GetInfo(); chk.E(err) {
		os.Exit(1)
	}
	log.I.S(nodeInfo)
}
