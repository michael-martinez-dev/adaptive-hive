package gossip

import (
	"testing"

	"github.com/michael-martinez-dev/adaptive-hive/pkg/types"
)

func TestCodec_RoundtripPing(t *testing.T) {
	// Create a Ping, encode it, decode it, compare
	codec := NewCodec()

	original := &Ping{
		MessageHeader: MessageHeader{
			Version:  1,
			Type:     MessageTypePing,
			SeqNo:    42,
			SourceID: "node1",
		},
		Gossip: []GossipEntry{
			{
				NodeID:      types.NodeID("node2"),
				Address:     "10.0.0.1:1234",
				State:       types.StateAlive,
				Incarnation: 5,
			},
		},
	}

	data, err := codec.Encode(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	decoded, err := codec.Decode(data)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	ping, ok := decoded.(*Ping)
	if !ok {
		t.Fatalf("decoded type = %T, want *Ping", decoded)
	}

	if ping.SeqNo != original.SeqNo {
		t.Errorf("SeqNo = %d, want %d", ping.SeqNo, original.SeqNo)
	}

	if ping.SourceID != original.SourceID {
		t.Errorf("SourceID = %s, want %s", ping.SourceID, original.SourceID)
	}

	if len(ping.Gossip) != len(original.Gossip) {
		t.Errorf("Gossip length = %d, want %d", len(ping.Gossip), len(original.Gossip))
	}
}

func TestCodec_RoundtripAck(t *testing.T) {
	// Create an Ack, encode it, decode it, compare
	codec := NewCodec()

	original := &Ack{
		MessageHeader: MessageHeader{
			Version:  1,
			Type:     MessageTypeAck,
			SeqNo:    42,
			SourceID: "node1",
		},
		Gossip: []GossipEntry{
			{
				NodeID:      types.NodeID("node2"),
				Address:     "10.0.0.1:1234",
				State:       types.StateAlive,
				Incarnation: 5,
			},
		},
	}

	data, err := codec.Encode(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	decoded, err := codec.Decode(data)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	ack, ok := decoded.(*Ack)
	if !ok {
		t.Fatalf("decoded type = %T, want *Ping", decoded)
	}

	if ack.SeqNo != original.SeqNo {
		t.Errorf("SeqNo = %d, want %d", ack.SeqNo, original.SeqNo)
	}

	if ack.SourceID != original.SourceID {
		t.Errorf("SourceID = %s, want %s", ack.SourceID, original.SourceID)
	}

	if len(ack.Gossip) != len(original.Gossip) {
		t.Errorf("Gossip length = %d, want %d", len(ack.Gossip), len(original.Gossip))
	}
}

func TestCodec_RoundtripPingReq(t *testing.T) {
	// Create a PingRequest, encode it, decode it, compare
	codec := NewCodec()

	original := &PingReq{
		Ping: Ping{
			MessageHeader: MessageHeader{
				Version:  1,
				Type:     MessageTypePingReq,
				SeqNo:    42,
				SourceID: "node1",
			},
			Gossip: []GossipEntry{
				{
					NodeID:      types.NodeID("node2"),
					Address:     "10.0.0.1:1234",
					State:       types.StateAlive,
					Incarnation: 5,
				},
			},
		},
		TargetAdder: "node3",
	}

	data, err := codec.Encode(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	decoded, err := codec.Decode(data)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	ping, ok := decoded.(*PingReq)
	if !ok {
		t.Fatalf("decoded type = %T, want *PingReq", decoded)
	}

	if ping.SeqNo != original.SeqNo {
		t.Errorf("SeqNo = %d, want %d", ping.SeqNo, original.SeqNo)
	}

	if ping.SourceID != original.SourceID {
		t.Errorf("SourceID = %s, want %s", ping.SourceID, original.SourceID)
	}

	if len(ping.Gossip) != len(original.Gossip) {
		t.Errorf("Gossip length = %d, want %d", len(ping.Gossip), len(original.Gossip))
	}
}

func TestCodec_RoundtripWithGossipEntries(t *testing.T) {
	// Message with piggybacked gossip data
	codec := NewCodec()

	original := &Ping{
		MessageHeader: MessageHeader{
			Version:  1,
			Type:     MessageTypePing,
			SeqNo:    42,
			SourceID: "node1",
		},
		Gossip: []GossipEntry{
			{
				NodeID:      types.NodeID("node2"),
				Address:     "10.0.0.1:1234",
				State:       types.StateAlive,
				Incarnation: 5,
			},
			{
				NodeID:      types.NodeID("node3"),
				Address:     "10.0.0.2:1234",
				State:       types.StateAlive,
				Incarnation: 3,
			},
			{
				NodeID:      types.NodeID("node4"),
				Address:     "10.0.0.3:1234",
				State:       types.StateSuspect,
				Incarnation: 9,
			},
		},
	}

	data, err := codec.Encode(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	decoded, err := codec.Decode(data)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	ping, ok := decoded.(*Ping)
	if !ok {
		t.Fatalf("decoded type = %T, want *Ping", decoded)
	}

	if ping.SeqNo != original.SeqNo {
		t.Errorf("SeqNo = %d, want %d", ping.SeqNo, original.SeqNo)
	}

	if ping.SourceID != original.SourceID {
		t.Errorf("SourceID = %s, want %s", ping.SourceID, original.SourceID)
	}

	if len(ping.Gossip) != len(original.Gossip) {
		t.Errorf("Gossip length = %d, want %d", len(ping.Gossip), len(original.Gossip))
	}
}

func TestCodec_RoundtripUnknownType(t *testing.T) {
	// First invalid MessageType, should error
	codec := NewCodec()

	original := &envelope{}

	_, err := codec.Encode(original)
	if err == nil {
		t.Fatalf("Encode should have failed but didn't")
	}
}

func TestCodec_DecodeCorruptedData(t *testing.T) {
	// Garbage bytes, should error not panic
	codec := NewCodec()

	garbage := []byte{0x00, 0xFF, 0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE}

	_, err := codec.Decode(garbage)
	if err == nil {
		t.Fatal("Decode should have failed on corrupt data")
	}
}

func TestCodec_EmptyData(t *testing.T) {
	// Empty slice, should error
	codec := NewCodec()

	var err error
	_, err = codec.Decode([]byte{})
	if err == nil {
		t.Fatalf("Encode should have failed but didn't")
	}

	_, err = codec.Decode(nil)
	if err == nil {
		t.Fatalf("Encode should have failed but didn't")
	}
}
