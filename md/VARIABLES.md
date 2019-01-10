# On Types in Express

Primitive types consist of _single-value_ types:
1) `byte`
2) `int`
3) `bool`
4) `char` - Default UTF-8 encoded
5) `float`
6) `stmt`
7) `expr`

<br>

The next type to learn in the Express type hierarchy are the _Composite types_. These consist of types that hold multiple values of different types:
1) `struct`
2) `map`
3) `object`
4) `tuple`
5) `function`
6) `macro`

<br>

In addition to these are what is known as _Repeated types_. These types allow for multiple values of a single type:
1) `bytes`
2) `string` - Default UTF-8 encoded
3) `array` - denoted with the `[<int>]` notation, can be applied to any type
4) `list` - denoted with the `[]` notation, can be applied to any type