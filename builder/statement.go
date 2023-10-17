package builder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/pkg/errors"
	ast "github.com/scottshotgg/express-ast"
	lex "github.com/scottshotgg/express-lex"
	token "github.com/scottshotgg/express-token"
)

var (
	ErrNoEqualsFoundAfterIdent = errors.New("No equals found after ident in assignment")
)

func (b *Builder) ParseDeferStmt() (*Node, error) {
	// Check ourselves
	if b.Tokens[b.Index].Type != token.Defer {
		return nil, b.AppendTokenToError("Could not get group of statements")
	}

	// Step over the defer token
	b.Index++

	var stmt, err = b.ParseStmt()
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
		stmt, err = b.ParseStmt()
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

func (b *Builder) ParseForPrepositionStmt() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.For {
		return nil, b.AppendTokenToError("Could not get for in")
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

	body, err := b.ParseBlockStmt()
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

func (b *Builder) ParseForStdStmt() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.For {
		return nil, b.AppendTokenToError("Could not get for std")
	}

	// Step over the for token
	b.Index++

	// Parse the declaration or assignment statement
	stmt, err := b.ParseStmt()
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

	node.Value, err = b.ParseBlockStmt()
	if err != nil {
		return nil, err
	}

	return &node, nil
}

func (b *Builder) ParseForEverStmt() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.For {
		return nil, b.AppendTokenToError("Could not get for ever")
	}

	// Step over the for token
	b.Index++

	val, err := b.ParseBlockStmt()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:  "forever",
		Value: val,
	}, nil
}

func (b *Builder) ParseForStmt() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.For {
		return nil, b.AppendTokenToError("Could not get for std")
	}

	if b.Index > len(b.Tokens)-2 {
		return nil, ErrOutOfTokens
	}

	// For right now just look ahead two
	if b.Tokens[b.Index+2].Type == token.Keyword {
		return b.ParseForPrepositionStmt()
	}

	// For-ever statement
	if b.Tokens[b.Index+1].Type == token.LBrace {
		return b.ParseForEverStmt()
	}

	return b.ParseForStdStmt()
}

func (b *Builder) ParseIfStmt() (*Node, error) {
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

	n.Left, err = b.ParseBlockStmt()
	if err != nil {
		return nil, err
	}

	if b.Index < len(b.Tokens)-1 && b.Tokens[b.Index].Type == token.Else {
		// Step over the else token
		b.Index++

		// Check for an else if
		if b.Tokens[b.Index].Type == token.If {
			n.Right, err = b.ParseIfStmt()
			if err != nil {
				return nil, err
			}
		} else {
			n.Right, err = b.ParseBlockStmt()
			if err != nil {
				return nil, err
			}
		}
	}

	return &n, nil
}

func (b *Builder) ParseMapBlockStmt() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.LBrace {
		return nil, b.AppendTokenToError("Could not get left brace of map")
	}

	// Increment over the left brace token
	b.Index++

	var (
		stmt  *Node
		stmts []*Node
		err   error
	)

	for b.Index < len(b.Tokens) && b.Tokens[b.Index].Type != token.RBrace {
		stmt, err = b.ParseExpression()
		if err != nil {
			return nil, err
		}

		blob, _ := json.Marshal(stmt)
		fmt.Println("kvstmtkv:", string(blob))

		// All statements in a map have to be key-value pairs
		if stmt.Type != "kv" {
			return nil, errors.Errorf("All statements in a map have to be key-value pairs: %+v\n", stmt)
		}

		fmt.Println("stmt:", stmt)

		stmts = append(stmts, stmt)
	}

	// Step over the right brace token
	b.Index++

	return &Node{
		Type:  "block",
		Value: stmts,
	}, nil
}

