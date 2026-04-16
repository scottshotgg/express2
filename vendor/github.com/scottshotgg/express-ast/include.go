package ast

import (
	"fmt"

	token "github.com/scottshotgg/express-token"
)

// Import is an import statement in the form of:
// `include` [ string_lit ]
type Include struct {
	Token token.Token
	Name  *Ident
	Path  string
}

// Implement Node and Statement
func (i *Include) statementNode() {}

// TokenLiteral returns the literal value of the token
func (i *Include) TokenLiteral() token.Token { return i.Token }

func (i *Include) Kind() NodeType { return IncludeNode }

func (i *Include) String() string {
	// FIXME: just doing this to get it to compile
	return fmt.Sprintf("%+v", *i)
}
