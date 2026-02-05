package core

import (
	"runtime"
	"time"
)

// Metrics contains detailed execution metrics for an operation
type Metrics struct {
	StartTime         time.Time
	EndTime           time.Time
	ExecutionTime     time.Duration
	MemoryBefore      uint64
	MemoryAfter       uint64
	MemoryAllocated   int64
	MemoryFreed       int64
	CacheHits         int
	CacheMisses       int
	PoolHits          int
	PoolMisses        int
	OptimizationsUsed []string
	WorkersUsed       int
}

// NewMetrics creates a new metrics collector
func NewMetrics() *Metrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &Metrics{
		StartTime:         time.Now(),
		MemoryBefore:      m.Alloc,
		OptimizationsUsed: make([]string, 0),
	}
}

// Stop finalizes the metrics collection
func (m *Metrics) Stop() {
	m.EndTime = time.Now()
	m.ExecutionTime = m.EndTime.Sub(m.StartTime)

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	m.MemoryAfter = mem.Alloc

	if m.MemoryAfter > m.MemoryBefore {
		m.MemoryAllocated = int64(m.MemoryAfter - m.MemoryBefore)
	} else {
		m.MemoryFreed = int64(m.MemoryBefore - m.MemoryAfter)
	}
}

// AddOptimization records that an optimization was used
func (m *Metrics) AddOptimization(name string) {
	m.OptimizationsUsed = append(m.OptimizationsUsed, name)
}

// RecordCacheHit records a cache hit
func (m *Metrics) RecordCacheHit() {
	m.CacheHits++
}

// RecordCacheMiss records a cache miss
func (m *Metrics) RecordCacheMiss() {
	m.CacheMisses++
}

// RecordPoolHit records a memory pool hit
func (m *Metrics) RecordPoolHit() {
	m.PoolHits++
}

// RecordPoolMiss records a memory pool miss
func (m *Metrics) RecordPoolMiss() {
	m.PoolMisses++
}

// SetWorkers records the number of workers used
func (m *Metrics) SetWorkers(n int) {
	m.WorkersUsed = n
}

// CacheHitRate returns the cache hit rate [0.0 to 1.0]
func (m *Metrics) CacheHitRate() float64 {
	total := m.CacheHits + m.CacheMisses
	if total == 0 {
		return 0.0
	}
	return float64(m.CacheHits) / float64(total)
}

// PoolHitRate returns the memory pool hit rate [0.0 to 1.0]
func (m *Metrics) PoolHitRate() float64 {
	total := m.PoolHits + m.PoolMisses
	if total == 0 {
		return 0.0
	}
	return float64(m.PoolHits) / float64(total)
}
