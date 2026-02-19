package api

import (
	"math"
	"sort"

	"github.com/kudmo/toporank/internal/potentials"
	"github.com/kudmo/toporank/types"
)

// RunTopoRank executes the TopoRank algorithm on graph g using the supplied
// configuration and optional anomaly potentials. The function updates node
// scores in-place and returns a slice of nodes ordered by descending score.
//
// Parameters:
//   - g: the graph to rank (nodes should be present in g.Nodes)
//   - config: walker configuration controlling iterations, retention and sigma
//   - anomalyPotentials: optional initial score overrides for specific nodes
func RunTopoRank(g *types.Graph, config types.RandomWalkConfig, anomalyPotentials map[string]float64) []*types.Node {
	// 1. Initialize node scores from anomaly potentials when provided.
	for id, val := range anomalyPotentials {
		if node, ok := g.Nodes[id]; ok {
			node.Score = val
		}
	}

	// 2. Compute topological potential used as neighbor weighting.
	tPot := potentials.ComputeTopologicalPotential(g, config.Sigma)

	// 3. Random-walk style score propagation.
	for iter := 0; iter < config.MaxIter; iter++ {
		diff := 0.0
		newScores := make(map[string]float64)

		for id, node := range g.Nodes {
			share := node.Score * (1 - config.SelfRetention)
			neighbors := node.Neighbors
			weightSum := 0.0
			weights := make([]float64, len(neighbors))

			// Distribute `share` among neighbors proportionally to their
			// topological potential.
			for i, nbrID := range neighbors {
				w := tPot[nbrID]
				weights[i] = w
				weightSum += w
			}

			for i, nbrID := range neighbors {
				if weightSum > 0 {
					newScores[nbrID] += share * (weights[i] / weightSum)
				} else {
					// If no neighbors or zero total weight, keep share on node.
					newScores[id] += share
				}
			}

			// Add retained self-score.
			newScores[id] += node.Score * config.SelfRetention
		}

		// 4. Check convergence by total L1 change and commit new scores.
		for id, node := range g.Nodes {
			d := math.Abs(newScores[id] - node.Score)
			diff += d
			node.Score = newScores[id]
		}

		if diff < config.ConvergenceTol {
			break
		}
	}

	// 5. Normalize scores to sum to 1.
	total := 0.0
	for _, node := range g.Nodes {
		total += node.Score
	}
	if total > 0 {
		for _, node := range g.Nodes {
			node.Score /= total
		}
	}

	// 6. Collect and sort nodes by score descending.
	nodes := make([]*types.Node, 0, len(g.Nodes))
	for _, node := range g.Nodes {
		nodes = append(nodes, node)
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Score > nodes[j].Score
	})

	return nodes
}
