package gossip

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/michael-martinez-dev/adaptive-hive/pkg/types"
)

var (
	ErrTransportStopped = errors.New("transport is stopped")
)

// TestNetwork simulates a network for testing gossip protocols
// It allows multiple TestTransports to communicate and supports
// fault injection like partitions, latency, and packet loss
type TestNetwork struct {
	mu         sync.RWMutex
	transports map[string]*TestTransport
	partitions map[string]map[string]bool // addr -> set of unreachable addrs
}

// NewTestNetwork creates a new sim network
func NewTestNetwork() *TestNetwork {
	return &TestNetwork{
		transports: make(map[string]*TestTransport),
		partitions: make(map[string]map[string]bool),
	}
}

// NewTransport creates a new TestTransport attached to this network
func (n *TestNetwork) NewTransport(addr string) *TestTransport {
	n.mu.Lock()
	defer n.mu.Unlock()

	t := &TestTransport{
		addr:    addr,
		msgCh:   make(chan *types.NetworkMessage, 100),
		network: n,
		running: false,
	}
	n.transports[addr] = t
	return t
}

// Partition makes 'from' unable to reach any of the 'unreachable' nodes
// Partitions are bidirectional by default
func (n *TestNetwork) Partition(from string, unreachable ...string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.partitions[from] == nil {
		n.partitions[from] = make(map[string]bool)
	}

	for _, addr := range unreachable {
		n.partitions[from][addr] = true

		if n.partitions[addr] == nil {
			n.partitions[addr] = make(map[string]bool)
		}
		n.partitions[addr][from] = true
	}
}

func (n *TestNetwork) Heal(a, b string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.partitions[a] != nil {
		delete(n.partitions[a], b)
	}

	if n.partitions[b] != nil {
		delete(n.partitions[b], a)
	}
}

// send delivers a message from one transport to another
// Returns nil even if the message is dropped (simulates UDP semantics)
// #nosec G404 -- test code
func (n *TestNetwork) send(from, to string, msg []byte) error {
	n.mu.RLock()
	defer n.mu.RUnlock()

	if parts, ok := n.partitions[from]; ok && parts[to] {
		return nil // partition between nodes - DROP
	}

	target, ok := n.transports[to]
	if !ok {
		return nil // target unknown - DROP
	}

	if !target.isRunning() {
		return nil // target not running - DROP
	}

	if target.latency > 0 {
		time.Sleep(target.latency) // sim latency
	}

	if target.dropRate > 0 && rand.Float64() < target.dropRate {
		return nil // simulated DROP
	}

	// Deliver MSG
	select {
	case target.msgCh <- &types.NetworkMessage{
		Payload: msg,
		From:    from,
		Time:    time.Now(),
	}:
	default:
		// channel full - DROP
	}
	return nil
}

// TestTransport implements Transport for testing
type TestTransport struct {
	addr     string
	msgCh    chan *types.NetworkMessage
	network  *TestNetwork
	dropRate float64 // 0.0-1.0 probability of dropping packets
	latency  time.Duration

	mu      sync.RWMutex
	running bool
}

func (t *TestTransport) SetDropRate(rate float64) {
	t.dropRate = rate
}

func (t *TestTransport) SetLatency(d time.Duration) {
	t.latency = d
}

func (t *TestTransport) Start(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.running = true
	return nil
}

func (t *TestTransport) Stop() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.running = false
	return nil
}

func (t *TestTransport) SendTo(addr string, msg []byte) error {
	if !t.isRunning() {
		return ErrTransportStopped
	}

	return t.network.send(t.addr, addr, msg)
}

func (t *TestTransport) Messages() <-chan *types.NetworkMessage {
	return t.msgCh
}

func (t *TestTransport) LocalAddr() string {
	return t.addr
}

func (t *TestTransport) isRunning() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.running
}
