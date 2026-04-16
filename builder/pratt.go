package builder

import (
	token "github.com/scottshotgg/express-token"
)

// Precedence levels for the Pratt parser, lowest to highest.
type Precedence int

const (
	PrecNone       Precedence = iota // 0 — stopping
	PrecEquality                     // 1 — ==
	PrecComparison                   // 2 — < > <= >=
	PrecAddition                     // 3 — + -
	PrecMultiply                     // 4 — * / % ^
	PrecUnary                        // 5 — prefix * &
	PrecPostfix                      // 6 — ++ . [] ()
)

type (
	prefixParseFn func() (*Node, error)
	infixParseFn  func(*Node) (*Node, error)
)

// precedenceOf returns the precedence for a given token type when it
// appears in infix (or postfix) position.
func precedenceOf(tok token.Token) Precedence {
	switch tok.Type {
	case token.IsEqual:
		return PrecEquality
	case token.LThan, token.GThan, token.EqOrLThan, token.EqOrGThan:
		return PrecComparison
	case token.SecOp:
		return PrecAddition
	case token.PriOp:
		return PrecMultiply
	case token.Increment, token.Decrement, token.Accessor, token.LBracket, token.LParen:
		return PrecPostfix
	default:
		return PrecNone
	}
}

// registerParseFns populates the prefix and infix function maps.
func (b *Builder) registerParseFns() {
	b.prefixParseFns = map[string]prefixParseFn{
		token.Literal:   b.parseLiteral,
		token.Ident:     b.parseIdentExpr,
		token.Type:      b.parseTypePrimary,
		token.PriOp:     b.parsePrefixDeref,
		token.Ampersand: b.parsePrefixRef,
		token.Bang:      b.parsePrefixNot,
		token.LParen:    b.parsePrefixGroup,
		token.LBracket:  b.parsePrefixArray,
		token.LBrace:    b.parsePrefixBlock,
	}

	b.infixParseFns = map[string]infixParseFn{
		token.PriOp:     b.parseInfixMul,
		token.SecOp:     b.parseInfixAdd,
		token.IsEqual:   b.parseInfixComp,
		token.LThan:     b.parseInfixComp,
		token.GThan:     b.parseInfixComp,
		token.EqOrLThan:   b.parseInfixComp,
		token.EqOrGThan:   b.parseInfixComp,
		token.Increment:   b.parsePostfixInc,
		token.Decrement:   b.parsePostfixDec,
		token.Accessor:    b.parseInfixSelection,
		token.LBracket:  b.parseInfixIndex,
		token.LParen:    b.parseInfixCall,
	}
}

// prattParse is the core Pratt parsing loop.
// After it returns, b.Index is on the last consumed token.
func (b *Builder) prattParse(minPrec Precedence) (*Node, error) {
	// 1. Parse prefix
	prefixFn, ok := b.prefixParseFns[b.peek().Type]
	if !ok {
		return nil, b.AppendTokenToError("Could not parse expression from token")
	}

	left, err := prefixFn()
	if err != nil {
		return nil, err
	}

	// 2. Infix loop: peek at the next token (Index+1)
	for !b.atEnd() {
		next := b.peekAt(1)
		nextPrec := precedenceOf(next)
		if nextPrec <= minPrec {
			break
		}

		// When left is a type and next is *, treat * as pointer type modifier (not multiplication).
		// e.g. `Student*` becomes a pointer-to-Student type node.
		if next.Type == token.PriOp && next.Value.String == "*" && left.Type == "type" {
			b.advance() // consume the *
			left, err = b.ParsePointerType(left)
			if err != nil {
				return nil, err
			}
			continue
		}

		// Don't consume * as infix multiplication after a ref (&x) or deref (*x) expression.
		// In Express, pointer expressions followed by * means the * starts a new deref statement.
		// Use parentheses to disambiguate pointer arithmetic: (*x) * y.
		if next.Type == token.PriOp && next.Value.String == "*" && (left.Type == "ref" || left.Type == "deref") {
			break
		}

		infixFn, ok := b.infixParseFns[next.Type]
		if !ok {
			break
		}

		// Advance onto the operator token
		b.advance()

		left, err = infixFn(left)
		if err != nil {
			return nil, err
		}
	}

	// b.Index is on the last consumed token
	return left, nil
}

// --- Prefix parse functions ---

func (b *Builder) parseLiteral() (*Node, error) {
	cur := b.peek()
	return &Node{
		Type:  "literal",
		Kind:  cur.Value.Type,
		Value: cur.Value.True,
	}, nil
}

func (b *Builder) parseIdentExpr() (*Node, error) {
	value := b.peek().Value.String

	// Check scope for type
	n := b.ScopeTree.GetType(value)
	if n != nil {
		// If next token is LBrace, parse struct/type literal
		if !b.atEnd() && b.peekAt(1).Type == token.LBrace {
			b.advance()
			nn, err := b.ParseExpression()
			if err != nil {
				return nil, err
			}

			return &Node{
				Type:  "literal",
				Kind:  n.Kind,
				Value: value,
				Right: nn,
			}, nil
		}

		return &Node{
			Type:  "type",
			Value: value,
		}, nil
	}

	// Check scope for package
	nv := b.ScopeTree.Get(value)
	if nv != nil {
		if nv.Type == "program" {
			return &Node{
				Type:  "package",
				Value: value,
			}, nil
		}
	}

	return &Node{
		Type:  "ident",
		Value: value,
	}, nil
}

