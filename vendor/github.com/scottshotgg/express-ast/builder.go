package ast

import (
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
	token "github.com/scottshotgg/express-token"
)

type ASTBuilder struct {
	Tokens []token.Token
	Index  int
}

func (a *ASTBuilder) GetFactor() (Expression, error) {
	currentToken := a.Tokens[a.Index]

	switch currentToken.Type {

	case token.Literal:
		switch currentToken.Value.Type {

		case token.IntType:
			return NewInt(currentToken, currentToken.Value.True.(int)), nil

		case token.BoolType:
			return NewBool(currentToken, currentToken.Value.True.(bool)), nil

		case token.FloatType:
			return NewFloat(currentToken, currentToken.Value.True.(float64)), nil

		case token.CharType:
			return NewChar(currentToken, currentToken.Value.True.(int)), nil

		case token.StringType:
			return NewString(currentToken, currentToken.Value.True.(string)), nil
		}

	case token.Ident:
		return NewIdent(currentToken, "")

		// TODO: consider changing this if we want to allow for a different syntax inside of objects
	case token.LBrace:
		return a.GetBlock()

	case token.LBracket:
		aToken := a.Tokens[a.Index]

		// Increment over the array opening
		a.Index++

		// The array should only contain expressions
		elements := []Expression{}
		for a.Tokens[a.Index].Type != token.RBracket {
			fmt.Println("expression")
			expr, err := a.GetExpression()
			if err != nil {
				return nil, err
			}

			elements = append(elements, expr)
			a.Index++
		}

		// Step over the RBracket
		// a.Index++

		return NewArray(aToken, elements), nil
	}

	return nil, errors.Errorf("Could not parse factor from token: %+v", currentToken)
}

func (a *ASTBuilder) GetTerm() (Expression, error) {
	factor, err := a.GetFactor()
	if err != nil {
		return nil, err
	}

	fmt.Println("factor", factor)

	if a.Index+1 > len(a.Tokens)-1 {
		return factor, nil
	}

	// FIXME: ideally, we should check for a `UNARY` operator class
	nextToken := a.Tokens[a.Index+1]
	if nextToken.Type == token.Increment {
		a.Index++
		factor = &UnaryOp{
			Token: a.Tokens[a.Index],
			Op:    Increment,
			Value: factor,
		}
		a.Index++
	}

	// FIXME: ideally, these should have a `CONDITIONAL` class that we can
	// test and then do a NewConditional with the token
	nextToken = a.Tokens[a.Index+1]
	if nextToken.Type == token.LThan || nextToken.Type == token.GThan {
		a.Index++
		a.Index++

		if a.Tokens[a.Index+2].Type == token.Assign {
			log.Println("wtf is this thing doing right here")
		}

		factor2, err := a.GetExpression()
		if err != nil {
			return nil, err
		}

		var ct ConditionType
		switch nextToken.Type {
		case token.LThan:
			ct = LessThan

		case token.GThan:
			ct = GreaterThan

		default:
			return nil, errors.New("Could not deduce condition type")
		}

		fmt.Println("therese a conditional")
		return &Condition{
			TypeOf: ct,
			Left:   factor,
			Right:  factor2,
		}, nil
	}

	if a.Index+1 < len(a.Tokens)-1 {
		for a.Tokens[a.Index+1].Type == token.PriOp {
			a.Index++

			operand := a.Tokens[a.Index]

			a.Index++
			factor2, err := a.GetExpression()
			if err != nil {
				return nil, err
			}

			fmt.Println("factor2", factor2, operand.Value.String)

			factor, err = NewBinaryOperation(operand, operand.Value.String, factor, factor2)
			if err != nil {
				return nil, err
			}

			if a.Index > len(a.Tokens)-1 {
				break
			}
		}
	}

	fmt.Println("returning")

	return factor, nil
}

