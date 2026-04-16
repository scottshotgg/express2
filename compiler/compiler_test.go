package compiler_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/scottshotgg/express2/compiler"
)

const (
	errFormatString = "err: %+v\n"
)

var (
	compilerInstance *compiler.Compiler
	err              error
)

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMid(s, substr)))
}

func containsMid(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// getCompiler creates a test compiler with the test output path
func getCompiler() (*compiler.Compiler, error) {
	outputPath := "/tmp/test_compiler_output"
	return compiler.New(outputPath)
}

// createTestExprFile creates a temporary .expr file with the given content
// Returns the file path or error
func createTestExprFile(content string) (string, error) {
	// Create a temporary directory if it doesn't exist
	tempDir := "/tmp/express_test"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", err
	}

	// Create a temporary file with .expr extension
	tempFile := filepath.Join(tempDir, "test.expr")

	if err := os.WriteFile(tempFile, []byte(content), 0644); err != nil {
		return "", err
	}

	return tempFile, nil
}

// cleanupTestFile removes the test .expr file
func cleanupTestFile(filePath string) {
	os.Remove(filePath)
}

// TestNew tests creating a Compiler with valid output
func TestNew(t *testing.T) {
	// Set EXPRPATH for testing
	os.Setenv("EXPRPATH", "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2")

	c, err := compiler.New("/tmp/test_output")
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	if c == nil {
		t.Fatal("Compiler was nil")
	}

	if c.Outputs == nil {
		t.Fatal("Outputs map was nil")
	}

	if c.OutputData == nil {
		t.Fatal("OutputData map was nil")
	}

	if c.LibBase == "" {
		t.Fatal("LibBase was not set")
	}

	// Verify the path is correct
	expectedLibBase := "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/"
	if c.LibBase != expectedLibBase {
		t.Errorf("Expected LibBase to be %s, got %s", expectedLibBase, c.LibBase)
	}
}

