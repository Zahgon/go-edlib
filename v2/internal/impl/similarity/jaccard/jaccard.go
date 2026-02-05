package jaccard

import (
	"context"
	"strings"

	"github.com/hbollon/go-edlib/v2/algorithms"
	"github.com/hbollon/go-edlib/v2/internal/core"
	"github.com/hbollon/go-edlib/v2/internal/utils"
)

// Calculator implements Jaccard similarity coefficient
type Calculator struct {
	caseInsensitive bool
	enableMetrics   bool
	metrics         *core.Metrics
	ngramSize       int // 0 = split on whitespace, >0 = use shingles
}

// New creates a new Jaccard calculator
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

// Distance returns 1 - Similarity (to fit the Distance interface)
func (c *Calculator) Distance(s1, s2 string) int {
	sim := c.Similarity(s1, s2)
	return int((1.0 - sim) * 100) // Scale to 0-100
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

// Similarity calculates the Jaccard similarity coefficient
func (c *Calculator) Similarity(s1, s2 string) float64 {
	if c.caseInsensitive {
		s1 = strings.ToLower(s1)
		s2 = strings.ToLower(s2)
	}

	if s1 == "" || s2 == "" {
		return 0.0
	}

	var splittedStr1, splittedStr2 []string
	if c.ngramSize == 0 {
		// Split on whitespaces
		splittedStr1 = strings.Split(s1, " ")
		splittedStr2 = strings.Split(s2, " ")
	} else {
		// Use shingle algorithm
		splittedStr1 = utils.ShingleSlice(s1, c.ngramSize)
		splittedStr2 = utils.ShingleSlice(s2, c.ngramSize)
	}

	// Convert splitted strings into rune arrays
	runeStr1 := make([][]rune, len(splittedStr1))
	for i, str := range splittedStr1 {
		runeStr1[i] = []rune(str)
	}
	runeStr2 := make([][]rune, len(splittedStr2))
	for i, str := range splittedStr2 {
		runeStr2[i] = []rune(str)
	}

	// Create union keywords slice between input strings
	unionStr := c.union(splittedStr1, splittedStr2)
	jacc := float64(len(runeStr1) + len(runeStr2) - len(unionStr))

	return jacc / float64(len(unionStr))
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

// union creates a union of two string slices and returns as rune matrix
func (c *Calculator) union(a, b []string) [][]rune {
	m := make(map[string]bool)
	for _, item := range a {
		m[item] = true
	}
	for _, item := range b {
		if _, ok := m[item]; !ok {
			a = append(a, item)
		}
	}

	// Convert to rune matrix
	out := make([][]rune, len(a))
	for i, word := range a {
		out[i] = []rune(word)
	}
	return out
}

// MaxDistance returns 100 (scaled max distance)
func (c *Calculator) MaxDistance(s1, s2 string) int {
	return 100
}

// Name returns the algorithm name
func (c *Calculator) Name() string {
	return "Jaccard"
}

// Properties returns the algorithm properties
func (c *Calculator) Properties() algorithms.AlgorithmProperties {
	return algorithms.Jaccard.Properties()
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

// SetNgramSize sets the n-gram size (0 = split on whitespace)
func (c *Calculator) SetNgramSize(size int) error {
	c.ngramSize = size
	return nil
}
