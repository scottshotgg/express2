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
2) `map`
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

There is one uncategorized type, which we will explain later:
1) `var`

<br>

A keyword that may seem to compose itself as a type, but actually isn't is `val`. This is actually a keyword that invokes a type inference to specifically create an _immutable_ variable.