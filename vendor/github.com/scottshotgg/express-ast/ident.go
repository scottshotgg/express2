package ast

import (
	"errors"
	"fmt"

	token "github.com/scottshotgg/express-token"
)

// Ident represents the following form:
// [ name ]
type Ident struct {
	Token  token.Token
	TypeOf *Type
	Name   string
}

func (i *Ident) expressionNode() {}

// TokenLiteral returns the literal value of the token
func (i *Ident) TokenLiteral() token.Token { return i.Token }

func (i *Ident) Kind() NodeType { return IdentNode }

func (i Ident) String() string {
	if i.TypeOf.Name == token.StringType {
		return "std::" + i.TypeOf.Name + " " + i.Name
	} else if i.TypeOf.Name == token.ObjectType {
		// Object is just a var on the backend
		return "object " + i.Name
	}

	return i.TypeOf.Name + " " + i.Name
}

func (i *Ident) Type() *Type {
	// TODO: this might be kind of weird to use since idents
	//  don't always have a type on the first pass
	return i.TypeOf
}

// Might need to make specific type-functions
// But I don't think identifiers here need to have a type, that's NOT what the AST is for; keep track of that in the parser, etc

// NewIdent returns a new identifier
// func NewIdent(t token.Token, it Type) (*Ident, error) {
func NewIdent(t token.Token, it string) (*Ident, error) {
	if t.Value.String == "" {
		return nil, errors.New("Cannot use empty string as identifier name")
	}

	fmt.Println("SHIT", t.Value.String)

	switch it {
	case "int":
		return NewIntIdent(t)

	case "bool":
		return NewBoolIdent(t)

	case "float":
		return NewFloatIdent(t)

	case "char":
		return NewCharIdent(t)

	case "string":
		return NewStringIdent(t)

	case "var":
		return NewVarIdent(t)

	case "object":
		return NewObjectIdent(t)

	case "int[]":
		id, err := NewIntIdent(t)
		if err != nil {
			return nil, err
		}
		id.TypeOf.Array = true

		return id, nil

	case "string[]":
		id, err := NewStringIdent(t)
		if err != nil {
			return nil, err
		}
		id.TypeOf.Array = true

		return id, nil

	default:
		// if strings.Contains(it, "[]") {
		// 	fmt.Println("i am here")
		// 	return &Ident{
		// 		Token:  t,
		// 		TypeOf: NewArrayType(TypeFromString(strings.Replace(it, "[]", "", 1)), true),
		// 		Name:   t.Value.String,
		// 	}, nil
		// } else {

		return &Ident{
			Token:  t,
			TypeOf: &Type{},
			Name:   t.Value.String,
		}, nil
		// }

		return nil, errors.New("NewIdent did not have this type: " + it)
	}
}

// NewIntIdent returns a new identifier for an int type
func NewIntIdent(t token.Token) (*Ident, error) {
	if t.Value.String == "" {
		return nil, errors.New("Cannot use empty string as identifier name")
	}

	return &Ident{
		Token:  t,
		TypeOf: NewIntType(),
		Name:   t.Value.String,
	}, nil
}

// NewBoolIdent returns a new identifier for an bool type
func NewBoolIdent(t token.Token) (*Ident, error) {
	if t.Value.String == "" {
		return nil, errors.New("Cannot use empty string as identifier name")
	}

	return &Ident{
		Token:  t,
		TypeOf: NewBoolType(),
		Name:   t.Value.String,
	}, nil
}

// NewFloatIdent returns a new identifier for an float type
func NewFloatIdent(t token.Token) (*Ident, error) {
	if t.Value.String == "" {
		return nil, errors.New("Cannot use empty string as identifier name")
	}

	return &Ident{
		Token:  t,
		TypeOf: NewFloatType(),
		Name:   t.Value.String,
	}, nil
}

// NewCharIdent returns a new identifier for an char type
func NewCharIdent(t token.Token) (*Ident, error) {
	if t.Value.String == "" {
		return nil, errors.New("Cannot use empty string as identifier name")
	}

	return &Ident{
		Token:  t,
		TypeOf: NewCharType(),
		Name:   t.Value.String,
	}, nil
}

// NewStringIdent returns a new identifier for an string type
func NewStringIdent(t token.Token) (*Ident, error) {
	if t.Value.String == "" {
		return nil, errors.New("Cannot use empty string as identifier name")
	}

	return &Ident{
		Token:  t,
		TypeOf: NewStringType(),
		Name:   t.Value.String,
	}, nil
}

// func NewStructIdent(t token.Token) (*Ident, error) {
// 	if t.Value.String == "" {
// 		return nil, errors.New("Cannot use empty string as identifier name")
// 	}

// 	return &Ident{
// 		Token: t,
// 		TypeOf:  NewStructType(),
// 		Name:  n,
// 	}, nil
// }

// NewObjectIdent returns a new identifier for an object type
func NewObjectIdent(t token.Token) (*Ident, error) {
	if t.Value.String == "" {
		return nil, errors.New("Cannot use empty string as identifier name")
	}

	return &Ident{
		Token:  t,
		TypeOf: NewObjectType(),
		Name:   t.Value.String,
	}, nil
}

// NewFunctionIdent returns a new identifier for an function type
func NewFunctionIdent(t token.Token) (*Ident, error) {
	if t.Value.String == "" {
		return nil, errors.New("Cannot use empty string as identifier name")
	}

	return &Ident{
		Token:  t,
		TypeOf: NewFunctionType(),
		Name:   t.Value.String,
	}, nil
}

func NewVarIdent(t token.Token) (*Ident, error) {
	if t.Value.String == "" {
		return nil, errors.New("Cannot use empty string as identifier name")
	}

	return &Ident{
		Token: t,
		// Set the var type to nothing; "0"
		TypeOf: NewVarType(NewNoneType()),
		Name:   t.Value.String,
	}, nil
}
