package typecheck_test

import (
	"fmt"
	"testing"

	"github.com/scottshotgg/express2/typecheck"
)

var (
	stack     *typecheck.Stack
	testScope = typecheck.Scope{
		"i": NewVariable("i", 6, typecheck.INT),
		// "a": typecheck.NewVariable("a", "hey its me", typecheck.STRING),
	}
)

func TestNewStack(t *testing.T) {
	stack = typecheck.NewStack()
	fmt.Printf("Stack: %+v\n", stack)
	fmt.Println()
}

func TestPush(t *testing.T) {
	TestNewStack(t)

	for k, v := range []string{"a", "b", "c"} {
		stack.Push(typecheck.Scope{
			v: typecheck.NewVariable(v, k, typecheck.INT),
		})
	}
	// stack.Push(testScope)
	// stack.Push(testScope)
	fmt.Printf("Stack: %+v\n", stack)
	fmt.Println()
}

func TestPop(t *testing.T) {
	TestPush(t)

	fmt.Printf("Stack: %+v\n", stack)

	pop, err := stack.Pop()
	if err != nil {
		fmt.Println("failed")
		t.Fail()
		return
	}

	fmt.Printf("Pop: %+v\n", pop)
	fmt.Printf("Stack: %+v\n", stack)
	fmt.Println()
}

func TestPeek(t *testing.T) {
	TestPush(t)

	peek, err := stack.Peek()
	if err != nil {
		t.Fail()
		return
	}

	fmt.Printf("Peek: %+v\n", peek)
	fmt.Printf("Stack: %+v\n", stack)
	fmt.Println()
}
