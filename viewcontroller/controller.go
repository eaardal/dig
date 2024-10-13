package viewcontroller

import (
	"github.com/eaardal/dig/logentry"
	"github.com/eaardal/dig/search"
)

type ViewEntry struct {
	Origin                       string
	LogEntry                     *logentry.LogEntry
	MatchLocations               []string
	LogEntriesBefore             []*logentry.LogEntry
	LogEntriesAfter              []*logentry.LogEntry
	NumLogEntriesToPreviousMatch int
	NumLogEntriesToNextMatch     int
}

type Options struct {
	NumLogEntriesBefore int
	NumLogEntriesAfter  int
}

type direction int

const (
	before direction = iota
	after
)

func PrepareSearchResultsForDisplay(searchResults []*search.Result, opts Options) ([]*ViewEntry, error) {
	viewEntries := make([]*ViewEntry, 0)

	for index, searchResult := range searchResults {
		if searchResult.IsMatch {
			viewEntries = append(viewEntries, &ViewEntry{
				Origin:                       searchResult.LogEntry.Origin,
				LogEntry:                     searchResult.LogEntry,
				MatchLocations:               searchResult.MatchLocations,
				LogEntriesBefore:             findNearbyLogEntries(searchResults, index, opts.NumLogEntriesBefore, before),
				LogEntriesAfter:              findNearbyLogEntries(searchResults, index, opts.NumLogEntriesAfter, after),
				NumLogEntriesToPreviousMatch: countLogEntriesToNearbyMatch(searchResults, index, before),
				NumLogEntriesToNextMatch:     countLogEntriesToNearbyMatch(searchResults, index, after),
			})
		}
	}

	return viewEntries, nil
}

func findNearbyLogEntries(searchResults []*search.Result, index int, numLogEntries int, directon direction) []*logentry.LogEntry {
	var logEntries []*logentry.LogEntry

	if directon == before {
		for i := index - 1; i >= 0 && len(logEntries) < numLogEntries; i-- {
			// Take log entries before the match, but stop if we encounter another match.
			if searchResults[i].IsMatch {
				break
			}

			logEntries = append(logEntries, searchResults[i].LogEntry)
		}
	} else {
		for i := index + 1; i < len(searchResults) && len(logEntries) < numLogEntries; i++ {
			// Take log entries after the match, but stop if we encounter another match.
			if searchResults[i].IsMatch {
				break
			}

			logEntries = append(logEntries, searchResults[i].LogEntry)
		}
	}

	return logEntries
}

func countLogEntriesToNearbyMatch(searchResults []*search.Result, index int, direction direction) int {
	numLogEntries := 0

	if direction == before {
		for i := index - 1; i >= 0; i-- {
			if searchResults[i].IsMatch {
				break
			}
			numLogEntries++
		}
	} else {
		for i := index + 1; i < len(searchResults); i++ {
			if searchResults[i].IsMatch {
				break
			}
			numLogEntries++
		}
	}

	return numLogEntries
}
