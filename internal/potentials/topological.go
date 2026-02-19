package potentials

import (
	"math"

	graphutils "github.com/kudmo/toporank/internal/graph"
	"github.com/kudmo/toporank/types"
)

// ComputeTopologicalPotential computes a topological potential for each node
// in the graph. The potential of a node i is the sum over all nodes j of
// j.Score * exp(-d(i,j)^2 / (2*sigma^2)), where d(i,j) is the shortest-path
// distance (in edges) from i to j. The sigma parameter controls distance decay.
//
// The function returns a map from node ID to its computed potential value.
func ComputeTopologicalPotential(graph *types.Graph, sigma float64) map[string]float64 {
	potentials := make(map[string]float64)

	for idI := range graph.Nodes {
		sum := 0.0
		for idJ, nodeJ := range graph.Nodes {
			// Compute shortest-path distance from i to j.
			d := graphutils.BFS(graph, idI, idJ)
			// Add contribution of j with Gaussian decay by distance.
			sum += nodeJ.Score * math.Exp(-float64(d*d)/(2*sigma*sigma))
		}
		potentials[idI] = sum
	}

	return potentials
}
