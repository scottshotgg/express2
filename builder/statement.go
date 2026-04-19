package builder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	ast "github.com/scottshotgg/express-ast"
	lex "github.com/scottshotgg/express-lex"
	token "github.com/scottshotgg/express-token"
	"github.com/scottshotgg/express2/pkg/logger"
)

var (
	ErrNoEqualsFoundAfterIdent = errors.New("No equals found after ident in assignment")
)

func (b *Builder) ParseDeferStatement() (*Node, error) {
	// Check ourselves
	if b.Tokens[b.Index].Type != token.Defer {
		return nil, b.AppendTokenToError("Could not get group of statements")
	}

	// Step over the defer token
	b.Index++

	var stmt, err = b.ParseStatement()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type: "defer",
		Left: stmt,
	}, nil
}

func (b *Builder) ParseGroupOfStatements() (*Node, error) {
	// Check ourselves
	if b.Tokens[b.Index].Type != token.LParen {
		return nil, b.AppendTokenToError("Could not get group of statements")
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
		return nil, b.AppendTokenToError("Could not get for in")
	}

	// Step over the for token
	b.Index++

	// Parse the first ident before the preposition
	var ident, err = b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Step over the ident token
	b.Index++

	// Check for a second ident (for i, j over x)
	var ident2 *Node
	if b.Tokens[b.Index].Type == token.Separator && b.Tokens[b.Index].Value.Type == "comma" {
		// Step over the comma
		b.Index++

		ident2, err = b.ParseExpression()
		if err != nil {
			return nil, err
		}

		// Step over the second ident
		b.Index++
	}

	var prepType string

	switch b.Tokens[b.Index].Value.String {
	case "in":
		prepType = "forin"

	case "of":
		prepType = "forof"

	case "over":
		prepType = "forover"

	default:
		return nil, b.AppendTokenToError("Could not get preposition")
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
	b.Index++ // step past `}`

	metadata := map[string]interface{}{
		"start": ident,
		"end":   array,
	}

	if ident2 != nil {
		metadata["start2"] = ident2
	}

	return &Node{
		Type:     prepType,
		Value:    body,
		Metadata: metadata,
	}, nil
}

func (b *Builder) ParseForStdStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.For {
		return nil, b.AppendTokenToError("Could not get for std")
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
	b.Index++ // step past `}`

	return &node, nil
}

func (b *Builder) ParseWhileStatement() (*Node, error) {
	if b.Tokens[b.Index].Type != token.Loop {
		return nil, b.AppendTokenToError("Could not get while token")
	}

	// Skip 'while'
	b.Index++

	cond, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	b.Index++

	body, err := b.ParseBlockStatement()
	if err != nil {
		return nil, err
	}
	b.Index++ // step past `}`

	return &Node{
		Type:  "while",
		Left:  cond,
		Value: body,
	}, nil
}

func (b *Builder) ParseForStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.For {
		return nil, b.AppendTokenToError("Could not get for std")
	}

	if b.Index > len(b.Tokens)-2 {
		return nil, ErrOutOfTokens
	}

	// Scan forward from Index+1 looking for a preposition keyword (in/of/over)
	// before hitting a block opener. This handles both `for i in x` and `for i, j over x`.
	for i := b.Index + 1; i < len(b.Tokens); i++ {
		if b.Tokens[i].Type == token.LBrace {
			break
		}
		if b.Tokens[i].Type == token.Keyword {
			return b.ParseForPrepositionStatement()
		}
	}

	return b.ParseForStdStatement()
}

func (b *Builder) ParseIfStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.If {
		return nil, b.AppendTokenToError("Could not get if")
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
	b.Index++ // step past `}`

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
			b.Index++ // step past `}`
		}
	}

	return &n, nil
}

func (b *Builder) ParseMapBlockStatement() (*Node, error) {
	b.log.Debug("=== ParseMapBlockStatement called ===")
	b.log.Debugf("Index=%d, token=%+v", b.Index, b.Tokens[b.Index])

	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.LBrace {
		return nil, b.AppendTokenToError("Could not get left brace")
	}

	// Increment over the left brace token
	b.Index++

	b.log.Debugf("After LBrace, Index=%d", b.Index)

	var (
		stmt  *Node
		stmts []*Node
		err   error
	)

	for b.Index < len(b.Tokens) && b.Tokens[b.Index].Type != token.RBrace {
		b.log.Debugf("Loop: Index=%d, token=%+v", b.Index, b.Tokens[b.Index])

		// Parse the key (expression)
		stmt, err = b.ParseExpression()
		if err != nil {
			return nil, err
		}

		b.log.Debugf("After ParseExpression for key: stmt.Type=%s, Index=%d", stmt.Type, b.Index)

		// Check if the NEXT token is the separator token (= or :)
		if b.Index+1 < len(b.Tokens) && (b.Tokens[b.Index+1].Type == token.Assign || b.Tokens[b.Index+1].Type == token.Set) {
			b.log.Debug("Creating kv pair!")
			b.log.Debugf("Next token is separator: %+v", b.Tokens[b.Index+1])

			// Move past both the key and separator to get to the value
			b.Index += 2

			// Parse the value
			var value *Node
			value, err = b.ParseExpression()
			if err != nil {
				return nil, err
			}

			// Move past the value token
			b.Index++

			stmt = &Node{
				Type:  "kv",
				Left:  stmt,
				Right: value,
			}
			b.log.Debugf("Created kv: %+v", stmt)
		} else {
			b.log.Debug("NOT creating kv pair, next token is not = or :")
		}

		blob, _ := json.Marshal(stmt)
		b.log.Debug("kvstmtkv:", string(blob))

		// All statements in a map have to be key-value pairs
		if stmt.Type != "kv" {
			return nil, errors.Errorf("All statements in a map have to be key-value pairs: %+v\n", stmt)
		}

		b.log.Debug("stmt:", stmt)

		stmts = append(stmts, stmt)
	}

	b.log.Debug("=== ParseMapBlockStatement done ===")

	// Step over the right brace token
	b.Index++

	return &Node{
		Type:  "block",
		Value: stmts,
	}, nil
}

