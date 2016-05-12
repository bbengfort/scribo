package scribo

import "fmt"

var nextNodeID uint64
var nextPingID uint64
var nodes Nodes
var pings Pings

func init() {
	RepoCreateNode(Node{Name: "alpha", Address: "127.0.0.11"})
	RepoCreateNode(Node{Name: "bravo", Address: "127.0.0.10"})
	RepoCreatePing(Ping{Source: 1, Target: 2, Latency: 11.142, Payload: 54})
	RepoCreatePing(Ping{Source: 1, Target: 2, Latency: 8.218, Payload: 54})
	RepoCreatePing(Ping{Source: 2, Target: 1, Timeout: true, Payload: 54})
}

// RepoFindNode searches for a node item by an id.
func RepoFindNode(id uint64) Node {
	for _, t := range nodes {
		if t.ID == id {
			return t
		}
	}

	// Return empty Node if not found
	return Node{}
}

// RepoFindPing searches for a node item by an id.
func RepoFindPing(id uint64) Ping {
	for _, t := range pings {
		if t.ID == id {
			return t
		}
	}

	// Return empty Node if not found
	return Ping{}
}

// RepoCreateNode inserts a node into the current nodes list
func RepoCreateNode(n Node) Node {
	nextNodeID++
	n.ID = nextNodeID
	nodes = append(nodes, n)
	return n
}

// RepoCreatePing inserts a ping into the current pings list
func RepoCreatePing(p Ping) Ping {
	nextPingID++
	p.ID = nextPingID
	pings = append(pings, p)
	return p
}

// RepoDestroyNode deletes the ping with the given ID.
func RepoDestroyNode(id uint64) error {
	for i, t := range nodes {
		if t.ID == id {
			nodes = append(nodes[:i], nodes[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("Could not find Node with id of %d to delete", id)
}

// RepoDestroyPing deletes the ping with the given ID.
func RepoDestroyPing(id uint64) error {
	for i, t := range pings {
		if t.ID == id {
			nodes = append(nodes[:i], nodes[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("Could not find Ping with id of %d to delete", id)
}
