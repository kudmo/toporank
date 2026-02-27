// Package types provides the core data structures for the TopoRank algorithm.
package types

// Node represents a graph vertex (microservice in the original paper).
// Each node has an ID, anomaly score, and stores computed values.
type Node struct {
	// ID is the unique identifier for the node
	ID string

	// AnomalyDegree is the initial anomaly score (ai in the paper)
	// Typically ranges from 0 (normal) to 1 (highly anomalous)
	AnomalyDegree float64

	// Rank is the final TopoRank score after PageRank computation
	Rank float64

	// Preference is the topological potential value (ui in the paper)
	// This is the personalization vector for PageRank
	Preference float64
}

// Edge represents a weighted directed edge in the correlation graph.
// The weight typically represents correlation strength from metrics.
type Edge struct {
	// From is the source node ID
	From string

	// To is the target node ID
	To string

	// Weight represents the correlation strength (typically 0-1)
	// Higher values indicate stronger influence
	Weight float64
}

// CorrelationGraph is a weighted directed graph representing service dependencies
// and their correlation strengths. This is the main input structure for TopoRank.
type CorrelationGraph struct {
	// Nodes maps node IDs to Node objects
	Nodes map[string]*Node

	// Edges maps source node IDs to their outgoing edges
	Edges map[string][]*Edge
}

// NewCorrelationGraph creates a new empty correlation graph.
func NewCorrelationGraph() *CorrelationGraph {
	return &CorrelationGraph{
		Nodes: make(map[string]*Node),
		Edges: make(map[string][]*Edge),
	}
}

// AddNode adds a node to the graph with the given anomaly degree.
// Returns an error if the anomaly degree is invalid (<0).
func (g *CorrelationGraph) AddNode(id string, anomalyDegree float64) error {
	if anomalyDegree < 0 {
		return ErrInvalidAnomaly
	}

	if _, exists := g.Nodes[id]; !exists {
		g.Nodes[id] = &Node{
			ID:            id,
			AnomalyDegree: anomalyDegree,
		}
	}
	return nil
}

// AddEdge adds a weighted directed edge to the graph.
// Returns an error if nodes don't exist or weight is invalid.
func (g *CorrelationGraph) AddEdge(from, to string, weight float64) error {
	// Validate nodes exist
	if _, exists := g.Nodes[from]; !exists {
		return ErrNodeNotFound
	}
	if _, exists := g.Nodes[to]; !exists {
		return ErrNodeNotFound
	}

	// Validate weight range (typically correlations are 0-1)
	if weight < 0 || weight > 1 {
		return ErrInvalidWeight
	}

	edge := &Edge{
		From:   from,
		To:     to,
		Weight: weight,
	}

	g.Edges[from] = append(g.Edges[from], edge)
	return nil
}

// GetOutEdges returns all outgoing edges from a node.
// Returns nil if the node has no outgoing edges.
func (g *CorrelationGraph) GetOutEdges(nodeID string) []*Edge {
	return g.Edges[nodeID]
}

// GetUpstreamNodes returns all nodes that can reach the target node through
// directed paths. This is used to compute topological potential from all
// upstream influences, as defined in the PBScaler paper.
func (g *CorrelationGraph) GetUpstreamNodes(targetID string) []string {
	visited := make(map[string]bool)
	var upstream []string

	// DFS to find all paths to target
	var dfs func(current string)
	dfs = func(current string) {
		for _, edge := range g.Edges[current] {
			if edge.To == targetID && !visited[edge.From] {
				visited[edge.From] = true
				upstream = append(upstream, edge.From)
				dfs(edge.From)
			}
		}
	}

	// Start DFS from all nodes
	for nodeID := range g.Nodes {
		if nodeID != targetID && !visited[nodeID] {
			dfs(nodeID)
		}
	}

	return upstream
}
