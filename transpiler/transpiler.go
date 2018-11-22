package transpiler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/scottshotgg/express-ast"
	"github.com/scottshotgg/express-token"
)

var includes = map[string]bool{}
var genMain = true

func TranslateExpression(e ast.Expression, name string) (string, error) {
	switch e.Kind() {
	case ast.IdentNode:
		// FIXME: need to check ok on all of these
		i := e.(*ast.Ident)

		typeOfString := i.TypeOf.Name
		switch i.TypeOf.Name {
		case token.StringType:
			includes["string"] = true
			typeOfString = "std::" + typeOfString

		case token.ObjectType:
			fallthrough

		case token.VarType:
			includes["lib/var.cpp"] = true
		}

		// if i.TypeOf.Array {
		// 	i.Name += "[]"
		// }

		return typeOfString + " " + i.Name, nil

	case ast.LiteralNode:
		return e.(ast.Literal).String(), nil

	case ast.BlockNode:
		// FIXME: this needs to translate a different way if it is going to be an expression
		// TODO: create a separate TranslateObject function
		// return TranspileBlock(e.(*ast.Block).Statements)
		// TODO: if we added the ident to the block too ... ??
		// return TranspileObject(e.(*ast.Block).Statements, e.(*ast.Block).Ident.Name)
		fmt.Println("name is", name, e)
		return TranspileObject(e.(*ast.Block).Statements, name)

	case ast.BinaryOperationNode:
		bo := e.(*ast.BinaryOperation)

		op := ""
		switch bo.Op {
		case ast.AdditionBinaryOp:
			op = "+"

		case ast.SubtractionBinaryOp:
			op = "-"

		case ast.MultiplicationBinaryOp:
			op = "*"

		case ast.DivisionBinaryOp:
			op = "/"

		default:
			return "", errors.Errorf("Binary operation not defined: %v", bo.Op)
		}

		lhs, err := TranslateExpression(bo.LeftNode, name)
		if err != nil {
			return "", err
		}

		rhs, err := TranslateExpression(bo.RightNode, name)
		if err != nil {
			return "", err
		}

		return lhs + op + rhs, nil

	case ast.ConditionNode:
		c := e.(*ast.Condition)

		// TODO: switch on this later
		op := "<"

		lhs, err := TranslateExpression(c.Left, name)
		if err != nil {
			return "", err
		}

		rhs, err := TranslateExpression(c.Right, name)
		if err != nil {
			return "", err
		}

		return lhs + op + rhs, nil

	case ast.UnaryNode:
		uo := e.(*ast.UnaryOp)

		rhs, err := TranslateExpression(uo.Value, name)
		if err != nil {
			return "", err
		}

		return rhs + "++", nil

	case ast.ArrayNode:
		var (
			a = e.(*ast.Array)
			// elements = make([]string, len(a.Elements))
			// elements = []string{}
			// err      error

			arrayString = "{};"
		)

		// if a.Homogenous && a.Type().Type != ast.VarType {
		// 	var (
		// 		err      error
		// 		elements = make([]string, len(a.Elements))
		// 	)

		// 	for i, elem := range a.Elements {
		// 		if elem.Type().Type == ast.ObjectType {
		// 			arrayString += name + "[" + strconv.Itoa(i) + "] = {};"
		// 			elements[i], err = TranslateExpression(elem, name+"["+strconv.Itoa(i)+"]")
		// 		} else {
		// 			elements[i], err = TranslateExpression(elem, name)
		// 		}

		// 		if err != nil {
		// 			return "", err
		// 		}
		// 	}

		// 	if a.Type().Type == ast.ObjectType {
		// 		arrayString += strings.Join(elements, "")
		// 	} else {
		// 		arrayString = "{" + strings.Join(elements, ", ") + "}"
		// 	}

		// 	return arrayString, nil
		// }

		var (
			// err      error
			elements = make([]string, len(a.Elements))
		)

		for i, elem := range a.Elements {
			fmt.Println("hey its me", elem)

			thing, err := TranslateExpression(elem, name+"["+strconv.Itoa(i)+"]")
			if elem.Type().Type == ast.ObjectType { //|| elem.Type().Type == ast.ArrayType {
				arrayString += name + "[" + strconv.Itoa(i) + "] = {};"
			} else {
				fmt.Println(thing)
				thing = name + "[" + strconv.Itoa(i) + "] = " + thing + ";"
			}

			elements[i] = thing

			if err != nil {
				return "", err
			}
		}

		// if !a.Homogenous {
		arrayString += strings.Join(elements, "")
		// } else {
		// 	arrayString = "{" + strings.Join(elements, ", ") + "}"
		// }

		return arrayString, nil
	}

	// TODO: just return this for now as the default value of the function
	return "", errors.Errorf("could not determine expression type: %v", e)
}

