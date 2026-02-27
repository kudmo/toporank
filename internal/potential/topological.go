// Package potential implements topological potential computation as defined
// in the PBScaler paper. The topological potential φ(i) for a node i is:
// φ(i) = a_i + Σ(a_j * e^(-(h_ji/σ)²)) for all upstream nodes j
package potential

import (
	"math"

	"github.com/kudmo/toporank/internal/graph"
	"github.com/kudmo/toporank/types"
)

// ComputeTopologicalPotential calculates the topological potential for every
// node in the graph. This becomes the personalization vector for PageRank.
//
// For each node, we start with its own anomaly score and add the decayed
// influence from all upstream nodes. The decay is exponential based on the
// distance (hops) and the impact factor σ.
func ComputeTopologicalPotential(
	g *types.CorrelationGraph,
	config types.TopoRankConfig,
) map[string]float64 {

	potential := make(map[string]float64, len(g.Nodes))

	// Calculate potential for each node
	for nodeID, node := range g.Nodes {
		// Start with own anomaly (Equation from Algorithm 2, line 3-4)
		phi := node.AnomalyDegree

		// Get all upstream nodes that can influence this node
		upstream := g.GetUpstreamNodes(nodeID)

		// Add influence from each upstream node (Algorithm 2, lines 5-9)
		for _, upstreamID := range upstream {
			upstreamNode, exists := g.Nodes[upstreamID]
			if !exists {
				continue
			}

			// Calculate minimum hops h_ji (Algorithm 2, line 6)
			hops := graph.ComputeMinHops(g, upstreamID, nodeID)
			if hops == math.MaxInt32 {
				continue // No path means no influence
			}

			// Exponential decay: e^(-(h/σ)²) (Algorithm 2, line 8)
			decay := math.Exp(-math.Pow(float64(hops)/config.ImpactFactor, 2))

			// Add influence: a_j * decay
			phi += upstreamNode.AnomalyDegree * decay
		}

		potential[nodeID] = phi
	}

	return potential
}

// NormalizePreferenceVector normalizes the topological potentials to sum to 1,
// creating a valid probability distribution for PageRank's personalization.
// This implements the normalization step before Algorithm 2, line 17.
func NormalizePreferenceVector(preferences map[string]float64) map[string]float64 {
	normalized := make(map[string]float64, len(preferences))

	// Calculate sum for normalization
	var sum float64
	for _, val := range preferences {
		sum += val
	}

	// Normalize or use uniform distribution if sum is zero
	if sum > 0 {
		for id, val := range preferences {
			normalized[id] = val / sum
		}
	} else {
		// If all potentials are zero, use uniform distribution
		uniform := 1.0 / float64(len(preferences))
		for id := range preferences {
			normalized[id] = uniform
		}
	}

	return normalized
}
