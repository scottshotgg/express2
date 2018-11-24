package builder_test

import (
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
			t.Errorf(errString, err)
		}

		programAST, err = b.ParseStatement()
		if err != nil {
			t.Errorf(errString, err)
		}

		fmt.Printf(programASTString, programAST)
	}
}

func TestAllExpressions(t *testing.T) {
	for name, expr := range tests[ExpressionTest] {
		fmt.Printf(runningString, name)
		b, err = getBuilderFromString(expr)
		if err != nil {
			t.Errorf(errString, err)
		}

		programAST, err = b.ParseExpression()
		if err != nil {
			t.Errorf(errString, err)
		}

		fmt.Printf(programASTString, programAST)
	}
}
