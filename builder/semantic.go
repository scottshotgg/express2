package builder

import (
	"github.com/pkg/errors"
)

var baseTypes = map[string]struct{}{
	"int":    {},
	"char":   {},
	"bool":   {},
	"string": {},
	"float":  {},
}

type Checker struct {
	ps        []Pass
	ast       *Node
	scopeTree *ScopeTree
}

func NewChecker(ast *Node, p ...Pass) *Checker {
	return &Checker{
		ps:  p,
		ast: ast,
	}
}

// AddPass appends a pass to the checker's pipeline.
func (c *Checker) AddPass(p Pass) {
	c.ps = append(c.ps, p)
}

// TODO: think about this returning a report
func (c *Checker) Execute() (*Node, error) {
	for _, pass := range c.ps {
		_, err := pass.Check(c.ast)
		if err != nil {
			return nil, err
		}
	}

	return c.ast, nil
}

type CheckFunc func(n *Node) (bool, error)

type Pass interface {
	Check(n *Node) (bool, error)
	Name() string
}

type DummyPass struct {
	name string
}

func (d *DummyPass) Check(n *Node) (bool, error) {
	return true, nil
}

func (d *DummyPass) Name() string {
	return d.name
}

func NewDummyPass(name string) *DummyPass {
	return &DummyPass{
		name: name,
	}
}

// ---------------------------- real thing ----------------------------

// TypeResolver is a pass that attempts to resolve idents into types in the case of declaration statements
type TypeResolver struct {
	scopeTree *ScopeTree
}

func NewTypeResolver() *TypeResolver {
	return &TypeResolver{
		scopeTree: NewScopeTree(),
	}
}

func NewTypeResolverWithScope(scopeTree *ScopeTree) *TypeResolver {
	return &TypeResolver{
		scopeTree: scopeTree,
	}
}

