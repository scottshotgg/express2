package builder

import (
	"github.com/pkg/errors"
	"github.com/scottshotgg/express-token"
)

// func (b *Builder) AddType(key, value *Node) error {
// 	switch key.Type {
// 	case "switch":
// 	}

// 	// b.TypeMap[key] =
// }

func (b *Builder) AddPrimitive(key string, value *Node) (*TypeValue, error) {
	// Check the value to make sure it is: int, char, byte, string, bool, float

	if value.Value == nil {
		return nil, errors.New("wtf shit was nil")
	}

	// Check whether this type has already been declared
	if b.ScopeTree.GetType(key) != nil {
		return nil, errors.Errorf("Type is already declared in type map: %s", key)
	}

	var sv, ok = value.Value.(string)
	if !ok {
		return nil, errors.New("Value was not a string")
	}

	// switch sv {
	// case "int":
	// case "float":
	// case "bool":
	// case "char":
	// case "byte":
	// case "string":

	// default:
	// 	return nil, errors.Errorf("Type not defined in AddPrimitive: %s", sv)
	// }

	// Check if the type we are trying to alias is in the type map
	if b.ScopeTree.GetType(sv) == nil {
		return nil, errors.Errorf("Type is not declared in type map: %s", sv)
	}

	var v = &TypeValue{
		Type: PrimitiveValue,
		Kind: sv,
	}

	return v, b.ScopeTree.NewType(key, v)
}

func (b *Builder) AddRepeated(key string, value *Node) (*TypeValue, error) {
	return nil, errors.New("Not implemented: AddRepeated")
}

func (b *Builder) AddStructured(key string, value *Node) (*TypeValue, error) {
	if value.Type != "block" {
		return nil, errors.Errorf("Value of type is not a block: %s", value.Kind)
	}

	var props, err = b.extractPropsFromComposite(value)
	if err != nil {
		return nil, err
	}

	var v = &TypeValue{
		Composite: true,
		Type:      StruturedValue,
		Kind:      value.Kind,
		Props:     props,
	}

	return v, b.ScopeTree.NewType(key, v)
}

func (b *Builder) extractPropsFromStruct(n *Node) (map[string]*TypeValue, error) {
	// The struct actually has decl statements inside of it,
	// so parse each of those nodes for the key : value
	// For now I don't think we need the actual value for anything

	var (
		propsRaw = n.Value.([]*Node)
		propMap  = map[string]*TypeValue{}
	)

	for _, prop := range propsRaw {
		var pv = prop.Value.(*Node)

		// Need to check the type that we extract from here as well
		// make a function for that

		var propType = b.ScopeTree.GetType(pv.Value.(string))
		if propType == nil {
			return nil, errors.Errorf("Type not defined: %s, %+v", pv.Value.(string), pv)
		}

		propMap[prop.Left.Value.(string)] = propType
	}

	return propMap, nil
}

// should really change this to use a 'structBlock' or something if the parser can
// determine the types
func (b *Builder) extractPropsFromComposite(n *Node) (map[string]*TypeValue, error) {
	switch n.Type {
	case "block":
		return b.extractPropsFromStruct(n)

	default:
		return nil, errors.Errorf("Not implemented in extractPropsFromComposite: %s", n.Kind)
	}
}

func (b *Builder) BuildNodeFromTypeValue(t *TypeValue) (*Node, error) {
	if t == nil {
		return nil, errors.Errorf("TypeValue was nil ...")
	}

	switch t.Type {
	case StruturedValue:
		return buildStructureFromTypeValue(t)

	default:
		return nil, errors.Errorf("Not implemented: %+v", t)
	}
}

func buildStructureFromTypeValue(t *TypeValue) (*Node, error) {
	return nil, nil
}

