package builder

// stmt_exprs.go — thin Pratt prefix wrappers for statement-level constructs.
//
// Each statement parser follows the "statement contract": it leaves b.Index
// one past the last consumed token (past the closing `}` for block statements,
// past the last expression token for others).
//
// Pratt's invariant is the opposite: b.Index ON the last consumed token.
// Every wrapper calls the underlying statement parser and then does b.Index--
// to convert from the statement contract to the Pratt invariant.

func (b *Builder) parseIfExpr() (*Node, error) {
	n, err := b.ParseIfStatement()
	if err != nil {
		return nil, err
	}
	b.Index-- // statement leaves one-past-`}` → Pratt invariant ON `}`
	return n, nil
}

func (b *Builder) parseLetExpr() (*Node, error) {
	n, err := b.ParseLetStatement()
	if err != nil {
		return nil, err
	}
	b.Index-- // statement leaves one-past-last-token → Pratt invariant ON last token
	return n, nil
}

func (b *Builder) parseForExpr() (*Node, error) {
	n, err := b.ParseForStatement()
	if err != nil {
		return nil, err
	}
	b.Index-- // statement leaves one-past-`}` → Pratt invariant ON `}`
	return n, nil
}

func (b *Builder) parseWhileExpr() (*Node, error) {
	n, err := b.ParseWhileStatement()
	if err != nil {
		return nil, err
	}
	b.Index-- // statement leaves one-past-`}` → Pratt invariant ON `}`
	return n, nil
}

func (b *Builder) parseFuncExpr() (*Node, error) {
	n, err := b.ParseFunctionStatement()
	if err != nil {
		return nil, err
	}
	b.Index-- // statement leaves one-past-`}` → Pratt invariant ON `}`
	return n, nil
}

func (b *Builder) parseDeferExpr() (*Node, error) {
	n, err := b.ParseDeferStatement()
	if err != nil {
		return nil, err
	}
	b.Index-- // statement leaves one-past-last-token → Pratt invariant ON last token
	return n, nil
}

func (b *Builder) parseReturnExpr() (*Node, error) {
	n, err := b.ParseReturnStatement()
	if err != nil {
		return nil, err
	}
	b.Index-- // statement leaves one-past-last-token → Pratt invariant ON last token
	return n, nil
}

func (b *Builder) parseBreakExpr() (*Node, error) {
	// b.Index is ON the break token — Pratt invariant: stay here.
	return &Node{Type: "break"}, nil
}

func (b *Builder) parseContinueExpr() (*Node, error) {
	// b.Index is ON the continue token — Pratt invariant: stay here.
	return &Node{Type: "continue"}, nil
}
