package compiler_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/scottshotgg/express2/compiler"
)

// exampleTest describes one end-to-end example program.
type exampleTest struct {
	file     string
	expected string
}

var exampleTests = []exampleTest{
	{
		"01-hello-world.expr",
		"Hello, World!\n42\nThe value of x is: 10 and the language is: Express\n",
	},
	{
		"02-variables.expr",
		"Dynamic var (int): 42\nDynamic var (string): \"Now I'm a string\"\nDynamic var (bool): true\nInferred int: 100\nInferred string: Hello\nInferred bool: 0\nInferred float: 3.14\nExplicit int: 200\nExplicit string: Explicitly typed\nExplicit bool: 1\nExplicit float: 2.718\n",
	},
	{
		"03-arithmetic.expr",
		"a = 10\nb = 3\nsum = 13\ndifference = 7\nproduct = 30\nquotient = 3\nremainder = 1\ncounter = 1\n",
	},
	{
		"04-control-flow.expr",
		"The weather is pleasant\nx is exactly 10\nIndex: 0\nIndex: 1\nIndex: 2\n",
	},
	{
		"05-functions.expr",
		"Hello, Alice\nHello, Bob\n5 + 10 = 15\nFactorial of 5: 120\n",
	},
	{
		"06-arrays.expr",
		"First element: 1\nThird element: 3\nFor-in loop (indices):\nIndex: 0\nIndex: 1\nIndex: 2\nIndex: 3\nIndex: 4\nFor-of loop (values):\nValue: 1\nValue: 2\nValue: 3\nValue: 4\nValue: 5\n",
	},
	{
		"07-structs.expr",
		"Alice's name: Alice\nBob's age: 25\nAlice's new age: 31\nHello, Alice\nAge: 31\n",
	},
	{
		"08-maps.expr",
		"Red hex: FF0000\nBlue hex: 0000FF\nWhite hex: FFFFFF\n",
	},
	{
		"09-pointers.expr",
		"Original value: 42\nValue via pointer: 42\nAfter modifying via pointer: 100\nAfter swap: x = 10 y = 5\n",
	},
	{
		"10-complex.expr",
		"Name: Alice Grade: 92 Passing: 1\nName: Bob Grade: 55 Passing: 0\nName: Carol Grade: 78 Passing: 1\nAlice passed!\nBob did not pass.\nTotal score: 225\nAlice grade via pointer: 92\n",
	},
}

// TestExamples compiles and runs the 10 verified example programs end-to-end.
// Requires EXPRPATH to be set (the root of the Express installation).
//
// NOTE: tests run sequentially because the builder uses package-level globals
// (scopeTree, currentTree) that are not safe for concurrent use.
//
// NOTE: the compiler always places the binary next to the source file
// (c.path is set but unused). The test cleans up the produced binary and
// .cpp file after each sub-test.
func TestExamples(t *testing.T) {
	if os.Getenv("EXPRPATH") == "" {
		t.Skip("EXPRPATH not set — skipping integration tests")
	}

	examplesDir, err := filepath.Abs("../docs/examples")
	if err != nil {
		t.Fatalf("could not resolve examples dir: %v", err)
	}

	for _, tt := range exampleTests {
		t.Run(tt.file, func(t *testing.T) {
			// NOT t.Parallel() — builder has package-level globals

			srcPath := filepath.Join(examplesDir, tt.file)
			// The compiler always writes the binary next to the source file.
			binPath := strings.TrimSuffix(srcPath, ".expr")

			// Clean up the produced artifacts after this sub-test.
			t.Cleanup(func() {
				os.Remove(binPath)
				os.Remove(binPath + ".cpp")
			})

			c, err := compiler.New(binPath)
			if err != nil {
				t.Fatalf("compiler.New: %v", err)
			}

			if err := c.CompileFile(srcPath); err != nil {
				t.Fatalf("CompileFile(%s): %v", tt.file, err)
			}

			var buf bytes.Buffer
			cmd := exec.Command(binPath)
			cmd.Stdout = &buf
			if err := cmd.Run(); err != nil {
				t.Fatalf("running binary %s: %v", binPath, err)
			}

			if got := buf.String(); got != tt.expected {
				t.Errorf("output mismatch for %s\ngot:\n%s\nwant:\n%s", tt.file, got, tt.expected)
			}
		})
	}
}