func TranslateAssignmentStatement(a *ast.Assignment, name string) (string, error) {
	var (
		lhs, rhs string
		err      error
	)

	if a.LHS.Kind() != ast.IdentNode {
		return "", errors.Errorf("LHS was not an ident: %v", lhs)
	}

	lhs, err = TranslateExpression(a.LHS, name)
	if err != nil {
		return "", err
	}

	// TODO: will probably need to check the type
	ident, ok := a.LHS.(*ast.Ident)
	if !ok {
		// for some reason we have an assignment expression where the left side is not an ident
		return "", errors.Errorf("Left side of assignment was not an ident somehow... : %v", a.LHS)
	}

	// switching on the type here will work if the assignment is not an inference
	// we may need to do a deeper check to resolve the type if it is inferred
	if ident.Type().Type == ast.ObjectType || ident.Type().Type == ast.StructType || ident.Type().Type == ast.VarType {
		if ident.Name == "" {
			// Somehow we processed an ident without a name ...
			return "", errors.New("Left side ident did not have a name")
		}

		// a.RHS.Type().Array = true

		rhs, err = TranslateExpression(a.RHS, name+ident.Name)
		if err != nil {
			return "", err
		}
	} else {
		rhs, err = TranslateExpression(a.RHS, name+ident.Name)
		if err != nil {
			return "", err
		}
	}

	fmt.Println("waddup", a.LHS)
	// Check if the array type is a var or not; we don't need to have an array of vars on the backend
	if a.LHS.Type().Array && a.LHS.Type().Type != ast.VarType && a.LHS.Type().Type != ast.ObjectType {
		arr, ok := a.RHS.(*ast.Array)
		if ok {
			fmt.Println("i am here2", a.LHS.TokenLiteral())
			lhs += "[" + strconv.Itoa(len(arr.Elements)) + "]"
		} else {
			lhs += "[]"
		}
	}

	// Always put "=" because there is no ":=" in C++; we are just using it for the compiler
	return lhs + "=" + rhs + ";", nil
}

func TranspileObject(statements []ast.Statement, name string) (string, error) {
	// TODO: implement all object logic here for the assignments and stuff; would like to keep it in the same function but w/e

	fmt.Println("i am here", statements, name)

	objectString := "{};\n"

	for _, stmt := range statements {
		switch stmt.Kind() {
		case ast.AssignmentNode:
			as := stmt.(*ast.Assignment)

			var rhs string
			var err error
			// FIXME: for now lets just test objects with idents, can make literals later
			// as.LHS.Type()
			fmt.Println("hi", as.RHS.Type())
			if as.RHS.Type().Type == ast.ObjectType {
				rhs, err = TranslateExpression(as.RHS, name+"[\""+as.LHS.(*ast.Ident).Name+"\"]")
			} else {
				rhs, err = TranslateExpression(as.RHS, name)
			}

			if err != nil {
				return "", err
			}

			objectString += name + "[\"" + as.LHS.(*ast.Ident).Name + "\"] = " + rhs + ";"
		}
	}

	return objectString, nil
}

func TranspileLoop(f *ast.Loop) (string, error) {
	// FIXME: this is a hack to get around the body not generating the right {}
	if len(f.Body.Statements) < 1 {
		return "", nil
	}

	switch f.Type {
	case ast.StdFor:
		as, err := TranslateAssignmentStatement(f.Init, "")
		if err != nil {
			return "", err
		}

		cond, err := TranslateExpression(f.Cond, "")
		if err != nil {
			return "", err
		}

		post, err := TranslateExpression(f.Post, "")
		if err != nil {
			return "", err
		}

		body, err := TranspileBlock(f.Body.Statements)
		if err != nil {
			return "", err
		}

		return "for (" + as + cond + ";" + post + ")" + body, nil

	case ast.ForIn:
		fallthrough
	case ast.ForOf:
		fallthrough
	case ast.ForOver:
		return "", errors.New("Preposition loops are not implemented")

		// get the ident
		//	this is used as a temp variable, name it something random
		// get the keyword
		//	this will determine whether you want i or array[i]
		// get the expression
		//	make sure its an array
		//	insert an array node so that the transpiler will generate the array
		// get the block
		// dump all of the aforementioned in another block

		/*
			{
				array = [ARRAY]
				i_random = 0
					or
				i = 0
				while i_random < len(array) {
					i = array[i_random]
						or
					i = i_random

					i_random++
				}
			}
		*/
	}

	return "", errors.Errorf("Could not transpile loop type: %v", f)
}

func TranspileBlock(statements []ast.Statement) (string, error) {
	cProgramJargon := ""

	for _, stmt := range statements {
		switch stmt.Kind() {
		case ast.AssignmentNode:
			cStmt, err := TranslateAssignmentStatement(stmt.(*ast.Assignment), "")
			if err != nil {
				return "", err
			}

			cProgramJargon += cStmt

		case ast.FunctionNode:
			f := stmt.(*ast.Function)
			// Technically don't have to do this since clang will probably
			// optimize out the '#include<functional>' anyways if it isn't used
			if f.Ident.Name != "main" {
				includes["functional"] = true
			}

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

		case ast.LoopNode:
			cStmt, err := TranspileLoop(stmt.(*ast.Loop))
			if err != nil {
				return "", err
			}

			cProgramJargon += cStmt

		default:
			return "", errors.Errorf("Transpilation for statement has not been implemented %v", stmt.String())
		}
	}

	if len(cProgramJargon) > 0 {
		cProgramJargon = "{" + cProgramJargon + "}"
	}

	return cProgramJargon, nil
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
		// TODO: this is kinda a hack to get around the double block in main
		if len(cProgramJargon) < 1 {
			cProgramJargon = "{" + cProgramJargon + "}"
		}
		cProgramJargon = "int main() " + cProgramJargon
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
