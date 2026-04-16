package builder

import (
	ast "github.com/scottshotgg/express-ast"
	fmt "fmt"
	token "github.com/scottshotgg/express-token"

	"github.com/pkg/errors"
)

var (
// primTypes = []string{
// 	"int",
// 	"float",
// 	"bool",
// 	"char",
// 	"byte",
// 	"string",
// }
)

func New(tokens []token.Token) *Builder {
	// Reset global scope tree before each compilation
	scopeTree = nil
	currentTree = nil

	fmt.Println("DEBUG: New() called, scopeTree reset")

	var b = Builder{
		Tokens:    tokens,
		ScopeTree: NewScopeTree(),
		// TypeMap: map[string]*TypeValue{},
	}

	fmt.Println("DEBUG: New() created ScopeTree, Children:", getKeys(b.ScopeTree.Children))

	// // Load the concrete types; we'll need to do more
	// // later if we have separate transpilers per file ...
	// for _, value := range primTypes {
	// 	b.TypeMap[value] = &TypeValue{
	// 		Type: PrimitiveValue,
	// 		Kind: value,
	// 	}
	// }

	b.registerParseFns()

	return &b
}

func (b *Builder) BuildAST() (*Node, error) {
	var (
		stmt        *Node
		stmts       []*Node
		err         error
		programNode = &Node{
			Type: "program",
		}
	)

	b.ScopeTree = NewScopeTree()
	fmt.Println("DEBUG: BuildAST() created ScopeTree, Children:", getKeys(b.ScopeTree.Children))

	// First pass: process typedefs to populate the scope tree
	// This allows structs and other statements to reference type aliases
	for b.Index < len(b.Tokens)-1 {
		// Peek at the next token to check if it's a typedef
		if b.Tokens[b.Index].Type == token.TypeDef {
			stmt, err = b.ParseTypeDeclarationStatement()
			if err != nil {
				return nil, err
			}

			// Process the typedef immediately to add it to the scope tree
			if err = b.ProcessTypeDeclaration(stmt); err != nil {
				return nil, err
			}
		} else {
			// Not a typedef, break out to normal parsing
			break
		}
	}

	// Second pass: parse remaining statements
	b.Index = 0 // Reset index and re-parse all tokens
	b.Tokens, err = ast.CompressTokens(b.Tokens) // Re-compress since we reset
	if err != nil {
		return nil, err
	}

	// Actually, let's use a different approach: parse all statements, but
	// process typedefs as we go. For each typedef, add it to the scope tree
	// immediately. Then when parsing structs, the types will be available.

	// Reset and parse everything in one pass, processing typedefs immediately
	b.Index = 0
	stmts = nil

	for b.Index < len(b.Tokens)-1 {
		stmt, err = b.ParseStatement()
		if err != nil {
			if err == ErrOutOfTokens {
				break
			}

			return nil, err
		}

		if stmt == nil {
			return nil, b.AppendTokenToError("Statement was nil")
		}

		// If this is a typedef, process it immediately to add to scope tree
		if stmt.Type == "typedef" {
			if err = b.ProcessTypeDeclaration(stmt); err != nil {
				return nil, err
			}
		}

		stmts = append(stmts, stmt)
	}

	programNode.Value = stmts

	return programNode, nil
}

// ProcessTypeDeclaration processes a typedef statement and adds the alias
// to the scope tree immediately, so it's available for subsequent statements.
func (b *Builder) ProcessTypeDeclaration(decl *Node) error {
	// decl is of the form: {Type: "typedef", Left: ident, Right: type_expr}
	// We need to look up the right-hand side type and add the alias

	var typeExpr = decl.Right
	var aliasName = decl.Left.Value.(string)

	// Find the type that we're aliasing
	var tv *TypeValue
	switch typeExpr.Type {
	case "ident":
		// Looking up a simple type like "int" or "float"
		tv = b.ScopeTree.GetType(typeExpr.Value.(string))
		if tv == nil {
			return errors.Errorf("could not find type to alias: %s", typeExpr.Value.(string))
		}

	case "type":
		// Already a type node, use its Kind
		tv = b.ScopeTree.GetType(typeExpr.Kind)
		if tv == nil {
			return errors.Errorf("could not find type: %s", typeExpr.Kind)
		}

	default:
		return errors.Errorf("unsupported type expression in typedef: %s", typeExpr.Type)
	}

	// Add the alias to the scope tree
	return b.ScopeTree.NewType(aliasName, tv)
}
