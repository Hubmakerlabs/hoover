package main

import (
	"context"
	"crypto/tls"
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
	endpoint = flag.String("endpoint", "hub-grpc.pinata.cloud", "The farcaster endpoint to connect to.")
	// port     = flag.Int("port", 2283, "The port to send gRPC requests to.")
)

func main() {
	flag.Parse()
	start := time.Now()

	// Set up the gRPC connection and client.
	creds := credentials.NewTLS(&tls.Config{})
	conn, err := grpc.NewClient(*endpoint, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to connect to %s: %v", *endpoint, err)
	}
	defer conn.Close()

	// Create the client.
	c := pb.NewHubServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Receive new messages as the hub merges them in.
	evts := []pb.HubEventType{pb.HubEventType_HUB_EVENT_TYPE_MERGE_MESSAGE}
	stream, err := c.Subscribe(ctx, &pb.SubscribeRequest{EventTypes: evts})
	if err != nil {
		log.Fatalf("failed to subscribe: %v", err)
	}
	defer stream.CloseSend()

	// Handle incoming messages.
	go func() {
		for {
			var msg pb.HubEvent
			if err := stream.RecvMsg(&msg); err != nil {
				log.Fatalf("failed to receive message: %v", err)
			}

			// If the message is a new cast, print out its text.
			data := msg.GetMergeMessageBody().GetMessage().GetData()
			if data.GetType() == pb.MessageType_MESSAGE_TYPE_CAST_ADD {
				castBody := data.GetCastAddBody()
				log.Println("New cast:", castBody.GetText())
			}
		}
	}()

	// Wait for a signal to shut down.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	fmt.Println("Shutting down. Ran for", time.Since(start))
}