func (b *Builder) ParseEnumBlockStmt() (*Node, error) {
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
		return nil, b.AppendTokenToError("Could not get left brace of enum")
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

func (b *Builder) ParseBlockStmt() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.LBrace {
		return nil, b.AppendTokenToError("Could not get left brace of block")
	}

	// Increment over the left brace token
	b.Index++ // Create a new child scope for the function

	var (
		stmt  *Node
		stmts []*Node
		err   error
		rt    *ast.Type
	)

	for b.Index < len(b.Tokens) &&
		b.Tokens[b.Index].Type != token.RBrace {
		stmt, err = b.ParseStmt()
		if err != nil {
			if err.Error() == "ParseLiteralStatement not implemented for: R_BRACE" {
				break
			}

			return nil, err
		}

		// // If we are returning something from the block then we need to ground the type
		// if stmt.Type == "return" {
		// 	if stmt.Left != nil {
		// 		var t = stmt.Left.Kind

		// 		if stmt.Left.Type == "ident" {
		// 			// // Find the original type
		// 			tv := b.ScopeTree.Local(stmt.Left.Value.(string))
		// 			if tv == nil {
		// 				panic(stmt)
		// 				// return false, errors.Errorf("could not alias to unfound type: %s", n.Left.Value.(string))
		// 			}

		// 			t = tv.Right.Kind
		// 		}

		// 		rt = greatestCommonType(rt, ast.TypeFromString(t))
		// 		fmt.Println("rt:", rt)
		// 	}
		// }

		fmt.Println("i am here", stmt)

		stmts = append(stmts, stmt)
	}

	// Step over the right brace token
	b.Index++

	// TODO: set node.Metadata["returns"] stmt from the rt

	return &Node{
		Type:       "block",
		Value:      stmts,
		ReturnType: rt,
	}, nil
}

func (b *Builder) ParseInterfaceBlock(rcvrName string) (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.LBrace {
		return nil, b.AppendTokenToError("Could not get left brace of block")
	}

	// Increment over the left brace token
	b.Index++ // Create a new child scope for the function

	var (
		stmt  *Node
		stmts []*Node
		err   error
		rt    *ast.Type
	)

	for b.Index < len(b.Tokens) &&
		b.Tokens[b.Index].Type != token.RBrace {
		stmt, err = b.ParseFunctionPartialDecl(rcvrName)
		if err != nil {
			return nil, err
		}

		// err = b.ScopeTree.Declare(stmt)
		// if err != nil {
		// 	panic(fmt.Sprintf("wtf: %s", err))
		// }

		fmt.Println("i am here", stmt)

		stmts = append(stmts, stmt)
	}

	// Step over the right brace token
	b.Index++

	// TODO: set node.Metadata["returns"] stmt from the rt

	return &Node{
		Type:       "block",
		Value:      stmts,
		ReturnType: rt,
	}, nil
}

// func greatestCommonLT(rt, lt ast.LiteralType) (ast.LiteralType, bool) {
// 	if rt == lt {
// 		return rt, false
// 	}

// 	// TODO: need to do some upgrade map?
// }

func greatestCommonType(rt, lt *ast.Type) *ast.Type {
	if rt == nil {
		return lt
	}

	if lt == nil {
		panic("WTF THE LEFT TYPE PASSED IN WAS NIL")
	}

	if lt.Type == rt.Type {
		return lt
	}

	if rt.UpgradesTo != nil && lt.Type == rt.UpgradesTo.Type {
		return lt
	}

	// default:
	// 	// TODO: we need to change this to a Type so that we can do an upgradeable comparison
	// 	if rt.ShadowType != nil {
	// 		switch *rt.ShadowType {
	// 		case
	// 			lt.Type,
	// 			lt.UpgradesTo:

	// 		}
	// 	}

	if lt.UpgradesTo != nil {
		if lt.UpgradesTo.Type == rt.Type {
			return rt
		}

		if lt.UpgradesTo.Type == rt.UpgradesTo.Type {
			return lt.UpgradesTo
		}
	}

	// TODO: at some point the compiler needs to generate code to apply this at runtime
	panic("not sure if ast.NewVarType(ast.NoneType) will work")
	// return ast.NewVarType(ast.NoneType)
}

