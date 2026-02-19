package walker

import (
	"math"

	"github.com/kudmo/toporank/internal/potentials"
	"github.com/kudmo/toporank/types"
)

// RunRandomWalk performs a random-walk style score propagation on graph g
// according to the supplied configuration. The function updates node
// scores in-place on the provided graph.
//
// The algorithm uses a topological potential (computed by the potentials
// package) to weight how score is distributed among neighbors. `SelfRetention`
// controls how much of the node score remains with the node between
// iterations. The routine stops early if the total L1 change across all
// node scores falls below `ConvergenceTol`.
func RunRandomWalk(g *types.Graph, config types.RandomWalkConfig) {
	maxIter := config.MaxIter
	selfRetention := config.SelfRetention
	tol := config.ConvergenceTol
	sigma := config.Sigma

	tPot := potentials.ComputeTopologicalPotential(g, sigma)

	for iter := 0; iter < maxIter; iter++ {
		diff := 0.0
		newScores := make(map[string]float64)

		for id, node := range g.Nodes {
			share := node.Score * (1 - selfRetention)
			neighbors := node.Neighbors
			weightSum := 0.0
			weights := make([]float64, len(neighbors))
			for i, nbrID := range neighbors {
				w := tPot[nbrID]
				weights[i] = w
				weightSum += w
			}
			for i, nbrID := range neighbors {
				if weightSum > 0 {
					newScores[nbrID] += share * (weights[i] / weightSum)
				} else {
					newScores[id] += share
				}
			}
			newScores[id] += node.Score * selfRetention
		}

		for id, node := range g.Nodes {
			d := math.Abs(newScores[id] - node.Score)
			diff += d
			node.Score = newScores[id]
		}

		if diff < tol {
			break
		}
	}
}
