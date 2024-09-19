package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	ports                = []string{"2281", "2282", "2283"}
	currentEndpointIndex = 0
	currentPortIndex     = 0
	outputFilePath       = "output.jsonl"
)

func connectToHub(url string, port string, is_port bool) (*grpc.ClientConn, pb.HubServiceClient, error) {
	creds := credentials.NewTLS(&tls.Config{})
	var conn *grpc.ClientConn
	var err error
	if is_port {
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

func main() {
	start := time.Now()
	flag.Parse()

	// Create file stream for writing JSONL output
	file, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	for {
		hub := hubRpcEndpoints[currentEndpointIndex]
		conn, client, err := connectToHub(hub.url, ports[currentPortIndex], hub.port)
		if err != nil {
			log.Printf("%v", err)
			if hub.port && currentPortIndex < len(ports) {
				currentPortIndex++
			} else if currentEndpointIndex < len(hubRpcEndpoints) {
				currentEndpointIndex++
				currentPortIndex = 0
			} else {
				log.Println("All connection attempts failed. Exiting.")
				return
			}
			continue
		}

		defer conn.Close()

		// Context and signal handling
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Subscribe to MERGE_MESSAGE events
		evts := []pb.HubEventType{pb.HubEventType_HUB_EVENT_TYPE_MERGE_MESSAGE}
		stream, err := client.Subscribe(ctx, &pb.SubscribeRequest{EventTypes: evts})
		if err != nil {
			log.Printf("failed to subscribe: %v", err)
			if hub.port && currentPortIndex < len(ports) {
				currentPortIndex++
			} else if currentEndpointIndex < len(hubRpcEndpoints) {
				currentEndpointIndex++
				currentPortIndex = 0
			} else {
				log.Println("All connection attempts failed. Exiting.")
				return
			}
			continue
		}

		// Handle incoming messages
		go func() {
			for {
				msg, err := stream.Recv()
				if err != nil {
					log.Printf("failed to receive message: %v", err)
					return
				}
				data := msg.GetMergeMessageBody().GetMessage().GetData()
				// Write cast to JSONL file
				eventJson, err := json.Marshal(data)
				if err != nil {
					log.Printf("failed to marshal event: %v", err)
					continue
				}

				file.WriteString(string(eventJson) + "\n")
				log.Println("New", data.GetType(), "added:", data.Body)
			}
		}()

		// Wait for a signal to shut down
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc
		fmt.Println("Shutting down. Ran for", time.Since(start))
		return
	}
}
