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

	// InterfaceType denotes a struct literal type
	InterfaceType

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

var (
	ltMap = map[LiteralType]string{
		NoneType:    "none",
		DefaultType: "default",
		IntType:     "int",
		BoolType:    "bool",
		FloatType:   "float",
		CharType:    "char",
		StringType:  "string",
		ObjectType:  "object",
		VarType:     "var",

		// // TODO: these needs sub-typeing
		// ArrayType:       "array::?",
		// StructType:      "struct::?",
		// FunctionType:    "function::?",
		// UserDefinedType: "user_defined::?",
	}
)

func (lt LiteralType) String() string {
	return ltMap[lt]
}

// func LiteralTypeFromString(s string) LiteralType {
// }

// Type is used to specify a variable type
type Type struct {
	Name  string
	Type  LiteralType
	Array bool

	// UpgradesTo is used to specify how/what a variable can upgrade to
	UpgradesTo *Type

	// ShadowType is mainly used in dynamic types to specify what the real type is
	ShadowType *Type
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

// func NewType()

// NewIntType is used to take some of the boilerplate code out of defining an int Type
func NewNoneType() *Type {
	return &Type{
		Name: "none",
		Type: NoneType,
	}
}

// NewIntType is used to take some of the boilerplate code out of defining an int Type
func NewIntType() *Type {
	return &Type{
		Name:       "int",
		Type:       IntType,
		UpgradesTo: NewFloatType(),
	}
}

// NewBoolType is used to take some of the boilerplate code out of defining a bool Type
func NewBoolType() *Type {
	return &Type{
		Name: "bool",
		Type: BoolType,
	}
}

// NewFloatType is used to take some of the boilerplate code out of defining a float Type
func NewFloatType() *Type {
	return &Type{
		Name: "float",
		Type: FloatType,
	}
}

// NewCharType is used to take some of the boilerplate code out of defining a char Type
func NewCharType() *Type {
	return &Type{
		Name:       "char",
		Type:       CharType,
		UpgradesTo: NewStringType(),
	}
}

// NewStringType is used to take some of the boilerplate code out of defining a string Type
func NewStringType() *Type {
	return &Type{
		Name: "string",
		Type: StringType,
	}
}

// NewVarType is used to take some of the boilerplate code out of defining an var Type
func NewVarType(lt *Type) *Type {
	// somehow need to guarantee that the shadow type is not `var`
	return &Type{
		Name:       "var",
		Type:       VarType,
		ShadowType: lt,
		// UpgradesTo: UpgradableTypesMap[lt],
	}
}

// NewObjectType is used to take some of the boilerplate code out of defining an object Type
func NewObjectType() *Type {
	return &Type{
		Name: "object",
		Type: ObjectType,
	}
}

// NewStructType is used to take some of the boilerplate code out of defining a struct Type
func NewStructType(lt LiteralType) *Type {
	var t = &Type{
		Name:       "struct",
		Type:       lt,
		ShadowType: NewStructType(0),
		UpgradesTo: NewObjectType(),
	}

	if lt != 0 {
		var ok bool
		t.ShadowType, ok = idToUserDefinedTypeMap[lt]
		if !ok {
			// FIXME: fix this later or something
			fmt.Printf("Not able to find %d in map during struct inititializer\n", lt)
			panic("oh shit brah")
		}
	}

	return t
}

// NewInterfaceType is used to take some of the boilerplate code out of defining a struct Type
func NewInterfaceType(lt LiteralType) *Type {
	panic("NewInterfaceType not implemented")
	// var t = &Type{
	// 	Name:       "interface",
	// 	Type:       lt,
	// 	ShadowType: NewStructType(0),
	// 	UpgradesTo: NewObjectType(),
	// }

	// if lt != 0 {
	// 	var ok bool
	// 	t.ShadowType, ok = idToUserDefinedTypeMap[lt]
	// 	if !ok {
	// 		// FIXME: fix this later or something
	// 		fmt.Printf("Not able to find %d in map during struct inititializer\n", lt)
	// 		panic("oh shit brah")
	// 	}
	// }

	// return t
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
