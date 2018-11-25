package builder_test

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseBinOpExpression(t *testing.T) {
	b, err = getBuilderFromString(tests[ExpressionTest]["binop"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseExpression()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
	// fmt.Printf(nodeString, *node.Metadata["args"].(*builder.Node).Value.([]*builder.Node)[0].Right)
}

func TestParseGroupOfExpressions(t *testing.T) {
	b, err = getBuilderFromString("(1, i, s, 9)")
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseGroupOfExpressions()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseDerefExpression(t *testing.T) {
	b, err = getBuilderFromString(tests[ExpressionTest]["deref"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseExpression()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseIdentExpression(t *testing.T) {
	b, err = getBuilderFromString(tests[ExpressionTest]["ident"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	// node, err =  b.ParseConditionExpression()
	node, err = b.ParseExpression()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseConditionExpression(t *testing.T) {
	b, err = getBuilderFromString(tests[ExpressionTest]["condition"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	// node, err =  b.ParseConditionExpression()
	node, err = b.ParseExpression()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseIncrementExpression(t *testing.T) {
	b, err = getBuilderFromString(tests[ExpressionTest]["inc"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	// node, err =  b.ParseIncrement(&Node{})
	node, err = b.ParseExpression()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	// fmt.Println(b.ParseBlockStatement())

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseArrayExpression(t *testing.T) {
	b, err = getBuilderFromString(tests[ExpressionTest]["array"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseArrayExpression()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseLiteral(t *testing.T) {
	b, err = getBuilderFromString(tests[ExpressionTest]["intLit"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseExpression()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseIdentIndexExpression(t *testing.T) {
	b, err = getBuilderFromString(tests[ExpressionTest]["identIndex"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseExpression()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	// Use DFS for this
	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
	fmt.Printf(astFormatString, node.Left.Left)
	fmt.Printf(astFormatString, node.Left.Right)
	fmt.Printf(astFormatString, node.Right)
}

func TestParseCallExpression(t *testing.T) {
	b, err = getBuilderFromString(tests[ExpressionTest]["identCall"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseExpression()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseBlockExpression(t *testing.T) {
	b, err = getBuilderFromString(tests[ExpressionTest]["blockExpr"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseExpression()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestParseSelectionExpression(t *testing.T) {
	b, err = getBuilderFromString(tests[ExpressionTest]["identSelect"])
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseExpression()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	// Remember: The left always provides the value...
	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

// func TestParseStructBlockExpression(t *testing.T) {}

// func TestParseAllowStatement(t *testing.T) {}

// func TestParseUsingStatement(t *testing.T) {}

// Not sure if we need this because we have the group of statements thing
// func TestParseMultipleStatements(t *testing.T) {}