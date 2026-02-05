package damerau

import (
	"context"
	"strings"

	"github.com/hbollon/go-edlib/v2/algorithms"
	"github.com/hbollon/go-edlib/v2/internal/core"
	"github.com/hbollon/go-edlib/v2/internal/utils"
)

// Calculator implements both OSA and full Damerau-Levenshtein distance algorithms
type Calculator struct {
	caseInsensitive bool
	enableMetrics   bool
	metrics         *core.Metrics
	useOSA          bool // true for OSA variant, false for full Damerau-Levenshtein
}

// New creates a new Damerau-Levenshtein calculator
func New(opts ...core.CalculatorOption) (*Calculator, error) {
	c := &Calculator{
		useOSA: false, // Default to full Damerau-Levenshtein
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

// NewOSA creates a new OSA Damerau-Levenshtein calculator
func NewOSA(opts ...core.CalculatorOption) (*Calculator, error) {
	c := &Calculator{
		useOSA: true,
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

// Distance calculates the Damerau-Levenshtein distance between two strings
func (c *Calculator) Distance(s1, s2 string) int {
	if c.caseInsensitive {
		s1 = strings.ToLower(s1)
		s2 = strings.ToLower(s2)
	}

	if c.useOSA {
		return c.calculateOSA(s1, s2)
	}
	return c.calculateFull(s1, s2)
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

// Similarity calculates similarity as 1 - (distance / maxPossibleDistance)
func (c *Calculator) Similarity(s1, s2 string) float64 {
	dist := c.Distance(s1, s2)
	maxDist := c.MaxDistance(s1, s2)
	if maxDist == 0 {
		return 1.0
	}
	return 1.0 - float64(dist)/float64(maxDist)
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

// MaxDistance returns the maximum possible distance (length of longer string)
func (c *Calculator) MaxDistance(s1, s2 string) int {
	if c.caseInsensitive {
		s1 = strings.ToLower(s1)
		s2 = strings.ToLower(s2)
	}
	return utils.Max(len([]rune(s1)), len([]rune(s2)))
}

// Name returns the algorithm name
func (c *Calculator) Name() string {
	if c.useOSA {
		return "OSADamerauLevenshtein"
	}
	return "DamerauLevenshtein"
}

// Properties returns the algorithm properties
func (c *Calculator) Properties() algorithms.AlgorithmProperties {
	if c.useOSA {
		return algorithms.OSADamerauLevenshtein.Properties()
	}
	return algorithms.DamerauLevenshtein.Properties()
}

// Metrics returns the collected metrics if enabled
func (c *Calculator) Metrics() *core.Metrics {
	return c.metrics
}

// calculateOSA computes OSA Damerau-Levenshtein distance
// Optimal string alignment distance variant using Wagner-Fisher dynamic programming
// Doesn't allow multiple transformations on the same substring
func (c *Calculator) calculateOSA(str1, str2 string) int {
	// Convert string parameters to rune arrays to be compatible with non-ASCII
	runeStr1 := []rune(str1)
	runeStr2 := []rune(str2)

	// Get and store length of these strings
	runeStr1len := len(runeStr1)
	runeStr2len := len(runeStr2)
	
	if runeStr1len == 0 {
		return runeStr2len
	} else if runeStr2len == 0 {
		return runeStr1len
	} else if utils.Equal(runeStr1, runeStr2) {
		return 0
	} else if runeStr1len < runeStr2len {
		return c.calculateOSA(str2, str1)
	}

	// 2D Array with memory optimization (only need 3 rows)
	row := utils.Min(runeStr1len+1, 3)
	matrix := make([][]int, row)
	for i := 0; i < row; i++ {
		matrix[i] = make([]int, runeStr2len+1)
		matrix[i][0] = i
	}

	for j := 0; j <= runeStr2len; j++ {
		matrix[0][j] = j
	}

	var count int
	for i := 1; i <= runeStr1len; i++ {
		matrix[i%3][0] = i
		for j := 1; j <= runeStr2len; j++ {
			if runeStr1[i-1] == runeStr2[j-1] {
				count = 0
			} else {
				count = 1
			}

			matrix[i%3][j] = utils.Min(utils.Min(matrix[(i-1)%3][j]+1, matrix[i%3][j-1]+1),
				matrix[(i-1)%3][j-1]+count) // insertion, deletion, substitution
			if i > 1 && j > 1 && runeStr1[i-1] == runeStr2[j-2] && runeStr1[i-2] == runeStr2[j-1] {
				matrix[i%3][j] = utils.Min(matrix[i%3][j], matrix[(i-2)%3][j-2]+1) // transposition
			}
		}
	}
	return matrix[runeStr1len%3][runeStr2len]
}

// calculateFull computes true Damerau-Levenshtein distance with adjacent transpositions
func (c *Calculator) calculateFull(str1, str2 string) int {
	// Convert string parameters to rune arrays to be compatible with non-ASCII
	runeStr1 := []rune(str1)
	runeStr2 := []rune(str2)

	// Get and store length of these strings
	runeStr1len := len(runeStr1)
	runeStr2len := len(runeStr2)
	
	if runeStr1len == 0 {
		return runeStr2len
	} else if runeStr2len == 0 {
		return runeStr1len
	} else if utils.Equal(runeStr1, runeStr2) {
		return 0
	}

	// Create alphabet based on input strings
	da := make(map[rune]int)
	for i := 0; i < runeStr1len; i++ {
		da[runeStr1[i]] = 0
	}
	for i := 0; i < runeStr2len; i++ {
		da[runeStr2[i]] = 0
	}

	// 2D Array for distance matrix : matrix[0..str1.length+2][0..s2.length+2]
	matrix := make([][]int, runeStr1len+2)
	for i := 0; i <= runeStr1len+1; i++ {
		matrix[i] = make([]int, runeStr2len+2)
		for j := 0; j <= runeStr2len+1; j++ {
			matrix[i][j] = 0
		}
	}

	// Maximum possible distance
	maxDist := runeStr1len + runeStr2len

	// Initialize matrix
	matrix[0][0] = maxDist
	for i := 0; i <= runeStr1len; i++ {
		matrix[i+1][0] = maxDist
		matrix[i+1][1] = i
	}
	for i := 0; i <= runeStr2len; i++ {
		matrix[0][i+1] = maxDist
		matrix[1][i+1] = i
	}

	// Process edit distance
	var cost int
	for i := 1; i <= runeStr1len; i++ {
		db := 0
		for j := 1; j <= runeStr2len; j++ {
			i1 := da[runeStr2[j-1]]
			j1 := db
			if runeStr1[i-1] == runeStr2[j-1] {
				cost = 0
				db = j
			} else {
				cost = 1
			}

			matrix[i+1][j+1] = utils.Min(
				utils.Min(
					matrix[i+1][j]+1,  // Addition
					matrix[i][j+1]+1), // Deletion
				utils.Min(
					matrix[i][j]+cost, // Substitution
					matrix[i1][j1]+(i-i1-1)+1+(j-j1-1))) // Transposition
		}

		da[runeStr1[i-1]] = i
	}

	return matrix[runeStr1len+1][runeStr2len+1]
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

// SetMemoryPool is a no-op for this algorithm (doesn't use pooling)
func (c *Calculator) SetMemoryPool(pool *core.MemoryPool) error {
	return nil
}
