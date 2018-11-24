package builder

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/scottshotgg/express-token"
)

type Node struct {
	Type     string
	Kind     string
	Value    interface{}
	Metadata map[string]interface{}
	Left     *Node
	Right    *Node
}

type Builder struct {
	Tokens []token.Token
	Index  int
	// [op_tier][op] -> func
	OpFuncMap map[int]map[string]func(n *Node) (*Node, error)
}

type Index struct {
	Type  string
	Value interface{}
}

var (
	ErrNotImplemented = errors.New("Not implemented")
)

func New(tokens []token.Token) *Builder {
	b := Builder{
		Tokens: tokens,
	}

	b.OpFuncMap = map[int]map[string]func(n *Node) (*Node, error){
		1: map[string]func(n *Node) (*Node, error){
			token.Increment: b.ParseIncrement,
			token.LThan:     b.ParseConditionExpression,
			token.LBracket:  b.ParseIndexExpression,
			token.LParen:    b.ParseCall,
			token.Accessor:  b.ParseSelection,
		},
	}

	return &b
}

func (b *Builder) GetNextToken() (*token.Token, error) {
	if b.Index > len(b.Tokens)-1 {
		return nil, errors.New("Out of tokens")
	}

	return &b.Tokens[b.Index], nil
}

func (b *Builder) ParseGroupOfExpressions() (*Node, error) {
	// Check ourselves
	if b.Tokens[b.Index].Type != token.LParen {
		return nil, errors.New("Could not get group of expressions")
	}

	// Skip over the left paren token
	b.Index++

	exprs := []*Node{}

	for b.Tokens[b.Index].Type != token.RParen {
		expr, err := b.ParseExpression()
		if err != nil {
			return nil, err
		}

		b.Index++

		exprs = append(exprs, expr)

		// Check and skip over the separator
		if b.Tokens[b.Index].Type == token.Separator {
			b.Index++
		}
	}

	// Step over the right brace token
	b.Index++

	return &Node{
		Type:  "egroup",
		Value: exprs,
	}, nil
}

func (b *Builder) ParseCall(n *Node) (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.LParen {
		return nil, errors.New("Could not get lparen in function call;")
	}

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

func (b *Builder) ParseBlockStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.LBrace {
		return nil, errors.Errorf("Could not get lbrace; %s", b.Tokens[b.Index].Type)
	}

	// Increment over the left brace token
	b.Index++

	stmts := []*Node{}

	for b.Tokens[b.Index].Type != token.RBrace {
		stmt, err := b.ParseStatement()
		if err != nil {
			return nil, err
		}

		stmts = append(stmts, stmt)
	}

	// Step over the right brace token
	b.Index++

	return &Node{
		Type:  "block",
		Value: stmts,
	}, nil
}

func (b *Builder) ParseReturnStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Return {
		return nil, errors.Errorf("Could not get return; %+v", b.Tokens[b.Index])
	}

	// Skip over the `return` token
	b.Index++

	// If there is a newline, the return is void typed
	if b.Index < len(b.Tokens) && b.Tokens[b.Index].Value.String == "\n" {
		return &Node{
			Type: "return",
		}, nil
	}

	// we are only supporting one return value for now
	expr, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Step over the expression token
	b.Index++

	return &Node{
		Type: "return",
		Left: expr,
	}, nil
}

func (b *Builder) ParseIdentStatement() (*Node, error) {
	return nil, nil
}

// ParseStatement ** does ** not look ahead
func (b *Builder) ParseStatement() (*Node, error) {
	switch b.Tokens[b.Index].Type {
	case token.TypeDef:
		return b.ParseTypeDeclarationStatement()

	case token.Type:
		fmt.Println("**type**")
		// get declaration statement
		// need to check if there is an asterisk after
		// need to check if there is brackets after
		// MAKE a ParseType function that will return the type, check against user types, pointer, array, etc

		// For now: no user defined types ...
		return b.ParseDeclarationStatement()

	// TODO: what if types were expressions ...
	case token.Ident:
		return b.ParseIdentStatement()
		// check if user defined type - ParseType
		// get assignment statement

		// just let ParseExpression/Term/Factor handle these ...
		// or call
		// selection operation
		// index operation

		// TODO: Make a function that will determine what kind of `IdentStatement`
		//

	case token.Function:
		fmt.Println("function**")
		// get a function statement
		return b.ParseFunctionStatement()

	case token.Block:
		// get a block
		return b.ParseBlockStatement()

	case "*":
		// get a deref assignment
		return nil, ErrNotImplemented

	case token.If:
		// get an if-else
		fmt.Println("if**")
		return b.ParseIfStatement()

	case token.For:
		// get a loop
		return b.ParseForStatement()

	case token.Return:
		// get a return statement
		return b.ParseReturnStatement()
	}

	// defer processing to next level higher
	return nil, nil
}

