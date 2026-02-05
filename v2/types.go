package edlib

import (
	"time"

	"github.com/hbollon/go-edlib/v2/algorithms"
)

// ComparisonResult contains the result of a string comparison operation
type ComparisonResult struct {
	// Distance is the edit distance between the strings (for distance-based algorithms)
	Distance int

	// Similarity is the similarity score between the strings [0.0 to 1.0]
	Similarity float64

	// Algorithm used for the comparison
	Algorithm algorithms.Algorithm

	// Metadata contains additional information about the comparison
	Metadata ComparisonMetadata

	// Metrics contains execution metrics if metrics were enabled
	Metrics *ExecutionMetrics

	// Error contains any error that occurred during the comparison
	Error error
}

// ComparisonMetadata contains metadata about a comparison operation
type ComparisonMetadata struct {
	// CaseSensitive indicates whether the comparison was case-sensitive
	CaseSensitive bool

	// NormalizedStrings indicates whether the strings were normalized
	NormalizedStrings bool

	// MaxDistance is the maximum possible distance for the comparison
	MaxDistance int

	// UsedOptimizations lists the optimizations that were applied
	UsedOptimizations []string
}

// ExecutionMetrics contains performance metrics for an operation
type ExecutionMetrics struct {
	// ExecutionTime is the total time taken for the operation
	ExecutionTime time.Duration

	// MemoryAllocated is the amount of memory allocated in bytes
	MemoryAllocated int64

	// MemoryFreed is the amount of memory freed in bytes
	MemoryFreed int64

	// CacheHits is the number of cache hits
	CacheHits int

	// CacheMisses is the number of cache misses
	CacheMisses int

	// PoolHits is the number of memory pool hits
	PoolHits int

	// PoolMisses is the number of memory pool misses
	PoolMisses int

	// OptimizationsUsed lists the optimizations that were applied
	OptimizationsUsed []string

	// WorkersUsed is the number of parallel workers used
	WorkersUsed int
}

// StringPair represents a pair of strings to compare
type StringPair struct {
	// First string in the pair
	First string

	// Second string in the pair
	Second string

	// ID is an optional identifier for the pair
	ID string
}

// CostWeights defines the costs for edit operations
type CostWeights struct {
	// Insertion cost (default: 1)
	Insertion int

	// Deletion cost (default: 1)
	Deletion int

	// Substitution cost (default: 1)
	Substitution int

	// Transposition cost (default: 1, only for Damerau-Levenshtein)
	Transposition int

	// Custom allows defining custom costs for specific rune pairs
	// Map structure: source rune -> target rune -> cost
	Custom map[rune]map[rune]int
}

// DefaultCostWeights returns the default cost weights (all operations cost 1)
func DefaultCostWeights() CostWeights {
	return CostWeights{
		Insertion:     1,
		Deletion:      1,
		Substitution:  1,
		Transposition: 1,
		Custom:        nil,
	}
}

// Validate checks if the cost weights are valid
func (w CostWeights) Validate() error {
	if w.Insertion < 0 {
		return NewConfigurationError("Insertion", w.Insertion, "cost must be non-negative")
	}
	if w.Deletion < 0 {
		return NewConfigurationError("Deletion", w.Deletion, "cost must be non-negative")
	}
	if w.Substitution < 0 {
		return NewConfigurationError("Substitution", w.Substitution, "cost must be non-negative")
	}
	if w.Transposition < 0 {
		return NewConfigurationError("Transposition", w.Transposition, "cost must be non-negative")
	}
	return nil
}
