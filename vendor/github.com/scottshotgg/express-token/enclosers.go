package token

// EncloserMap holds all valid encloser tokens
var EncloserMap = map[string]Token{
	"(": {
		WSNotRequired: true,
		Type:          "L_PAREN",
		Value: Value{
			Type:   "op_3",
			String: "(",
		},
	},
	")": {
		WSNotRequired: true,
		Type:          "R_PAREN",
		Value: Value{
			Type:   "op_3",
			String: ")",
		},
	},

	"{": {
		WSNotRequired: true,
		Type:          "L_BRACE",
		Value: Value{
			Type:   "op_3",
			String: "{",
		},
	},
	"}": {
		WSNotRequired: true,
		Type:          "R_BRACE",
		Value: Value{
			Type:   "op_3",
			String: "}",
		},
	},

	"[": {
		WSNotRequired: true,
		Type:          "L_BRACKET",
		Value: Value{
			Type:   "op_3",
			String: "[",
		},
	},
	"]": {
		WSNotRequired: true,
		Type:          "R_BRACKET",
		Value: Value{
			Type:   "lthan",
			String: "]",
		},
	},

	"<": {
		WSNotRequired: true,
		Type:          "L_THAN",
		Value: Value{
			Type:   "lthan",
			String: "<",
		},
	},
	">": {
		WSNotRequired: true,
		Type:          "G_THAN",
		Value: Value{
			Type:   "rthan",
			String: ">",
		},
	},

	"`": {
		WSNotRequired: true,
		Type:          "GRAVE",
		Value: Value{
			Type:   "op_3",
			String: "`",
		},
	},
	"~": {
		WSNotRequired: true,
		Type:          "TILDE",
		Value: Value{
			Type:   "op_3",
			String: "~",
		},
	},
	"'": {
		WSNotRequired: true,
		Type:          SQuote,
		Value: Value{
			Type:   "squote",
			String: "'",
		},
	},
	"\"": {
		WSNotRequired: true,
		Type:          DQuote,
		Value: Value{
			Type:   "dquote",
			String: "\"",
		},
	},
	"@": {
		WSNotRequired: true,
		Type:          At,
		Value: Value{
			Type:   "op_3",
			String: "@",
		},
	},
}
