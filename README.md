# Express2

> A compiler for the Express programming language that transpiles to C++

![Go Version](https://img.shields.io/badge/go-1.26+-blue.svg)
[![Build Status](https://img.shields.io/badge/build-status-green.svg)](https://github.com/scottshotgg/express2/actions)
[![Coverage](https://img.shields.io/badge/coverage-coverage-yellow.svg)](https://github.com/scottshotgg/express2)

## Description

Express2 is a compiler for the Express programming language that transpiles to C++. It features a modern, expressive syntax with gradual typing, powerful type inference, and support for both low-level memory manipulation and high-level abstractions.

The project uses a multi-stage compilation pipeline:
1. **Lexer/Tokenizer** - Converts source code into tokens
2. **Parser (Pratt parser)** - Builds an Abstract Syntax Tree (AST)
3. **Semantic Analyzer** - Type checking and scope analysis
4. **Tree Flattener** - Converts high-level constructs to simpler forms
5. **Transpiler** - Generates C++ code from the AST
6. **C++ Compiler** - Compiles to native binary using clang++

## Features

| Feature | Description |
|---------|-------------|
| **Gradual Typing** | Supports `var` (dynamic), `let` (inferred), and explicit type declarations |
| **Structs** | First-class structured data with nested initialization |
| **Maps** | Associative arrays with key-value pairs |
| **Arrays** | Static and dynamic arrays with type safety |
| **Pointers** | Direct memory manipulation with `*` and `&` operators |
| **For-in/For-of** | Iterate over keys (indices) or values in collections |
| **Functions** | Named, lambda, variadic, and named return values |
| **Control Flow** | if/else if/else with block scoping |
| **Defer** | Execute functions when the current scope exits |
| **Enums** | Named integer constants with auto-increment |
| **String Support** | Multi-line, escaped characters, concatenation |
| **Operator Overloading** | Arithmetic, relational, bitwise operators |

## Requirements

- **Go 1.26+** - For the compiler toolchain
- **clang++** - For C++ compilation (C++2a standard)
- **make** - For building the project
- **clang-format** - Optional, for code formatting

## Building

The Express compiler is built as a Go application that produces a binary capable of compiling `.expr` source files.

### Initial Setup

```bash
# Set up the library path (if needed)
export EXPRPATH=/path/to/express2

# Build the compiler itself
go build -o express ./main.go
```

### Building a Program

```bash
# Build an Express program to a binary
./express build program.expr

# Build with specific output name
./express build program.expr -o output_binary

# Compile to C++ only (no binary generation)
./express build program.expr --emit-cpp
```

## Quick Start

### 1. Create a Simple Program

Create a file named `hello.expr`:

```express
func main() {
  Println("Hello, Express!")
  int x = 42
  Println("The answer is:", x)
}
```

### 2. Compile and Run

```bash
# Build and run in one step
./express run hello.expr

# Or build first, then run
./express build hello.expr
./hello
```

### 3. Use the Makefile

If available:

```bash
make build      # Build the compiler
make test       # Run tests
make clean      # Clean build artifacts
```

## Language Overview

### Variable Declarations

Express uses three keywords for variable declarations:

| Keyword | Type | Description |
|---------|------|-------------|
| `var` | Dynamic | Can hold any value, type can change |
| `let` | Inferred | Type inferred from initialization, immutable |
| `type` | Explicit | Explicit type declaration |

```express
func main() {
  var v = 42        // Dynamic typing
  v = "hello"       // Type can change
  
  let x = 99        // Type inferred as int
  // x = "wrong"     // Error: type mismatch
  
  int y = 100       // Explicit type
}
```

### Types

**Primitive Types:**
- `int` - Integer values
- `float` - Floating-point numbers
- `string` - Text strings
- `bool` - Boolean values (`true`/`false`)

**Composite Types:**

```express
// Arrays
int[] nums = [1, 2, 3]
var[] mixed = [1, "two", false]

// Maps
map m = {
  "key" : "value"
  42 : true
}

// Structs
struct Person = {
  string name = ""
  int age = 0
}

Person p = {
  name = "Alice"
  age = 30
}

// Pointers
int a = 42
int* ptr = &a
*ptr = 100
```

### Control Flow

```express
func main() {
  int x = 10
  
  if x > 5 {
    Println("x is greater than 5")
  } else if x == 5 {
    Println("x equals 5")
  } else {
    Println("x is less than 5")
  }
  
  // For-in: iterates over keys/indices
  for i in [10, 20, 30] {
    Println("Index:", i)
  }
  
  // For-of: iterates over values
  for value of [10, 20, 30] {
    Println("Value:", value)
  }
}
```

### Functions

```express
// Basic function
func greet(name string) {
  Println("Hello,", name)
}

// Function with return value
func add(a int, b int) int {
  return a + b
}

// Named return values
func compute() (sum int, product int) {
  sum = 10 + 20
  product = 10 * 20
  return
}

// Variadic function
func sum(nums int...) int {
  int total = 0
  for n of nums {
    total += n
  }
  return total
}

// Lambda (anonymous function)
func main() {
  var multiply = func(a int, b int) int {
    return a * b
  }
  Println(multiply(5, 3))  // 15
}
```

### Operators

**Arithmetic:**
```express
int a = 10 + 5   // Addition
int b = 10 - 3   // Subtraction
int c = 4 * 7    // Multiplication
int d = 20 / 4   // Division
int e = 10 % 3   // Modulo
int f = 2 ^ 8    // Exponent
```

**Relational:**
```express
bool a = 5 == 5  // Equal
bool b = 5 < 10  // Less than
bool c = 10 > 5  // Greater than
bool d = 5 <= 5  // Less than or equal
bool e = 10 >= 9 // Greater than or equal
```

**Logical:**
```express
bool and = true && false
bool or = true || false
bool not = !true
```

### Comments

```express
// Single-line comment

/* Multi-line comment
   spanning multiple lines */

func main() {
  int x = 10  // Inline comment
}
```

### Defer

```express
func main() {
  defer Println("First")
  defer Println("Second")
  Println("Main body")
  
  // Output:
  // Main body
  // Second
  // First
}
```

### Enums

```express
func main() {
  enum Colors {
    red
    green = red + 5
    blue
  }
  // red = 0, green = 5, blue = 6
}
```

## Project Structure

```
/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/
├── main.go                 # Entry point
├── cmd/                    # CLI commands (build, run)
│   ├── root.go
│   ├── build.go
│   └── run.go
├── compiler/               # Main compiler pipeline
│   ├── compiler.go
│   └── util.go
├── builder/                # Parser and AST builder
│   ├── expression.go       # Expression parsing (Pratt parser)
│   ├── statement.go        # Statement parsing
│   ├── type.go             # Type system
│   ├── types.go            # Type definitions
│   ├── scope_tree.go       # Scope management
│   └── semantic.go         # Semantic analysis
├── tree_flattener/         # Tree flattening
│   ├── tree_flattener.go   # Converts high-level to low-level
├── transpiler/             # C++ code generation
│   └── transpiler.go       # AST to C++ transpiler
├── semantic/               # Semantic analysis
│   ├── stack.go
│   ├── treeWalk.go
│   └── type_check.go
├── scope_tree/             # Scope tree implementation
│   └── scope_tree.go
├── test/                   # Test files
│   └── main/
│       └── source/
│           ├── types/      # Type tests
│           ├── loops/      # Loop tests
│           ├── control/    # Control flow tests
│           ├── operators/  # Operator tests
│           ├── function/   # Function tests
│           └── keywords/   # Keyword tests
├── vendor/                 # Dependencies
└── go.mod                  # Go module
```

## Compilation Pipeline

The Express compiler uses a multi-stage pipeline:

```
Source (.expr file)
    ↓
[1] Lexer/Tokenizer
    Converts source code into a stream of tokens
    ↓
[2] Pratt Parser
    Builds an Abstract Syntax Tree (AST) using recursive descent
    ↓
[3] Semantic Analyzer
    Performs type checking and scope analysis
    Builds type map and checks for semantic errors
    ↓
[4] Tree Flattener
    Converts high-level constructs to simpler forms
    (e.g., for-in loops become while loops with indices)
    ↓
[5] Transpiler
    Generates C++ code from the flattened AST
    Handles: types, functions, loops, expressions
    ↓
[6] C++ Compiler (clang++)
    Compiles the generated C++ to a native binary
    Links against libmill library
    ↓
Binary executable
```

### Pipeline Details

**Stage 1 - Tokenization:**
- Uses the `express-lex` library for lexical analysis
- Produces compressed tokens (e.g., `:=` instead of `: =`)

**Stage 2 - Parsing:**
- Implements a Pratt parser for expression parsing
- Handles operator precedence correctly
- Builds hierarchical AST nodes

**Stage 3 - Semantic Analysis:**
- Type checking for all expressions
- Scope tracking for variables
- Function signature validation

**Stage 4 - Tree Flattening:**
- Converts `for-in` loops to index-based `while` loops
- Converts `for-of` loops to value-based `while` loops
- Flattens nested blocks

**Stage 5 - Transpilation:**
- Maps Express types to C++ types
- Generates appropriate C++ code for each AST node
- Handles includes and imports

**Stage 6 - Compilation:**
- Uses `clang++` with `-std=c++2a`
- Links against `libmill` library
- Produces optimized binary with `-Ofast`

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./builder/...
go test ./compiler/...
go test ./tree_flattener/...
go test ./semantic/...

# Run tests with coverage
go test -cover ./...
```

### Test Structure

Tests are organized by component:

| Test Directory | Purpose |
|----------------|---------|
| `test/main/source/types/` | Type system tests |
| `test/main/source/loops/` | Loop construct tests |
| `test/main/source/control/` | Control flow tests |
| `test/main/source/operators/` | Operator tests |
| `test/main/source/function/` | Function tests |
| `test/main/source/keywords/` | Keyword behavior tests |
| `builder/` | Parser and AST builder tests |
| `compiler/` | Compiler pipeline tests |
| `tree_flattener/` | Tree flattener tests |
| `semantic/` | Semantic analysis tests |

### Example Test Flow

```bash
# Test a specific source file
go test -run TestCompile ./test/main/source/types/array.expr

# Run all tests with verbose output
go test -v ./...

# Check coverage for the builder package
go test -coverprofile=builder.cover ./builder/
go tool cover -html=builder.cover
```

## Contributing

Contributions are welcome. Please ensure:

1. All tests pass: `go test ./...`
2. Code is properly formatted
3. New features have corresponding tests
4. Follow the existing code style

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.
