package sorensendice

import (
	"github.com/hbollon/go-edlib/v2/algorithms"
	"github.com/hbollon/go-edlib/v2/internal/core"
)

// Factory creates SorensenDice calculators
type Factory struct{}

// CreateCalculator creates a new SorensenDice calculator
func (f *Factory) CreateCalculator(opts ...core.CalculatorOption) (core.Calculator, error) {
	return New(opts...)
}

// Algorithm returns the algorithm identifier
func (f *Factory) Algorithm() algorithms.Algorithm {
	return algorithms.SorensenDice
}

// Properties returns the algorithm properties
func (f *Factory) Properties() algorithms.AlgorithmProperties {
	return algorithms.SorensenDice.Properties()
}

var _ core.AlgorithmFactory = (*Factory)(nil)

func init() {
	if err := core.Register(&Factory{}); err != nil {
		panic(err)
	}
}
