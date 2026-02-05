package levenshtein

import (
	"context"
	"strings"

	"github.com/hbollon/go-edlib/v2/algorithms"
	"github.com/hbollon/go-edlib/v2/internal/core"
	"github.com/hbollon/go-edlib/v2/internal/utils"
)

// Calculator implements the Levenshtein distance algorithm
type Calculator struct {
	threshold        float64
	caseInsensitive  bool
	weights          weights
	pool             *core.MemoryPool
	earlyTermination bool
	maxDistance      int
	enableMetrics    bool
	metrics          *core.Metrics
}

type weights struct {
	insertion     int
	deletion      int
	substitution  int
	hasCustom     bool
	custom        map[rune]map[rune]int
}

// New creates a new Levenshtein calculator with default settings
func New(opts ...core.CalculatorOption) (*Calculator, error) {
	c := &Calculator{
		threshold:        0.0,
		caseInsensitive:  false,
		weights: weights{
			insertion:    1,
			deletion:     1,
			substitution: 1,
			hasCustom:    false,
		},
		pool:             nil,
		earlyTermination: false,
		maxDistance:      0,
		enableMetrics:    false,
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Distance computes the Levenshtein distance between two strings
func (c *Calculator) Distance(s1, s2 string) int {
	// Normalize case if needed
	if c.caseInsensitive {
		s1 = strings.ToLower(s1)
		s2 = strings.ToLower(s2)
	}

	// Convert to runes for Unicode support
	r1, r2 := []rune(s1), []rune(s2)

	// Quick checks
	if len(r1) == 0 {
		return len(r2) * c.weights.insertion
	}
	if len(r2) == 0 {
		return len(r1) * c.weights.deletion
	}
	if utils.Equal(r1, r2) {
		return 0
	}

	// Use optimized two-row implementation
	return c.distanceOptimized(r1, r2)
}

// DistanceWithContext computes distance with context for cancellation
func (c *Calculator) DistanceWithContext(ctx context.Context, s1, s2 string) (int, error) {
	// Check context before starting
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	// For now, compute synchronously and check context after
	// In production, we'd check context periodically during computation
	dist := c.Distance(s1, s2)

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		return dist, nil
	}
}

// distanceOptimized uses two-row algorithm for O(min(n,m)) space complexity
// distanceOptimized uses two-row algorithm for O(min(n,m)) space complexity
func (c *Calculator) distanceOptimized(r1, r2 []rune) int {
	// Ensure r1 is the shorter string for space optimization
	swapped := false
	if len(r1) > len(r2) {
		r1, r2 = r2, r1
		swapped = true
	}
	
	// Get weights (swap if strings were swapped)
	insertionWeight := c.weights.insertion
	deletionWeight := c.weights.deletion
	if swapped {
		insertionWeight, deletionWeight = c.weights.deletion, c.weights.insertion
	}

	// Get or allocate rows
	var prevRow, currRow []int
	if c.pool != nil {
		prevRow = c.pool.GetIntSlice(len(r1) + 1)
		currRow = c.pool.GetIntSlice(len(r1) + 1)
		defer func() {
			c.pool.PutIntSlice(prevRow)
			c.pool.PutIntSlice(currRow)
		}()
		if c.metrics != nil {
			c.metrics.RecordPoolHit()
		}
	} else {
		prevRow = make([]int, len(r1)+1)
		currRow = make([]int, len(r1)+1)
		if c.metrics != nil {
			c.metrics.RecordPoolMiss()
		}
	}

	// Initialize first row
	for i := 0; i <= len(r1); i++ {
		prevRow[i] = i * deletionWeight
	}

	// Compute distances
	for i := 1; i <= len(r2); i++ {
		currRow[0] = i * insertionWeight

		minInRow := currRow[0]

		for j := 1; j <= len(r1); j++ {
			cost := c.getCost(r1[j-1], r2[i-1])

			currRow[j] = utils.Min3(
				currRow[j-1]+deletionWeight,
				prevRow[j]+insertionWeight,
				prevRow[j-1]+cost,
			)

			if currRow[j] < minInRow {
				minInRow = currRow[j]
			}
		}

		// Early termination check
		if c.earlyTermination && c.shouldTerminateEarly(minInRow, len(r1), len(r2)) {
			if c.metrics != nil {
				c.metrics.AddOptimization("early_termination")
			}
			return c.estimateFinalDistance(minInRow, i, len(r2))
		}

		// Swap rows
		prevRow, currRow = currRow, prevRow
	}

	return prevRow[len(r1)]
}

// getCost returns the cost for substituting r1 with r2
func (c *Calculator) getCost(r1, r2 rune) int {
	if r1 == r2 {
		return 0
	}

	// Check custom weights if available
	if c.weights.hasCustom {
		if costs, ok := c.weights.custom[r1]; ok {
			if cost, ok := costs[r2]; ok {
				return cost
			}
		}
	}

	return c.weights.substitution
}

// shouldTerminateEarly checks if we should stop early based on threshold
func (c *Calculator) shouldTerminateEarly(minInRow, len1, len2 int) bool {
	if c.maxDistance > 0 && minInRow > c.maxDistance {
		return true
	}

	if c.threshold > 0 {
		maxLen := utils.Max(len1, len2)
		threshold := int(c.threshold * float64(maxLen))
		return minInRow > threshold
	}

	return false
}

// estimateFinalDistance estimates the final distance when terminating early
func (c *Calculator) estimateFinalDistance(minInRow, currentRow, totalRows int) int {
	remainingRows := totalRows - currentRow
	return minInRow + remainingRows
}

// Similarity computes the similarity score [0.0 to 1.0]
func (c *Calculator) Similarity(s1, s2 string) float64 {
	dist := c.Distance(s1, s2)
	maxDist := c.MaxDistance(s1, s2)

	if maxDist == 0 {
		return 1.0
	}

	similarity := 1.0 - (float64(dist) / float64(maxDist))
	return utils.ClampFloat64(similarity, 0.0, 1.0)
}

// SimilarityWithContext computes similarity with context for cancellation
func (c *Calculator) SimilarityWithContext(ctx context.Context, s1, s2 string) (float64, error) {
	dist, err := c.DistanceWithContext(ctx, s1, s2)
	if err != nil {
		return 0, err
	}

	maxDist := c.MaxDistance(s1, s2)
	if maxDist == 0 {
		return 1.0, nil
	}

	similarity := 1.0 - (float64(dist) / float64(maxDist))
	return utils.ClampFloat64(similarity, 0.0, 1.0), nil
}

// MaxDistance returns the maximum possible distance
func (c *Calculator) MaxDistance(s1, s2 string) int {
	len1, len2 := len([]rune(s1)), len([]rune(s2))
	
	if len1 > len2 {
		return len1 * utils.Max(c.weights.deletion, c.weights.insertion)
	}
	return len2 * utils.Max(c.weights.deletion, c.weights.insertion)
}

// Name returns the calculator name
func (c *Calculator) Name() string {
	return "Levenshtein"
}

// Properties returns the algorithm properties
func (c *Calculator) Properties() algorithms.AlgorithmProperties {
	return algorithms.Levenshtein.Properties()
}

// Optimization interface implementations

func (c *Calculator) WithThreshold(threshold float64) core.OptimizedCalculator {
	c.threshold = threshold
	return c
}

func (c *Calculator) WithMemoryPool(pool *core.MemoryPool) core.OptimizedCalculator {
	c.pool = pool
	return c
}

func (c *Calculator) WithEarlyTermination(enabled bool) core.OptimizedCalculator {
	c.earlyTermination = enabled
	return c
}

func (c *Calculator) WithMaxDistance(maxDist int) core.OptimizedCalculator {
	c.maxDistance = maxDist
	return c
}


// SetCaseInsensitive configures case sensitivity
func (c *Calculator) SetCaseInsensitive(insensitive bool) error {
	c.caseInsensitive = insensitive
	return nil
}

// SetMetricsEnabled enables or disables metrics collection
func (c *Calculator) SetMetricsEnabled(enabled bool) error {
	c.enableMetrics = enabled
	if enabled && c.metrics == nil {
		c.metrics = core.NewMetrics()
	}
	return nil
}

// SetMemoryPool sets the memory pool
func (c *Calculator) SetMemoryPool(pool *core.MemoryPool) error {
	c.pool = pool
	return nil
}

// Ensure Calculator implements the interfaces
var (
	_ core.Calculator          = (*Calculator)(nil)
	_ core.OptimizedCalculator = (*Calculator)(nil)
)
