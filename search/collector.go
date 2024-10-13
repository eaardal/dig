package search

import "context"

func GroupSearchResults(ctx context.Context, sourceCh <-chan *Result) ([]*Result, error) {
	var results []*Result

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()

		case result, ok := <-sourceCh:
			if !ok {
				return results, nil
			}

			results = append(results, result)
		}
	}
}
