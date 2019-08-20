package compiler_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/scottshotgg/express2/compiler"
)

var (
	c   *compiler.Compiler
	err error
)

func init() {
	os.Setenv("EXPRPATH", "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2")

	c, err = compiler.New("output.something")
	if err != nil {
		// t.Fatalf("err %+v", err)
		fmt.Printf("err %+v", err)
	}
}

func TestRun(t *testing.T) {
	err = c.RunFile("test/simple_import.expr")
	if err != nil {
		t.Fatalf("err %+v", err)
	}
}

func TestCompile(t *testing.T) {
	err = c.CompileFile("test/test.expr")
	if err != nil {
		t.Fatalf("err %+v", err)
	}
}

// //////

// // Check the scope map for the variable name that was returned
// var node = b.ScopeTree.Get(ident.Value.(string))

// // If we didn't get anything from the scope tree then the assignment can't proceed
// // The Get method might need to change since how will index operations work ...
// if node == nil {
// 	// TODO: do other type checking here later
// 	return nil, errors.Errorf("Variable has not been declared yet: %+v", ident)
// }

// // Set the current ident to what we got from the scope tree
// *ident = *node.Left

// /////
