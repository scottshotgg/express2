package builder


func (b *Builder) ParseSet(n *Node) (*Node, error) {
	// This will be encountered when we have:
	// <expr> `:` <expr>

	// Step over the set token
	b.Index++

	right, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Step over the expression
	b.Index++

	return &Node{
		Type:  "kv",
		Left:  n,
		Right: right,
	}, nil
}

func (b *Builder) ParseCall(n *Node) (*Node, error) {
	// We are not allowing for named arguments right now
	args, err := b.ParseGroupOfExpressions()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:  "call",
		Value: n,
		Metadata: map[string]interface{}{
			"args": args,
		},
	}, nil
}

