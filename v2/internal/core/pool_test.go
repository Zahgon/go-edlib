package core

import (
	"sync"
	"testing"
)

func TestNewMemoryPool(t *testing.T) {
	tests := []struct {
		name         string
		maxSliceSize int
		wantMax      int
	}{
		{"Default size", 0, 10000},
		{"Custom size", 5000, 5000},
		{"Negative size", -100, 10000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := NewMemoryPool(tt.maxSliceSize)
			if pool.maxSliceSize != tt.wantMax {
				t.Errorf("NewMemoryPool() maxSliceSize = %v, want %v", pool.maxSliceSize, tt.wantMax)
			}
		})
	}
}

func TestMemoryPool_GetPutIntSlice(t *testing.T) {
	pool := NewMemoryPool(10000)

	tests := []struct {
		name     string
		capacity int
	}{
		{"Small capacity", 10},
		{"Medium capacity", 100},
		{"Large capacity", 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slice := pool.GetIntSlice(tt.capacity)
			
			if len(slice) != tt.capacity {
				t.Errorf("GetIntSlice() len = %v, want %v", len(slice), tt.capacity)
			}
			
			if cap(slice) < tt.capacity {
				t.Errorf("GetIntSlice() cap = %v, want >= %v", cap(slice), tt.capacity)
			}
			
			// Modify slice to test clearing
			for i := range slice {
				slice[i] = i + 1
			}
			
			// Return to pool
			pool.PutIntSlice(slice)
			
			// Get again and verify it's cleared
			slice2 := pool.GetIntSlice(tt.capacity)
			for i, v := range slice2 {
				if v != 0 {
					t.Errorf("GetIntSlice() after Put: slice[%d] = %v, want 0 (should be cleared)", i, v)
					break
				}
			}
			pool.PutIntSlice(slice2)
		})
	}
}

func TestMemoryPool_OversizedIntSlice(t *testing.T) {
	pool := NewMemoryPool(100)
	
	// Request slice larger than maxSliceSize
	slice := pool.GetIntSlice(200)
	
	if len(slice) != 200 {
		t.Errorf("GetIntSlice() for oversized len = %v, want 200", len(slice))
	}
	
	// This should not panic, but won't be pooled
	pool.PutIntSlice(slice)
	
	// Check that misses were recorded
	hits, misses := pool.Stats()
	if misses == 0 {
		t.Error("Expected misses > 0 for oversized allocation")
	}
	if hits != 0 {
		t.Errorf("Expected hits = 0, got %v", hits)
	}
}

func TestMemoryPool_GetPutInt2DSlice(t *testing.T) {
	pool := NewMemoryPool(10000)

	tests := []struct {
		name string
		rows int
		cols int
	}{
		{"Small matrix", 5, 5},
		{"Medium matrix", 10, 20},
		{"Large matrix", 50, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matrix := pool.GetInt2DSlice(tt.rows, tt.cols)
			
			if len(matrix) != tt.rows {
				t.Errorf("GetInt2DSlice() rows = %v, want %v", len(matrix), tt.rows)
			}
			
			for i, row := range matrix {
				if len(row) != tt.cols {
					t.Errorf("GetInt2DSlice() matrix[%d] cols = %v, want %v", i, len(row), tt.cols)
				}
			}
			
			// Modify matrix to test clearing
			for i := range matrix {
				for j := range matrix[i] {
					matrix[i][j] = i*tt.cols + j
				}
			}
			
			// Return to pool
			pool.PutInt2DSlice(matrix)
			
			// Get again and verify it's cleared
			matrix2 := pool.GetInt2DSlice(tt.rows, tt.cols)
			for i := range matrix2 {
				for j, v := range matrix2[i] {
					if v != 0 {
						t.Errorf("GetInt2DSlice() after Put: matrix[%d][%d] = %v, want 0", i, j, v)
						return
					}
				}
			}
			pool.PutInt2DSlice(matrix2)
		})
	}
}

func TestMemoryPool_GetPutRuneSlice(t *testing.T) {
	pool := NewMemoryPool(10000)

	tests := []struct {
		name     string
		capacity int
	}{
		{"Small capacity", 10},
		{"Medium capacity", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slice := pool.GetRuneSlice(tt.capacity)
			
			if len(slice) != tt.capacity {
				t.Errorf("GetRuneSlice() len = %v, want %v", len(slice), tt.capacity)
			}
			
			// Modify slice
			for i := range slice {
				slice[i] = rune('a' + i)
			}
			
			// Return to pool
			pool.PutRuneSlice(slice)
			
			// Get again and verify it's cleared
			slice2 := pool.GetRuneSlice(tt.capacity)
			for i, v := range slice2 {
				if v != 0 {
					t.Errorf("GetRuneSlice() after Put: slice[%d] = %v, want 0", i, v)
					break
				}
			}
			pool.PutRuneSlice(slice2)
		})
	}
}

