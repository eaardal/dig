package search

import (
	"context"
	"github.com/eaardal/dig/logentry"
	"strings"
)

func Search(ctx context.Context, params Params, sourceCh <-chan *logentry.LogEntry, sinkCh chan<- *Result) error {
	defer close(sinkCh)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case logEntry, ok := <-sourceCh:
			if !ok {
				return nil
			}

			result := search(logEntry, params)
			sendToSink(ctx, sinkCh, result)
		}
	}
}

func search(logEntry *logentry.LogEntry, params Params) *Result {
	result := &Result{
		IsMatch:        false,
		Params:         params,
		LogEntry:       logEntry,
		MatchLocations: []string{},
	}

	if params.InMessage || params.Anywhere {
		if isValueMatch(logEntry.Message, params.Query, params.CaseSensitive) {
			result.IsMatch = true
			result.MatchLocations = append(result.MatchLocations, MatchLocationMessage)
			return result
		}
	}

	if params.InFields || params.Anywhere {
		for fieldKey, fieldValue := range logEntry.Fields {
			if isValueMatch(fieldValue, params.Query, params.CaseSensitive) {
				result.IsMatch = true
				result.MatchLocations = append(result.MatchLocations, fieldKey)
				return result
			}
		}
	}

	if params.FieldName != "" {
		fieldsToSearch := make(map[string]string)

		for key, value := range logEntry.Fields {
			if strings.HasPrefix(params.FieldName, "*") && strings.HasSuffix(params.FieldName, "*") {
				if strings.Contains(key, params.FieldName[1:len(params.FieldName)-1]) {
					fieldsToSearch[key] = value
				}
			}
			if strings.HasPrefix(params.FieldName, "*") {
				if strings.HasSuffix(key, params.FieldName[1:]) {
					fieldsToSearch[key] = value
				}
			}
			if strings.HasSuffix(params.FieldName, "*") {
				if strings.HasPrefix(key, params.FieldName[:len(params.FieldName)-1]) {
					fieldsToSearch[key] = value
				}
			}
			if key == params.FieldName {
				fieldsToSearch[key] = value
			}
		}

		for fieldName, fieldValue := range fieldsToSearch {
			if isValueMatch(fieldValue, params.Query, params.CaseSensitive) {
				result.IsMatch = true
				result.MatchLocations = append(result.MatchLocations, fieldName)
				return result
			}
		}
	}

	return result
}

func isValueMatch(valueToSearch string, valueToFind string, caseSensitive bool) bool {
	if !caseSensitive {
		valueToSearch = strings.ToLower(valueToSearch)
		valueToFind = strings.ToLower(valueToFind)
	}

	if strings.HasPrefix(valueToFind, "*") && strings.HasSuffix(valueToFind, "*") {
		return strings.Contains(valueToSearch, valueToFind[1:len(valueToFind)-1])
	}
	if strings.HasPrefix(valueToFind, "*") {
		return strings.HasSuffix(valueToSearch, valueToFind[1:])
	}
	if strings.HasSuffix(valueToFind, "*") {
		return strings.HasPrefix(valueToSearch, valueToFind[:len(valueToFind)-1])
	}
	return valueToSearch == valueToFind
}

func sendToSink(ctx context.Context, resultCh chan<- *Result, result *Result) {
	select {
	case <-ctx.Done():
		return
	case resultCh <- result:
	}
}
