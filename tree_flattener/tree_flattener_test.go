package tree_flattener

import (
	"testing"

	ast "github.com/scottshotgg/express-ast"
	lex "github.com/scottshotgg/express-lex"
	"github.com/scottshotgg/express2/builder"
	"github.com/scottshotgg/express2/pkg/logger"
	"github.com/scottshotgg/express2/test"
)

// getBuilderFromString creates a builder from a test string
func getBuilderFromString(source string) (*builder.Builder, error) {
	tokens, err := lex.New(source).Lex()
	if err != nil {
		return nil, err
	}

	// Compress tokens (i.e., `:` and `=` compress into `:=`)
	tokens, err = ast.CompressTokens(tokens)
	if err != nil {
		return nil, err
	}

	return builder.New(tokens, logger.Noop()), nil
}

// TestGetArrayType tests extracting array type from nodes
func TestGetArrayType(t *testing.T) {
	f := New()

	// Test with ident node - returns the ident value as the type
	identNode := &builder.Node{
		Type:  "ident",
		Value: "myArray",
	}
	arrayType, err := f.getArrayType(identNode)
	if err != nil {
		t.Fatalf("getArrayType ident error: %v", err)
	}
	if arrayType != "myArray" {
		t.Errorf("Expected arrayType to be 'myArray', got '%s'", arrayType)
	}

	// Test with array node containing int type nodes
	arrayNode := &builder.Node{
		Type: "array",
		Value: []*builder.Node{
			{
				Type: "type",
				Kind: "int",
			},
			{
				Type: "type",
				Kind: "int",
			},
			{
				Type: "type",
				Kind: "int",
			},
		},
	}
	arrayType, err = f.getArrayType(arrayNode)
	if err != nil {
		t.Fatalf("getArrayType array error: %v", err)
	}
	if arrayType != "int" {
		t.Errorf("Expected arrayType to be 'int', got '%s'", arrayType)
	}

	// Test with single element array
	singleArrayNode := &builder.Node{
		Type: "array",
		Value: []*builder.Node{
			{
				Type: "type",
				Kind: "string",
			},
		},
	}
	arrayType, err = f.getArrayType(singleArrayNode)
	if err != nil {
		t.Fatalf("getArrayType single error: %v", err)
	}
	if arrayType != "string" {
		t.Errorf("Expected arrayType to be 'string', got '%s'", arrayType)
	}

	// Test with empty array - now returns an error instead of os.Exit
	t.Run("emptyArrayReturnsError", func(t *testing.T) {
		emptyArrayNode := &builder.Node{
			Type:  "array",
			Value: []*builder.Node{},
		}
		_, err := f.getArrayType(emptyArrayNode)
		if err == nil {
			t.Errorf("Expected error for empty array, got nil")
		}
	})
}

// TestMakeLengthCall tests generating std::size calls
func TestMakeLengthCall(t *testing.T) {
	f := New()

	// Create a simple array ident node
	arrayNode := &builder.Node{
		Type:  "ident",
		Value: "myArray",
	}

	lengthCall := f.makeLengthCall(arrayNode)

	if lengthCall.Type != "call" {
		t.Errorf("Expected call type, got '%s'", lengthCall.Type)
	}

	value := lengthCall.Value.(*builder.Node)
	if value.Type != "ident" {
		t.Errorf("Expected ident type for value, got '%s'", value.Type)
	}
	if value.Value != "std::size" {
		t.Errorf("Expected 'std::size', got '%s'", value.Value)
	}

	args := lengthCall.Metadata["args"].(*builder.Node)
	if args.Type != "egroup" {
		t.Errorf("Expected egroup type for args, got '%s'", args.Type)
	}

	argsValues := args.Value.([]*builder.Node)
	if len(argsValues) != 1 {
		t.Errorf("Expected 1 arg, got %d", len(argsValues))
	}
	if argsValues[0].Value != "myArray" {
		t.Errorf("Expected 'myArray' as arg, got '%s'", argsValues[0].Value)
	}

	// Test with literal node
	literalNode := &builder.Node{
		Type:  "literal",
		Value: 42,
	}
	lengthCall = f.makeLengthCall(literalNode)
	args = lengthCall.Metadata["args"].(*builder.Node)
	argsValues = args.Value.([]*builder.Node)
	if argsValues[0].Value != 42 {
		t.Errorf("Expected 42 as arg, got '%v'", argsValues[0].Value)
	}
}

