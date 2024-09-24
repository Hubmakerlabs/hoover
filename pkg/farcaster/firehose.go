package farcaster

import (
	"context"
	"crypto/tls"
	"encoding/hex"
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
	ports           = []string{"2281", "2282", "2283"}
	maxHashCount    = 5000 // Maximum number of hashes to keep in memory
	hashesOrder     = make([]string, 0, maxHashCount)
	hashesOrderLock sync.Mutex
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

// manageHashCapacity ensures that the hash storage doesn't exceed the specified limit
func manageHashCapacity(hash string, seenPosts *sync.Map) {
	hashesOrderLock.Lock()
	defer hashesOrderLock.Unlock()

	if len(hashesOrder) >= maxHashCount {
		oldestHash := hashesOrder[0]
		hashesOrder = hashesOrder[1:] // Remove the oldest hash
		seenPosts.Delete(oldestHash)
	}
	hashesOrder = append(hashesOrder, hash)
}

// subscribeToHub listens for messages from the given hub and sends them to the provided channel
func subscribeToHub(ctx context.Context, hub struct {
	url  string
	port bool
}, port string, bundleStream chan<- *types.BundleItem, seenPosts *sync.Map, sem chan struct{}, wg *sync.WaitGroup, remainingHubs *sync.Map, connLock *sync.Mutex) {
	defer wg.Done()
	defer func() { <-sem }()

	conn, client, err := connectToHub(hub.url, port, hub.port)
	if err != nil {
		log.Printf("Failed to connect to hub %s:%s - %v", hub.url, port, err)
		replaceFailedConnection(ctx, bundleStream, seenPosts, sem, wg, remainingHubs, connLock)
		return
	}
	defer conn.Close()

	// Subscribe to MERGE_MESSAGE events
	evts := []pb.HubEventType{pb.HubEventType_HUB_EVENT_TYPE_MERGE_MESSAGE}
	stream, err := client.Subscribe(ctx, &pb.SubscribeRequest{EventTypes: evts})
	if err != nil {
		log.Printf("Failed to subscribe to hub %s:%s - %v", hub.url, port, err)
		replaceFailedConnection(ctx, bundleStream, seenPosts, sem, wg, remainingHubs, connLock)
		return
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Printf("Failed to receive message from hub %s:%s - %v", hub.url, port, err)
			replaceFailedConnection(ctx, bundleStream, seenPosts, sem, wg, remainingHubs, connLock)
			return
		}
		message := msg.GetMergeMessageBody().GetMessage()
		hash := hex.EncodeToString(message.GetHash())

		// Check if the post has already been seen
		if _, loaded := seenPosts.LoadOrStore(hash, true); loaded {
			continue
		}

		// Manage the capacity of the seenPosts map
		manageHashCapacity(hash, seenPosts)

		bundle, err := MessageToBundleItem(message)
		if err != nil {
			log.Printf("Failed to convert message to bundle item: %v", err)
			continue
		}
		bundleStream <- bundle
	}
}

// replaceFailedConnection replaces a failed connection with a new one from the remaining pool
func replaceFailedConnection(ctx context.Context, bundleStream chan<- *types.BundleItem, seenPosts *sync.Map, sem chan struct{}, wg *sync.WaitGroup, remainingHubs *sync.Map, connLock *sync.Mutex) {
	connLock.Lock()
	defer connLock.Unlock()

	// Find next available hub and port combination
	var found bool
	remainingHubs.Range(func(key, value interface{}) bool {
		hub := key.(string)
		port := value.(string)

		// Start a new subscription
		wg.Add(1)
		go subscribeToHub(ctx, struct {
			url  string
			port bool
		}{url: hub, port: true}, port, bundleStream, seenPosts, sem, wg, remainingHubs, connLock)

		remainingHubs.Delete(key)
		found = true
		return false
	})

	if !found {
		<-sem
	}
}

// Firehose function connects to multiple hubs concurrently and streams BundleItems
func Firehose(ctx context.Context, bundleStream chan<- *types.BundleItem) error {
	var wg sync.WaitGroup
	seenPosts := &sync.Map{}
	sem := make(chan struct{}, 3)
	remainingHubs := &sync.Map{}
	connLock := &sync.Mutex{}

	for _, hub := range hubRpcEndpoints {
		for _, port := range ports {
			remainingHubs.Store(fmt.Sprintf(hub.url), port)
		}
	}

	// Start initial three connections
	for i := 0; i < 3; i++ {
		sem <- struct{}{}
		wg.Add(1)
		go replaceFailedConnection(ctx, bundleStream, seenPosts, sem, &wg, remainingHubs, connLock)
	}

	go func() {
		wg.Wait()
		close(bundleStream)
	}()

	return nil
}