func (b *Builder) ParseForPrepositionStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.For {
		return nil, errors.New("Could not get for in")
	}

	// Step over the for token
	b.Index++

	// Parse the ident before the `in` token
	expr1, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Step over the ident token
	b.Index++

	var prepType string

	switch b.Tokens[b.Index].Value.String {
	case "in":
		prepType = "forin"

	case "of":
		prepType = "forof"

	default:
		return nil, errors.New("Could not get preposition; " + b.Tokens[b.Index].Value.String)
	}

	// Step over the preposition
	b.Index++

	// Parse the array/expression
	expr2, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	body, err := b.ParseBlockStatement()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:  prepType,
		Value: body,
		Metadata: map[string]interface{}{
			"start": expr1,
			"end":   expr2,
		},
	}, nil
}

func (b *Builder) ParseForStdStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.For {
		return nil, errors.Errorf("Could not get for std; %+v", b.Tokens[b.Index])
	}

	// Step over the for token
	b.Index++

	// Parse the declaration or assignment statement
	stmt, err := b.ParseStatement()
	if err != nil {
		return nil, err
	}

	node := Node{
		Type: "forstd",
		Metadata: map[string]interface{}{
			"start": stmt,
		},
	}

	// Check and skip over the separator
	if b.Tokens[b.Index].Type == token.Separator {
		b.Index++
	}

	// Parse the bounding conditional (expression)
	// Might want to make specific functions like `ParseConditional`
	// if we know we need it
	node.Metadata["end"], err = b.ParseExpression()
	if err != nil {
		return nil, err
	}

	b.Index++

	// Check and skip over the separator
	if b.Tokens[b.Index].Type == token.Separator {
		b.Index++
	}

	// Parse the increment
	node.Metadata["step"], err = b.ParseExpression()
	if err != nil {
		return nil, err
	}

	b.Index++

	// Check and skip over the separator
	if b.Tokens[b.Index].Type == token.Separator {
		b.Index++
	}

	node.Value, err = b.ParseBlockStatement()
	if err != nil {
		return nil, err
	}

	return &node, nil
}

func (b *Builder) ParseForStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.For {
		return nil, errors.New("Could not get for std")
	}

	if b.Index > len(b.Tokens)-2 {
		return nil, errors.New("Out of tokens")
	}

	// Clone the builder to backtrack
	// jsonClone, _ := json.Marshal(b)
	// // Click back to the last save
	// bp := &builder.Builder{}
	// json.Unmarshal(jsonClone, bp)

	// For right now just look ahead two
	if b.Tokens[b.Index+2].Type == token.Keyword {
		return b.ParseForPrepositionStatement()
	}

	return b.ParseForStdStatement()
}

func (b *Builder) ParseIfStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.If {
		return nil, errors.New("Could not get if")
	}

	// Step over the if token
	b.Index++

	// if EXPR BLOCK
	condition, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Step over the expression
	// TODO: this will have to move inside the expression I think
	b.Index++

	block, err := b.ParseBlockStatement()
	if err != nil {
		return nil, err
	}

	var elseBlock *Node

	// Check for an else block
	// nt, err := b.GetNextToken()
	// if err != nil {
	// 	// TODO: need to return a better error here
	// 	return nil, err
	// }

	if b.Index < len(b.Tokens)-1 && b.Tokens[b.Index].Type == token.Else {
		// Step over the else token
		b.Index++

		// Check for an else if
		if b.Tokens[b.Index].Type == token.If {
			elseBlock, err = b.ParseIfStatement()
			if err != nil {
				return nil, err
			}
		} else {
			elseBlock, err = b.ParseBlockStatement()
			if err != nil {
				return nil, err
			}
		}
	}

	return &Node{
		Type:  "if",
		Value: condition,
		Left:  block,
		Right: elseBlock,
	}, nil
}

