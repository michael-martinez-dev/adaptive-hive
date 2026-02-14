package gossip

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/michael-martinez-dev/adaptive-hive/pkg/types"
)

// --- Helpers --- //

func waitForMessage(t *testing.T, ch <-chan *types.NetworkMessage, timeout time.Duration) *types.NetworkMessage {
	t.Helper()
	select {
	case msg := <-ch:
		return msg
	case <-time.After(timeout):
		t.Fatal("timed out waiting for message")
		return nil
	}
}

func assertNoMessage(t *testing.T, ch <-chan *types.NetworkMessage, duration time.Duration) {
	t.Helper()
	select {
	case msg := <-ch:
		t.Fatalf("unexpected message received: %v", msg)
	case <-time.After(duration):
	}
}

func TestTestNetwork_SendReceive(t *testing.T) {
	network := NewTestNetwork()
	t1 := network.NewTransport("node1")
	t2 := network.NewTransport("node2")

	ctx := context.Background()
	if err := t1.Start(ctx); err != nil {
		t.Errorf("failed to start t1")
	}
	if err := t2.Start(ctx); err != nil {
		t.Errorf("failed to start t2")
	}

	payload := []byte("hello")
	if err := t1.SendTo("node2", payload); err != nil {
		t.Errorf("failed to send message to node2")
	}

	msg := waitForMessage(t, t2.Messages(), time.Second)

	if string(msg.Payload) != "hello" {
		t.Errorf("payload = %q, want %q", msg.Payload, "hello")
	}

	if msg.From != "node1" {
		t.Errorf("from = %q, want %q", msg.From, "node1")
	}
}

func TestTestNetwork_PartitionIsBidirectional(t *testing.T) {
	network := NewTestNetwork()
	t1 := network.NewTransport("node1")
	t2 := network.NewTransport("node2")

	ctx := context.Background()
	if err := t1.Start(ctx); err != nil {
		t.Errorf("failed to start t1")
	}
	if err := t2.Start(ctx); err != nil {
		t.Errorf("failed to start t2")
	}

	network.Partition("node1", "node2")

	payload := []byte("hello")

	if err := t1.SendTo("node2", payload); err != nil {
		t.Error("failed to send message to node2")
	}
	assertNoMessage(t, t2.Messages(), 100*time.Millisecond)

	if err := t2.SendTo("node1", payload); err != nil {
		t.Error("failed to send message to node1")
	}
	assertNoMessage(t, t1.Messages(), 100*time.Millisecond)
}

func TestTestNetwork_HealRestoresConnectivity(t *testing.T) {
	network := NewTestNetwork()
	t1 := network.NewTransport("node1")
	t2 := network.NewTransport("node2")

	ctx := context.Background()
	if err := t1.Start(ctx); err != nil {
		t.Errorf("failed to start t1")
	}
	if err := t2.Start(ctx); err != nil {
		t.Errorf("failed to start t2")
	}

	network.Partition("node1", "node2")
	network.Heal("node1", "node2")

	if err := t1.SendTo("node2", []byte("hello")); err != nil {
		t.Error("could not send from node1 to node2")
	}
	msg := waitForMessage(t, t2.Messages(), time.Second)

	if string(msg.Payload) != "hello" {
		t.Errorf("payload = %q, want %q", msg.Payload, "hello")
	}
}

func TestTestNetwork_SendBeforeStartFails(t *testing.T) {
	network := NewTestNetwork()
	t1 := network.NewTransport("node1")

	err := t1.SendTo("node2", []byte("hello"))
	if err != ErrTransportStopped {
		t.Errorf("err = %v, want ErrTransportStopped", err)
	}
}

func TestTestNetwork_SendAfterStopFails(t *testing.T) {
	network := NewTestNetwork()
	t1 := network.NewTransport("node1")

	ctx := context.Background()
	if err := t1.Start(ctx); err != nil {
		t.Error("could not start node1")
	}
	if err := t1.Stop(); err != nil {
		t.Error("could not stop node1")
	}

	err := t1.SendTo("node2", []byte("hello"))
	if err != ErrTransportStopped {
		t.Errorf("err = %v, want %q", t1.LocalAddr(), "node1")
	}
}

func TestTestNetwork_LocalAddr(t *testing.T) {
	network := NewTestNetwork()
	t1 := network.NewTransport("node1")

	if t1.LocalAddr() != "node1" {
		t.Errorf("LocalAddr() = %q, want %q", t1.LocalAddr(), "node1")
	}
}

func TestTestNetwork_DropRate(t *testing.T) {
	network := NewTestNetwork()
	t1 := network.NewTransport("node1")
	t2 := network.NewTransport("node2")

	ctx := context.Background()
	if err := t1.Start(ctx); err != nil {
		t.Errorf("failed to start t1")
	}
	if err := t2.Start(ctx); err != nil {
		t.Errorf("failed to start t2")
	}

	t2.SetDropRate(1.0)

	var err error
	for i := 0; i < 10; i++ {
		err = t1.SendTo("node2", []byte(fmt.Sprintf("hello%q", i)))
		if err != nil {
			t.Error("node1 could not send to node2")
		}
	}

	assertNoMessage(t, t2.Messages(), 100*time.Millisecond)
}

func TestTestNetwork_MessageToUnknownAddrIsDropped(t *testing.T) {
	network := NewTestNetwork()
	t1 := network.NewTransport("node1")

	ctx := context.Background()
	if err := t1.Start(ctx); err != nil {
		t.Error("node1 failed to start")
	}

	err := t1.SendTo("unknown", []byte("hello"))
	if err != nil {
		t.Error("send to unknown addr returned error: %V", err)
	}
}
