package typeCheck

import (
	"errors"
	"fmt"

	"github.com/scottshotgg/express-ast"
)

// This will need to check types on:
//	- assignments
//	- function calls
//	- function returns
//		- this has to analyze return statements
//	- change type inference to type assignment

func getTypeOfExpression(e ast.Expression) (*ast.Type, error) {
	switch e.Kind() {
	case ast.IdentNode:
		// TODO: make a TranslateIdent node
		return e.(*ast.Ident).TypeOf, nil

	case ast.LiteralNode:
		return e.(ast.Literal).Type(), nil

		// // FIXME: fill out the switch statement
		// switch l.Type().Type {
		// // case ast.IntType:
		// // 	// FIXME: this def needs to be checked
		// // 	return strconv.Itoa(l.(*ast.IntLiteral).Value), nil

		// // case ast.StringType:
		// default:
		// 	return l.Type(), nil
		// }
	}

	// TODO: just return this for now as the default value of the function
	fmt.Println(e.Kind())
	// FIXME: This should be able to return nil
	return nil, errors.New("could not determine expression type")
}

func setTypeOfExpression(e1 ast.Expression, e2 ast.Expression) error {

	return nil
}

// TODO: need to have a map of variables that are used to track the type checking
// Port over the variable mapping algorithm/scheme from the first Express
// Might just do a monolithic type thing and then change the AST and output that
func TypeCheck(p *ast.Program) error {
	fmt.Println("typeCheck", p)

	for _, file := range p.Files {
		for _, stmt := range file.Statements {
			fmt.Println("stmt", stmt)

			switch stmt.Kind() {
			case ast.AssignmentNode:
				// TODO: if it is a declaration then we need to check that the variable is not already in the variable map
				//

				fmt.Println("got an assignment")
				as := stmt.(*ast.Assignment)

				if as.Type == ast.Equals {
					fmt.Println("got an equals")
					type1, err := getTypeOfExpression(as.LHS)
					if err != nil {
						return err
					}

					type2, err := getTypeOfExpression(as.RHS)
					if err != nil {
						return err
					}

					fmt.Println("something", type1.Type, type2.Type, type2.UpgradesTo)

					// If the types are not directly the same then check whether the right hand side can upgrade
					if type1.Type != type2.Type {
						if type1.Type != type2.UpgradesTo { // || type2.UpgradesTo == 0 {
							return errors.New("Types did not match")
						}
					}
				} else {
					fmt.Println("got a declaration")
					type1, err := getTypeOfExpression(as.LHS)
					if err != nil {
						return err
					}

					type2, err := getTypeOfExpression(as.RHS)
					if err != nil {
						return err
					}

					fmt.Println("something2", type1.Type, type2.Type, type2.UpgradesTo)
					// Make sure the left hand side has a type
					if type1.Type != 0 {
						return errors.New("Left hand side already has type")
					}

					// Make sure right hand side has a type
					if type2.Type == 0 {
						return errors.New("Right hand side has no type specified")
					}

					*type1 = *type2

					fmt.Println(type1)
				}
			}
		}

	}
	return nil
}
