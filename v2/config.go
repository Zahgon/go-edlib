package edlib

import (
	"context"
	"time"

	"github.com/hbollon/go-edlib/v2/algorithms"
	"github.com/hbollon/go-edlib/v2/internal/core"
)

// Config contains all configuration options for string comparison and search operations
type Config struct {
	// Algorithm specifies which algorithm to use
	Algorithm algorithms.Algorithm

	// Context for cancellation and timeout
	Context context.Context

	// Timeout for operations (creates context if Context is nil)
	Timeout time.Duration

	// Threshold for similarity filtering [0.0 to 1.0]
	Threshold float64

	// CaseSensitive determines if comparisons are case-sensitive
	CaseSensitive bool

	// Normalize determines if strings should be normalized (Unicode normalization)
	Normalize bool

	// Weights defines the costs for edit operations
	Weights CostWeights

	// Parallel enables parallel processing
	Parallel bool

	// Workers specifies the number of parallel workers (0 = auto-detect)
	Workers int

	// EnableMetrics enables detailed execution metrics
	EnableMetrics bool

	// EnableCache enables result caching
	EnableCache bool

	// CacheSize specifies the maximum cache size (0 = default: 1000)
	CacheSize int

	// UseMemoryPool enables memory pooling for matrix allocations
	UseMemoryPool bool

	// EarlyTermination enables early termination optimization
	EarlyTermination bool

	// MaxDistance allows early termination when distance exceeds this value
	MaxDistance int
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Algorithm:        algorithms.Levenshtein,
		Context:          context.Background(),
		Timeout:          0, // No timeout
		Threshold:        0.0,
		CaseSensitive:    true,
		Normalize:        false,
		Weights:          DefaultCostWeights(),
		Parallel:         false,
		Workers:          0, // Auto-detect
		EnableMetrics:    false,
		EnableCache:      false,
		CacheSize:        1000,
		UseMemoryPool:    true,
		EarlyTermination: false,
		MaxDistance:      0, // No limit
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Threshold < 0.0 || c.Threshold > 1.0 {
		return NewConfigurationError("Threshold", c.Threshold, "must be between 0.0 and 1.0")
	}

	if c.Workers < 0 {
		return NewConfigurationError("Workers", c.Workers, "must be non-negative")
	}

	if c.CacheSize < 0 {
		return NewConfigurationError("CacheSize", c.CacheSize, "must be non-negative")
	}

	if c.MaxDistance < 0 {
		return NewConfigurationError("MaxDistance", c.MaxDistance, "must be non-negative")
	}

	if err := c.Weights.Validate(); err != nil {
		return err
	}

	return nil
}

// Clone creates a deep copy of the configuration
func (c *Config) Clone() *Config {
	clone := *c

	// Deep copy weights if custom weights exist
	if c.Weights.Custom != nil {
		clone.Weights.Custom = make(map[rune]map[rune]int)
		for k, v := range c.Weights.Custom {
			clone.Weights.Custom[k] = make(map[rune]int)
			for kk, vv := range v {
				clone.Weights.Custom[k][kk] = vv
			}
		}
	}

	return &clone
}

// WithTimeout creates a context with timeout if needed
func (c *Config) WithTimeout() (context.Context, context.CancelFunc) {
	if c.Context == nil {
		c.Context = context.Background()
	}

	if c.Timeout > 0 {
		return context.WithTimeout(c.Context, c.Timeout)
	}

	return c.Context, func() {}
}

// ToCalculatorOptions converts config to calculator options
func (c *Config) ToCalculatorOptions() []core.CalculatorOption {
	var opts []core.CalculatorOption

	if c.Threshold > 0 {
		opts = append(opts, core.WithThreshold(c.Threshold))
	}

	if !c.CaseSensitive {
		opts = append(opts, core.WithCaseInsensitive())
	}

	if c.Weights.Insertion != 1 || c.Weights.Deletion != 1 || c.Weights.Substitution != 1 || c.Weights.Transposition != 1 || c.Weights.Custom != nil {
		opts = append(opts, core.WithWeights(c.Weights))
	}

	if c.UseMemoryPool {
		opts = append(opts, core.WithMemoryPool())
	}

	if c.EarlyTermination {
		opts = append(opts, core.WithEarlyTermination())
	}

	if c.MaxDistance > 0 {
		opts = append(opts, core.WithMaxDistance(c.MaxDistance))
	}

	if c.EnableMetrics {
		opts = append(opts, core.WithMetrics())
	}

	return opts
}
