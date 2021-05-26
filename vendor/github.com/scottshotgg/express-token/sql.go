package token

// FIXME: leave this for fixing until we need it

// SQLMap is a map of all the SQL specific tokens
var SQLMap = map[string]Token{
	"SELECT": {
		ID:   9,
		Type: Keyword,
		Value: Value{
			Type:   "SQL", // TODO: what to put here?
			String: "SELECT",
		},
	},
	"FROM": {
		ID:   9,
		Type: Keyword,
		Value: Value{
			Type:   "SQL", // TODO: what to put here?
			String: "FROM",
		},
	},
	"WHERE": {
		ID:   9,
		Type: Keyword,
		Value: Value{
			Type:   "SQL", // TODO: what to put here?
			String: "WHERE",
		},
	},
}
