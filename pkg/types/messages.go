package types

import "time"

// NetworkMessage represents a raw message received from the metwork
// This is the low-level transport representation before protocol parsing
type NetworkMessage struct {
	Payload []byte
	From    string
	Time    time.Time
}
