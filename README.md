# TopoRank – Topological Node Ranking for Microservice Anomaly Detection

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Lightweight Go 1.21 library for **ranking nodes in directed weighted graphs** using
**topological potentials** and **personalized PageRank**. This library implements the 
**exact TopoRank algorithm** from the [PBScaler: Prioritized Scaling of Microservices via Topological Ranking](https://arxiv.org/html/2303.14620v3) 
paper, designed for detecting anomalous or influential nodes in microservice architectures
and general directed graphs (up to ~100 nodes).

The algorithm combines two key concepts from the paper:
- **Topological potential** computation with exponential distance decay from upstream nodes
- **Personalized PageRank** with correlation-based transition probabilities

### Install

```bash
go get github.com/kudmo/toporank
```

### Core concepts

**CorrelationGraph (types.CorrelationGraph)**: Main structure representing nodes with anomaly scores and directed weighted edges. Edge weights represent correlation strengths from metrics monitoring.

**Node (types.Node)**: Holds ID, anomaly degree, computed preference (topological potential), and final rank score.

**TopoRankConfig (types.TopoRankConfig)**: Configuration parameters controlling impact factor (σ), damping factor, convergence tolerance, and maximum iterations.

**Functional API (api.RunTopoRank)**: Runs the exact PBScaler TopoRank algorithm and returns nodes sorted by final rank.

**Internal modules (internal/potential, internal/pagerank, internal/graph)**: Handle topological potential computation, personalized PageRank, and graph utilities – completely hidden from the user.

### Quick start

```go
package main

import (
    "fmt"
    "log"

    "github.com/kudmo/toporank/api"
    "github.com/kudmo/toporank/types"
)

func main() {
	g := types.NewCorrelationGraph()

	// Nodes (anomaly degree aligned with example table)
	g.AddNode("A", 12)
	g.AddNode("B", 4)
	g.AddNode("C", 8)
	g.AddNode("D", 4)
	g.AddNode("E", 4) // Bottleneck microservice
	// Nodes F,G,H aren't in abnormal subgraph

	// Edges (topology strictly follows the diagram)

	// A dependencies
	g.AddEdge("A", "B", 0.31)
	g.AddEdge("A", "C", 0.42)

	// B dependencies
	g.AddEdge("B", "D", 0.37)

	// C dependencies
	g.AddEdge("C", "D", 0.29)
	g.AddEdge("C", "E", 0.51)

	// D dependencies
	g.AddEdge("D", "E", 0.62)

	// TopoRank configuration
	config := types.TopoRankConfig{
		ImpactFactor:  1.0,  // σ = 1
		DampingFactor: 0.85, // α = 0.85
		MaxIterations: 1000,
		Tolerance:     1e-8,
	}

	ranked := api.RunTopoRank(g, config)

	fmt.Println("Graph Example Results:")
	fmt.Println("Rank | Service | Anomaly | Potential | Final Rank")
	fmt.Println("-----|---------|---------|-----------|-----------")

	for i, node := range ranked {
		fmt.Printf("%-4d | %-7s | %.2f    | %.2f      |%.6f\n",
			i+1, node.ID, node.AnomalyDegree, node.Preference, node.Rank)
	}
}
```

### Configuration

The main configuration struct is `types.TopoRankConfig`:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `ImpactFactor` | σ (sigma) – Controls distance decay in topological potential. Larger values allow influence to propagate further through the graph. | 1.0 |
| `DampingFactor` | Standard PageRank damping factor. Probability of following edges vs. teleporting according to preference vector. | 0.85 |
| `MaxIterations` | Maximum number of PageRank power iterations. | 100 |
| `Tolerance` | Convergence threshold – iterations stop when total L1 change in ranks is below this value. | 1e-6 |

### Algorithm Details

TopoRank implements Algorithm 2 from the PBScaler paper:

1. **Topological Potential Computation** (Preference Vector 𝒖)
   - For each node, start with its own anomaly degree
   - Add decayed influence from all upstream nodes: `a_j * e^(-(h_ji/σ)²)`
   - Where `h_ji` is the minimum hops from node j to node i

2. **Transition Matrix Construction** (Matrix P)
   - Build from correlation weights in the graph
   - Each row normalized to sum to 1 (stochastic matrix)

3. **Personalized PageRank**
   - `rank = (1-d) * 𝒖 + d * P^T * rank`
   - Power iteration until convergence

4. **Ranked Output**
   - Nodes sorted by descending final rank score

### Features

- **Exact PBScaler implementation** – Faithful to Algorithm 2 from the research paper
- **Correlation graph input** – Works with pre-computed weighted graphs from metrics
- **Pure Go** – No external dependencies, uses only standard library
- **Clean architecture** – Separated internal packages for maintainability
- **Comprehensive documentation** – Every step documented with paper references
- **Error handling** – Input validation with meaningful errors
- **Examples included** – Multiple usage scenarios in `/examples`

### Project Structure

```
toporank/
├── api/            # Public API (RunTopoRank)
├── types/          # Core data structures (Node, Graph, Config)
├── internal/       
│   ├── graph/      # Graph utilities (path finding, matrix building)
│   ├── potential/  # Topological potential computation
│   └── pagerank/   # Personalized PageRank implementation
└── examples/       # Usage examples
```

### Requirements

- Go 1.21 or higher
- Input graph must be directed with weighted edges (0-1)

### License

MIT License – See [LICENSE](LICENSE) file for details.

### References

- [PBScaler: Prioritized Scaling of Microservices via Topological Ranking](https://arxiv.org/html/2303.14620v3) (arXiv:2303.14620)
- Original paper presented at IEEE ICWS 2023
```
