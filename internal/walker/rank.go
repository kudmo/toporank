package walker

import (
	"sort"

	"github.com/kudmo/toporank/types"
)

// RankNodes returns a slice of node pointers sorted by descending score. The
// original graph is not modified; nodes are collected and ordered for
// presentation or downstream processing.
func RankNodes(g *types.Graph) []*types.Node {
	nodes := make([]*types.Node, 0, len(g.Nodes))
	for _, n := range g.Nodes {
		nodes = append(nodes, n)
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Score > nodes[j].Score
	})
	return nodes
}
