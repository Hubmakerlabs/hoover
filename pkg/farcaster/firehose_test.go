package farcaster

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	ao "github.com/Hubmakerlabs/hoover/pkg/arweave"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
)

var (
	outputFilePath = "bundled.jsonl"
)

func TestFirehose(t *testing.T) {

	// Create file stream for writing JSONL output
	file, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bundleStream := make(chan *types.BundleItem)

	// Start the Firehose to stream BundleItems
	go Firehose(ctx, bundleStream)

	// Consume the stream and print to console and save to bundled.jsonl
	for bundle := range bundleStream {
		eventJson, err := json.Marshal(bundle)
		if err != nil {
			log.Printf("Failed to marshal bundle: %v", err)
		}

		file.WriteString(string(eventJson) + "\n")
		ao.PrintBundleItem(bundle)
	}

	// Simulate running for 10 seconds and then stop
	time.Sleep(30 * time.Second)
	cancel() // Cancel the context to stop the Firehose
}
