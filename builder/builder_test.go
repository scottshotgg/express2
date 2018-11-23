package builder_test

import (
	"fmt"
	"testing"

	"github.com/scottshotgg/express-ast"
	"github.com/scottshotgg/express-lex"
	"github.com/scottshotgg/express-token"
	"github.com/scottshotgg/express2/builder"
)

func getTokensFromString(s string) ([]token.Token, error) {
	// Lex and tokenize the source code
	tokens, err := lex.New(s).Lex()
	if err != nil {
		fmt.Println("err", err)
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

func printTokensFromBuilder(b *builder.Builder) {
	for _, token := range b.Tokens {
		fmt.Println(token)
	}
}

func TestParseTerm(t *testing.T) {
	test := "i++"

	b, err := getBuilderFromString(test)
	if err != nil {
		fmt.Println("err", err)
		t.Fatal()
	}

	programAST, err := b.ParseExpression()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	fmt.Printf("programAST %+v\n", programAST)
}

func TestParseExpression(t *testing.T) {
	test := "thisIsAnIdent"

	tokens, err := getTokensFromString(test)
	if err != nil {
		fmt.Println("err", err)
		t.Fatal()
	}

	programAST, err := builder.New(tokens).ParseExpression()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	fmt.Printf("programAST %+v\n", programAST)
}

func TestParseDeclarationStatement(t *testing.T) {
	// TODO: we need the rest of the declaration types and stuff
	test := "int i = 10"

	tokens, err := getTokensFromString(test)
	if err != nil {
		fmt.Println("err", err)
		t.Fatal()
	}

	programAST, err := builder.New(tokens).ParseDeclarationStatement()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	fmt.Printf("programAST %+v\n", programAST)
}

func TestParseAssignmentStatement(t *testing.T) {
	test := "i = 10"

	tokens, err := getTokensFromString(test)
	if err != nil {
		fmt.Println("err", err)
		t.Fatal()
	}

	programAST, err := builder.New(tokens).ParseAssignmentStatement()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	fmt.Printf("programAST %+v\n", programAST)
}

func TestIfElseStatement(t *testing.T) {
	test := "if something { } else if somethingElse { }"

	b, err := getBuilderFromString(test)
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseIfStatement()
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	fmt.Printf("programAST %+v\n", programAST)
}

func TestParseGroupOfStatements(t *testing.T) {
	test := "(int i, string s)"

	b, err := getBuilderFromString(test)
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseGroupOfStatements()
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	fmt.Printf("\nprogramAST %+v\n", programAST)
}

func TestParseFunctionStatement(t *testing.T) {
	test := "function something(int i, string s) { float f = 10.1 }"

	b, err := getBuilderFromString(test)
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseFunctionStatement()
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	fmt.Printf("\nprogramAST %+v\n", programAST)
}

func TestParseGroupOfExpressions(t *testing.T) {
	test := "(1, i, s, 9)"

	b, err := getBuilderFromString(test)
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseGroupOfExpressions()
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	fmt.Printf("\nprogramAST %+v\n", programAST)
}

func TestParseCallStatement(t *testing.T) {
	// TODO: put all strings into a map[string]string
	// and pull them out so that other tests can use them,
	// like ParseStatement
	test := "something(i, s)"

	b, err := getBuilderFromString(test)
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseCall()
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	fmt.Printf("\nprogramAST %+v\n", programAST)
}

// TODO: this is an expression too
func TestParseBlockStatement(t *testing.T) {
	test := "{ int i = 10 int j = 99 }"

	b, err := getBuilderFromString(test)
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseBlockStatement()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	fmt.Printf("programAST %+v\n", programAST)
}

func TestParseImportStatement(t *testing.T) {
	test := "import \"somethingHere.expr\""

	b, err := getBuilderFromString(test)
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseImportStatement()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	fmt.Printf("programAST %+v\n", programAST)
}

func TestParseIncludeStatement(t *testing.T) {
	test := "include \"somethingHere.expr\""

	b, err := getBuilderFromString(test)
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseIncludeStatement()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	fmt.Printf("programAST %+v\n", programAST)
}

func TestParseConditionExpression(t *testing.T) {
	test := "something < 10 < (7)"

	b, err := getBuilderFromString(test)
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseExpression()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	fmt.Printf("programAST %+v\n", programAST)
}

func TestParseIncrementExpression(t *testing.T) {
	test := "i++"

	b, err := getBuilderFromString(test)
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseExpression()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	// fmt.Println(b.ParseBlockStatement())

	fmt.Printf("programAST %+v\n", programAST)
}

func TestParseForStdStatement(t *testing.T) {
	test := "for int i = 1; i < 10; i++ { }"

	b, err := getBuilderFromString(test)
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseForStdStatement()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	fmt.Printf("programAST %+v\n", programAST)
}

func TestParseArrayExpression(t *testing.T) {
	test := "[ \"something\", [8, 8], 9, i ]"

	b, err := getBuilderFromString(test)
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseArrayExpression()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	fmt.Printf("programAST %+v\n", programAST)
}

func TestParseType(t *testing.T) {
	test := "int[][5]"

	b, err := getBuilderFromString(test)
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseType()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	fmt.Printf("programAST %+v\n", programAST)
}

func TestParseArrayDeclaration(t *testing.T) {
	test := "int[] i = [ 8, 9, 0 ]"

	b, err := getBuilderFromString(test)
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseStatement()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	fmt.Printf("programAST %+v\n", programAST)
}

func TestParseForInStatement(t *testing.T) {
	test := "for i in [ 7, 8, 9 ] { }"

	b, err := getBuilderFromString(test)
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseForPrepositionStatement()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	fmt.Printf("programAST %+v\n", programAST)
}

func TestParseForOfStatement(t *testing.T) {
	test := "for i of [ 7, 8, 9 ] { }"

	b, err := getBuilderFromString(test)
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseForPrepositionStatement()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	fmt.Printf("programAST %+v\n", programAST)
}

func TestParseLiteral(t *testing.T) {
	test := "7"

	b, err := getBuilderFromString(test)
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseExpression()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	fmt.Printf("programAST %+v\n", programAST)
}

// TODO: later
// func TestParseIdent(t *testing.T) {
// 	test := "something[7]"

// 	b, err := getBuilderFromString(test)
// 	if err != nil {
// 		t.Errorf("err %+v\n", err)
// 	}

// 	programAST, err := b.ParseIdent()
// 	if err != nil {
// 		fmt.Printf("err %+v\n", err)
// 		t.Fatal()
// 	}

// 	fmt.Printf("programAST %+v\n", programAST)
// }

func TestParseIndexAssignmentStatement(t *testing.T) {
	test := "something[7] = \"hey its me\""

	b, err := getBuilderFromString(test)
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseAssignmentStatement()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	fmt.Printf("programAST %+v\n", programAST)
}

func TestParseAssignmentFromIndexExpression(t *testing.T) {
	test := "something = something[9][0]"

	b, err := getBuilderFromString(test)
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseAssignmentStatement()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	// Use DFS for this
	fmt.Printf("programAST %+v\n", programAST)
	fmt.Printf("programAST %+v\n", programAST.Left)
	fmt.Printf("programAST %+v\n", programAST.Right)
}

// func TestParseIndexExpression(t *testing.T) {
// 	test := "[7]"

// 	b, err := getBuilderFromString(test)
// 	if err != nil {
// 		t.Errorf("err %+v\n", err)
// 	}

// 	programAST, err := b.ParseExpression()
// 	if err != nil {
// 		fmt.Printf("err %+v\n", err)
// 		t.Fatal()
// 	}

// 	fmt.Printf("programAST %+v\n", programAST)
// }

func TestParseIdentIndexExpression(t *testing.T) {
	test := "something[9][0]"

	b, err := getBuilderFromString(test)
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseExpression()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	// Use DFS for this
	fmt.Printf("programAST %+v\n", programAST)
	fmt.Printf("programAST %+v\n", programAST.Left.Left)
	fmt.Printf("programAST %+v\n", programAST.Left.Right)
	fmt.Printf("programAST %+v\n", programAST.Right)
}

// TODO: this is an expression too
func TestParseSelectionStatement(t *testing.T) {}

func TestParseTypedefStatement(t *testing.T) {}

// TODO: later
func TestParseStructStatement(t *testing.T) {}

func TestParseCallExpression(t *testing.T) {
	test := "funcYou(6, 7)"

	b, err := getBuilderFromString(test)
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseCall()
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	fmt.Printf("\nprogramAST %+v\n", programAST)
}

// func TestParseAllowStatement(t *testing.T) {}

// func TestParseUsingStatement(t *testing.T) {}

// func TestParseTypedefStatement(t *testing.T) {}

func TestParseStatement(t *testing.T) {
	test := "function something(int i, string s) { float f = 10.1 }"

	b, err := getBuilderFromString(test)
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseStatement()
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	fmt.Printf("\nprogramAST %+v\n", programAST)
}

// Not sure if we need this because we have the group of statements thing
// func TestParseMultipleStatements(t *testing.T) {}

func TestParseStructBlockExpression(t *testing.T) {}

// TODO: this is an object; wait till later to do it
func TestParseBlockExpression(t *testing.T) {}

func TestParseLetStatement(t *testing.T) {}
