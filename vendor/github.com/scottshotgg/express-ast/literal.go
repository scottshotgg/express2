package ast

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	token "github.com/scottshotgg/express-token"
)

// Literal is an abstract type that represents a literal value, in constrast with a value-producer, such as an expression
type Literal interface {
	Expression
}

// Literals should have acting types and acting values that get set when the value is set

// IntLiteral represents any non floating-point number
type IntLiteral struct {
	Token  token.Token
	TypeOf *Type
	Value  int
}

func (il *IntLiteral) expressionNode() {}

// TokenLiteral returns the literal value of the token
func (il *IntLiteral) TokenLiteral() token.Token { return il.Token }

// Type implements literal
func (il *IntLiteral) Type() *Type { return il.TypeOf }

func (il *IntLiteral) Kind() NodeType { return LiteralNode }

func (il *IntLiteral) String() string {
	return strconv.Itoa(il.Value)
}

// BoolLiteral represents a variable that is restricted to either a true or false value
type BoolLiteral struct {
	Token  token.Token
	TypeOf *Type
	Value  bool
}

func (bl *BoolLiteral) expressionNode() {}

// TokenLiteral returns the literal value of the token
func (bl *BoolLiteral) TokenLiteral() token.Token { return bl.Token }

// Type implements literal
func (bl *BoolLiteral) Type() *Type { return bl.TypeOf }

func (bl *BoolLiteral) Kind() NodeType { return LiteralNode }

func (bl *BoolLiteral) String() string {
	return strconv.FormatBool(bl.Value)
}

// FloatLiteral represents any floating point number
type FloatLiteral struct {
	Token  token.Token
	TypeOf *Type
	Value  float64
}

func (fl *FloatLiteral) expressionNode() {}

// TokenLiteral returns the literal value of the token
func (fl *FloatLiteral) TokenLiteral() token.Token { return fl.Token }

// Type implements literal
func (fl *FloatLiteral) Type() *Type { return fl.TypeOf }

func (fl *FloatLiteral) Kind() NodeType { return LiteralNode }

func (fl *FloatLiteral) String() string {
	// FIXME: %+v has a default precision of 5 or 6 - we need to fix this
	return strconv.FormatFloat(fl.Value, 'f', -1, 64)
}

// CharLiteral represents a single-character capped string:
// `'` [ _single_character_ ] `'`
type CharLiteral struct {
	Token  token.Token
	TypeOf *Type
	Value  int
}

func (cl *CharLiteral) expressionNode() {}

// TokenLiteral returns the literal value of the token
func (cl *CharLiteral) TokenLiteral() token.Token { return cl.Token }

// Type implements literal
func (cl *CharLiteral) Type() *Type { return cl.TypeOf }

func (cl *CharLiteral) Kind() NodeType { return LiteralNode }

func (cl *CharLiteral) String() string {
	return fmt.Sprintf("%d", cl.Value)
}

// StringLiteral represents a double quoted body of text:
// TODO: how to do a backtick quoted body of text
// `"` [ _text_ ] `"`
type StringLiteral struct {
	Token  token.Token
	TypeOf *Type
	Value  string
}

func (sl *StringLiteral) expressionNode() {}

// TokenLiteral returns the literal value of the token
func (sl *StringLiteral) TokenLiteral() token.Token { return sl.Token }

// Type implements literal
func (sl *StringLiteral) Type() *Type { return sl.TypeOf }

func (sl *StringLiteral) Kind() NodeType { return LiteralNode }

func (sl *StringLiteral) String() string {
	return fmt.Sprintf("\"%s\"", sl.Value)
}

// VarLiteral represents a dynamically typed variable; it can hold anything
type VarLiteral struct {
	Token  token.Token
	TypeOf *Type
	// TODO: could either do it this way or this can reference another literal type
	Value interface{}
}

func (vl *VarLiteral) expressionNode() {}

// TokenLiteral returns the literal value of the token
func (vl *VarLiteral) TokenLiteral() token.Token { return vl.Token }

// Type implements literal
func (vl *VarLiteral) Type() *Type { return vl.TypeOf }

func (vl *VarLiteral) Kind() NodeType { return LiteralNode }

