package search

import (
	"context"
	"sync"

	"github.com/hbollon/go-edlib/v2"
	"github.com/hbollon/go-edlib/v2/internal/utils"
	searchpkg "github.com/hbollon/go-edlib/v2/search"
)

// ParallelProcessor handles parallel processing of string comparisons
type ParallelProcessor struct {
	workers int
	ctx     context.Context
}

// NewParallelProcessor creates a new parallel processor
func NewParallelProcessor(workers int, ctx context.Context) *ParallelProcessor {
	if workers <= 0 {
		workers = utils.NumCPU()
	}

	if ctx == nil {
		ctx = context.Background()
	}

	return &ParallelProcessor{
		workers: workers,
		ctx:     ctx,
	}
}

// ProcessBatch processes multiple string pairs in parallel
func (p *ParallelProcessor) ProcessBatch(
	pairs []edlib.StringPair,
	processor func(edlib.StringPair) edlib.ComparisonResult,
) []edlib.ComparisonResult {
	n := len(pairs)
	if n == 0 {
		return nil
	}

	results := make([]edlib.ComparisonResult, n)
	
	// For small batches, process sequentially
	if n < p.workers*2 {
		for i, pair := range pairs {
			results[i] = processor(pair)
		}
		return results
	}

	// Calculate optimal chunk size
	chunkSize := (n + p.workers - 1) / p.workers
	if chunkSize < 1 {
		chunkSize = 1
	}

	var wg sync.WaitGroup
	
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()

			for j := start; j < end; j++ {
				// Check context cancellation
				select {
				case <-p.ctx.Done():
					results[j] = edlib.ComparisonResult{
						Error: p.ctx.Err(),
					}
					return
				default:
				}

				results[j] = processor(pairs[j])
			}
		}(i, end)
	}

	wg.Wait()
	return results
}

// ProcessSearch processes search operations in parallel
func (p *ParallelProcessor) ProcessSearch(
	query string,
	targets []string,
	processor func(string, string, int) searchpkg.Result,
) []searchpkg.Result {
	n := len(targets)
	if n == 0 {
		return nil
	}

	results := make([]searchpkg.Result, n)
	
	// For small target lists, process sequentially
	if n < p.workers*2 {
		for i, target := range targets {
			results[i] = processor(query, target, i)
		}
		return results
	}

	// Calculate optimal chunk size
	chunkSize := (n + p.workers - 1) / p.workers
	if chunkSize < 1 {
		chunkSize = 1
	}

	var wg sync.WaitGroup
	
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()

			for j := start; j < end; j++ {
				// Check context cancellation
				select {
				case <-p.ctx.Done():
					return
				default:
				}

				results[j] = processor(query, targets[j], j)
			}
		}(i, end)
	}

	wg.Wait()
	return results
}

// Workers returns the number of workers
func (p *ParallelProcessor) Workers() int {
	return p.workers
}
