package ast

import (
	"fmt"

	token "github.com/scottshotgg/express-token"
)

// Import is an import statement in the form of:
// `import` [ string_lit ]
type Import struct {
	Token token.Token
	Name  *Ident
	Path  string
}

// Implement Node and Statement
func (i *Import) statementNode() {}

// TokenLiteral returns the literal value of the token
func (i *Import) TokenLiteral() token.Token { return i.Token }

func (i *Import) Kind() NodeType { return ImportNode }

func (i *Import) String() string {
	// FIXME: just doing this to get it to compile
	return fmt.Sprintf("%+v", *i)
}
