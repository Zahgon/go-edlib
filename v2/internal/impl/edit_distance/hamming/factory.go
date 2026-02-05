package hamming

import (
	"github.com/hbollon/go-edlib/v2/algorithms"
	"github.com/hbollon/go-edlib/v2/internal/core"
)

// Factory creates Hamming distance calculators
type Factory struct{}

// CreateCalculator creates a new Hamming calculator
func (f *Factory) CreateCalculator(opts ...core.CalculatorOption) (core.Calculator, error) {
	return New(opts...)
}

// Algorithm returns the algorithm type
func (f *Factory) Algorithm() algorithms.Algorithm {
	return algorithms.Hamming
}

// Properties returns the algorithm properties
func (f *Factory) Properties() algorithms.AlgorithmProperties {
	return algorithms.Hamming.Properties()
}

// Ensure Factory implements AlgorithmFactory
var _ core.AlgorithmFactory = (*Factory)(nil)

// init registers the Hamming algorithm in the global registry
func init() {
	if err := core.Register(&Factory{}); err != nil {
		panic(err)
	}
}
