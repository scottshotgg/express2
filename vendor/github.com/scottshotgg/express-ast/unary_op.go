package ast

import (
	"fmt"

	token "github.com/scottshotgg/express-token"
)

type UnaryType int

const (
	_ UnaryType = iota
	Increment
	Decrement
)

type UnaryOp struct {
	Token token.Token
	Op    UnaryType
	Value Expression
}

func (u *UnaryOp) expressionNode() {}

// TokenLiteral returns the literal value of the token
func (u *UnaryOp) TokenLiteral() token.Token { return u.Token }

func (u *UnaryOp) Kind() NodeType { return UnaryNode }

func (u *UnaryOp) String() string {
	// FIXME: just doing this to get it to compile
	return fmt.Sprintf("%+v", *u)
}

func (u *UnaryOp) Type() *Type {
	return u.Value.Type()
}
