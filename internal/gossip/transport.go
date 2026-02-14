package gossip

import (
	"context"

	"github.com/michael-martinez-dev/adaptive-hive/pkg/types"
)

// Transport handles the network layer for gossip comms
type Transport interface {
	// TStart begins listening for incoming messages
	Start(ctx context.Context) error

	// Stop shuts down the transport gracefully
	Stop() error

	// SendTo sends a message to a specific address
	// Returns an error if the transport is not running
	SendTo(addr string, msg []byte) error

	// Messages returns a channel of received messages
	// The channel is closed when the transport stops
	Messages() <-chan *types.NetworkMessage

	// LocalAddr returns the address this transport is listening on
	LocalAddr() string
}
