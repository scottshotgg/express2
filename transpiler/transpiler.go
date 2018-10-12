package transpiler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/scottshotgg/express-ast"
	"github.com/scottshotgg/express-token"
)

var includes = map[string]bool{}

func TranslateExpression(e ast.Expression) (string, error) {
	switch e.Kind() {
	case ast.IdentNode:
		// FIXME: need to check ok on all of these
		i := e.(*ast.Ident)

		switch i.TypeOf.Name {
		case token.StringType:
			includes["string"] = true
			return "std::" + i.TypeOf.Name + " " + i.Name, nil

		case token.ObjectType:
			fallthrough

		case token.VarType:
			includes["lib/var.cpp"] = true
		}

		// TODO: make a TranslateIdent node
		return i.TypeOf.Name + " " + i.Name, nil

	case ast.LiteralNode:
		l := e.(ast.Literal)

		// FIXME: fill out the switch statement
		switch l.Type().Type {
		// case ast.IntType:
		// 	// FIXME: this def needs to be checked
		// 	return strconv.Itoa(l.(*ast.IntLiteral).Value), nil

		// case ast.StringType:
		default:
			return l.String(), nil
		}

	case ast.BlockNode:
		// FIXME: this needs to translate a different way if it is going to be an expression
		// TODO: create a separate TranslateObject function
		// return TranspileBlock(e.(*ast.Block).Statements)
		// TODO: if we added the ident to the block too ... ??
		// return TranspileObject(e.(*ast.Block).Statements, e.(*ast.Block).Ident.Name)
		return TranspileObject(e.(*ast.Block).Statements)
	}

	// TODO: just return this for now as the default value of the function
	fmt.Println(e.Kind())
	return "", errors.New("could not determine expression type")
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

	fmt.Println("lhs", lhs)
	return lhs + "=" + rhs + ";", nil
}

var genMain = true

func TranspileObject(statements []ast.Statement) (string, error) {
	// TODO: implement all object logic here for the assignments and stuff; would like to keep it in the same function but w/e

	objectString := "{}"

	// for _, stmt := range statements {
	// 	switch stmt.Kind() {
	// 	case ast.AssignmentNode:
	// 		cStmt, err := TranslateAssignmentStatement(stmt.(*ast.Assignment))
	// 		if err != nil {
	// 			return "", err
	// 		}

	// 		objectString += cStmt
	// 	}
	// }

	return objectString, nil
}

func TranspileBlock(statements []ast.Statement) (string, error) {
	cProgramJargon := "{"

	for _, stmt := range statements {
		switch stmt.Kind() {
		case ast.AssignmentNode:
			cStmt, err := TranslateAssignmentStatement(stmt.(*ast.Assignment))
			if err != nil {
				return "", err
			}

			cProgramJargon += cStmt

		case ast.FunctionNode:
			includes["functional"] = true

			f := stmt.(*ast.Function)
			blockString, err := TranspileBlock(f.Body.Statements)
			if err != nil {
				return "", err
			}
			return1 := "void"

			// Don't know if we need this, just being cautious rn
			if f.Returns != nil && f.Returns.Elements[0] != nil {
				return1 = f.Returns.Elements[0].(*ast.Ident).Name
			}

			functionString := ""
			if f.Ident.Name == "main" {
				if genMain == false {
					return "", errors.New("Cannot have two main functions")
				}

				genMain = false
				functionString = "int main()" + blockString
			} else {
				// FIXME: put all the functions at the top of the C++ file
				functionString = return1 + " " + f.Ident.Name + f.Arguments.String() + blockString
			}

			functions = append(functions, functionString)

		default:
			cProgramJargon += stmt.String()
		}
	}

	// if len(cProgramJargon) > 0 {
	// 	cProgramJargon = cProgramJargon + "}"
	// }
	return cProgramJargon + "}", nil
}

var functions []string

func Transpile(p *ast.Program) (string, error) {
	fmt.Println(p)

	// Put all these functions and crap into a struct that has channels/readers, etc

	cProgramJargon := ""

	for _, file := range p.Files {
		// return file.String(), nil

		// FIXME: make an array the size of the statements
		// this should really transpile a 'BLOCK'
		// scatter/gather the statements
		// - do a parallelize the statement parsing after that and then recombine
		blockString, err := TranspileBlock(file.Statements)
		if err != nil {
			return "", err
		}

		cProgramJargon += blockString
	}

	includesArray := []string{}
	for include := range includes {
		if strings.Contains(include, ".cpp") || strings.Contains(include, ".h") {
			includesArray = append(includesArray, "#include \""+include+"\"")
		} else {
			includesArray = append(includesArray, "#include <"+include+">")
		}
	}

	if genMain {
		cProgramJargon = "\nint main() " + cProgramJargon
	}

	if len(includesArray) > 0 {
		return strings.Join(includesArray, "\n") + "\n" + strings.Join(functions, "\n") + cProgramJargon, nil
	}

	return strings.Join(functions, "\n") + cProgramJargon, nil
}

/*
	Before making any more advancements on the transpiler, we need to think about the impacts that the prior stages of the compiler will have
	on the organization of the nodes in the AST.

	- Turn a file into a block with a name
	- Add a []Functions attribute to the block to make transpiling easier
	- Add a []Imports attribute to the block to make transpiling easier
	- It is the transpilers responsibility to add C++ includes

*/
