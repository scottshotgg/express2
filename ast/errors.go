package ast

import "github.com/pkg/errors"

const (
	formatString = "; %+v"
)

var (
	ErrNotImplemented = errors.New("Not implemented")
	ErrMultDimArrInit = errors.New("Cannot use multiple expression inside array type initializer")
	ErrOutOfTokens    = errors.New("Out of tokens")
)

func (b *Builder) AppendTokenToError(errText string) (*Node, error) {
	if b.Index < len(b.Tokens)-1 {
		return nil, errors.Errorf(errText+formatString, b.Tokens[b.Index])
	}

	return nil, errors.New(errText + "; No token to print")
}