func (a *ASTBuilder) GetExpression() (Expression, error) {
	if a.Tokens[a.Index].Type == token.Separator {
		// TODO: just skip the separator for now
		a.Index++
		return a.GetExpression()
	}

	term, err := a.GetTerm()
	if err != nil {
		return nil, err
	}

	if a.Index+1 > len(a.Tokens)-1 {
		return term, nil
	}

	// FIXME: ideally, these should have a `CONDITIONAL` class that we can
	// test and then do a NewConditional with the token
	nextToken := a.Tokens[a.Index+1]
	if nextToken.Type == token.LThan || nextToken.Type == token.GThan {
		a.Index += 2

		term2, err := a.GetExpression()
		if err != nil {
			return nil, err
		}

		var ct ConditionType
		switch nextToken.Type {
		case token.LThan:
			ct = LessThan

		case token.GThan:
			ct = GreaterThan

		default:
			return nil, errors.New("Could not deduce condition type")
		}

		fmt.Println("therese a conditional")
		return &Condition{
			TypeOf: ct,
			Left:   term,
			Right:  term2,
		}, nil
	}

	if a.Index+1 < len(a.Tokens)-1 {
		for a.Tokens[a.Index+1].Type == token.SecOp {
			a.Index++

			operand := a.Tokens[a.Index].Value.String

			a.Index++
			term2, err := a.GetExpression()
			if err != nil {
				return nil, err
			}

			fmt.Println("term2", term2, operand)

			term, err = NewBinaryOperation(a.Tokens[a.Index], operand, term, term2)
			if err != nil {
				return nil, err
			}

			if a.Index >= len(a.Tokens)-1 {
				break
			}
		}
	}

	// FIXME: should probably check for secondary operations right here

	return term, nil
}

func (a *ASTBuilder) GetGroup() (*Group, error) {
	// `(` [ expression ]* `)`
	fmt.Println("getting group", a.Tokens[a.Index])

	if a.Tokens[a.Index].Type != token.LParen {
		return nil, errors.New("Function declaration requires left paren after function identifier")
	}

	groupToken := a.Tokens[a.Index]

	elements := []Expression{}

	a.Index++
	for a.Tokens[a.Index].Type != token.RParen {
		expr, err := a.GetExpression()
		if err != nil {
			return nil, err
		}

		fmt.Println("hey its me", expr, err)
		elements = append(elements, expr)

		a.Index++
	}

	return &Group{
		// Not sure if we should create a `group` token
		Token: groupToken,
		// Not sure what type to put, maybe make a group type?
		// TypeOf: Type(0),
		Elements: elements,
	}, nil
}

func (a *ASTBuilder) GetIf() (*IfElse, error) {
	var err error

	ifElse := IfElse{
		Token: a.Tokens[a.Index],
	}

	a.Index++
	// FIXME: make a GetCondition() function later m8
	ifElse.Condition, err = a.GetExpression()
	if err != nil {
		return nil, err
	}

	a.Index++
	ifElse.Body, err = a.GetBlock()
	if err != nil {
		return nil, err
	}
	fmt.Println("ifElse", ifElse)
	// Check for an else branch
	elseToken := a.Tokens[a.Index+1]
	if elseToken.Type == token.Else {
		a.Index++

		// Check whether there is an if or just another block
		switch a.Tokens[a.Index+1].Type {
		case token.If:
			ifElse.Else, err = a.GetIf()
			if err != nil {
				return nil, err
			}

		case token.LBrace:
			a.Index++

			ifElse.Else = &IfElse{
				Token: elseToken,
			}

			ifElse.Else.Body, err = a.GetBlock()
			if err != nil {
				return nil, err
			}

		default:
			return nil, errors.New("Empty else branch")
		}
	}

	return &ifElse, nil
}

func (a *ASTBuilder) GetBlock() (*Block, error) {
	if a.Tokens[a.Index].Type != token.LBrace {
		return nil, errors.New("Could not find block opening")
	}

	lb := a.Tokens[a.Index]

	statements := []Statement{}

	a.Index++
	for a.Tokens[a.Index].Type != token.RBrace {
		stmt, err := a.GetStatement()
		if err != nil {
			return nil, err
		}

		statements = append(statements, stmt)

		a.Index++
	}

	return &Block{
		Token:      lb,
		Statements: statements,
		// Scope: //FIXME: I don't think this should be here
	}, nil
}

