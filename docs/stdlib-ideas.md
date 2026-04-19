# Express Standard Library — Ideas

This is a working document for designing the Express stdlib. Iterate here before planning implementation.

---

## Design Direction: Model After Go's stdlib

Key principles to adopt from Go:
- **Packages are directories** — `lib/fmt/fmt.expr`, `lib/os/os.expr`, not flat files
- **Small, focused packages** — each does one thing well
- **Explicit imports** — `import "fmt"`, `import "os"`, etc.
- **Package name = directory name** — `import "fmt"` → package `fmt`

---

## Open Design Questions

### 1. `fmt` vs builtin `Println`
- Go deprecated bare prints in favor of `fmt.Println`, `fmt.Printf`, `fmt.Sprintf`
- Express currently has `Println` as a compiler builtin
- Do we migrate to `fmt.Println`? Keep both? Deprecate the builtin eventually?

### 2. Package directory structure
- Flat: `lib/fmt.expr`, `lib/os.expr` (simpler, what we have now with file.expr)
- Directory-per-package: `lib/fmt/fmt.expr`, `lib/os/os.expr` (Go convention, scales better)
- Decision needed before implementing multiple packages

### 3. Error handling
- Go returns `(value, error)` — Express doesn't have multi-return yet
- Options:
  - Return `-1`/`nil` on failure (C style, what file.expr does now)
  - Design a `Result` type now: `struct Result { bool ok; string err }`
  - Defer until multi-return is implemented
- Decision affects ALL stdlib API design

### 4. String methods
- Express strings are C++ `std::string` underneath
- `.c_str()`, `.length()`, `.substr()` etc. already accessible via method call syntax
- Do we wrap these in a `strings` package or leave as direct method calls?

---

## Proposed Packages (Initial Wave)

### `fmt` — formatted I/O
Wraps printf-family. The workhorse.

```
func Printf(string fmt, ...)    // c.printf(fmt, ...)
func Println(...)               // c.printf + newline
func Sprintf(string fmt, ...) string  // snprintf into a buffer
func Fprintf(File f, string fmt, ...)  // c.fprintf
```

**Depends on:** `import c`  
**C headers needed:** `stdio.h` (already in `import c`)

### `os` — operating system interface
```
func Exit(int code)             // c.exit(code)
func Getenv(string key) string  // c.getenv(key)
func Getcwd() string            // c.getcwd(buf, size)
var Args []string               // constructed from argc/argv (harder — needs main() integration)
```

**Depends on:** `import c`  
**C headers needed:** `stdlib.h`, `unistd.h`

### `math` — math functions
```
func Sqrt(float x) float        // c.sqrt(x)
func Abs(float x) float         // c.fabs(x)
func Pow(float base, float exp) float  // c.pow(base, exp)
func Floor(float x) float       // c.floor(x)
func Ceil(float x) float        // c.ceil(x)
func Log(float x) float         // c.log(x)
const Pi = 3.14159265358979323846
const E  = 2.71828182845904523536
```

**Depends on:** `import c`  
**C headers needed:** `math.h` (needs to be added to `import c` list)

### `strings` — string operations
```
func Contains(string s, string sub) bool  // use c.strstr
func HasPrefix(string s, string pre) bool
func HasSuffix(string s, string suf) bool
func ToUpper(string s) string   // c.toupper loop or std::transform
func ToLower(string s) string
func Trim(string s) string
func Split(string s, string sep) []string
func Join([]string parts, string sep) string
func Length(string s) int       // c.strlen or s.length()
```

**Depends on:** `import c`  
**C headers needed:** `string.h`, `ctype.h`

### `file` — file I/O (already written)
`lib/file.expr` — just needs to compile (see compiler fixes plan).

### `time` — time and sleep
```
func Sleep(int ms)              // c.usleep(ms * 1000)
func Now() int                  // c.time(NULL)
func Format(int t) string       // c.strftime
```

**Depends on:** `import c`  
**C headers needed:** `time.h`, `unistd.h`

---

## Longer-Term / Harder Packages

### `io` — reader/writer interfaces
Requires interface types (not yet designed for Express).

### `net` — networking
Requires socket APIs, error handling design, possibly async model.

### `sync` — concurrency primitives
Express has `launch` (goroutine analog) and `libmill` — needs design.

### `json` — JSON encode/decode
Requires a C JSON library (e.g., cJSON) or a pure Express implementation.

---

## `bindc` Tool (separate project)

**Location:** `lib/bindc/`  
**Purpose:** Parse C headers (`.h` files) and generate `.expr` type signature files.  
**Effect:** Enables IDE type-checking for `c.X` calls. No runtime impact — transpiler still strips `c.` prefix.  
**Example:** `bindc sqlite3.h` → `lib/bindc/sqlite3.expr` with type signatures for all sqlite3 functions.

---

## Notes

- All initial stdlib packages only need `import c` — no Express-to-Express dependencies
- This makes them compilable as soon as the compiler fixes are done
- `os.Args` is the one exception — it needs access to `argc`/`argv` from `main()`, which requires a special compiler mechanism
- `math.h` needs to be added to the `import c` header list (currently only stdio, stdlib, string, unistd, libgen)
