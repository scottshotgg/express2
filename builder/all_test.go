package builder_test

import (
	"encoding/json"
	"fmt"
	"testing"
)

const (
	runningString = "Running %s test ...\n"
)

func TestAllStatements(t *testing.T) {
	for name, stmt := range tests[StatementTest] {
		fmt.Printf(runningString, name)
		b, err = getBuilderFromString(stmt)
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
}

func TestAllExpressions(t *testing.T) {
	for name, expr := range tests[ExpressionTest] {
		fmt.Printf(runningString, name)
		b, err = getBuilderFromString(expr)
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
}
