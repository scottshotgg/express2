# new TODO

## dont focus on dynamic variables for now

0. Branch the compiler
    - checkout new branch: `rearch_parser`
    - delete everything referring to checking in `builder`, rely on tests heavily
    - build a testing framework

1. Rebuild `builder` with no checking at all unless it is related to syntax.
    - i.e, no type or ident checking
    - create a `parser` package and change `builder` to `syntax`

2. Create a `semantic` stage that does ident, type, struct field, array bounds etc.
    - first make a dummy stage that does nothing; this will represent the golden path
      and is only expected to work with the proper input
    - make this a modular architecture and each one of these a `pass`

3. Make `transpiler` an interface
    - rename `transpiler` to `ir` 
    - convert existing C++ transpiler to the new interface architecture
    - ultimately fix up the conversion to C++
    - make an LLVM IR implementation