package token

// SeparatorMap holds all defined statement separators
var SeparatorMap = map[string]Token{
	",": {
		WSNotRequired: true,
		Type:          "SEPARATOR",
		Value: Value{
			Type:   "comma",
			String: ",",
		},
	},
	";": {
		WSNotRequired: true,
		Type:          "SEPARATOR",
		Value: Value{
			Type:   "semicolon",
			String: ";",
		},
	},
}
