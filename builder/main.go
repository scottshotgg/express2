package builder

import (
	ast "github.com/scottshotgg/express-ast"
	token "github.com/scottshotgg/express-token"

	"github.com/pkg/errors"
	"github.com/scottshotgg/express2/pkg/logger"
)

func New(tokens []token.Token, log logger.Logger) *Builder {
	// Reset global scope tree before each compilation
	scopeTree = nil
	currentTree = nil

	if log == nil {
		log = logger.Noop()
	}

	b := &Builder{
		Tokens:    tokens,
		ScopeTree: NewScopeTree(),
		log:       log,
	}

	b.registerParseFns()

	return b
}

func (b *Builder) BuildAST() (*Node, error) {
	var (
		stmt  *Node
		stmts []*Node
		err   error
	)

	// ScopeTree is already initialized by New(); do NOT reset it here.

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

		// Process typedefs immediately so subsequent statements can reference them.
		if stmt.Type == "typedef" {
			if err = b.ProcessTypeDeclaration(stmt); err != nil {
				return nil, err
			}
		}

		stmts = append(stmts, stmt)
	}

	return &Node{Type: "program", Value: stmts}, nil
}

// ProcessTypeDeclaration processes a typedef statement and adds the alias
// to the scope tree immediately, so it's available for subsequent statements.
func (b *Builder) ProcessTypeDeclaration(decl *Node) error {
	var typeExpr = decl.Right
	var aliasName = decl.Left.Value.(string)

	var tv *TypeValue
	switch typeExpr.Type {
	case "ident":
		tv = b.ScopeTree.GetType(typeExpr.Value.(string))
		if tv == nil {
			return errors.Errorf("could not find type to alias: %s", typeExpr.Value.(string))
		}

	case "type":
		tv = b.ScopeTree.GetType(typeExpr.Kind)
		if tv == nil {
			return errors.Errorf("could not find type: %s", typeExpr.Kind)
		}

	default:
		return errors.Errorf("unsupported type expression in typedef: %s", typeExpr.Type)
	}

	return b.ScopeTree.NewType(aliasName, tv)
}

// CompressTokens is a convenience wrapper around the ast package function.
func CompressTokens(tokens []token.Token) ([]token.Token, error) {
	return ast.CompressTokens(tokens)
}
