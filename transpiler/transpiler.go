package transpiler

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

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

	wg.Add(1)

	go appendWorker(&wg)

	t.ASTCloneJSON, _ = json.Marshal(ast)

	return &t
}

var appendChan = make(chan string, 5)

func appendWorker(wg *sync.WaitGroup) {
	defer wg.Done()

	var totalFile string

	for a := range appendChan {
		totalFile += a
	}

	fmt.Println("totalFile", totalFile)
}

func emit(line string) {
	appendChan <- line
}

var wg sync.WaitGroup

/*
	Transpile needs to work like this:
	- recurse through each statement
	- if the statement contains ANY block, then flatten on the node
*/

func (t *Transpiler) Transpile() (string, error) {
	// Extract the nodes
	var (
		nodes   = t.AST.Value.([]*builder.Node)
		stringP *string
		err     error

		cpp string
	)

	for _, node := range nodes {
		fmt.Println("node", node)

		// Switch on the statement type to figure out how to process it
		switch node.Type {
		case "function":
			stringP, err = TranspileFunctionStatement(node)
			if err != nil {
				return cpp, err
			}

			cpp += *stringP
		}
	}

	return cpp, nil
}

// // just grab the first one for now
// var node = t.AST.Value.([]*builder.Node)[0]

// fmt.Println("ndoe", node)

// // // first flatten the nodes
// tree_flattener.Flatten(node.Value.(*builder.Node))
// fmt.Println("nodes", node.Value.(*builder.Node))

// for _, stmt := range node.Value.(*builder.Node).Value.([]*builder.Node) {
// 	fmt.Printf("stmt %+v\n", stmt)
// 	stmtSTringP, err := TranspileStatement(stmt)
// 	if err != nil {
// 		return "", err
// 	}

// 	fmt.Println(*stmtSTringP)
// }

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

	case "function":
		return TranspileFunctionStatement(n)

	case "return":
		return TranspileReturnStatement(n)

	case "forof":
		return TranspileForOfStatement(n)

	case "forin":
		return TranspileForOfStatement(n)

	case "forstd":
	}

	return nil, errors.New("Not implemented: " + n.Type)
}

func TranspileReturnStatement(n *builder.Node) (*string, error) {
	if n.Type != "return" {
		return nil, errors.New("Node is not a return statement")
	}

	// Return statments come in the form `return` { expr }

	var nString = "return"

	fmt.Printf("n: %+v\n", n)

	// LHS (the return expression) is allowed to be empty
	if n.Left != nil {
		exprString, err := TranspileExpression(n.Left)
		if err != nil {
			return nil, err
		}

		nString += " " + *exprString
	}

	nString += ";"

	return &nString, nil
}

func TranspileFunctionStatement(n *builder.Node) (*string, error) {
	if n.Type != "function" {
		return nil, errors.New("Node is not an function")
	}

	/*
		A map with keys for `returns` and `args` will be egroups in the Metadata
		`Kind` is the name of the function
		`Value` is the block than needs to be translated
	*/

	if n.Kind == "" {
		return nil, errors.New("Somehow we parsed a function without a name ...")
	}

	// Start out with just the name; we will put the return type later
	var nString = n.Kind

	// args is an `sgroup`
	argsString, err := TranspileSGroup(n.Metadata["args"].(*builder.Node))
	if err != nil {
		return nil, err
	}

	// Append the args
	nString += *argsString

	var (
		returns = n.Metadata["returns"]

		returnsString = "void"

		// Start returns off as void
		returnsStringP = &returnsString
	)

	if returns != nil {
		// returns is a `type` for now; multiple returns are not supported right now
		returnsStringP, err = TranspileType(returns.(*builder.Node))
		if err != nil {
			return nil, err
		}
	}

	// Prepend the return string with a space
	nString = *returnsStringP + " " + nString

	blockString, err := TranspileBlockStatement(n.Value.(*builder.Node))
	if err != nil {
		return nil, err
	}

	nString += *blockString

	return &nString, nil
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

	// LHS is not allowed to be nil
	if n.Left == nil {
		return nil, errors.New("nil Left hand side")
	}

	// Translate the ident expression (lhs)
	vString, err = TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	nString += *vString

	// RHS is allowed to be nil to support declarations without values like `string s`
	if n.Right == nil {
		return &nString, nil
	}

	// Translate the ident expression (lhs)
	vString, err = TranspileExpression(n.Right)
	if err != nil {
		return nil, err
	}

	nString += " = " + *vString + ";"

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

	nString += *vString + "++;"

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
		nString string
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

	if len(nString) > 0 {
		nString = nString[:len(nString)-1]
	}

	nString = "(" + nString + ")"

	return &nString, nil
}

func TranspileSGroup(n *builder.Node) (*string, error) {
	if n.Type != "sgroup" {
		return nil, errors.New("Node is not a sgroup")
	}

	var (
		nString string
		vString *string
		err     error
	)

	for _, s := range n.Value.([]*builder.Node) {
		vString, err = TranspileStatement(s)
		if err != nil {
			return nil, err
		}

		nString += *vString + ","
	}

	nString = "(" + nString + ")"

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

	return nil, errors.New("not implemented: forin")

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

	return nil, errors.New("not implemented: forof")

	if n.Type != "forof" {
		return nil, errors.New("Node is not a forof")
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
