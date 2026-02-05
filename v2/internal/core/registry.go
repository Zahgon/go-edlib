package core

import (
	"fmt"
	"sync"

	"github.com/hbollon/go-edlib/v2/algorithms"
)

// AlgorithmFactory creates calculator instances for an algorithm
type AlgorithmFactory interface {
	// CreateCalculator creates a new calculator instance
	CreateCalculator(opts ...CalculatorOption) (Calculator, error)

	// Algorithm returns the algorithm type
	Algorithm() algorithms.Algorithm

	// Properties returns the algorithm properties
	Properties() algorithms.AlgorithmProperties
}

// Registry manages algorithm factories
type Registry struct {
	factories map[algorithms.Algorithm]AlgorithmFactory
	mu        sync.RWMutex
}

var (
	globalRegistry     *Registry
	globalRegistryOnce sync.Once
)

// GetGlobalRegistry returns the global registry instance
func GetGlobalRegistry() *Registry {
	globalRegistryOnce.Do(func() {
		globalRegistry = NewRegistry()
	})
	return globalRegistry
}

// NewRegistry creates a new algorithm registry
func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[algorithms.Algorithm]AlgorithmFactory),
	}
}

// Register registers an algorithm factory
func (r *Registry) Register(factory AlgorithmFactory) error {
	if factory == nil {
		return fmt.Errorf("factory cannot be nil")
	}

	algo := factory.Algorithm()
	if !algo.IsValid() {
		return fmt.Errorf("invalid algorithm: %v", algo)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.factories[algo]; exists {
		return fmt.Errorf("algorithm %s already registered", algo)
	}

	r.factories[algo] = factory
	return nil
}

// Unregister removes an algorithm factory
func (r *Registry) Unregister(algo algorithms.Algorithm) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.factories, algo)
}

// Get retrieves an algorithm factory
func (r *Registry) Get(algo algorithms.Algorithm) (AlgorithmFactory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factory, exists := r.factories[algo]
	if !exists {
		return nil, fmt.Errorf("algorithm %s not registered", algo)
	}

	return factory, nil
}

// CreateCalculator creates a calculator for the specified algorithm
func (r *Registry) CreateCalculator(algo algorithms.Algorithm, opts ...CalculatorOption) (Calculator, error) {
	factory, err := r.Get(algo)
	if err != nil {
		return nil, err
	}

	return factory.CreateCalculator(opts...)
}

// IsRegistered checks if an algorithm is registered
func (r *Registry) IsRegistered(algo algorithms.Algorithm) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.factories[algo]
	return exists
}

// ListAlgorithms returns all registered algorithms
func (r *Registry) ListAlgorithms() []algorithms.Algorithm {
	r.mu.RLock()
	defer r.mu.RUnlock()

	algos := make([]algorithms.Algorithm, 0, len(r.factories))
	for algo := range r.factories {
		algos = append(algos, algo)
	}

	return algos
}

// Clear removes all registered factories
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories = make(map[algorithms.Algorithm]AlgorithmFactory)
}

// Global registry convenience functions

// Register registers an algorithm factory in the global registry
func Register(factory AlgorithmFactory) error {
	return GetGlobalRegistry().Register(factory)
}

// CreateCalculator creates a calculator using the global registry
func CreateCalculator(algo algorithms.Algorithm, opts ...CalculatorOption) (Calculator, error) {
	return GetGlobalRegistry().CreateCalculator(algo, opts...)
}

// IsRegistered checks if an algorithm is registered in the global registry
func IsRegistered(algo algorithms.Algorithm) bool {
	return GetGlobalRegistry().IsRegistered(algo)
}

// ListAlgorithms returns all algorithms registered in the global registry
func ListAlgorithms() []algorithms.Algorithm {
	return GetGlobalRegistry().ListAlgorithms()
}
