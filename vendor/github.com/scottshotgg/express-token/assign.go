package token

// AssignMap holds every assignment operator
var AssignMap = map[string]Token{
	"=": Token{
		WSNotRequired: true,
		Type:          Assign,
		Value: Value{
			Type:   "assign",
			String: "=",
		},
	},
	":": Token{
		WSNotRequired: true,
		Type:          Set,
		Value: Value{
			Type:   "set",
			String: ":",
		},
	},
	// ":=": Token{
	// 	Type: "ASSIGN",
	// 	Value: Value{
	// 		Type:   "init",
	// 		String: ":=",
	// 	},
	// },
}
