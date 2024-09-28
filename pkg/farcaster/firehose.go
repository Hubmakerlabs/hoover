package farcaster

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"log"
	"sync"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/context"
	pb "github.com/juiceworks/hubble-grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	hubRpcEndpoints = []struct {
		url        string
		needs_port bool
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
	cancel_global   context.F
	firstSub        = false
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
func subscribeToHub(ctx context.T, hub struct {
	url        string
	needs_port bool
}, port string, bundleStream chan<- *types.BundleItem, seenPosts *sync.Map, sem chan struct{}, wg *sync.WaitGroup, remainingHubs *sync.Map) {
	defer wg.Done()
	wg.Add(1)
	defer func() { <-sem }()
	if !firstSub {
		firstSub = true
	}
	conn, client, err := connectToHub(hub.url, port, hub.needs_port)
	if err != nil {
		log.Printf("Failed to connect to hub %s:%s - %v", hub.url, port, err)
		replaceFailedConnection(ctx, bundleStream, seenPosts, sem, wg, remainingHubs)
		return
	}
	defer conn.Close()

	// Subscribe to MERGE_MESSAGE events
	evts := []pb.HubEventType{pb.HubEventType_HUB_EVENT_TYPE_MERGE_MESSAGE}
	stream, err := client.Subscribe(ctx, &pb.SubscribeRequest{EventTypes: evts})
	if err != nil {
		log.Printf("Failed to subscribe to hub %s:%s - %v", hub.url, port, err)
		replaceFailedConnection(ctx, bundleStream, seenPosts, sem, wg, remainingHubs)
		return
	}

	for {
		select {
		case <-ctx.Done():
			close(bundleStream)
			close(sem)
			cancel_global()
			return
		default:
			msg, err := stream.Recv()
			if err != nil {
				log.Printf("Failed to receive message from hub %s:%s - %v", hub.url, port, err)
				replaceFailedConnection(ctx, bundleStream, seenPosts, sem, wg, remainingHubs)
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
}

// replaceFailedConnection replaces a failed connection with a new one from the remaining pool
func replaceFailedConnection(ctx context.T, bundleStream chan<- *types.BundleItem, seenPosts *sync.Map, sem chan struct{}, wg *sync.WaitGroup, remainingHubs *sync.Map) {
	select {
	case <-ctx.Done():
		close(bundleStream)
		close(sem)
		cancel_global()
		return
	default:
		var remainingHubsEmpty bool
		remainingHubs.Range(func(key, value interface{}) bool {
			remainingHubsEmpty = false
			return false
		})

		if remainingHubsEmpty {
			for _, hub := range hubRpcEndpoints {
				for _, port := range ports {
					if hub.needs_port {
						remainingHubs.Store(hub, port)
					} else if _, ok := remainingHubs.Load(hub); !ok {
						remainingHubs.Store(hub, port)
					}
				}
			}
		}
		var found bool

		remainingHubs.Range(func(key any, value interface{}) bool {
			hub := key.(struct {
				url        string
				needs_port bool
			})
			port := value.(string)

			go subscribeToHub(ctx, hub, port, bundleStream, seenPosts, sem, wg, remainingHubs)

			remainingHubs.Delete(key)
			found = true
			return false

		})
		if !found {
			<-sem
		}
	}

}

// Firehose function connects to multiple hubs concurrently and streams BundleItems
func Firehose(ctx context.T, cancel context.F, wg_parent *sync.WaitGroup,
	fn func(bundle *types.BundleItem) (err error)) {
	wg_parent.Add(1)
	var ready bool
	var wg sync.WaitGroup
	seenPosts := &sync.Map{}
	sem := make(chan struct{}, 3)
	remainingHubs := &sync.Map{}
	bundleStream := make(chan *types.BundleItem)
	cancel_global = cancel

	for _, hub := range hubRpcEndpoints {
		for _, port := range ports {
			if hub.needs_port {
				remainingHubs.Store(hub, port)
			} else if _, ok := remainingHubs.Load(hub); !ok {
				remainingHubs.Store(hub, port)
			}
		}
	}

	// Start initial three connections
	for i := 0; i < 3; i++ {
		sem <- struct{}{}
		go replaceFailedConnection(ctx, bundleStream, seenPosts, sem, &wg, remainingHubs)
	}

	// Close the bundleStream when all subscriptions are done
	for {
		if firstSub {
			go func() {
				wg.Wait()
				close(bundleStream)
				cancel()
			}()
			break
		}
	}

	for bundle := range bundleStream {
		if !ready {
			ready = true
			wg_parent.Done()
		}
		wg_parent.Wait()
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovered from panic in fn: %v", r)
				}
			}()
			select {
			case <-ctx.Done():
				close(bundleStream)
				close(sem)
				cancel()
				return
			default:
				if fn == nil {
					fmt.Println("No function provided. Printing bundle:")
					fmt.Println(bundle)
				} else if bundle == nil {
					log.Printf("Bundle is nil")
				} else if err := fn(bundle); err != nil {
					log.Printf("Error processing bundle: %v", err)
				}
			}
		}()
	}

}
