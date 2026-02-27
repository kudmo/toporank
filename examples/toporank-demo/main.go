// Package main provides a simple example of using the TopoRank library.
package main

import (
	"fmt"
	"log"

	"github.com/kudmo/toporank/api"
	"github.com/kudmo/toporank/types"
)

func main() {
	// Example 1: Simple graph with three services
	fmt.Println("=== Example 1: Simple Microservice Graph ===")
	simpleExample()

	fmt.Println("\n=== Example 2: Multi-tier Architecture ===")
	multiTierExample()

	fmt.Println("\n=== Example 3: Handling Disconnected Components ===")
	disconnectedExample()

	fmt.Println("\n=== Example 4: Article Graph Example ===")
	topoRankGraphExample()
}

// simpleExample demonstrates a basic 3-node graph with linear dependencies.
func simpleExample() {
	// Create a new correlation graph
	g := types.NewCorrelationGraph()

	// Add nodes with their anomaly scores (from monitoring)
	// Higher scores indicate more anomalous behavior
	if err := g.AddNode("payment", 0.8); err != nil {
		log.Fatal(err)
	}
	if err := g.AddNode("inventory", 0.3); err != nil {
		log.Fatal(err)
	}
	if err := g.AddNode("shipping", 0.2); err != nil {
		log.Fatal(err)
	}

	// Add weighted edges (correlations computed from metrics)
	// Edge weights represent influence strength (0-1)
	if err := g.AddEdge("payment", "inventory", 0.9); err != nil {
		log.Fatal(err)
	}
	if err := g.AddEdge("inventory", "shipping", 0.7); err != nil {
		log.Fatal(err)
	}
	if err := g.AddEdge("payment", "shipping", 0.4); err != nil {
		log.Fatal(err)
	}

	// Run TopoRank with default configuration
	config := types.DefaultConfig()
	// Adjust impact factor to control distance decay
	config.ImpactFactor = 1.5

	ranked := api.RunTopoRank(g, config)

	// Print results
	fmt.Println("Ranked services (most anomalous first):")
	fmt.Println("Rank | Service   | Anomaly | Potential | Final Rank")
	fmt.Println("-----|-----------|---------|-----------|-----------")
	for i, node := range ranked {
		fmt.Printf("%-4d | %-9s | %.2f    | %.4f    | %.6f\n",
			i+1, node.ID, node.AnomalyDegree, node.Preference, node.Rank)
	}
}

// multiTierExample demonstrates a more complex microservice architecture.
func multiTierExample() {
	g := types.NewCorrelationGraph()

	// Frontend services
	g.AddNode("web-ui", 0.2)
	g.AddNode("mobile-api", 0.3)

	// Business logic services
	g.AddNode("order-service", 0.6)
	g.AddNode("user-service", 0.1)
	g.AddNode("product-service", 0.4)

	// Data services
	g.AddNode("pay-processor", 0.9)
	g.AddNode("inventory-db", 0.2)
	g.AddNode("user-db", 0.1)

	// Add dependencies with correlation strengths
	// Frontend to business logic
	g.AddEdge("web-ui", "order-service", 0.8)
	g.AddEdge("web-ui", "product-service", 0.7)
	g.AddEdge("mobile-api", "order-service", 0.9)
	g.AddEdge("mobile-api", "user-service", 0.6)

	// Business logic to data services
	g.AddEdge("order-service", "pay-processor", 0.95)
	g.AddEdge("order-service", "inventory-db", 0.7)
	g.AddEdge("user-service", "user-db", 0.8)
	g.AddEdge("product-service", "inventory-db", 0.6)

	// Cross-service dependencies
	g.AddEdge("order-service", "user-service", 0.4)
	g.AddEdge("pay-processor", "inventory-db", 0.3)

	// Custom configuration
	config := types.TopoRankConfig{
		ImpactFactor:  2.0, // Allow influence to propagate further
		DampingFactor: 0.85,
		MaxIterations: 100,
		Tolerance:     1e-8,
	}

	ranked := api.RunTopoRank(g, config)

	fmt.Println("Multi-tier Architecture Results:")
	fmt.Println("Rank | Service          | Anomaly | Final Rank")
	fmt.Println("-----|------------------|---------|-----------")
	for i, node := range ranked[:5] { // Show top 5
		fmt.Printf("%-4d | %-16s | %.2f    | %.6f\n",
			i+1, node.ID, node.AnomalyDegree, node.Rank)
	}
}

// disconnectedExample shows how the algorithm handles disconnected components.
func disconnectedExample() {
	g := types.NewCorrelationGraph()

	// Component A (high anomaly, interconnected)
	g.AddNode("service-a1", 0.8)
	g.AddNode("service-a2", 0.7)
	g.AddNode("service-a3", 0.6)
	g.AddEdge("service-a1", "service-a2", 0.9)
	g.AddEdge("service-a2", "service-a3", 0.8)
	g.AddEdge("service-a1", "service-a3", 0.7)

	// Component B (low anomaly, isolated)
	g.AddNode("service-b1", 0.1)
	g.AddNode("service-b2", 0.1)
	g.AddEdge("service-b1", "service-b2", 0.5)

	// Component C (medium anomaly, but no outgoing edges)
	g.AddNode("service-c1", 0.5)

	config := types.DefaultConfig()
	ranked := api.RunTopoRank(g, config)

	fmt.Println("Disconnected Components Results:")
	fmt.Println("Rank | Service    | Anomaly | Final Rank | Component")
	fmt.Println("-----|------------|---------|------------|----------")
	for i, node := range ranked {
		// Determine component (simplified)
		component := "A"
		if node.ID == "service-b1" || node.ID == "service-b2" {
			component = "B"
		} else if node.ID == "service-c1" {
			component = "C"
		}
		fmt.Printf("%-4d | %-10s | %.2f    |  %.6f  | %s\n",
			i+1, node.ID, node.AnomalyDegree, node.Rank, component)
	}
}

// topoRankGraphExample reproduces the graph structure from the provided diagram.
func topoRankGraphExample() {
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
