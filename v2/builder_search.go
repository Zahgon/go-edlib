package edlib

import (
	"github.com/hbollon/go-edlib/v2/algorithms"
	"github.com/hbollon/go-edlib/v2/search"
)

// SearchBuilder provides a fluent interface for fuzzy search operations
type SearchBuilder struct {
	builder       *Builder
	maxResults    int
	rankingMethod search.RankingMethod
}

// WithAlgorithm sets the algorithm to use for search
func (s *SearchBuilder) WithAlgorithm(algo algorithms.Algorithm) *SearchBuilder {
	s.builder.config.Algorithm = algo
	return s
}

// WithThreshold sets the minimum similarity threshold
func (s *SearchBuilder) WithThreshold(threshold float64) *SearchBuilder {
	s.builder.config.Threshold = threshold
	return s
}

// WithMaxResults limits the number of results returned
func (s *SearchBuilder) WithMaxResults(max int) *SearchBuilder {
	s.maxResults = max
	return s
}

// WithRanking sets the ranking method for results
func (s *SearchBuilder) WithRanking(method search.RankingMethod) *SearchBuilder {
	s.rankingMethod = method
	return s
}

// WithParallel enables parallel processing with the specified number of workers
func (s *SearchBuilder) WithParallel(workers int) *SearchBuilder {
	s.builder.config.Parallel = true
	s.builder.config.Workers = workers
	return s
}

// WithCaseSensitive sets whether search is case-sensitive
func (s *SearchBuilder) WithCaseSensitive(sensitive bool) *SearchBuilder {
	s.builder.config.CaseSensitive = sensitive
	return s
}

// FindBest finds the single best match for the query
func (s *SearchBuilder) FindBest(query string, targets []string) (*search.Result, error) {
	engine := s.createEngine()
	return engine.FindBest(query, targets)
}

// FindAll finds all matches above the threshold
func (s *SearchBuilder) FindAll(query string, targets []string) ([]search.Result, error) {
	engine := s.createEngine()
	return engine.FindAll(query, targets)
}

// FindTopN finds the top N matches
func (s *SearchBuilder) FindTopN(query string, targets []string, n int) ([]search.Result, error) {
	s.maxResults = n
	engine := s.createEngine()
	return engine.FindTopN(query, targets, n)
}

// FindWithThreshold finds all matches with at least the specified similarity
func (s *SearchBuilder) FindWithThreshold(query string, targets []string, threshold float64) ([]search.Result, error) {
	s.builder.config.Threshold = threshold
	engine := s.createEngine()
	return engine.FindWithThreshold(query, targets, threshold)
}

// createEngine creates a search engine from the builder configuration
func (s *SearchBuilder) createEngine() *search.Engine {
	config := &search.Config{
		Algorithm:  s.builder.config.Algorithm,
		Threshold:  s.builder.config.Threshold,
		MaxResults: s.maxResults,
		Ranking:    s.rankingMethod,
		Parallel:   s.builder.config.Parallel,
		Workers:    s.builder.config.Workers,
		Context:    s.builder.config.Context,
	}

	return search.New(config)
}