// GetStatement needs to switch and capture these:
//   - assignment
//   - type
//   - ident
//   - block
//   - call
//   - ident
//   - func / fn
//   - if/else
//   - loop
//   - return
func (a *ASTBuilder) GetStatement() (Statement, error) {
	typeOf := ""
	currentToken := a.Tokens[a.Index]

	switch currentToken.Type {
	case token.Separator:
		// TODO: just skip the separator for now
		a.Index++
		return a.GetStatement()

	case token.Type:
		// Look for an ident as the next thing for now
		// fallthrough to the next block for now
		typeOf = currentToken.Value.String
		a.Index++

		// Expect an ident to always follow a token for now
		fallthrough

	case token.Ident:
		// Here we will want to look at what is next and handle it
		// If it is an assignment statement then we are looking for an expression afterwards

		// FIXME: need to implement Type() so that we can get the var type
		ident, err := NewIdent(a.Tokens[a.Index], typeOf)
		if err != nil {
			return nil, err
		}

		a.Index++

		switch a.Tokens[a.Index].Type {
		case token.Assign:
			assignmentToken := a.Tokens[a.Index]

			a.Index++
			expr, err := a.GetExpression()
			if err != nil {
				return nil, err
			}
			fmt.Println("expr", expr)
			// os.Exit(9)

			// TODO: figure out why i put this here
			if expr == nil {
				return nil, nil
			}

			// TODO: could make a new boolean assignment here?
			as, err := NewAssignment(assignmentToken, ident, expr)
			if err != nil {
				return nil, err
			}

			if typeOf != "" {
				as.SetDeclaration(true)
			}

			// TODO: add statement here later
			return as, nil

		// This is a function call as a statement
		// TODO: implement function call expressions
		case token.LParen:
			// Get the group for the args
			args, err := a.GetGroup()
			if err != nil {
				return nil, err
			}

			return &Call{
				Token:     ident.Token,
				Ident:     ident,
				Arguments: args,
			}, nil

		default:
			// assignmentToken, ok := token.TokenMap["="]
			// if !ok {
			// 	return nil, errors.New("Could not find assignment token in tokenmap")
			// }

			fmt.Println("shit", currentToken.Value.String)
			// os.Exit(9)
			// If the token that comes afterwards is none of these
			// apply a default value for a declaration
			as, err := NewAssignment(token.TokenMap["="], ident, NewDefault(token.Token{
				Type: "DEFAULT",
				Value: token.Value{
					String: currentToken.Value.String,
				},
			}))
			if err != nil {
				return nil, err
			}

			as.SetDeclaration(true)

			a.Index--

			// TODO: add statement here later
			return as, nil
		}

		// By default, if there is a `[ type ] [ ident]` combination, that is a default valued initialization
		// return nil, errors.Errorf("Expected assignment token, got %+v", a.Tokens[a.Index])

	case token.LBrace:
		// Here we will want to recursively call GetStatement()
		// however, a block should be able to be parsed for an expression as well
		return a.GetBlock()

	// TODO: break this out into the individual keywords
	// - switch, etc
	case token.Keyword:
		return nil, errors.Errorf("token.Keyword statements are not implemented yet %+v", currentToken)

	case token.Function:
		// Next things we look for after the Function token is:
		//	[ ident ] [ group ] { group } [ block ]
		fmt.Println("Found a function token")

		// Get the function token
		functionToken := a.Tokens[a.Index]

		// Get the ident token
		a.Index++
		identToken := a.Tokens[a.Index]

		// Get the group for the args
		a.Index++
		args, err := a.GetGroup()
		if err != nil {
			return nil, err
		}

		// FIXME: skip getting the returns for now

		// Get the body of the function
		// the body is essentially just a list of statements
		// this is the exact same as a file in our definition

		a.Index++
		block, err := a.GetBlock()
		if err != nil {
			return nil, err
		}

		return NewFunction(functionToken, identToken, args, block)

	case token.If:
		// TODO:
		// look for a conditional/expression
		// get a block
		// check for an else
		// if theres an else, look for another if or a block

		return a.GetIf()

	// // FIXME: maybe this needs to switch to token.Loop later on
	case token.For:
		fmt.Println("we found a for loop")
		// We need to be able to parse different types of loops here:
		// - standard loops
		// - preposition loops

		// Save the `for` token
		forToken := a.Tokens[a.Index]

		a.Index++

		// Figure out what type of loop it is by the next token
		switch a.Tokens[a.Index].Type {
		// support declaring static typed variables as well
		case token.Type:
			typeOf = a.Tokens[a.Index].Value.String
			a.Index++
			fallthrough

		case token.Ident:
			ident, err := NewIdent(a.Tokens[a.Index], typeOf)
			if err != nil {
				return nil, err
			}

			// Look ahead one token to determine what type of loop it is
			switch a.Tokens[a.Index+1].Type {

			// For now just keep it like this:
			// Later we can change it to actually get specific nodes:
			// like:
			//	- GetAssignmentStatement()
			//	- GetConditionalExpression()
			//	- GetArithmeticExpression()

			case token.Assign:
				stmt, err := a.GetStatement()
				if err != nil {
					return nil, err
				}

				a.Index++

				// For now just check for the separator here
				if a.Tokens[a.Index].Type == token.Separator {
					a.Index++
				}

				expr, err := a.GetExpression()
				if err != nil {
					return nil, err
				}

				a.Index++

				// For now just check for the separator here
				if a.Tokens[a.Index].Type == token.Separator {
					a.Index++
				}

				expr2, err := a.GetExpression()
				if err != nil {
					return nil, err
				}

				body, err := a.GetBlock()
				if err != nil {
					return nil, err
				}

				// FIXME: should make a new function for this
				return &Loop{
					Token: forToken,
					Type:  StdFor,
					Init:  stmt.(*Assignment),
					Cond:  expr,
					Post:  expr2,
					Body:  body,
				}, nil

			case token.Keyword:
				a.Index++
				preposition := a.Tokens[a.Index]

				a.Index++
				expr, err := a.GetExpression()
				if err != nil {
					return nil, err
				}
				fmt.Println("expr me", expr)

				iter, err := NewIterable(forToken, preposition, ident, expr)
				if err != nil {
					return nil, err
				}

				a.Index++

				body, err := a.GetBlock()
				if err != nil {
					return nil, err
				}

				prepType := ForIn
				if preposition.Value.String == "of" {
					prepType = ForOf
				} else if preposition.Value.String == "over" {
					prepType = ForOver
				}

				// FIXME: should make a new function for this
				return &Loop{
					Token: forToken,
					Type:  prepType,
					Iter:  iter,
					Body:  body,
				}, nil

			default:
				fmt.Println("preposition", a.Tokens[a.Index+1])
			}
		}

		return nil, errors.New("Could not parse loop")

	case token.Return:
		// For now just look for a single expression afterwards
		a.Index++
		expr, err := a.GetExpression()
		if err != nil {
			return nil, err
		}

		fmt.Println("return return")

		return NewReturn(token.Token{}, expr), nil

	default:
		return nil, errors.Errorf("Could not get statement from token: %+v", currentToken)
	}

	return nil, errors.Errorf("Could not deduce statement starting at: %+v", a.Tokens[a.Index])
}

