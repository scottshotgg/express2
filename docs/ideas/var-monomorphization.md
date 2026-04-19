# Idea: Transparent `var` Monomorphization

## Summary

When a `var` binding is assigned values of only one concrete type across all write
points in its scope, the compiler can transparently emit a statically-typed C++
variable instead of a dynamic one — giving native performance without the programmer
doing anything.

## The Observation

Given:

```
var x = 5
x++
x = 99
int y = 7
x += y
```

At every write point (`= 5`, `++`, `= 99`, `+= y`) the type is `int`. The compiler
can prove this statically. Instead of emitting `std::any x = 5`, it can emit
`int x = 5` and get full native int performance.

The programmer wrote a dynamically-typed declaration and got a statically-typed
binary — transparently.

## The Rule

Collect all write points for a `var` binding within its scope. Infer the type at
each write point. If all write points unify to a single concrete type, emit that
type directly in C++. If they don't unify, fall back to the dynamic `var`
representation.

## What Breaks Monomorphization

Any heterogeneous assignment forces the binding to remain dynamic:

```
var x = 5
x = "hello"   // int and string do not unify → must stay var
```

Control flow joining two different types also breaks it:

```
var x = 5
if condition {
    x = "hello"
}
// x : int | string → must stay var
```

Even if only one branch assigns a different type, the join after the `if` is
ambiguous and must be dynamic.

## Flow-Sensitive Analysis

For control flow, the type of `x` after a branch is the join of its type on each
incoming path. If all paths agree, monomorphization holds:

```
var x = 5
if condition {
    x = 10      // int on this path
} else {
    x = 99      // int on this path
}
// join: int — still monomorphizable
```

This is the same analysis TypeScript does for narrowing and Julia does for method
specialization, but applied transparently at transpile time.

## Relationship to the Mutability Design

This pairs directly with the `var int x` / `var x` / `int x` mutability design:

- `int x = 5` — immutable, always a C++ `const int`
- `var int x = 5` — mutable, type fixed, always a C++ `int`
- `var x = 5` — mutable, type varies; monomorphize if possible, else dynamic

The programmer sees a spectrum from dynamic to static. The compiler closes the gap
between "I wrote `var x`" and "I get dynamic overhead" whenever it can prove the
type is uniform.

## Implementation Sketch

This would live in the semantic pass, after the AST is built and before transpilation:

1. Walk all `var` declarations in a function scope.
2. For each `var` binding, collect every assignment node (including `++`, `+=`, etc.).
3. Infer the concrete type of each assignment's RHS.
4. If all inferred types are identical and concrete (not `var`), annotate the
   declaration node with that type.
5. The transpiler reads the annotation and emits the concrete C++ type instead of
   the dynamic representation.

## Open Questions

- What is the dynamic `var` representation in C++? (`std::any`, tagged union, or
  something from `lib/var.cpp`?) The answer determines what "falling back" costs.
- Does monomorphization interact with function call boundaries? If `var x` is passed
  to a function that accepts `var`, does the callee need to know the concrete type?
- Should monomorphization be visible to the programmer (e.g., a warning: "this var
  was monomorphized to int") or fully silent?
- How does this interact with `for-of` loops where the loop variable is implicitly
  typed from the container element type?
