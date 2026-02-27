package types

// TopoRankConfig contains configurable parameters for the TopoRank algorithm.
// These parameters control both the topological potential calculation and
// the personalized PageRank execution.
type TopoRankConfig struct {
	// ImpactFactor (σ) controls the distance decay in topological potential.
	// Larger values mean influence propagates further through the graph.
	ImpactFactor float64

	// DampingFactor is the standard PageRank damping parameter.
	// Typically set to 0.85, it controls the probability of following edges
	// vs. teleporting according to the preference vector.
	DampingFactor float64

	// MaxIterations limits the number of PageRank iterations.
	MaxIterations int

	// Tolerance determines convergence - iterations stop when the total
	// L1 change in ranks is below this threshold.
	Tolerance float64
}

// DefaultConfig returns a TopoRankConfig with recommended default values.
// These values work well for most microservice anomaly detection scenarios.
func DefaultConfig() TopoRankConfig {
	return TopoRankConfig{
		ImpactFactor:  1.0,
		DampingFactor: 0.85,
		MaxIterations: 100,
		Tolerance:     1e-6,
	}
}