func (b *Builder) ParseEnumBlockStatement() (*Node, error) {
	// Increment over the enum keyword
	b.Index++

	var (
		ident *Node
		err   error
	)

	// Allow for named/typed enums
	if b.Tokens[b.Index].Type == token.Ident {
		ident, err = b.ParseExpression()
		if err != nil {
			return nil, errors.New("an error trying to get ident-type from enum")
		}

		b.Index++
	}

	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.LBrace {
		return nil, b.AppendTokenToError("Could not get left brace")
	}

	// Increment over the left brace token
	b.Index++

	var (
		stmt  *Node
		stmts []*Node
	)

	// FIXME: for now, setting values to enums is prohibited
	for b.Index < len(b.Tokens) &&
		b.Tokens[b.Index].Type != token.RBrace {
		stmt, err = b.ParseExpression()
		if err != nil {
			// Recover the parse if it gets the right error
			if err != ErrNoEqualsFoundAfterIdent {
				return nil, err
			}
		}

		b.Index++

		// All statements in a map have to be key-value
		// This isn't true wtf
		if stmt.Type != "assignment" && stmt.Type != "ident" {
			return nil, errors.Errorf("All statements in a enum have to be assignment or ident: %+v\n", stmt)
		}

		stmts = append(stmts, stmt)
	}

	// Step over the right brace token
	b.Index++

	var node = &Node{
		Type: "enum",
		Left: &Node{
			Type:  "block",
			Value: stmts,
		},
	}

	// Assert the type
	if ident != nil {
		node.Value = ident
	}

	return node, nil
}

func (b *Builder) ParseBlockStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.LBrace {
		return nil, b.AppendTokenToError("Could not get left brace")
	}

	// Increment over the left brace token
	b.Index++

	// Create a new child scope for the block with a unique name
	blockName := fmt.Sprintf("block_%d", b.BlockCounter)
	b.BlockCounter++

	newScope, err := b.ScopeTree.NewChildScope(blockName)
	if err != nil {
		return nil, err
	}
	b.ScopeTree = newScope

	var (
		stmt  *Node
		stmts []*Node
		err2  error
	)

	for b.Index < len(b.Tokens) &&
		b.Tokens[b.Index].Type != token.RBrace {
		// Skip statement separators (; and ,) between statements
		if b.Tokens[b.Index].Type == token.EOS ||
			b.Tokens[b.Index].Type == token.Separator ||
			b.Tokens[b.Index].Type == token.Comma {
			b.Index++
			continue
		}

		stmt, err2 = b.ParseStatement()
		if err2 != nil {
			return nil, err2
		}

		b.log.Debug("i am here", stmt)

		stmts = append(stmts, stmt)
	}

	// Leave the block scope (b.Index is ON the closing `}` — Pratt invariant)
	b.ScopeTree, err = b.ScopeTree.Leave()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:  "block",
		Value: stmts,
	}, nil
}

