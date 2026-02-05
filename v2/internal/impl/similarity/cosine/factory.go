package cosine

import (
	"github.com/hbollon/go-edlib/v2/algorithms"
	"github.com/hbollon/go-edlib/v2/internal/core"
)

// Factory creates Cosine similarity calculator instances
type Factory struct{}

// CreateCalculator creates a new Cosine calculator
func (f *Factory) CreateCalculator(opts ...core.CalculatorOption) (core.Calculator, error) {
	return New(opts...)
}

// Algorithm returns the algorithm type
func (f *Factory) Algorithm() algorithms.Algorithm {
	return algorithms.Cosine
}

// Properties returns the algorithm properties
func (f *Factory) Properties() algorithms.AlgorithmProperties {
	return algorithms.Cosine.Properties()
}

// Ensure Factory implements AlgorithmFactory
var _ core.AlgorithmFactory = (*Factory)(nil)

// Register registers the Cosine algorithm
func init() {
	if err := core.Register(&Factory{}); err != nil {
		panic(err)
	}
}
