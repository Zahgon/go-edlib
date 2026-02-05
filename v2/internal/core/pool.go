package core

import (
	"sync"
)

// MemoryPool manages reusable memory buffers to reduce allocations
type MemoryPool struct {
	intSlicePool    sync.Pool
	int2DSlicePool  sync.Pool
	runeSlicePool   sync.Pool
	boolSlicePool   sync.Pool
	maxSliceSize    int
	hits            int64
	misses          int64
	mu              sync.RWMutex
}

// NewMemoryPool creates a new memory pool with the specified maximum slice size
func NewMemoryPool(maxSliceSize int) *MemoryPool {
	if maxSliceSize <= 0 {
		maxSliceSize = 10000 // Default max size
	}

	return &MemoryPool{
		maxSliceSize: maxSliceSize,
		intSlicePool: sync.Pool{
			New: func() interface{} {
				return make([]int, 0)
			},
		},
		int2DSlicePool: sync.Pool{
			New: func() interface{} {
				return make([][]int, 0)
			},
		},
		runeSlicePool: sync.Pool{
			New: func() interface{} {
				return make([]rune, 0)
			},
		},
		boolSlicePool: sync.Pool{
			New: func() interface{} {
				return make([]bool, 0)
			},
		},
	}
}

var (
	globalPool     *MemoryPool
	globalPoolOnce sync.Once
)

// GetGlobalMemoryPool returns the global memory pool instance
func GetGlobalMemoryPool() *MemoryPool {
	globalPoolOnce.Do(func() {
		globalPool = NewMemoryPool(10000)
	})
	return globalPool
}

// GetIntSlice retrieves an int slice from the pool with at least the specified capacity
func (p *MemoryPool) GetIntSlice(capacity int) []int {
	if capacity > p.maxSliceSize {
		p.recordMiss()
		return make([]int, capacity)
	}

	slice := p.intSlicePool.Get().([]int)
	
	if cap(slice) < capacity {
		p.recordMiss()
		return make([]int, capacity)
	}

	p.recordHit()
	return slice[:capacity]
}

// PutIntSlice returns an int slice to the pool
func (p *MemoryPool) PutIntSlice(slice []int) {
	if cap(slice) <= p.maxSliceSize {
		// Clear the slice
		for i := range slice {
			slice[i] = 0
		}
		p.intSlicePool.Put(slice[:0])
	}
}

// GetInt2DSlice retrieves a 2D int slice from the pool
func (p *MemoryPool) GetInt2DSlice(rows, cols int) [][]int {
	if rows > p.maxSliceSize || cols > p.maxSliceSize {
		p.recordMiss()
		matrix := make([][]int, rows)
		for i := range matrix {
			matrix[i] = make([]int, cols)
		}
		return matrix
	}

	slice := p.int2DSlicePool.Get().([][]int)
	
	if cap(slice) < rows {
		p.recordMiss()
		matrix := make([][]int, rows)
		for i := range matrix {
			matrix[i] = make([]int, cols)
		}
		return matrix
	}

	p.recordHit()
	slice = slice[:rows]
	for i := range slice {
		if cap(slice[i]) < cols {
			slice[i] = make([]int, cols)
		} else {
			slice[i] = slice[i][:cols]
		}
	}

	return slice
}

// PutInt2DSlice returns a 2D int slice to the pool
func (p *MemoryPool) PutInt2DSlice(slice [][]int) {
	if cap(slice) <= p.maxSliceSize {
		// Clear the slices
		for i := range slice {
			if cap(slice[i]) <= p.maxSliceSize {
				for j := range slice[i] {
					slice[i][j] = 0
				}
			}
		}
		p.int2DSlicePool.Put(slice[:0])
	}
}

// GetRuneSlice retrieves a rune slice from the pool
func (p *MemoryPool) GetRuneSlice(capacity int) []rune {
	if capacity > p.maxSliceSize {
		p.recordMiss()
		return make([]rune, capacity)
	}

	slice := p.runeSlicePool.Get().([]rune)
	
	if cap(slice) < capacity {
		p.recordMiss()
		return make([]rune, capacity)
	}

	p.recordHit()
	return slice[:capacity]
}

// PutRuneSlice returns a rune slice to the pool
func (p *MemoryPool) PutRuneSlice(slice []rune) {
	if cap(slice) <= p.maxSliceSize {
		// Clear the slice
		for i := range slice {
			slice[i] = 0
		}
		p.runeSlicePool.Put(slice[:0])
	}
}

// GetBoolSlice retrieves a bool slice from the pool
func (p *MemoryPool) GetBoolSlice(capacity int) []bool {
	if capacity > p.maxSliceSize {
		p.recordMiss()
		return make([]bool, capacity)
	}

	slice := p.boolSlicePool.Get().([]bool)
	
	if cap(slice) < capacity {
		p.recordMiss()
		return make([]bool, capacity)
	}

	p.recordHit()
	return slice[:capacity]
}

// PutBoolSlice returns a bool slice to the pool
func (p *MemoryPool) PutBoolSlice(slice []bool) {
	if cap(slice) <= p.maxSliceSize {
		// Clear the slice
		for i := range slice {
			slice[i] = false
		}
		p.boolSlicePool.Put(slice[:0])
	}
}

// recordHit increments the hit counter
func (p *MemoryPool) recordHit() {
	p.mu.Lock()
	p.hits++
	p.mu.Unlock()
}

// recordMiss increments the miss counter
func (p *MemoryPool) recordMiss() {
	p.mu.Lock()
	p.misses++
	p.mu.Unlock()
}

// Stats returns pool statistics
func (p *MemoryPool) Stats() (hits, misses int64) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.hits, p.misses
}

// HitRate returns the pool hit rate [0.0 to 1.0]
func (p *MemoryPool) HitRate() float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	total := p.hits + p.misses
	if total == 0 {
		return 0.0
	}
	return float64(p.hits) / float64(total)
}

// Reset resets the pool statistics
func (p *MemoryPool) Reset() {
	p.mu.Lock()
	p.hits = 0
	p.misses = 0
	p.mu.Unlock()
}