func (b *Builder) ParseReturnStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Return {
		return nil, b.AppendTokenToError("Could not get return")
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

// func (b *Builder) ParseDeclarationStatement(typeHint *TypeValue) (*Node, error) {
// 	var typeOf, err = b.ParseType(typeHint)
// 	if err != nil {
// 		return nil, err
// 	}

// 	fmt.Println("typeOf outside", typeOf)

// 	// Check that the next token is an ident
// 	if b.Tokens[b.Index].Type != token.Ident {
// 		return nil, b.AppendTokenToError("Could not get ident in declaration statement")
// 	}

// 	// Create the ident
// 	ident, err := b.ParseExpression()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var typeString = typeOf.Value.(string)
// 	if typeString == "map" || typeString == "object" || typeString == "struct" {
// 		b.ScopeTree, err = b.ScopeTree.NewChildScope(ident.Value.(string))
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	// // Check the scope map to make sure this hasn't been declared for the current scope
// 	// var node = b.ScopeTree.Local(ident.Value.(string))

// 	// // If the return value isn't nil then that means we found something in the local scope
// 	// if node != nil {
// 	// 	return nil, errors.Errorf("Variable already declared: %+v\n", node)
// 	// }

// 	// err = b.ScopeTree.Declare(ident)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// Increment over the ident token
// 	b.Index++

// 	// Check for the equals token
// 	if b.Tokens[b.Index].Type != token.Assign {
// 		return &Node{
// 			Type:  "decl",
// 			Value: typeOf,
// 			Left:  ident,
// 		}, nil

// 		// return nil, errors.New("No equals found after ident")
// 	}

// 	// Increment over the equals
// 	b.Index++

// 	// Parse the right hand side
// 	expr, err := b.ParseExpression()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Increment over the first part of the expression
// 	b.Index++

// 	// Leave the scope if we entered it above
// 	if typeString == "map" || typeString == "object" || typeString == "struct" {
// 		// Assign our scope back to the current one
// 		b.ScopeTree, err = b.ScopeTree.Leave()
// 		if err != nil {
// 			return nil, err
// 		}

// 		if typeString == "struct" {
// 			var v = &TypeValue{
// 				Composite: true,
// 				Type:      StruturedValue,
// 				Kind:      expr.Kind,
// 			}

// 			v.Props, err = b.extractPropsFromComposite(expr)
// 			if err != nil {
// 				return nil, err
// 			}

// 			err = b.ScopeTree.NewType(ident.Value.(string), v)
// 			if err != nil {
// 				return nil, err
// 			}
// 		}

// 		// Could defer this and then exit when we error?
// 	}

// 	var node = &Node{
// 		Type:  "decl",
// 		Value: typeOf,
// 		Left:  ident,
// 		Right: expr,
// 	}

// 	return node, b.ScopeTree.Declare(node)
// }

func (b *Builder) ParseTypeDeclarationStatement() (*Node, error) {
	// Check ourselves ...
	if b.Index < len(b.Tokens)-1 &&
		b.Tokens[b.Index].Type != token.TypeDef {
		return nil, b.AppendTokenToError("Could not get type declaration statement")
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
		return nil, b.AppendTokenToError("No equals found after ident in typedef")
	}

	// Increment over the equals
	b.Index++

	// // Parse the right hand side
	// typeOf, err := b.ParseType(nil)
	// if err != nil {
	// 	return nil, err
	// }

	// _, err = b.AddPrimiypeive(ident.Value.(string), typeOf)
	// if err != nil {
	// 	return nil, err
	// }

	typeOf, err := b.ParseExpression()
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

// ParseCBlock parses a raw C/C++ code injection block (`c { ... }`).
// THIS SHOULD NOT BE IN THE TRANSPILER; THIS SHOULD BE TAKEN CARE OF BY THE SEMANTIC STAGE
// BEFORE EVER REACHING THIS FAR. IT IS MERELY HERE AS A PLACEHOLDER SO I DO NOT FORGET
// MY OWN DECISIONS BECAUSE I CANT REMEMBER THINGS.
// For now the C block will be a direct injection of code into the final source.
func (b *Builder) ParseCBlock() (*Node, error) {
	var errStatement = "`c` blocks, as the are oh-so affectionately known within the Express community, are only implemented as a direct code injection at time. This will take some thinking; the compiler will have to `back-compile` the C/C++ code FROM the AST output of Clang and then translate that back into Express code essentially to check it"
	b.log.Warn(errStatement)

	// Skip over the `c` token
	b.Index++

	// Skip the opening brace
	if b.Index >= len(b.Tokens) || b.Tokens[b.Index].Type != token.LBrace {
		return nil, errors.New("Expected `{` after `c`")
	}
	b.Index++

	var (
		total []string
		found bool
	)

	// Gobble up all the code until the matching right brace; track nesting depth
	depth := 0
	for _, t := range b.Tokens[b.Index:] {
		b.log.Debug("t", t)

		if t.Type == token.LBrace {
			depth++
			total = append(total, "{")
			b.Index++
			continue
		}

		if t.Type == token.RBrace {
			if depth == 0 {
				found = true
				break
			}
			depth--
			total = append(total, "}")
			b.Index++
			continue
		}

		// Reconstruct the token as it appeared in source.
		// The lexer strips quotes from string/char literals, so we add them back.
		var raw string
		if t.Type == token.Literal {
			switch t.Value.Type {
			case token.StringType:
				raw = `"` + t.Value.String + `"`
			case token.CharType:
				raw = `'` + t.Value.String + `'`
			default:
				raw = t.Value.String
			}
		} else {
			raw = t.Value.String
		}

		total = append(total, raw)

		// Increment the index so that the gobbling reflects when we jump out of scope
		b.Index++
	}

	if !found {
		return nil, errors.New("No matching right brace found for c block")
	}

	// Skip over the rbrace
	b.Index++

	return &Node{
		Type:  "c",
		Value: strings.Join(total, " "),
	}, nil
}

func (b *Builder) ParseObjectStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Object {
		return nil, b.AppendTokenToError("Could not get object declaration statement")
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

	// Create a new child scope for the function
	b.ScopeTree, err = b.ScopeTree.NewChildScope(ident.Value.(string))
	if err != nil {
		return nil, err
	}

	// Check for the equals token
	if b.Tokens[b.Index].Type != token.Assign {
		return nil, b.AppendTokenToError("No equals found after ident in object def")
	}

	// Increment over the equals
	b.Index++

	// Parse the right hand side
	body, err := b.ParseBlockStatement()
	if err != nil {
		return nil, err
	}
	b.Index++ // step past `}`

	body.Kind = "object"

	// _, err = b.AddStructured(ident.Value.(string), body)
	// if err != nil {
	// 	return nil, err
	// }

	// Object does not get a type ... yet
	// var v = &TypeValue{
	// 	Composite: true,
	// 	Type:      StruturedValue,
	// 	Kind:      body.Kind,
	// }
	// v.Props, err = b.extractPropsFromComposite(body)
	// if err != nil {
	// 	return nil, err
	// }

	// // Increment over the first part of the expression
	// b.Index++

	// Assign our scope back to the current one
	b.ScopeTree, err = b.ScopeTree.Leave()
	if err != nil {
		return nil, err
	}

	// Again about the object not creating a type ...
	// err = b.ScopeTree.NewType(ident.Value.(string), v)
	// if err != nil {
	// 	return nil, err
	// }

	return &Node{
		Type:  "object",
		Left:  ident,
		Right: body,
	}, nil
}

func (b *Builder) ParseStructStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Struct {
		return nil, b.AppendTokenToError("Could not get struct declaration statement")
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

	// Create a new child scope for the function
	b.ScopeTree, err = b.ScopeTree.NewChildScope(ident.Value.(string))
	if err != nil {
		return nil, err
	}

	// Check for the equals token
	if b.Tokens[b.Index].Type != token.Assign {
		return nil, b.AppendTokenToError("No equals found after ident in struct def")
	}

	// Increment over the equals
	b.Index++

	// Parse the right hand side
	body, err := b.ParseBlockStatement()
	if err != nil {
		return nil, err
	}
	b.Index++ // step past `}`

	body.Kind = "struct"

	// Tag field declarations so the transpiler skips the const qualifier
	if children, ok := body.Value.([]*Node); ok {
		for _, child := range children {
			if child.Type == "decl" {
				if child.Metadata == nil {
					child.Metadata = map[string]interface{}{}
				}
				child.Metadata["is_field"] = true
			}
		}
	}

	// _, err = b.AddStructured(ident.Value.(string), body)
	// if err != nil {
	// 	return nil, err
	// }

	var v = &TypeValue{
		Composite: true,
		Type:      StruturedValue,
		Kind:      body.Kind,
	}

	v.Props, err = b.extractPropsFromComposite(body)
	if err != nil {
		return nil, err
	}

	// // Increment over the first part of the expression
	// b.Index++

	// Assign our scope back to the current one
	b.ScopeTree, err = b.ScopeTree.Leave()
	if err != nil {
		return nil, err
	}

	err = b.ScopeTree.NewType(ident.Value.(string), v)
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:  "struct",
		Left:  ident,
		Right: body,
	}, nil
}

func (b *Builder) ParseMapStatement() (*Node, error) {
	b.log.Debug("=== ParseMapStatement called ===")
	b.log.Debugf("Index=%d, token=%+v", b.Index, b.Tokens[b.Index])

	// Check ourselves ... (map can be token.Map or token.Type with Value.Type == "map")
	if b.Tokens[b.Index].Type != token.Map && !(b.Tokens[b.Index].Type == token.Type && b.Tokens[b.Index].Value.Type == "map") {
		return nil, b.AppendTokenToError("Could not get map declaration statement")
	}

	// Skip over the `map` token
	b.Index++

	b.log.Debugf("After map, Index=%d", b.Index)

	// Optional [K -> V] or [K, K -> V] type annotation
	var keyNode, valueNode *Node
	if b.Tokens[b.Index].Type == token.LBracket {
		var err2 error
		keyNode, valueNode, err2 = b.parseMapTypeAnnotation()
		if err2 != nil {
			return nil, err2
		}
	}

	// Create the ident
	ident, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	b.log.Debugf("After ParseExpression for ident: Index=%d", b.Index)

	// Increment over the ident token
	b.Index++

	// Create a new child scope for the function
	b.ScopeTree, err = b.ScopeTree.NewChildScope(ident.Value.(string))
	if err != nil {
		return nil, err
	}

	// Check for the equals token
	if b.Index >= len(b.Tokens) || b.Tokens[b.Index].Type != token.Assign {
		if keyNode != nil {
			// zero-init typed map declaration: map[string -> int] m
			b.ScopeTree, err = b.ScopeTree.Leave()
			if err != nil {
				return nil, err
			}
			return &Node{
				Type:     "map",
				Left:     ident,
				Metadata: map[string]interface{}{"key_node": keyNode, "value_node": valueNode},
			}, nil
		}
		return nil, b.AppendTokenToError("No equals found after ident in map declaration")
	}

	// Increment over the equals
	b.Index++

	b.log.Debugf("About to call ParseMapBlockStatement: Index=%d", b.Index)

	// Parse the right hand side
	body, err := b.ParseMapBlockStatement()
	if err != nil {
		return nil, err
	}

	body.Kind = "map"

	// Assign our scope back to the current one
	b.ScopeTree, err = b.ScopeTree.Leave()
	if err != nil {
		return nil, err
	}

	node := &Node{
		Type:  "map",
		Left:  ident,
		Right: body,
	}
	if keyNode != nil {
		node.Metadata = map[string]interface{}{"key_node": keyNode, "value_node": valueNode}
	}
	return node, nil
}

// parseMapTypeAnnotation parses the [K -> V] or [K, K -> V] type annotation.
// Returns the first key type node and the value type node (or a nested map node
// for multi-dimensional maps).  After return b.Index is one past the closing ].
func (b *Builder) parseMapTypeAnnotation() (keyNode, valueNode *Node, err error) {
	// consume [
	b.Index++

	// Collect key dimension types until we hit ->
	var keyNodes []*Node
	for {
		if b.Tokens[b.Index].Type != token.Type {
			return nil, nil, b.AppendTokenToError("expected type for map key in [K -> V]")
		}
		keyNodes = append(keyNodes, &Node{
			Type:  "type",
			Kind:  b.Tokens[b.Index].Value.Type,
			Value: b.Tokens[b.Index].Value.Type,
		})
		b.Index++

		switch {
		case b.Tokens[b.Index].Type == token.Separator && b.Tokens[b.Index].Value.Type == "comma":
			// Another key dimension: consume , and continue
			b.Index++
		case b.Tokens[b.Index].Type == token.Arrow:
			// Found the -> separator: consume and parse value
			b.Index++
			goto parseValue
		default:
			return nil, nil, b.AppendTokenToError("expected ',' or '->' in map type annotation")
		}
	}

parseValue:
	// Parse the value type
	valueNode, err = b.parseMapValueType()
	if err != nil {
		return nil, nil, err
	}

	if b.Tokens[b.Index].Type != token.RBracket {
		return nil, nil, b.AppendTokenToError("expected ] after map type annotation")
	}
	b.Index++ // consume ]

	// Right-fold extra key dimensions into nested map type nodes
	// e.g. keys=[string, string], value=int → key=string, value=map{string->int}
	for i := len(keyNodes) - 1; i >= 1; i-- {
		valueNode = &Node{
			Type: "type",
			Kind: "map",
			Metadata: map[string]interface{}{
				"key_node":   keyNodes[i],
				"value_node": valueNode,
			},
		}
	}

	return keyNodes[0], valueNode, nil
}

// parseMapValueType parses the value type in a map annotation.
// Handles plain types (token.Type) and nested map types (token.Map).
// After return b.Index is one past the last consumed value-type token.
func (b *Builder) parseMapValueType() (*Node, error) {
	switch b.Tokens[b.Index].Type {
	case token.Type:
		n := &Node{
			Type:  "type",
			Kind:  b.Tokens[b.Index].Value.Type,
			Value: b.Tokens[b.Index].Value.Type,
		}
		b.Index++
		return n, nil

	case token.Map:
		// nested map: map[K -> V]
		b.Index++ // consume map
		if b.Tokens[b.Index].Type != token.LBracket {
			return nil, b.AppendTokenToError("expected [ after nested map in map type annotation")
		}
		kn, vn, err := b.parseMapTypeAnnotation()
		if err != nil {
			return nil, err
		}
		return &Node{
			Type: "type",
			Kind: "map",
			Metadata: map[string]interface{}{
				"key_node":   kn,
				"value_node": vn,
			},
		}, nil

	default:
		return nil, b.AppendTokenToError("expected type for map value in map annotation")
	}
}

func (b *Builder) ParseLetStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Let {
		return nil, b.AppendTokenToError("Could not get let statement")
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
		return nil, b.AppendTokenToError("No equals found after ident in let")
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


func (b *Builder) ParseLiteralStatement() (*Node, error) {
	b.log.Debug("=== ParseLiteralStatement called ===")
	b.log.Debugf("Current token at Index(%d): %+v", b.Index, b.Tokens[b.Index])

	// Parse an expession
	// check the next token for a `:`
	// Parse another expression
	// Return a key-value pair

	// Get the expression
	var left, err = b.ParseExpression()
	if err != nil {
		return nil, err
	}

	blob, _ := json.Marshal(left)
	b.log.Debug("leftblobby:", string(blob))

	b.Index++

	b.log.Debugf("After parsing left, Index=%d, token: %+v", b.Index, b.Tokens[b.Index])

	switch b.Tokens[b.Index].Type {
	case token.Set:
		return b.ParseSet(left)

	default:
		return nil, errors.Errorf("ParseLiteralStatement not implemented for: %+v", b.Tokens[b.Index].Type)
	}
}

// ParseIdentStatement: Although idents are not statements, they do start many statements
// and this function serves to disambiguate those statements
func (b *Builder) ParseIdentStatement() (*Node, error) {
	// into: {type} [expr] = [expr]
	// Check that the next token is an ident
	// if b.Tokens[b.Index].Type != token.Ident {
	// 	return nil, b.AppendTokenToError("Could not get assignment statement without ident")
	// }

	// TODO(scottshotgg): this is a super stupid and hacky way of doing the function call parsing. It really needs to be part of the expression parsing
	// var value = b.Tokens[b.Index].Value.String
	// if cFuncs[value] {
	// 	defer func() { b.Index++ }()
	// 	b.Index++
	// 	return b.ParseCall(&Node{
	// 		Type:  "ident",
	// 		Value: value,
	// 	})
	// }

	// Parse the first ident; this COULD be a type
	identOrType, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}
	// Increment over the ident token
	b.Index++

	if identOrType.Type == "call" || identOrType.Type == "inc" || identOrType.Type == "dec" {
		return identOrType, nil
	}

	if b.Index > len(b.Tokens)-1 {
		return identOrType, nil
	}

	// Default to assignment so that we only have to write the assignment
	// logic once and can use a fallthrough case
	var node = &Node{
		Type: "assignment",
		// Value: This will be a type set to identOrType if it is declaration
		Left: identOrType,
		// Right: expr,
	}

	b.log.Debug("identOrType, err", identOrType, err, b.Tokens[b.Index].Type)

	switch b.Tokens[b.Index].Type {
	case token.Type:
		// var <primitive-type> <ident> = <expr> pattern (e.g. `var int x = 5`)
		if identOrType.Kind != "var" {
			return nil, b.AppendTokenToError("unexpected type after non-var type")
		}
		node.Type = "decl"
		actualType, err := b.ParseExpression()
		if err != nil {
			return nil, err
		}
		node.Value = actualType
		if node.Metadata == nil {
			node.Metadata = map[string]interface{}{}
		}
		node.Metadata["mutable"] = true
		b.Index++
		if b.Index > len(b.Tokens)-1 {
			return nil, b.AppendTokenToError("expected identifier after type in var declaration")
		}
		if b.Tokens[b.Index].Type != token.Ident {
			return nil, b.AppendTokenToError("expected identifier after type in var declaration")
		}
		node.Left, err = b.ParseExpression()
		if err != nil {
			return nil, err
		}
		b.Index++
		if b.Index > len(b.Tokens)-1 || b.Tokens[b.Index].Type != token.Assign {
			return node, nil
		}
		b.Index++ // step over =
		node.Right, err = b.ParseExpression()
		if err != nil {
			return nil, err
		}
		b.Index++
		return node, nil

	case token.Ident:
		// var <UserType> <ident> = <expr> pattern (e.g. `var Person alice = { ... }`)
		if identOrType.Kind == "var" {
			typeName := b.Tokens[b.Index].Value.String
			if b.ScopeTree.GetType(typeName) != nil {
				node.Type = "decl"
				actualType, err := b.ParseExpression()
				if err != nil {
					return nil, err
				}
				node.Value = actualType
				if node.Metadata == nil {
					node.Metadata = map[string]interface{}{}
				}
				node.Metadata["mutable"] = true
				b.Index++
				if b.Index > len(b.Tokens)-1 {
					return nil, b.AppendTokenToError("expected identifier after type in var declaration")
				}
				if b.Tokens[b.Index].Type != token.Ident {
					return nil, b.AppendTokenToError("expected identifier after type in var declaration")
				}
				node.Left, err = b.ParseExpression()
				if err != nil {
					return nil, err
				}
				b.Index++
				if b.Index > len(b.Tokens)-1 || b.Tokens[b.Index].Type != token.Assign {
					return node, nil
				}
				b.Index++ // step over =
				node.Right, err = b.ParseExpression()
				if err != nil {
					return nil, err
				}
				b.Index++
				return node, nil
			}
		}

		/*
			In this case, we have two idents back to back which leads us
			to make the only informed decision we can; that the first ident
			was a type, like in cases such as:
				int i = 0
		*/

		// Set the proper node values
		node.Type = "decl"
		node.Value = identOrType

		b.log.Debug("got another ident", b.Tokens[b.Index], node)

		node.Left, err = b.ParseExpression()
		if err != nil {
			return nil, err
		}

		b.log.Debug("node.Left", node.Left)

		// Step over the real ident
		b.Index++

		/*
			 This is the case where we do not have any more tokens but still
			 could be valid for cases like:
					int i
		*/
		if b.Index > len(b.Tokens)-1 {
			return node, nil
		}

		if b.Tokens[b.Index].Type != token.Assign {
			return node, nil
		}

		fallthrough

	case token.Assign:
		/*
			This is the assignment case where a simple assignment is as:
				i = 0
		*/
		b.log.Debug("got assign")
		b.log.Debugf("Index before ParseExpression: %d", b.Index)

		// Step over the assign
		b.Index++

		node.Right, err = b.ParseExpression()
		if err != nil {
			return nil, err
		}

		b.log.Debug("node.Right:", node.Right)
		b.log.Debugf("Index after ParseExpression: %d", b.Index)

		b.Index++

		b.log.Debugf("Index before return: %d", b.Index)
		return node, nil

	case token.LBrace:
		// Set the proper node values
		node.Type = "decl"
		node.Value = identOrType

		node.Right, err = b.ParseExpression()
		if err != nil {
			return nil, err
		}

		b.Index++

		blob, _ := json.Marshal(node)
		b.log.Debug("blobby:", string(blob))

		return node, nil

	case token.Set:
		b.Index++

		// Set the proper node values
		node.Type = "kv"
		node.Left = identOrType

		node.Right, err = b.ParseExpression()
		if err != nil {
			return nil, err
		}

		b.Index++

		blob, _ := json.Marshal(node)
		b.log.Debug("setblobby:", string(blob))

		return node, nil

	case token.AddAssign, token.SubAssign, token.MulAssign, token.DivAssign:
		// Desugar x += y  →  x = x + y
		opMap := map[string]string{
			token.AddAssign: "+",
			token.SubAssign: "-",
			token.MulAssign: "*",
			token.DivAssign: "/",
		}
		op := opMap[b.Tokens[b.Index].Type]

		// Skip the compound operator
		b.Index++

		rhs, err := b.ParseExpression()
		if err != nil {
			return nil, err
		}

		b.Index++

		node.Type = "assignment"
		node.Left = identOrType
		node.Right = &Node{
			Type:  "binop",
			Value: op,
			Left:  identOrType,
			Right: rhs,
		}
		return node, nil

	// Just return the ident if you don't know what to do
	// this will defer the judgement to the next statement up
	default:
		return identOrType, nil
	}
}

func (b *Builder) ParsePackageStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Package {
		return nil, b.AppendTokenToError("Could not get package statement")
	}

	// Step over the package token
	b.Index++

	expr, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Step over the literal
	b.Index++

	// Get the rest of the statements
	// We will need to get all of the files in the folder
	// Grab the rest of the statements in the folder and assign them to this node
	// Use the semantic compiler to sort out multiple files in a package

	var stmts []*Node
	for b.Index < len(b.Tokens)-1 {
		stmt, err := b.ParseStatement()
		if err != nil {
			if err == ErrOutOfTokens {
				break
			}

			return nil, err
		}

		// Just a fallback; probably won't need it later
		if stmt == nil {
			return nil, b.AppendTokenToError("Statement was nil")
		}

		stmts = append(stmts, stmt)
		b.log.Debug("STMT", stmt)
	}

	return &Node{
		Type: "package",
		Left: expr,
		Right: &Node{
			Type:  "block",
			Value: stmts,
		},
	}, nil
}

