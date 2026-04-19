# Known Limitations

This document lists the current limitations of the Express compiler. Each item
either produces a C++ compilation error or a parse error rather than a clear
Express-level diagnostic, because the semantic pass is not yet complete.

---

## 1. No type checking

The semantic pass is stubbed. Type errors are caught by clang, not by Express.
Error messages will reference generated C++ code rather than the original `.expr` file.

## 2. No `*ptr * 5` dereference-then-multiply

The `*` operator after a pointer always starts a new dereference statement.
Use parentheses to work around this: `(*ptr) * 5`.

## 3. No `nil`

Nil comparisons and nil literals are not supported. Use `nullptr` directly in a
C block if needed.

## 4. No multiple return values end-to-end

The syntax parses but transpilation of multiple return values is incomplete.

## 5. Chained pointer access

Only one level of `->` is supported: `s->field` works, but `s->field->subfield`
does not. Full chaining requires a type-checking pass to know the type of `field`.

## 6. `len()` not supported for fixed-size arrays

`len(arr)` works for vectors (`int[]`) and strings. It does NOT work for
C-style arrays (`int[5]`) — their size is a compile-time constant and they do
not expose a `.size()` method.

## 7. Struct values cannot be stored in a map

`var` (the map value type) does not have `operator=` for user-defined struct
types. Assigning a struct into a `map<string,var>` fails at C++ compile time.

## 8. Fixed-size arrays cannot be stored in a map

C-style arrays are not first-class objects in C++. They cannot be assigned into
a `var` and therefore cannot be stored as map values.

## 10. Map literal values are limited to strings

`map m = { "k" : v }` infers `std::map<std::string, std::string>` — the value must be a string literal.
For mixed-type values (int, bool, vector, etc.) use the assignment form: `map m; m["k"] = 42`.

## 11. Non-string map keys require explicit typed maps

`map` always uses `std::string` keys. Non-string keys (e.g. `int`,
`bool`) are not supported with untyped maps. Typed maps (`map[K, V]`) are a
planned future feature.

## 13. `val` is a reserved type keyword and cannot be used as an identifier

`val` appears in the type map as a planned future keyword for immutable variable
declarations. It lexes as `token.Type`, not an identifier, so it cannot be used as
a variable or parameter name. Use a different name (e.g. `str`, `v`, `value`).

## 12. `{ }` map literal is only valid in expression context

An inline map literal like `{ "k" : v }` can only appear as the right-hand
side of an assignment or declaration, not as a standalone statement.
