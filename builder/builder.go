package builder

import (
	"github.com/pkg/errors"

	"github.com/scottshotgg/express-token"
)

type (
	opCallbackFn func(n *Node) (*Node, error)

	Node struct {
		Type     string
		Kind     string
		Value    interface{}
		Metadata map[string]interface{}
		Left     *Node
		Right    *Node
	}

	Builder struct {
		Tokens []token.Token
		Index  int
		// [op_tier][op] -> func
		OpFuncMap []map[string]opCallbackFn
	}

	Index struct {
		Type  string
		Value interface{}
	}
)

var (
	ErrNotImplemented = errors.New("Not implemented")
	ErrMultDimArrInit = errors.New("Cannot use multiple expression inside array type initializer")
	ErrOutOfTokens    = errors.New("Out of tokens")
)

func New(tokens []token.Token) *Builder {
	b := Builder{
		Tokens: tokens,
	}

	b.OpFuncMap = []map[string]opCallbackFn{
		0: map[string]opCallbackFn{
			token.Increment: b.ParseIncrement,
			token.Accessor:  b.ParseSelection,
			token.LBracket:  b.ParseIndexExpression,
			token.LParen:    b.ParseCall,
			token.LThan:     b.ParseConditionExpression,
			// token.PriOp: b.ParsePriOp,
		},

		1: map[string]opCallbackFn{
			// token.SecOp: b.ParseSecOp,
		},
	}

	return &b
}

func (b *Builder) GetNextToken() (*token.Token, error) {
	if b.Index > len(b.Tokens)-1 {
		return nil, ErrOutOfTokens
	}

	return &b.Tokens[b.Index], nil
}

func (b *Builder) ParseGroupOfExpressions() (*Node, error) {
	// Check ourselves
	if b.Tokens[b.Index].Type != token.LParen {
		return b.AppendTokenToError("Could not get group of expressions")
	}

	// Skip over the left paren token
	b.Index++

	var (
		expr  *Node
		exprs []*Node
		err   error
	)

	for b.Tokens[b.Index].Type != token.RParen {
		expr, err = b.ParseExpression()
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
		return b.AppendTokenToError("Could not get lparen in function call")
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
		return b.AppendTokenToError("Could not get left brace")
	}

	// Increment over the left brace token
	b.Index++

	var (
		stmt  *Node
		stmts []*Node
		err   error
	)

	for b.Tokens[b.Index].Type != token.RBrace {
		stmt, err = b.ParseStatement()
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
		return b.AppendTokenToError("Could not get return")
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

func (b *Builder) ParseDeref() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.PriOp &&
		b.Tokens[b.Index].Value.String == "*" {
		return b.AppendTokenToError("Could not get deref")
	}

	// Look ahead and make sure it is an ident;you can't deref just anything...
	if b.Tokens[b.Index+1].Type != token.Ident {
		return b.AppendTokenToError("Could not get ident to deref")
	}

	// Step over the deref
	b.Index++

	ident, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type: "deref",
		Left: ident,
	}, nil
}

// TODO: what if types were expressions ...

// ParseStatement ** does ** not look ahead
func (b *Builder) ParseStatement() (*Node, error) {
	switch b.Tokens[b.Index].Type {
	case token.PriOp:
		return b.ParseDeref()

	case token.Package:
		return b.ParsePackageStatement()

	case token.Import:
		return b.ParseImportStatement()

	case token.Include:
		return b.ParseIncludeStatement()

	case token.TypeDef:
		return b.ParseTypeDeclarationStatement()

	case token.Type:
		return b.ParseDeclarationStatement()

	case token.Ident:
		return b.ParseAssignmentStatement()

	case token.Function:
		return b.ParseFunctionStatement()

	case token.LBrace:
		return b.ParseBlockStatement()

	case token.Struct:
		return b.ParseStructStatement()

	case token.Let:
		return b.ParseLetStatement()

	case token.If:
		return b.ParseIfStatement()

	case token.For:
		return b.ParseForStatement()

	case token.Return:
		return b.ParseReturnStatement()
	}

	return b.AppendTokenToError("Could not get statement from")
}

func (b *Builder) ParseForPrepositionStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.For {
		return b.AppendTokenToError("Could not get for in")
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
		return b.AppendTokenToError("Could not get preposition")
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
		return b.AppendTokenToError("Could not get for std")
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

func (b *Builder) AppendTokenToError(errText string) (*Node, error) {
	if b.Index < len(b.Tokens)-1 {
		return nil, errors.Errorf(errText+"; %+v", b.Tokens[b.Index])
	}

	return nil, errors.New(errText)
}

func (b *Builder) ParseForStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.For {
		return b.AppendTokenToError("Could not get for std")
	}

	if b.Index > len(b.Tokens)-2 {
		return nil, ErrOutOfTokens
	}

	// For right now just look ahead two
	if b.Tokens[b.Index+2].Type == token.Keyword {
		return b.ParseForPrepositionStatement()
	}

	return b.ParseForStdStatement()
}

