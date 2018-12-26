package transpiler_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/pkg/errors"
	"github.com/scottshotgg/express-ast"
	"github.com/scottshotgg/express-lex"
	"github.com/scottshotgg/express-token"
	"github.com/scottshotgg/express2/builder"
	"github.com/scottshotgg/express2/test"
	"github.com/scottshotgg/express2/transpiler"
)

var (
	cpp *string
	err error
)

func getTokensFromString(s string) ([]token.Token, error) {
	// Lex and tokenize the source code
	tokens, err := lex.New(s).Lex()
	if err != nil {
		return nil, err
	}

	fmt.Println("\nCompressing tokens ...")

	// Compress certain tokens;
	// i.e: `:` and `=` compress into `:=`
	return ast.CompressTokens(tokens)
}

func getBuilderFromString(test string) (*builder.Builder, error) {
	tokens, err := getTokensFromString(test)
	if err != nil {
		return nil, err
	}

	for _, token := range tokens {
		fmt.Println(token)
	}

	return builder.New(tokens), nil
}

func getASTFromString(test string) (*builder.Node, error) {
	b, err := getBuilderFromString(test)
	if err != nil {
		return nil, err
	}

	node, err := b.BuildAST()

	fmt.Println("scopeTree", b.ScopeTree)

	return node, err
}

func getTranspilerFromString(test, name string) (*transpiler.Transpiler, error) {
	ast, err := getASTFromString(test)
	if err != nil {
		return nil, errors.Errorf("Could not create AST: %+v", err)
	}

	return transpiler.New(ast, name), nil
}

func getExpressionASTFromString(test string) (*builder.Node, error) {
	b, err := getBuilderFromString(test)
	if err != nil {
		return nil, err
	}

	return b.ParseExpression()
}

func getTypeASTFromString(test string) (*builder.Node, error) {
	b, err := getBuilderFromString(test)
	if err != nil {
		return nil, err
	}

	return b.ParseType()
}

func getStatementASTFromString(test string) (*builder.Node, error) {
	b, err := getBuilderFromString(test)
	if err != nil {
		return nil, err
	}

	return b.ParseStatement()
}

// func getTranspilerFromFilename(filename string) {
// 	testBytes, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		return nil, errors.Errorf("Could not read file: %+v", err)
// 	}

// 	test := string(testBytes)

// 	ast, err := getASTFromString(test)
// 	if err != nil {
// 		return nil, errors.Errorf("Could not create AST: %+v", err)
// 	}

// 	t, err := getTranspilerFromString(test)
// 	if err != nil {
// 		return nil, err
// 	}

// }

// TODO: do this; gonna have to make a cross-testing package or something
// func TestTranspileExpression(t *testing.T) {
// 	test :=
// }

func TestTranspiler(t *testing.T) {
	testBytes, err := ioutil.ReadFile("test.expr")
	if err != nil {
		t.Fatalf("Could not read file: %+v", err)
	}

	var test = string(testBytes)

	tr, err := getTranspilerFromString(test, "main")
	if err != nil {
		t.Fatalf("Could not create transpiler: %+v", err)
	}

	cpp, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Could not transpile to C++: %+v", err)
	}

	fmt.Printf("\nC++: %s\n\n", cpp)
}

func TestTranspileIdentExpression(t *testing.T) {
	ast, err := getExpressionASTFromString(test.Tests[test.ExpressionTest]["ident"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err = transpiler.TranspileIdentExpression(ast)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileLiteralExpression(t *testing.T) {
	ast, err := getExpressionASTFromString(test.Tests[test.ExpressionTest]["intLit"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err = transpiler.TranspileLiteralExpression(ast)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

// ALL expressions are the same when transpiled except for arrays,
// some selections (implicit derefs), and blocks

func TestTranspileArrayExpression(t *testing.T) {
	ast, err := getExpressionASTFromString(test.Tests[test.ExpressionTest]["intLitArray"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err = transpiler.TranspileArrayExpression(ast)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileType(t *testing.T) {
	testt := "float"

	ast, err := getTypeASTFromString(testt)
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err = transpiler.TranspileType(ast)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileTypeDeclaration(t *testing.T) {
	var ast, err = getStatementASTFromString(test.Tests[test.StatementTest]["typeDef"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err = transpiler.TranspileTypeDeclaration(ast)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileAssignmentStatement(t *testing.T) {
	ast, err := getStatementASTFromString(test.Tests[test.StatementTest]["simpleAssign"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err = transpiler.TranspileAssignmentStatement(ast)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileDeclarationStatement(t *testing.T) {
	ast, err := getStatementASTFromString(test.Tests[test.StatementTest]["decl"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	fmt.Printf("ast: %+v\n", ast.Right)

	cpp, err = transpiler.TranspileDeclarationStatement(ast)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileIncrementExpression(t *testing.T) {
	ast, err := getStatementASTFromString(test.Tests[test.ExpressionTest]["inc"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err = transpiler.TranspileIncrementExpression(ast)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileConditionExpression(t *testing.T) {
	ast, err := getStatementASTFromString(test.Tests[test.ExpressionTest]["condition"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err = transpiler.TranspileConditionExpression(ast)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileBinOpExpression(t *testing.T) {
	ast, err := getExpressionASTFromString(test.Tests[test.ExpressionTest]["binop"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err = transpiler.TranspileBinOpExpression(ast)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileBlockStatement(t *testing.T) {
	ast, err := getStatementASTFromString(test.Tests[test.StatementTest]["block"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err = transpiler.TranspileBlockStatement(ast)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileForInStatement(t *testing.T) {
	ast, err := getStatementASTFromString(test.Tests[test.StatementTest]["forin"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err = transpiler.TranspileForInStatement(ast)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileCallExpression(t *testing.T) {
	ast, err := getStatementASTFromString(test.Tests[test.ExpressionTest]["identCall"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	fmt.Printf("ast %+v\n", ast)

	cpp, err = transpiler.TranspileCallExpression(ast)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

// this dont work because sgroup is not in the ParseStatement switch
func TestTranspileSGroup(t *testing.T) {
	ast, err := getStatementASTFromString(test.Tests[test.StatementTest]["sgroup"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err = transpiler.TranspileSGroup(ast)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileFunctionStatement(t *testing.T) {
	ast, err := getStatementASTFromString(test.Tests[test.StatementTest]["funcDef"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err = transpiler.TranspileFunctionStatement(ast)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileIndexExpression(t *testing.T) {
	ast, err := getStatementASTFromString(test.Tests[test.ExpressionTest]["identIndex"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	fmt.Printf("ast %+v\n", ast)

	cpp, err = transpiler.TranspileIndexExpression(ast)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileImportStatement(t *testing.T) {
	ast, err := getStatementASTFromString(test.Tests[test.ExpressionTest]["import"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	fmt.Printf("ast %+v\n", ast)

	cpp, err = transpiler.TranspileImportStatement(ast)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileIncludeStatement(t *testing.T) {
	ast, err := getStatementASTFromString(test.Tests[test.StatementTest]["include"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	fmt.Printf("ast %+v\n", ast)

	cpp, err = transpiler.TranspileIncludeStatement(ast)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}
