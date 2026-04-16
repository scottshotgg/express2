package ast

import (
	"fmt"

	token "github.com/scottshotgg/express-token"
)

// IfElse represents the following form:
// if [ condition ] [ block ] { [ else ] [ statement ] }
type IfElse struct {
	Token token.Token
	// IfCondition   *Condition
	Condition Expression
	Body      *Block
	// ElseCondition *Condition
	// ElseCondition Expression
	// TODO: Hmmm this is supposed to only be a block or another if statement
	// but should we try to bound it?
	// Else *ElseStatement
	Else *IfElse
}

func (ie *IfElse) statementNode()     {}
func (ie *IfElse) elseStatementNode() {}

// TokenLiteral returns the literal value of the token
func (ie *IfElse) TokenLiteral() token.Token { return ie.Token }

func (ie *IfElse) Kind() NodeType { return IfElseNode }

func (ie *IfElse) String() string {
	// FIXME: just doing this to get it to compile
	return fmt.Sprintf("%+v", *ie)
}