func (vl *VarLiteral) String() string {
	// This should actually call the stringer for the right type
	// do a switch? ORRRR make it's type hold another Literal node
	return fmt.Sprintf("%+v", vl.Value)
}

// ObjectLiteral represents a named block : this produces a variable
type ObjectLiteral struct {
	Token  token.Token
	TypeOf *Type
	// TODO: could either do it this way or make block implement literal and then it can be directly used as a literal
	Value Block

	// Only allow assignment operations inside objects for now
	// Value map[string]Literal
}

func (ol *ObjectLiteral) expressionNode() {}

// TokenLiteral returns the literal value of the token
func (ol *ObjectLiteral) TokenLiteral() token.Token { return ol.Token }

// Type implements literal
func (ol *ObjectLiteral) Type() *Type { return ol.TypeOf }

func (ol *ObjectLiteral) Kind() NodeType { return LiteralNode }

func (ol *ObjectLiteral) String() string {
	// Might want to use JSON for this
	return ol.Value.String()
}

// StructLiteral represents a named object : this produces a type
// TODO: this might need to be moved to the type.go file
// FIXME: this might need to be fixed or something
type StructLiteral struct {
	Token  token.Token
	TypeOf *Type
	Value  map[string]Expression
}

func (sl *StructLiteral) expressionNode() {}

// TokenLiteral returns the literal value of the token
func (sl *StructLiteral) TokenLiteral() token.Token { return sl.Token }

// Type implements literal
func (sl *StructLiteral) Type() *Type { return sl.TypeOf }

func (sl *StructLiteral) Kind() NodeType { return LiteralNode }

func (sl *StructLiteral) String() string {
	// Might want to use JSON for this
	return fmt.Sprintf("%+v", sl.Value)
}

// FunctionLiteral represents a named object : this produces a type
type FunctionLiteral struct {
	Token  token.Token
	TypeOf *Type
	// TODO: could either do it this way or make block implement literal and then it can be directly used as a literal

	// On the backend, a function would essentially just be a block (i.e, object) that is able to be called
	Value Block
}

func (fl *FunctionLiteral) expressionNode() {}

// TokenLiteral returns the literal value of the token
func (fl *FunctionLiteral) TokenLiteral() token.Token { return fl.Token }

// Type implements literal
func (fl *FunctionLiteral) Type() *Type { return fl.TypeOf }

func (fl *FunctionLiteral) Kind() NodeType { return LiteralNode }

func (fl *FunctionLiteral) String() string {
	// This should probably have some specific string representation
	return fmt.Sprintf("%+v", fl.Value)
}

func TypeFromString(t string) *Type {
	if strings.Contains(t, "[]") {
		return NewArrayType(TypeFromString(strings.Replace(t, "[]", "", 1)), true)
	}

	switch t {
	case "int":
		return NewIntType()

	case "bool":
		return NewBoolType()

	case "float":
		return NewFloatType()

	case "char":
		return NewCharType()

	case "string":
		return NewStringType()

	case "var":
		// TODO: fix this - we are just passing an int type here for now
		return NewVarType(NewIntType())

	case "object":
		return NewObjectType()
	}

	fmt.Printf("TYPE WAS NOT DEFINED IN TypeFromString() %+v\n", t)
	os.Exit(9)
	return nil
}

func NewLiteral(t token.Token, ty *Type) Literal {
	if ty.Array {
		fmt.Println("fucking array", ty.Name)
		return &Array{
			Token:       t,
			TypeOf:      NewArrayType(ty, true),
			ElementType: TypeFromString(ty.Name),
			Length:      0,
			Elements:    []Expression{},
			Homogenous:  true,
		}
	}

	switch ty.Type {
	case IntType:
		return NewInt(t, 0)

	case BoolType:
		return NewBool(t, false)

	case FloatType:
		return NewFloat(t, 0.0)

	case CharType:
		return NewChar(t, 0)

	case StringType:
		return NewString(t, "")

	case VarType:
		// The default value for a var is the integer 0 because it reduces memory
		// footprint and is the least complicated value to containerize
		return NewVarFromInt(t, 0)

	case ObjectType:
		return NewObject(t, Block{})
	}

	fmt.Println("TYPE WAS NOT DEFINED IN NewLiteral()")
	os.Exit(9)
	return nil
}

