package jaro

import (
	"context"
	"strings"

	"github.com/hbollon/go-edlib/v2/algorithms"
	"github.com/hbollon/go-edlib/v2/internal/core"
	"github.com/hbollon/go-edlib/v2/internal/utils"
)

// Calculator implements both Jaro and Jaro-Winkler similarity algorithms
type Calculator struct {
	useWinkler       bool
	prefixScale      float64
	caseInsensitive  bool
	pool             *core.MemoryPool
	enableMetrics    bool
	metrics          *core.Metrics
}

const (
	defaultPrefixScale = 0.1
	maxPrefixLength    = 4
)

// New creates a new Jaro calculator
func New(opts ...core.CalculatorOption) (*Calculator, error) {
	c := &Calculator{
		useWinkler:      false,
		prefixScale:     defaultPrefixScale,
		caseInsensitive: false,
		pool:            nil,
		enableMetrics:   false,
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// NewWinkler creates a new Jaro-Winkler calculator
func NewWinkler(opts ...core.CalculatorOption) (*Calculator, error) {
	c, err := New(opts...)
	if err != nil {
		return nil, err
	}
	c.useWinkler = true
	return c, nil
}

// Distance returns 1 - similarity (to implement DistanceCalculator interface)
func (c *Calculator) Distance(s1, s2 string) int {
	similarity := c.Similarity(s1, s2)
	// Convert similarity [0..1] to distance
	return int((1.0 - similarity) * 100)
}

// DistanceWithContext computes distance with context
func (c *Calculator) DistanceWithContext(ctx context.Context, s1, s2 string) (int, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	dist := c.Distance(s1, s2)

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		return dist, nil
	}
}

// Similarity computes the Jaro or Jaro-Winkler similarity
func (c *Calculator) Similarity(s1, s2 string) float64 {
	if c.caseInsensitive {
		s1 = strings.ToLower(s1)
		s2 = strings.ToLower(s2)
	}

	r1, r2 := []rune(s1), []rune(s2)

	if len(r1) == 0 && len(r2) == 0 {
		return 1.0
	}
	if len(r1) == 0 || len(r2) == 0 {
		return 0.0
	}
	if utils.Equal(r1, r2) {
		return 1.0
	}

	jaroSim := c.jaroSimilarity(r1, r2)

	if !c.useWinkler || jaroSim == 0.0 || jaroSim == 1.0 {
		return jaroSim
	}

	// Apply Winkler modification
	return c.applyWinklerBonus(r1, r2, jaroSim)
}

// jaroSimilarity computes the Jaro similarity
func (c *Calculator) jaroSimilarity(r1, r2 []rune) float64 {
	len1, len2 := len(r1), len(r2)

	// Maximum distance for matching
	maxDist := utils.Max(len1, len2)/2 - 1
	if maxDist < 0 {
		maxDist = 0
	}

	// Get match tables from pool or allocate
	var matches1, matches2 []bool
	if c.pool != nil {
		matches1 = c.pool.GetBoolSlice(len1)
		matches2 = c.pool.GetBoolSlice(len2)
		defer func() {
			c.pool.PutBoolSlice(matches1)
			c.pool.PutBoolSlice(matches2)
		}()
		if c.metrics != nil {
			c.metrics.RecordPoolHit()
		}
	} else {
		matches1 = make([]bool, len1)
		matches2 = make([]bool, len2)
		if c.metrics != nil {
			c.metrics.RecordPoolMiss()
		}
	}

	// Find matches
	matches := 0
	for i := 0; i < len1; i++ {
		start := utils.Max(0, i-maxDist)
		end := utils.Min(len2, i+maxDist+1)

		for j := start; j < end; j++ {
			if matches2[j] || r1[i] != r2[j] {
				continue
			}
			matches1[i] = true
			matches2[j] = true
			matches++
			break
		}
	}

	if matches == 0 {
		return 0.0
	}

	// Count transpositions
	transpositions := 0
	k := 0
	for i := 0; i < len1; i++ {
		if !matches1[i] {
			continue
		}
		for !matches2[k] {
			k++
		}
		if r1[i] != r2[k] {
			transpositions++
		}
		k++
	}

	m := float64(matches)
	return (m/float64(len1) + m/float64(len2) + (m-float64(transpositions)/2.0)/m) / 3.0
}

// applyWinklerBonus applies the Jaro-Winkler prefix bonus
func (c *Calculator) applyWinklerBonus(r1, r2 []rune, jaroSim float64) float64 {
	// Find common prefix length (max 4)
	prefixLen := 0
	minLen := utils.Min(len(r1), len(r2))
	for i := 0; i < minLen && i < maxPrefixLength; i++ {
		if r1[i] == r2[i] {
			prefixLen++
		} else {
			break
		}
	}

	return jaroSim + float64(prefixLen)*c.prefixScale*(1.0-jaroSim)
}

// SimilarityWithContext computes similarity with context
func (c *Calculator) SimilarityWithContext(ctx context.Context, s1, s2 string) (float64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	sim := c.Similarity(s1, s2)

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		return sim, nil
	}
}

// MaxDistance returns the maximum possible distance
func (c *Calculator) MaxDistance(s1, s2 string) int {
	return 100 // Since we use percentage-based distance
}

// Name returns the calculator name
func (c *Calculator) Name() string {
	if c.useWinkler {
		return "JaroWinkler"
	}
	return "Jaro"
}

// Properties returns the algorithm properties
func (c *Calculator) Properties() algorithms.AlgorithmProperties {
	if c.useWinkler {
		return algorithms.JaroWinkler.Properties()
	}
	return algorithms.Jaro.Properties()
}

// WithThreshold is a no-op for Jaro (doesn't support threshold optimization)
func (c *Calculator) WithThreshold(threshold float64) core.OptimizedCalculator {
	return c
}

// WithMemoryPool enables memory pooling
func (c *Calculator) WithMemoryPool(pool *core.MemoryPool) core.OptimizedCalculator {
	c.pool = pool
	return c
}

// WithEarlyTermination is a no-op for Jaro
func (c *Calculator) WithEarlyTermination(enabled bool) core.OptimizedCalculator {
	return c
}

// WithMaxDistance is a no-op for Jaro
func (c *Calculator) WithMaxDistance(maxDist int) core.OptimizedCalculator {
	return c
}

// Ensure Calculator implements interfaces
var (
	_ core.Calculator          = (*Calculator)(nil)
	_ core.OptimizedCalculator = (*Calculator)(nil)
)
