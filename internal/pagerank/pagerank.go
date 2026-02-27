// Package pagerank implements personalized PageRank for the TopoRank algorithm.
// The implementation follows the standard power iteration method with
// personalization vector u and transition matrix P from the correlation graph.
package pagerank

import (
	"math"

	"github.com/kudmo/toporank/types"
)

// PageRank computes personalized PageRank scores using the power iteration method.
// The formula is: r = (1-d) * u + d * P^T * r
// where:
//   - r is the rank vector
//   - d is the damping factor
//   - u is the personalization vector (topological potential)
//   - P is the transition matrix from the correlation graph
//
// This implements Algorithm 2, line 17 from the PBScaler paper.
func PageRank(
	g *types.CorrelationGraph,
	transitionMatrix map[string]map[string]float64,
	preferenceVector map[string]float64,
	config types.TopoRankConfig,
) map[string]float64 {

	// Initialize rank vectors
	ranks := make(map[string]float64, len(g.Nodes))
	newRanks := make(map[string]float64, len(g.Nodes))

	// Start with the preference vector as initial ranks (Algorithm 2, line 1)
	for id, pref := range preferenceVector {
		ranks[id] = pref
	}

	n := float64(len(g.Nodes))
	teleportBase := (1 - config.DampingFactor)

	// Power iteration (Algorithm 2, line 17)
	for iter := 0; iter < config.MaxIterations; iter++ {
		// Initialize with teleportation component: (1-d) * u_i
		for id := range g.Nodes {
			newRanks[id] = teleportBase * preferenceVector[id]
		}

		// Add contribution from following edges: d * Σ(rank[j] * P[j][i])
		for fromID := range g.Nodes {
			if ranks[fromID] == 0 {
				continue
			}

			transitions := transitionMatrix[fromID]
			if len(transitions) == 0 {
				// No outgoing edges - distribute evenly to all nodes
				// (standard PageRank sink handling)
				contribution := config.DampingFactor * ranks[fromID] / n
				for toID := range g.Nodes {
					newRanks[toID] += contribution
				}
			} else {
				// Normal case - follow weighted edges
				for toID, prob := range transitions {
					newRanks[toID] += config.DampingFactor * ranks[fromID] * prob
				}
			}
		}

		// Check for convergence (Algorithm 2 convergence check)
		diff := 0.0
		for id := range g.Nodes {
			diff += math.Abs(newRanks[id] - ranks[id])
			ranks[id] = newRanks[id]
		}

		if diff < config.Tolerance {
			break
		}
	}

	return ranks
}
