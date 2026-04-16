package builder

import (
	"testing"

	ast "github.com/scottshotgg/express-ast"
	lex "github.com/scottshotgg/express-lex"
	token "github.com/scottshotgg/express-token"
	"github.com/scottshotgg/express2/pkg/logger"
)

func newBuilderFromSource(src string) (*Builder, error) {
	tokens, err := lex.New(src).Lex()
	if err != nil {
		return nil, err
	}
	tokens, err = ast.CompressTokens(tokens)
	if err != nil {
		return nil, err
	}
	return New(tokens, logger.Noop()), nil
}

func TestCursorPeek(t *testing.T) {
	b, err := newBuilderFromSource("7")
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	tok := b.peek()
	if tok.Type != token.Literal {
		t.Errorf("peek().Type = %q, want %q", tok.Type, token.Literal)
	}
	if b.Index != 0 {
		t.Errorf("peek() advanced Index to %d, want 0", b.Index)
	}
}

func TestCursorPeekAt(t *testing.T) {
	b, err := newBuilderFromSource("7 + 8")
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	tok0 := b.peekAt(0)
	tok1 := b.peekAt(1)
	if tok0.Type != token.Literal {
		t.Errorf("peekAt(0).Type = %q, want Literal", tok0.Type)
	}
	if tok1.Type != token.SecOp {
		t.Errorf("peekAt(1).Type = %q, want SecOp (+)", tok1.Type)
	}
	if b.Index != 0 {
		t.Errorf("peekAt() changed Index to %d, want 0", b.Index)
	}
}

func TestCursorAdvance(t *testing.T) {
	b, err := newBuilderFromSource("7 + 8")
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	if b.Index != 0 {
		t.Fatalf("initial Index = %d, want 0", b.Index)
	}
	tok := b.advance()
	if tok.Type != token.Literal {
		t.Errorf("advance().Type = %q, want Literal", tok.Type)
	}
	if b.Index != 1 {
		t.Errorf("after advance, Index = %d, want 1", b.Index)
	}
}

func TestCursorAdvanceMultiple(t *testing.T) {
	b, err := newBuilderFromSource("7 + 8")
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	b.advance() // index → 1
	b.advance() // index → 2
	if b.Index != 2 {
		t.Errorf("after 2 advances, Index = %d, want 2", b.Index)
	}
}

func TestCursorExpect_Success(t *testing.T) {
	b, err := newBuilderFromSource("7")
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	tok, err := b.expect(token.Literal)
	if err != nil {
		t.Fatalf("expect(Literal) error: %v", err)
	}
	if tok.Type != token.Literal {
		t.Errorf("returned token Type = %q, want Literal", tok.Type)
	}
	if b.Index != 1 {
		t.Errorf("after expect, Index = %d, want 1", b.Index)
	}
}

func TestCursorExpect_Failure(t *testing.T) {
	b, err := newBuilderFromSource("7")
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	_, err = b.expect(token.Ident)
	if err == nil {
		t.Fatal("expect(Ident) should have returned an error for a Literal token")
	}
	if b.Index != 0 {
		t.Errorf("after failed expect, Index = %d, want 0 (no advance)", b.Index)
	}
}

func TestCursorAtEnd(t *testing.T) {
	b, err := newBuilderFromSource("7 + 8")
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	if b.atEnd() {
		t.Error("atEnd() = true at Index=0, want false")
	}
	// Advance to the last token
	for b.Index < len(b.Tokens)-1 {
		b.advance()
	}
	if !b.atEnd() {
		t.Errorf("atEnd() = false at Index=%d (len=%d), want true", b.Index, len(b.Tokens))
	}
}