// ParseUseStatement parses a `use "path" as alias` statement.
// TODO: add this back in to ParseStatement once the `as` keyword and aliasing semantics are defined.
func (b *Builder) ParseUseStatement() (*Node, error) {
	// Step over the use token
	b.Index++

	// This expression takes the same rules as import/include with quotes and no quotes
	expr, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Step over the literal
	b.Index++

	// With a use statement, we are expecting an as operation afterwards and then another _ident_
	// I don't know or want to add "as" as a keyword right now, not sure it has much use; however, it needs to be checked nevertheless
	expr1, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	if expr1.Type != "ident" {
		return nil, errors.Errorf("Expecting \"as\" keyword after use expression, found: %+v", expr)
	}

	// Hop over the "as"
	b.Index++

	// Next up: we are expecting an _ident_; parse it as an expression so operation rules will apply
	expr1, err = b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// May have to mangle the names for this ;_;
	if expr1.Type != "ident" {
		return nil, errors.Errorf("Expecting ident expression after as keyword, found: %+v", expr)
	}

	// And finally, hop over the ending ident
	b.Index++

	return &Node{
		Type:  "use",
		Left:  expr,
		Right: expr1,
	}, nil
}

func (b *Builder) parseFileImport(filename string) (*Node, *ScopeTree, error) {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		// Fall back to $EXPRPATH/lib/<name>.expr for stdlib packages
		if libpath := os.Getenv("EXPRPATH"); libpath != "" {
			libFile := filepath.Join(libpath, "lib", filename+".expr")
			source, err = ioutil.ReadFile(libFile)
		}
		if err != nil {
			return nil, nil, err
		}
	}

	b.log.Debug("source", string(source))

	// Lex and tokenize the source code
	tokens, err := lex.New(string(source)).Lex()
	if err != nil {
		return nil, nil, err
	}

	// Compress certain tokens;
	tokens, err = ast.CompressTokens(tokens)
	if err != nil {
		return nil, nil, err
	}

	// Build the AST
	b2 := New(tokens, logger.Noop())
	ast, err := b2.BuildAST()
	if err != nil {
		return nil, nil, err
	}

	// fmt.Printf("ast %+v\n", ast.Value.([]*Node)[0].Left.Value.(string))

	// TODO: extremely unsafe, fix this
	return ast, b2.ScopeTree, nil
}

