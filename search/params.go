package search

type Params struct {
	// The name of the app where the logs originated.
	Origin string
	// The query to search for.
	Query string
	// The log entry's field name to search in.
	FieldName string
	// Should we only search in the message field?
	InMessage bool
	// Should we only search in the data fields?
	InFields bool
	// Should we search in all fields and the message text?
	Anywhere bool
	// Should the search be case-sensitive?
	CaseSensitive bool
}