package types

// RandomWalkConfig contains tunable parameters for the random-walk based
// ranking algorithm used by TopoRank.
type RandomWalkConfig struct {
	// MaxIter is the maximum number of iterations to run the walker.
	MaxIter int

	// SelfRetention controls how much score a node retains between iterations
	// (value in [0,1]). A typical value is 0.15.
	SelfRetention float64

	// ConvergenceTol is the threshold for the total L1 change across all
	// node scores below which the algorithm will stop early.
	ConvergenceTol float64

	// Sigma is a distance decay scale used when computing topological
	// potentials.
	Sigma float64
}
