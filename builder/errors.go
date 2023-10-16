package builder

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	formatString = "; %+v"
)

var (
	ErrNotImplemented = errors.New("Not implemented")
	ErrMultDimArrInit = errors.New("Cannot use multiple expression inside array type initializer")
	ErrOutOfTokens    = errors.New("Out of tokens")
	ErrContinue       = errors.New("continue")
)

func (b *Builder) AppendTokenToError(errText string) error {
	var err = errors.New(errText)

	if b.Index < len(b.Tokens)-1 {
		return errors.Wrap(err, fmt.Sprintf("%+v", b.Tokens[b.Index]))
	}

	return errors.Wrap(err, "no token to print")
}
