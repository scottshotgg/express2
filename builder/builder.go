package builder

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/scottshotgg/express-token"
)

type Node struct {
	Type     string
	Value    interface{}
	Metadata map[string]interface{}
	Left     *Node
	Right    *Node
}

type Builder struct {
	Tokens []token.Token
	Index  int
}

func New(tokens []token.Token) *Builder {
	return &Builder{Tokens: tokens}
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

func (b *Builder) ParseCall() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Ident {
		return nil, errors.New("Could not get ident after type")
	}

	// Create the ident
	ident := &Node{
		Type:  "ident",
		Value: b.Tokens[b.Index].Value.String,
	}

	// Skip over the ident token
	b.Index++

	// We are not allowing for named arguments right now
	args, err := b.ParseGroupOfExpressions()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:  "call",
		Value: ident,
		Metadata: map[string]interface{}{
			"args": args,
		},
	}, nil
}

func (b *Builder) ParseBlockStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.LBrace {
		return nil, errors.New("Could not get ident after type")
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

// ParseStatement ** does ** not look ahead
func (b *Builder) ParseStatement() (*Node, error) {
	switch b.Tokens[b.Index].Type {
	// case token.TypeDef:
	case token.Type:
		fmt.Println("**type**")
		// get declaration statement
		// need to check if there is an asterisk after
		// need to check if there is brackets after
		// MAKE a ParseType function that will return the type, check against user types, pointer, array, etc

		typeOf := b.Tokens[b.Index].Value.Type
		// For now: no user defined types ...
		decl, err := b.ParseDeclarationStatement()
		if decl != nil {
			decl.Value = typeOf
		}

		fmt.Println("decl", decl)

		return decl, err

	// TODO: what if types were expressions ...
	case token.Ident:
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

	// case token.Block:
	// 	// get a block

	case "*":
		// get a deref assignment

	case token.If:
		// get an if-else
		fmt.Println("if**")
		return b.ParseIfStatement()

	case token.For:
		// get a loop

	case token.Return:
		// get a return statement

	}

	// defer processing to next level higher
	return nil, nil
}

func (b *Builder) GetNextToken() (*token.Token, error) {
	if b.Index > len(b.Tokens)-1 {
		return nil, errors.New("Out of tokens")
	}

	return &b.Tokens[b.Index], nil
}

func (b *Builder) ParseForInStatement() (*Node, error) {
	return nil, nil
}

func (b *Builder) ParseForOfStatement() (*Node, error) {
	return nil, nil
}

func (b *Builder) ParseForStdStatement() (*Node, error) {
	return nil, nil
}

func (b *Builder) ParseForStatement() (*Node, error) {
	return nil, nil
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

func (b *Builder) ParseDeclarationStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Type {
		return nil, errors.New("Could not get ident after type")
	}

	// Increment over the type token
	b.Index++

	// Check that the next token is an ident
	if b.Tokens[b.Index].Type != token.Ident {
		return nil, errors.New("Could not get declaration statement without ident")
	}

	// Create the ident
	ident := &Node{
		Type:  "ident",
		Value: b.Tokens[b.Index].Value.String,
	}

	// Increment over the ident token
	b.Index++

	// Check for the equals token
	if b.Tokens[b.Index].Type != token.Assign {
		return &Node{
			Type: "decl",
			Left: ident,
		}, nil

		// return nil, errors.New("No equals found after ident")
	}

	// Increment over the equals
	b.Index++

	// Parse the right hand side
	expr, err := b.ParseExpression()
	if err != nil {
		return nil, errors.New("Could not get expression")
	}

	// Increment over the first part of the expression
	b.Index++

	return &Node{
		Type:  "decl",
		Left:  ident,
		Right: expr,
	}, nil
}

func (b *Builder) ParseAssignmentStatement() (*Node, error) {
	// Check that the next token is an ident
	if b.Tokens[b.Index].Type != token.Ident {
		return nil, errors.New("Could not get assignment statement without ident")
	}

	// Create the ident
	ident := &Node{
		Type:  "ident",
		Value: b.Tokens[b.Index].Value.String,
	}

	// Increment over the ident token
	b.Index++

	// Check for the equals token
	if b.Tokens[b.Index].Type != token.Assign {
		// return &Node{
		// 	Type: "decl",
		// 	Left: ident,
		// }, nil

		// This is where we would implement variable declarations without values
		// Leave it alone for now
		return nil, errors.New("No equals found after ident")
	}

	// Increment over the equals
	b.Index++

	// Parse the right hand side
	expr, err := b.ParseExpression()
	if err != nil {
		return nil, errors.New("Could not get expression")
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
	// This is where we will implement primary tier operators (* , /)

	return b.ParseTerm()
}

func (b *Builder) ParseTerm() (*Node, error) {
	// This is where we will implement secondary tier operators (+ , -)

	return b.ParseFactor()
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
	case token.Literal:
		// Increment over the assignment token
		return &Node{
			Type:  "literal",
			Value: b.Tokens[b.Index].Value.True,
			Metadata: map[string]interface{}{
				"type": b.Tokens[b.Index].Value.Type,
			},
		}, nil

	case token.Ident:
		// fmt.Println(b.Tokens[b.Index])
		return &Node{
			Type:  "ident",
			Value: b.Tokens[b.Index].Value.String,
		}, nil

	default:
		return nil, errors.Errorf("Could not parse statement from token; %+v", b.Tokens[b.Index])
	}
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
