package typeCheck

import (
	"errors"
	"fmt"

	ast "github.com/scottshotgg/express-ast"
)

type WalkerFunc func(n ast.Node) error

var TexasRanger = func(n ast.Node) error {
	fmt.Println("got an assignment")
	as := n.(*ast.Assignment)

	type1, err := getTypeOfExpression(as.LHS)
	if err != nil {
		return err
	}

	type2, err := getTypeOfExpression(as.RHS)
	if err != nil {
		return err
	}

	// If the types are not directly the same then check whether the right hand side can upgrade
	if type1.Type != type2.Type {
		if type1.Type != type2.UpgradesTo {
			return errors.New("Types did not match")
		}
	}

	// If it is inferred then set the type of the left side
	return errors.New("w/e")
}

func TreeWalk(p *ast.Program, texasRanger WalkerFunc) error {
	fmt.Println("typeCheck", p)

	for _, file := range p.Files {
		for _, stmt := range file.Statements {
			fmt.Println("stmt", stmt)

			if stmt.Kind() == ast.AssignmentNode {
				return texasRanger(stmt)
			}
		}
	}

	return nil
}
