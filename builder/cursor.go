package builder

import (
	"fmt"

	token "github.com/scottshotgg/express-token"
)

// peek returns the current token without advancing.
func (b *Builder) peek() token.Token {
	return b.Tokens[b.Index]
}

// peekAt returns the token at Index+offset without advancing.
func (b *Builder) peekAt(offset int) token.Token {
	return b.Tokens[b.Index+offset]
}

// advance returns the current token and increments Index.
func (b *Builder) advance() token.Token {
	t := b.Tokens[b.Index]
	b.Index++
	return t
}

// expect asserts the current token has the given type, advances, and returns it.
func (b *Builder) expect(tokenType string) (token.Token, error) {
	if b.Tokens[b.Index].Type != tokenType {
		return token.Token{}, b.AppendTokenToError(
			fmt.Sprintf("expected %s but got %s", tokenType, b.Tokens[b.Index].Type))
	}
	return b.advance(), nil
}

// atEnd returns true if there are no more tokens to consume.
func (b *Builder) atEnd() bool {
	return b.Index >= len(b.Tokens)-1
}
