package scribo

import "time"

// Node is a model that represents a participant in the network
type Node struct {
	ID      uint64    `json:"id"`      // Unique ID of the node
	Name    string    `json:"name"`    // Name/DNS of the node
	Address string    `json:"address"` // IP Address of the node
	DNS     string    `json:"dns"`     // DNS Lookup for the node
	Key     string    `json:"-"`       // Authentication key of the node
	Created time.Time `json:"created"` // Datetime the node was created
	Update  time.Time `json:"updated"` // Datetime the node was updated
}

// Ping is a model that represents a latency report.
type Ping struct {
	ID      uint64    `json:"id"`      // Unique ID of the ping
	Source  uint64    `json:"source"`  // The ID of the source node
	Target  uint64    `json:"target"`  // The ID of the target node
	Payload uint      `json:"payload"` // The size in bytes of the payload
	Latency float64   `json:"latency"` // The time in ms of the round trip
	Timeout bool      `json:"timeout"` // Whether or not the request timed out
	Created time.Time `json:"created"` // Datetime the ping was created
	Update  time.Time `json:"updated"` // Datetime the ping was updated
}

// Nodes is a collection of node items for use elsewhere.
type Nodes []Node

// Pings is a collection of latency reports for use elsewhere.
type Pings []Ping
