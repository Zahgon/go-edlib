package jaro

import (
	"github.com/hbollon/go-edlib/v2/algorithms"
	"github.com/hbollon/go-edlib/v2/internal/core"
)

// Factory creates Jaro calculator instances
type Factory struct{}

// CreateCalculator creates a new Jaro calculator
func (f *Factory) CreateCalculator(opts ...core.CalculatorOption) (core.Calculator, error) {
	return New(opts...)
}

// Algorithm returns the algorithm type
func (f *Factory) Algorithm() algorithms.Algorithm {
	return algorithms.Jaro
}

// Properties returns the algorithm properties
func (f *Factory) Properties() algorithms.AlgorithmProperties {
	return algorithms.Jaro.Properties()
}

// WinklerFactory creates Jaro-Winkler calculator instances
type WinklerFactory struct{}

// CreateCalculator creates a new Jaro-Winkler calculator
func (f *WinklerFactory) CreateCalculator(opts ...core.CalculatorOption) (core.Calculator, error) {
	return NewWinkler(opts...)
}

// Algorithm returns the algorithm type
func (f *WinklerFactory) Algorithm() algorithms.Algorithm {
	return algorithms.JaroWinkler
}

// Properties returns the algorithm properties
func (f *WinklerFactory) Properties() algorithms.AlgorithmProperties {
	return algorithms.JaroWinkler.Properties()
}

// Ensure factories implement AlgorithmFactory
var (
	_ core.AlgorithmFactory = (*Factory)(nil)
	_ core.AlgorithmFactory = (*WinklerFactory)(nil)
)

// Register registers both Jaro and Jaro-Winkler algorithms
func init() {
	if err := core.Register(&Factory{}); err != nil {
		panic(err)
	}
	if err := core.Register(&WinklerFactory{}); err != nil {
		panic(err)
	}
}
