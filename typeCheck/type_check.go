package typeCheck

import (
	"errors"
	"fmt"
	"os"

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

type VariableNode struct {
	Node    ast.Node
	Ident   *ast.Ident
	Type    *ast.Type
	UsedYet bool
}

type Meta struct {
	// global       Scope
	CurrentScope Scope
	scopes       *Stack
}

type Scope map[string]*VariableNode

func (m *Meta) NewScope() {
	m.scopes.Push(m.CurrentScope)
	m.CurrentScope = Scope{}
}

func (m *Meta) ExitScope() (Scope, error) {
	scope, err := m.scopes.Pop()
	if err != nil {
		// TODO:
		return Scope{}, err
	}

	m.CurrentScope = scope.(Scope)
	return m.CurrentScope, nil
}

func (m *Meta) GetVariable(variableName string) (*VariableNode, bool) {
	// Might have problems with the pointer here
	variable, ok := m.CurrentScope[variableName]
	if ok {
		return variable, true
	}

	currentScope := m.CurrentScope
	defer func(m *Meta, current Scope) {
		m.CurrentScope = current
	}(m, currentScope)

	pop, err := m.ExitScope()
	defer m.scopes.Push(pop)
	if err != nil {
		return nil, false
	}

	return m.GetVariable(variableName)
}

var m = &Meta{
	CurrentScope: Scope{},
	scopes:       NewStack(),
}

var something *ast.Statement

func TypeCheck(p *ast.Program) error {

	fmt.Println("typeCheck", p)

	var err error
	for _, file := range p.Files {
		// TODO: solve this later
		// Create a new scope for each file
		// m.NewScope()
		err = CheckStatements(file.Statements)
		if err != nil {
			return err
		}

		something = &file.Statements[3]
	}

	return nil
}

func CheckStatements(statements []ast.Statement) error {
	for _, stmt := range statements {

		switch stmt.Kind() {
		case ast.AssignmentNode:
			// TODO: if it is a declaration then we need to check that the variable is not already in the variable map
			//
			fmt.Println("got an assignment")
			as := stmt.(*ast.Assignment)

			if as.LHS.Kind() == ast.IdentNode {
				if as.Declaration {
					fmt.Println("checking")
					// We should make an interface called Assignable
					if as.LHS.Kind() == ast.IdentNode {
						_, ok := m.GetVariable(as.LHS.(*ast.Ident).Name)
						if ok {
							return errors.New("variable already declared")
						}

						if as.Inferred {
							as.LHS.(*ast.Ident).TypeOf = as.RHS.(ast.Literal).Type()
						}

						// token := as.RHS.TokenLiteral()
						// if token.Type == "DEFAULT" {
						// 	as.RHS.(*ast.DefaultLiteral).Value = defaultsMap[as.RHS.(*ast.DefaultLiteral).Value]
						// 	fmt.Println("hey", as.RHS.(*ast.DefaultLiteral).Value)
						// 	stmt = as
						// 	something = &stmt
						// }

						m.CurrentScope[as.LHS.(*ast.Ident).Name] = &VariableNode{
							Ident: as.LHS.(*ast.Ident),
							Type:  as.LHS.(*ast.Ident).TypeOf,
						}

						continue

					} else {
						fmt.Println("change this to be an assignable")
						os.Exit(9)
					}
				}

				// We should make an interface called Assignable
				// TODO: Port over the method of recursing up
				variable, ok := m.GetVariable(as.LHS.(*ast.Ident).Name)
				if !ok {
					return errors.New("Use of undeclared variable")
				}

				// as.LHS = variable.Ident

				type2, err := getTypeOfExpression(as.RHS)
				if err != nil {
					return err
				}

				fmt.Println("something", variable.Type, type2.Type, type2.UpgradesTo)

				// If the types are not directly the same then check whether the right hand side can upgrade
				if variable.Type.Type != type2.Type {
					if variable.Type.Type != type2.UpgradesTo { // || type2.UpgradesTo == 0 {
						return errors.New("Types did not match")
					}
				}

				// Right here we would check whether something was already used
				// If it hasn't been used since declaration then just initialize it to this value

			} else {
				fmt.Println("change this to be an assignable")
				os.Exit(9)
			}

		case ast.BlockNode:
			m.NewScope()
			CheckStatements(stmt.(*ast.Block).Statements)
			_, err := m.ExitScope()
			if err != nil {
				return err
			}

		}
	}

	return nil
}
