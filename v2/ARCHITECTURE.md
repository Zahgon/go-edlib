# go-edlib v2.0 Architecture

## 📋 Overview

This document describes the architecture and design decisions for go-edlib v2.0, a complete rewrite of the library with modern Go patterns and enterprise features.

## 🏗️ Core Design Principles

### 1. **Fluent API Pattern**
- Intuitive chainable methods for configuration
- Type-safe builders with compile-time validation
- Sensible defaults with easy customization

### 2. **Registry Pattern**
- Extensible algorithm registration via `init()`
- Factory pattern for lazy instantiation
- Plugin-like architecture for third-party algorithms

### 3. **Interface Segregation**
- Small, focused interfaces (Calculator, OptimizedCalculator, MetricsCalculator)
- Composable capabilities via interface composition
- No forced dependencies on unused features

### 4. **Zero Dependencies**
- Pure Go implementation for all core algorithms
- Optional OpenTelemetry for metrics (not yet implemented)
- No external dependencies for basic usage

### 5. **Performance First**
- Memory pooling via `sync.Pool`
- Two-row algorithm for space optimization O(min(n,m))
- Early termination support
- Context-aware cancellation
- Parallel processing for batch operations

## 📁 Directory Structure

```
v2/
├── edlib.go                 # Main package entry point
├── builder.go               # Core fluent API builder
├── builder_specialized.go   # Specialized builders (EditDistance, Similarity, Quick)
├── builder_search.go        # Search builder
├── config.go                # Configuration types and validation
├── types.go                 # Public types (ComparisonResult, StringPair, etc.)
├── errors.go                # Error types and constructors
│
├── algorithms/              # Algorithm type definitions
│   └── algorithm.go         # Algorithm enum, properties, metadata
│
├── search/                  # Search engine package
│   ├── search.go           # Search engine implementation
│   └── types.go            # Search-specific types (Result, etc.)
│
├── internal/
│   ├── core/               # Core interfaces and infrastructure
│   │   ├── interfaces.go   # Calculator interfaces
│   │   ├── registry.go     # Algorithm registry
│   │   ├── pool.go         # Memory pool
│   │   ├── metrics.go      # Metrics collection
│   │   └── errors.go       # Internal errors
│   │
│   ├── impl/               # Algorithm implementations
│   │   ├── edit_distance/
│   │   │   └── levenshtein/
│   │   │       ├── levenshtein.go  # Implementation
│   │   │       └── factory.go      # Factory + auto-registration
│   │   │
│   │   └── similarity/
│   │       ├── jaro/
│   │       │   ├── jaro.go        # Jaro + Jaro-Winkler
│   │       │   └── factory.go     # Factory + auto-registration
│   │       └── cosine/
│   │           ├── cosine.go      # Implementation
│   │           └── factory.go     # Factory + auto-registration
│   │
│   ├── search/
│   │   └── parallel.go     # Parallel processing utilities
│   │
│   └── utils/
│       └── utils.go        # Utility functions
│
└── examples/
    └── basic/
        └── main.go         # Usage examples
```

## 🔌 Core Interfaces

### Calculator Interfaces

```go
// Base calculator interface
type Calculator interface {
    Distance(s1, s2 string) int
    DistanceWithContext(ctx context.Context, s1, s2 string) (int, error)
    Similarity(s1, s2 string) float64
    SimilarityWithContext(ctx context.Context, s1, s2 string) (float64, error)
    MaxDistance(s1, s2 string) int
    Name() string
    Properties() algorithms.AlgorithmProperties
}

// Optimization capabilities
type OptimizedCalculator interface {
    Calculator
    WithThreshold(threshold float64) OptimizedCalculator
    WithMemoryPool(pool *MemoryPool) OptimizedCalculator
    WithEarlyTermination(enabled bool) OptimizedCalculator
    WithMaxDistance(maxDist int) OptimizedCalculator
}

// Metrics capabilities (not yet implemented)
type MetricsCalculator interface {
    Calculator
    DistanceWithMetrics(s1, s2 string) (int, *Metrics)
    SimilarityWithMetrics(s1, s2 string) (float64, *Metrics)
}
```

### Factory Interface

```go
type AlgorithmFactory interface {
    CreateCalculator(opts ...CalculatorOption) (Calculator, error)
    Algorithm() algorithms.Algorithm
    Properties() algorithms.AlgorithmProperties
}
```

## 🎯 API Design

### Quick API (One-liners)
```go
dist, _ := edlib.Quick().Levenshtein("hello", "helo")
sim, _ := edlib.Quick().JaroWinkler("martha", "marhta")
```

