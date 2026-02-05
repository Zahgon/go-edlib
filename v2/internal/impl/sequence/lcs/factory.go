package lcs

import (
	"github.com/hbollon/go-edlib/v2/algorithms"
	"github.com/hbollon/go-edlib/v2/internal/core"
)

// Factory creates LCS calculators
type Factory struct{}

// CreateCalculator creates a new LCS calculator
func (f *Factory) CreateCalculator(opts ...core.CalculatorOption) (core.Calculator, error) {
	return New(opts...)
}

// Algorithm returns the algorithm identifier
func (f *Factory) Algorithm() algorithms.Algorithm {
	return algorithms.LCS
}

// Properties returns the algorithm properties
func (f *Factory) Properties() algorithms.AlgorithmProperties {
	return algorithms.LCS.Properties()
}

var _ core.AlgorithmFactory = (*Factory)(nil)

func init() {
	if err := core.Register(&Factory{}); err != nil {
		panic(err)
	}
}