func (b *Builder) ParseImportStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Import {
		return nil, b.AppendTokenToError("Could not get import statement")
	}

	// Step over the import token
	b.Index++

	// Special case: `import c` means "import standard C library headers"
	// It's handled before ParseExpression because `c` is a keyword token
	// and the Pratt parser doesn't know how to handle keywords as expressions.
	if b.Index < len(b.Tokens) && b.Tokens[b.Index].Type == token.C {
		b.Index++
		return &Node{
			Type: "import",
			Kind: "c",
		}, nil
	}

	expr, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Now that we have the expression, we need to go parse that file
	// 1. Parse the file
	// 2. Use a variable to link the file
	// 3. Normal selection checking after that
	// 4. Take special care for transpileImportStatement

	b.log.Debug("expr.Kind", expr.Value.(string))

	if expr.Value.(string) == "c" {
		b.Index++
		return &Node{
			Type: "import",
			Kind: "c",
		}, nil
	}

	// TODO: Later on we will need to check this whether it is a module, file, or remote
	ast, scope, err := b.parseFileImport(expr.Value.(string))
	if err != nil {
		return nil, err
	}

	var split = strings.Split(expr.Value.(string), "/")
	var namespace = split[len(split)-1]

	if strings.HasSuffix(namespace, ".expr") {
		expr.Value = namespace[:len(namespace)-5]
	}

	// Set the new scope trees value to the scope retrieved from the file
	b.ScopeTree.Imports[expr.Value.(string)] = scope
	b.ScopeTree.Vars[expr.Value.(string)] = ast

	// Step over the literal
	b.Index++

	return &Node{
		Type:  "import",
		Left:  expr,
		Right: ast,
	}, nil
}

