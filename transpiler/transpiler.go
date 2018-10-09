package transpiler

import (
	"errors"
	"fmt"
	"strconv"

	ast "github.com/scottshotgg/express-ast"
)

func TranslateExpression(e ast.Expression) (string, error) {
	switch e.Kind() {
	case ast.IdentNode:
		// FIXME: need to check ok on all of these
		i := e.(*ast.Ident)

		// TODO: make a TranslateIdent node
		return i.TypeOf.Name + " " + i.Name, nil

	case ast.LiteralNode:
		l := e.(ast.Literal)

		// FIXME: fill out the switch statement
		switch l.Type().Type {
		case ast.IntType:
			// FIXME: this def needs to be checked
			return strconv.Itoa(l.(*ast.IntLiteral).Value), nil
		}
	}

	// TODO: just return this for now as the default value of the function
	fmt.Println(e.Kind())
	return "", errors.New("couldnt determine")
}

func TranslateAssignmentStatement(a *ast.Assignment) (string, error) {
	// TODO: Would be nice to have a type indication for array here ...

	// Always put "=" because there is no ":=" in C++; we are just using it for the compiler
	lhs, err := TranslateExpression(a.LHS)
	if err != nil {
		return "", err
	}

	rhs, err := TranslateExpression(a.RHS)
	if err != nil {
		return "", err
	}

	return lhs + "=" + rhs + ";", nil
}

func Transpile(p *ast.Program) (string, error) {
	fmt.Println(p)

	for _, file := range p.Files {
		return file.String(), nil

		// // FIXME: make an array the size of the statements
		// // this should really transpile a 'BLOCK'
		// // scatter/gather the statements
		// // - do a parallelize the statement parsing after that and then recombine
		// for _, stmt := range file.Statements {
		// 	fmt.Println("stmt", stmt)

		// 	switch stmt.Kind() {
		// 	case ast.AssignmentNode:
		// 		fmt.Println("I got an assignment")
		// 		cStmt, err := TranslateAssignmentStatement(stmt.(*ast.Assignment))
		// 		if err != nil {
		// 			return "", nil
		// 		}

		// 	cProgramJargon += cStmt
		// }
		// }
	}

	return "}", nil
}
