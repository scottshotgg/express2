package typecheck_test

import (
	"testing"

	typecheck "github.com/scottshotgg/express2/semantic"
)

func TestNewStack(t *testing.T) {
	s := typecheck.NewStack()
	if s == nil {
		t.Fatal("NewStack returned nil")
	}
	if s.Length() != 0 {
		t.Errorf("Length() = %d, want 0", s.Length())
	}
}

func TestPush(t *testing.T) {
	s := typecheck.NewStack()
	s.Push("hello")
	if s.Length() != 1 {
		t.Errorf("Length() = %d, want 1 after Push", s.Length())
	}
}

func TestPop(t *testing.T) {
	s := typecheck.NewStack()
	s.Push("hello")
	val, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop error: %v", err)
	}
	if val.(string) != "hello" {
		t.Errorf("Pop() = %v, want hello", val)
	}
	if s.Length() != 0 {
		t.Errorf("Length() = %d, want 0 after Pop", s.Length())
	}
}

func TestPeek(t *testing.T) {
	s := typecheck.NewStack()
	s.Push("world")
	val, err := s.Peek()
	if err != nil {
		t.Fatalf("Peek error: %v", err)
	}
	if val.(string) != "world" {
		t.Errorf("Peek() = %v, want world", val)
	}
	// Peek should not remove the item
	if s.Length() != 1 {
		t.Errorf("Length() = %d, want 1 after Peek", s.Length())
	}
}

func TestPopEmpty(t *testing.T) {
	s := typecheck.NewStack()
	_, err := s.Pop()
	if err == nil {
		t.Fatal("Pop on empty stack should return an error")
	}
	if err != typecheck.ErrEmptyStack {
		t.Errorf("Pop error = %v, want ErrEmptyStack", err)
	}
}

func TestPeekEmpty(t *testing.T) {
	s := typecheck.NewStack()
	_, err := s.Peek()
	if err == nil {
		t.Fatal("Peek on empty stack should return an error")
	}
	if err != typecheck.ErrEmptyStack {
		t.Errorf("Peek error = %v, want ErrEmptyStack", err)
	}
}

func TestPushMultiple(t *testing.T) {
	s := typecheck.NewStack()
	s.Push("a")
	s.Push("b")
	s.Push("c")
	if s.Length() != 3 {
		t.Errorf("Length() = %d, want 3", s.Length())
	}
	// Stack is LIFO: pop should return "c", then "b", then "a"
	c, _ := s.Pop()
	b, _ := s.Pop()
	a, _ := s.Pop()
	if c.(string) != "c" {
		t.Errorf("first pop = %v, want c", c)
	}
	if b.(string) != "b" {
		t.Errorf("second pop = %v, want b", b)
	}
	if a.(string) != "a" {
		t.Errorf("third pop = %v, want a", a)
	}
}

func TestPushPopMixed(t *testing.T) {
	s := typecheck.NewStack()
	s.Push(1)
	s.Push(2)
	v, err := s.Pop()
	if err != nil || v.(int) != 2 {
		t.Errorf("pop after push(1),push(2) = %v, want 2", v)
	}
	s.Push(3)
	v, err = s.Pop()
	if err != nil || v.(int) != 3 {
		t.Errorf("pop after push(3) = %v, want 3", v)
	}
	v, err = s.Pop()
	if err != nil || v.(int) != 1 {
		t.Errorf("pop remaining = %v, want 1", v)
	}
	if s.Length() != 0 {
		t.Errorf("Length() = %d, want 0", s.Length())
	}
}
