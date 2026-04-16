# Known Limitations

This document lists the current limitations of the Express compiler. Each item
produces a C++ compilation error rather than a clear Express-level diagnostic
because the semantic pass is not yet complete.

---

## 1. No type checking

The semantic pass is stubbed. Type errors are caught by clang, not by Express.
Error messages will reference generated C++ code rather than the original `.expr` file.

## 2. No `*ptr * 5` dereference-then-multiply

The `*` operator after a pointer always starts a new dereference statement.
Use parentheses to work around this: `(*ptr) * 5`.

## 3. No `nil`

Nil comparisons and nil literals are not supported.

## 4. No compound assignment operators

`+=`, `-=`, `*=`, `/=` are not implemented. Use explicit assignment:
```
x = x + 1
```

## 5. No `len()`

The built-in `len()` function is not supported. Array sizes must be tracked
manually or via a known constant.

## 6. No nested struct fields

Struct fields must be primitive types (`int`, `string`, `bool`, `float`, `char`)
or type aliases of primitives. Structs containing other structs are not yet supported.

## 7. No struct methods

Method syntax (functions defined on types) is not implemented. Use free functions
that take a struct parameter instead.

## 8. No typed map fields

`map K -> V` as a struct field does not work. Maps at the top level are supported.

## 9. No multiple return values end-to-end

The syntax parses but transpilation of multiple return values is incomplete.

## 10. Chained pointer access

Only one level of `->` is supported: `s->field` works, but `s->field->subfield`
does not. Full chaining requires a type-checking pass to know the type of `field`.