func (b *Builder) ParseIfStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.If {
		return b.AppendTokenToError("Could not get if")
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

	n := Node{
		Type:  "if",
		Value: condition,
	}

	n.Left, err = b.ParseBlockStatement()
	if err != nil {
		return nil, err
	}

	if b.Index < len(b.Tokens)-1 && b.Tokens[b.Index].Type == token.Else {
		// Step over the else token
		b.Index++

		// Check for an else if
		if b.Tokens[b.Index].Type == token.If {
			n.Right, err = b.ParseIfStatement()
			if err != nil {
				return nil, err
			}
		} else {
			n.Right, err = b.ParseBlockStatement()
			if err != nil {
				return nil, err
			}
		}
	}

	return &n, nil
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
		return b.AppendTokenToError("Could not get type")
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
		return nil, ErrOutOfTokens
	}

	if b.Tokens[b.Index].Type != token.LBracket {
		return b.AppendTokenToError("Could not get left bracket")
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
		return b.AppendTokenToError("Could not get declaration statement")
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
		return b.AppendTokenToError("Could not get type declaration statement")
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
		return b.AppendTokenToError("No equals found after ident in typedef")
	}

	// Increment over the equals
	b.Index++

	// Parse the right hand side
	typeOf, err := b.ParseType()
	if err != nil {
		return nil, err
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
		return b.AppendTokenToError("Could not get struct declaration statement")
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
		return b.AppendTokenToError("No equals found after ident in struct def")
	}

	// Increment over the equals
	b.Index++

	// Parse the right hand side
	body, err := b.ParseBlockStatement()
	if err != nil {
		return nil, err
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
		return b.AppendTokenToError("Could not get let statement")
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
		return b.AppendTokenToError("No equals found after ident in let")
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
		return b.AppendTokenToError("Could not get assignment statement without ident")
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
		return b.AppendTokenToError("No equals found after ident")
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

func (b *Builder) ParsePackageStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Package {
		return b.AppendTokenToError("Could not get package statement")
	}

	// Step over the package token
	b.Index++

	expr, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type: "package",
		Left: expr,
	}, nil
}

func (b *Builder) ParseImportStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Import {
		return b.AppendTokenToError("Could not get import statement")
	}

	// Step over the import token
	b.Index++

	expr, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type: "import",
		Left: expr,
	}, nil
}

func (b *Builder) ParseIncludeStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Include {
		return b.AppendTokenToError("Could not get include statement")
	}

	// Step over the import token
	b.Index++

	expr, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type: "include",
		Left: expr,
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
		Type: "comp",
		// Value: "",
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
		opFunc, ok := b.OpFuncMap[0][b.Tokens[b.Index+1].Type]
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

	// Deref operator
	case token.PriOp:
		return b.ParseDeref()

	// Nested expression
	case token.LParen:
		return b.ParseNestedExpression()

	// Array expression
	case token.LBracket:
		return b.ParseArrayExpression()

	// Named block
	case token.LBrace:
		return b.ParseBlockStatement()
	}

	return b.AppendTokenToError("Could not parse expression from token")
}

func (b *Builder) ParseNestedExpression() (*Node, error) {
	// Check ourselves
	if b.Tokens[b.Index].Type != token.LParen {
		return b.AppendTokenToError("Could not get nested expression")
	}

	// Skip over the left paren
	b.Index++

	expr, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Skip over the expression
	b.Index++

	if b.Tokens[b.Index].Type != token.RParen {
		return b.AppendTokenToError("No rparen found at end of factor-expression")
	}

	// Skip over the right paren
	b.Index++

	return expr, nil
}

func (b *Builder) ParseSelection(n *Node) (*Node, error) {
	if b.Index > len(b.Tokens)-1 {
		return nil, ErrOutOfTokens
	}

	if b.Tokens[b.Index].Type != token.Accessor {
		return b.AppendTokenToError("Could not get selection operator")
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
		return b.AppendTokenToError("Could not get array expression")
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
		return b.AppendTokenToError("Could not get group of statements")
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
		return b.AppendTokenToError("Could not get function")
	}

	// Step over the function token
	b.Index++

	node := &Node{
		Type:     "function",
		Metadata: map[string]interface{}{},
	}

	// Named function
	if b.Tokens[b.Index].Type != token.Ident {
		return b.AppendTokenToError("Could not get ident after function token")
	}

	node.Value = b.Tokens[b.Index].Value.String
	// Step over the ident token
	b.Index++

	if b.Tokens[b.Index].Type != token.LParen {
		return b.AppendTokenToError("Could not get left paren")
	}

	args, err := b.ParseGroupOfStatements()
	if err != nil {
		return nil, err
	}

	// Might want to avoid putting this here if we don't have any
	node.Metadata["args"] = args

	// We are not supporting multiple returns for now
	// // Check for multiple returns;another left paren
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

func (b *Builder) BuildAST() (*Node, error) {
	for b.Index < len(b.Tokens)-1 {
		// switch b.Tokens[i].Type {
		// case token.Function:
		// 	i++
		// 	fmt.Println(b.ParseFunction())
		// }
		stmt, err := b.ParseStatement()
		if err != nil {
			return nil, err
		}

		// Just a fallback
		if stmt == nil {
			return b.AppendTokenToError("Could not get statement")
		}
	}

	// fmt.Println(b.ParseBlock())

	return &Node{
		Type: "program",
	}, nil
}
