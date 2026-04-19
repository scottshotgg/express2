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
		"The weather is pleasant\nx is exactly 10\nYou can drive legally.\nIndex: 0\nIndex: 1\nIndex: 2\nValue: 10\nValue: 20\nValue: 30\n( 0 , 0 )\n( 0 , 1 )\n( 1 , 0 )\n( 1 , 1 )\n1\n2\n1\n2\n4\n5\n",
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
		"Red hex: \"FF0000\"\nBlue hex: \"0000FF\"\nWhite hex: \"FFFFFF\"\n",
	},
	{
		"09-pointers.expr",
		"Original value: 42\nValue via pointer: 42\nAfter modifying via pointer: 100\nAfter swap: x = 10 y = 5\n",
	},
	{
		"10-complex.expr",
		"Name: Alice Grade: 92 Passing: 1\nName: Bob Grade: 55 Passing: 0\nName: Carol Grade: 78 Passing: 1\nAlice passed!\nBob did not pass.\nTotal score: 225\nAlice grade via pointer: 92\n",
	},
	{
		"11-enum.expr",
		"Red: 0\nGreen: 1\nBlue: 2\n",
	},
	{
		"12-defer.expr",
		"first\nsecond\ndeferred\n",
	},
	{
		"13-type-aliases.expr",
		"x: 10\ngreeting: Hello\nsum: 15\n",
	},
	{
		"14-gradual-typing.expr",
		"int: 10\nstring: \"hello\"\nbool: true\nlet int: 42\nlet string: world\n",
	},
	{
		"16-nested-structs.expr",
		"Name: Alice\nAge: 30\nCity: Portland\nCountry: USA\n",
	},
	{
		"18-strings.expr",
		"Hello, World!\nLength check: Hello, World! Hello, World!\nThe answer is: 42\n",
	},
	{
		"19-c-blocks.expr",
		"From C directly: Hello from C block!\nx=10, y=20, sum=30\nGoing to sleep...\nWoke up!\nDone sleeping.\nThat is how you import C!\n",
	},
	{
		"20-c-namespace.expr",
		"Hello from c.printf!\nSEEK_SET value: 0\nLength of \"Express\": 7\n",
	},
	{
		"21-zero-init.expr",
		"int: 0\nfloat: 0\nbool: 0\nstring: \nouter.name: \nouter.count: 0\nouter.inner.value: 0\narr[0]: 0\narr[2]: 0\nvec len: 0\nafter int: 42\nafter str: hello\nafter name: Express\n",
	},
	{
		"22-combinations.expr",
		"int[] len: 0\nfloat[] len: 0\nbool[] len: 0\nstring[] len: 0\nchar[] len: 0\nint[] after push: 7\nstring[] after push: hello\nfloat arr[0]: 0\nbool arr[0]: 0\nstring arr[0]: \nfloat arr[0] after: 1.5\nbool arr[0] after: 1\nstring arr[0] after: test\nstruct[] len: 0\nstruct[] len after push: 1\npoints[0].x: 3\npoints[0].y: 4\nstruct arr[0].x: 0\nstruct arr[0].x after: 10\nwithvec.label: \nwithvec.items len: 0\nwithvec.label after: data\nwithvec.items[0]: 99\nwitharr.name: \nwitharr.buf[0]: 0\nwitharr.buf[0] after: 55\nwithmap.label: \nwithmap.data[key]: \"value\"\nmap vector set: ok\nnested map set: ok\nmap[] len: 1\n",
	},
	{
		"23-map-every-type.expr",
		"int: 42\nfloat: 3.14\nbool: true\nstring: \"hello\"\nchar: 122\nvar: 99\nvec set: ok\nnested map set: ok\n",
	},
	{
		"24-for-over.expr",
		"0\n1\n2\n0 10\n1 20\n2 30\n",
	},
	{
		"25-typed-maps.expr",
		"95\n0\n",
	},
	{
		"26-nested-maps.expr",
		"95\n88\n77\n",
	},
	{
		"27-math.expr",
		"3\n4\n3\n",
	},
	{
		"28-fmt.expr",
		"hello world\ndone\n",
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
