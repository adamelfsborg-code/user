package server

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Nats struct {
	host string
}

func ConnectNats(n *Nats) (*nats.Conn, error) {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a channel to receive the NATS connection
	ncCh := make(chan *nats.Conn)

	go func() {
		// Establish a connection to the NATS server with nats.DefaultURL
		nc, err := nats.Connect(nats.DefaultURL)
		if err != nil {
			log.Fatal(err)
		}
		// Send the NATS connection via the channel
		ncCh <- nc
	}()

	select {
	case <-ctx.Done():
		// The context was canceled or timed out.
		// You can perform cleanup or take action accordingly.
	case nc := <-ncCh:
		// Return the NATS connection when it becomes available
		return nc, nil
	}

	return nil, errors.New("NATS connection not established")
}

func ConnectJetstream(n *nats.Conn) (nats.JetStreamContext, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ncc, _ := n.JetStream(nats.Context(ctx))

	return ncc, nil
}
