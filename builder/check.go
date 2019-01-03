package builder

// Check attempts to ensure that node makes sense when viewed with a particular scope tree
func Check(n *Node, st *ScopeTree) {
	/*
		Need to do:

		1) Type checking

		2) Index/bounds checking
		3) Selection checking
		4) Nil checking
		5) Unused type/function/variable checking
		6) Unused path execution
		7 Optimization?
	*/

	// Type check works generically on a node and just decends recursively down all lineage paths

	// Switch on the type to figure out what scope tree we need to use next when checking the children
	switch n.Type {
	// Anything with a block has to go here
	default:
		// do nothing and just pass st
	}

	// Check the left node
	if n.Left != nil {
		Check(n.Left, st)
	}

	// Check the right node
	if n.Right != nil {
		Check(n.Right, st)
	}

	// Value could be a node as well, in which case we need to check it
	var node, ok = n.Value.(*Node)
	if ok {
		Check(node, st)
	} else {
		// Could be an array of nodes, in which case we need to check all of them
		nodes, ok := n.Value.([]*Node)
		if ok {
			for i := range nodes {
				Check(nodes[i], st)
			}
		}
	}
}
