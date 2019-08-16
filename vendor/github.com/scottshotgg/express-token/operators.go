package token

// OperatorMap holds all defined operator tokens
var OperatorMap = map[string]Token{
	"+": Token{
		WSNotRequired: true,
		Type:          SecOp,
		Value: Value{
			Type:   "add",
			String: "+",
		},
	},

	"-": Token{
		WSNotRequired: true,
		Type:          SecOp,
		Value: Value{
			Type:   "sub",
			String: "-",
		},
	},

	"*": Token{
		WSNotRequired: true,
		Type:          PriOp,
		Value: Value{
			Type:   "mult",
			String: "*",
		},
	},

	"/": Token{
		WSNotRequired: true,
		Type:          PriOp,
		Value: Value{
			Type:   "div",
			String: "/",
		},
	},

	"\\": Token{
		WSNotRequired: true,
		Type:          PriOp,
		Value: Value{
			Type:   "backslash",
			String: "\\",
		},
	},

	"%": Token{
		WSNotRequired: true,
		Type:          PriOp,
		Value: Value{
			Type:   "mod",
			String: "%",
		},
	},

	"^": Token{
		WSNotRequired: true,
		Type:          PriOp,
		Value: Value{
			Type:   "exp",
			String: "^",
		},
	},

	"!": Token{
		WSNotRequired: true,
		Type:          Bang,
		Value: Value{
			Type:   "bang",
			String: "!",
		},
	},

	"?": Token{
		WSNotRequired: true,
		Type:          QuestionMark,
		Value: Value{
			Type:   "qm",
			String: "!",
		},
	},

	// "_": Token{
	// 	Type: Underscore,
	// 	Value: Value{
	// 		Type:   "underscore",
	// 		String: "_",
	// 	},
	// },
	// FIXME: DOLLA DOLLA BILLS YALL: define this
	"$": Token{
		WSNotRequired: true,
		Type:          DDBY,
		Value: Value{
			Type:   "ddby",
			String: "$",
		},
	},

	"&": Token{
		WSNotRequired: true,
		Type:          Ampersand,
		Value: Value{
			Type:   "op_3",
			String: "&",
		},
	},

	"|": Token{
		WSNotRequired: true,
		Type:          Pipe,
		Value: Value{
			Type:   "op_3",
			String: "|",
		},
	},

	"#": Token{
		WSNotRequired: true,
		Type:          Hash,
		Value: Value{
			Type:   "op_3",
			String: "#",
		},
	},

	".": Token{
		WSNotRequired: true,
		Type:          Accessor,
		Value: Value{
			Type:   "period",
			String: ".",
		},
	},

	"==": Token{
		WSNotRequired: true,
		Type:          IsEqual,
		Value: Value{
			Type:   "is_equal",
			String: "==",
		},
	},

	">=": Token{
		WSNotRequired: true,
		Type:          EqOrGThan,
		Value: Value{
			Type:   "eq or greater than",
			String: ">=",
		},
	},

	"<=": Token{
		WSNotRequired: true,
		Type:          EqOrLThan,
		Value: Value{
			Type:   "eq or less than",
			String: "<=",
		},
	},

	// Increment
	"++": Token{
		WSNotRequired: true,
		Type:          Increment,
		Value: Value{
			Type:   "increment",
			String: "++",
		},
	},

	// TODO: add the templated operators ability to the parser and remove the tokens completely
	// VECTOR OPERANDS
	".+": Token{
		WSNotRequired: true,
		Type:          "VEC_ADD",
		Value: Value{
			Type:   "op_3",
			String: ".+",
		},
	},

	".-": Token{
		WSNotRequired: true,
		Type:          "VEC_SUB",
		Value: Value{
			Type:   "op_4",
			String: ".-",
		},
	},

	".*": Token{
		WSNotRequired: true,
		Type:          "VEC_MULT",
		Value: Value{
			Type:   "op_3",
			String: ".*",
		},
	},

	"./": Token{
		WSNotRequired: true,
		Type:          "VEC_DIV",
		Value: Value{
			Type:   "op_3",
			String: "./",
		},
	},
}