### Fluent API (Configuration)
```go
result, _ := edlib.New().
    Using(algorithms.Levenshtein).
    WithThreshold(0.8).
    WithCaseSensitive(false).
    WithMemoryOptimization(true).
    Compare("Hello", "hello")
```

### Search API
```go
results, _ := edlib.New().
    Search().
    WithAlgorithm(algorithms.JaroWinkler).
    WithThreshold(0.7).
    FindAll("test", targets)
```

## 🚀 Performance Optimizations

### 1. Memory Pooling
- `sync.Pool` for matrix and buffer reuse
- Reduces GC pressure
- Configurable pool sizes
- Automatic pool management

### 2. Space Optimization
- Two-row Levenshtein: O(min(n,m)) instead of O(n*m)
- Ensures shorter string is used for column dimension
- Minimal allocations

### 3. Early Termination
- Stop computation when threshold exceeded
- Distance-based early exit
- Row-by-row minimum tracking

### 4. Parallel Processing
- Worker pools for batch operations
- Optimal chunk size calculation
- Context-aware cancellation
- Load balancing

## 📊 Implemented Algorithms

### Edit Distance
- ✅ **Levenshtein** - Full implementation with optimizations
  - Two-row space optimization
  - Custom weights support
  - Early termination
  - Memory pooling

### Similarity
- ✅ **Jaro** - Complete implementation
- ✅ **Jaro-Winkler** - Complete with prefix bonus
- ✅ **Cosine** - N-gram based similarity

### Search
- ✅ **FindBest** - Single best match
- ✅ **FindAll** - All matches above threshold
- ✅ **FindTopN** - Top N ranked results
- ✅ **Ranking** - Score, Distance, Hybrid methods

## 🔧 Configuration System

### Config Structure
```go
type Config struct {
    Algorithm        algorithms.Algorithm
    Context          context.Context
    Timeout          time.Duration
    Threshold        float64
    CaseSensitive    bool
    Weights          CostWeights
    Parallel         bool
    Workers          int
    EnableMetrics    bool
    UseMemoryPool    bool
    EarlyTermination bool
    MaxDistance      int
}
```

### Validation
- Automatic validation on execution
- Clear error messages
- Type-safe constraints

## 🧪 Testing

### Test Coverage
- ✅ Quick API tests
- ✅ Fluent API tests
- ✅ Algorithm correctness tests
- ✅ Search functionality tests
- ✅ Unicode support tests
- ✅ Context cancellation tests
- ✅ Batch processing tests
- ⚠️ Metrics tests (implementation pending)

### Benchmarks
- ✅ Algorithm benchmarks (short, medium, long strings)
- ✅ Optimization comparisons (with/without pool, early termination)
- ✅ Search benchmarks
- ✅ Batch processing benchmarks
- ✅ Quick API benchmarks

## 📝 Known Limitations

### Not Yet Implemented
1. **Case-insensitive comparison** - Configuration exists but not applied in calculators
2. **Metrics collection** - Metrics structure exists but not populated
3. **Hamming distance** - Not implemented
4. **Damerau-Levenshtein variants** - Not implemented
5. **LCS** - Not implemented
6. **Phonetic algorithms** - Not implemented
7. **Advanced parallel processing** - Basic structure exists but not fully integrated
8. **Cache system** - Not implemented

### Intentional Design Decisions
- **No backwards compatibility** with v1 - Complete redesign
- **Minimal dependencies** - Only OpenTelemetry for metrics (optional)
- **Registry auto-registration** - Uses `init()` which some consider anti-pattern
- **Builder pattern complexity** - More verbose than simple functions but more flexible

## 🔮 Future Enhancements (Phase 4-6)

### Phase 4: Remaining Algorithms
- Hamming distance
- Damerau-Levenshtein (OSA and true)
- LCS and diff
- Jaccard
- Sorensen-Dice
- Q-Gram

### Phase 5: Advanced Performance
- SIMD optimizations
- Result caching (LRU)
- Advanced parallelization strategies
- Memory-mapped large string handling

### Phase 6: Enterprise Features
- OpenTelemetry integration
- Prometheus metrics
- Structured logging
- Index structures (tries, suffix trees)
- ML-based ranking
- Hybrid algorithm strategies

## 📚 References

### Design Patterns
- Fluent Interface Pattern
- Builder Pattern
- Factory Pattern
- Registry Pattern
- Strategy Pattern
- Template Method Pattern

### Algorithms
- Wagner-Fischer algorithm (Levenshtein)
- Jaro-Winkler distance
- Cosine similarity with n-grams
- Space-optimized dynamic programming

### Go Best Practices
- Effective Go
- Go Code Review Comments
- Go Proverbs
- Standard library patterns

## 👥 Contributors

- Hugo Bollon (@hbollon) - Original author and v2 architect

## 📄 License

MIT License - See LICENSE.md
