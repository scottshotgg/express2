package ast

import (
	"fmt"

	token "github.com/scottshotgg/express-token"
)

// Access is the action ([ expression ] `.` [ expression ]) represented by the
// dot operator that allows the internals of a struct to be utilized
type Access struct {
	Token   token.Token
	Name    string
	Parents []string
}

func (a *Access) expressionNode() {}

// TokenLiteral returns the literal value of the token
func (a *Access) TokenLiteral() token.Token { return a.Token }

func (a *Access) Kind() NodeType { return AccessNode }

func (a *Access) String() string {
	// FIXME: just doing this to get it to compile
	return fmt.Sprintf("%+v", *a)
}
