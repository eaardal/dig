package search

import "github.com/eaardal/dig/logentry"

const (
	MatchLocationMessage = "message"
)

type Result struct {
	IsMatch        bool
	Params         Params
	LogEntry       *logentry.LogEntry
	MatchLocations []string
}
