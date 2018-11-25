package builder

import (
	"github.com/scottshotgg/express-token"
)

func New(tokens []token.Token) *Builder {
	b := Builder{
		Tokens: tokens,
	}

	b.OpFuncMap = []map[string]opCallbackFn{
		0: map[string]opCallbackFn{
			token.Increment: b.ParseIncrement,
			token.Accessor:  b.ParseSelection,
			token.LBracket:  b.ParseIndexExpression,
			token.LParen:    b.ParseCall,
			token.LThan:     b.ParseConditionExpression,
			token.PriOp:     b.ParseBinOp,
		},

		1: map[string]opCallbackFn{
			token.SecOp: b.ParseBinOp,
		},
	}

	return &b
}

func (b *Builder) BuildAST() (*Node, error) {
	var (
		stmt  *Node
		stmts []*Node
		err   error
	)

	for b.Index < len(b.Tokens)-1 {
		stmt, err = b.ParseStatement()
		if err != nil {
			if err == ErrOutOfTokens {
				break
			}

			return nil, err
		}

		// Just a fallback; probably won't need it later
		if stmt == nil {
			return b.AppendTokenToError("Statement was nil")
		}

		stmts = append(stmts, stmt)
	}

	return &Node{
		Type:  "program",
		Value: stmts,
	}, nil
}