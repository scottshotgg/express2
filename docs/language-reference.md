# Express Language Reference

Express is a gradually-typed, C++-transpiled language that blends the ergonomics of Go and JavaScript with the performance of native code. It compiles to C++ and leverages the C++ standard library at runtime.

---

## Table of Contents

1. [Type System](#type-system)
2. [Variable Declarations](#variable-declarations)
3. [Operators](#operators)
4. [Control Flow](#control-flow)
5. [Functions](#functions)
6. [Data Structures](#data-structures)
7. [Module System](#module-system)
8. [Concurrency](#concurrency)
9. [Lifecycle Hooks](#lifecycle-hooks)
10. [Built-in Functions](#built-in-functions)
11. [C++ Interop](#c-interop)
12. [Planned Features](#planned-features)

---

## Type System

Express uses a **gradual type system**. You can be fully explicit, use inference, or go fully dynamic.

### Primitive Types

| Type     | Description              | Example            |
|----------|--------------------------|---------------------|
| `int`    | Integer                  | `int x = 5`        |
| `float`  | Floating point           | `float f = 3.14`   |
| `string` | String                   | `string s = "hi"`  |
| `bool`   | Boolean                  | `bool b = true`    |
| `char`   | Single character         | `char c = 'a'`     |

### `let` — Type Inference

`let` infers the type from the right-hand side. The variable is still statically typed — the compiler determines the type at compile time.

```
let x = 99        // int
let s = "hello"   // string
let b = true      // bool
```

### `var` — Dynamic Type

`var` is a true dynamic type. A `var` variable can hold any value and can change types at runtime. It transpiles to a variant/any type in C++.

```
var v              // uninitialized var
var v = 42         // holds an int
v = "now a string" // now holds a string
v = false          // now holds a bool
```

`var` arrays can hold mixed types:

```
var[] vv = [666, "something_here", false, 73.986622195]
```

`var` is useful for generic containers, maps with heterogeneous values, and function parameters that accept any type.

### Arrays and Vectors

Express has two collection types based on whether the size is known at compile time:

| Express syntax | Name       | C++ backing            | Notes                          |
|----------------|------------|------------------------|--------------------------------|
| `int[5] arr`   | **array**  | `int arr[5]`           | Fixed-size, stack-allocated    |
| `int[] arr`    | **vector** | `std::vector<int>`     | Dynamic, growable              |

```
int[] i = [1, 2, 3, 4, 5]       // vector — dynamic, growable
int[5] fixed                      // array — fixed-size (5 ints, zero-initialized)
var[] mixed = [1, "two", false]  // vector of mixed types (var elements)
```

"Vector" is used because the dynamic form IS `std::vector` under the hood. "Array" is used for the fixed-size C-style form.

### Pointers

Express supports C-style pointers with `&` (address-of) and `*` (dereference).

```
int a = 0
int* b = &a    // b points to a
int c = *b     // c = 0 (dereferenced)
*b = 42        // a is now 42
```

### Type Aliases

Create new type names with `type`:

```
type myInt = int
```

---

## Variable Declarations

### Explicit Type

```
int i = 10
string name = "express"
bool flag = true
float pi = 3.14159
```

### Inferred Type (`let`)

```
let count = 42
let greeting = "hello"
```

### Dynamic Type (`var`)

```
var anything = 99
anything = "now a string"
```

### Uninitialized Declarations

Variables can be declared without initialization. Express **auto-zero-initializes** all uninitialized typed declarations — there is no undefined behavior.

```
int i       // → int i = 0;
float f     // → float f = 0;
bool b      // → bool b = false;
char c      // → char c = '\0';
string s    // → std::string s;   (default-constructs to "")
var v       // → var v;           (default-constructs to nullType)
```

Zero values by type:

| Type     | Zero value | Notes                                      |
|----------|------------|--------------------------------------------|
| `int`    | `0`        |                                            |
| `float`  | `0`        |                                            |
| `bool`   | `false`    |                                            |
| `char`   | `'\0'`     | Null character                             |
| `string` | `""`       | `std::string` default-constructs to empty  |
| `var`    | `null`     | `var` class default-constructs to nullType |
| `map`    | `{}`       | `std::map` default-constructs to empty     |
| `int[5]` | `= {}`     | C-style array, aggregate-init to all zeros |
| `int[]`  | `[]`       | `std::vector` default-constructs to empty  |
| struct   | `= {}`     | All fields recursively zeroed              |
| pointer  | `nullptr`  |                                            |

`let` always requires an initializer — `let x` is a parse error.

### Assignment Operators

| Operator | Meaning        | Example                          |
|----------|----------------|----------------------------------|
| `=`      | Assignment     | `int i = 10`, `i = 20`          |
| `:`      | Key-value set  | `thing : "value"` (in maps)     |
| `+=`     | Add-assign     | `x += 5`  → `x = x + 5`        |
| `-=`     | Sub-assign     | `x -= 2`  → `x = x - 2`        |
| `*=`     | Mul-assign     | `x *= 3`  → `x = x * 3`        |
| `/=`     | Div-assign     | `x /= 4`  → `x = x / 4`        |

The `=` operator is used for all variable declarations and assignments. The `:` operator creates key-value pairs and is used inside map literals. Compound assignment operators (`+=`, `-=`, `*=`, `/=`) are desugared at parse time.

---

## Operators

### Arithmetic

| Operator | Description    |
|----------|----------------|
| `+`      | Addition       |
| `-`      | Subtraction    |
| `*`      | Multiplication |
| `/`      | Division       |
| `%`      | Modulo         |
| `^`      | Exponentiation |

### Comparison

| Operator | Description           |
|----------|-----------------------|
| `==`     | Equal to              |
| `<`      | Less than             |
| `>`      | Greater than          |
| `<=`     | Less than or equal    |
| `>=`     | Greater than or equal |

### Unary

| Operator | Description    |
|----------|----------------|
| `++`     | Increment      |
| `!`      | Logical NOT    |
| `&`      | Address-of     |
| `*`      | Dereference    |

### Member Access

| Operator | Description         | Example               |
|----------|---------------------|-----------------------|
| `.`      | Field/method access | `s.value.size()`      |
| `[]`     | Index access        | `arr[0]`, `m["key"]`  |

---

## Control Flow

### If / Else

```
if something > 10 {
  int x = 7
} else if some_flag {
  string y = "other"
} else {
  int z = 0
}
```

Conditions do not require parentheses.

### Standard For Loop

```
for int i = 0; i < 10; i++ {
  // body
}
```

### For-In (Keys)

Iterates over the **keys** (indices) of a collection:

```
for j in [1, 2, 3] {
  // j = 0, 1, 2 (indices)
}
```

### For-Of (Values)

Iterates over the **values** of a collection:

```
for j of [1, 2, 3] {
  // j = 1, 2, 3 (values)
}
```

### For-Over (Key + Value)

Iterates over **both key and value** of a collection:

```
// Single variable — receives a tuple of (key, value)
for item over [10, 20, 30] {
  // item is a (key, value) tuple
}

// Two variables — destructured key and value
for i, v over [10, 20, 30] {
  // i = 0, 1, 2 (keys)
  // v = 10, 20, 30 (values)
}
```

### Nested Loops

```
for j in [1, 2, 3] {
  for k in [4, 5, 6] {
    i = k
  }
}
```

---

## Functions

### Named Functions

Use `func` to define named functions:

```
func something(int i, string s) {
  int j = 6666666
}
```

With a return type:

```
func something(int i, string s) int {
  return 10
}
```

### Return Types

```
func echo(var i) var {
  return i + i
}
```

### The `main` Function

The `main` function is the entry point. It does not require a return type — `int` is automatically injected if not supplied.

```
func main() {
  // program starts here
}
```

### Functions Do Not Need Pre-Declaration

Functions can be called before they are defined in the source:

```
func main() {
  something()   // works even though something() is defined below
}

func something() {
  int j = 42
}
```

### Multiple Return Values (Planned)

```
let something, _ = callThisFunction()
```

The `_` discards unwanted return values (Go-style).

---

## Data Structures

### Structs

Structs define typed data structures with default values:

```
struct AnotherOne = {
  bool abc = true
}

struct myStruct = {
  int i = 10
  string something = "something"
  AnotherOne ayy
}
```

### Struct Instantiation

Override defaults when creating instances:

```
myStruct s = {
  i = 100 * 7 / 3
  something = "else"
  ayy = {
    abc = true
  }
}
```

### Maps

Maps use the `:` operator for key-value pairs. Keys must be strings — `map` is always `std::map<std::string, var>`:

```
map m = {
  "name" : "Alice"
  "city" : "Portland"
}
```

Values in map literals must be strings (the literal form infers `std::string` for both key and value).
For mixed-type values, use assignment form instead:

```
map m
m["name"]  = "Alice"
m["score"] = 42
m["active"] = true
```

An inline map literal `{ "k" : v }` can appear in expression context:

```
map m
m["nested"] = { "x" : 7 }
```

Uninitialized `map m` defaults to `std::map<std::string, var>` — safe to use without an explicit initializer.

#### Typed Maps

Use `map[K -> V]` to specify explicit key and value types. The `->` reads as "maps to":

```
map[string -> int] scores = { "Alice" : 95  "Bob" : 87 }
map[string -> bool] flags = {}
```

Transpiles to `std::map<std::string, int>` and `std::map<std::string, bool>` respectively.

#### N-Dimensional Maps

Multi-key maps use comma-separated key types followed by `->` and the value type. The last type is always the value; all preceding types are nested key dimensions, folded right-to-left:

```
map[string, string -> int] scores
scores["alice"]["math"] = 95
```

`map[string, string -> int]` desugars to `std::map<std::string, std::map<std::string, int>>`.

### Enums

```
enum {
  some
  one = some + 2
  here
}
```

Enum values auto-increment. Values can reference previous members.

### Tuples (Planned)

Tuples are fixed-size, heterogeneous collections:

```
(int, string) t = (5, "hello")
let a, b = t   // destructure
```

Used for multiple return values and `for-over` iteration.

### Union / Result Types (Planned)

Tagged unions for safe error handling:

```
union Result = {
  int ok
  string err
}
```

---

## Module System

### `package`

Declares the package name for the current file:

```
package something
```

### `import`

Imports another Express source file:

```
import "path/to/file.expr"
```

### `include`

Includes C/C++ headers for interop:

```
include (
  cl.h
  std
)
```

### `use` (Shelved)

Intended for importing and aliasing external code — extending implementations that are not originally yours:

```
// use os as os2
```

Currently not implemented.

---

## Concurrency

### `launch`

Launches a function call as a concurrent coroutine, similar to Go's `go` keyword:

```
launch something()
```

Launch with inline function:

```
launch func() {
  // concurrent work
}
```

Launch and capture a result (promise/future):

```
let result = launch func() string {
  return "async result"
}

// Pass a compatible function to handle the result
result.then(handler)

// Block and wait for the result
let value = result.future()
```

### `defer`

Defers execution of a statement until the enclosing scope exits:

```
defer Println("--- ENDING ---")
Println("--- STARTING ---")
// "--- ENDING ---" prints after the scope exits
```

### Channels (Planned)

Channels for communication between concurrent coroutines, similar to Go channels:

```
chan int c              // unbuffered channel
chan int c = chan(10)   // buffered channel with capacity 10
c <- 42                // send
let v = <- c           // receive
```

---

## Lifecycle Hooks

Lifecycle hooks control deferred behavior at specific exit points:

| Hook       | Description                                    |
|------------|------------------------------------------------|
| `onexit`   | Runs when the program exits                    |
| `onreturn` | Runs when the enclosing function returns       |
| `onleave`  | Runs when leaving the enclosing scope/block    |

These provide finer-grained control than `defer`, which runs at scope exit. Currently defined as keywords but not yet parsed.

---

## Built-in Functions

| Function       | Description                                                |
|----------------|------------------------------------------------------------|
| `Println(...)` | Variadic print with newline (like Go's `fmt.Println`)     |
| `sleep(n)`     | Sleep for `n` seconds                                      |
| `len(x)`       | Returns the number of elements in a vector or string       |

`len()` works on vectors (`int[]`, `string[]`, etc.) and `std::string` values.
It transpiles to `(x).size()`. It does **not** work on fixed-size arrays (`int[5]`) — their size is a compile-time constant.

```
int[] v = [1, 2, 3]
Println(len(v))        // 3

string s = "hello"
Println(len(s))        // 5
```

---

## C++ Interop

Express transpiles to C++ and can use C++ standard library types and methods directly:

```
// Using C++ STL methods on Express types
s.value.size()          // std::vector::size()
s.value.push_back(v)    // std::vector::push_back()
s.value.pop_back()      // std::vector::pop_back()
m[res].to_string()      // converts to std::string
```

Including C/C++ headers:

```
include (
  cl.h
  std
)

std.vector<cl.Platform> platform
cl.Platform.get(&platform)
```

---

## Comments

```
// Single-line comment

/* 
  Multi-line
  comment 
*/
```

---

## Planned Features

The following features are defined in the compiler infrastructure but not yet fully implemented:

### Likely Next
- **`for x over`** — iteration with key+value (keyword defined, not parsed yet)
- **Tuples** — for multiple return values and over-loop destructuring
- **Channels** — goroutine-style communication with `launch`
- **Union/Result types** — tagged unions for error handling

### Future
- **Interfaces** — Go-style structural typing
- **`pub`/`priv`** — visibility modifiers for fields and functions
- **`val`** — immutable variable declaration (like `const`)
- **`use X as Y`** — import aliasing and code extension
- **`async`/`await`** — alternative concurrency model (currently `launch` is primary)
- **Multiple return values** — Go-style `let a, b = func()`

### Operator Roadmap
- **Vector ops**: `.+`, `.-`, `.*`, `./` — element-wise operations
- **Spread**: `...` — expand collections
- **Range**: `..` — numeric ranges
- **Pipe**: `|>` — function composition
- **Null coalesce**: `?` — safe access
- **Unwrap**: `?:` — optional unwrapping

### Removed
- **`fn`** — was lambda keyword, removed
- **`select`** — removed
- **`switch`** — removed (may return later)
- **`object`** — shelved, struct covers the use case for now