func (b *Builder) ParseArrayType(typeOf string) (*Node, error) {
	var dim []*Index

	// Look ahead at the next token here
	for b.Index < len(b.Tokens)-1 && b.Tokens[b.Index+1].Type == token.LBracket {
		// Increment over the type token
		b.Index++

		expr, err := b.ParseArrayExpression()
		if err != nil {
			return nil, err
		}

		b.Index--

		if expr.Value == nil {
			return nil, errors.Errorf("Array parse value was nil; %+v", expr)
		}

		nodesAssert, ok := expr.Value.([]*Node)
		if !ok {
			return nil, errors.Errorf("Invalid assertion; %+v", expr.Value)
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
			return nil, errors.New("Cannot use multiple expression inside array type initializer")
		}

		dim = append(dim, &dimValue)
	}

	b.Index++

	return &Node{
		Type:  "type",
		Value: "array",
		Metadata: map[string]interface{}{
			// "type": typeOf,
			"dim": dim,
		},
	}, nil
}

func (b *Builder) ParseType() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Type {
		return nil, errors.New("Could not get type; " + b.Tokens[b.Index].Value.String)
	}

	// TODO: we would need to implement something like this
	// TODO: this is where we would also do pointers, need to do function types, etc
	// if typeOf == "map" {

	// }

	typeOf := b.Tokens[b.Index].Value.String

	if b.Index < len(b.Tokens)-1 && b.Tokens[b.Index+1].Type == token.LBracket {
		return b.ParseArrayType(typeOf)
	}

	// Increment over the type
	b.Index++

	return &Node{
		Type:  "type",
		Value: typeOf,
	}, nil
}

func (b *Builder) ParseIndexExpression(n *Node) (*Node, error) {
	if b.Index > len(b.Tokens)-1 {
		return nil, errors.Errorf("Out of tokens")
	}

	if b.Tokens[b.Index].Type != token.LBracket {
		return nil, errors.Errorf("Could not get LBracket")
	}

	b.Index++

	expr, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Step over the expression
	b.Index++

	return &Node{
		Type: "index",
		// Value: n,
		Left:  n,
		Right: expr,
	}, nil
}

func (b *Builder) ParseDeclarationStatement() (*Node, error) {
	typeOf, err := b.ParseType()
	if err != nil {
		return nil, err
	}

	// Check that the next token is an ident
	if b.Tokens[b.Index].Type != token.Ident {
		return nil, errors.New("Could not get declaration statement without ident")
	}

	// Create the ident
	ident, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Increment over the ident token
	b.Index++

	// Check for the equals token
	if b.Tokens[b.Index].Type != token.Assign {
		return &Node{
			Type:  "decl",
			Value: typeOf,
			Left:  ident,
		}, nil

		// return nil, errors.New("No equals found after ident")
	}

	// Increment over the equals
	b.Index++

	// Parse the right hand side
	expr, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Increment over the first part of the expression
	b.Index++

	return &Node{
		Type:  "decl",
		Value: typeOf,
		Left:  ident,
		Right: expr,
	}, nil
}

func (b *Builder) ParseTypeDeclarationStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.TypeDef {
		return nil, errors.New("Could not get type declaration statement")
	}

	// Skip over the `type` token
	b.Index++

	// Create the ident
	ident, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Increment over the ident token
	b.Index++

	// Check for the equals token
	if b.Tokens[b.Index].Type != token.Assign {
		return nil, errors.New("No equals found after ident in typedef")
	}

	// Increment over the equals
	b.Index++

	// Parse the right hand side
	typeOf, err := b.ParseType()
	if err != nil {
		return nil, errors.New("Could not get expression")
	}

	// Increment over the first part of the expression
	b.Index++

	return &Node{
		Type:  "typedef",
		Left:  ident,
		Right: typeOf,
	}, nil
}

func (b *Builder) ParseStructStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Struct {
		return nil, errors.New("Could not get struct declaration statement")
	}

	// Skip over the `struct` token
	b.Index++

	// Create the ident
	ident, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Increment over the ident token
	b.Index++

	// Check for the equals token
	if b.Tokens[b.Index].Type != token.Assign {
		return nil, errors.New("No equals found after ident in struct def")
	}

	// Increment over the equals
	b.Index++

	// Parse the right hand side
	body, err := b.ParseBlockStatement()
	if err != nil {
		return nil, errors.New("Could not get expression")
	}

	// Increment over the first part of the expression
	b.Index++

	return &Node{
		Type:  "struct",
		Left:  ident,
		Right: body,
	}, nil
}

