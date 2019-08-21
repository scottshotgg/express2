package ast

import token "github.com/scottshotgg/express-token"

func New(tokens []token.Token) *AST {
	var a = AST{
		Tokens: tokens,
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

	// b.OpFuncMap = []map[string]opCallbackFn{
	// 	// Tier 1 operators
	// 	0: map[string]opCallbackFn{
	// 		token.Increment: b.ParseIncrement,
	// 		token.Accessor:  b.ParseSelection,
	// 		token.LBracket:  b.ParseIndexExpression,
	// 		token.LParen:    b.ParseCall,
	// 		token.LThan:     b.ParseLessThanExpression,
	// 		token.GThan:     b.ParseGreaterThanExpression,
	// 		token.EqOrLThan: b.ParseLessOrEqualThanExpression,
	// 		token.EqOrGThan: b.ParseGreaterOrEqualThanExpression,
	// 		token.IsEqual:   b.ParseEqualityExpression,
	// 		token.PriOp:     b.ParseBinOp,
	// 	},

	// 	// Tier 2 operators
	// 	1: map[string]opCallbackFn{
	// 		token.SecOp: b.ParseBinOp,
	// 		// token.Set:   b.ParseSet,
	// 	},
	// }

	return &a
}

func (a *AST) BuildAST() (*Node, error) {
	var (
		// stmt        *Node
		// err         error
		stmts       []*Node
		programNode = &Node{
			Type: "program",
		}
	)

	// for a.Index < len(a.Tokens)-1 {
	// 	stmt, err = a.ParseStatement()
	// 	if err != nil {
	// 		if err == ErrOutOfTokens {
	// 			break
	// 		}

	// 		return nil, err
	// 	}

	// 	// Just a fallback; probably won't need it later
	// 	if stmt == nil {
	// 		return a.AppendTokenToError("Statement was nil")
	// 	}

	// 	stmts = append(stmts, stmt)
	// }

	programNode.Value = stmts

	return programNode, nil
}
