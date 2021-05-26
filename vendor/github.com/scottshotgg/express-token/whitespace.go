package token

// WhitespaceMap holds all defined Whitespace tokens
var WhitespaceMap = map[string]Token{
	" ": {
		Type: Whitespace,
		Value: Value{
			Type:   "space",
			String: " ",
		},
	},
	"\t": {
		Type: Whitespace,
		Value: Value{
			Type:   "tab",
			String: "\t",
		},
	},
	"\n": {
		Type: Whitespace,
		Value: Value{
			Type:   "newline",
			String: "\n",
		},
	},
}
