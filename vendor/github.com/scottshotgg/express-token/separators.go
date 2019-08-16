package token

// SeparatorMap holds all defined statement separators
var SeparatorMap = map[string]Token{
	",": Token{
		WSNotRequired: true,
		Type:          "SEPARATOR",
		Value: Value{
			Type:   "comma",
			String: ",",
		},
	},
	";": Token{
		WSNotRequired: true,
		Type:          "SEPARATOR",
		Value: Value{
			Type:   "semicolon",
			String: ";",
		},
	},
}
