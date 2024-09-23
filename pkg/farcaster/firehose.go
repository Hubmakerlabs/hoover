package farcaster

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"sync"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	pb "github.com/juiceworks/hubble-grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	hubRpcEndpoints = []struct {
		url  string
		port bool
	}{
		{"hub.farcaster.standardcrypto.vc", true},
		{"hub.pinata.cloud", false},
		{"hoyt.farcaster.xyz", true},
		{"lamia.farcaster.xyz", true},
		{"api.farcasthub.com", true},
		{"nemes.farcaster.xyz", true},
		{"api.hub.wevm.dev", false},
	}
	ports = []string{"2281", "2282", "2283"}
)

// connectToHub establishes a connection to the specified Farcaster hub
func connectToHub(url string, port string, isPort bool) (*grpc.ClientConn, pb.HubServiceClient, error) {
	creds := credentials.NewTLS(&tls.Config{})
	var conn *grpc.ClientConn
	var err error
	if isPort {

		conn, err = grpc.NewClient(fmt.Sprintf("%s:%s", url, port), grpc.WithTransportCredentials(creds))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to connect to %s:%s: %v", url, port, err)
		}
	} else {

		conn, err = grpc.NewClient(url, grpc.WithTransportCredentials(creds))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to connect to %s: %v", url, err)
		}
	}

	client := pb.NewHubServiceClient(conn)
	return conn, client, nil
}

// subscribeToHub listens for messages from the given hub and sends them to the provided channel
func subscribeToHub(ctx context.Context, hub struct {
	url  string
	port bool
}, port string, bundleStream chan<- *types.BundleItem, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, client, err := connectToHub(hub.url, port, hub.port)
	if err != nil {
		log.Printf("Failed to connect to hub %s:%s - %v", hub.url, port, err)
		return
	}
	defer conn.Close()

	// Subscribe to MERGE_MESSAGE events
	evts := []pb.HubEventType{pb.HubEventType_HUB_EVENT_TYPE_MERGE_MESSAGE}
	stream, err := client.Subscribe(ctx, &pb.SubscribeRequest{EventTypes: evts})
	if err != nil {
		log.Printf("Failed to subscribe to hub %s:%s - %v", hub.url, port, err)
		return
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Printf("Failed to receive message from hub %s:%s - %v", hub.url, port, err)
			return
		}
		message := msg.GetMergeMessageBody().GetMessage()
		bundle, err := MessageToBundleItem(message)
		if err != nil {
			log.Printf("Failed to convert message to bundle item: %v", err)
			continue
		}
		bundleStream <- bundle
	}
}

// Firehose function connects to multiple hubs concurrently and streams BundleItems
func Firehose(ctx context.Context, bundleStream chan<- *types.BundleItem) error {
	var wg sync.WaitGroup

	// Start a Goroutine for each hub and port combination
	for _, hub := range hubRpcEndpoints {
		for _, port := range ports {
			wg.Add(1)
			go subscribeToHub(ctx, hub, port, bundleStream, &wg)
		}
	}

	// Wait for all Goroutines to complete
	go func() {
		wg.Wait()
		close(bundleStream)
	}()

	return nil
}
