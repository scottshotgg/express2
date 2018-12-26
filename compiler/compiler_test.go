package compiler_test

import (
	"testing"

	"github.com/scottshotgg/express2/compiler"
)

func TestRun(t *testing.T) {
	var err = compiler.Run("test/test.expr")
	if err != nil {
		t.Fatalf("err %+v", err)
	}
}

func TestCompile(t *testing.T) {
	var err = compiler.Compile("test/test.expr")
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
