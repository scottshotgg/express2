package builder_test

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseBinOpAssignmentStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["binop"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseAssignmentStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseDeclarationStatement(t *testing.T) {
	// TODO: we need the rest of the declaration types and stuff
	b, err = getBuilderFromString(tests[StatementTest]["decl"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseDeclarationStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseAssignmentFromIndexStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["assignFromIndex"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseAssignmentStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	// Use DFS for this
	// 	nodeJSON, _ = json.Marshal(node) fmt.Printf(jsonFormatString, nodeJSON)
	// fmt.Printf(astFormatString, node.Left)
	// fmt.Printf(astFormatString, node.Right)
	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseAssignmentStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["simpleAssign"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseAssignmentStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestIfElseStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["ifElse"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseIfStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseGroupOfStatements(t *testing.T) {
	test := "(int i, string s)"

	b, err = getBuilderFromString(test)
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseGroupOfStatements()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseFunctionStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["funcDef"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseFunctionStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseCallAssignmentStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["callAssign"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseAssignmentStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseBlockStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["block"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseBlockStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseDerefAssignmentStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["derefAssign"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParsePackageStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["package"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParsePackageStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseImportStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["import"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseImportStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseIncludeStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["include"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseIncludeStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseForStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["stdFor"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseForStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseForStdStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["stdFor"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseForStdStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseArrayDeclaration(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["arrayDef"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseDeclarationStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseForInStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["forin"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseForPrepositionStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseForOfStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["forin"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseForPrepositionStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseIndexAssignmentStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["indexAssign"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseAssignmentStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}
func TestParseStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["funcDef"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseSelectionAssignmentStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["selectionAssign"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseAssignmentStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	// Remember: The left always provides the value...
	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseAssignmentFromSelectionStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["assignFromSelect"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseAssignmentStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	// Remember: The left always provides the value...
	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseTypedefStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["typeDef"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseTypeDeclarationStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	// Remember: The left always provides the value...
	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseReturnStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["returnSomething"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseReturnStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	// Remember: The left always provides the value...
	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseStructStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["struct"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseStructStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	// Remember: The left always provides the value...
	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseLetStatement(t *testing.T) {
	b, err = getBuilderFromString(tests[StatementTest]["simpleLet"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseLetStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	// Remember: The left always provides the value...
	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}
