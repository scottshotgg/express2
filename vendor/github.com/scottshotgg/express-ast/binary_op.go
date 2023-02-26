package ast

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	token "github.com/scottshotgg/express-token"
)

// BinaryOpType encompasses the types of binary operations
type BinaryOpType int

const (
	// AdditionBinaryOp is the + operator
	AdditionBinaryOp BinaryOpType = iota + 1

	// SubtractionBinaryOp is the - operator
	SubtractionBinaryOp

	// MultiplicationBinaryOp is the * operator
	MultiplicationBinaryOp

	// DivisionBinaryOp is the / operator
	DivisionBinaryOp
)

// type BinaryOp interface {
// 	Expression

// 	Type() BinaryOpType
// 	Right() *Expression
// 	Left() *Expression
// 	Evaluate() *Literal
// }

// BinaryOperation represents the following form:
// [ expression ] [ binary_op ] [ expression ]
type BinaryOperation struct {
	Token     token.Token
	Op        BinaryOpType
	LeftNode  Expression
	RightNode Expression
	Value     Literal
}

// func (b *BinaryOperation) Type() *Expression     { return b.Kind }
// func (b *BinaryOperation) Right() *Expression    { return b.RightExpr }
// func (b *BinaryOperation) Left() *Expression     { return b.LeftExpr }
// func (b *BinaryOperation) Evaluate() *Expression { return b.Value }

func (b *BinaryOperation) expressionNode() {}

// TokenLiteral returns the literal value of the token
func (b *BinaryOperation) TokenLiteral() token.Token { return b.Token }

func (b *BinaryOperation) Kind() NodeType { return BinaryOperationNode }

func (b *BinaryOperation) String() string {
	// FIXME: just doing this to get it to compile
	return fmt.Sprintf("%+v", *b)
}

func (b *BinaryOperation) Type() *Type {
	// FIXME: this needs some thought
	fmt.Println("Type() for binary operations is not implemented yet")
	os.Exit(9)

	return nil
}

// NewBinaryOperation returns a BinaryOperation with the evaluation value
func NewBinaryOperation(t token.Token, binOpString string, l Expression, r Expression) (*BinaryOperation, error) {
	var k BinaryOpType

	switch binOpString {
	case "+":
		k = AdditionBinaryOp

	case "-":
		k = SubtractionBinaryOp

	case "*":
		k = MultiplicationBinaryOp

	case "/":
		k = DivisionBinaryOp

	default:
		return nil, errors.Errorf("Could not decifer operation from supplied operand: %s", binOpString)
	}

	return &BinaryOperation{
		Token:     t,
		Op:        k,
		LeftNode:  l,
		RightNode: r,
	}, nil
}
