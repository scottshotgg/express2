package builder_test

import (
	"fmt"
	"testing"
)

func TestParseDeclarationStatement(t *testing.T) {
	// TODO: we need the rest of the declaration types and stuff
	b, err := getBuilderFromString(tests[StatementTest]["decl"])
	if err != nil {
		fmt.Println("err", err)
		t.Fatal()
	}

	programAST, err := b.ParseDeclarationStatement()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	fmt.Printf("programAST %+v\n", programAST)
}

func TestParseAssignmentFromIndexStatement(t *testing.T) {
	b, err := getBuilderFromString(tests[StatementTest]["assignFromIndex"])
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

func TestParseAssignmentStatement(t *testing.T) {
	b, err := getBuilderFromString(tests[StatementTest]["simpleAssign"])
	if err != nil {
		fmt.Println("err", err)
		t.Fatal()
	}

	programAST, err := b.ParseAssignmentStatement()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	fmt.Printf("programAST %+v\n", programAST)
}

func TestIfElseStatement(t *testing.T) {
	b, err := getBuilderFromString(tests[StatementTest]["ifElse"])
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
	b, err := getBuilderFromString(tests[StatementTest]["funcDef"])
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseFunctionStatement()
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	fmt.Printf("\nprogramAST %+v\n", programAST)
}

func TestParseCallAssignmentStatement(t *testing.T) {
	b, err := getBuilderFromString(tests[StatementTest]["callAssign"])
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseAssignmentStatement()
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	fmt.Printf("\nprogramAST %+v\n", programAST)
}

func TestParseBlockStatement(t *testing.T) {
	b, err := getBuilderFromString(tests[StatementTest]["block"])
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

func TestParseDerefAssignmentStatement(t *testing.T) {
	b, err := getBuilderFromString(tests[StatementTest]["derefAssign"])
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

func TestParsePackageStatement(t *testing.T) {
	b, err := getBuilderFromString(tests[StatementTest]["package"])
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParsePackageStatement()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	fmt.Printf("programAST %+v\n", programAST)
}

func TestParseImportStatement(t *testing.T) {
	b, err := getBuilderFromString(tests[StatementTest]["import"])
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
	b, err := getBuilderFromString(tests[StatementTest]["include"])
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

func TestParseForStatement(t *testing.T) {
	b, err := getBuilderFromString(tests[StatementTest]["stdFor"])
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseForStatement()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	fmt.Printf("programAST %+v\n", programAST)
}

func TestParseForStdStatement(t *testing.T) {
	b, err := getBuilderFromString(tests[StatementTest]["stdFor"])
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

func TestParseArrayDeclaration(t *testing.T) {
	b, err := getBuilderFromString(tests[StatementTest]["arrayDef"])
	if err != nil {
		t.Errorf("err %+v\n", err)
	}

	programAST, err := b.ParseDeclarationStatement()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		t.Fatal()
	}

	fmt.Printf("programAST %+v\n", programAST)
}

func TestParseForInStatement(t *testing.T) {
	b, err := getBuilderFromString(tests[StatementTest]["forin"])
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
	b, err := getBuilderFromString(tests[StatementTest]["forin"])
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

func TestParseIndexAssignmentStatement(t *testing.T) {
	b, err := getBuilderFromString(tests[StatementTest]["indexAssign"])
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
