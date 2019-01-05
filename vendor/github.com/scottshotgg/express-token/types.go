package token

// TypeMap holds all defined type tokens
var TypeMap = map[string]Token{
	VarType: Token{
		Type: Type,
		Value: Value{
			Type:   VarType,
			String: VarType,
		},
	},
	"val": Token{
		Type: Type,
		Value: Value{
			Type:   "val",
			String: "val",
		},
	},
	IntType: Token{
		Type: Type,
		Value: Value{
			Type:   IntType,
			String: IntType,
		},
	},
	FloatType: Token{
		Type: Type,
		Value: Value{
			Type:   FloatType,
			String: FloatType,
		},
	},
	CharType: Token{
		Type: Type,
		Value: Value{
			Type:   CharType,
			String: CharType,
		},
	},
	StringType: Token{
		Type: Type,
		Value: Value{
			Type:   StringType,
			String: StringType,
		},
	},
	BoolType: Token{
		Type: Type,
		Value: Value{
			Type:   BoolType,
			String: BoolType,
		},
	},

	"map": Token{
		Type: Type,
		Value: Value{
			Type:   "map", // this doesn't create a var
			String: "map",
		},
	},

	"stmt": Token{
		Type: Type,
		Value: Value{
			Type:   "stmt", // this doesn't create a var
			String: "stmt",
		},
	},

	// Make object a keyword like struct
	// Left this in here for express-ast to continue working
	ObjectType: Token{
		Type: Type,
		Value: Value{
			Type:   ObjectType,
			String: ObjectType,
		},
	},
	// StructType: Token{
	// 	Type: Type,
	// 	Value: Value{
	// 		Type:   StructType,
	// 		String: StructType,
	// 	},
	// },
	ArrayType: Token{
		Type: Type,
		Value: Value{
			Type:   ArrayType,
			String: ArrayType,
		},
	},
	IntArrayType: Token{
		Type: Type,
		Value: Value{
			Type:   ArrayType,
			Acting: IntType,
			String: ArrayType,
		},
	},
	StringArrayType: Token{
		Type: Type,
		Value: Value{
			Type:   ArrayType,
			Acting: StringType,
			String: ArrayType,
		},
	},
}
