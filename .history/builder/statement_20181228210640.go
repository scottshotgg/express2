package builder

import (
	"github.com/pkg/errors"
	"github.com/scottshotgg/express-token"
)

func (b *Builder) ParseGroupOfStatements() (*Node, error) {
	// Check ourselves
	if b.Tokens[b.Index].Type != token.LParen {
		return b.AppendTokenToError("Could not get group of statements")
	}

	// Skip over the left paren token
	b.Index++

	var (
		stmt  *Node
		stmts []*Node
		err   error
	)

	for b.Tokens[b.Index].Type != token.RParen {
		stmt, err = b.ParseStatement()
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

func (b *Builder) ParseForPrepositionStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.For {
		return b.AppendTokenToError("Could not get for in")
	}

	// Step over the for token
	b.Index++

	// Parse the ident before the `in` token
	var ident, err = b.ParseExpression()
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
	array, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// TODO: we also need to parse here to figure out if this array is an ident

	b.Index++

	body, err := b.ParseBlockStatement()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:  prepType,
		Value: body,
		Metadata: map[string]interface{}{
			"start": ident,
			"end":   array,
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

	for b.Index < len(b.Tokens) &&
		b.Tokens[b.Index].Type != token.RBrace {
		stmt, err = b.ParseStatement()
		if err != nil {
			return nil, err
		}
		// fmt.Println("i am here", stmt)

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
	if b.Index < len(b.Tokens) &&
		b.Tokens[b.Index].Value.Type == "newline" {
		b.Index++

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

func (b *Builder) ParseDeclarationStatement() (*Node, error) {
	var typeOf, err = b.ParseType()
	if err != nil {
		return nil, err
	}

	// Check that the next token is an ident
	if b.Tokens[b.Index].Type != token.Ident {
		return b.AppendTokenToError("Could not get ident in declaration statement")
	}

	// Create the ident
	ident, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// // Check the scope map to make sure this hasn't been declared for the current scope
	// var node = b.ScopeTree.Local(ident.Value.(string))

	// // If the return value isn't nil then that means we found something in the local scope
	// if node != nil {
	// 	return nil, errors.Errorf("Variable already declared: %+v\n", node)
	// }

	// err = b.ScopeTree.Declare(ident)
	// if err != nil {
	// 	return nil, err
	// }

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
	if b.Index < len(b.Tokens)-1 &&
		b.Tokens[b.Index].Type != token.TypeDef {
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

	// // Increment over the first part of the expression
	// b.Index++

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

	b.TypeMap[ident] = *body

	// // Increment over the first part of the expression
	// b.Index++

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

	// // Check the scope map to make sure this hasn't been declared for the current scope
	// var node = b.ScopeTree.Get(ident.Value.(string))

	// // If the return value isn't nil then that means we found something in the local scope
	// if node == nil {
	// 	return nil, errors.Errorf("Use of undeclared identifier: %+v\n", ident)
	// }

	// *ident = *node

	// Increment over the ident token
	b.Index++

	if b.Index > len(b.Tokens)-1 {
		return ident, nil
	}

	// Check for the equals token
	if b.Tokens[b.Index].Type != token.Assign {
		if ident.Type == "call" {
			return ident, nil
		}

		return b.AppendTokenToError("No equals found after ident in assignment")
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

	// Step over the literal
	b.Index++

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

	// Step over the literal
	b.Index++

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

	// Step over the literal
	b.Index++

	return &Node{
		Type: "include",
		Left: expr,
	}, nil
}

func (b *Builder) ParseFunctionStatement() (*Node, error) {
	// Check ourselves
	if b.Tokens[b.Index].Type != token.Function {
		return b.AppendTokenToError("Could not get function")
	}

	// Step over the function token
	b.Index++

	var (
		err  error
		node = Node{
			Type:     "function",
			Metadata: map[string]interface{}{},
		}
	)

	// Named function
	if b.Tokens[b.Index].Type != token.Ident {
		return b.AppendTokenToError("Could not get ident after function token")
	}

	node.Kind = b.Tokens[b.Index].Value.String

	// Step over the ident token
	b.Index++

	if b.Tokens[b.Index].Type != token.LParen {
		return b.AppendTokenToError("Could not get left paren")
	}

	args, err := b.ParseGroupOfStatements()
	if err != nil {
		return nil, err
	}

	if args != nil {
		node.Metadata["args"] = args
	}

	// Might want to avoid putting this here if we don't have any

	// We are not supporting multiple returns for now
	// // Check for multiple returns;another left paren
	// if b.Tokens[b.Index].Type == token.LParen {
	// 	return nil, errors.New("Could not get returns")
	// }

	var returnType = b.Tokens[b.Index].Value.String

	// We are not supporting named arguments for now
	// Check for the return type token
	if b.Tokens[b.Index].Type == token.Type {
		node.Metadata["returns"] = &Node{
			Type:  "type",
			Value: returnType,
		}

		// Step over the type token
		b.Index++
	}

	// If the function is named main then check that it returns an int
	// If it doesn't have any return type then apply an int return
	// If it already has a return type that is not int then that is an error
	if node.Kind == "main" {
		if node.Metadata["returns"] != nil {
			// Add this later
			// if len(node.Metadata["returns"].([]*Node)) > 1 {
			// 	return nil, errors.New("main can only have one return")
			// }

			if returnType != "int" {
				return nil, errors.New("main can only return an int type")
			}
		}

		// Apply the int return
		node.Metadata["returns"] = &Node{
			Type:  "type",
			Value: "int",
		}
	}

	node.Value, err = b.ParseBlockStatement()
	if err != nil {
		return nil, err
	}

	return &node, nil
}

func (b *Builder) ParseDerefStatement() (*Node, error) {
	// if b.Tokens[b.Index].Type != token.Ident {
	// 	return b.AppendTokenToError("Could not get assignment statement without ident")
	// }

	deref, err := b.ParseDerefExpression()
	if err != nil {
		return nil, err
	}

	// // Increment over the ident token
	// b.Index++

	if b.Index > len(b.Tokens) {
		return deref, nil
	}

	// Check for the equals token
	if b.Tokens[b.Index].Type != token.Assign {
		if deref.Type == "call" {
			return deref, nil
		}

		return b.AppendTokenToError("No equals found after ident in assignment")
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
		Left:  deref,
		Right: expr,
	}, nil
}

// TODO: what if types were expressions ...

// ParseStatement ** does ** not look ahead
func (b *Builder) ParseStatement() (*Node, error) {
	var (
		node *Node
		err  error
	)

	switch b.Tokens[b.Index].Type {

	case token.PriOp:
		node, err = b.ParseDerefStatement()

	case token.Package:
		node, err = b.ParsePackageStatement()

	case token.Import:
		node, err = b.ParseImportStatement()

	case token.Include:
		node, err = b.ParseIncludeStatement()

	case token.TypeDef:
		node, err = b.ParseTypeDeclarationStatement()

	case token.Type:
		node, err = b.ParseDeclarationStatement()

	case token.Ident:
		node, err = b.ParseAssignmentStatement()

	case token.Function:
		node, err = b.ParseFunctionStatement()

	case token.LBrace:
		node, err = b.ParseBlockStatement()

	case token.Struct:
		node, err = b.ParseStructStatement()

	case token.Let:
		node, err = b.ParseLetStatement()

	case token.If:
		node, err = b.ParseIfStatement()

	case token.For:
		node, err = b.ParseForStatement()

	case token.Return:
		node, err = b.ParseReturnStatement()

	default:
		return b.AppendTokenToError("Could not get statement from")
	}

	return node, err
}
