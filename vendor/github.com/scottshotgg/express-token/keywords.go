package token

// KeywordMap is a map of all the keywords
var KeywordMap = map[string]Token{
	"let": Token{
		Type: Let,
		Value: Value{
			Type:   "keyword",
			String: "let",
		},
	},

	"type": Token{
		Type: TypeDef,
		Value: Value{
			Type:   "keyword",
			String: "type",
		},
	},

	"struct": Token{
		Type: Struct,
		Value: Value{
			Type:   "keyword",
			String: "struct",
		},
	},

	"object": Token{
		Type: Object,
		Value: Value{
			Type:   "keyword",
			String: "object",
		},
	},

	"package": Token{
		Type: Package,
		Value: Value{
			Type:   "keyword",
			String: "package",
		},
	},

	// "c": Token{
	// 	Type: C,
	// 	Value: Value{
	// 		Type:   "keyword",
	// 		String: "c",
	// 	},
	// },

	"use": Token{
		Type: Use,
		Value: Value{
			Type:   "keyword",
			String: "Use",
		},
	},

	"import": Token{
		Type: Import,
		Value: Value{
			Type:   "keyword",
			String: "import",
		},
	},

	"include": Token{
		Type: Include,
		Value: Value{
			Type:   "keyword",
			String: "include",
		},
	},

	"launch": Token{
		Type: Launch,
		Value: Value{
			Type:   "keyword",
			String: "launch",
		},
	},

	"enum": Token{
		Type: Enum,
		Value: Value{
			Type:   "keyword",
			String: "enum",
		},
	},

	"select": Token{
		ID:   9,
		Type: "SELECT",
		Value: Value{
			Type:   "keyword", // TODO: what to put here?
			String: "select",
		},
	},

	"for": Token{
		ID:   9,
		Type: "FOR",
		Value: Value{
			Type:   "loop", // TODO: what to put here?
			String: "for",
		},
	},

	"if": Token{
		ID:   9,
		Type: "IF",
		Value: Value{
			Type:   "keyword", // TODO: what to put here?
			String: "if",
		},
	},

	"else": Token{
		ID:   9,
		Type: "ELSE",
		Value: Value{
			Type:   "keyword", // TODO: what to put here?
			String: "else",
		},
	},

	"in": Token{
		ID:   9,
		Type: "KEYWORD",
		Value: Value{
			Type:   "keyword", // TODO: what to put here?
			String: "in",
		},
	},

	"of": Token{
		ID:   9,
		Type: "KEYWORD",
		Value: Value{
			Type:   "keyword", // TODO: what to put here?
			String: "of",
		},
	},

	"over": Token{
		ID:   9,
		Type: "KEYWORD",
		Value: Value{
			Type:   "keyword", // TODO: what to put here?
			String: "over",
		},
	},

	"function": Token{
		ID:   9,
		Type: "FUNCTION",
		Value: Value{
			Type:   "keyword", // TODO: what to put here?
			String: "function",
		},
	},

	"func": Token{
		ID:   9,
		Type: "FUNCTION",
		Value: Value{
			Type:   "keyword", // TODO: what to put here?
			String: "func",
		},
	},

	"fn": Token{
		ID:   9,
		Type: "FN",
		Value: Value{
			Type:   "keyword", // TODO: what to put here?
			String: "fn",
		},
	},

	"return": Token{
		ID:   9,
		Type: "RETURN",
		Value: Value{
			Type:   "keyword", // TODO: what to put here?
			String: "return",
		},
	},

	"onexit": Token{
		ID:   9,
		Type: "ONEXIT",
		Value: Value{
			Type: "keyword", // TODO: what to put here?
			// String: OnExit,
			String: "onexit",
		},
	},

	"onreturn": Token{
		ID:   9,
		Type: "ONRETURN",
		Value: Value{
			Type: "keyword", // TODO: what to put here?
			// String: OnReturn,
			String: "onreturn",
		},
	},

	"onleave": Token{
		ID:   9,
		Type: "ONLEAVE",
		Value: Value{
			Type: "keyword", // TODO: what to put here?
			// String: OnLeave,
			String: "onleave",
		},
	},

	"defer": Token{
		ID:   9,
		Type: "DEFER",
		Value: Value{
			Type: "keyword", // TODO: what to put here?
			// String: Defer,
			String: "defer",
		},
	},
}
