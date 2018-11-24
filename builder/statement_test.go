package builder_test

import (
	"fmt"
	"testing"
)

func TestParseBinOpAssignmentStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["binop"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseAssignmentStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	fmt.Printf(programASTString, programAST)
}

func TestParseDeclarationStatement(t *testing.T) {
	// TODO: we need the rest of the declaration types and stuff
	b, err = getBuilderFromString(tests[StatementTest]["decl"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseDeclarationStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	fmt.Printf(programASTString, programAST)
}

func TestParseAssignmentFromIndexStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["assignFromIndex"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseAssignmentStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	// Use DFS for this
	fmt.Printf(programASTString, programAST)
	fmt.Printf(programASTString, programAST.Left)
	fmt.Printf(programASTString, programAST.Right)
}

func TestParseAssignmentStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["simpleAssign"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseAssignmentStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	fmt.Printf(programASTString, programAST)
}

func TestIfElseStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["ifElse"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseIfStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	fmt.Printf(programASTString, programAST)
}

func TestParseGroupOfStatements(t *testing.T) {
	test := "(int i, string s)"

	b, err = getBuilderFromString(test)
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseGroupOfStatements()
	if err != nil {
		t.Errorf(errString, err)
	}

	fmt.Printf("programASTString", programAST)
}

func TestParseFunctionStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["funcDef"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseFunctionStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	fmt.Printf("programASTString", programAST)
}

func TestParseCallAssignmentStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["callAssign"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseAssignmentStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	fmt.Printf("programASTString", programAST)
}

func TestParseBlockStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["block"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseBlockStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	fmt.Printf(programASTString, programAST)
}

func TestParseDerefAssignmentStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["derefAssign"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	fmt.Printf(programASTString, programAST)
}

// func TestProgram(t *testing.T) {
// 	test := `
// 		for name of names {
// 			string s = getFileContentsFromName(name)
// 			println(s)
// 	}
// 	`

// 	b, err = getBuilderFromString(test)
// 	if err != nil {
// 		t.Errorf(errString, err)
// 	}

// 	programAST, err = b.ParseStatement()
// 	if err != nil {
// 		t.Errorf(errString, err)
// 	}

// 	fmt.Printf(programASTString, programAST)
// }

func TestParsePackageStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["package"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParsePackageStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	fmt.Printf(programASTString, programAST)
}

func TestParseImportStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["import"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseImportStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	fmt.Printf(programASTString, programAST)
}

func TestParseIncludeStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["include"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseIncludeStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	fmt.Printf(programASTString, programAST)
}

func TestParseForStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["stdFor"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseForStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	fmt.Printf(programASTString, programAST)
}

func TestParseForStdStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["stdFor"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseForStdStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	fmt.Printf(programASTString, programAST)
}

func TestParseArrayDeclaration(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["arrayDef"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseDeclarationStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	fmt.Printf(programASTString, programAST)
}

func TestParseForInStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["forin"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseForPrepositionStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	fmt.Printf(programASTString, programAST)
}

func TestParseForOfStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["forin"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseForPrepositionStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	fmt.Printf(programASTString, programAST)
}

func TestParseIndexAssignmentStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["indexAssign"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseAssignmentStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	fmt.Printf(programASTString, programAST)
}
func TestParseStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["funcDef"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	fmt.Printf(programASTString, programAST)
}

func TestParseSelectionAssignmentStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["selectionAssign"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseAssignmentStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	// Remember: The left always provides the value...
	fmt.Printf(programASTString, programAST)
}

func TestParseAssignmentFromSelectionStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["assignFromSelect"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseAssignmentStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	// Remember: The left always provides the value...
	fmt.Printf(programASTString, programAST)
}

func TestParseTypedefStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["typeDef"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseTypeDeclarationStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	// Remember: The left always provides the value...
	fmt.Printf(programASTString, programAST)
}

func TestParseReturnStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["returnSomething"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseReturnStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	// Remember: The left always provides the value...
	fmt.Printf(programASTString, programAST)
}

func TestParseStructStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["struct"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseStructStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	// Remember: The left always provides the value...
	fmt.Printf(programASTString, programAST)
}

func TestParseLetStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["simpleLet"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseLetStatement()
	if err != nil {
		t.Errorf(errString, err)
	}

	// Remember: The left always provides the value...
	fmt.Printf(programASTString, programAST)
}
