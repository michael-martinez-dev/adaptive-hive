package gossip

import (
	"math/rand"
	"slices"
	"sync"
	"time"

	"github.com/michael-martinez-dev/adaptive-hive/pkg/types"
)

type Membership interface {
	GetNode(id types.NodeID) (types.Node, bool)
	AllNodes() []types.Node
	GetNodesByState(state types.NodeState) []types.Node
	Len() int
	Merge(entry GossipEntry) bool
	Suspect(id types.NodeID) bool
	Dead(id types.NodeID) bool
	Remove(id types.NodeID)
	RandomNode(...types.NodeID) (types.Node, bool)
	RandomNodes(k int, exclude ...types.NodeID) []types.Node
}

type membership struct {
	mu    sync.RWMutex
	nodes map[string]*types.Node
}

func NewMembership() Membership {
	return &membership{
		nodes: map[string]*types.Node{},
	}
}

func (m *membership) GetNode(id types.NodeID) (types.Node, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	node, exists := m.nodes[string(id)]
	if !exists {
		return types.Node{}, false
	}
	return *node, exists
}

func (m *membership) AllNodes() []types.Node {
	return make([]types.Node, 0)
}

func (m *membership) GetNodesByState(state types.NodeState) []types.Node {
	return make([]types.Node, 0)
}

func (m *membership) Len() int {
	return len(m.nodes)
}

func (m *membership) Merge(entry GossipEntry) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	node, exists := m.nodes[string(entry.NodeID)]
	if !exists {
		node = &types.Node{
			ID:          entry.NodeID,
			Address:     entry.Address,
			State:       entry.State,
			Incarnation: entry.Incarnation,
			LastUpdated: time.Now(),
		}
		m.nodes[string(entry.NodeID)] = node
		return true
	}

	if entry.Incarnation > node.Incarnation {
		node.State = entry.State
		node.Incarnation = entry.Incarnation
		return true
	}

	if entry.Incarnation == node.Incarnation && entry.State > node.State {
		node.State = entry.State
		node.Incarnation = entry.Incarnation
		return true
	}

	return false
}

func (m *membership) Suspect(id types.NodeID) bool {
	return false
}

func (m *membership) Dead(id types.NodeID) bool {
	return false
}

func (m *membership) Remove(id types.NodeID) {
}

func (m *membership) RandomNode(exclude ...types.NodeID) (types.Node, bool) {
	rnode := m.RandomNodes(1, exclude...)
	if len(rnode) == 0 {
		return types.Node{}, false
	}
	return rnode[0], true
}

func (m *membership) RandomNodes(k int, exclude ...types.NodeID) []types.Node {
	if k < 1 || len(m.nodes) == 0 {
		return make([]types.Node, 0)
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	rnodes := make([]types.Node, 0)
	for _, val := range m.nodes {
		if isExcluded(exclude, val) {
			continue
		}
		if !isLiveState(val.State) {
			continue
		}

		rnodes = append(rnodes, *val)
	}
	rnodes = randomize(rnodes)
	if len(rnodes) > k {
		return rnodes[:k]
	}
	return rnodes
}

func isExcluded(excludedNodes []types.NodeID, node *types.Node) bool {
	return slices.Contains(excludedNodes, node.ID)
}

func isLiveState(state types.NodeState) bool {
	return state == types.StateAlive || state == types.StateSuspect
}

func randomize(nodes []types.Node) []types.Node {
	rand.Shuffle(len(nodes), func(i, j int) { nodes[i], nodes[j] = nodes[j], nodes[i] })
	return nodes
}