// BuildAST builds an AST from the tokens provided by the lexer
func (a *ASTBuilder) BuildAST() (*Program, error) {
	p := NewProgram()

	// FIXME: Spoof this name for now
	file := NewFile("main.expr")

	for {
		// We know that the file can only consist of statements
		stmt, err := a.GetStatement()
		if err != nil {
			return nil, err
		}

		file.AddStatement(stmt)

		a.Index++

		if a.Index > len(a.Tokens)-1 {
			break
		}
	}

	p.AddFile(file)

	return p, nil
}

func CompressTokens(lexTokens []token.Token) ([]token.Token, error) {
	compressedTokens := []token.Token{}

	// alreadyChecked := false

	// Combine operators
	for i := 0; i < len(lexTokens); i++ {
		// fmt.Printf("%+v\n", lexTokens[i])

		currentToken := lexTokens[i]

		if i < len(lexTokens)-1 {
			nextToken := lexTokens[i+1]

			// This needs to be simplified
			if currentToken.Type == token.Assign || currentToken.Type == token.SecOp || currentToken.Type == token.PriOp && nextToken.Type == token.Assign || nextToken.Type == token.SecOp || nextToken.Type == token.PriOp {
				compressedToken, ok := token.TokenMap[currentToken.Value.String+nextToken.Value.String]
				// fmt.Println("added \"" + lexTokens[i].Value.String + nextToken.Value.String + "\"")
				if ok {
					compressedTokens = append(compressedTokens, compressedToken)
					i++

					// If we were able to combine the last two tokens and make a new one, mark it
					if i == len(lexTokens)-1 {
						// alreadyChecked = true
					}

					continue
				}
			}

			if currentToken.Type == token.GThan || currentToken.Type == token.LThan && nextToken.Type == token.Assign {
				compressedToken, ok := token.TokenMap[currentToken.Value.String+nextToken.Value.String]
				// fmt.Println("added \"" + lexTokens[i].Value.String + nextToken.Value.String + "\"")
				if ok {
					compressedTokens = append(compressedTokens, compressedToken)
					i++

					// If we were able to combine the last two tokens and make a new one, mark it
					if i == len(lexTokens)-1 {
						// alreadyChecked = true
					}

					continue
				}
			}
		}

		compressedTokens = append(compressedTokens, lexTokens[i])
	}

	// // If it hasn't been already checked and the last token is not a white space, then append it
	// if !alreadyChecked && lexTokens[len(lexTokens)-1].Type != token.Whitespace {
	// 	compressedTokens = append(compressedTokens, lexTokens[len(lexTokens)-1])
	// }

	compressedTokens2 := []token.Token{}
	// Combine array type tokens
	for i := 0; i < len(compressedTokens); i++ {
		currentToken := compressedTokens[i]

		if currentToken.Type == token.Ident {
			if strings.Contains(currentToken.Value.String, ".") {
				identSplit := strings.Split(currentToken.Value.String, ".")

				for i, is := range identSplit {
					if is == "" {
						// add an accessor
						compressedTokens2 = append(compressedTokens2,
							token.Token{
								Type: token.Accessor,
								Value: token.Value{
									Type:   "period",
									String: ".",
								},
							})
					} else {
						// add the ident
						compressedTokens2 = append(compressedTokens2,
							token.Token{
								Type: token.Ident,
								Value: token.Value{
									String: is,
								},
							})

						if i < len(identSplit)-1 {
							// add an accessor
							compressedTokens2 = append(compressedTokens2,
								token.Token{
									Type: token.Accessor,
									Value: token.Value{
										Type:   "period",
										String: ".",
									},
								})
						}
					}
				}

				continue
			}
		}

		compressedTokens2 = append(compressedTokens2, currentToken)
	}

	compressedTokens3 := []token.Token{}
	// Filter out the _un-needed_ white space
	for i := 0; i < len(compressedTokens2); i++ {
		if compressedTokens2[i].Type == token.Whitespace {
			continue
		}

		if compressedTokens2[i].Type == token.Return &&
			compressedTokens2[i+1].Value.String == "\n" {
			compressedTokens3 = append(compressedTokens3, compressedTokens2[i])
			compressedTokens3 = append(compressedTokens3, compressedTokens2[i+1])
			i++

			continue
		}

		compressedTokens3 = append(compressedTokens3, compressedTokens2[i])
	}

	return compressedTokens3, nil
}