// TestMakeLTComp tests less-than comparison generation
func TestMakeLTComp(t *testing.T) {
	f := New()

	lhs := &builder.Node{
		Type:  "ident",
		Value: "i",
	}
	rhs := &builder.Node{
		Type:  "literal",
		Value: 10,
	}

	comp := f.makeLTComp(lhs, rhs)

	if comp.Type != "comp" {
		t.Errorf("Expected comp type, got '%s'", comp.Type)
	}
	if comp.Value != "<" {
		t.Errorf("Expected '<' as value, got '%s'", comp.Value)
	}
	if comp.Left.Value != "i" {
		t.Errorf("Expected 'i' as left, got '%s'", comp.Left.Value)
	}
	if comp.Right.Value != 10 {
		t.Errorf("Expected 10 as right, got '%v'", comp.Right.Value)
	}

	// Test with different types
	lhs2 := &builder.Node{
		Type:  "literal",
		Value: 0,
	}
	rhs2 := &builder.Node{
		Type:  "ident",
		Value: "len",
	}
	comp = f.makeLTComp(lhs2, rhs2)
	if comp.Left.Value != 0 {
		t.Errorf("Expected 0 as left, got '%v'", comp.Left.Value)
	}
	if comp.Right.Value != "len" {
		t.Errorf("Expected 'len' as right, got '%v'", comp.Right.Value)
	}
}

// TestMakeIncrementOp tests increment operation generation
func TestMakeIncrementOp(t *testing.T) {
	f := New()

	node := &builder.Node{
		Type:  "ident",
		Value: "i",
	}

	inc := f.makeIncrementOp(node)

	if inc.Type != "inc" {
		t.Errorf("Expected inc type, got '%s'", inc.Type)
	}
	if inc.Left.Value != "i" {
		t.Errorf("Expected 'i' as left, got '%s'", inc.Left.Value)
	}

	// Test with literal node
	literalNode := &builder.Node{
		Type:  "literal",
		Value: 5,
	}
	inc = f.makeIncrementOp(literalNode)
	if inc.Left.Value != 5 {
		t.Errorf("Expected 5 as left, got '%v'", inc.Left.Value)
	}
}

// TestFlattenNode tests flattening various node types
func TestFlattenNode(t *testing.T) {
	f := New()

	// Test with literal node (should not error, default case)
	literalNode := &builder.Node{
		Type:  "literal",
		Value: 42,
	}
	err := f.FlattenNode(literalNode)
	if err != nil {
		t.Errorf("FlattenNode failed for literal: %v", err)
	}

	// Test with block node
	blockNode := &builder.Node{
		Type: "block",
		Value: []*builder.Node{
			{
				Type:  "literal",
				Value: 1,
			},
			{
				Type:  "literal",
				Value: 2,
			},
		},
	}
	err = f.FlattenNode(blockNode)
	if err != nil {
		t.Errorf("FlattenNode failed for block: %v", err)
	}

	// Test with forin node
	forinTest := test.Tests[test.StatementTest]["forin"]
	b, err := getBuilderFromString(forinTest)
	if err != nil {
		t.Fatalf("Failed to parse forin test: %v", err)
	}

	node, err := b.ParseForPrepositionStatement()
	if err != nil {
		t.Fatalf("Failed to parse forin statement: %v", err)
	}

	err = f.FlattenNode(node)
	if err != nil {
		t.Errorf("FlattenNode failed for forin: %v", err)
	}

	// Test with forof node
	forofTest := test.Tests[test.StatementTest]["forof"]
	b, err = getBuilderFromString(forofTest)
	if err != nil {
		t.Fatalf("Failed to parse forof test: %v", err)
	}

	node, err = b.ParseForPrepositionStatement()
	if err != nil {
		t.Fatalf("Failed to parse forof statement: %v", err)
	}

	err = f.FlattenNode(node)
	if err != nil {
		t.Errorf("FlattenNode failed for forof: %v", err)
	}

	// Test with function node
	funcTest := test.Tests[test.StatementTest]["funcDef"]
	b, err = getBuilderFromString(funcTest)
	if err != nil {
		t.Fatalf("Failed to parse function test: %v", err)
	}

	node, err = b.ParseFunctionStatement()
	if err != nil {
		t.Fatalf("Failed to parse function statement: %v", err)
	}

	err = f.FlattenNode(node)
	if err != nil {
		t.Errorf("FlattenNode failed for function: %v", err)
	}
}

