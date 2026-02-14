package gossip

import (
	"sync"
	"testing"

	"github.com/michael-martinez-dev/adaptive-hive/pkg/types"
)

func TestMembership_MergeHigherIncarnationWins(t *testing.T) {
	var (
		entry   GossipEntry
		changed bool
		node    types.Node
		exists  bool
	)

	// Setup: node exists at incarnation 5, state alive
	membership := NewMembership()

	entry = GossipEntry{
		NodeID:      "node1",
		Address:     "10.0.0.1:1234",
		State:       types.StateAlive,
		Incarnation: 5,
	}
	changed = membership.Merge(entry)
	if !changed {
		t.Error("initial merge unsuccessful")
	}

	node, exists = membership.GetNode("node1")
	if !exists {
		t.Error("node1 does not exist as expected")
	}
	if node.Incarnation != 5 {
		t.Errorf("incarnation = %d, want 5", node.Incarnation)
	}

	// Action: merge entry with incarnation 6, state alive
	entry = GossipEntry{
		NodeID:      "node1",
		Address:     "10.0.0.1:1234",
		State:       types.StateAlive,
		Incarnation: 6,
	}
	changed = membership.Merge(entry)
	if !changed {
		t.Error("update merge unsuccessful")
	}

	// Assert: node updated to incarnation 6
	node, exists = membership.GetNode("node1")
	if !exists {
		t.Error("node1 does not exist as expected")
	}
	if node.Incarnation != 6 {
		t.Errorf("incarnation = %d, want 6", node.Incarnation)
	}
}

func TestMembership_MergeSameIncarnationWorseStateWins(t *testing.T) {
	var (
		entry   GossipEntry
		changed bool
		node    types.Node
		exists  bool
	)

	// Setup: node exists at incarnation 5, state alive
	membership := NewMembership()

	entry = GossipEntry{
		NodeID:      "node1",
		Address:     "10.0.0.1:1234",
		State:       types.StateAlive,
		Incarnation: 5,
	}
	changed = membership.Merge(entry)
	if !changed {
		t.Error("initial merge unsuccessful")
	}

	node, exists = membership.GetNode("node1")
	if !exists {
		t.Error("node1 does not exist as expected")
	}
	if node.Incarnation != 5 {
		t.Errorf("incarnation = %d, want 5", node.Incarnation)
	}

	// Action: merge entry with incarnation 5, state suspect
	entry = GossipEntry{
		NodeID:      "node1",
		Address:     "10.0.0.1:1234",
		State:       types.StateSuspect,
		Incarnation: 5,
	}
	changed = membership.Merge(entry)
	if !changed {
		t.Error("update merge unsuccessful")
	}

	// Assert: node updated to suspect (same incarnation)
	node, exists = membership.GetNode("node1")
	if !exists {
		t.Error("node1 does not exist as expected")
	}
	if node.Incarnation != 5 {
		t.Errorf("incarnation = %d, want 5", node.Incarnation)
	}

	if node.State != types.StateSuspect {
		t.Errorf("state = %s, want suspect", node.State)
	}
}

func TestMembership_MergeSameIncarnationBetterStateIgnored(t *testing.T) {
	var (
		entry   GossipEntry
		changed bool
		node    types.Node
		exists  bool
	)

	// Setup: node exists at incarnaiton 5, state suspect
	membership := NewMembership()

	entry = GossipEntry{
		NodeID:      "node1",
		Address:     "10.0.0.1:1234",
		State:       types.StateSuspect,
		Incarnation: 5,
	}
	changed = membership.Merge(entry)
	if !changed {
		t.Error("initial merge unsuccessful")
	}

	node, exists = membership.GetNode("node1")
	if !exists {
		t.Error("node1 does not exist as expected")
	}
	if node.Incarnation != 5 {
		t.Errorf("incarnation = %d, want 5", node.Incarnation)
	}

	// Action: merge entry with incarnation 5, state alive
	entry = GossipEntry{
		NodeID:      "node1",
		Address:     "10.0.0.1:1234",
		State:       types.StateAlive,
		Incarnation: 5,
	}
	changed = membership.Merge(entry)
	if changed {
		t.Error("update should not have merge")
	}

	node, exists = membership.GetNode("node1")
	if !exists {
		t.Error("node1 does not exist as expected")
	}
	// Assert: node updated to suspect (same incarnation)
	if node.Incarnation != 5 {
		t.Errorf("incarnation = %d, want 5", node.Incarnation)
	}

	// Assert: node still suspect
	if node.State != types.StateSuspect {
		t.Errorf("state = %s, want suspect", node.State)
	}
}

