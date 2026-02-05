package qgram

import (
	"context"
	"math"
	"strings"

	"github.com/hbollon/go-edlib/v2/algorithms"
	"github.com/hbollon/go-edlib/v2/internal/core"
	"github.com/hbollon/go-edlib/v2/internal/utils"
)

// Calculator implements Q-gram distance and similarity
type Calculator struct {
	caseInsensitive bool
	enableMetrics   bool
	metrics         *core.Metrics
	ngramSize       int
}

// New creates a new Q-gram calculator
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

// Distance calculates the Q-gram distance
func (c *Calculator) Distance(s1, s2 string) int {
	if c.caseInsensitive {
		s1 = strings.ToLower(s1)
		s2 = strings.ToLower(s2)
	}

	return c.qgramDistance(s1, s2)
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

// Similarity calculates Q-gram similarity (between 0 and 1)
func (c *Calculator) Similarity(s1, s2 string) float64 {
	if c.caseInsensitive {
		s1 = strings.ToLower(s1)
		s2 = strings.ToLower(s2)
	}

	splittedStr1 := utils.Shingle(s1, c.ngramSize)
	splittedStr2 := utils.Shingle(s2, c.ngramSize)
	
	res := float64(c.qgramDistanceCustom(splittedStr1, splittedStr2))
	
	totalShingles := 0
	for _, i := range splittedStr1 {
		totalShingles += i
	}
	for _, i := range splittedStr2 {
		totalShingles += i
	}
	
	if totalShingles == 0 {
		return 0.0
	}
	
	return 1.0 - (res / float64(totalShingles))
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

// qgramDistance computes Q-gram distance between two strings
func (c *Calculator) qgramDistance(str1, str2 string) int {
	splittedStr1 := utils.Shingle(str1, c.ngramSize)
	splittedStr2 := utils.Shingle(str2, c.ngramSize)
	return c.qgramDistanceCustom(splittedStr1, splittedStr2)
}

// qgramDistanceCustom computes Q-gram distance from custom n-gram maps
func (c *Calculator) qgramDistanceCustom(splittedStr1, splittedStr2 map[string]int) int {
	union := make(map[string]int)
	for i := range splittedStr1 {
		union[i] = 0
	}
	for i := range splittedStr2 {
		union[i] = 0
	}

	res := 0
	for i := range union {
		res += int(math.Abs(float64(splittedStr1[i] - splittedStr2[i])))
	}

	return res
}

// MaxDistance returns the maximum possible distance
func (c *Calculator) MaxDistance(s1, s2 string) int {
	if c.caseInsensitive {
		s1 = strings.ToLower(s1)
		s2 = strings.ToLower(s2)
	}
	// Max distance is the sum of all q-grams
	return len([]rune(s1)) + len([]rune(s2))
}

// Name returns the algorithm name
func (c *Calculator) Name() string {
	return "QGram"
}

// Properties returns the algorithm properties
func (c *Calculator) Properties() algorithms.AlgorithmProperties {
	return algorithms.QGram.Properties()
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