// NewInt returns a new int literal
func NewInt(t token.Token, value int) *IntLiteral {
	return &IntLiteral{
		Token:  t,
		TypeOf: NewIntType(),
		Value:  value,
	}
}

// Make this take a type and initialize the default from the default map
// func NewDefault(t token.Token) *IntLiteral {
// 	return NewIntFromValue(token, 0)
// }

// NewBool returns a new bool literal
func NewBool(t token.Token, value bool) *BoolLiteral {
	return &BoolLiteral{
		Token:  t,
		TypeOf: NewBoolType(),
		Value:  value,
	}
}

// NewFloat returns a new float literal
func NewFloat(t token.Token, value float64) *FloatLiteral {
	return &FloatLiteral{
		Token:  t,
		TypeOf: NewFloatType(),
		Value:  value,
	}
}

// NewChar returns a new char literal
func NewChar(t token.Token, value int) *CharLiteral {
	return &CharLiteral{
		Token:  t,
		TypeOf: NewCharType(),
		Value:  value,
	}
}

// NewString returns a new string literal
func NewString(t token.Token, value string) *StringLiteral {
	return &StringLiteral{
		Token:  t,
		TypeOf: NewStringType(),
		Value:  value,
	}
}

// NewStruct returns a new struct literal
func NewStruct(t token.Token, structType LiteralType, value map[string]Expression) *StructLiteral {
	return &StructLiteral{
		Token:  t,
		TypeOf: NewStructType(structType),

		// This is for the properties of the struct, but somehow we probably need to have a
		// UserDefinedValueMap like we have for the UserDefinedTypeMap
		Value: value,
	}
}

// NewObject returns a new object literal
func NewObject(t token.Token, value Block) *ObjectLiteral {
	return &ObjectLiteral{
		Token:  t,
		TypeOf: NewObjectType(),
		Value:  value,
	}
}

// NewVarFromInt returns a new int shadow-typed var
func NewVarFromInt(t token.Token, value int) *VarLiteral {
	return &VarLiteral{
		Token:  t,
		TypeOf: NewVarType(NewIntType()),
		Value:  value,
	}
}

// NewVarFromBool returns a new bool shadow-typed var
func NewVarFromBool(t token.Token, value bool) *VarLiteral {
	return &VarLiteral{
		Token:  t,
		TypeOf: NewVarType(NewBoolType()),
		Value:  value,
	}
}

// NewVarFromFloat returns a new float shadow-typed var
func NewVarFromFloat(t token.Token, value float64) *VarLiteral {
	return &VarLiteral{
		Token:  t,
		TypeOf: NewVarType(NewFloatType()),
		Value:  value,
	}
}

// NewVarFromChar returns a new char shadow-typed var
func NewVarFromChar(t token.Token, value rune) *VarLiteral {
	return &VarLiteral{
		Token:  t,
		TypeOf: NewVarType(NewCharType()),
		Value:  value,
	}
}

// NewVarFromString returns a new string shadow-typed var
func NewVarFromString(t token.Token, value string) *VarLiteral {
	return &VarLiteral{
		Token:  t,
		TypeOf: NewVarType(NewStringType()),
		Value:  value,
	}
}

// NewVarFromObject returns a new object shadow-typed var
func NewVarFromObject(t token.Token, value Block) *VarLiteral {
	return &VarLiteral{
		Token:  t,
		TypeOf: NewVarType(NewObjectType()),
		Value:  value,
	}
}

// // NewVarFromStruct returns a new struct shadow-typed var
// func NewVarFromStruct(t token.Token, structType LiteralType, value map[string]Expression) *VarLiteral {
// 	return &VarLiteral{
// 		Token: t,
// 		TypeOf: NewVarType(StructType),
// 		Value:  value,
// 	}
// }

// NewVarFromFunction returns a new function shadow-typed var
func NewVarFromFunction(t token.Token, value Block) *VarLiteral {
	return &VarLiteral{
		Token:  t,
		TypeOf: NewVarType(NewFunctionType()),
		Value:  value,
	}
}
