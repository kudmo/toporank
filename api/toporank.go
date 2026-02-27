// Package api provides the public interface for the TopoRank algorithm.
// It exposes a simple, functional API that takes a correlation graph and
// configuration, then returns ranked nodes.
package api

import (
	"sort"

	"github.com/kudmo/toporank/internal/graph"
	"github.com/kudmo/toporank/internal/pagerank"
	"github.com/kudmo/toporank/internal/potential"
	"github.com/kudmo/toporank/types"
)

// RunTopoRank executes the TopoRank algorithm on a pre-computed correlation graph.
// The algorithm follows these steps:
//  1. Compute topological potential for each node (preference vector u)
//  2. Normalize the preference vector
//  3. Build transition matrix P from graph edge weights
//  4. Run personalized PageRank with preference vector u and matrix P
//  5. Return nodes sorted by descending rank
//
// Parameters:
//   - g: A weighted directed graph with node anomaly scores and edge correlations
//   - config: Algorithm parameters (impact factor, damping factor, etc.)
//
// Returns:
//
//	A slice of nodes sorted by descending TopoRank score.
//	The nodes themselves are modified to store their computed Preference and Rank.
func RunTopoRank(g *types.CorrelationGraph, config types.TopoRankConfig) []*types.Node {
	// Use default config if none provided
	if config == (types.TopoRankConfig{}) {
		config = types.DefaultConfig()
	}

	// Step 1: Compute topological potential (preference vector u)
	// This implements Algorithm 2, lines 2-10
	preferences := potential.ComputeTopologicalPotential(g, config)

	// Step 2: Normalize preference vector for PageRank
	normalizedPrefs := potential.NormalizePreferenceVector(preferences)

	// Store preference scores in nodes for debugging/analysis
	for id, pref := range preferences {
		if node, ok := g.Nodes[id]; ok {
			node.Preference = pref
		}
	}

	// Step 3: Build transition matrix P from weighted edges
	// This corresponds to the correlation-based transition probabilities
	transitionMatrix := graph.BuildTransitionMatrix(g)

	// Step 4: Run personalized PageRank (Algorithm 2, line 17)
	ranks := pagerank.PageRank(g, transitionMatrix, normalizedPrefs, config)

	// Store final ranks in nodes
	for id, rank := range ranks {
		if node, ok := g.Nodes[id]; ok {
			node.Rank = rank
		}
	}

	// Step 5: Sort nodes by descending rank (Algorithm 2, line 18)
	result := make([]*types.Node, 0, len(g.Nodes))
	for _, node := range g.Nodes {
		result = append(result, node)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Rank > result[j].Rank
	})

	return result
}
