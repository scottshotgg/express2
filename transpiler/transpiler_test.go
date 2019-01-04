package transpiler_test

import (
	"fmt"
	"io/ioutil"
	"testing"

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
	var b, err = getBuilderFromString(test)
	if err != nil {
		return nil, err
	}

	ast, err := b.BuildAST()
	if err != nil {
		return nil, err
	}

	return transpiler.New(ast, b, name), nil
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

func getStatementTranspilerFromString(test string) (*transpiler.Transpiler, error) {
	var b, err = getBuilderFromString(test)
	if err != nil {
		return nil, err
	}

	ast, err := b.ParseStatement()
	if err != nil {
		return nil, err
	}

	return transpiler.New(ast, b, "main"), nil
}

func getExpressionTranspilerFromString(test string) (*transpiler.Transpiler, error) {
	var b, err = getBuilderFromString(test)
	if err != nil {
		return nil, err
	}

	ast, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	return transpiler.New(ast, b, "main"), nil
}

func getTypeTranspilerFromString(test string) (*transpiler.Transpiler, error) {
	var b, err = getBuilderFromString(test)
	if err != nil {
		return nil, err
	}

	ast, err := b.ParseType()
	if err != nil {
		return nil, err
	}

	return transpiler.New(ast, b, "main"), nil
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
	tr, err := getExpressionTranspilerFromString(test.Tests[test.ExpressionTest]["ident"])
	if err != nil {
		t.Fatalf("Could not create Transpiler: %+v", err)
	}

	cpp, err := tr.TranspileIdentExpression(tr.AST)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileLiteralExpression(t *testing.T) {
	tr, err := getExpressionTranspilerFromString(test.Tests[test.ExpressionTest]["intLit"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err := tr.TranspileLiteralExpression(tr.AST)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

// ALL expressions are the same when transpiled except for arrays,
// some selections (implicit derefs), and blocks

func TestTranspileArrayExpression(t *testing.T) {
	tr, err := getExpressionTranspilerFromString(test.Tests[test.ExpressionTest]["intLitArray"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err := tr.TranspileArrayExpression(tr.AST)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileType(t *testing.T) {
	testt := "float"

	tr, err := getTypeTranspilerFromString(testt)
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err := tr.TranspileType(tr.AST)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileTypeDeclarationStatement(t *testing.T) {
	var tr, err = getStatementTranspilerFromString(test.Tests[test.StatementTest]["typeDef"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err := tr.TranspileTypeDeclaration(tr.AST)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileAssignmentStatement(t *testing.T) {
	var tr, err = getStatementTranspilerFromString(test.Tests[test.StatementTest]["simpleAssign"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err := tr.TranspileAssignmentStatement(tr.AST)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileDeclarationStatement(t *testing.T) {
	var tr, err = getStatementTranspilerFromString(test.Tests[test.StatementTest]["decl"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err := tr.TranspileDeclarationStatement(tr.AST)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileStructDeclarationStatement(t *testing.T) {
	var tr, err = getStatementTranspilerFromString(test.Tests[test.StatementTest]["struct"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err := tr.TranspileStructDeclaration(tr.AST)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileIncrementExpression(t *testing.T) {
	var tr, err = getStatementTranspilerFromString(test.Tests[test.ExpressionTest]["inc"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err := tr.TranspileIncrementExpression(tr.AST)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileConditionExpression(t *testing.T) {
	var tr, err = getStatementTranspilerFromString(test.Tests[test.ExpressionTest]["condition"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err := tr.TranspileConditionExpression(tr.AST)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileBinOpExpression(t *testing.T) {
	tr, err := getExpressionTranspilerFromString(test.Tests[test.ExpressionTest]["binop"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err := tr.TranspileBinOpExpression(tr.AST)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileBlockStatement(t *testing.T) {
	var tr, err = getStatementTranspilerFromString(test.Tests[test.StatementTest]["block"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err := tr.TranspileBlockStatement(tr.AST)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileForInStatement(t *testing.T) {
	var tr, err = getStatementTranspilerFromString(test.Tests[test.StatementTest]["forin"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err := tr.TranspileForInStatement(tr.AST)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileCallExpression(t *testing.T) {
	var tr, err = getStatementTranspilerFromString(test.Tests[test.ExpressionTest]["identCall"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err := tr.TranspileCallExpression(tr.AST)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

// this dont work because sgroup is not in the ParseStatement switch
func TestTranspileSGroup(t *testing.T) {
	var tr, err = getStatementTranspilerFromString(test.Tests[test.StatementTest]["sgroup"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err := tr.TranspileSGroup(tr.AST)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileFunctionStatement(t *testing.T) {
	var tr, err = getStatementTranspilerFromString(test.Tests[test.StatementTest]["funcDef"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err := tr.TranspileFunctionStatement(tr.AST)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileIndexExpression(t *testing.T) {
	var tr, err = getStatementTranspilerFromString(test.Tests[test.ExpressionTest]["identIndex"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err := tr.TranspileIndexExpression(tr.AST)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileImportStatement(t *testing.T) {
	var tr, err = getStatementTranspilerFromString(test.Tests[test.ExpressionTest]["import"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err := tr.TranspileImportStatement(tr.AST)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}

func TestTranspileIncludeStatement(t *testing.T) {
	var tr, err = getStatementTranspilerFromString(test.Tests[test.StatementTest]["include"])
	if err != nil {
		t.Fatalf("Could not create AST: %+v", err)
	}

	cpp, err := tr.TranspileIncludeStatement(tr.AST)
	if err != nil {
		t.Errorf("err: %+v", err)
	}

	fmt.Println("C++:", *cpp)
}
