# Idea: Typed Maps

## Design decision (locked)

`map` is sugar for `map[string, var]`. All untyped maps are `std::map<std::string, var>`.
This matches the universal dynamic-language convention: JS objects, JSON, Python namespaces —
untyped maps are always string-keyed. Non-string keys require an explicit type annotation.

## Motivation

- Type safety: catch key/value type mismatches at compile time
- Performance: `std::map<std::string, int>` is more efficient than `std::map<std::string, var>`
- C++ output quality: typed maps emit idiomatic C++ without the `var` overhead
- Enables typed map fields in structs

## Proposed syntax

```
map[string, int] scores
map[string, Point] positions = { "origin" : { x = 0  y = 0 } }
```

`map` without annotation remains `map[string, var]`.

## Open questions

- Parser: `map[K, V]` conflicts with `map m; m[key]` — need to distinguish
  type-annotation brackets from index brackets (same challenge as `Point[]`)
- Scope tree: typed maps need a new TypeValue representation tracking K and V
- Inference: could the compiler infer `map[string, int]` from a literal where all
  values are int? Probably defer — explicit annotation is clearer
- Interaction with `{ }` literal syntax: `map[string, int] m = { "a" : 1 }` should
  work naturally once typed maps are parsed
- Struct fields: `map[string, int] scores = {}` as a struct field

## Non-goals

- Changing existing `map` (untyped) semantics — already locked in
- HashMap / unordered map — separate idea