// FIXME: rewrite me ffs
func (b *Builder) ParseType() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Type {
		// If the token is not a type then check the type map to see if it is
		// var t = b.ScopeTree.GetType(b.Tokens[b.Index].Value.String)
		// if t == nil {
		return b.AppendTokenToError("Could not get type")
		// }
	}

	var (
		injectedType = ""

		// We need to inject the original type as well
		t = b.ScopeTree.GetType(b.Tokens[b.Index].Value.String)
	)

	if t == nil {
		return nil, errors.Errorf("Type could not be found in scope: %s", b.Tokens[b.Index].Value.String)
	}

	injectedType = t.Kind

	var (
		err    error
		typeOf = b.Tokens[b.Index].Value.String
		node   = &Node{
			Type:  "type",
			Value: typeOf,
			Kind:  injectedType,
		}
	)

	defer func() {
		b.Index++
	}()

	for b.Index < len(b.Tokens)-1 {
		switch b.Tokens[b.Index+1].Type {

		// Array operator
		case token.LBracket:
			node, err = b.ParseArrayType(typeOf)

		// Pointer operator
		case token.PriOp:
			node, err = b.ParsePointerType(node)

		// TODO: reworking typing from a more expression oriented architecture
		// almost as if they were expressions
		// // Type annotation
		// case token.LThan:
		// 	var n *Node
		// 	n, err = b.ParseAnnotatedType(node)

		default:
			return node, nil
		}

		if err != nil {
			return nil, err
		}

		// Increment over the type
		b.Index++
	}

	return node, nil
}

// TODO: ParseType needs to be completely redone to have the same sort of
// architecture as Expression, that way we can use some of the same
// techniques for parsing as we did in array, group, etc

// // Embed the incoming type node in an annotation
// func (b *Builder) ParseAnnotatedType() (*Node, error) {
// 	// Need to parse the types
// 	// Can have:
// 	//	- <k:v>
// 	//	- a,b

// 	// Increment over the type
// 	b.Index++

// 	// Increment over the lthan
// 	b.Index++

// 	// Parse the type that should be inside of it
// 	var node, err = b.ParseType()
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer func() {
// 		b.Index++
// 	}()

// 	for b.Index < len(b.Tokens)-1 {
// 		switch b.Tokens[b.Index].Type {

// 		// Pair operator
// 		case token.Set:
// 			// Increment over the set
// 			b.Index++
// 			node.Left, err = b.ParseType()
// 			if err != nil {
// 				return nil, err
// 			}

// 			node = &Node{
// 				Type:  "type",
// 				Kind:  "pair",
// 				Value: "pair",
// 				Left:  node,
// 				Right: n,
// 			}

// 		// List operator
// 		// TODO: this needs to be revisited when its actually needed ...
// 		case token.Comma:
// 			n.Left, err = b.ParseType()
// 			if err != nil {
// 				return nil, err
// 			}

// 			node = &Node{
// 				Type:  "type",
// 				Kind:  "list",
// 				Value: "list",
// 				Left:  node,
// 				Right: n,
// 			}

// 		default:
// 			return node, nil
// 		}

// 		if err != nil {
// 			return nil, err
// 		}

// 		// Increment over the type
// 		b.Index++
// 	}

// 	return node, nil

// 	return &Node{
// 		Type:  "type",
// 		Kind:  "annotation",
// 		Value: "annotation",
// 		Left:  n,
// 	}, nil
// }

// Embed the incoming type node in another type node that has pointer
func (b *Builder) ParsePointerType(n *Node) (*Node, error) {
	if n.Type != "type" {
		return nil, errors.Errorf("Node was not a type: %+v", n)
	}

	return &Node{
		Type:  "type",
		Kind:  "pointer",
		Value: "pointer",
		Left:  n,
	}, nil
}

// TODO: this need to be rewritten to take the node type and embed it in
// the array type
func (b *Builder) ParseArrayType(typeOf string) (*Node, error) {
	var dim []*Index

	// Look ahead at the next token here
	for b.Index < len(b.Tokens)-1 &&
		b.Tokens[b.Index+1].Type == token.LBracket {
		// Increment over the type token
		b.Index++

		expr, err := b.ParseArrayExpression()
		if err != nil {
			return nil, err
		}

		if expr.Value == nil {
			return nil, errors.Errorf("Array parse value was nil; %+v", expr)
		}

		nodesAssert, ok := expr.Value.([]*Node)
		if !ok {
			return nil, errors.Errorf("Invalid assertion; %+v", expr)
		}

		var dimValue Index

		switch len(nodesAssert) {

		case 1:
			dimValue.Type = nodesAssert[0].Kind
			dimValue.Value, ok = nodesAssert[0].Value.(int)
			if !ok {
				return nil, errors.Errorf("Could not assert array value to int; %+v", nodesAssert[0].Value)
			}

		case 0:
			dimValue.Type = "none"
			dimValue.Value = -1

		default:
			return nil, ErrMultDimArrInit
		}

		dim = append(dim, &dimValue)
	}

	b.Index++

	return &Node{
		Type:  "type",
		Kind:  typeOf,
		Value: "array",
		Metadata: map[string]interface{}{
			// "type": typeOf,
			"dim": dim,
		},
	}, nil
}
