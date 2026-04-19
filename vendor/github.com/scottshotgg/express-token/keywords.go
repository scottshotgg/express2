package token

// KeywordMap is a map of all the keywords
var KeywordMap = map[string]Token{
	"let": {
		Type: Let,
		Value: Value{
			Type:   "keyword",
			String: "let",
		},
	},

	"type": {
		Type: TypeDef,
		Value: Value{
			Type:   "keyword",
			String: "type",
		},
	},

	"struct": {
		Type: Struct,
		Value: Value{
			Type:   "keyword",
			String: "struct",
		},
	},

	"interface": {
		Type: Interface,
		Value: Value{
			Type:   "keyword",
			String: "interface",
		},
	},

	"object": {
		Type: Object,
		Value: Value{
			Type:   "keyword",
			String: "object",
		},
	},

	"package": {
		Type: Package,
		Value: Value{
			Type:   "keyword",
			String: "package",
		},
	},

	// "c" is intentionally not a lexer keyword — it would split identifiers
	// containing the letter c (e.g. func, struct, object, launch).
	// C blocks are recognised in ParseStatement by checking for an ident
	// whose value is "c" followed by an LBrace.

	"use": {
		Type: Use,
		Value: Value{
			Type:   "keyword",
			String: "Use",
		},
	},

	"import": {
		Type: Import,
		Value: Value{
			Type:   "keyword",
			String: "import",
		},
	},

	"include": {
		Type: Include,
		Value: Value{
			Type:   "keyword",
			String: "include",
		},
	},

	"thread": {
		Type: Thread,
		Value: Value{
			Type:   "keyword",
			String: "thread",
		},
	},

	"link": {
		Type: Link,
		Value: Value{
			Type:   "keyword",
			String: "link",
		},
	},

	"enum": {
		Type: Enum,
		Value: Value{
			Type:   "keyword",
			String: "enum",
		},
	},

	"launch": {
		Type: Launch,
		Value: Value{
			Type:   "keyword",
			String: "launch",
		},
	},

	"select": {
		ID:   9,
		Type: "SELECT",
		Value: Value{
			Type:   "keyword", // TODO: what to put here?
			String: "select",
		},
	},

	"for": {
		ID:   9,
		Type: "FOR",
		Value: Value{
			Type:   "loop", // TODO: what to put here?
			String: "for",
		},
	},

	"if": {
		ID:   9,
		Type: "IF",
		Value: Value{
			Type:   "keyword", // TODO: what to put here?
			String: "if",
		},
	},

	"else": {
		ID:   9,
		Type: "ELSE",
		Value: Value{
			Type:   "keyword", // TODO: what to put here?
			String: "else",
		},
	},

	"in": {
		ID:   9,
		Type: "KEYWORD",
		Value: Value{
			Type:   "keyword", // TODO: what to put here?
			String: "in",
		},
	},

	"of": {
		ID:   9,
		Type: "KEYWORD",
		Value: Value{
			Type:   "keyword", // TODO: what to put here?
			String: "of",
		},
	},

	"over": {
		ID:   9,
		Type: "KEYWORD",
		Value: Value{
			Type:   "keyword", // TODO: what to put here?
			String: "over",
		},
	},

	"function": {
		ID:   9,
		Type: "FUNCTION",
		Value: Value{
			Type:   "keyword", // TODO: what to put here?
			String: "function",
		},
	},

	"func": {
		ID:   9,
		Type: "FUNCTION",
		Value: Value{
			Type:   "keyword", // TODO: what to put here?
			String: "func",
		},
	},

	"fn": {
		ID:   9,
		Type: "FN",
		Value: Value{
			Type:   "keyword", // TODO: what to put here?
			String: "fn",
		},
	},

	"return": {
		ID:   9,
		Type: "RETURN",
		Value: Value{
			Type:   "keyword", // TODO: what to put here?
			String: "return",
		},
	},

	"onexit": {
		ID:   9,
		Type: "ONEXIT",
		Value: Value{
			Type: "keyword", // TODO: what to put here?
			// String: OnExit,
			String: "onexit",
		},
	},

	"onreturn": {
		ID:   9,
		Type: "ONRETURN",
		Value: Value{
			Type: "keyword", // TODO: what to put here?
			// String: OnReturn,
			String: "onreturn",
		},
	},

	"onleave": {
		ID:   9,
		Type: "ONLEAVE",
		Value: Value{
			Type: "keyword", // TODO: what to put here?
			// String: OnLeave,
			String: "onleave",
		},
	},

	"defer": {
		ID:   9,
		Type: "DEFER",
		Value: Value{
			Type: "keyword", // TODO: what to put here?
			// String: Defer,
			String: "defer",
		},
	},

	"while": {
		Type: Loop,
		Value: Value{
			Type:   "keyword",
			String: "while",
		},
	},

	"break": {
		ID:   9,
		Type: Break,
		Value: Value{
			Type:   "keyword",
			String: "break",
		},
	},

	"continue": {
		ID:   9,
		Type: Continue,
		Value: Value{
			Type:   "keyword",
			String: "continue",
		},
	},
}