// TestFlattenForIn tests converting for-in to while loops with array copy
func TestFlattenForIn(t *testing.T) {
	f := New()

	// Parse a forin statement using test data
	forinTest := test.Tests[test.StatementTest]["forin"]
	b, err := getBuilderFromString(forinTest)
	if err != nil {
		t.Fatalf("Failed to parse forin test: %v", err)
	}

	node, err := b.ParseForPrepositionStatement()
	if err != nil {
		t.Fatalf("Failed to parse forin statement: %v", err)
	}

	result, err := f.FlattenForIn(node)
	if err != nil {
		t.Fatalf("FlattenForIn error: %v", err)
	}

	if result == nil {
		t.Fatalf("FlattenForIn returned nil")
	}

	if len(result) < 2 {
		t.Fatalf("Expected at least 2 statements (decl + while), got %d", len(result))
	}

	// FlattenForIn returns [keyVar, ...extraDecls, while]; the last element is the while loop.
	whileStmt := result[len(result)-1]
	if whileStmt.Type != "while" {
		t.Errorf("Expected while type, got '%s'", whileStmt.Type)
	}

	// Check the condition is a comparison
	if whileStmt.Left.Type != "comp" {
		t.Errorf("Expected comp type for condition, got '%s'", whileStmt.Left.Type)
	}
	if whileStmt.Left.Value != "<" {
		t.Errorf("Expected '<' in condition, got '%s'", whileStmt.Left.Value)
	}

	// Verify the loop body has statements
	whileValue := whileStmt.Value.(*builder.Node)
	if whileValue.Type != "block" {
		t.Errorf("Expected block type for while value, got '%s'", whileValue.Type)
	}

	stmts := whileValue.Value.([]*builder.Node)
	if len(stmts) < 1 {
		t.Errorf("Expected at least 1 statement in while body, got %d", len(stmts))
	}
}

// TestTransformIdentToDecl tests the transformIdentToDecl helper.
func TestTransformIdentToDecl(t *testing.T) {
	f := New()

	node := &builder.Node{Type: "ident", Value: "x"}

	// "int" returns a valid decl
	decl := f.transformIdentToDecl("int", 0, node)
	if decl == nil {
		t.Fatal("expected non-nil for int type")
	}
	if decl.Type != "decl" {
		t.Errorf("expected type 'decl', got %q", decl.Type)
	}

	// "auto" returns a valid decl
	autoDecl := f.transformIdentToDecl("auto", "hello", node)
	if autoDecl == nil {
		t.Fatal("expected non-nil for auto type")
	}
	if autoDecl.Type != "decl" {
		t.Errorf("expected type 'decl', got %q", autoDecl.Type)
	}

	// unsupported types return nil
	for _, unsupported := range []string{"float", "bool", "char", "string"} {
		result := f.transformIdentToDecl(unsupported, nil, node)
		if result != nil {
			t.Errorf("transformIdentToDecl(%q): expected nil, got %+v", unsupported, result)
		}
	}
}

// TestTransformArrayToDecl tests the transformArrayToDecl helper.
func TestTransformArrayToDecl(t *testing.T) {
	f := New()

	// Case: node.Value is a string — uses the string directly as ident name.
	namedNode := &builder.Node{
		Type:  "ident",
		Value: "myArr",
	}
	decl := f.transformArrayToDecl("int", namedNode)
	if decl == nil {
		t.Fatal("expected non-nil decl")
	}
	if decl.Left.Value.(string) != "myArr" {
		t.Errorf("expected ident 'myArr', got %q", decl.Left.Value)
	}

	// Case: node.Value is not a string — generates arr_N name.
	counterBefore := f.IdentCounter
	literalNode := &builder.Node{
		Type:  "array",
		Value: []*builder.Node{{Type: "literal", Kind: "int", Value: 1}},
	}
	decl2 := f.transformArrayToDecl("int", literalNode)
	if decl2 == nil {
		t.Fatal("expected non-nil decl for generated name")
	}
	if f.IdentCounter != counterBefore+1 {
		t.Errorf("IdentCounter not incremented: before=%d after=%d", counterBefore, f.IdentCounter)
	}
}

