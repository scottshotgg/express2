package ast

import token "github.com/scottshotgg/express-token"

// Function represents the following form:
// [ `func` | `fn` ] [ ident ] [ group ] { group } [ block ]
type Function struct {
	Lambda    bool
	Async     bool
	Token     token.Token
	Ident     *Ident
	Arguments *Group
	Returns   *Group
	Body      *Block
}

// Implement statement
func (f *Function) statementNode() {}

// Implement expression
func (f *Function) expressionNode() {}

// TokenLiteral returns the literal value of the token
func (f *Function) TokenLiteral() token.Token { return f.Token }

// Type implements expression so that functions can be assigned to idents
func (f *Function) Type() *Type { return NewFunctionType() }

func (f *Function) Kind() NodeType { return FunctionNode }

// // TODO: we are only supporting on return for now
// func (f Function) String() string {
// 	return1 := "void"
// 	// Don't know if we need this, just being cautious rn
// 	if f.Returns != nil && f.Returns.Elements[0] != nil {
// 		return1 = f.Returns.Elements[0].(*Ident).Name
// 	}

// 	// FIXME: put all the functions at the top of the C++ file
// 	return return1 + " " + f.Ident.Name + f.Arguments.String() + f.Body.String()
// }

// TODO: we are only supporting on return for now
// std::function
// func (f Function) String() string {
// 	return1 := "void"
// 	// Don't know if we need this, just being cautious rn
// 	if f.Returns != nil && f.Returns.Elements[0] != nil {
// 		return1 = f.Returns.Elements[0].(*Ident).Name
// 	}

// 	// FIXME: put all the functions at the top of the C++ file
// 	return "std::function<" + return1 + f.Arguments.String() + ">" + f.Ident.Name + "= []" + f.Arguments.String() + f.Body.String()
// }

func (f Function) String() string {
	return1 := "void"
	// Don't know if we need this, just being cautious rn
	if f.Returns != nil && f.Returns.Elements[0] != nil {
		return1 = f.Returns.Elements[0].(*Ident).Name
	}

	// TODO: should probably check the returns and arguments in here
	if f.Ident.Name == "main" {
		return "int main()" + f.Body.String()
	}

	// FIXME: put all the functions at the top of the C++ file
	return return1 + " " + f.Ident.Name + f.Arguments.String() + f.Body.String()
}

func NewFunction(ft, it token.Token, args *Group, body *Block) (*Function, error) {
	ident, err := NewIdent(it, "")
	if err != nil {
		return nil, err
	}

	return &Function{
		Token:     ft,
		Ident:     ident,
		Arguments: args,
		Body:      body,
	}, nil
}
