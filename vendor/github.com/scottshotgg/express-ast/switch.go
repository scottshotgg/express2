package ast

import (
	"fmt"

	"github.com/scottshotgg/express-token"
)

// Switch statements represents the following form:
// `switch` { expression } [ case_block ]
type Switch struct {
	Token      token.Token
	Expression Expression
	Cases      *CaseBlock
	Default    Statement
}

// CaseBlock represents the following form:
// `{` [ case ]* `}`
type CaseBlock struct {
	Cases []Case
}

// Case represents the following form:
// `case` [ expression ] `:` [ block ]
type Case struct {
	Token      token.Token
	Expression Expression
	Body       Statement
}

// Implement Node and Statement

func (s *Switch) statementNode() {}

// TokenLiteral returns the literal value of the token
func (s *Switch) TokenLiteral() token.Token { return s.Token }

func (s *Switch) Kind() NodeType { return SwitchNode }

func (s *Switch) String() string {
	// FIXME: just doing this to get it to compile
	return fmt.Sprintf("%+v", *s)
}
