package types

import "time"

// NodeID uniquely identifies a node in a cluster
type NodeID string

// NodeState represents the membership state of a node
type NodeState int

const (
	StateAlive NodeState = iota
	StateSuspect
	StateDead
	StateLeft
)

// String returns a human-readable representation of NodeState
func (s NodeState) String() string {
	switch s {
	case StateAlive:
		return "alive"
	case StateSuspect:
		return "suspect"
	case StateDead:
		return "dead"
	case StateLeft:
		return "left"
	default:
		return "unknown"
	}
}

// Node represents a member of the cluster
type Node struct {
	ID          NodeID
	Address     string // host:port for gossip comms
	State       NodeState
	Incarnation uint64 // Lanport-like counter for state consistency
	Metadata    NodeMetadata
	LastUpdated time.Time
}

// NodeMetadata contains additional information about a node's capabilities
type NodeMetadata struct {
	Resources Resources
	Labels    map[string]string // for scheduling constraints
	Priority  int
}

// Resources describes the available resources on a node
type Resources struct {
	CPUMillis   int64
	MemoryBytes int64
	DiskBytes   int64
}
