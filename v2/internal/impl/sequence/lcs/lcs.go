package lcs

import (
	"context"
	"errors"
	"strings"

	"github.com/hbollon/go-edlib/v2/algorithms"
	"github.com/hbollon/go-edlib/v2/internal/core"
	"github.com/hbollon/go-edlib/v2/internal/utils"
)

// Calculator implements Longest Common Subsequence algorithms
type Calculator struct {
	caseInsensitive bool
	enableMetrics   bool
	metrics         *core.Metrics
	pool            *core.MemoryPool
}

// New creates a new LCS calculator
func New(opts ...core.CalculatorOption) (*Calculator, error) {
	c := &Calculator{}

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

// Distance calculates the LCS length between two strings
func (c *Calculator) Distance(s1, s2 string) int {
	if c.caseInsensitive {
		s1 = strings.ToLower(s1)
		s2 = strings.ToLower(s2)
	}

	return c.calculateLCS(s1, s2)
}

// DistanceWithContext calculates LCS with context support
func (c *Calculator) DistanceWithContext(ctx context.Context, s1, s2 string) (int, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		return c.Distance(s1, s2), nil
	}
}

// Similarity calculates similarity based on LCS
func (c *Calculator) Similarity(s1, s2 string) float64 {
	lcs := c.Distance(s1, s2)
	maxLen := c.MaxDistance(s1, s2)
	if maxLen == 0 {
		return 1.0
	}
	return float64(lcs) / float64(maxLen)
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

// MaxDistance returns the length of the longer string
func (c *Calculator) MaxDistance(s1, s2 string) int {
	if c.caseInsensitive {
		s1 = strings.ToLower(s1)
		s2 = strings.ToLower(s2)
	}
	return utils.Max(len([]rune(s1)), len([]rune(s2)))
}

// Name returns the algorithm name
func (c *Calculator) Name() string {
	return "LCS"
}

// Properties returns the algorithm properties
func (c *Calculator) Properties() algorithms.AlgorithmProperties {
	return algorithms.LCS.Properties()
}

// Metrics returns the collected metrics if enabled
func (c *Calculator) Metrics() *core.Metrics {
	return c.metrics
}

// calculateLCS computes the LCS length
func (c *Calculator) calculateLCS(str1, str2 string) int {
	runeStr1 := []rune(str1)
	runeStr2 := []rune(str2)

	if len(runeStr1) == 0 || len(runeStr2) == 0 {
		return 0
	} else if utils.Equal(runeStr1, runeStr2) {
		return len(runeStr1)
	}

	lcsMatrix := c.lcsProcess(runeStr1, runeStr2)
	return lcsMatrix[len(runeStr1)][len(runeStr2)]
}

// lcsProcess returns the computed LCS matrix
func (c *Calculator) lcsProcess(runeStr1, runeStr2 []rune) [][]int {
	lcsMatrix := make([][]int, len(runeStr1)+1)
	for i := 0; i <= len(runeStr1); i++ {
		lcsMatrix[i] = make([]int, len(runeStr2)+1)
	}

	for i := 1; i <= len(runeStr1); i++ {
		for j := 1; j <= len(runeStr2); j++ {
			if runeStr1[i-1] == runeStr2[j-1] {
				lcsMatrix[i][j] = lcsMatrix[i-1][j-1] + 1
			} else {
				lcsMatrix[i][j] = utils.Max(lcsMatrix[i][j-1], lcsMatrix[i-1][j])
			}
		}
	}

	return lcsMatrix
}

// Backtrack returns one LCS string
func (c *Calculator) Backtrack(str1, str2 string) (string, error) {
	if c.caseInsensitive {
		str1 = strings.ToLower(str1)
		str2 = strings.ToLower(str2)
	}

	runeStr1 := []rune(str1)
	runeStr2 := []rune(str2)

	if len(runeStr1) == 0 || len(runeStr2) == 0 {
		return "", errors.New("cannot process and backtrack LCS with empty string")
	} else if utils.Equal(runeStr1, runeStr2) {
		return str1, nil
	}

	lcsMatrix := c.lcsProcess(runeStr1, runeStr2)
	return c.processBacktrack(str1, str2, lcsMatrix, len(runeStr1), len(runeStr2)), nil
}

func (c *Calculator) processBacktrack(str1, str2 string, lcsMatrix [][]int, m, n int) string {
	runeStr1 := []rune(str1)
	runeStr2 := []rune(str2)

	if m == 0 || n == 0 {
		return ""
	} else if runeStr1[m-1] == runeStr2[n-1] {
		return c.processBacktrack(str1, str2, lcsMatrix, m-1, n-1) + string(runeStr1[m-1])
	} else if lcsMatrix[m][n-1] > lcsMatrix[m-1][n] {
		return c.processBacktrack(str1, str2, lcsMatrix, m, n-1)
	}

	return c.processBacktrack(str1, str2, lcsMatrix, m-1, n)
}

// BacktrackAll returns all common subsequences
func (c *Calculator) BacktrackAll(str1, str2 string) ([]string, error) {
	if c.caseInsensitive {
		str1 = strings.ToLower(str1)
		str2 = strings.ToLower(str2)
	}

	runeStr1 := []rune(str1)
	runeStr2 := []rune(str2)

	if len(runeStr1) == 0 || len(runeStr2) == 0 {
		return nil, errors.New("cannot process and backtrack LCS with empty string")
	} else if utils.Equal(runeStr1, runeStr2) {
		return []string{str1}, nil
	}

	lcsMatrix := c.lcsProcess(runeStr1, runeStr2)
	return c.processBacktrackAll(str1, str2, lcsMatrix, len(runeStr1), len(runeStr2)).ToArray(), nil
}

func (c *Calculator) processBacktrackAll(str1, str2 string, lcsMatrix [][]int, m, n int) utils.StringHashMap {
	runeStr1 := []rune(str1)
	runeStr2 := []rune(str2)

	substrings := make(utils.StringHashMap)

	if m == 0 || n == 0 {
		substrings[""] = struct{}{}
	} else if runeStr1[m-1] == runeStr2[n-1] {
		for key := range c.processBacktrackAll(str1, str2, lcsMatrix, m-1, n-1) {
			substrings[key+string(runeStr1[m-1])] = struct{}{}
		}
	} else {
		if lcsMatrix[m-1][n] >= lcsMatrix[m][n-1] {
			substrings.AddAll(c.processBacktrackAll(str1, str2, lcsMatrix, m-1, n))
		}
		if lcsMatrix[m][n-1] >= lcsMatrix[m-1][n] {
			substrings.AddAll(c.processBacktrackAll(str1, str2, lcsMatrix, m, n-1))
		}
	}

	return substrings
}

// Diff backtraces through the LCS matrix and returns the diff between two sequences
func (c *Calculator) Diff(str1, str2 string) ([]string, error) {
	if c.caseInsensitive {
		str1 = strings.ToLower(str1)
		str2 = strings.ToLower(str2)
	}

	runeStr1 := []rune(str1)
	runeStr2 := []rune(str2)

	if len(runeStr1) == 0 || len(runeStr2) == 0 {
		return nil, errors.New("cannot process LCS diff with empty string")
	} else if utils.Equal(runeStr1, runeStr2) {
		return []string{str1}, nil
	}

	lcsMatrix := c.lcsProcess(runeStr1, runeStr2)
	return c.processDiff(str1, str2, lcsMatrix, len(runeStr1), len(runeStr2)), nil
}

func (c *Calculator) processDiff(str1, str2 string, lcsMatrix [][]int, m, n int) []string {
	runeStr1 := []rune(str1)
	runeStr2 := []rune(str2)

	diff := make([]string, 2)

	if m > 0 && n > 0 && runeStr1[m-1] == runeStr2[n-1] {
		diff = c.processDiff(str1, str2, lcsMatrix, m-1, n-1)
		diff[0] = diff[0] + " " + string(runeStr1[m-1])
		diff[1] = diff[1] + "  "
		return diff
	} else if n > 0 && (m == 0 || lcsMatrix[m][n-1] > lcsMatrix[m-1][n]) {
		diff = c.processDiff(str1, str2, lcsMatrix, m, n-1)
		diff[0] = diff[0] + " " + string(runeStr2[n-1])
		diff[1] = diff[1] + " +"
		return diff
	} else if m > 0 && (n == 0 || lcsMatrix[m][n-1] <= lcsMatrix[m-1][n]) {
		diff = c.processDiff(str1, str2, lcsMatrix, m-1, n)
		diff[0] = diff[0] + " " + string(runeStr1[m-1])
		diff[1] = diff[1] + " -"
		return diff
	}

	return diff
}

// EditDistance determines the edit distance using LCS
// (allows only insert and delete operations)
func (c *Calculator) EditDistance(str1, str2 string) int {
	if c.caseInsensitive {
		str1 = strings.ToLower(str1)
		str2 = strings.ToLower(str2)
	}

	runeStr1 := []rune(str1)
	runeStr2 := []rune(str2)

	if len(runeStr1) == 0 {
		return len(runeStr2)
	} else if len(runeStr2) == 0 {
		return len(runeStr1)
	} else if str1 == str2 {
		return 0
	}

	lcs := c.calculateLCS(str1, str2)
	return (len(runeStr1) - lcs) + (len(runeStr2) - lcs)
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

// SetMemoryPool sets the memory pool for matrix allocation
func (c *Calculator) SetMemoryPool(pool *core.MemoryPool) error {
	c.pool = pool
	return nil
}
