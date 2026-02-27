package types

import "errors"

// Common errors returned by the library
var (
	// ErrNodeNotFound indicates that a referenced node doesn't exist in the graph
	ErrNodeNotFound = errors.New("node not found in graph")

	// ErrInvalidWeight indicates that an edge weight is invalid (e.g., negative)
	ErrInvalidWeight = errors.New("edge weight must be between 0 and 1")

	// ErrInvalidAnomaly indicates that an anomaly score is out of valid range
	ErrInvalidAnomaly = errors.New("anomaly degree must be between 0 and 1")
)
