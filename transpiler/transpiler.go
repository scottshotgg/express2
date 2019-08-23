package transpiler

import (
	"fmt"

	"strings"

	"strconv"

	"github.com/pkg/errors"
	"github.com/scottshotgg/express2/builder"
)

// Transpiler should only deal in terms of statements
// everything else that you need should be behind the scenes
type Transpiler interface {
	// This has a bunch of function
	Transpile() (*string, error)
	Statement(n *builder.Node) (*string, error)
	Expression(n *builder.Node) (*string, error)

	// Block(n *builder.Node) (*string, error) // This might have to have both a Statement and Expression
	// Function(n *builder.Node) (*string, error)
	// Return(n *builder.Node) (*string, error)
	// Call(n *builder.Node) (*string, error)
	// IfElse(n *builder.Node) (*string, error)
	// TypeDef(n *builder.Node) (*string, error)
	// Struct(n *builder.Node) (*string, error)
	// Import(n *builder.Node) (*string, error)
	// Include(n *builder.Node) (*string, error)
	// Package(n *builder.Node) (*string, error)

	// Increment(n *builder.Node) (*string, error)
	// Decrement(n *builder.Node) (*string, error)
	// Index(n *builder.Node) (*string, error)
	// Select(n *builder.Node) (*string, error)
	// Ref(n *builder.Node) (*string, error)
	// Deref(n *builder.Node) (*string, error)

	// Map(n *builder.Node) (*string, error)
	// Object(n *builder.Node) (*string, error)
	// Enum(n *builder.Node) (*string, error)
	// Conditional(n *builder.Node) (*string, error)
	// BinOp(n *builder.Node) (*string, error)

	// Defer(n *builder.Node) (*string, error)

	// Ident(n *builder.Node) (*string, error)
	// Literal(n *builder.Node) (*string, error)

	// Declaration(n *builder.Node) (*string, error)
	// KeyValue(n *builder.Node) (*string, error)
	// Launch(n *builder.Node) (*string, error)
	// Assignment(n *builder.Node) (*string, error)
	// Array(n *builder.Node) (*string, error) // This might have to have both a Statement and Expression
	// EGroup(n *builder.Node) (*string, error)
	// SGroup(n *builder.Node) (*string, error)
	// ForStd(n *builder.Node) (*string, error)
	// ForIn(n *builder.Node) (*string, error)
	// ForOf(n *builder.Node) (*string, error)
	// While(n *builder.Node) (*string, error)

	// Utility functions need to be added; things like GenerateLengthCall from node
}

// Make a C implementation; JUST C, not C++
// Make an LLVM IR implementation

type C99 struct {
	ast *builder.Node
}

func NewC99(ast *builder.Node) *C99 {
	return &C99{
		ast: ast,
	}
}

func (t *C99) Transpile() (*string, error) {
	// try make later
	var stmts []string

	for _, n := range t.ast.Value.([]*builder.Node) {
		var stmt, err = t.Statement(n)
		if err != nil {
			return nil, err
		}

		stmts = append(stmts, *stmt)
	}

	var code = strings.Join(stmts, "\n")
	return &code, nil
}

func (t *C99) Statement(n *builder.Node) (*string, error) {
	switch n.Type {
	case "function":
		return t.function(n)

	case "decl":
		return t.declaration(n)

		// case "assignment":
		// 	return t.assign(n)
	}

	return nil, errors.Errorf("could not transpile statement: %+v", *n)
}

func (t *C99) Expression(n *builder.Node) (*string, error) {
	/*
		An expression could be many things:
		- literal
		- ident
		- function call
		- operation
	*/

	var code string

	switch n.Type {
	case "ident":
		code = n.Value.(string)

	case "literal":
		switch n.Kind {
		case "int":
			code = strconv.Itoa(n.Value.(int))

		case "bool":
			code = strconv.FormatBool(n.Value.(bool))

		case "float":
			code = strconv.FormatFloat(n.Value.(float64), 'f', -1, 64)

		case "char":
			fallthrough
		case "string":
			code = n.Value.(string)

		}

	default:
		pn(n)
		return nil, errors.Errorf("could not deduce literal from: %+v", *n)
	}

	return &code, nil
}

func (t *C99) declaration(n *builder.Node) (*string, error) {
	// [type] [assignment]

	// The value is the type of the declaration
	var typeOf *string
	var assignment *string
	var err error

	typeOf, err = t.typeOfDecl(n)
	if err != nil {
		return nil, err
	}

	assignment, err = t.assignment(n)
	if err != nil {
		return nil, err
	}

	var code = *typeOf + " " + *assignment

	return &code, nil

}

func (t *C99) function(n *builder.Node) (*string, error) {
	// func [ident](args) (returns) (body)

	// Get the return type kind -> support multi args later

	// body
	body, err := t.sblock(n.Value.(*builder.Node))
	if err != nil {
		return nil, err
	}

	// args
	args, err := t.sgroup(n.Metadata["args"].(*builder.Node))
	if err != nil {
		return nil, err
	}

	// returns; for now the first one will be taken as the return
	returns, err := t.egroup(n.Metadata["returns"].(*builder.Node))
	if err != nil {
		return nil, err
	}

	var code = fmt.Sprintf("%s %s(%s) %s", *returns, n.Kind, *args, *body)

	return &code, nil
}

func (t *C99) assignment(n *builder.Node) (*string, error) {
	// [expression] = [expression]
	var expr, err = t.Expression(n.Left)
	if err != nil {
		return nil, err
	}

	var code = *expr

	// This is for the case:
	//	int i
	if n.Right == nil {
		return &code, nil
	}

	expr, err = t.Expression(n.Right)
	if err != nil {
		return nil, err
	}

	code += " = " + *expr

	return &code, nil
}

func (t *C99) typeOfDecl(n *builder.Node) (*string, error) {
	var typeOf string

	switch n.Value.(*builder.Node).Type {
	case "type":
		typeOf = n.Value.(*builder.Node).Value.(string)
		return &typeOf, nil
		// case "selection"
	}

	return nil, errors.Errorf("unknown type for type %+v", n.Value.(*builder.Node).Type)
}

// These will give back blank strings if they are nil
func (t *C99) sgroup(n *builder.Node) (*string, error) {
	var args []string

	for _, node := range n.Value.([]*builder.Node) {
		arg, err := t.declaration(node)
		if err != nil {
			return nil, err
		}

		args = append(args, *arg)
	}

	var code = strings.Join(args, ", ")

	return &code, nil
}

func (t *C99) egroup(n *builder.Node) (*string, error) {
	// Only get the first return now since we are not worried about multiple returns rn
	var code = n.Value.([]*builder.Node)[0].Value.(string)

	return &code, nil
}

func (t *C99) sblock(n *builder.Node) (*string, error) {
	var stmts = []string{"{"}

	for _, node := range n.Value.([]*builder.Node) {
		stmt, err := t.Statement(node)
		if err != nil {
			return nil, err
		}

		stmts = append(stmts, *stmt)
	}

	var code = strings.Join(append(stmts, "}"), "\n")

	return &code, nil
}

// func (t *C99) eblock(n *builder.Node) (*string, error) {
// 	for _, node := range {

// 	}

// 	return nil, nil
// }

func pn(n *builder.Node) {
	fmt.Printf("node %+v\n", *n)
}
