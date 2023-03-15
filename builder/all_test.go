package builder_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scottshotgg/express2/test"
)

const (
	runningString = "Running %s test ...\n"
)

func TestAllStatements(t *testing.T) {
	for name, stmt := range test.Tests[test.StatementTest] {
		fmt.Printf(runningString, name)
		b, err = getBuilderFromString(stmt)
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
}

func TestAllExpressions(t *testing.T) {
	for name, expr := range test.Tests[test.ExpressionTest] {
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
