package main

import (
	"bytes"
	"os"

	"github.com/Hubmakerlabs/replicatr/pkg/slog"
)

type (
	B = []byte
	S = string
)

var (
	log, chk = slog.New(os.Stderr)
	equals   = bytes.Equal
)