func TestMembership_MergeLowerIncarnationIgnored(t *testing.T) {
	var (
		entry   GossipEntry
		changed bool
		node    types.Node
		exists  bool
	)

	// Setup: node exists at incarnation 5, state alive
	membership := NewMembership()

	entry = GossipEntry{
		NodeID:      "node1",
		Address:     "10.0.0.1:1234",
		State:       types.StateAlive,
		Incarnation: 5,
	}
	changed = membership.Merge(entry)
	if !changed {
		t.Error("update merge unsuccessful")
	}

	node, exists = membership.GetNode("node1")
	if !exists {
		t.Error("node1 does not exist as expected")
	}

	// Action: merge entry with incarnation 4, state dead
	entry = GossipEntry{
		NodeID:      "node1",
		Address:     "10.0.0.1:1234",
		State:       types.StateDead,
		Incarnation: 4,
	}
	changed = membership.Merge(entry)
	if changed {
		t.Error("update should not have merge")
	}

	// Assert: node still alive at incarnation 5 (stale update ignored)
	node, exists = membership.GetNode("node1")
	if !exists {
		t.Error("node1 does not exist as expected")
	}

	if node.Incarnation != 5 {
		t.Errorf("incarnation = %d, want 5", node.Incarnation)
	}

	if node.State != types.StateAlive {
		t.Errorf("state = %s, want alive", node.State)
	}
}

func TestMembership_ConcurrentAccess(t *testing.T) {
	membership := NewMembership()

	// Send with node
	membership.Merge(GossipEntry{
		NodeID:      "node1",
		Address:     "10.0.0.1:1234",
		State:       types.StateAlive,
		Incarnation: 1,
	})

	var wg sync.WaitGroup

	// Concurrent readers
	var i uint64
	for i = range 10 {
		wg.Add(1)
		go func(inc uint64) {
			defer wg.Done()
			for range 100 {
				membership.Merge(GossipEntry{
					NodeID:      "node1",
					Address:     "10.0.0.1:1234",
					State:       types.StateAlive,
					Incarnation: inc,
				})
			}
		}(i)
	}

	wg.Wait()
}

func TestMembership_RandomNodeExcludesDeadAndLeft(t *testing.T) {
	membership := NewMembership()

	// Add one alive, one dead node
	_ = membership.Merge(GossipEntry{NodeID: "alive", State: types.StateAlive, Incarnation: 1})
	_ = membership.Merge(GossipEntry{NodeID: "dead", State: types.StateDead, Incarnation: 1})
	_ = membership.Merge(GossipEntry{NodeID: "suspect", State: types.StateSuspect, Incarnation: 1})
	_ = membership.Merge(GossipEntry{NodeID: "left", State: types.StateLeft, Incarnation: 1})

	for range 100 {
		node, ok := membership.RandomNode()
		if !ok {
			t.Fatal("RandomNode returned not ok")
		}
		if node.State == types.StateDead || node.State == types.StateLeft {
			t.Fatal("RandomNode unexpectedly returned node in state dead/left")
		}
	}
}

func TestMembership_RandomNodesExcludesDeadAndLeft(t *testing.T) {
	membership := NewMembership()

	_ = membership.Merge(GossipEntry{NodeID: "alive", State: types.StateAlive, Incarnation: 1})
	_ = membership.Merge(GossipEntry{NodeID: "dead", State: types.StateDead, Incarnation: 1})
	_ = membership.Merge(GossipEntry{NodeID: "suspect", State: types.StateSuspect, Incarnation: 1})
	_ = membership.Merge(GossipEntry{NodeID: "left", State: types.StateLeft, Incarnation: 1})

	for range 100 {
		nodes := membership.RandomNodes(4)
		if len(nodes) != 2 {
			t.Fatalf("RandomNodes returned %d nodes, expected 3", len(nodes))
		}
		for _, n := range nodes {
			if n.State == types.StateDead || n.State == types.StateLeft {
				t.Fatalf("RandomNodes unexpectedly return %s node", n.State)
			}
		}
	}
}

func TestMembership_RandomNodeExcludesSpecified(t *testing.T) {
	membership := NewMembership()

	_ = membership.Merge(GossipEntry{NodeID: "node1", State: types.StateAlive, Incarnation: 1})
	_ = membership.Merge(GossipEntry{NodeID: "node2", State: types.StateAlive, Incarnation: 1})
	_ = membership.Merge(GossipEntry{NodeID: "node3", State: types.StateAlive, Incarnation: 1})
	_ = membership.Merge(GossipEntry{NodeID: "node4", State: types.StateAlive, Incarnation: 1})

	var excNode types.NodeID = "node2"

	for range 100 {
		nodes := membership.RandomNodes(3, excNode)
		if len(nodes) != 3 {
			t.Fatalf("RandomNodes returned %d nodes, expected 3", len(nodes))
		}
		for _, n := range nodes {
			if n.ID == excNode {
				t.Fatal("RandomNodes returned excluded node")
			}
		}
	}
}
