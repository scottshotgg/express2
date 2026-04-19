# Generics, Error Handling, and Multiple Returns

## The problem

Multiple return values in Express are mainly motivated by error handling. Designing
them well requires first deciding the error handling story, which in turn requires
deciding how much (if any) parameterization the type system should support.

---

## Error handling approaches

### Go-style `(value, error)` multiple returns
- Pros: simple, no new types needed
- Cons: repetitive `if err != nil` everywhere; positional returns are untyped at the call site

### Rust-style `Result<T, E>`
- Pros: explicit, composable, `?` operator for propagation is ergonomic
- Cons: requires generics (`T`, `E` are type parameters); hard to implement without
  a real type system
- Could be approximated with a fixed `Result` struct using `var` for the value —
  but loses type safety

### Union/tagged union (already planned in language-reference.md)
```
union Result = {
  int ok
  string err
}
```
- Pros: no generics needed; readable
- Cons: a different `Result` type per function signature; no reuse without codegen

### C-style: return codes + output params
- `-1` / `nullptr` / `0` on failure (already used in lib/file.expr)
- Simple but not composable

---

## Generics question

Rust `Result<T, E>` needs generics. Options:

### Compile-time generics (C++ templates / Rust monomorphization)
- Fast: each instantiation is specialized code
- Hard to reason about: error messages are notoriously bad, mental model is complex
- User explicitly does NOT want this

### Runtime generics (Java/Go interfaces, boxing)
- Simpler mental model
- Performance cost (heap allocation, indirection)
- `var` in Express is already a runtime "any" type — this is essentially what
  `map m` does today

### Code generation (Go `go generate` style)
- Write a template or generator; run it to produce concrete types
- Explicit: you can read the generated code
- No compiler magic, no hard-to-read error messages
- Fits Express philosophy well
- Could be a separate `expr generate` tool

---

## Tentative direction

1. **Short term**: use `var` as the value type in a fixed `Result` struct for
   error-returning functions. Loses type safety but unblocks error handling patterns.

2. **Medium term**: design a codegen tool (`expr generate` or similar) that
   can instantiate typed Result variants from a template.

3. **Long term**: evaluate whether a limited form of compile-time generics
   (only for standard library types like Result/Option, not user-defined) is
   worth the complexity.

4. **Multiple returns**: defer until error handling is settled. If `Result` covers
   the main use case, full multiple return syntax may not be needed. If it is added,
   limit to 2 return values (value + error) to constrain scope.

---

## Open questions

- Should `?` propagation syntax be adopted (auto-return on error)?
- Is `var` acceptable as the untyped value in a `Result` struct?
- Is `expr generate` the right codegen model, or something closer to C++ X-macros?
- Runtime generics via interface/any — acceptable performance trade-off?
