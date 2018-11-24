package builder_test

import (
	"fmt"
	"testing"
)

func TestParseBinOpExpression(t *testing.T) {
	b, err = getBuilderFromString(tests[ExpressionTest]["binop"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseExpression()
	if err != nil {
		t.Errorf(errString, err)
	}

	fmt.Printf(programASTString, programAST)

	// fmt.Printf(programASTString, *programAST.Metadata["args"].(*builder.Node).Value.([]*builder.Node)[0].Right)
}

func TestParseGroupOfExpressions(t *testing.T) {
	b, err = getBuilderFromString("(1, i, s, 9)")
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseGroupOfExpressions()
	if err != nil {
		t.Errorf(errString, err)
	}

	fmt.Printf(programASTString, programAST)
}

func TestParseDerefExpression(t *testing.T) {
	b, err = getBuilderFromString(tests[ExpressionTest]["deref"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseExpression()
	if err != nil {
		fmt.Printf(errString, err)
		t.Fatal()
	}

	fmt.Printf(programASTString, programAST)
}

func TestParseIdentExpression(t *testing.T) {
	b, err = getBuilderFromString(tests[ExpressionTest]["ident"])
	if err != nil {
		t.Errorf(errString, err)
	}

	// programAST, err =  b.ParseConditionExpression()
	programAST, err = b.ParseExpression()
	if err != nil {
		fmt.Printf(errString, err)
		t.Fatal()
	}

	fmt.Printf(programASTString, programAST)
}

func TestParseConditionExpression(t *testing.T) {
	b, err = getBuilderFromString(tests[ExpressionTest]["condition"])
	if err != nil {
		t.Errorf(errString, err)
	}

	// programAST, err =  b.ParseConditionExpression()
	programAST, err = b.ParseExpression()
	if err != nil {
		fmt.Printf(errString, err)
		t.Fatal()
	}

	fmt.Printf(programASTString, programAST)
}

func TestParseIncrementExpression(t *testing.T) {
	b, err = getBuilderFromString(tests[ExpressionTest]["inc"])
	if err != nil {
		t.Errorf(errString, err)
	}

	// programAST, err =  b.ParseIncrement(&Node{})
	programAST, err = b.ParseExpression()
	if err != nil {
		fmt.Printf(errString, err)
		t.Fatal()
	}

	// fmt.Println(b.ParseBlockStatement())

	fmt.Printf(programASTString, programAST)
}

func TestParseArrayExpression(t *testing.T) {
	b, err = getBuilderFromString(tests[ExpressionTest]["array"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseArrayExpression()
	if err != nil {
		fmt.Printf(errString, err)
		t.Fatal()
	}

	fmt.Printf(programASTString, programAST)
}

func TestParseType(t *testing.T) {
	test := "int[][5]"

	b, err = getBuilderFromString(test)
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseType()
	if err != nil {
		fmt.Printf(errString, err)
		t.Fatal()
	}

	fmt.Printf(programASTString, programAST)
}

func TestParseLiteral(t *testing.T) {
	b, err = getBuilderFromString(tests[ExpressionTest]["intLit"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseExpression()
	if err != nil {
		fmt.Printf(errString, err)
		t.Fatal()
	}

	fmt.Printf(programASTString, programAST)
}

func TestParseIdentIndexExpression(t *testing.T) {
	b, err = getBuilderFromString(tests[ExpressionTest]["identIndex"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseExpression()
	if err != nil {
		fmt.Printf(errString, err)
		t.Fatal()
	}

	// Use DFS for this
	fmt.Printf(programASTString, programAST)
	fmt.Printf(programASTString, programAST.Left.Left)
	fmt.Printf(programASTString, programAST.Left.Right)
	fmt.Printf(programASTString, programAST.Right)
}

func TestParseCallExpression(t *testing.T) {
	b, err = getBuilderFromString(tests[ExpressionTest]["identCall"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseExpression()
	if err != nil {
		t.Errorf(errString, err)
	}

	fmt.Printf(programASTString, programAST)
}

func TestParseBlockExpression(t *testing.T) {
	b, err = getBuilderFromString(tests[ExpressionTest]["blockExpr"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseExpression()
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

func TestParseSelectionExpression(t *testing.T) {
	b, err = getBuilderFromString(tests[ExpressionTest]["identSelect"])
	if err != nil {
		t.Errorf(errString, err)
	}

	programAST, err = b.ParseExpression()
	if err != nil {
		t.Errorf(errString, err)
	}

	// Remember: The left always provides the value...
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

// func TestParseStructBlockExpression(t *testing.T) {}

// func TestParseAllowStatement(t *testing.T) {}

// func TestParseUsingStatement(t *testing.T) {}

// Not sure if we need this because we have the group of statements thing
// func TestParseMultipleStatements(t *testing.T) {}