func (b *Builder) ParseLetStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Let {
		return nil, errors.New("Could not get let statement")
	}

	// Skip over the let token
	b.Index++

	ident, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Increment over the ident token
	b.Index++

	// Check for the equals token
	if b.Tokens[b.Index].Type != token.Assign {
		// This is where we would implement variable declarations
		// without values, other types of assignment, etc
		// Leave it alone for now
		return nil, errors.New("No equals found after ident in let")
	}

	// Increment over the equals
	b.Index++

	// Parse the right hand side
	expr, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Increment over the first part of the expression
	b.Index++

	return &Node{
		Type:  "let",
		Left:  ident,
		Right: expr,
	}, nil
}

func (b *Builder) ParseAssignmentStatement() (*Node, error) {
	// into: [expr] = [expr]
	// Check that the next token is an ident
	if b.Tokens[b.Index].Type != token.Ident {
		return nil, errors.New("Could not get assignment statement without ident")
	}

	ident, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Increment over the ident token
	b.Index++

	// Check for the equals token
	if b.Tokens[b.Index].Type != token.Assign {
		// This is where we would implement variable declarations
		// without values, other types of assignment, etc
		// Leave it alone for now
		return nil, errors.New("No equals found after ident")
	}

	// Increment over the equals
	b.Index++

	// Parse the right hand side
	expr, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Increment over the first part of the expression
	b.Index++

	return &Node{
		Type:  "assignment",
		Left:  ident,
		Right: expr,
	}, nil
}

func (b *Builder) ParseImportStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Import {
		return nil, errors.New("Could not get import statement")
	}

	// Step over the import token
	b.Index++

	expr, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:  "import",
		Value: expr,
	}, nil
}

func (b *Builder) ParseIncludeStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Include {
		return nil, errors.New("Could not get include statement")
	}

	// Step over the import token
	b.Index++

	expr, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:  "include",
		Value: expr,
	}, nil
}

func (b *Builder) ParseExpression() (*Node, error) {
	// This is where we will implement secondary tier operators (+ , -)

	return b.ParseTerm()
}

func (b *Builder) ParseConditionExpression(expr *Node) (*Node, error) {
	// Step over the conditional operator token
	b.Index++

	right, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:  "comp",
		Value: "",
		Left:  expr,
		Right: right,
	}, nil
}

func (b *Builder) ParseIncrement(n *Node) (*Node, error) {
	return &Node{
		Type:  "inc",
		Value: n,
	}, nil
}

func (b *Builder) ParseTerm() (*Node, error) {
	// This is where we will implement primary tier operators (* , /)

	factor, err := b.ParseFactor()
	if err != nil {
		return nil, err
	}

	// var ok = true
	// var opFunc func(n *Node) (*Node, error)

	// LOOKAHEAD performed to figure out whether the expression is done
	for b.Index < len(b.Tokens)-1 {
		// Look for a tier1 operator in the func map
		opFunc, ok := b.OpFuncMap[1][b.Tokens[b.Index+1].Type]
		if !ok {
			break
		}

		// Step over the factor
		b.Index++

		factor, err = opFunc(factor)
		if err != nil {
			return nil, err
		}
	}

	return factor, nil
}

func (b *Builder) ParseFactor() (*Node, error) {
	// Here we will switch on the type and determine whether we have:
	// - literal
	// - ident
	// - call
	// - index operation
	// - selection operation
	// - block
	// - array
	// - nil

	switch b.Tokens[b.Index].Type {
	// Any literal value
	case token.Literal:
		return &Node{
			Type:  "literal",
			Kind:  b.Tokens[b.Index].Value.Type,
			Value: b.Tokens[b.Index].Value.True,
		}, nil

	// Variable identifier
	case token.Ident:
		return &Node{
			Type:  "ident",
			Value: b.Tokens[b.Index].Value.String,
		}, nil

	// Nested expression
	case token.LParen:
		// Skip over the left paren
		b.Index++

		expr, err := b.ParseExpression()
		if err != nil {
			return nil, err
		}

		// Skip over the expression
		b.Index++

		if b.Tokens[b.Index].Type != token.RParen {
			return nil, errors.Errorf("No rparen found at end of factor-expression: %+v", b.Tokens[b.Index])
		}

		// Skip over the right paren
		b.Index++

		return expr, nil

	// Array expression
	case token.LBracket:
		return b.ParseArrayExpression()

	// Named block
	case token.LBrace:
		return b.ParseBlockStatement()

	default:
		return nil, errors.Errorf("Could not parse expression from token; %+v", b.Tokens[b.Index])
	}
}