// TestFlattenForOf tests converting for-of to while loops with indexed access
func TestFlattenForOf(t *testing.T) {
	f := New()

	// Parse a forof statement using test data
	forofTest := test.Tests[test.StatementTest]["forof"]
	b, err := getBuilderFromString(forofTest)
	if err != nil {
		t.Fatalf("Failed to parse forof test: %v", err)
	}

	node, err := b.ParseForPrepositionStatement()
	if err != nil {
		t.Fatalf("Failed to parse forof statement: %v", err)
	}

	result, err := f.FlattenForOf(node)
	if err != nil {
		t.Fatalf("FlattenForOf error: %v", err)
	}

	if result == nil {
		t.Fatalf("FlattenForOf returned nil")
	}

	if len(result) < 2 {
		t.Fatalf("Expected at least 2 statements (decl + while), got %d", len(result))
	}

	// The result should start with an int index variable declaration.
	incVar := result[0]
	if incVar.Type != "decl" {
		t.Errorf("Expected decl type for incVar, got '%s'", incVar.Type)
	}

	// Check the value type (should be int)
	valueNode := incVar.Value.(*builder.Node)
	if valueNode.Value != "int" {
		t.Errorf("Expected 'int' as value type, got '%s'", valueNode.Value)
	}

	// The last element is the while loop.
	whileStmt := result[len(result)-1]
	if whileStmt.Type != "while" {
		t.Errorf("Expected while type, got '%s'", whileStmt.Type)
	}

	// Check the while loop body
	whileValue := whileStmt.Value.(*builder.Node)
	if whileValue.Type != "block" {
		t.Errorf("Expected block type for while value, got '%s'", whileValue.Type)
	}

	// The body should contain array assignment and increment
	stmts := whileValue.Value.([]*builder.Node)
	if len(stmts) < 2 {
		t.Errorf("Expected at least 2 statements in while body, got %d", len(stmts))
	}
}

// TestFlattenNode_FunctionCase exercises the "function" case in FlattenNode.
func TestFlattenNode_FunctionCase(t *testing.T) {
	f := New()

	funcNode := &builder.Node{
		Type: "function",
		Kind: "myFunc",
		Value: &builder.Node{
			Type:  "block",
			Value: []*builder.Node{},
		},
		Metadata: map[string]interface{}{},
	}

	err := f.FlattenNode(funcNode)
	if err != nil {
		t.Fatalf("FlattenNode(function) error: %v", err)
	}
}

// TestFlattenNode_DefaultCase exercises the default case in FlattenNode (no-op).
func TestFlattenNode_DefaultCase(t *testing.T) {
	f := New()

	// Any node type not in the switch → default → returns nil immediately
	n := &builder.Node{Type: "typedef", Value: []*builder.Node{}}
	err := f.FlattenNode(n)
	if err != nil {
		t.Fatalf("FlattenNode(default) error: %v", err)
	}
}

// TestFlatten_NonProgramNode exercises the error path in Flatten.
func TestFlatten_NonProgramNode(t *testing.T) {
	f := New()

	n := &builder.Node{Type: "block", Value: []*builder.Node{}}
	_, err := f.Flatten(n)
	if err == nil {
		t.Error("expected error for non-program node, got nil")
	}
}

// TestFlattenForIn_NilStart exercises the nil-start early-return path in FlattenForIn.
func TestFlattenForIn_NilStart(t *testing.T) {
	f := New()

	// Metadata["start"] is nil → FlattenForIn returns nil, nil immediately
	forinNode := &builder.Node{
		Type: "forin",
		Metadata: map[string]interface{}{
			"start": nil,
			// "end" is also nil/absent — won't be accessed
		},
		Value: &builder.Node{Type: "block", Value: []*builder.Node{}},
	}

	result, err := f.FlattenForIn(forinNode)
	if err != nil {
		t.Fatalf("FlattenForIn(nil start) error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result for nil start, got %v", result)
	}
}
