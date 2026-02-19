package types

// Package types contains simple data structures used across the TopoRank
// implementation. These types are intentionally lightweight and intended for
// examples and small graphs.

// Node represents a vertex in the graph. Fields are exported for simplicity
// in examples and tests.
type Node struct {
	ID        string
	Score     float64
	Neighbors []string
}

// Graph is a minimal adjacency-style graph container mapping node IDs to
// Node objects.
type Graph struct {
	Nodes map[string]*Node
}

// NewGraph creates and returns an empty Graph instance.
func NewGraph() *Graph {
	return &Graph{Nodes: make(map[string]*Node)}
}

// AddNode ensures a node with the provided id exists in the graph.
// If the node is already present this is a no-op.
func (g *Graph) AddNode(id string) {
	if _, exists := g.Nodes[id]; !exists {
		g.Nodes[id] = &Node{ID: id, Neighbors: []string{}}
	}
}

// AddEdge appends a directed edge from `from` to `to`. Both nodes are
// assumed to already exist; if `from` is missing the call is a no-op.
func (g *Graph) AddEdge(from, to string) {
	if node, exists := g.Nodes[from]; exists {
		node.Neighbors = append(node.Neighbors, to)
	}
}
