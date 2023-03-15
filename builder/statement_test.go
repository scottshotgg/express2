package builder_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scottshotgg/express2/test"
)

func TestParseBinOpAssignmentStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["binop"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseIdentStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseIsEqualBoolDeclarationStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["isEqualBool"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseIdentStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

// func TestParseDeclarationStmt(t *testing.T) {
// 	// TODO: we need the rest of the declaration types and stuff
// 	b, err = getBuilderFromString(test.Tests[test.StatementTest]["decl"])
// 	if err != nil {
// 		t.Errorf(errFormatString, err)
// 	}

// 	node, err = b.ParseDeclarationStmt(nil)
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

func TestParseAssignmentFromIndexStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["assignFromIndex"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseIdentStmt()
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

func TestParsePointerDeclarationStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["pointerDeclaration"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseStmt()
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

// func TestScopeTreeAssignmentStmt(t *testing.T) {

// }

func TestParseIdentStmt(t *testing.T) {
	// var totalTest = test.Tests[test.StatementTest]["decl"] + " " + test.Tests[test.StatementTest]["simpleAssign"]

	var tests = map[string]error{
		"int i = 0":   nil,
		"*int o = &i": nil,
	}

	// TODO: Figure out how we can run test like the above

	for test := range tests {
		b, err = getBuilderFromString(test)
		if err != nil {
			t.Errorf(errFormatString, err)
		}

		node, tests[test] = b.ParseIdentStmt()

		nodeJSON, _ = json.Marshal(node)
		fmt.Printf(jsonFormatString, nodeJSON)
	}

	fmt.Println("Report:", tests)
}

func TestIfElseStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["ifElse"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseIfStmt()
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

func TestParseFunctionStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["funcDef"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseFunctionStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

var testt = `
func main() {
	int i = now() + 1000
}`

func TestTestt(t *testing.T) {
	b, err = getBuilderFromString(testt)
	if err != nil {
		t.Fatalf(errFormatString, err)
	}

	ast, err := b.BuildAST()
	if err != nil {
		t.Fatalf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(ast)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseAnotherFunctionStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["anotherFuncDef"])
	if err != nil {
		t.Fatalf(errFormatString, err)
	}

	node, err = b.ParseFunctionStmt()
	if err != nil {
		t.Fatalf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseCallStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["callNonAssign"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseIdentStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseCallAssignmentStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["callAssign"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseIdentStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseBlockStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["block"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseBlockStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseStatement_1622331331(t *testing.T) {
	var testt = `
	func results(map m, int x) {
		m[x] = x * x
		m[res] = m[res] + m[x]
	  
		Println(
		  "square: ", m[x], 
		  "\\nresult: ", 
		  m[res], 
		  "\\n"
		)
	  
		m[x] = m[x].to_string()
	  }
	`
	b, err = getBuilderFromString(testt)
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseDerefAssignmentStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["derefAssign"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParsePackageStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["package"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParsePackageStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseCImportStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["cimport"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseImportStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseImportStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["import"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseImportStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseIncludeStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["include"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseIncludeStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseForStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["stdFor"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseForStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseForStdStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["stdFor"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseForStdStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

// func TestParseArrayDecl(t *testing.T) {
// 	b, err = getBuilderFromString(test.Tests[test.StatementTest]["arrayDef"])
// 	if err != nil {
// 		t.Errorf(errFormatString, err)
// 	}

// 	node, err = b.ParseDeclarationStmt(nil)
// 	if err != nil {
// 		t.Errorf(errFormatString, err)
// 	}

// 	nodeJSON, _ = json.Marshal(node)
// 	fmt.Printf(jsonFormatString, nodeJSON)
// }

func TestParseForInStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["forin"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseForPrepositionStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseForOfStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["forin"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseForPrepositionStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseIndexAssignmentStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["indexAssign"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseIdentStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["funcDef"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseSelectionAssignmentStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["selectionAssign"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseIdentStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	// Remember: The left always provides the value...
	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseAssignmentFromSelectionStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["assignFromSelect"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseIdentStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	// Remember: The left always provides the value...
	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseTypeDeclarationStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["typeDef"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseTypeDeclStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	// Remember: The left always provides the value...
	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseReturnStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["returnSomething"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseReturnStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	// Remember: The left always provides the value...
	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseStructStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["struct"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseStructStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	// Remember: The left always provides the value...
	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)

	nodeJSON, _ = json.Marshal(b.ScopeTree)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseStructDeclarationStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["struct"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseStructDeclarationStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	// Remember: The left always provides the value...
	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseLetStmt(t *testing.T) {
	b, err = getBuilderFromString(test.Tests[test.StatementTest]["simpleLet"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseLetStmt()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	// Remember: The left always provides the value...
	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}
