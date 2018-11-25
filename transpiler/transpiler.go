package transpiler

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/scottshotgg/express2/builder"
)

type Transpiler struct {
	Name         string
	ASTCloneJSON []byte
	AST          *builder.Node
	Functions    map[string]*builder.Node
	Types        map[string]*builder.Node
	Includes     []*builder.Node
	Imports      []*builder.Node
}

func New(ast *builder.Node, name string) *Transpiler {
	t := Transpiler{
		Name: name,
		AST:  ast,
	}

	t.ASTCloneJSON, _ = json.Marshal(ast)

	return &t
}

func (t *Transpiler) Transpile() (string, error) {
	node := t.AST.Value.([]*builder.Node)[0]

	switch node.Type {

	case "function":
		fmt.Println(node)
		functionString := ""
		returnType := node.Metadata["returns"]
		if returnType == nil {
			if node.Value.(string) == "main" {
				functionString += "int "
			} else {
				functionString += "void "
			}
		}

		functionString += node.Value.(string) + "("

		args := node.Metadata["args"].(*builder.Node).Value.([]*builder.Node)
		for _, arg := range args {
			fmt.Println("arg", arg)
		}

		functionString += ") {}"

		return functionString, nil
	}

	return "nothing", nil
}

func TranspileExpression(n *builder.Node) (*string, error) {
	switch n.Type {

	case "literal":
		return TranspileLiteralExpression(n)

	case "ident":
		return TranspileIdentExpression(n)
	}

	return nil, errors.New("Not implemented: " + n.Type)
}

func TranspileStatement(n *builder.Node) (*string, error) {
	switch n.Type {

	case "literal":
		return TranspileLiteralExpression(n)

	case "ident":
		return TranspileIdentExpression(n)
	}

	return nil, errors.New("Not implemented: " + n.Type)
}

func TranspileIdentExpression(n *builder.Node) (*string, error) {
	if n.Type != "ident" {
		return nil, errors.New("Node is not an ident")
	}

	nString, ok := n.Value.(string)
	if !ok {
		return nil, errors.New("Node value was not a string")
	}

	return &nString, nil
}

func TranspileType(n *builder.Node) (*string, error) {
	if n.Type != "type" {
		return nil, errors.New("Node is not an type")
	}

	nString, ok := n.Value.(string)
	if !ok {
		return nil, errors.New("Node value was not a string")
	}

	return &nString, nil
}

func TranspileLiteralExpression(n *builder.Node) (*string, error) {
	if n.Type != "literal" {
		return nil, errors.New("Node is not an literal")
	}

	nString := fmt.Sprintf("%v", n.Value)

	return &nString, nil
}

func TranspileArrayExpression(n *builder.Node) (*string, error) {
	if n.Type != "array" {
		return nil, errors.New("Node is not an array")
	}

	var (
		nString = "{ "
		vString *string
		err     error
	)

	value := n.Value.([]*builder.Node)
	for _, v := range value {
		vString, err = TranspileExpression(v)
		if err != nil {
			return nil, err
		}

		nString += *vString + ", "
	}

	// Cut off the last comma and space
	nString = nString[:len(nString)-2] + " }"

	return &nString, nil
}

func TranspileAssignmentStatement(n *builder.Node) (*string, error) {
	if n.Type != "assignment" {
		return nil, errors.New("Node is not an assignment")
	}

	var (
		nString = ""
		vString *string
		err     error
	)

	// Left should be ident
	// Right should be general expression
	// This will require some prepping atleast to figure out
	// if we need any pre-statements

	// Translate the ident expression (lhs)
	vString, err = TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	nString = *vString + " = "

	// Translate the ident expression (lhs)
	vString, err = TranspileExpression(n.Right)
	if err != nil {
		return nil, err
	}

	nString += *vString

	return &nString, nil
}

func TranspileDeclarationStatement(n *builder.Node) (*string, error) {
	if n.Type != "decl" {
		return nil, errors.New("Node is not an declaration")
	}

	var (
		nString = ""
		vString *string
		err     error
	)

	// Left should be ident
	// Right should be general expression
	// This will require some prepping atleast to figure out
	// if we need any pre-statements

	vString, err = TranspileType(n.Value.(*builder.Node))
	if err != nil {
		return nil, err
	}

	nString = *vString + " "

	// Translate the ident expression (lhs)
	vString, err = TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	nString += *vString + " = "

	// Translate the ident expression (lhs)
	vString, err = TranspileExpression(n.Right)
	if err != nil {
		return nil, err
	}

	nString += *vString

	return &nString, nil
}
