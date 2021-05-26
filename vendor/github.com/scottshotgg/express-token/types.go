package token

// TypeMap holds all defined type tokens
var TypeMap = map[string]Token{
	VarType: {
		Type: Type,
		Value: Value{
			Type:   VarType,
			String: VarType,
		},
	},
	"val": {
		Type: Type,
		Value: Value{
			Type:   "val",
			String: "val",
		},
	},
	IntType: {
		Type: Type,
		Value: Value{
			Type:   IntType,
			String: IntType,
		},
	},
	FloatType: {
		Type: Type,
		Value: Value{
			Type:   FloatType,
			String: FloatType,
		},
	},
	CharType: {
		Type: Type,
		Value: Value{
			Type:   CharType,
			String: CharType,
		},
	},
	StringType: {
		Type: Type,
		Value: Value{
			Type:   StringType,
			String: StringType,
		},
	},
	BoolType: {
		Type: Type,
		Value: Value{
			Type:   BoolType,
			String: BoolType,
		},
	},

	"map": {
		Type: Type,
		Value: Value{
			Type:   "map", // this doesn't create a var
			String: "map",
		},
	},

	"stmt": {
		Type: Type,
		Value: Value{
			Type:   "stmt", // this doesn't create a var
			String: "stmt",
		},
	},

	// Make object a keyword like struct
	// Left this in here for express-ast to continue working
	// ObjectType: {
	// 	Type: Type,
	// 	Value: Value{
	// 		Type:   ObjectType,
	// 		String: ObjectType,
	// 	},
	// },
	// StructType: {
	// 	Type: Type,
	// 	Value: Value{
	// 		Type:   StructType,
	// 		String: StructType,
	// 	},
	// },
	ArrayType: {
		Type: Type,
		Value: Value{
			Type:   ArrayType,
			String: ArrayType,
		},
	},
	IntArrayType: {
		Type: Type,
		Value: Value{
			Type:   ArrayType,
			Acting: IntType,
			String: ArrayType,
		},
	},
	StringArrayType: {
		Type: Type,
		Value: Value{
			Type:   ArrayType,
			Acting: StringType,
			String: ArrayType,
		},
	},
}