func (b *Builder) ParseIncludeStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Include {
		return nil, b.AppendTokenToError("Could not get include statement")
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

func (b *Builder) ParseLaunchStatement() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Launch {
		return nil, b.AppendTokenToError("Could not get launch statement")
	}

	// Step over the import token
	b.Index++

	// Might need to make this an explicit function call later
	expr, err := b.ParseStatement()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type: "launch",
		Left: expr,
	}, nil
}

func (b *Builder) ParseFunctionStatement() (*Node, error) {
	// Check ourselves
	if b.Tokens[b.Index].Type != token.Function {
		return nil, b.AppendTokenToError("Could not get function")
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
		return nil, b.AppendTokenToError("Could not get ident after function token")
	}

	// Set the name of the function
	node.Kind = b.Tokens[b.Index].Value.String

	// Step over the ident token
	b.Index++

	// Method declaration: func Receiver.MethodName(...)
	if b.Tokens[b.Index].Type == token.Accessor {
		b.Index++ // skip '.'
		if b.Tokens[b.Index].Type != token.Ident {
			return nil, b.AppendTokenToError("expected method name after '.'")
		}
		node.Metadata["receiver"] = node.Kind
		node.Kind = b.Tokens[b.Index].Value.String
		b.Index++
	}

	if b.Tokens[b.Index].Type != token.LParen {
		return nil, b.AppendTokenToError("Could not get left paren")
	}

	// Create a new child scope for function arguments so they don't collide
	// with args from other functions at the same level.
	b.ScopeTree, err = b.ScopeTree.NewChildScope("fn_" + node.Kind)
	if err != nil {
		return nil, err
	}

	args, err := b.ParseGroupOfStatements()
	if err != nil {
		return nil, err
	}

	if args != nil {
		node.Metadata["args"] = args
	}

	// If the next token is not a left brace, then we have returns
	if b.Tokens[b.Index].Type != token.LBrace {
		// Check for multiple returns; another left paren
		if b.Tokens[b.Index].Type == token.LParen {
			// This should be a group of idents for the types
			node.Metadata["returns"], err = b.ParseGroupOfExpressions()
			if err != nil {
				return nil, err
			}

			b.Index++

		} else {
			// Single return type: primitive or user-defined
			typeNode, typeErr := b.ParseTypeExpr()
			if typeErr != nil {
				return nil, errors.Errorf("could not parse return type for %s: %v", node.Kind, typeErr)
			}
			node.Metadata["returns"] = &Node{
				Type:  "egroup",
				Value: []*Node{typeNode},
			}
			b.Index++
		}
	}

	b.log.Debug("node.Metadata[returns]:", node.Metadata["returns"])

	// We are not supporting named arguments for now
	// Check for the return type token

	// If the function is named main then check that it returns an int
	// If it doesn't have any return type then apply an int return
	// If it already has a return type that is not int then that is an error
	if node.Kind == "main" {
		if node.Metadata["returns"] != nil {
			return nil, errors.New("main can not have any return type specified")
		}
	}

	node.Value, err = b.ParseBlockStatement()
	if err != nil {
		return nil, err
	}
	b.Index++ // step past `}`

	// Leave the function argument scope
	b.ScopeTree, err = b.ScopeTree.Leave()
	if err != nil {
		return nil, err
	}

	// node.Value = addDeferDeclarationToBlock(block)

	// Declare the type in the upper scope
	err = b.ScopeTree.Declare(&node)
	if err != nil {
		return nil, err
	}

	return &node, nil
}

