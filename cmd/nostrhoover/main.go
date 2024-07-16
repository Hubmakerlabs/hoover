package main

import (
	"os"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
)

const arNode = "https://arweave.net"

func main() {
	var err error
	cl := goar.NewClient(arNode)
	var nodeInfo *types.NetworkInfo
	if nodeInfo, err = cl.GetInfo(); chk.E(err) {
		os.Exit(1)
	}
	log.I.S(nodeInfo)
}
