package main

import (
	"fmt"

	"github.com/kudmo/toporank/api"
	"github.com/kudmo/toporank/types"
)

// A tiny example demonstrating how to build a graph and run TopoRank.
func main() {
	g := types.NewGraph()
	g.AddNode("A")
	g.AddNode("B")
	g.AddNode("C")

	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	g.AddEdge("A", "C")

	// Optional initial anomaly potentials (seed scores).
	anomalyPotentials := map[string]float64{
		"A": 0.5,
		"B": 0.3,
		"C": 0.2,
	}

	// Random walk configuration used by the algorithm.
	config := types.RandomWalkConfig{
		MaxIter:        20,
		SelfRetention:  0.15,
		ConvergenceTol: 1e-6,
		Sigma:          1.0,
	}

	// Run the ranking algorithm and print results.
	res := api.RunTopoRank(g, config, anomalyPotentials)

	fmt.Println("TopoRank results:")
	for i, p := range res {
		fmt.Printf("%d. %s = %.6f\n", i+1, p.ID, p.Score)
	}
}
