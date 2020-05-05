# On Types in Express

Basic types that consist of single values are known as _Primitive_ types:
1) `byte`
2) `int`
3) `bool`
4) `char` - Default UTF-8 encoded
5) `float`
6) `stmt`
7) `expr`
8) `ref` - invoked with the `&` operator
9) `ptr` - denoted by the `*` operator

<br>

The next type to learn in the Express type hierarchy are the _Composite_ types. These consist of types that hold multiple values of different types:
1) `struct`
2) `map<type,type>`
3) `object`
4) `tuple`
5) `function`
6) `macro`

<br>

In addition to these are what is known as _Repeated_ types. These types allow for multiple values of a single type:
1) `bytes`
2) `string` - Default UTF-8 encoded
3) `array` - denoted with the `[<int>]` notation, can be applied to any type
4) `list` - denoted with the `[]` notation, can be applied to any type

<br>

There is one non-type:
1) `var` - dynamic type used to provide dynamicism for some computational dynamic situations.

There is also two pseudo-type that is ultimately resolved by the compiler.
1) `let` - invokes a type inference in which the compiler is directed to solve for the "ground state" type, which is the most-primitive type that it can be. The only guarantee from the compiler is that the type will be non-dynamic; `var` and `object` are logically resolvable as the ground state for any type but the compiler will not make the decision for you as the run time penalty must be a concious decision. 
For composite types, a let statement will crimp the _inner types_; property types, field types, etc.

2) `const` - a construct that may seem to compose itself as a type, but actually isn't is `val`. This is actually a keyword that invokes a type inference to specifically declare an _immutable_ variable. From here an assumption can be founded that you cannot modify any value _within_ a value and as such it will maintain static objects, structs, tuples, etc

<br>
