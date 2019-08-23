package builder_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scottshotgg/express2/test"
)

func TestParseBinOpAssignmentStatement(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["binop"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseIdentStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

// func TestParseDeclarationStatement(t *testing.T) {
// 	// TODO: we need the rest of the declaration types and stuff
// 	b, err = getBuilderFromString(test.Tests[test.StatementTest]["decl"])
// 	if err != nil {
// 		t.Errorf(errFormatString, err)
// 	}

// 	node, err = b.ParseDeclarationStatement(nil)
// 	if err != nil {
// 		t.Errorf(errFormatString, err)
// 	}

// 	nodeJSON, _ = json.Marshal(node)
// 	fmt.Printf(jsonFormatString, nodeJSON)

// 	var v = b.ScopeTree.Get("i")
// 	if v == nil {
// 		t.Fatalf("Could not find variable after insertion")
// 	}

// 	nodeJSON, _ = json.Marshal(v)
// 	fmt.Printf(jsonFormatString, nodeJSON)
// }

func TestParseAssignmentFromIndexStatement(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["assignFromIndex"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseIdentStatement()
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

// func TestScopeTreeAssignmentStatement(t *testing.T) {

// }

func TestParseIdentStatement(t *testing.T) {
	// var totalTest = test.Tests[test.StatementTest]["decl"] + " " + test.Tests[test.StatementTest]["simpleAssign"]

	var tests = map[string]error{
		"int i = 0":   nil,
		"int i":       nil,
		"i = 0":       nil,
		"c.fputs()":   nil,
		"c.int i = 0": nil,
	}

	// TODO: Figure out how we can run test like the above

	for test := range tests {
		b, err = getBuilderFromString(test)
		if err != nil {
			t.Errorf(errFormatString, err)
		}

		node, tests[test] = b.ParseIdentStatement()

		nodeJSON, _ = json.Marshal(node)
		fmt.Printf(jsonFormatString, nodeJSON)
	}

	fmt.Println("Report:", tests)
}

func TestIfElseStatement(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["ifElse"])
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
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["sgroup"])
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
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["funcDef"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseFunctionStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(b.ScopeTree.Get("something"))
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseCallStatement(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["callNonAssign"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseIdentStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseCallAssignmentStatement(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["callAssign"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseIdentStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseBlockStatement(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["block"])
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
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["derefAssign"])
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
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["package"])
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

func TestParseCImportStatement(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["cimport"])
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

func TestParseImportStatement(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["import"])
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
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["include"])
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
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["stdFor"])
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
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["stdFor"])
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

// func TestParseArrayDeclaration(t *testing.T) {
// 	b, err = getBuilderFromString(test.Tests[test.StatementTest]["arrayDef"])
// 	if err != nil {
// 		t.Errorf(errFormatString, err)
// 	}

// 	node, err = b.ParseDeclarationStatement(nil)
// 	if err != nil {
// 		t.Errorf(errFormatString, err)
// 	}

// 	nodeJSON, _ = json.Marshal(node)
// 	fmt.Printf(jsonFormatString, nodeJSON)
// }

func TestParseForInStatement(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["forin"])
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
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["forin"])
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
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["indexAssign"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseIdentStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseStatement(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["funcDef"])
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
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["selectionAssign"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseIdentStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	// Remember: The left always provides the value...
	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseAssignmentFromSelectionStatement(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["assignFromSelect"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseIdentStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	// Remember: The left always provides the value...
	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseTypeDeclarationStatement(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["typeDef"])
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
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["returnSomething"])
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
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["struct"])
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

	nodeJSON, _ = json.Marshal(b.ScopeTree)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseStructDeclarationStatement(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["struct"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseStructDeclarationStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	// Remember: The left always provides the value...
	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseLetStatement(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["simpleLet"])
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
