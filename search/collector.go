package search

import (
	"context"
	"sort"
	"time"
)

func GroupSearchResults(ctx context.Context, sourceCh <-chan *Result) ([]*Result, error) {
	var results []*Result

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()

		case result, ok := <-sourceCh:
			if !ok {
				return sortResults(results), nil
			}

			results = append(results, result)
		}
	}
}

func sortResults(results []*Result) []*Result {
	// Copy the results to avoid modifying the original slice
	sorted := make([]*Result, len(results))
	copy(sorted, results)

	sort.Slice(sorted, func(i, j int) bool {
		timeI, errI := time.Parse(time.RFC3339, sorted[i].LogEntry.Time)
		timeJ, errJ := time.Parse(time.RFC3339, sorted[j].LogEntry.Time)

		if errI != nil {
			timeI = time.Time{}
		}
		if errJ != nil {
			timeJ = time.Time{}
		}

		// Sort by the timestamp (ascending order)
		return timeI.Before(timeJ)
	})

	return sorted
}