func TestMemoryPool_GetPutBoolSlice(t *testing.T) {
	pool := NewMemoryPool(10000)

	tests := []struct {
		name     string
		capacity int
	}{
		{"Small capacity", 10},
		{"Medium capacity", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slice := pool.GetBoolSlice(tt.capacity)
			
			if len(slice) != tt.capacity {
				t.Errorf("GetBoolSlice() len = %v, want %v", len(slice), tt.capacity)
			}
			
			// Modify slice
			for i := range slice {
				slice[i] = true
			}
			
			// Return to pool
			pool.PutBoolSlice(slice)
			
			// Get again and verify it's cleared
			slice2 := pool.GetBoolSlice(tt.capacity)
			for i, v := range slice2 {
				if v != false {
					t.Errorf("GetBoolSlice() after Put: slice[%d] = %v, want false", i, v)
					break
				}
			}
			pool.PutBoolSlice(slice2)
		})
	}
}

func TestMemoryPool_Stats(t *testing.T) {
	pool := NewMemoryPool(100)
	
	// Initially should be zero
	hits, misses := pool.Stats()
	if hits != 0 || misses != 0 {
		t.Errorf("Initial Stats() = (%v, %v), want (0, 0)", hits, misses)
	}
	
	// Get and put a slice within limit (should be a miss first, then hit)
	slice1 := pool.GetIntSlice(50)
	pool.PutIntSlice(slice1)
	
	hits, misses = pool.Stats()
	if misses == 0 {
		t.Error("Expected at least one miss after first GetIntSlice()")
	}
	
	// Get again (should reuse from pool and record a hit)
	slice2 := pool.GetIntSlice(50)
	pool.PutIntSlice(slice2)
	
	hits2, _ := pool.Stats()
	if hits2 <= hits {
		t.Errorf("Expected hits to increase after reusing from pool")
	}
	
	// Get oversized slice (should record a miss)
	pool.GetIntSlice(200)
	
	_, misses3 := pool.Stats()
	if misses3 <= misses {
		t.Error("Expected misses to increase for oversized allocation")
	}
}

func TestMemoryPool_HitRate(t *testing.T) {
	pool := NewMemoryPool(100)
	
	// Initially should be zero
	if rate := pool.HitRate(); rate != 0.0 {
		t.Errorf("Initial HitRate() = %v, want 0.0", rate)
	}
	
	// Generate some hits and misses
	slice1 := pool.GetIntSlice(50) // miss
	pool.PutIntSlice(slice1)
	
	slice2 := pool.GetIntSlice(50) // hit
	pool.PutIntSlice(slice2)
	
	rate := pool.HitRate()
	if rate <= 0.0 || rate > 1.0 {
		t.Errorf("HitRate() = %v, want value in (0.0, 1.0]", rate)
	}
}

func TestMemoryPool_Reset(t *testing.T) {
	pool := NewMemoryPool(100)
	
	// Generate some activity
	slice := pool.GetIntSlice(50)
	pool.PutIntSlice(slice)
	
	hits, misses := pool.Stats()
	if hits == 0 && misses == 0 {
		t.Skip("No stats to reset")
	}
	
	// Reset pool
	pool.Reset()
	
	hits, misses = pool.Stats()
	if hits != 0 || misses != 0 {
		t.Errorf("After Reset() Stats() = (%v, %v), want (0, 0)", hits, misses)
	}
}

func TestGetGlobalMemoryPool(t *testing.T) {
	pool1 := GetGlobalMemoryPool()
	pool2 := GetGlobalMemoryPool()
	
	if pool1 != pool2 {
		t.Error("GetGlobalMemoryPool() should return same instance")
	}
	
	if pool1.maxSliceSize != 10000 {
		t.Errorf("Global pool maxSliceSize = %v, want 10000", pool1.maxSliceSize)
	}
}

func TestMemoryPool_Concurrent(t *testing.T) {
	pool := NewMemoryPool(10000)
	
	var wg sync.WaitGroup
	numGoroutines := 100
	
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			// Each goroutine gets and puts slices multiple times
			for j := 0; j < 10; j++ {
				slice := pool.GetIntSlice(50)
				
				// Do some work with the slice
				for k := range slice {
					slice[k] = id*100 + j*10 + k
				}
				
				pool.PutIntSlice(slice)
			}
		}(i)
	}
	
	wg.Wait()
	
	// Verify stats were recorded correctly
	hits, misses := pool.Stats()
	total := hits + misses
	expectedTotal := int64(numGoroutines * 10)
	
	if total != expectedTotal {
		t.Errorf("Total pool operations = %v, want %v", total, expectedTotal)
	}
}
