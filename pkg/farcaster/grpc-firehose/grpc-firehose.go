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
	hubRpcEndpoints = []string{
		"hub.farcaster.standardcrypto.vc:2283",
		"hub-grpc.pinata.cloud",
		"api.farcasthub.com:2283",
		"nemes.farcaster.xyz:2283",
		"hub.pinata.cloud",
		"lamia.farcaster.xyz:2283",
		"hoyt.farcaster.xyz:2283",
		"api.hub.wevm.dev",
	}
	currentEndpointIndex = 0
	outputFilePath       = "output.jsonl"
)

func connectToHub() (*grpc.ClientConn, pb.HubServiceClient, error) {
	creds := credentials.NewTLS(&tls.Config{})
	endpoint := hubRpcEndpoints[currentEndpointIndex]

	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to %s: %v", endpoint, err)
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
		conn, client, err := connectToHub()
		if err != nil {
			log.Printf("%v", err)
			currentEndpointIndex++
			if currentEndpointIndex >= len(hubRpcEndpoints) {
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

				// Check if the message is a new cast
				if data := msg.GetMergeMessageBody().GetMessage().GetData(); data.GetType() == pb.MessageType_MESSAGE_TYPE_CAST_ADD {
					castBody := data.GetCastAddBody()

					// Write cast to JSONL file
					eventJson, err := json.Marshal(castBody)
					if err != nil {
						log.Printf("failed to marshal event: %v", err)
						continue
					}
					file.WriteString(string(eventJson) + "\n")

					// Print the cast to the console
					log.Println("New cast added:", castBody.GetText())
				}
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
