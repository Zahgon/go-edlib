package cosine

import (
	"context"
	"math"
	"strings"

	"github.com/hbollon/go-edlib/v2/algorithms"
	"github.com/hbollon/go-edlib/v2/internal/core"
)

// Calculator implements the Cosine similarity algorithm
type Calculator struct {
	ngramSize       int
	caseInsensitive bool
	splitOnSpace    bool
	enableMetrics   bool
	metrics         *core.Metrics
}

const (
	defaultNgramSize = 2
)

// New creates a new Cosine similarity calculator
func New(opts ...core.CalculatorOption) (*Calculator, error) {
	c := &Calculator{
		ngramSize:       defaultNgramSize,
		caseInsensitive: false,
		splitOnSpace:    false,
		enableMetrics:   false,
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Distance returns 1 - similarity (scaled to 100)
func (c *Calculator) Distance(s1, s2 string) int {
	similarity := c.Similarity(s1, s2)
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

// Similarity computes the cosine similarity between two strings
func (c *Calculator) Similarity(s1, s2 string) float64 {
	if s1 == "" || s2 == "" {
		return 0.0
	}

	if c.caseInsensitive {
		s1 = strings.ToLower(s1)
		s2 = strings.ToLower(s2)
	}

	// Generate n-grams or word tokens
	var tokens1, tokens2 []string
	if c.splitOnSpace {
		tokens1 = strings.Fields(s1)
		tokens2 = strings.Fields(s2)
	} else {
		tokens1 = c.generateNgrams(s1)
		tokens2 = c.generateNgrams(s2)
	}

	if len(tokens1) == 0 || len(tokens2) == 0 {
		return 0.0
	}

	// Build frequency vectors
	vec1 := c.buildVector(tokens1)
	vec2 := c.buildVector(tokens2)

	// Compute cosine similarity
	return c.cosineSimilarity(vec1, vec2)
}

// generateNgrams creates n-grams from a string
func (c *Calculator) generateNgrams(s string) []string {
	runes := []rune(s)
	if len(runes) < c.ngramSize {
		return []string{s}
	}

	ngrams := make([]string, 0, len(runes)-c.ngramSize+1)
	for i := 0; i <= len(runes)-c.ngramSize; i++ {
		ngrams = append(ngrams, string(runes[i:i+c.ngramSize]))
	}

	return ngrams
}

// buildVector creates a frequency vector from tokens
func (c *Calculator) buildVector(tokens []string) map[string]int {
	vec := make(map[string]int)
	for _, token := range tokens {
		vec[token]++
	}
	return vec
}

// cosineSimilarity computes cosine similarity between two frequency vectors
func (c *Calculator) cosineSimilarity(vec1, vec2 map[string]int) float64 {
	// Build union of keys
	union := make(map[string]bool)
	for k := range vec1 {
		union[k] = true
	}
	for k := range vec2 {
		union[k] = true
	}

	// Compute dot product and magnitudes
	var dotProduct float64
	var mag1, mag2 float64

	for token := range union {
		v1 := float64(vec1[token])
		v2 := float64(vec2[token])

		dotProduct += v1 * v2
		mag1 += v1 * v1
		mag2 += v2 * v2
	}

	// Compute cosine
	magnitude := math.Sqrt(mag1) * math.Sqrt(mag2)
	if magnitude == 0 {
		return 0.0
	}

	return dotProduct / magnitude
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
	return 100
}

// Name returns the calculator name
func (c *Calculator) Name() string {
	return "Cosine"
}

// Properties returns the algorithm properties
func (c *Calculator) Properties() algorithms.AlgorithmProperties {
	return algorithms.Cosine.Properties()
}

// WithThreshold is a no-op for Cosine
func (c *Calculator) WithThreshold(threshold float64) core.OptimizedCalculator {
	return c
}

// WithMemoryPool is a no-op for Cosine
func (c *Calculator) WithMemoryPool(pool *core.MemoryPool) core.OptimizedCalculator {
	return c
}

// WithEarlyTermination is a no-op for Cosine
func (c *Calculator) WithEarlyTermination(enabled bool) core.OptimizedCalculator {
	return c
}

// WithMaxDistance is a no-op for Cosine
func (c *Calculator) WithMaxDistance(maxDist int) core.OptimizedCalculator {
	return c
}

// Ensure Calculator implements interfaces
var (
	_ core.Calculator          = (*Calculator)(nil)
	_ core.OptimizedCalculator = (*Calculator)(nil)
)
