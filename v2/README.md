# go-edlib v2.0

> Modern, high-performance string comparison and edit distance library for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/hbollon/go-edlib/v2.svg)](https://pkg.go.dev/github.com/hbollon/go-edlib/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/hbollon/go-edlib/v2)](https://goreportcard.com/report/github.com/hbollon/go-edlib/v2)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/hbollon/go-edlib/blob/master/LICENSE.md)

## ✨ What's New in v2.0

go-edlib v2 is a **complete redesign** with modern Go idioms, performance optimizations, and enterprise features:

- 🔥 **Fluent API** - Intuitive, type-safe builder pattern
- ⚡ **Performance** - Memory pooling, parallel processing, early termination
- 🎯 **Zero Dependencies** - Pure Go implementation (except optional OpenTelemetry)
- 🌍 **Unicode Support** - Full support for emojis, CJK, RTL text
- 📊 **Observability** - Built-in metrics and execution profiling
- 🔍 **Advanced Search** - Fuzzy search with ranking and filtering
- ⚙️ **Extensible** - Plugin architecture via registry pattern

## 🚀 Quick Start

### Installation

```bash
go get github.com/hbollon/go-edlib/v2
```

### Simple Usage

```go
import "github.com/hbollon/go-edlib/v2"

// Quick one-liners
distance, _ := edlib.Quick().Levenshtein("hello", "helo")
similarity, _ := edlib.Quick().JaroWinkler("martha", "marhta")

// Fluent API
result, _ := edlib.New().
    Using(algorithms.Levenshtein).
    WithThreshold(0.8).
    Compare("hello", "helo")

fmt.Printf("Distance: %d, Similarity: %.2f\n", result.Distance, result.Similarity)
```

## 📖 Features

### Supported Algorithms

#### Edit Distance
- **Levenshtein** - Classic edit distance (insertions, deletions, substitutions)
- **Damerau-Levenshtein** - Includes transpositions
- **Hamming** - Fast distance for equal-length strings

#### Similarity Metrics
- **Jaro** - Optimal for short strings
- **Jaro-Winkler** - Jaro with prefix bonus
- **Cosine** - Vector-based similarity
- **Jaccard** - Set-based similarity
- **Sorensen-Dice** - Token overlap similarity

#### Sequence Algorithms
- **LCS** - Longest Common Subsequence

#### Phonetic (Coming Soon)
- **Soundex** - Phonetic encoding
- **Metaphone** - Advanced phonetic algorithm
- **Double Metaphone** - Improved Metaphone

## 💡 Usage Examples

### 1. Basic Comparison

```go
// Simple distance calculation
dist, _ := edlib.Quick().Levenshtein("kitten", "sitting")
// dist = 3

// Full comparison with metadata
result, _ := edlib.New().
    Using(algorithms.Levenshtein).
    Compare("hello", "helo")

fmt.Printf("Distance: %d\n", result.Distance)       // 1
fmt.Printf("Similarity: %.2f\n", result.Similarity) // 0.80
fmt.Printf("Algorithm: %s\n", result.Algorithm)     // Levenshtein
```

### 2. Advanced Configuration

```go
result, _ := edlib.New().
    Using(algorithms.Levenshtein).
    WithThreshold(0.7).
    WithCaseSensitive(false).
    WithMemoryOptimization(true).
    WithEarlyTermination(true).
    WithWeights(edlib.CostWeights{
        Insertion:    1,
        Deletion:     1,
        Substitution: 2,
    }).
    Compare("Hello World", "hello wrold")
```

### 3. Fuzzy Search

```go
targets := []string{"apple", "application", "apply", "banana", "bandana"}

// Find best match
best, _ := edlib.New().
    Search().
    WithAlgorithm(algorithms.JaroWinkler).
    FindBest("aplication", targets)
// best.Value = "application"

// Find top 3 matches
top3, _ := edlib.New().
    Search().
    WithAlgorithm(algorithms.Levenshtein).
    FindTopN("aplication", targets, 3)

// Find all matches above threshold
matches, _ := edlib.New().
    Search().
    WithThreshold(0.7).
    FindAll("appl", targets)
```

### 4. Batch Processing

```go
pairs := []edlib.StringPair{
    {First: "hello", Second: "helo"},
    {First: "world", Second: "wrold"},
    {First: "test", Second: "testing"},
}

// Sequential processing
results, _ := edlib.New().
    Using(algorithms.Levenshtein).
    Batch(pairs)

// Parallel processing
results, _ := edlib.New().
    Using(algorithms.Levenshtein).
    WithParallel(4).
    Batch(pairs)
```

### 5. Context and Timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := edlib.New().
    WithContext(ctx).
    Using(algorithms.Levenshtein).
    Compare(longString1, longString2)

if err == context.DeadlineExceeded {
    fmt.Println("Operation timed out")
}
```

### 6. Performance Metrics

```go
result, _ := edlib.New().
    Using(algorithms.Levenshtein).
    WithMetrics(true).
    WithMemoryOptimization(true).
    Compare("hello", "world")

if result.Metrics != nil {
    fmt.Printf("Execution time: %v\n", result.Metrics.ExecutionTime)
    fmt.Printf("Memory allocated: %d bytes\n", result.Metrics.MemoryAllocated)
    fmt.Printf("Pool hit rate: %d/%d\n", result.Metrics.PoolHits, result.Metrics.PoolMisses)
}
```

### 7. Unicode Support

```go
// Emojis
dist, _ := edlib.Quick().Levenshtein("hello 👋", "hello 👋")

// CJK characters
sim, _ := edlib.Quick().JaroWinkler("こんにちは", "こんにちわ")

// Mixed scripts
result, _ := edlib.New().
    Using(algorithms.Levenshtein).
    Compare("Hello 世界", "Hello 世界!")
```

## 🏗️ Architecture

### Design Principles

1. **Fluent API** - Chainable methods for intuitive configuration
2. **Registry Pattern** - Extensible algorithm registration
3. **Factory Pattern** - Lazy calculator instantiation
4. **Builder Pattern** - Complex object construction
5. **Strategy Pattern** - Swappable algorithms

### Directory Structure

```
v2/
├── edlib.go              # Main entry point
├── builder.go            # Fluent API builders
├── config.go             # Configuration types
├── types.go              # Public types
├── errors.go             # Error types
├── algorithms/           # Algorithm types and properties
├── search/               # Search engine
├── internal/
│   ├── core/            # Core interfaces and registry
│   ├── impl/            # Algorithm implementations
│   │   ├── edit_distance/
│   │   ├── similarity/
│   │   └── sequence/
│   └── utils/           # Internal utilities
└── examples/            # Usage examples
```

## 📊 Performance

### Benchmarks

```
BenchmarkLevenshtein/short-8                5000000    250 ns/op     0 B/op    0 allocs/op
BenchmarkLevenshtein/short_with_pool-8      8000000    180 ns/op     0 B/op    0 allocs/op
BenchmarkLevenshtein/medium-8                500000   2800 ns/op     0 B/op    0 allocs/op
BenchmarkJaroWinkler/short-8                3000000    400 ns/op   128 B/op    2 allocs/op
BenchmarkCosine/short-8                     2000000    650 ns/op   256 B/op    5 allocs/op
BenchmarkSearch/find_best-8                  100000  12000 ns/op  2048 B/op   50 allocs/op
```

### Optimizations

- **Memory Pooling**: Reuses buffers to reduce GC pressure
- **Early Termination**: Stops computation when threshold exceeded
- **Space Optimization**: O(min(n,m)) space for Levenshtein
- **Parallel Processing**: Worker pools for batch operations
- **Context-aware**: Respects cancellation and timeouts

## 🔧 Migration from v1

### Breaking Changes

v2 is a complete rewrite with **no backwards compatibility** with v1. Key changes:

| v1 | v2 |
|----|-----|
| `edlib.LevenshteinDistance(s1, s2)` | `edlib.Quick().Levenshtein(s1, s2)` |
| `edlib.StringsSimilarity(s1, s2, algo)` | `edlib.New().Using(algo).Similarity(s1, s2)` |
| `edlib.FuzzySearch(q, targets, algo)` | `edlib.New().Search().WithAlgorithm(algo).FindBest(q, targets)` |

### Migration Guide

```go
// v1
dist := edlib.LevenshteinDistance("hello", "helo")
sim, _ := edlib.StringsSimilarity("hello", "helo", edlib.Levenshtein)
match, _ := edlib.FuzzySearch("test", targets, edlib.Levenshtein)

// v2
dist, _ := edlib.Quick().Levenshtein("hello", "helo")
sim, _ := edlib.New().Using(algorithms.Levenshtein).Similarity("hello", "helo")
match, _ := edlib.New().Search().WithAlgorithm(algorithms.Levenshtein).FindBest("test", targets)
```

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details.

### Development

```bash
# Run tests
go test ./...

# Run benchmarks
go test -bench=. -benchmem ./...

# Run with coverage
go test -cover ./...
```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## 👤 Author

**Hugo Bollon**

- GitHub: [@hbollon](https://github.com/hbollon)
- LinkedIn: [Hugo Bollon](https://www.linkedin.com/in/hugo-bollon-68a2381a4/)
- Website: [hugobollon.me](https://www.hugobollon.me)

## ⭐ Show your support

Give a ⭐️ if this project helped you!

## 📚 References

- [Levenshtein Distance](https://en.wikipedia.org/wiki/Levenshtein_distance)
- [Jaro-Winkler Similarity](https://en.wikipedia.org/wiki/Jaro%E2%80%93Winkler_distance)
- [Cosine Similarity](https://en.wikipedia.org/wiki/Cosine_similarity)
