package shortestpath

import (
	"math"

	"github.com/kudmo/toporank/types"
)

// BFS returns the length (number of edges) of the shortest path from
// fromID to toID in an unweighted graph using breadth-first search. If
// fromID equals toID, zero is returned. If no path exists, a large
// sentinel value (1000) is returned.
func BFS(graph *types.Graph, fromID, toID string) int {
	if fromID == toID {
		return 0
	}

	visited := map[string]bool{fromID: true}
	queue := []string{fromID}
	depth := map[string]int{fromID: 0}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		currDepth := depth[curr]

		for _, nbr := range graph.Nodes[curr].Neighbors {
			if nbr == toID {
				return currDepth + 1
			}
			if !visited[nbr] {
				visited[nbr] = true
				queue = append(queue, nbr)
				depth[nbr] = currDepth + 1
			}
		}
	}

	return math.MaxInt32
}
