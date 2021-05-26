package builder

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

var baseTypes = map[string]struct{}{
	"int":    {},
	"char":   {},
	"bool":   {},
	"string": {},
	"float":  {},
}

// type Checker interface {
// 	Passes() []Pass
// 	AST() *Node
// 	ScopeTree() *ScopeTree
// }

type Checker struct {
	ps        []Pass
	ast       *Node
	scopeTree *ScopeTree
}

func NewChecker(ast *Node, p ...Pass) *Checker {
	return &Checker{
		ps:        p,
		ast:       ast,
		scopeTree: NewScopeTree(),
	}
}

func (c *Checker) AddPass(p Pass) {
	c.ps = append(c.ps, p)
}

// TODO: think about this returning a report
func (c *Checker) Execute() (*Node, error) {
	for _, pass := range c.ps {
		changed, err := pass.Check(c.ast)
		if err != nil {
			return nil, err
		}

		fmt.Println("changed, err:", changed, err)
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
	fmt.Println("my name is:", d.Name())
	fmt.Println("i pass anything that comes my way because im a dummy")

	astJSON, err := json.Marshal(n)
	if err != nil {
		return false, err
	}

	fmt.Println(string(astJSON))

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

// Define out the way that passes/checker/s will work

// TypeResolver is a pass that attempts to resolve idents into types in the case of declaration statements
type TypeResolver struct {
	scopeTree *ScopeTree
}

func NewTypeResolver() *TypeResolver {
	return &TypeResolver{
		scopeTree: NewScopeTree(),
	}
}

func (t *TypeResolver) Check(n *Node) (bool, error) {
	// Look for anything with a body
	fmt.Printf("node_me %+v", *n)
	switch n.Type {
	case "call":
	case "selection":
	case "assignment":
	case "return":
		// later on when we return structs, user-defined types, and object
		// then we can do that
		// Nothing needs to be done RIGHT NOW in this case for a return
	case "egroup":
		// egroup has to be here for function returns
		// In this case it is just a list of idents that are actually types
		for i, returnType := range n.Value.([]*Node) {

			// Find the original type
			tv := t.scopeTree.GetType(returnType.Value.(string))
			if tv == nil {
				return false, errors.Errorf("could not alias to unfound type: %s", n.Left.Value.(string))
			}

			n.Value.([]*Node)[i] = &Node{
				Type: "type",
				Kind: tv.Kind,
			}
		}

	case "sgroup":
		// For the sgroup, this represents function arguments, so just
		// process them like normal statements
		fmt.Printf("n %+v\n", *n)
		for _, stmt := range n.Value.([]*Node) {
			changed, err := t.Check(stmt)
			if err != nil {
				return changed, err
			}
		}

	// TODO: this also needs to check that the type is not already defined
	case "typedef":
		fmt.Println("got a typedef")
		// Typedef will come out as:
		// type [ident] = [ident|selection]
		switch n.Right.Type {
		case "ident":
			// ensure that this is a type
			// if it is, add n.Left.Value.(string) to the types
			fmt.Println("found a type", n.Left.Value.(string))

			// Find the original type
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

			// The selection has to come from an imported package
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
			fmt.Println("stmt", *stmt)
			changed, err := t.Check(stmt)
			if err != nil {
				return changed, err
			}
		}

	case "decl":
		// resolve the type
		fmt.Println("found a decl statement")
		fmt.Printf("node: %+v\n", *n.Value.(*Node))

		var shouldBeType = *n.Value.(*Node)
		fmt.Println("shouldBeType", shouldBeType.Value)

		switch shouldBeType.Type {
		case "selection":
			// look in imports as that is the only other place types can be declared
			// this should only ever be [package].[type]

			var packageOf = shouldBeType.Left.Value.(string)
			var typeOf = shouldBeType.Right.Value.(string)
			fmt.Println("packageOf, typeOf", packageOf, typeOf)

			// Skip the C types for now, change this later when the meta package is taken out
			if packageOf == "c" {
				n.Value = shouldBeType.Right
			} else {
				fmt.Printf("imports %+v\n", t.scopeTree.GetImports())

				// Set the type
				var tv = t.scopeTree.GetImportedType(packageOf, typeOf)
				shouldBeType.Type = tv.Kind
				n.Value = &Node{
					Type:  "type",
					Kind:  "imported",
					Value: shouldBeType,
				}
				fmt.Printf("n.Value %+v\n", n.Value)

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

			// Grab the pointer, set it as a type, and reset the node value
			var node = n.Value.(*Node)
			node.Type = "type"
			n.Value = node
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

		// Check the args
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

		// Check the returns
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

		// Check the body; this should never be nil
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

		// fix this shit: n.Right.Value.([]*Node)[0].Left.Value
		if n.Kind == "c" {
			return false, nil
		}

		packageName = n.Right.Value.([]*Node)[0].Left.Value.(string)

		t.scopeTree, err = t.scopeTree.NewPackageScope(packageName)
		if err != nil {
			return false, err
		}

		// Right side has the AST for that file
		changed, err := t.Check(n.Right)
		if err != nil {
			return changed, err
		}

		t.scopeTree, err = t.scopeTree.Leave()
		if err != nil {
			return false, err
		}

	case "package":
		// check the right side
		var changed, err = t.Check(n.Right)
		if err != nil {
			return changed, err
		}

	default:
		fmt.Println("got type:", n.Type)
		return false, errors.Errorf("type not implemented in %s: %+v", t.Name(), *n)
	}

	astJSON, err := json.Marshal(n)
	if err != nil {
		return false, err
	}

	fmt.Println(string(astJSON))

	return false, nil
}

func (t *TypeResolver) Name() string {
	return "type_resolver"
}
