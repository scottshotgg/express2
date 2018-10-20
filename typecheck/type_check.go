package typecheck

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
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

	case ast.BlockNode:
		m.NewScope()

		// Need to check the block
		_, err := CheckStatements(e.(*ast.Block).Statements)
		if err != nil {
			return nil, err
		}

		_, err = m.ExitScope()
		if err != nil {
			return nil, err
		}

		return e.(*ast.Block).Type(), nil

	case ast.ArrayNode:
		// FIXME: actually check the type
		// TODO: This is going to need to be determined by whether or not it is homogenous
		return e.(*ast.Array).Type(), nil
	}

	// TODO: just return this for now as the default value of the function
	fmt.Println(e.Kind())
	// FIXME: This should be able to return nil
	return nil, errors.Errorf("could not determine expression type in type checker %v", e)
}

func setTypeOfExpression(e1 ast.Expression, e2 ast.Expression) error { return nil }

// TODO: need to have a map of variables that are used to track the type checking
// Port over the variable mapping algorithm/scheme from the first Express
// Might just do a monolithic type thing and then change the AST and output that

type VariableNode struct {
	Statement ast.Statement
	Ident     *ast.Ident
	Type      *ast.Type
	IsUsed    bool
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

func TypeCheck(p *ast.Program) error {
	var err error

	for _, file := range p.Files {
		// TODO: solve this later
		// Create a new scope for each file
		// m.NewScope()
		file.Statements, err = CheckStatements(file.Statements)
		if err != nil {
			return err
		}
	}

	return nil
}

func CheckStatements(statements []ast.Statement) ([]ast.Statement, error) {
	for i, stmt := range statements {
		switch stmt.Kind() {
		case ast.AssignmentNode:
			// TODO: if it is a declaration then we need to check that the variable is not already in the variable map
			//
			as := stmt.(*ast.Assignment)
			fmt.Println("got an assignment", as)

			if as.LHS.Kind() == ast.IdentNode {
				if as.Declaration {
					fmt.Println("checking")
					// We should make an interface called Assignable
					if as.LHS.Kind() == ast.IdentNode {
						_, ok := m.CurrentScope[as.LHS.(*ast.Ident).Name]
						if ok {
							return nil, errors.New("variable already declared")
						}

						type2, err := getTypeOfExpression(as.RHS)
						if err != nil {
							return nil, err
						}

						if as.Inferred {
							as.LHS.(*ast.Ident).TypeOf = type2
						}

						m.CurrentScope[as.LHS.(*ast.Ident).Name] = &VariableNode{
							Statement: statements[i],
							Ident:     as.LHS.(*ast.Ident),
							Type:      as.LHS.(*ast.Ident).TypeOf,
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
					return nil, errors.New("Use of undeclared variable")
				}

				type2, err := getTypeOfExpression(as.RHS)
				if err != nil {
					return nil, err
				}

				fmt.Println("something", variable.Type, type2.Type, type2.UpgradesTo)

				fmt.Println("checking types ", variable, type2)
				// If the types are not directly the same then check whether the right hand side can upgrade
				if variable.Type.Type != ast.VarType && variable.Type.Type != type2.Type {
					if variable.Type.Type != type2.UpgradesTo {
						return nil, errors.Errorf("Types did not match %v %v", as.LHS, as.RHS)
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

			_, err := CheckStatements(stmt.(*ast.Block).Statements)
			if err != nil {
				return nil, err
			}

			_, err = m.ExitScope()
			if err != nil {
				return nil, err
			}

		case ast.FunctionNode:
			// TODO: need to look into local scoping
			m.NewScope()

			_, err := CheckStatements(stmt.(*ast.Function).Body.Statements)
			if err != nil {
				return nil, err
			}

			_, err = m.ExitScope()
			if err != nil {
				return nil, err
			}

		case ast.LoopNode:
			m.NewScope()

			// TODO: add , ok around these
			loop := stmt.(*ast.Loop)

			// if loop.Type == ast.ForIn || ast.Type == ast.ForOf || ast.Type == ast.ForOver {

			// }

			if loop.Type == ast.StdFor {
				_, err := CheckStatements(append([]ast.Statement{loop.Init}, loop.Body.Statements...))
				if err != nil {
					return nil, err
				}
			} else {
				_, err := CheckStatements(loop.Body.Statements)
				if err != nil {
					return nil, err
				}
			}

			_, err := m.ExitScope()
			if err != nil {
				return nil, err
			}
		}
	}

	// for _, variable := range m.CurrentScope {
	// 	if !variable.IsUsed {
	// 		// stmts = append(stmts, variable.Statement)
	// 		variable.Ident = nil
	// 	}

	// 	// fmt.Println(variable)
	// }

	return statements, nil
}
