// Package graph provides internal graph utilities for the TopoRank algorithm.
package graph

import (
	"math"

	"github.com/kudmo/toporank/types"
)

// ComputeMinHops finds the shortest path distance between two nodes in a
// directed graph using BFS. Returns MaxInt32 if no path exists.
// This is used to calculate the exponential decay factor e^(-(h/σ)²)
// in topological potential computation.
func ComputeMinHops(g *types.CorrelationGraph, from, to string) int {
	if from == to {
		return 0
	}

	// BFS initialization
	visited := make(map[string]bool)
	queue := []struct {
		node string
		dist int
	}{{from, 0}}

	visited[from] = true

	// BFS loop
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		edges := g.GetOutEdges(current.node)
		for _, edge := range edges {
			if edge.To == to {
				return current.dist + 1
			}
			if !visited[edge.To] {
				visited[edge.To] = true
				queue = append(queue, struct {
					node string
					dist int
				}{edge.To, current.dist + 1})
			}
		}
	}

	// No path found
	return math.MaxInt32
}

// BuildTransitionMatrix converts the weighted correlation graph into a
// stochastic transition matrix for PageRank. Each row is normalized so that
// the sum of outgoing probabilities equals 1.
//
// For nodes with no outgoing edges, an empty map is returned - PageRank will
// handle these as teleportation sinks.
func BuildTransitionMatrix(g *types.CorrelationGraph) map[string]map[string]float64 {
	transitions := make(map[string]map[string]float64)

	for fromID := range g.Nodes {
		edges := g.GetOutEdges(fromID)
		transitions[fromID] = make(map[string]float64)

		if len(edges) == 0 {
			// No outgoing edges - PageRank will handle via teleportation
			continue
		}

		// Calculate total weight for normalization
		var totalWeight float64
		for _, edge := range edges {
			totalWeight += edge.Weight
		}

		// Normalize weights to probabilities
		for _, edge := range edges {
			if totalWeight > 0 {
				transitions[fromID][edge.To] = edge.Weight / totalWeight
			} else {
				// If all weights are zero, distribute uniformly
				transitions[fromID][edge.To] = 1.0 / float64(len(edges))
			}
		}
	}

	return transitions
}