func (b *Builder) ParseSelection(n *Node) (*Node, error) {
	if b.Index > len(b.Tokens)-1 {
		return nil, errors.Errorf("Out of tokens")
	}

	if b.Tokens[b.Index].Type != token.Accessor {
		return nil, errors.Errorf("Could not get selection operator; %+v", b.Tokens[b.Index])
	}

	b.Index++

	expr, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type: "selection",
		// Value: n,
		Left:  n,
		Right: expr,
	}, nil
}

func (b *Builder) ParseArrayExpression() (*Node, error) {
	// Check ourselves
	if b.Tokens[b.Index].Type != token.LBracket {
		return nil, errors.New("Could not get array expression")
	}

	// Skip over the left bracket token
	b.Index++

	exprs := []*Node{}

	for b.Index < len(b.Tokens) && b.Tokens[b.Index].Type != token.RBracket {
		expr, err := b.ParseExpression()
		if err != nil {
			return nil, err
		}

		b.Index++

		exprs = append(exprs, expr)

		// Check and skip over the separator
		if b.Tokens[b.Index].Type == token.Separator {
			b.Index++
		}
	}

	// Step over the right bracket token
	b.Index++

	return &Node{
		Type:  "array",
		Value: exprs,
	}, nil
}

func (b *Builder) ParseGroupOfStatements() (*Node, error) {
	// Check ourselves
	if b.Tokens[b.Index].Type != token.LParen {
		return nil, errors.New("Could not get group of statements")
	}

	// Skip over the left paren token
	b.Index++

	stmts := []*Node{}

	for b.Tokens[b.Index].Type != token.RParen {
		stmt, err := b.ParseStatement()
		if err != nil {
			return nil, err
		}

		stmts = append(stmts, stmt)

		// Check and skip over the separator
		if b.Tokens[b.Index].Type == token.Separator {
			b.Index++
		}
	}

	// Step over the right brace token
	b.Index++

	return &Node{
		Type:  "sgroup",
		Value: stmts,
	}, nil
}

func (b *Builder) ParseFunctionStatement() (*Node, error) {
	// Check ourselves
	if b.Tokens[b.Index].Type != token.Function {
		return nil, errors.New("Could not get function")
	}

	// Step over the function token
	b.Index++

	node := &Node{
		Type:     "function",
		Metadata: map[string]interface{}{},
	}

	// Named function
	if b.Tokens[b.Index].Type != token.Ident {
		return nil, errors.New("Could not get ident after function token")
	}

	node.Value = b.Tokens[b.Index].Value.String
	// Step over the ident token
	b.Index++

	if b.Tokens[b.Index].Type != token.LParen {
		return nil, errors.New("Could not get lparen")
	}

	args, err := b.ParseGroupOfStatements()
	if err != nil {
		return nil, err
	}

	// Might want to avoid putting this here if we don't have any
	node.Metadata["args"] = args

	// We are not supporting multiple returns for now
	// // Check for multiple returns; another left paren
	// if b.Tokens[b.Index].Type == token.LParen {
	// 	return nil, errors.New("Could not get returns")
	// }

	// We are not supporting named arguments for now
	// Check for the return type token
	if b.Tokens[b.Index].Type == token.Type {
		node.Metadata["returns"] = &Node{
			Type:  "type",
			Value: b.Tokens[b.Index].Type,
		}

		// Step over the type token
		b.Index++
	}

	node.Value, err = b.ParseBlockStatement()
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (b *Builder) BuildAST() (Node, error) {
	for b.Index < len(b.Tokens)-1 {
		// switch b.Tokens[i].Type {
		// case token.Function:
		// 	i++
		// 	fmt.Println(b.ParseFunction())
		// }
		stmt, err := b.ParseStatement()
		if err != nil {
			return Node{}, err
		}

		if stmt == nil {
			return Node{}, errors.New("Could not get statement")
		}

		fmt.Println("stmt", stmt)
	}

	// fmt.Println(b.ParseBlock())

	return Node{
		Type: "program",
	}, nil
}