func (b *Builder) parseTypePrimary() (*Node, error) {
	return b.ParseType(nil)
}

func (b *Builder) parsePrefixDeref() (*Node, error) {
	if b.peek().Value.String != "*" {
		return nil, b.AppendTokenToError("Could not get deref")
	}

	// Step over the *
	b.advance()

	ident, err := b.prattParse(PrecUnary)
	if err != nil {
		return nil, err
	}

	kind := "ident"
	if ident.Type == "type" {
		kind = "type"
	}

	return &Node{
		Type: "deref",
		Left: ident,
		Kind: kind,
	}, nil
}

func (b *Builder) parsePrefixNot() (*Node, error) {
	// Step over the ! token
	b.advance()

	operand, err := b.prattParse(PrecUnary)
	if err != nil {
		return nil, err
	}

	return &Node{
		Type: "not",
		Left: operand,
	}, nil
}

func (b *Builder) parsePrefixRef() (*Node, error) {
	if b.peek().Value.String != "&" {
		return nil, b.AppendTokenToError("Could not get ref")
	}

	if b.atEnd() || b.peekAt(1).Type != token.Ident {
		return nil, b.AppendTokenToError("Could not get ident to ref")
	}

	// Step over the &
	b.advance()

	ident, err := b.prattParse(PrecUnary)
	if err != nil {
		return nil, err
	}

	return &Node{
		Type: "ref",
		Left: ident,
	}, nil
}

func (b *Builder) parsePrefixGroup() (*Node, error) {
	// ParseGroupOfExpressions expects b.Index on LParen and leaves it on RParen.
	return b.ParseGroupOfExpressions()
}

func (b *Builder) parsePrefixArray() (*Node, error) {
	// Delegate to existing ParseArrayExpression which expects b.Index on LBracket
	return b.ParseArrayExpression()
	// ParseArrayExpression leaves b.Index on RBracket — that's the last consumed token. Good.
}

func (b *Builder) parsePrefixBlock() (*Node, error) {
	// Delegate to existing ParseBlockStatement which expects b.Index on LBrace
	a, c := b.ParseBlockStatement()
	// ParseBlockStatement leaves b.Index one past the RBrace.
	// We need b.Index on the last consumed token (the RBrace), so back up.
	b.Index--
	return a, c
}

// --- Infix parse functions ---

func (b *Builder) parseInfixMul(left *Node) (*Node, error) {
	// b.Index is on the operator token
	op := b.peek().Value.String

	// Step over the operator
	b.advance()

	right, err := b.prattParse(PrecMultiply)
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:  "binop",
		Value: op,
		Left:  left,
		Right: right,
	}, nil
}

func (b *Builder) parseInfixAdd(left *Node) (*Node, error) {
	op := b.peek().Value.String

	b.advance()

	right, err := b.prattParse(PrecAddition)
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:  "binop",
		Value: op,
		Left:  left,
		Right: right,
	}, nil
}

func (b *Builder) parseInfixComp(left *Node) (*Node, error) {
	cur := b.peek()

	// Determine operator value from token
	var op string
	switch cur.Type {
	case token.IsEqual:
		op = "=="
	case token.LThan:
		op = "<"
	case token.GThan:
		op = ">"
	case token.EqOrLThan:
		op = "<="
	case token.EqOrGThan:
		op = ">="
	}

	// Determine precedence for this level
	prec := precedenceOf(cur)

	b.advance()

	right, err := b.prattParse(prec)
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:  "comp",
		Value: op,
		Left:  left,
		Right: right,
	}, nil
}

func (b *Builder) parsePostfixInc(left *Node) (*Node, error) {
	// b.Index is on the ++ token, which is the last consumed token
	return &Node{
		Type: "inc",
		Left: left,
	}, nil
}

func (b *Builder) parsePostfixDec(left *Node) (*Node, error) {
	// b.Index is on the -- token, which is the last consumed token
	return &Node{
		Type: "dec",
		Left: left,
	}, nil
}

func (b *Builder) parseInfixSelection(left *Node) (*Node, error) {
	// b.Index is on the . token
	// Step over the accessor
	b.advance()

	right, err := b.prattParse(PrecPostfix)
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:  "selection",
		Left:  left,
		Right: right,
	}, nil
}

func (b *Builder) parseInfixIndex(left *Node) (*Node, error) {
	// b.Index is on the [ token
	// Step over the [
	b.advance()

	expr, err := b.prattParse(PrecNone)
	if err != nil {
		return nil, err
	}

	// Step over the expression's last consumed token onto the ]
	b.advance()

	if b.peek().Type != token.RBracket {
		return nil, b.AppendTokenToError("expected ] in index expression")
	}

	// b.Index is on ], which is the last consumed token
	return &Node{
		Type:  "index",
		Left:  left,
		Right: expr,
	}, nil
}

func (b *Builder) parseInfixCall(left *Node) (*Node, error) {
	// b.Index is on the ( token
	// ParseGroupOfExpressions expects b.Index on LParen
	args, err := b.ParseGroupOfExpressions()
	if err != nil {
		return nil, err
	}

	// ParseGroupOfExpressions leaves b.Index on RParen
	return &Node{
		Type:  "call",
		Value: left,
		Metadata: map[string]interface{}{
			"args": args,
		},
	}, nil
}