// func addDeferDeclarationToBlock(n *Node) *Node {
// 	var stmts, ok = n.Value.([]*builder.Node)
// 	stmts = append([]*builder.Node(&Node{
// 		Type: "defer"
// 	}, stmts...))
// }

func (b *Builder) ParseDerefStatement() (*Node, error) {
	if b.Tokens[b.Index].Type != token.PriOp {
		return nil, b.AppendTokenToError("Could not get deref statement without *")
	}

	deref, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Increment over the ident token
	b.Index++

	if b.Index > len(b.Tokens) {
		return deref, nil
	}

	var t = "decl"

	if deref.Kind != "type" {
		t = "assignment"
		// Check for the equals token
		if b.Tokens[b.Index].Type != token.Assign {
			if deref.Type == "call" {
				return deref, nil
			}

			return nil, b.AppendTokenToError(fmt.Sprintf("No equals found after ident in deref: %+v", b.Tokens[b.Index]))
		}

		// Increment over the equals
		b.Index++
	}

	// Parse the right hand side
	expr, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Increment over the first part of the expression
	b.Index++

	return &Node{
		Type:  t,
		Left:  deref,
		Right: expr,
	}, nil
}

// TODO: what if types were expressions ...

// ParseStatement ** does ** not look ahead
func (b *Builder) ParseStatement() (*Node, error) {
	b.log.Debugf("=== ParseStatement called, Index=%d, token=%+v ===", b.Index, b.Tokens[b.Index])
	switch b.Tokens[b.Index].Type {

	case token.Launch:
		return b.ParseLaunchStatement()

	case token.Enum:
		return b.ParseEnumBlockStatement()

	case token.Map:
		return b.ParseMapStatement()

	case token.PriOp:
		return b.ParseDerefStatement()

	case token.Package:
		return b.ParsePackageStatement()

	case token.Import:
		return b.ParseImportStatement()

	case token.Include:
		return b.ParseIncludeStatement()

	case token.TypeDef:
		return b.ParseTypeDeclarationStatement()

	case token.Struct:
		return b.ParseStructStatement()

	case token.Object:
		return b.ParseObjectStatement()

	// For literal and idents, we will need to figure out what
	// kind of statement it is
	case token.Literal:
		return b.ParseLiteralStatement()

	case token.Ident:
		// `c` is not a lexer keyword (it would split identifiers containing the
		// letter c).  Detect a c block by looking for an ident whose value is
		// "c" followed by an LBrace.
		if b.Tokens[b.Index].Value.String == "c" &&
			b.Index+1 < len(b.Tokens) &&
			b.Tokens[b.Index+1].Type == token.LBrace {
			return b.ParseCBlock()
		}

		if cFuncs[b.Tokens[b.Index].Value.String] {
			var value = b.Tokens[b.Index].Value.String
			b.Index++

			defer func() {
				if b.Index < len(b.Tokens) && b.Tokens[b.Index].Type == token.RParen {
					b.Index++
				}
			}()

			return b.ParseCall(&Node{
				Type:  "ident",
				Value: value,
			})
		}

		var nn = b.ScopeTree.Get(b.Tokens[b.Index].Value.String)
		if nn != nil {
			switch nn.Type {
			case "function":
				var value = b.Tokens[b.Index].Value.String
				b.Index++

				defer func() {
					if b.Index < len(b.Tokens) && b.Tokens[b.Index].Type == token.RParen {
						b.Index++
					}
				}()

				return b.ParseCall(&Node{
					Type:  "ident",
					Value: value,
				})
			}
		}

		var n, err = b.ParseIdentStatement()
		if err != nil {
			return nil, err
		}

		if n.Type == "decl" {
			blob, _ := json.Marshal(n)
			b.log.Debug("nblobn:", string(blob))

			err = b.ScopeTree.Declare(n)
			if err != nil {
				return nil, err
			}
			blob, _ = json.Marshal(b.ScopeTree)
			b.log.Debug("ScopeTree:", string(blob))
		}

		return n, nil

	case token.Type:
		// Check if this is a special type keyword that represents a statement
		// like `map m = { ... }` or `struct Point = { ... }`
		// or `map string -> int scores = { ... }`
		if b.Index+1 < len(b.Tokens) {
			nextToken := b.peekAt(1)
			// If the next token is an ident and followed by =, it might be a statement
			if nextToken.Type == token.Ident && b.Index+2 < len(b.Tokens) && b.peekAt(2).Type == token.Assign {
				// Check if this is a known type keyword that should be handled as a statement
				if b.Tokens[b.Index].Value.Type == "map" || b.Tokens[b.Index].Value.Type == "struct" || b.Tokens[b.Index].Value.Type == "object" {
					// These should be parsed as statements, not type declarations
					switch b.Tokens[b.Index].Value.Type {
					case "map":
						return b.ParseMapStatement()
					case "struct":
						return b.ParseStructStatement()
					case "object":
						return b.ParseObjectStatement()
					}
				}
			}
			// Handle full type annotation: map string -> int ident = ...
			b.log.Debugf("DEBUG: Checking full type annotation... Token=%s Value.Type=%s", b.Tokens[b.Index].Value.String, b.Tokens[b.Index].Value.Type)
			if b.Tokens[b.Index].Value.Type == "map" || b.Tokens[b.Index].Value.Type == "struct" || b.Tokens[b.Index].Value.Type == "object" {
				b.log.Debug("DEBUG: Is map/struct/object")
				// Check if this looks like a type annotation (type -> type pattern)
				if nextToken.Type == token.Type && b.Index+2 < len(b.Tokens) {
					b.log.Debug("DEBUG: Next token is Type and we have enough tokens")
					nextNextToken := b.peekAt(2)
					// Check for -> pattern (SecOp = sub, GThan = >)
					// The -> is tokenized as two separate tokens: SecOp (-) and GThan (>)
					if nextNextToken.Type == token.SecOp && b.Index+3 < len(b.Tokens) {
						b.log.Debug("DEBUG: Found SecOp (-)")
						nextNextNextToken := b.peekAt(3)
						if nextNextNextToken.Type == token.GThan && b.Index+4 < len(b.Tokens) {
							b.log.Debug("DEBUG: Found GThan (>)")
							nextNextNextNextToken := b.peekAt(4)
							if nextNextNextNextToken.Type == token.Type && b.Index+5 < len(b.Tokens) {
								b.log.Debug("DEBUG: Found second type")
								nextNextNextNextNextToken := b.peekAt(5)
								if nextNextNextNextNextToken.Type == token.Ident && b.Index+6 < len(b.Tokens) && b.peekAt(6).Type == token.Assign {
									b.log.Debug("DEBUG: Calling ParseMapStatement!")
									switch b.Tokens[b.Index].Value.Type {
									case "map":
										return b.ParseMapStatement()
									case "struct":
										return b.ParseStructStatement()
									case "object":
										return b.ParseObjectStatement()
									}
								}
							}
						}
					}
				}
			}

			// Handle map[K, V] typed map syntax: next token is LBracket followed by a type keyword
			if b.Tokens[b.Index].Value.Type == "map" && nextToken.Type == token.LBracket &&
				b.Index+2 < len(b.Tokens) && b.peekAt(2).Type == token.Type {
				return b.ParseMapStatement()
			}
		}

		var n, err = b.ParseIdentStatement()
		if err != nil {
			return nil, err
		}

		err = b.ScopeTree.Declare(n)
		if err != nil {
			return nil, err
		}

		return n, nil

	}

	// Everything else — if, let, for, while, func, blocks, defer, return — flows through Pratt.
	// prattParse uses the Pratt invariant (ON last token); convert to statement contract (one past).
	n, err := b.prattParse(PrecNone)
	if err != nil {
		return nil, err
	}
	b.Index++ // Pratt invariant → statement contract
	return n, nil
}
