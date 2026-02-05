package edlib

import (
	"context"
	"time"

	"github.com/hbollon/go-edlib/v2/algorithms"
	"github.com/hbollon/go-edlib/v2/internal/core"
)

// Builder is the main fluent API builder for string comparison operations
type Builder struct {
	config *Config
	err    error
}

// New creates a new Builder with default configuration
func New() *Builder {
	return &Builder{
		config: DefaultConfig(),
	}
}

// Quick returns a QuickBuilder for simple one-line operations
func Quick() *QuickBuilder {
	return &QuickBuilder{
		config: DefaultConfig(),
	}
}

// Configuration methods

// Using sets the algorithm to use
func (b *Builder) Using(algo algorithms.Algorithm) *Builder {
	b.config.Algorithm = algo
	return b
}

// WithContext sets the context for cancellation
func (b *Builder) WithContext(ctx context.Context) *Builder {
	b.config.Context = ctx
	return b
}

// WithTimeout sets the timeout for operations
func (b *Builder) WithTimeout(timeout time.Duration) *Builder {
	b.config.Timeout = timeout
	return b
}

// WithThreshold sets the similarity threshold
func (b *Builder) WithThreshold(threshold float64) *Builder {
	b.config.Threshold = threshold
	return b
}

// WithCaseSensitive sets whether comparisons are case-sensitive
func (b *Builder) WithCaseSensitive(sensitive bool) *Builder {
	b.config.CaseSensitive = sensitive
	return b
}

// WithNormalize enables string normalization
func (b *Builder) WithNormalize(normalize bool) *Builder {
	b.config.Normalize = normalize
	return b
}

// WithWeights sets the cost weights for edit operations
func (b *Builder) WithWeights(weights CostWeights) *Builder {
	b.config.Weights = weights
	return b
}

// WithParallel enables parallel processing with the specified number of workers
func (b *Builder) WithParallel(workers int) *Builder {
	b.config.Parallel = true
	b.config.Workers = workers
	return b
}

// WithMetrics enables detailed execution metrics
func (b *Builder) WithMetrics(enable bool) *Builder {
	b.config.EnableMetrics = enable
	return b
}

// WithCache enables result caching with the specified size
func (b *Builder) WithCache(size int) *Builder {
	b.config.EnableCache = true
	b.config.CacheSize = size
	return b
}

// WithMemoryOptimization enables memory pooling
func (b *Builder) WithMemoryOptimization(enable bool) *Builder {
	b.config.UseMemoryPool = enable
	return b
}

// WithEarlyTermination enables early termination optimization
func (b *Builder) WithEarlyTermination(enable bool) *Builder {
	b.config.EarlyTermination = enable
	return b
}

// WithMaxDistance sets the maximum distance for early termination
func (b *Builder) WithMaxDistance(maxDist int) *Builder {
	b.config.MaxDistance = maxDist
	return b
}

// Algorithm-specific builders

// EditDistance returns an EditDistanceBuilder for edit distance algorithms
func (b *Builder) EditDistance() *EditDistanceBuilder {
	return &EditDistanceBuilder{builder: b}
}

// Similarity returns a SimilarityBuilder for similarity algorithms
func (b *Builder) UsingSimilarity() *SimilarityBuilder {
	return &SimilarityBuilder{builder: b}
}

// Search returns a SearchBuilder for fuzzy search operations
func (b *Builder) Search() *SearchBuilder {
	return &SearchBuilder{builder: b}
}

// Compare executes a comparison between two strings using the configured algorithm
func (b *Builder) Compare(s1, s2 string) (*ComparisonResult, error) {
	if b.err != nil {
		return nil, b.err
	}

	if err := b.config.Validate(); err != nil {
		return nil, err
	}

	// Create calculator from registry
	calc, err := core.CreateCalculator(b.config.Algorithm, b.config.ToCalculatorOptions()...)
	if err != nil {
		return nil, err
	}

	// Setup context with timeout if needed
	ctx, cancel := b.config.WithTimeout()
	defer cancel()

	// Start metrics if enabled
	var metrics *core.Metrics
	if b.config.EnableMetrics {
		metrics = core.NewMetrics()
		defer metrics.Stop()
	}

	// Compute distance and similarity
	var distance int
	var similarity float64

	if b.config.Context != nil || b.config.Timeout > 0 {
		distance, err = calc.DistanceWithContext(ctx, s1, s2)
		if err != nil {
			return nil, err
		}
		similarity, err = calc.SimilarityWithContext(ctx, s1, s2)
		if err != nil {
			return nil, err
		}
	} else {
		distance = calc.Distance(s1, s2)
		similarity = calc.Similarity(s1, s2)
	}

	result := &ComparisonResult{
		Distance:   distance,
		Similarity: similarity,
		Algorithm:  b.config.Algorithm,
		Metadata: ComparisonMetadata{
			CaseSensitive:     b.config.CaseSensitive,
			NormalizedStrings: b.config.Normalize,
			MaxDistance:       calc.MaxDistance(s1, s2),
		},
	}

	if metrics != nil {
		result.Metrics = &ExecutionMetrics{
			ExecutionTime:     metrics.ExecutionTime,
			MemoryAllocated:   metrics.MemoryAllocated,
			MemoryFreed:       metrics.MemoryFreed,
			CacheHits:         metrics.CacheHits,
			CacheMisses:       metrics.CacheMisses,
			PoolHits:          metrics.PoolHits,
			PoolMisses:        metrics.PoolMisses,
			OptimizationsUsed: metrics.OptimizationsUsed,
			WorkersUsed:       metrics.WorkersUsed,
		}
	}

	return result, nil
}

// Distance computes only the edit distance between two strings
func (b *Builder) Distance(s1, s2 string) (int, error) {
	result, err := b.Compare(s1, s2)
	if err != nil {
		return 0, err
	}
	return result.Distance, nil
}

// Similarity computes only the similarity score between two strings
func (b *Builder) Similarity(s1, s2 string) (float64, error) {
	result, err := b.Compare(s1, s2)
	if err != nil {
		return 0, err
	}
	return result.Similarity, nil
}

// Batch processes multiple string pairs in parallel
func (b *Builder) Batch(pairs []StringPair) ([]ComparisonResult, error) {
	if b.err != nil {
		return nil, b.err
	}

	if err := b.config.Validate(); err != nil {
		return nil, err
	}

	results := make([]ComparisonResult, len(pairs))

	if b.config.Parallel && len(pairs) > 1 {
		// Parallel processing
		return b.batchParallel(pairs)
	}

	// Sequential processing
	for i, pair := range pairs {
		result, err := b.Compare(pair.First, pair.Second)
		if err != nil {
			results[i] = ComparisonResult{Error: err}
		} else {
			results[i] = *result
		}
	}

	return results, nil
}

// batchParallel processes pairs in parallel
func (b *Builder) batchParallel(pairs []StringPair) ([]ComparisonResult, error) {
	// Implementation will be in parallel processing phase
	return b.Batch(pairs) // Fallback to sequential for now
}

// Config returns the current configuration
func (b *Builder) Config() *Config {
	return b.config.Clone()
}
