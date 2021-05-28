package builder

import (
	token "github.com/scottshotgg/express-token"
)

var (
// primTypes = []string{
// 	"int",
// 	"float",
// 	"bool",
// 	"char",
// 	"byte",
// 	"string",
// }
)

func New(tokens []token.Token) *Builder {
	var b = Builder{
		Tokens:    tokens,
		ScopeTree: NewScopeTree(),
		// TypeMap: map[string]*TypeValue{},
	}

	// // Load the concrete types; we'll need to do more
	// // later if we have separate transpilers per file ...
	// for _, value := range primTypes {
	// 	b.TypeMap[value] = &TypeValue{
	// 		Type: PrimitiveValue,
	// 		Kind: value,
	// 	}
	// }

	b.OpFuncMap = []map[string]opCallbackFn{
		// Tier 1 operators
		0: {
			token.Increment: b.ParseIncrement,
			token.Accessor:  b.ParseSelection,
			token.LBracket:  b.ParseIndexExpression,
			// token.LBrace:    b.ParseBlockExpression,
			token.LParen:    b.ParseCall,
			token.LThan:     b.ParseLessThanExpression,
			token.GThan:     b.ParseGreaterThanExpression,
			token.EqOrLThan: b.ParseLessOrEqualThanExpression,
			token.EqOrGThan: b.ParseGreaterOrEqualThanExpression,
			token.PriOp:     b.ParseBinOp,
			// token.Set:       b.ParseSet,
		},

		// Tier 2 operators
		1: {
			token.SecOp: b.ParseBinOp,
			// token.Set:   b.ParseSet,
		},

		// Tier 2 operators
		2: {
			token.IsEqual: b.ParseEqualityExpression,
		},
	}

	return &b
}

func (b *Builder) BuildAST() (*Node, error) {
	var (
		stmt        *Node
		stmts       []*Node
		err         error
		programNode = &Node{
			Type: "program",
		}
	)

	b.ScopeTree = NewScopeTree()

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

	programNode.Value = stmts

	return programNode, nil
}
