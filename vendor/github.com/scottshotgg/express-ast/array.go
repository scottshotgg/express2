package ast

import (
	"github.com/scottshotgg/express-token"
)

// type ArrayType int

// const (
// 	_ ArrayType = iota

// 	Homogenous

// 	Heterogeneous
// )

// Array represents array type data structures
type Array struct {
	Token token.Token
	// How will this act with `var` elements?
	TypeOf      *Type
	ElementType *Type
	Length      int
	Elements    []Expression
	Homogenous  bool
}

// TODO: this should implement iterable.... no?

func (a *Array) expressionNode() {}

// TokenLiteral returns the literal value of the token
func (a *Array) TokenLiteral() token.Token { return a.Token }

// Type implements Literal
func (a *Array) Type() *Type { return a.TypeOf }

func (a *Array) Kind() NodeType { return ArrayNode }

func (a *Array) String() string {
	// FIXME: just doing this to get it to compile
	// return fmt.Sprintf("%+v", *a)
	// FIXME: this space is a hack kinda
	arrayString := "{ "

	for _, elem := range a.Elements {
		arrayString += elem.String() + ","
	}

	return arrayString[:len(arrayString)-1] + "}"
}

func NewArray(t token.Token, elements []Expression) *Array {
	homogenous := true

	var typeOf *Type
	if len(elements) > 0 {
		typeOf = elements[0].Type()
		for _, e := range elements[1:] {
			// if e.Kind() == IdentNode {
			// 	continue
			// }

			// Compare to figure out if we need to upgrade the array type or not
			if e.Type().Type != typeOf.Type && e.Type().UpgradesTo != typeOf.Type {
				// if the collected types can upgrade to the expression type
				if e.Type().Type == typeOf.UpgradesTo {
					typeOf = e.Type()
				} else {
					homogenous = false
					break
				}
			}
		}
	}

	typeOf.Array = true

	return &Array{
		Token:       t,
		TypeOf:      NewArrayType(typeOf, homogenous),
		ElementType: typeOf,
		Length:      len(elements),
		Elements:    elements,
		Homogenous:  homogenous,
	}
}
