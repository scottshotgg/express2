package token

// Value ...
type (
	Value struct {
		// FIXME: to maintain compatibility, just add an 'actingType' var for now, use a struct later
		Name       string                 `json:",omitempty"`
		Type       string                 `json:",omitempty"`
		Acting     string                 `json:",omitempty"`
		True       interface{}            `json:",omitempty"`
		OpMap      interface{}            `json:",omitempty"`
		String     string                 `json:",omitempty"`
		AccessType string                 `json:",omitempty"`
		Metadata   map[string]interface{} `json:",omitempty"`
	}

	// Token ...
	Token struct {
		ID            int    `json:",omitempty"`
		Type          string `json:",omitempty"`
		Expected      string `json:",omitempty"`
		Value         Value  `json:",omitempty"`
		WSNotRequired bool   `json:",omitempty"`
	}
)

// TokenMap ...
var (
	mapArray = []map[string]Token{
		AssignMap,
		EncloserMap,
		KeywordMap,
		OperatorMap,
		SeparatorMap,
		SQLMap,
		TypeMap,
		WhitespaceMap,
	}

	TokenMap = map[string]Token{}
)

// These public consts are to make the entire compiler consistent without having to use
// string literals. These may be changed to ints in the future
const (
	Var          = "VAR"
	Ident        = "IDENT"
	Type         = "TYPE"
	Let          = "LET"
	TypeDef      = "TYPEDEF"
	Struct       = "STRUCT"
	Interface    = "INTERFACE"
	Object       = "OBJECT"
	Map          = "MAP"
	Whitespace   = "WS"
	Literal      = "LITERAL"
	Attribute    = "ATTRIBUTE"
	Keyword      = "KEYWORD"
	SQL          = "SQL"
	Comma        = "COMMA"
	EOS          = "EOS"
	Separator    = "SEPARATOR"
	Bang         = "BANG"
	At           = "AT"
	Hash         = "HASH"
	Block        = "BLOCK"
	Function     = "FUNCTION"
	Call         = "CALL"
	Return       = "RETURN"
	OnExit       = "ONEXIT"
	OnReturn     = "ONRETURN"
	OnLeave      = "ONLEAVE"
	Defer        = "DEFER"
	Group        = "GROUP"
	Array        = "ARRAY"
	Set          = "SET"
	Assign       = "ASSIGN"
	Init         = "INIT"
	PriOp        = "PRI_OP"
	SecOp        = "SEC_OP"
	Mult         = "MULT"
	LBrace       = "L_BRACE"
	LBracket     = "L_BRACKET"
	LParen       = "L_PAREN"
	LThan        = "L_THAN"
	RBrace       = "R_BRACE"
	RBracket     = "R_BRACKET"
	RParen       = "R_PAREN"
	GThan        = "G_THAN"
	DQuote       = "D_QUOTE"
	SQuote       = "S_QUOTE"
	Pipe         = "PIPE"
	Ampersand    = "AMPERSAND"
	DDBY         = "DDBY"
	Underscore   = "UNDERSCORE"
	QuestionMark = "QM"
	Accessor     = "ACCESSOR"
	IsEqual      = "IS_EQUAL"
	EqOrGThan    = "EQ_OR_GT"
	EqOrLThan    = "EQ_OR_LT"
	Increment    = "INCREMENT"
	Package      = "PACKAGE"
	Use          = "USE"
	C            = "C"
	Import       = "IMPORT"
	Include      = "INCLUDE"
	Thread       = "THREAD"
	Link         = "LINK"
	Enum         = "ENUM"

	VarType         = "var"
	IntType         = "int"
	FloatType       = "float"
	StringType      = "string"
	BoolType        = "bool"
	CharType        = "char"
	ObjectType      = "object"
	StructType      = "struct"
	InterfaceType   = "interface"
	MapType         = "map"
	ArrayType       = "array"
	FunctionType    = "func"
	SetType         = "set"
	IntArrayType    = "int[]"
	StringArrayType = "string[]"

	PublicAccessType  = "public"
	PrivateAccessType = "private"

	Loop = "LOOP"
	For  = "FOR"
	If   = "IF"
	Else = "ELSE"
)

func init() {
	// Load all maps in
	for _, tMap := range mapArray {
		for key, value := range tMap {
			TokenMap[key] = value
		}
	}
}
