package token

// OperatorMap holds all defined operator tokens
var OperatorMap = map[string]Token{
	"+": {
		WSNotRequired: true,
		Type:          SecOp,
		Value: Value{
			Type:   "add",
			String: "+",
		},
	},

	"-": {
		WSNotRequired: true,
		Type:          SecOp,
		Value: Value{
			Type:   "sub",
			String: "-",
		},
	},

	"*": {
		WSNotRequired: true,
		Type:          PriOp,
		Value: Value{
			Type:   "mult",
			String: "*",
		},
	},

	"/": {
		WSNotRequired: true,
		Type:          PriOp,
		Value: Value{
			Type:   "div",
			String: "/",
		},
	},

	"\\": {
		WSNotRequired: true,
		Type:          PriOp,
		Value: Value{
			Type:   "backslash",
			String: "\\",
		},
	},

	"%": {
		WSNotRequired: true,
		Type:          PriOp,
		Value: Value{
			Type:   "mod",
			String: "%",
		},
	},

	"^": {
		WSNotRequired: true,
		Type:          PriOp,
		Value: Value{
			Type:   "exp",
			String: "^",
		},
	},

	"!": {
		WSNotRequired: true,
		Type:          Bang,
		Value: Value{
			Type:   "bang",
			String: "!",
		},
	},

	"?": {
		WSNotRequired: true,
		Type:          QuestionMark,
		Value: Value{
			Type:   "qm",
			String: "!",
		},
	},

	// "_": {
	// 	Type: Underscore,
	// 	Value: Value{
	// 		Type:   "underscore",
	// 		String: "_",
	// 	},
	// },
	// FIXME: DOLLA DOLLA BILLS YALL: define this
	"$": {
		WSNotRequired: true,
		Type:          DDBY,
		Value: Value{
			Type:   "ddby",
			String: "$",
		},
	},

	"&": {
		WSNotRequired: true,
		Type:          Ampersand,
		Value: Value{
			Type:   "op_3",
			String: "&",
		},
	},

	"|": {
		WSNotRequired: true,
		Type:          Pipe,
		Value: Value{
			Type:   "op_3",
			String: "|",
		},
	},

	"#": {
		WSNotRequired: true,
		Type:          Hash,
		Value: Value{
			Type:   "op_3",
			String: "#",
		},
	},

	".": {
		WSNotRequired: true,
		Type:          Accessor,
		Value: Value{
			Type:   "period",
			String: ".",
		},
	},

	"==": {
		WSNotRequired: true,
		Type:          IsEqual,
		Value: Value{
			Type:   "is_equal",
			String: "==",
		},
	},

	">=": {
		WSNotRequired: true,
		Type:          EqOrGThan,
		Value: Value{
			Type:   "eq or greater than",
			String: ">=",
		},
	},

	"<=": {
		WSNotRequired: true,
		Type:          EqOrLThan,
		Value: Value{
			Type:   "eq or less than",
			String: "<=",
		},
	},

	// Increment
	"++": {
		WSNotRequired: true,
		Type:          Increment,
		Value: Value{
			Type:   "increment",
			String: "++",
		},
	},

	// TODO: add the templated operators ability to the parser and remove the tokens completely
	// VECTOR OPERANDS
	".+": {
		WSNotRequired: true,
		Type:          "VEC_ADD",
		Value: Value{
			Type:   "op_3",
			String: ".+",
		},
	},

	".-": {
		WSNotRequired: true,
		Type:          "VEC_SUB",
		Value: Value{
			Type:   "op_4",
			String: ".-",
		},
	},

	".*": {
		WSNotRequired: true,
		Type:          "VEC_MULT",
		Value: Value{
			Type:   "op_3",
			String: ".*",
		},
	},

	"./": {
		WSNotRequired: true,
		Type:          "VEC_DIV",
		Value: Value{
			Type:   "op_3",
			String: "./",
		},
	},
}
