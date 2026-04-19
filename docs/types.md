# Express Type Reference

Complete reference for every type in Express: syntax, C++ backing, zero value, and usage notes.

---

## Primitive Types

| Express syntax | Canonical name | C++ backing  | Zero value | Notes                        |
|----------------|----------------|--------------|------------|------------------------------|
| `int x`        | int            | `int`        | `0`        | Platform-width integer       |
| `float x`      | float          | `double`     | `0`        | Double-precision float       |
| `bool x`       | bool           | `bool`       | `false`    |                              |
| `char x`       | char           | `char`       | `'\0'`     | Single character             |
| `string x`     | string         | `std::string`| `""`       | Default-constructs to empty  |

### Examples

```
int count = 42
float pi = 3.14159
bool active = true
char initial = 'A'
string name = "express"
```

---

## Gradual Typing

Express has three tiers of typing, from fully static to fully dynamic.

### `int` / `float` / `bool` / `char` / `string` — Explicit

The type is fixed at declaration and enforced throughout.

### `let` — Type Inference

`let` infers the type from the right-hand side. The variable is statically typed — the compiler determines the type once at compile time. Requires an initializer.

```
let x = 99        // inferred: int
let s = "hello"   // inferred: string
let b = true      // inferred: bool
```

`let` without an initializer is a **parse error**.

### `var` — Dynamic Type

`var` is a true dynamic type. A `var` variable can hold any value and change types at runtime. Transpiles to a `var` class (variant/any) in C++. Default-constructs to `nullType`.

```
var v              // v = null (nullType)
var v = 42         // v holds int 42
v = "now a string" // v now holds string
v = false          // v now holds bool
```

| Feature           | `int`/`float`/etc | `let`          | `var`          |
|-------------------|-------------------|----------------|----------------|
| Type known at...  | declaration       | compile time   | runtime        |
| Can change type   | no                | no             | yes            |
| Requires init     | no                | yes            | no             |
| C++ backing       | native type       | native type    | `var` class    |

---

## Collections

### Vector — `int[]`

**Dynamic, growable array.** Backed by `std::vector<T>`. Size is not known at compile time.

| Express syntax | Canonical name | C++ backing          | Zero value | `len()` |
|----------------|----------------|----------------------|------------|---------|
| `int[] arr`    | vector         | `std::vector<int>`   | empty `[]` | yes     |
| `string[] arr` | vector         | `std::vector<std::string>` | empty `[]` | yes |
| `bool[] arr`   | vector         | `std::vector<bool>`  | empty `[]` | yes     |
| `char[] arr`   | vector         | `std::vector<char>`  | empty `[]` | yes     |
| `var[] arr`    | vector         | `std::vector<var>`   | empty `[]` | yes     |

```
int[] nums = [1, 2, 3, 4, 5]     // vector of ints
var[] mixed = [1, "two", false]   // vector of var (mixed types)

int[] v
v.push_back(42)
Println(len(v))    // 1
```

`std::vector` default-constructs to empty — no action needed for declarations without an initializer. `= []` is also accepted as an explicit empty initializer.

**Note:** `T[]` syntax only works for primitive types and `var`. User-defined struct types (e.g. `Point[]`) are not yet supported — declare the vector inside a struct field instead.

### Array — `int[N]`

**Fixed-size array.** Backed by a C-style array `T name[N]`. Size is fixed at compile time.

| Express syntax | Canonical name | C++ backing         | Zero value | `len()` |
|----------------|----------------|---------------------|------------|---------|
| `int[5] arr`   | array          | `int arr[5]`        | `= {}`     | no      |
| `char[8] arr`  | array          | `char arr[8]`       | `= {}`     | no      |
| `string[3] arr`| array          | `std::string arr[3]`| `= {}`     | no      |

```
int[5] buf             // buf[0..4] all = 0  (zero-initialized)
char[256] scratch      // scratch[0..255] all = '\0'
```

Uninitialized arrays and arrays with `= {}` are aggregate-initialized — all elements zero. `len()` is **not** supported; the size is a compile-time constant known from the declaration.

**Naming rationale:** "vector" because the dynamic form IS `std::vector`. "array" for the fixed C-style form. "slice" (Go-only jargon) and "list" (linked-list connotations) were considered and rejected.

### Map

**Key-value store.** Keys are always `std::string`. Non-string keys require explicit typed maps (a planned future feature).

| Express form     | C++ backing                          | Zero value |
|------------------|--------------------------------------|------------|
| `map m` (uninit) | `std::map<std::string, var>`         | empty `{}` |
| `map m = { ... }`| `std::map<std::string, std::string>` | —          |

```
// Literal form — values must be strings
map m = {
  "name" : "Alice"
  "city" : "Portland"
}
```

```
// Assignment form — values can be any type (stored as var)
map m
m["name"]  = "Alice"
m["score"] = 42
m["flag"]  = true
```

