package ast

import (
	"fmt"

	token "github.com/scottshotgg/express-token"
)

// LoopType encompasses all types of loops
type LoopType int

const (
	_ LoopType = iota

	// StdFor is a standard for loop containing a condition and expression
	// operating around a declared variable:
	// `for` [ statement ] [ condition ] [ expression ] [ block ]
	StdFor

	// ForEver is the result of a blank StdFor loop:
	// `for` [ block ]
	ForEver

	// ForIn is a for loop that auto iterates over the keys of an iterable:
	// `for` [ literal ] `in` [ iterable ] [ block ]
	ForIn

	// ForOf is a for loop that auto iterates over the values of an iterable:
	// `for` [ literal ] `in` [ iterable ] [ block ]
	ForOf

	// ForOver is a for loop that auto iterates over the keys and values of an iterable:
	// `for` [ literal ] `,` [ literal ] `in` [ iterable ] [ block ]
	// `for` [ -object- ] `in` [ iterable ] [ block ]
	ForOver

	// While is a loop that operates only on a single required condition:
	// `while` [ condition ] [ block ]
	While

	// Until is a reverse-logic while loop:
	// `until` [ condition ] [ block ]
	Until
)

// Loop represents the following form:
// [ loop_type ] { iterable } [ block ]
type Loop struct {
	Token token.Token
	Type  LoopType
	Init  *Assignment
	Cond  Expression
	Post  Expression
	Body  *Block
	Iter  *Iterable
	Temps map[string]*Ident
}

// Implement expression
func (l *Loop) statementNode() {}

// TokenLiteral returns the literal value of the token
func (l *Loop) TokenLiteral() token.Token { return l.Token }

func (l *Loop) Kind() NodeType { return LoopNode }

func (l *Loop) String() string {
	// FIXME: just doing this to get it to compile
	return fmt.Sprintf("%+v", *l)
}
