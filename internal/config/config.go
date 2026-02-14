package config

import (
	"time"
)

// Config is the root configuration for the entire system
type Config struct {
	Node            NodeConfig
	Gossip          GossipConfig
	FailureDetector FailureDetectorConfig
}

// DefaultConfig returns defaults for all components
func DefaultConfig() Config {
	config := Config{
		Node:            NodeConfig{},
		Gossip:          GossipConfig{},
		FailureDetector: FailureDetectorConfig{},
	}
	config.Validate()
	return config
}

func (c *Config) LoadFromFile(filepath string) error {
	return nil
}

func (c *Config) LoadFromEnv() error {
	return nil
}

func (c *Config) Validate() error {
	var err error
	if err = c.Node.Validate(); err != nil {
		return err
	}
	if err = c.Gossip.Validate(); err != nil {
		return err
	}
	if err = c.FailureDetector.Validate(); err != nil {
		return err
	}
	return nil
}

// NodeConfig identifies this node
type NodeConfig struct {
	ID       string
	BindAddr string
	BindPort string
}

func (c *NodeConfig) Validate() error {
	return nil
}

func (c *Config) OverwriteStringProperty(key string, value any) {
}

// FailureDetectorConfig holds failure deterction tuning
type FailureDetectorConfig struct {
	ProbeInterval time.Duration
	ProbeTimeout  time.Duration
	IndirectNodes int
	SuspicionMult int
}

func (c *FailureDetectorConfig) Validate() error {
	return nil
}

func (c *FailureDetectorConfig) OverwriteStringProperty(key string, value any) {
}

// GossipConfig tunes the SWIM protocol
type GossipConfig struct {
	MaxBroadcast     int
	MaxGossipEntries int
	MaxPacketSize    int
}

func (c *GossipConfig) Validate() error {
	return nil
}

func (c *GossipConfig) OverwriteStringProperty(key string, value any) {
}
