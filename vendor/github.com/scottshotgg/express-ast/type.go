package ast

import (
	"fmt"
	"os"
)

// LiteralType encompasses all types of literals
type LiteralType int

const (
	NoneType LiteralType = iota + 1

	// DefaultType directs the parser to apply the default value
	// for the declaration and is used in places where the value
	// tokens for the declaration are not specified
	DefaultType

	// IntType denotes an integer literal type
	IntType

	// BoolType denotes a bool literal type
	BoolType

	// FloatType denotes a float literal type
	FloatType

	// CharType denotes a char literal type
	CharType

	// StringType denotes a string literal type
	StringType

	// StructType denotes a struct literal type
	StructType

	// ObjectType denotes an object literal type
	ObjectType

	// FunctionType denotes a function literal type
	FunctionType

	// VarType denotes a var literal type
	VarType

	ArrayType

	// UserDefinedType denotes a type user defined type
	UserDefinedType
)

// Type is used to specify a variable type
type Type struct {
	Name  string
	Type  LiteralType
	Array bool

	// UpgradesTo is used to specify how/what a variable can upgrade to
	UpgradesTo LiteralType

	// ShadowType is mainly used in dynamic types to specify what the real type is
	ShadowType *LiteralType
}

func (t *Type) Kind() NodeType { return TypeNode }

func (t *Type) String() string {
	// FIXME: just doing this to get it to compile
	return fmt.Sprintf("%+v", *t)
}

var (
	nameToUserDefinedTypeMap = map[string]*Type{}
	idToUserDefinedTypeMap   = map[LiteralType]*Type{}

	// Any type ID greater than 99 is a user defined type
	typeIndex LiteralType = 99

	// UpgradableTypesMap allows definitions of upgradable types
	UpgradableTypesMap = map[LiteralType]LiteralType{
		IntType:    FloatType,
		CharType:   StringType,
		StructType: ObjectType,
	}
)

// DeclareUserDefinedType declares a user defined type in the
// type map and returns a type ID
func DeclareUserDefinedType(udt *Type) LiteralType {
	typeIndex++

	udt.Type = typeIndex

	nameToUserDefinedTypeMap[udt.Name] = udt
	idToUserDefinedTypeMap[udt.Type] = udt

	return typeIndex
}

// NewIntType is used to take some of the boilerplate code out of defining an int Type
func NewIntType() *Type {
	return &Type{
		Name:       "int",
		Type:       IntType,
		UpgradesTo: FloatType,
	}
}

// NewBoolType is used to take some of the boilerplate code out of defining a bool Type
func NewBoolType() *Type {
	return &Type{
		Name:       "bool",
		Type:       BoolType,
		UpgradesTo: NoneType,
	}
}

// NewFloatType is used to take some of the boilerplate code out of defining a float Type
func NewFloatType() *Type {
	return &Type{
		Name:       "float",
		Type:       FloatType,
		UpgradesTo: NoneType,
	}
}

// NewCharType is used to take some of the boilerplate code out of defining a char Type
func NewCharType() *Type {
	return &Type{
		Name:       "char",
		Type:       CharType,
		UpgradesTo: StringType,
	}
}

// NewStringType is used to take some of the boilerplate code out of defining a string Type
func NewStringType() *Type {
	return &Type{
		Name:       "string",
		Type:       StringType,
		UpgradesTo: NoneType,
	}
}

// NewVarType is used to take some of the boilerplate code out of defining an var Type
func NewVarType(lt LiteralType) *Type {
	// somehow need to gaurantee that the shadow type is not `var`
	return &Type{
		Name:       "var",
		Type:       VarType,
		ShadowType: &lt,
		// UpgradesTo: UpgradableTypesMap[lt],
	}
}

// NewObjectType is used to take some of the boilerplate code out of defining an object Type
func NewObjectType() *Type {
	return &Type{
		Name:       "object",
		Type:       ObjectType,
		UpgradesTo: NoneType,
	}
}

// NewStructType is used to take some of the boilerplate code out of defining a struct Type
func NewStructType(lt LiteralType) *Type {
	if _, ok := idToUserDefinedTypeMap[lt]; !ok {
		// FIXME: fix this later or something
		fmt.Printf("Not able to find %d in map during struct inititializer\n", lt)
		panic("oh shit brah")
	}

	thing := StructType
	return &Type{
		Name:       "struct",
		Type:       lt,
		ShadowType: &thing,
		UpgradesTo: ObjectType,
	}
}

// NewFunctionType is used to take some of the boilerplate code out of defining a function Type
func NewFunctionType() *Type {
	return &Type{
		Name: "function",
		Type: FunctionType,
	}
}

func NewArrayType(t *Type, homogenous bool) *Type {
	var ty *Type

	// FIXME: is it even needed to return the type
	if homogenous {
		switch t.Type {
		case IntType:
			ty = NewIntType()

		case FloatType:
			ty = NewFloatType()

		case ObjectType:
			ty = NewObjectType()

		case StringType:
			ty = NewStringType()

		// This is here for a "list" type array essentially
		case VarType:
			ty = t

		default:
			fmt.Printf("This type was not implemented in NewArrayType: %+v", t)
			os.Exit(9)

		}
	} else {
		// TODO: just make it a var for now
		// FIXME: we will need to change the runtime to accept this
		// FIXME: this probably needs to be fixed ...
		ty = &Type{
			Name: "var",
			Type: VarType,
		}
	}

	if ty == nil {
		fmt.Printf("something happened %+v %+v", t, homogenous)
		os.Exit(9)
	}

	ty.Array = true
	return ty
}
