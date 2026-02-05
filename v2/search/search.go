package search

import (
	"context"
	"sort"

	"github.com/hbollon/go-edlib/v2/algorithms"
	"github.com/hbollon/go-edlib/v2/internal/core"
)

// RankingMethod defines how search results are ranked
type RankingMethod int

const (
	// ScoreRanking ranks by similarity score (highest first)
	ScoreRanking RankingMethod = iota

	// DistanceRanking ranks by edit distance (lowest first)
	DistanceRanking

	// HybridRanking uses a combination of score and distance
	HybridRanking
)

// Config contains configuration for search operations
type Config struct {
	Algorithm   algorithms.Algorithm
	Threshold   float64
	MaxResults  int
	Ranking     RankingMethod
	Parallel    bool
	Workers     int
	Context     context.Context
}

// DefaultConfig returns default search configuration
func DefaultConfig() *Config {
	return &Config{
		Algorithm:  algorithms.Levenshtein,
		Threshold:  0.0,
		MaxResults: 0, // No limit
		Ranking:    ScoreRanking,
		Parallel:   false,
		Workers:    0,
		Context:    context.Background(),
	}
}

// Engine performs fuzzy string searching
type Engine struct {
	config *Config
}

// New creates a new search engine
func New(config *Config) *Engine {
	if config == nil {
		config = DefaultConfig()
	}
	return &Engine{config: config}
}

// FindBest finds the single best match for the query in the target list
func (e *Engine) FindBest(query string, targets []string) (*Result, error) {
	if len(targets) == 0 {
		return nil, ErrNoResults
	}

	results, err := e.FindAll(query, targets)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, ErrNoResults
	}

	return &results[0], nil
}

// FindAll finds all matches for the query in the target list
func (e *Engine) FindAll(query string, targets []string) ([]Result, error) {
	if len(targets) == 0 {
		return nil, ErrNoResults
	}

	// Create calculator
	calc, err := core.CreateCalculator(e.config.Algorithm)
	if err != nil {
		return nil, err
	}

	// Compute similarities
	results := make([]Result, 0, len(targets))

	for i, target := range targets {
		// Check context cancellation
		select {
		case <-e.config.Context.Done():
			return nil, e.config.Context.Err()
		default:
		}

		var distance int
		var similarity float64

		if e.config.Context != nil {
			distance, err = calc.DistanceWithContext(e.config.Context, query, target)
			if err != nil {
				continue
			}
			similarity, err = calc.SimilarityWithContext(e.config.Context, query, target)
			if err != nil {
				continue
			}
		} else {
			distance = calc.Distance(query, target)
			similarity = calc.Similarity(query, target)
		}

		// Apply threshold filter
		if e.config.Threshold > 0 && similarity < e.config.Threshold {
			continue
		}

		result := Result{
			Value:    target,
			Score:    similarity,
			Distance: distance,
			Index:    i,
			Metadata: ResultMetadata{
				Algorithm:      e.config.Algorithm,
				AboveThreshold: similarity >= e.config.Threshold,
			},
		}

		results = append(results, result)
	}

	if len(results) == 0 {
		return nil, ErrNoResults
	}

	// Sort results based on ranking method
	e.sortResults(results)

	// Limit results if maxResults is set
	if e.config.MaxResults > 0 && len(results) > e.config.MaxResults {
		results = results[:e.config.MaxResults]
	}

	// Set ranks
	for i := range results {
		results[i].Metadata.Rank = i + 1
	}

	return results, nil
}

// sortResults sorts the results based on the ranking method
func (e *Engine) sortResults(results []Result) {
	switch e.config.Ranking {
	case ScoreRanking:
		sort.Slice(results, func(i, j int) bool {
			if results[i].Score == results[j].Score {
				return results[i].Distance < results[j].Distance
			}
			return results[i].Score > results[j].Score
		})

	case DistanceRanking:
		sort.Slice(results, func(i, j int) bool {
			if results[i].Distance == results[j].Distance {
				return results[i].Score > results[j].Score
			}
			return results[i].Distance < results[j].Distance
		})

	case HybridRanking:
		// Hybrid: combine normalized score and distance
		sort.Slice(results, func(i, j int) bool {
			scoreI := results[i].Score
			scoreJ := results[j].Score
			
			// Higher score is better, lower distance is better
			// Normalize and combine (60% weight on score, 40% on distance)
			hybridI := scoreI*0.6 + (1.0-float64(results[i].Distance)/100.0)*0.4
			hybridJ := scoreJ*0.6 + (1.0-float64(results[j].Distance)/100.0)*0.4
			
			return hybridI > hybridJ
		})
	}
}

// FindWithThreshold is a convenience method for threshold-based search
func (e *Engine) FindWithThreshold(query string, targets []string, threshold float64) ([]Result, error) {
	e.config.Threshold = threshold
	return e.FindAll(query, targets)
}

// FindTopN finds the top N matches for the query
func (e *Engine) FindTopN(query string, targets []string, n int) ([]Result, error) {
	e.config.MaxResults = n
	return e.FindAll(query, targets)
}

// ErrNoResults is returned when no results match the search criteria
var ErrNoResults = core.ErrNoResults