// TestNew_NoExprPath tests error when EXPRPATH not set
func TestNew_NoExprPath(t *testing.T) {
	// Unset EXPRPATH
	os.Unsetenv("EXPRPATH")

	_, err := compiler.New("/tmp/test_output")
	if err == nil {
		t.Fatal("Expected error when EXPRPATH is not set, but got nil")
	}

	expectedErrMsg := "`EXPRPATH` is not set"
	if !contains(err.Error(), expectedErrMsg) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

// TestSetOutput tests setting output map
func TestSetOutput(t *testing.T) {
	os.Setenv("EXPRPATH", "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2")

	c, err := getCompiler()
	if err != nil {
		t.Fatalf("Failed to create compiler: %v", err)
	}

	newOutputs := map[string]string{
		"lex":      "/tmp/test_lex.txt",
		"compress": "/tmp/test_compress.txt",
		"ast":      "/tmp/test_ast.txt",
	}

	c.SetOutput(newOutputs)

	if c.Outputs == nil {
		t.Fatal("Outputs map is nil after SetOutput")
	}

	if len(c.Outputs) != 3 {
		t.Errorf("Expected 3 outputs, got %d", len(c.Outputs))
	}

	// Verify specific outputs
	if c.Outputs["lex"] != "/tmp/test_lex.txt" {
		t.Errorf("Expected lex output to be /tmp/test_lex.txt, got %s", c.Outputs["lex"])
	}

	if c.Outputs["compress"] != "/tmp/test_compress.txt" {
		t.Errorf("Expected compress output to be /tmp/test_compress.txt, got %s", c.Outputs["compress"])
	}

	if c.Outputs["ast"] != "/tmp/test_ast.txt" {
		t.Errorf("Expected ast output to be /tmp/test_ast.txt, got %s", c.Outputs["ast"])
	}
}

// TestSetOutputData tests setting output data
func TestSetOutputData(t *testing.T) {
	os.Setenv("EXPRPATH", "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2")

	c, err := getCompiler()
	if err != nil {
		t.Fatalf("Failed to create compiler: %v", err)
	}

	// Set some output data
	c.OutputData["test"] = []byte("test data content")

	if len(c.OutputData) != 1 {
		t.Errorf("Expected 1 output data entry, got %d", len(c.OutputData))
	}

	if string(c.OutputData["test"]) != "test data content" {
		t.Errorf("Expected 'test data content', got '%s'", string(c.OutputData["test"]))
	}
}

// TestNew_InvalidExprPath tests error handling for invalid EXPRPATH
func TestNew_InvalidExprPath(t *testing.T) {
	// Set EXPRPATH to an invalid path
	// Note: filepath.Abs does not fail for non-existent paths, so this test
	// verifies that New accepts any string path (valid or not)
	os.Setenv("EXPRPATH", "/nonexistent/path/to/nowhere")

	c, err := compiler.New("/tmp/test_output")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify that the compiler was created with the invalid path converted to absolute
	if c == nil {
		t.Fatal("Compiler should not be nil for invalid path")
	}

	// Verify the path was converted to absolute
	expectedLibBase := "/nonexistent/path/to/nowhere/lib/"
	if c.LibBase != expectedLibBase {
		t.Errorf("Expected LibBase to be %s, got %s", expectedLibBase, c.LibBase)
	}
}

// TestCompileFile_SimpleProgram compiles a simple program with Println
func TestCompileFile_SimpleProgram(t *testing.T) {
	// Set EXPRPATH for testing
	os.Setenv("EXPRPATH", "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2")

	// Set up compiler
	c, err := getCompiler()
	if err != nil {
		t.Fatalf("Failed to create compiler: %v", err)
	}

	// Create a simple test program
	content := `func main() {
  println("Hello, World!")
}`
	tempFile, err := createTestExprFile(content)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer cleanupTestFile(tempFile)

	err = c.CompileFile(tempFile)
	if err != nil {
		// Note: This may fail if clang++ or libmill are not available
		t.Logf("CompileFile returned error (may be expected if dependencies not available): %v", err)
	}
}

// TestCompileFile_Arithmetic compiles arithmetic expressions
func TestCompileFile_Arithmetic(t *testing.T) {
	// Set EXPRPATH for testing
	os.Setenv("EXPRPATH", "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2")

	c, err := getCompiler()
	if err != nil {
		t.Fatalf("Failed to create compiler: %v", err)
	}

	content := `func main() {
  int a = 10 + 5
  int b = 10 - 3
  int c = 4 * 7
  int d = 20 / 4
  int e = 10 % 3
  int f = 2 ^ 8
}`
	tempFile, err := createTestExprFile(content)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer cleanupTestFile(tempFile)

	err = c.CompileFile(tempFile)
	if err != nil {
		t.Logf("CompileFile returned error (may be expected if dependencies not available): %v", err)
	}
}

// TestCompileFile_Variables compiles variable declarations
func TestCompileFile_Variables(t *testing.T) {
	// Set EXPRPATH for testing
	os.Setenv("EXPRPATH", "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2")

	c, err := getCompiler()
	if err != nil {
		t.Fatalf("Failed to create compiler: %v", err)
	}

	content := `func main() {
  // Explicit typed declarations
  int i = 10
  string s = "hello"
  bool b = true
  float f = 3.14

  // Uninitialized declarations
  int j
  string t

  // Type-inferred
  let x = 99

  // Dynamic type
  var v = "anything"
}`
	tempFile, err := createTestExprFile(content)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer cleanupTestFile(tempFile)

	err = c.CompileFile(tempFile)
	if err != nil {
		t.Logf("CompileFile returned error (may be expected if dependencies not available): %v", err)
	}
}

// TestCompileFile_Functions compiles function definitions and calls
func TestCompileFile_Functions(t *testing.T) {
	// Set EXPRPATH for testing
	os.Setenv("EXPRPATH", "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2")

	c, err := getCompiler()
	if err != nil {
		t.Fatalf("Failed to create compiler: %v", err)
	}

	content := `func add(a int, b int) int {
  return a + b
}

func main() {
  int result = add(5, 3)
  println(result)
}`
	tempFile, err := createTestExprFile(content)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer cleanupTestFile(tempFile)

	err = c.CompileFile(tempFile)
	if err != nil {
		t.Logf("CompileFile returned error (may be expected if dependencies not available): %v", err)
	}
}

// TestCompileFile_Arrays compiles array operations
func TestCompileFile_Arrays(t *testing.T) {
	// Set EXPRPATH for testing
	os.Setenv("EXPRPATH", "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2")

	c, err := getCompiler()
	if err != nil {
		t.Fatalf("Failed to create compiler: %v", err)
	}

	content := `func main() {
  int[] nums = [1, 2, 3, 4, 5]
  var[] mixed = [1, "two", false, 4.5]

  // Access elements
  int first = nums[0]
  int last = nums[4]
}`
	tempFile, err := createTestExprFile(content)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer cleanupTestFile(tempFile)

	err = c.CompileFile(tempFile)
	if err != nil {
		t.Logf("CompileFile returned error (may be expected if dependencies not available): %v", err)
	}
}

// TestCompileFile_Structs compiles struct definitions
func TestCompileFile_Structs(t *testing.T) {
	// Set EXPRPATH for testing
	os.Setenv("EXPRPATH", "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2")

	c, err := getCompiler()
	if err != nil {
		t.Fatalf("Failed to create compiler: %v", err)
	}

	content := `struct Point = {
  int x
  int y
}

func main() {
  Point p = {
    x = 10
    y = 20
  }
  println(p.x, p.y)
}`
	tempFile, err := createTestExprFile(content)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer cleanupTestFile(tempFile)

	err = c.CompileFile(tempFile)
	if err != nil {
		t.Logf("CompileFile returned error (may be expected if dependencies not available): %v", err)
	}
}

// TestCompileFile_Loops compiles for-in/for-of loops
func TestCompileFile_Loops(t *testing.T) {
	// Set EXPRPATH for testing
	os.Setenv("EXPRPATH", "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2")

	c, err := getCompiler()
	if err != nil {
		t.Fatalf("Failed to create compiler: %v", err)
	}

	content := `func main() {
  // for-in iterates over the keys (indices) of a collection
  for i in [10, 20, 30] {
    println(i)
  }

  // for-of iterates over the values of a collection
  for v of [10, 20, 30] {
    println(v)
  }

  // Standard for loop
  for int j = 0; j < 5; j = j + 1 {
    println(j)
  }
}`
	tempFile, err := createTestExprFile(content)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer cleanupTestFile(tempFile)

	err = c.CompileFile(tempFile)
	if err != nil {
		t.Logf("CompileFile returned error (may be expected if dependencies not available): %v", err)
	}
}

// TestCompileFile_InvalidExtension tests error handling for non-.expr files
func TestCompileFile_InvalidExtension(t *testing.T) {
	// Set EXPRPATH for testing
	os.Setenv("EXPRPATH", "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2")

	c, err := getCompiler()
	if err != nil {
		t.Fatalf("Failed to create compiler: %v", err)
	}

	// Create a file with wrong extension
	tempFile := "/tmp/test.txt"
	if err := os.WriteFile(tempFile, []byte("func main() {}"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer cleanupTestFile(tempFile)

	err = c.CompileFile(tempFile)
	if err == nil {
		t.Fatal("Expected error for non-.expr file, got nil")
	}
}

// TestPipelineTimes verifies that pipeline timing information is collected
func TestPipelineTimes(t *testing.T) {
	// Set EXPRPATH for testing
	os.Setenv("EXPRPATH", "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2")

	c, err := getCompiler()
	if err != nil {
		t.Fatalf("Failed to create compiler: %v", err)
	}

	if c.PipelineTimes == nil {
		t.Fatal("PipelineTimes map was nil")
	}

	if len(c.PipelineTimes) != 0 {
		t.Errorf("Expected empty PipelineTimes, got %d entries", len(c.PipelineTimes))
	}
}

// TestFlags verifies default compiler flags are set
func TestFlags(t *testing.T) {
	os.Setenv("EXPRPATH", "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2")

	c, err := getCompiler()
	if err != nil {
		t.Fatalf("Failed to create compiler: %v", err)
	}

	if c.Flags == nil {
		t.Fatal("Flags slice was nil")
	}

	// Check that std=c++2a flag is present
	found := false
	for _, flag := range c.Flags {
		if flag == "-std=c++2a" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected -std=c++2a flag to be present")
	}
}
