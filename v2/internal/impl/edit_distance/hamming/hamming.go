package hamming

import (
	"context"
	"fmt"
	"strings"

	"github.com/hbollon/go-edlib/v2/algorithms"
	"github.com/hbollon/go-edlib/v2/internal/core"
	"github.com/hbollon/go-edlib/v2/internal/utils"
)

// Calculator implements the Hamming distance algorithm
// Hamming distance only works on strings of equal length
type Calculator struct {
	caseInsensitive bool
	enableMetrics   bool
	metrics         *core.Metrics
}

// New creates a new Hamming calculator with default settings
func New(opts ...core.CalculatorOption) (*Calculator, error) {
	c := &Calculator{
		caseInsensitive: false,
		enableMetrics:   false,
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Distance computes the Hamming distance between two strings
// Returns error if strings are not of equal length
func (c *Calculator) Distance(s1, s2 string) int {
	// Normalize case if needed
	if c.caseInsensitive {
		s1 = strings.ToLower(s1)
		s2 = strings.ToLower(s2)
	}

	// Convert to runes for Unicode support
	r1, r2 := []rune(s1), []rune(s2)

	// Hamming distance requires equal length
	if len(r1) != len(r2) {
		// Return max possible distance as error indicator
		// The actual error will be available through DistanceWithContext
		return utils.Max(len(r1), len(r2))
	}

	// Quick check for identical strings
	if utils.Equal(r1, r2) {
		return 0
	}

	// Count differing positions
	distance := 0
	for i := 0; i < len(r1); i++ {
		if r1[i] != r2[i] {
			distance++
		}
	}

	return distance
}

// DistanceWithContext computes distance with context for cancellation
// Returns error if strings are not of equal length
func (c *Calculator) DistanceWithContext(ctx context.Context, s1, s2 string) (int, error) {
	// Check context before starting
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	// Normalize case if needed
	if c.caseInsensitive {
		s1 = strings.ToLower(s1)
		s2 = strings.ToLower(s2)
	}

	// Convert to runes for Unicode support
	r1, r2 := []rune(s1), []rune(s2)

	// Hamming distance requires equal length
	if len(r1) != len(r2) {
		return 0, fmt.Errorf("hamming distance undefined for strings of unequal length: %d != %d", len(r1), len(r2))
	}

	// Quick check for identical strings
	if utils.Equal(r1, r2) {
		return 0, nil
	}

	// Count differing positions
	distance := 0
	for i := 0; i < len(r1); i++ {
		if r1[i] != r2[i] {
			distance++
		}
	}

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		return distance, nil
	}
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
// For Hamming, this is the length of the strings (if equal)
func (c *Calculator) MaxDistance(s1, s2 string) int {
	len1, len2 := len([]rune(s1)), len([]rune(s2))
	
	// For equal length strings, max distance is the length
	if len1 == len2 {
		return len1
	}
	
	// For unequal lengths, return the max (even though it's technically undefined)
	return utils.Max(len1, len2)
}

// Name returns the calculator name
func (c *Calculator) Name() string {
	return "Hamming"
}

// Properties returns the algorithm properties
func (c *Calculator) Properties() algorithms.AlgorithmProperties {
	return algorithms.Hamming.Properties()
}

// Ensure Calculator implements the interfaces
var (
	_ core.Calculator = (*Calculator)(nil)
)