func (b *Builder) ParseReturnStmt() (*Node, error) {
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

// func (b *Builder) ParseDeclarationStmt(typeHint *TypeValue) (*Node, error) {
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

func (b *Builder) ParseTypeDeclStmt() (*Node, error) {
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

// THIS SHOULD NOT BE IN THE TRANSPILER; THIS SHOULD BE TAKEN CARE OF BY THE SEMANTIC STAGE
// BEFORE EVER REACHING THIS FAR. IT IS MERELY HERE AS A PLACEHOLDER SO I DO NOT FORGET
// MY OWN DECISIONS BECAUSE I CANT REMEMBER THINGS
func (b *Builder) ParseCBlock() (*Node, error) {
	var errStatement = "`c` blocks, as the are oh-so affectionately known within the Express community, are only implemented as a direct code injection  at time. This will take some thinking; the compiler will have to `back-compile` the C/C++ code FROM the AST output of Clang and then translate that back into Express code essentially to check it"
	// return nil, errors.New()
	log.Println(errStatement)

	// For now the C block will be a direct injection of code into the final source. This is the best we can get at this point

	// ADD THIS BACK IN
	// // Check ourselves ...
	// if b.Tokens[b.Index].Type != token.C {
	// 	return nil, b.AppendTokenToError("Could not get c block")
	// }

	// Skip over the `c` token
	b.Index++

	var (
		total []string
		found bool
	)

	// Gobble up all the code until the next left brace; use a simple array as a stack to know when we are done
	for _, t := range b.Tokens[b.Index:] {
		fmt.Println("t", t)

		if t.Type == token.RBrace {
			found = true
			break
		}

		// Append the string value of the token
		total = append(total, t.Value.String)

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

func (b *Builder) ParseObjectStmt() (*Node, error) {
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
	body, err := b.ParseBlockStmt()
	if err != nil {
		return nil, err
	}

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

func (b *Builder) ParseStructStmt() (*Node, error) {
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

	// // Check for the equals token
	// if b.Tokens[b.Index].Type != token.Assign {
	// 	return nil, b.AppendTokenToError("No equals found after ident in struct def")
	// }

	// // Increment over the equals
	// b.Index++

	// Parse the right hand side
	body, err := b.ParseBlockStmt()
	if err != nil {
		return nil, err
	}

	body.Kind = "struct"

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

func (b *Builder) ParseInterfaceStmt() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Interface {
		return nil, b.AppendTokenToError("Could not get struct declaration statement")
	}

	// Skip over the `interface` token
	b.Index++

	// Create the ident
	ident, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Increment over the ident token
	b.Index++

	// // Create a new child scope for the function
	// b.ScopeTree, err = b.ScopeTree.NewChildScope(ident.Value.(string))
	// if err != nil {
	// 	return nil, err
	// }

	// // Check for the equals token
	// if b.Tokens[b.Index].Type != token.Assign {
	// 	return nil, b.AppendTokenToError("No equals found after ident in struct def")
	// }

	// // Increment over the equals
	// b.Index++

	var v = &TypeValue{
		Composite: true,
		Type:      StruturedValue,
		Kind:      "interface",
		Props:     map[string]*TypeValue{},
	}

	err = b.ScopeTree.NewType(ident.Value.(string), v)
	if err != nil {
		return nil, err
	}

	// Parse the right hand side
	body, err := b.ParseInterfaceBlock(ident.Value.(string))
	if err != nil {
		return nil, err
	}

	body.Kind = "interface"

	// _, err = b.AddStructured(ident.Value.(string), body)
	// if err != nil {
	// 	return nil, err
	// }

	v.Props, err = b.extractPropsFromComposite(body)
	if err != nil {
		return nil, err
	}

	// // Increment over the first part of the expression
	// b.Index++

	// // Assign our scope back to the current one
	// b.ScopeTree, err = b.ScopeTree.Leave()
	// if err != nil {
	// 	return nil, err
	// }

	var n = Node{
		Type:  "interface",
		Left:  ident,
		Right: body,
	}

	return &n, nil
}

func (b *Builder) ParseMapStmt() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Map {
		return nil, b.AppendTokenToError("Could not get map declaration statement")
	}

	// Skip over the `map` token
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
		return nil, b.AppendTokenToError("No equals found after ident in map declaration")
	}

	// Increment over the equals
	b.Index++

	// Parse the right hand side
	body, err := b.ParseMapBlockStmt()
	if err != nil {
		return nil, err
	}

	body.Kind = "map"

	// Assign our scope back to the current one
	b.ScopeTree, err = b.ScopeTree.Leave()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:  "map",
		Left:  ident,
		Right: body,
	}, nil
}

func (b *Builder) ParseStructDeclarationStmt() (*Node, error) {
	return nil, errors.New("Not implemented: ParseStructDeclarationStatement")
}

func (b *Builder) ParseLetStmt() (*Node, error) {
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

func (b *Builder) ParseLiteralStmt() (*Node, error) {
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
	fmt.Println("leftblobby:", string(blob))

	b.Index++

	switch b.Tokens[b.Index].Type {
	case token.Set:
		return b.ParseSet(left)

	default:
		return nil, errors.Errorf("ParseLiteralStatement not implemented for: %+v", b.Tokens[b.Index].Type)
	}
}

// ParseIdentStmt: Although idents are not statements, they do start many statements
// and this function serves to disambiguate those statements
func (b *Builder) ParseIdentStmt() (*Node, error) {
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

	fmt.Println(b.Tokens[b.Index])

	// Parse the first ident; this COULD be a type
	identOrType, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Increment over the ident token
	b.Index++

	switch identOrType.Type {
	case "call":
		return identOrType, nil

	case "selection":
		if b.Tokens[b.Index].Type != token.Assign && b.Tokens[b.Index].Type != token.Set {
			return identOrType, nil
		}
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

	fmt.Println("identOrType, err", identOrType, err, b.Tokens[b.Index].Type)

	switch b.Tokens[b.Index].Type {
	case token.Ident:
		/*
			In this case, we have two idents back to back which leads us
			to make the only informed decision we can; that the first ident
			was a type, like in cases such as:
				int i = 0
		*/

		// Set the proper node values
		node.Type = "decl"
		node.Value = identOrType

		fmt.Println("got another ident", b.Tokens[b.Index], node)

		node.Left, err = b.ParseExpression()
		if err != nil {
			return nil, err
		}

		fmt.Println("node.Left", node.Left)

		if b.Index > len(b.Tokens)-1 && b.Tokens[b.Index+1].Type != token.Assign {
			fmt.Println("b.Tokens[b.Index+1]", b.Tokens[b.Index+1])
			return node, nil
		}

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
		fmt.Println("got assign")

		// Step over the assign
		b.Index++

		node.Right, err = b.ParseExpression()
		if err != nil {
			return nil, err
		}

		fmt.Println("node.Right:", node.Right)

		// TODO : scottshotgg : need to have d = h here as well; without a type
		if node.Value != nil {
			var nv = node.Value.(*Node)
			var kind = nv.Metadata["kind"]
			if nv.Type == "type" && kind != nil && kind.(string) == "interface" {
				// TODO: we have an interface assignment
				node.Metadata = map[string]interface{}{
					"isIfaceAssign": "true",
				}
			}
		}

		// TODO: here
		// b.Index++

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
		fmt.Println("blobby:", string(blob))

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
		fmt.Println("setblobby:", string(blob))

		return node, nil

	// Just return the ident if you don't know what to do
	// this will defer the judgement to the next statement up
	default:
		return identOrType, nil
	}

	// 	// If there is an ident after the ident, then we have what should be a type
	// 	// If there is assignment, then we have an assign statement

	// 	// Increment over the ident token
	// 	b.Index++

	// 	if b.Index > len(b.Tokens)-1 {
	// 		return ident, nil
	// 	}

	// 	if b.Tokens[b.Index].Type == token.Set {
	// 		return b.ParseSet(ident)
	// 	}

	// 	// Check for the equals token
	// 	if b.Tokens[b.Index].Type != token.Assign {
	// 		if ident.Type == "call" {
	// 			return ident, nil
	// 		}

	// 		// TODO: this is where we need to check for `:`

	// 		// return nil, b.AppendTokenToError(fmt.Sprintf("No equals found after ident in assignment: %+v", b.Tokens[b.Index]))
	// 		// This need to return the token in case the parse needs to be recovered! Look at ParseEnumBlock for an example of parse recovery
	// 		return ident, ErrNoEqualsFoundAfterIdent
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

	// 	var node = &Node{
	// 		Type:  "assignment",
	// 		Left:  ident,
	// 		Right: expr,
	// 	}

	// 	// Do one pass for declarations, and check that the assignments
	// 	// and usages corraborate in the type checker
	// 	// return node, b.ScopeTree.Assign(node)
	// 	return node, nil
	// }

	return nil, errors.Errorf("could not parse ident statement: %+v", b.Tokens[b.Index])
}

func (b *Builder) ParsePackageStmt() (*Node, error) {
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
		stmt, err := b.ParseStmt()
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
		fmt.Println("STMT", stmt)
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

func (b *Builder) ParseUseStmt() (*Node, error) {

	// TODO: add this back in
	// Check ourselves ...
	// if b.Tokens[b.Index].Type != token.Use {
	// 	return nil, b.AppendTokenToError("Could not get use statement")
	// }

	// Step over the import token
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
		// Just print out the entire expression for now
		return nil, errors.Errorf("Expecting \"as\" keyword after use expression, found: %+v", expr)
	}

	// Hop over the "as"
	b.Index++

	// Next up: we are expecting an _ident_; parse it as an expression so operation rules will apply
	// Not sure if that is needed (operation rules), but we'll see; could be a fun/fucky experiment
	expr1, err = b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// May have to mangle the names for this ;_; noooooo
	if expr1.Type != "ident" {
		// Just print out the entire expression for now
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
	// var path, err = os.Getwd()
	// if err != nil {
	// 	return nil, err
	// }

	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("source", string(source))

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
	b2 := New(tokens)
	ast, err := b2.BuildAST()
	if err != nil {
		return nil, nil, err
	}

	// fmt.Printf("ast %+v\n", ast.Value.([]*Node)[0].Left.Value.(string))

	// TODO: extremely unsafe, fix this
	return ast, b2.ScopeTree, nil
}

func (b *Builder) ParseImportStmt() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Import {
		return nil, b.AppendTokenToError("Could not get import statement")
	}

	// Step over the import token
	b.Index++

	expr, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Now that we have the expression, we need to go parse that file
	// 1. Parse the file
	// 2. Use a variable to link the file
	// 3. Normal selection checking after that
	// 4. Take special care for transpileImportStatement

	fmt.Println("expr.Kind", expr.Value.(string))

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

	if namespace[len(namespace)-5:] == ".expr" {
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

func (b *Builder) ParseIncludeStmt() (*Node, error) {
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

func (b *Builder) ParseThreadStmt() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Thread {
		return nil, b.AppendTokenToError("Could not get launch statement")
	}

	// Step over the import token
	b.Index++

	// Might need to make this an explicit function call later
	expr, err := b.ParseStmt()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type: "thread",
		Left: expr,
	}, nil
}

func (b *Builder) ParseFunctionStmt() (*Node, error) {
	// Check ourselves
	if b.Tokens[b.Index].Type != token.Function {
		return nil, b.AppendTokenToError("Could not get function")
	}

	// Step over the function token
	b.Index++

	return b.ParseFunctionPartialDecl("")
}

func (b *Builder) ParseFunctionPartialDecl(rcvrType string) (*Node, error) {
	var (
		err  error
		node = Node{
			Type: "function",
			Metadata: map[string]interface{}{
				"args": []*Node{},
			},
		}
	)

	// TODO: check for the square brackets here

	var rcvr *Node
	var isMethod bool
	var typeName = rcvrType

	// TODO: scottshogg : 04/16/23 : we do need to make an `isInterfaceMethod` ...
	if rcvrType != "" {
		isMethod = true
	}

	// Check if it is a method; func [Type Ident] ...
	if b.Tokens[b.Index].Type == token.LBracket {
		if rcvrType != "" {
			panic("WE HAVE AN ISSUE 12345")
		}

		// Step over the left bracket
		b.Index++

		rcvr, err = b.ParseIdentStmt()
		if err != nil {
			return nil, err
		}

		if b.Tokens[b.Index].Type == token.RBracket {
			b.Index++
		}

		var rcvrVal = rcvr.Value.(*Node)

		if rcvrVal.Type == "deref" {
			// if rcvrVal.Kind == "pointer" {
			typeName = rcvrVal.Left.Value.(string)
			rcvrType = "*" + typeName
		} else {
			typeName = rcvr.Value.(*Node).Value.(string)
			rcvrType = typeName
		}

		isMethod = true
	}

	// Named function
	if b.Tokens[b.Index].Type != token.Ident {
		return nil, b.AppendTokenToError("Could not get ident after function token")
	}

	// Set the name of the function
	var kind = b.Tokens[b.Index].Value.String
	node.Kind = kind

	if rcvrType != "" {
		node.Kind = fmt.Sprintf("%s.%s",
			rcvrType,
			kind,
		)
	}

	// Create a new child scope for the function
	b.ScopeTree, err = b.ScopeTree.NewChildScope(node.Kind)
	if err != nil {
		return nil, err
	}

	// Step over the ident token
	b.Index++

	if b.Tokens[b.Index].Type != token.LParen {
		return nil, b.AppendTokenToError("Could not get left paren")
	}

	args, err := b.ParseGroupOfStatements()
	if err != nil {
		return nil, err
	}

	if args != nil {
		if rcvr != nil {
			var argsValue = args.Value.([]*Node)
			args.Value = append([]*Node{rcvr}, argsValue...)
		}

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

		} else if b.Tokens[b.Index].Type == token.Type {
			// Make an egroup with one return in it
			node.Metadata["returns"] = &Node{
				Type: "egroup",
				Value: []*Node{
					{
						Type:  "type",
						Value: b.Tokens[b.Index].Value.String,
					},
				},
			}

			// Step over the type token
			b.Index++
		} // else we have a _function header_

		// else {
		// 	return nil, errors.Errorf("could not parse returns on %s: %v", node.Kind, b.Tokens[b.Index])
		// }
	}

	fmt.Println("node.Metadata[returns]:", node.Metadata["returns"])

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

	if b.Tokens[b.Index].Type == token.LBrace {
		fmt.Println("got a function body")

		bs, err := b.ParseBlockStmt()
		if err != nil {
			return nil, err
		}

		// // TODO: probably need to do a more concerned check here later
		// if node.Metadata["returns"] == nil && bs.ReturnType != nil {
		// 	node.Metadata["returns"] = &Node{
		// 		Type:  "type",
		// 		Value: bs.Kind,
		// 	}
		// }

		node.Value = bs
	}

	// node.Value = addDeferDeclarationToBlock(block)

	// Assign our scope back to the current one
	b.ScopeTree, err = b.ScopeTree.Leave()
	if err != nil {
		return nil, err
	}

	// scottshotgg : 3/27/23 : if it's a method then declare it on the struct
	if isMethod {
		var n = b.ScopeTree.GetType(typeName)
		if n == nil {
			return nil, errors.Errorf("n == nil: %s %+v", typeName, n)
		}

		n.Props[kind] = &TypeValue{
			Type:  FunctionValue,
			Value: &node,
		}

		/*
			NOTE: scottshotgg : 04/16/23 :
				map is going to be a pointer so we don't need to do
				any re-assignment here
		*/
		// b.ScopeTree.Types[typeName] = n
	}

	// Declare the type in the upper scope after leaving
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

func (b *Builder) ParseDerefStmt() (*Node, error) {
	if b.Tokens[b.Index].Type != token.PriOp {
		return nil, b.AppendTokenToError("Could not get deref statement without *")
	}

	deref, err := b.ParseDerefExpression()
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

	// TODO: scottshotgg : 10.16.23 : this is where we need to pick up from
	// swapping the pointers around is confusing the compiler; pointer vs deref

	return &Node{
		Type:  t,
		Left:  deref,
		Right: expr,
	}, nil
}

func (b *Builder) ParseArrayDeclStmt() (*Node, error) {
	// Check ourselves ...
	if b.Index < len(b.Tokens)-1 &&
		b.Tokens[b.Index].Type != token.LBracket {
		return nil, b.AppendTokenToError("Could not get type declaration statement")
	}

	// Parse the initial array expression declaration
	arrDef, err := b.ParseArrayExpression()
	if err != nil {
		return nil, err
	}

	// Increment over the R_BRACKET
	b.Index++

	// Parse the non-repeated type
	typeOf, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Increment over the type
	b.Index++

	// Parse the ident
	ident, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Increment over the ident
	b.Index++

	// Check for the equals token
	if b.Tokens[b.Index].Type != token.Assign {
		return nil, b.AppendTokenToError("No equals found after ident in array decl; required for now")
	}

	// Increment over the equals
	b.Index++

	// Parse the actual array expression
	arrExp, err := b.ParseArrayExpression()
	if err != nil {
		return nil, err
	}

	var specificType = deriveTypeFromArray(arrExp)
	_ = specificType

	// Increment over the R_BRACKET
	b.Index++

	return &Node{
		Type: "array_decl",
		// Kind: // TODO: kind should be the ground-state of the array
		Left:  typeOf,
		Right: ident,
		Value: arrExp,
		Metadata: map[string]interface{}{
			"def":           arrDef,
			"specific_type": "",
		},
	}, nil
}

type SpecificType struct {
	Type string
	Kind []SpecificType
}

/*
{
	"type": "array",
	"kind": [
		{
			"type": "int"
		},
		{
			"type": "int"
		}
		]
}
*/

func deriveTypeFromArray(n *Node) *SpecificType {
	if n.Type != "array" {
		return nil
	}

	if n.Value == nil {
		return nil
	}

	var nodes, ok = n.Value.([]*Node)
	if !ok {
		panic("wtf not an array")
	}

	var st SpecificType

	for _, node := range nodes {
		switch node.Type {
		// We have a map
		case "block":
			for _, node1 := range node.Value.([]*Node) {
				fmt.Println("node1:", node1)
				if node1.Type != "kv" {
					panic("node type is not a key-value pair for map")
				}

			}
		}
	}

	return &st
}

// TODO: what if types were expressions ...

// ParseStatement ** does not ** look ahead
func (b *Builder) ParseStmt() (*Node, error) {
	switch b.Tokens[b.Index].Type {
	case token.LBracket:
		return b.ParseArrayDeclStmt()

	case token.Thread:
		return b.ParseThreadStmt()

	case token.Defer:
		return b.ParseDeferStmt()

	case token.Enum:
		return b.ParseEnumBlockStmt()

	case token.Map:
		return b.ParseMapStmt()

	case token.PriOp:
		return b.ParseDerefStmt()

	case token.Package:
		return b.ParsePackageStmt()

	case token.Import:
		return b.ParseImportStmt()

	// case token.Use:
	// 	return b.ParseUseStmt()

	case token.Include:
		return b.ParseIncludeStmt()

	case token.TypeDef:
		return b.ParseTypeDeclStmt()

	case token.Struct:
		return b.ParseStructStmt()

	case token.Interface:
		return b.ParseInterfaceStmt()

	case token.Object:
		return b.ParseObjectStmt()

	case token.C:
		return b.ParseCBlock()

	// For literal and idents, we will need to figure out what
	// kind of statement it is
	case token.Literal:
		return b.ParseLiteralStmt()

	case token.Ident:
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

		var n, err = b.ParseIdentStmt()
		if err != nil {
			return nil, err
		}

		if n.Type == "decl" {
			blob, _ := json.Marshal(n)
			fmt.Println("nblobn:", string(blob))

			err = b.ScopeTree.Declare(n)
			if err != nil {
				return nil, err
			}
			blob, _ = json.Marshal(b.ScopeTree)
			fmt.Println("ScopeTree:", string(blob))
		}

		return n, nil

	case token.Type:
		var n, err = b.ParseIdentStmt()
		if err != nil {
			return nil, err
		}

		err = b.ScopeTree.Declare(n)
		if err != nil {
			return nil, err
		}

		return n, nil

	case token.Function:
		return b.ParseFunctionStmt()

	case token.LBrace:
		return b.ParseBlockStmt()

	case token.Let:
		return b.ParseLetStmt()

	case token.If:
		return b.ParseIfStmt()

	case token.For:
		return b.ParseForStmt()

	case token.Return:
		return b.ParseReturnStmt()

	case token.Link:
		return b.ParseLinkStmt()
	}

	return nil, b.AppendTokenToError(fmt.Sprintf("Could not create statement from: %+v", b.Tokens[b.Index].Type))
}

func (b *Builder) ParseLinkStmt() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Link {
		return nil, b.AppendTokenToError("Could not create link statement")
	}

	// Skip over the `link` token
	b.Index++

	// TODO: this should actually dump the function headers into the current scope
	// Create the ident
	ident, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Increment over the ident token
	b.Index++

	// Create the ident
	block, err := b.ParseIdentStmt()
	if err != nil {
		return nil, err
	}

	// // Create a new child scope for the function
	// b.ScopeTree, err = b.ScopeTree.NewChildScope(ident.Value.(string))
	// if err != nil {
	// 	return nil, err
	// }

	// // Check for the equals token
	// if b.Tokens[b.Index].Type != token.Assign {
	// 	return nil, b.AppendTokenToError("No equals found after ident in map declaration")
	// }

	// // Increment over the equals
	// b.Index++

	// // Parse the right hand side
	// body, err := b.ParseMapBlockStmt()
	// if err != nil {
	// 	return nil, err
	// }

	// body.Kind = "link"

	// // Assign our scope back to the current one
	// b.ScopeTree, err = b.ScopeTree.Leave()
	// if err != nil {
	// 	return nil, err
	// }

	for _, kv := range block.Value.([]*Node) {
		// Function header to link TO
		var funcHeader = b.ScopeTree.Get(kv.Right.Value.(string))
		funcHeader.Value = &Node{
			Type: "block",
			Value: []*Node{
				{
					Type: "return",
					Left: &Node{
						Type: "call",
						Metadata: map[string]interface{}{
							"args": &Node{
								Type:  "egroup",
								Value: convertArgsLibC(funcHeader),
							},
						},
						Value: &Node{
							Type: "ident",
							// Function from library to link FROM
							Value: kv.Left.Value.(string),
						},
					},
				},
			},
		}
	}

	return &Node{
		Type:  "link",
		Left:  ident,
		Right: block,
	}, nil
}

// func EGroupFromSGroup() {

// }

func convertArgsLibC(fh *Node) []*Node {
	var libcArgs []*Node

	var args = fh.Metadata["args"]
	for _, v := range args.(*Node).Value.([]*Node) {
		if v.Value.(*Node).Kind == "string" {
			libcArgs = append(libcArgs, &Node{
				Type: "selection",
				Left: v.Left,
				Right: &Node{
					Type: "call",
					Value: &Node{
						Type:  "ident",
						Value: "c_str",
						Metadata: map[string]interface{}{
							"args": &Node{
								Type:  "egroup",
								Value: []*Node{},
							},
						},
					},
				},
			})

			continue
		}

		libcArgs = append(libcArgs, v)
	}

	return libcArgs
}
