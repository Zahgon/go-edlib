package edlib

import "github.com/hbollon/go-edlib/v2/algorithms"

// EditDistanceBuilder provides a fluent interface for edit distance algorithms
type EditDistanceBuilder struct {
	builder *Builder
}

// Levenshtein configures the builder to use Levenshtein distance
func (b *EditDistanceBuilder) Levenshtein() *Builder {
	b.builder.config.Algorithm = algorithms.Levenshtein
	return b.builder
}

// SimilarityBuilder provides a fluent interface for similarity algorithms
type SimilarityBuilder struct {
	builder *Builder
}

// Jaro configures the builder to use Jaro similarity
func (b *SimilarityBuilder) Jaro() *Builder {
	b.builder.config.Algorithm = algorithms.Jaro
	return b.builder
}

// JaroWinkler configures the builder to use Jaro-Winkler similarity
func (b *SimilarityBuilder) JaroWinkler() *Builder {
	b.builder.config.Algorithm = algorithms.JaroWinkler
	return b.builder
}

// Cosine configures the builder to use Cosine similarity
func (b *SimilarityBuilder) Cosine() *Builder {
	b.builder.config.Algorithm = algorithms.Cosine
	return b.builder
}

// QuickBuilder provides a simple interface for one-line operations
type QuickBuilder struct {
	config *Config
}

// Levenshtein computes Levenshtein distance
func (q *QuickBuilder) Levenshtein(s1, s2 string) (int, error) {
	b := New().Using(algorithms.Levenshtein)
	return b.Distance(s1, s2)
}

// Jaro computes Jaro similarity
func (q *QuickBuilder) Jaro(s1, s2 string) (float64, error) {
	b := New().Using(algorithms.Jaro)
	return b.Similarity(s1, s2)
}

// JaroWinkler computes Jaro-Winkler similarity
func (q *QuickBuilder) JaroWinkler(s1, s2 string) (float64, error) {
	b := New().Using(algorithms.JaroWinkler)
	return b.Similarity(s1, s2)
}

// Cosine computes Cosine similarity
func (q *QuickBuilder) Cosine(s1, s2 string) (float64, error) {
	b := New().Using(algorithms.Cosine)
	return b.Similarity(s1, s2)
}
