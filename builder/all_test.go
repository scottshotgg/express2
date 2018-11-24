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
			fmt.Printf(errString, err)
			t.Fatal()
		}

		programAST, err = b.ParseStatement()
		if err != nil {
			fmt.Printf(errString, err)
			t.Fatal()
		}

		fmt.Printf(programASTString, programAST)
	}
}

func TestAllExpressions(t *testing.T) {
	for name, expr := range tests[ExpressionTest] {
		fmt.Printf(runningString, name)
		b, err = getBuilderFromString(expr)
		if err != nil {
			fmt.Printf(errString, err)
			t.Fatal()
		}

		programAST, err = b.ParseExpression()
		if err != nil {
			fmt.Printf(errString, err)
			t.Fatal()
		}

		fmt.Printf(programASTString, programAST)
	}
}
