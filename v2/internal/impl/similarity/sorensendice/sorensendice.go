package sorensendice

import (
	"context"
	"strings"

	"github.com/hbollon/go-edlib/v2/algorithms"
	"github.com/hbollon/go-edlib/v2/internal/core"
	"github.com/hbollon/go-edlib/v2/internal/utils"
)

// Calculator implements Sorensen-Dice coefficient
type Calculator struct {
	caseInsensitive bool
	enableMetrics   bool
	metrics         *core.Metrics
	ngramSize       int
}

// New creates a new Sorensen-Dice calculator
func New(opts ...core.CalculatorOption) (*Calculator, error) {
	c := &Calculator{
		ngramSize: 2, // Default 2-gram
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	if c.enableMetrics {
		c.metrics = core.NewMetrics()
	}

	return c, nil
}

// Distance returns 1 - Similarity (scaled to 0-100)
func (c *Calculator) Distance(s1, s2 string) int {
	sim := c.Similarity(s1, s2)
	return int((1.0 - sim) * 100)
}

// DistanceWithContext calculates distance with context support
func (c *Calculator) DistanceWithContext(ctx context.Context, s1, s2 string) (int, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		return c.Distance(s1, s2), nil
	}
}

// Similarity calculates the Sorensen-Dice coefficient
func (c *Calculator) Similarity(s1, s2 string) float64 {
	if c.caseInsensitive {
		s1 = strings.ToLower(s1)
		s2 = strings.ToLower(s2)
	}

	if s1 == "" && s2 == "" {
		return 0.0
	}

	shingle1 := utils.Shingle(s1, c.ngramSize)
	shingle2 := utils.Shingle(s2, c.ngramSize)

	intersection := float64(0)
	for i := range shingle1 {
		if _, ok := shingle2[i]; ok {
			intersection++
		}
	}

	return 2.0 * intersection / float64(len(shingle1)+len(shingle2))
}

// SimilarityWithContext calculates similarity with context support
func (c *Calculator) SimilarityWithContext(ctx context.Context, s1, s2 string) (float64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		return c.Similarity(s1, s2), nil
	}
}

// MaxDistance returns 100 (scaled max distance)
func (c *Calculator) MaxDistance(s1, s2 string) int {
	return 100
}

// Name returns the algorithm name
func (c *Calculator) Name() string {
	return "SorensenDice"
}

// Properties returns the algorithm properties
func (c *Calculator) Properties() algorithms.AlgorithmProperties {
	return algorithms.SorensenDice.Properties()
}

// Metrics returns the collected metrics if enabled
func (c *Calculator) Metrics() *core.Metrics {
	return c.metrics
}

// SetCaseInsensitive sets whether comparisons should be case-insensitive
func (c *Calculator) SetCaseInsensitive(v bool) error {
	c.caseInsensitive = v
	return nil
}

// SetMetricsEnabled enables or disables metrics collection
func (c *Calculator) SetMetricsEnabled(v bool) error {
	c.enableMetrics = v
	if v && c.metrics == nil {
		c.metrics = core.NewMetrics()
	}
	return nil
}

// SetMemoryPool is a no-op for this algorithm
func (c *Calculator) SetMemoryPool(pool *core.MemoryPool) error {
	return nil
}

// SetNgramSize sets the n-gram size
func (c *Calculator) SetNgramSize(size int) error {
	c.ngramSize = size
	return nil
}
