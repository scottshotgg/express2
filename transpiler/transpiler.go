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

	case "comp":
		return TranspileConditionExpression(n)

	case "binop":
		return TranspileBinOpExpression(n)

	case "array":
		return TranspileArrayExpression(n)

	case "call":
		return TranspileCallExpression(n)
	}

	return nil, errors.New("Not implemented: " + n.Type)
}

func TranspileStatement(n *builder.Node) (*string, error) {
	switch n.Type {

	case "literal":
		return TranspileLiteralExpression(n)

	case "ident":
		return TranspileIdentExpression(n)

	case "assignment":
		return TranspileAssignmentStatement(n)

	case "decl":
		return TranspileDeclarationStatement(n)
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

	nString += *vString + ";"

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

	nString += *vString + ";"

	return &nString, nil
}

func TranspileIncrementExpression(n *builder.Node) (*string, error) {
	if n.Type != "inc" {
		return nil, errors.New("Node is not an inc")
	}

	var (
		nString = ""
		vString *string
		err     error
	)

	// Translate the ident expression (lhs)
	vString, err = TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	nString += *vString + "++"

	return &nString, nil
}

func TranspileConditionExpression(n *builder.Node) (*string, error) {
	if n.Type != "comp" {
		return nil, errors.New("Node is not an comp")
	}

	var (
		nString = ""
		vString *string
		err     error
	)

	// Translate the lhs
	vString, err = TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	nString += *vString

	// Translate the rhs
	vString, err = TranspileExpression(n.Right)
	if err != nil {
		return nil, err
	}

	nString += n.Value.(string) + *vString

	return &nString, nil
}

func TranspileBinOpExpression(n *builder.Node) (*string, error) {
	if n.Type != "binop" {
		return nil, errors.New("Node is not a binop")
	}

	var (
		nString = ""
		vString *string
		err     error
	)

	// Translate the ident expression (lhs)
	vString, err = TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	nString += *vString + n.Value.(string)

	// Translate the ident expression (lhs)
	vString, err = TranspileExpression(n.Right)
	if err != nil {
		return nil, err
	}

	nString += *vString

	return &nString, nil
}

func TranspileBlockStatement(n *builder.Node) (*string, error) {
	if n.Type != "block" {
		return nil, errors.New("Node is not a block")
	}

	var (
		nString = ""
		vString *string
		err     error
	)

	for _, stmt := range n.Value.([]*builder.Node) {
		vString, err = TranspileStatement(stmt)
		if err != nil {
			return nil, err
		}

		nString += *vString
	}

	nString = "{" + nString + "}"

	return &nString, nil
}

func TranspileEGroup(n *builder.Node) (*string, error) {
	if n.Type != "egroup" {
		return nil, errors.New("Node is not a egroup")
	}

	var (
		nString = "("
		vString *string
		err     error
	)

	for _, e := range n.Value.([]*builder.Node) {
		vString, err = TranspileExpression(e)
		if err != nil {
			return nil, err
		}

		nString += *vString + ","
	}

	nString = nString[:len(nString)-1] + ")"

	return &nString, nil
}

func TranspileCallExpression(n *builder.Node) (*string, error) {
	if n.Type != "call" {
		return nil, errors.New("Node is not a call")
	}

	var (
		nString = ""
		vString *string
		err     error
	)

	vString, err = TranspileIdentExpression(n.Value.(*builder.Node))
	if err != nil {
		return nil, err
	}

	nString += *vString

	vString, err = TranspileEGroup(n.Metadata["args"].(*builder.Node))
	if err != nil {
		return nil, err
	}

	nString += *vString

	return &nString, nil
}

func TranspileForInStatement(n *builder.Node) (*string, error) {
	// Change forin to be a block statement containing:
	//	- declare var
	//	- declare array/iter
	//	- while var < iter.length
	//	- loop_block
	//	-	increment var

	if n.Type != "forin" {
		return nil, errors.New("Node is not a forin")
	}

	var (
		nString = "{"
		vString *string
		err     error
	)

	// Make and translate the ident into a declaration
	ds := TransformIdentToDefaultDeclaration(n.Metadata["start"].(*builder.Node))
	vString, err = TranspileDeclarationStatement(ds)
	if err != nil {
		return nil, err
	}

	nString += *vString

	// Make and translate the array expression into a declaration
	dss := TransformExpressionToDeclaration(n.Metadata["end"].(*builder.Node))
	vString, err = TranspileDeclarationStatement(dss)
	if err != nil {
		return nil, err
	}

	nString += *vString

	// Make and translate the less than operation
	vString, err = TranspileConditionExpression(&builder.Node{
		Type:  "comp",
		Value: "<",
		Left:  n.Metadata["start"].(*builder.Node),
		Right: GenerateLengthCall(dss),
	})
	if err != nil {
		return nil, err
	}

	nString += fmt.Sprintf("while(%s)", *vString)

	// Translate the block statement
	vString, err = TranspileBlockStatement(n.Value.(*builder.Node))
	if err != nil {
		return nil, err
	}

	nString += *vString

	// Lastly, make and translate an increment statement for the ident
	vString, err = TranspileIncrementExpression(&builder.Node{
		Type: "inc",
		Left: n.Metadata["start"].(*builder.Node),
	})
	if err != nil {
		return nil, err
	}

	nString += *vString + "}"

	return &nString, nil
}

func TranspileForOfStatement(n *builder.Node) (*string, error) {
	// Change forin to be a block statement containing:
	//	- declare temp var
	//	- declare array/iter
	//	- while tempvar < iter.length
	//	- var = iter[tempvar]
	//	- loop_block
	//	-	increment var

	if n.Type != "forin" {
		return nil, errors.New("Node is not a forin")
	}

	var (
		nString = "{"
		vString *string
		err     error
	)

	// Make and translate the ident into a declaration
	ds := TransformIdentToDefaultDeclaration(n.Metadata["start"].(*builder.Node))
	vString, err = TranspileDeclarationStatement(ds)
	if err != nil {
		return nil, err
	}

	nString += *vString

	// Make and translate the array expression into a declaration
	dss := TransformExpressionToDeclaration(n.Metadata["end"].(*builder.Node))
	vString, err = TranspileDeclarationStatement(dss)
	if err != nil {
		return nil, err
	}

	nString += *vString

	// Make and translate the less than operation
	vString, err = TranspileConditionExpression(&builder.Node{
		Type:  "comp",
		Value: "<",
		Left:  n.Metadata["start"].(*builder.Node),
		Right: GenerateLengthCall(dss),
	})
	if err != nil {
		return nil, err
	}

	nString += fmt.Sprintf("while(%s)", *vString)

	// Translate the block statement
	vString, err = TranspileBlockStatement(n.Value.(*builder.Node))
	if err != nil {
		return nil, err
	}

	nString += *vString

	// Lastly, make and translate an increment statement for the ident
	vString, err = TranspileIncrementExpression(&builder.Node{
		Type: "inc",
		Left: n.Metadata["start"].(*builder.Node),
	})
	if err != nil {
		return nil, err
	}

	nString += *vString + "}"

	return &nString, nil
}

// func ConstructGreaterThanOperation(op string, lhs, rhs *builder.Node)

func GenerateLengthCall(n *builder.Node) *builder.Node {
	// TODO: need to have a switch here for arrays and stuff
	return &builder.Node{
		Type: "call",
		Value: &builder.Node{
			Type:  "ident",
			Value: "std::size",
		},
		Metadata: map[string]interface{}{
			"args": &builder.Node{
				Type:  "egroup",
				Value: []*builder.Node{n.Left},
			},
		},
	}
}

func TransformExpressionToDeclaration(n *builder.Node) *builder.Node {
	fmt.Println("n", n)

	// TODO: Type checker would give type here; use auto for now
	return &builder.Node{
		Type: "decl",
		Value: &builder.Node{
			Type:  "type",
			Value: "auto",
		},
		Left: &builder.Node{
			Type:  "ident",
			Value: "SOMETHING",
		},
		Right: n,
	}
}

func TransformIdentToDefaultDeclaration(n *builder.Node) *builder.Node {
	return &builder.Node{
		Type: "decl",
		Value: &builder.Node{
			Type:  "type",
			Value: "int",
		},
		Left: n,
		Right: &builder.Node{
			Type:  "literal",
			Value: 0,
		},
	}
}
