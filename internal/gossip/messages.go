package gossip

import "github.com/michael-martinez-dev/adaptive-hive/pkg/types"

// MessageType identifies the kind of gossip protocol message
type MessageType uint8

const (
	MessageTypePing MessageType = iota + 1
	MessageTypePingReq
	MessageTypeAck
	MessageTypeNack
	MessageTypeSync
	MessageTypeSyncResponse
	MessageTypeLeave
)

// String returns a human-readable representation of MessageType
func (m MessageType) String() string {
	switch m {
	case MessageTypePing:
		return "ping"
	case MessageTypePingReq:
		return "ping-req"
	case MessageTypeAck:
		return "ack"
	case MessageTypeNack:
		return "nack"
	case MessageTypeSync:
		return "sync"
	case MessageTypeSyncResponse:
		return "sync-response"
	case MessageTypeLeave:
		return "leave"
	default:
		return "unknown"
	}
}

// MessageHeader is embedded in all gossip protocol message
type MessageHeader struct {
	Version  uint8
	Type     MessageType
	SeqNo    uint32 // For correlating requests with responses
	SourceID string // NodeID of the sender
}

// Ping checks if a target node is alive
type Ping struct {
	MessageHeader
	Target string        // NodeID of who to ping
	Gossip []GossipEntry // Piggybacked dissemination
}

// Ack responds to a successful ping
type Ack struct {
	MessageHeader
	Gossip []GossipEntry
}

// Nack indicates a failed indirect ping
type Nack struct {
	MessageHeader
}

// PingReq asks another node to ping a target on our behalf (indirect ping)
type PingReq struct {
	Ping
	TargetAdder string // Address to reach to the target
}

// GossipEntry is a single piece of membership information being disseminated
type GossipEntry struct {
	NodeID      types.NodeID
	Address     string
	State       types.NodeState
	Incarnation uint64
	Metadata    *types.NodeMetadata // nil if metadata unchanged
	Timestamp   int64               // Unix nanos
}
