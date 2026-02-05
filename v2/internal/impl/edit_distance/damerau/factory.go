package damerau

import (
	"github.com/hbollon/go-edlib/v2/algorithms"
	"github.com/hbollon/go-edlib/v2/internal/core"
)

// Factory creates Damerau-Levenshtein calculators (full variant)
type Factory struct{}

// CreateCalculator creates a new Damerau-Levenshtein calculator
func (f *Factory) CreateCalculator(opts ...core.CalculatorOption) (core.Calculator, error) {
	return New(opts...)
}

// Algorithm returns the algorithm identifier
func (f *Factory) Algorithm() algorithms.Algorithm {
	return algorithms.DamerauLevenshtein
}

// Properties returns the algorithm properties
func (f *Factory) Properties() algorithms.AlgorithmProperties {
	return algorithms.DamerauLevenshtein.Properties()
}

var _ core.AlgorithmFactory = (*Factory)(nil)

// OSAFactory creates OSA Damerau-Levenshtein calculators
type OSAFactory struct{}

// CreateCalculator creates a new OSA Damerau-Levenshtein calculator
func (f *OSAFactory) CreateCalculator(opts ...core.CalculatorOption) (core.Calculator, error) {
	return NewOSA(opts...)
}

// Algorithm returns the algorithm identifier
func (f *OSAFactory) Algorithm() algorithms.Algorithm {
	return algorithms.OSADamerauLevenshtein
}

// Properties returns the algorithm properties
func (f *OSAFactory) Properties() algorithms.AlgorithmProperties {
	return algorithms.OSADamerauLevenshtein.Properties()
}

var _ core.AlgorithmFactory = (*OSAFactory)(nil)

func init() {
	if err := core.Register(&Factory{}); err != nil {
		panic(err)
	}
	if err := core.Register(&OSAFactory{}); err != nil {
		panic(err)
	}
}