`std::map` default-constructs to empty — safe to use uninitialized.

An inline map literal `{ "k" : v }` can be used in expression context (stored as `var` via `initializer_list`):

```
m["nested"] = { "x" : 7 }
```

---

## Other Types

### Struct

Typed data structures with named fields and default values.

```
struct myStruct = {
  int i = 10
  string name = "default"
}

myStruct s = {
  i = 100
  name = "override"
}
```

Uninitialized struct declarations emit `= {}` — aggregate-init, all fields recursively zeroed to their types' zero values.

### Pointer

C-style pointer. Backed by a raw C++ pointer `T*`.

| Express syntax | C++ backing | Zero value  |
|----------------|-------------|-------------|
| `int* p`       | `int*`      | `nullptr`   |

```
int a = 0
int* b = &a    // b points to a
int c = *b     // c = 0 (dereferenced)
*b = 42        // a is now 42
```

Uninitialized pointer declarations are zero-initialized to `nullptr`.

### Enum

```
enum {
  some
  one = some + 2
  here
}
```

Enum values auto-increment from 0. Values can reference previous members.

---

## Zero Values

Complete table of zero values for uninitialized declarations:

| Type       | Zero value  | C++ emission              |
|------------|-------------|---------------------------|
| `int`      | `0`         | `int x = 0;`              |
| `float`    | `0`         | `double x = 0;`           |
| `bool`     | `false`     | `bool x = false;`         |
| `char`     | `'\0'`      | `char x = '\0';`          |
| `string`   | `""`        | `std::string x;`          |
| `var`      | `null`      | `var x;`                  |
| `map`      | `{}`        | `std::map<std::string, var> x;` |
| `int[]`    | `[]`        | `std::vector<int> x;`     |
| `int[5]`   | `= {}`      | `int x[5] = {};`          |
| struct     | `= {}`      | `MyStruct x = {};`        |
| pointer    | `nullptr`   | `int* x = nullptr;`       |

---

## Uninitialized Declarations

Express **auto-zero-initializes** all typed declarations without an explicit initializer. There is no undefined behavior.

```
int x           // x = 0
bool active     // active = false
char ch         // ch = '\0'
int[5] buf      // buf[0..4] = 0
```

**Design rationale:** Matches Go's zero-value philosophy. Eliminates C++ UB from uninitialized reads. `let` requires an initializer (parse error if omitted). `var` has a defined runtime default (`nullType`/`nullptr`).

**Future:** If definite-assignment analysis is added, the compiler could instead require an explicit initializer and report an error on reads before write. This would catch "forgot to assign a meaningful value" bugs at compile time. For now, auto-zero-init is the chosen tradeoff.

---

## Type Combinations

Verified results from `docs/examples/22-combinations.expr`. Status reflects whether the combination compiles and produces correct output.

### Struct field types

| Field type         | Status | Notes                                      |
|--------------------|--------|--------------------------------------------|
| `int`              | ✓      |                                            |
| `float`            | ✓      |                                            |
| `bool`             | ✓      |                                            |
| `char`             | ✓      |                                            |
| `string`           | ✓      |                                            |
| `var`              | ✓      |                                            |
| nested struct      | ✓      | See example 16                             |
| vector (`int[]`)   | ✓      | `= []` in field def is handled correctly   |
| array (`int[N]`)   | ✓      | `= {}` in field def emits correct C-style  |
| map                | ✓      |                                            |

### Vector element types (`T[]`)

| Element type       | Status | Notes                                               |
|--------------------|--------|-----------------------------------------------------|
| `int`              | ✓      |                                                     |
| `float`            | ✓      |                                                     |
| `bool`             | ✓      |                                                     |
| `char`             | ✓      |                                                     |
| `string`           | ✓      |                                                     |
| `var`              | ✓      |                                                     |
| user-defined struct| ✓      | `Point[]` works — parser fix in Part 1              |
| `map`              | ✓      | `map[]` works — emits `std::vector<std::map<std::string, var>>` |
| vector (`int[][]`) | ✗      | Not yet implemented — see KNOWN_LIMITATIONS #9      |

### Map value types

| Value type         | Status | Notes                                             |
|--------------------|--------|---------------------------------------------------|
| `int`              | ✓      |                                                   |
| `float`            | ✓      |                                                   |
| `bool`             | ✓      |                                                   |
| `char`             | ✓      |                                                   |
| `string`           | ✓      | See example 08                                    |
| `var`              | ✓      |                                                   |
| vector             | ✓      | `mv["nums"] = nums` works                         |
| nested map         | ✓      | `outer["nested"] = inner` works                   |
| user-defined struct| ✗      | `var` lacks `operator=` for struct — see KNOWN_LIMITATIONS #7 |
| fixed-size array   | ✗      | Arrays are not first-class objects — see KNOWN_LIMITATIONS #8 |