func (t *TypeResolver) Check(n *Node) (bool, error) {
	switch n.Type {
	case "call":
	case "selection":
	case "assignment":
	case "inc":
	case "dec":
	case "if":
	case "for":
	case "forin":
	case "forof":
	case "binop":
	case "comp":
	case "deref":
	case "ref":
	case "index":
	case "ident":
	case "type":
	case "return":
	case "struct":
	case "while":
	case "kv":
	case "map":
	case "not":
		// later on when we return structs, user-defined types, and object
		// then we can do that
		// Nothing needs to be done RIGHT NOW in this case for a return

	case "egroup":
		// egroup has to be here for function returns
		// In this case it is just a list of idents that are actually types
		for i, returnType := range n.Value.([]*Node) {
			tv := t.scopeTree.GetType(returnType.Value.(string))
			if tv == nil {
				return false, errors.Errorf("could not alias to unfound type: %s", n.Left.Value.(string))
			}

			n.Value.([]*Node)[i] = &Node{
				Type: "type",
				Kind: tv.Kind,
			}
		}

	case "literal":
		// Literals are already typed, no processing needed
		return false, nil

	case "sgroup":
		// For the sgroup, this represents function arguments, so just
		// process them like normal statements
		for _, stmt := range n.Value.([]*Node) {
			changed, err := t.Check(stmt)
			if err != nil {
				return changed, err
			}
		}

	// TODO: this also needs to check that the type is not already defined
	case "typedef":
		switch n.Right.Type {
		case "ident":
			tv := t.scopeTree.GetType(n.Right.Value.(string))
			if tv == nil {
				return false, errors.Errorf("could not alias to unfound type: %s", n.Left.Value.(string))
			}

			var err = t.scopeTree.NewType(n.Left.Value.(string), tv)
			if err != nil {
				return false, err
			}

		case "selection":
			var (
				packageName = n.Right.Value.(*Node).Left.Value.(string)
				typeName    = n.Right.Value.(*Node).Right.Value.(string)
			)

			tv := t.scopeTree.GetImportedType(packageName, typeName)
			if tv != nil {
				return false, errors.Errorf("could not find type %s from package %s", typeName, packageName)
			}

			var err = t.scopeTree.NewType(n.Left.Value.(string), tv)
			if err != nil {
				return false, err
			}

			return false, errors.New("selection types: not implemented")
		}

	// create the scope trees for these in the if/loop/etc
	case "block":
		fallthrough

	case "program":
		var stmts = n.Value.([]*Node)

		for _, stmt := range stmts {
			changed, err := t.Check(stmt)
			if err != nil {
				return changed, err
			}
		}

	case "decl":
		var shouldBeType = *n.Value.(*Node)

		switch shouldBeType.Type {
		case "selection":
			// look in imports as that is the only other place types can be declared
			// this should only ever be [package].[type]
			var packageOf = shouldBeType.Left.Value.(string)
			var typeOf = shouldBeType.Right.Value.(string)

			// Skip the C types for now, change this later when the meta package is taken out
			if packageOf == "c" {
				n.Value = shouldBeType.Right
			} else {
				var tv = t.scopeTree.GetImportedType(packageOf, typeOf)
				shouldBeType.Type = tv.Kind
				n.Value = &Node{
					Type:  "type",
					Kind:  "imported",
					Value: shouldBeType,
				}

				var err = t.scopeTree.Declare(n)
				if err != nil {
					return false, err
				}
			}

		case "ident":
			tv := t.scopeTree.GetType(shouldBeType.Value.(string))
			if tv == nil {
				return false, errors.Errorf("could not find type: %s", n.Left.Value.(string))
			}

			var node = n.Value.(*Node)
			node.Type = "type"
			n.Value = node
		}

	case "let":
		// For let statements, we infer the type from the right-hand side expression
		_, err := t.Check(n.Right)
		if err != nil {
			return false, err
		}

		var typeNode *Node

		switch n.Right.Type {
		case "literal":
			typeNode = &Node{
				Type: "type",
				Kind: n.Right.Kind,
			}

		case "ident":
			tv := t.scopeTree.GetType(n.Right.Value.(string))
			if tv != nil {
				typeNode = &Node{
					Type: "type",
					Kind: tv.Kind,
				}
			} else {
				decl := t.scopeTree.Get(n.Right.Value.(string))
				if decl != nil && decl.Type == "decl" && decl.Value != nil {
					typeNode = decl.Value.(*Node)
				} else {
					return false, errors.Errorf("could not infer type for let variable: %s", n.Left.Value.(string))
				}
			}

		case "call":
			typeNode = &Node{
				Type: "type",
				Kind: "unknown",
			}

		case "type":
			typeNode = n.Right

		default:
			typeNode = &Node{
				Type: "type",
				Kind: "unknown",
			}
		}

		n.Value = typeNode

		err = t.scopeTree.Declare(n)
		if err != nil {
			return false, err
		}

	case "function":
		var (
			err     error
			changed bool
		)

		t.scopeTree, err = t.scopeTree.NewChildScope(n.Kind)
		if err != nil {
			return false, err
		}

		if n.Metadata["args"] != nil {
			changed, err = t.Check(n.Metadata["args"].(*Node))
			if err != nil {
				return changed, err
			}
		} else {
			n.Metadata["args"] = &Node{
				Type:  "sgroup",
				Value: []*Node{},
			}
		}

		if n.Metadata["returns"] != nil {
			changed, err = t.Check(n.Metadata["returns"].(*Node))
			if err != nil {
				return changed, err
			}
		} else {
			var returnType = "void"
			if n.Kind == "main" {
				returnType = "int"
			}

			n.Metadata["returns"] = &Node{
				Type: "egroup",
				Value: []*Node{
					{
						Type: "type",
						Kind: returnType,
					},
				},
			}
		}

		changed, err = t.Check(n.Value.(*Node))
		if err != nil {
			return changed, err
		}

		t.scopeTree, err = t.scopeTree.Leave()
		if err != nil {
			return false, err
		}

	case "import":
		var err error
		var packageName string

		if n.Kind == "c" {
			return false, nil
		}

		packageName = n.Right.Value.([]*Node)[0].Left.Value.(string)

		t.scopeTree, err = t.scopeTree.NewPackageScope(packageName)
		if err != nil {
			return false, err
		}

		changed, err := t.Check(n.Right)
		if err != nil {
			return changed, err
		}

		t.scopeTree, err = t.scopeTree.Leave()
		if err != nil {
			return false, err
		}

	case "package":
		var changed, err = t.Check(n.Right)
		if err != nil {
			return changed, err
		}

	default:
		return false, errors.Errorf("type not implemented in %s: %+v", t.Name(), *n)
	}

	return false, nil
}

func (t *TypeResolver) Name() string {
	return "type_resolver"
}
