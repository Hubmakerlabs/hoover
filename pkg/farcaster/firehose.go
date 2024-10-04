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
	Urls = []string{
		"hub.pinata.cloud",
		"api.hub.wevm.dev",
		"hoyt.farcaster.xyz:2283",
		"lamia.farcaster.xyz:2283",
		"api.farcasthub.com:2283",
		"nemes.farcaster.xyz:2283",
		"hub.farcaster.standardcrypto.vc:2281",
		"hoyt.farcaster.xyz:2281",
		"lamia.farcaster.xyz:2281",
		"api.farcasthub.com:2281",
		"nemes.farcaster.xyz:2281",
		"hub.farcaster.standardcrypto.vc:2282",
		"hoyt.farcaster.xyz:2282",
		"lamia.farcaster.xyz:2282",
		"api.farcasthub.com:2282",
		"nemes.farcaster.xyz:2282",
		"hub.farcaster.standardcrypto.vc:2283",
	}
	totalUrls = len(Urls)
	currUrl   = struct {
		curr int
		mu   sync.Mutex
	}{curr: 0}

	maxHashCount    = 5000 // Maximum number of hashes to keep in memory
	hashesOrder     = make([]string, 0, maxHashCount)
	hashesOrderLock sync.Mutex
	cancel_global   context.F
	firstSub        = false
)

// connectToHub establishes a connection to the specified Farcaster hub
func connectToHub(url string) (*grpc.ClientConn, pb.HubServiceClient, error) {
	creds := credentials.NewTLS(&tls.Config{})
	var conn *grpc.ClientConn
	var err error
	conn, err = grpc.NewClient(url, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to %s: %v", url, err)
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
func subscribeToHub(ctx context.T, url string, bundleStream chan<- *types.BundleItem,
	seenPosts *sync.Map, wg *sync.WaitGroup, endpoints []string) {
	defer wg.Done()
	wg.Add(1)
	if !firstSub {
		firstSub = true
	}
	conn, client, err := connectToHub(url)
	if err != nil {
		log.Printf("Failed to connect to hub %s - %v", url, err)
		replaceFailedConnection(ctx, bundleStream, seenPosts, wg, endpoints)
		return
	}
	defer conn.Close()

	// Subscribe to MERGE_MESSAGE events
	evts := []pb.HubEventType{pb.HubEventType_HUB_EVENT_TYPE_MERGE_MESSAGE}
	stream, err := client.Subscribe(ctx, &pb.SubscribeRequest{EventTypes: evts})
	if err != nil {
		log.Printf("Failed to subscribe to hub %s - %v", url, err)
		replaceFailedConnection(ctx, bundleStream, seenPosts, wg, endpoints)
		return
	}

	for {
		select {
		case <-ctx.Done():
			cancel_global()
			return
		default:
			msg, err := stream.Recv()
			if err != nil {
				log.Printf("Failed to receive message from hub %s - %v", url, err)
				replaceFailedConnection(ctx, bundleStream, seenPosts, wg, endpoints)
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
func replaceFailedConnection(ctx context.T, bundleStream chan<- *types.BundleItem,
	seenPosts *sync.Map, wg *sync.WaitGroup, endpoints []string) {
	select {
	case <-ctx.Done():
		cancel_global()
		return
	default:

		currUrl.mu.Lock()
		defer currUrl.mu.Unlock()
		go subscribeToHub(ctx, endpoints[currUrl.curr], bundleStream, seenPosts, wg, endpoints)
		currUrl.curr = (currUrl.curr + 1) % len(endpoints)

	}

}

// Firehose function connects to multiple hubs concurrently and streams BundleItems
func Firehose(
	ctx context.T,
	cancel context.F,
	wg_parent *sync.WaitGroup,
	endpoints []string,
	fn func(bundle *types.BundleItem) (err error),
) {
	wg_parent.Add(1)
	var ready bool
	var wg sync.WaitGroup
	seenPosts := &sync.Map{}
	bundleStream := make(chan *types.BundleItem)
	cancel_global = cancel

	// Start initial three connections
	for i := 0; i < 3; i++ {
		go replaceFailedConnection(ctx, bundleStream, seenPosts, &wg, endpoints)
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
			select {
			case <-ctx.Done():
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
