# On length in Express:

String:
  - Use the `.length` property which will compose into `strlen(<string>)`

Array:
  - Use the `.length` property

List:
  - Use the `.length` property which will not return the `capacity` property calculated by `<vector>.size()`

Struct:
  - Use the `.length` property which will return the number of attributes the struct has

Object:
  - Use the `.length` property which will compose into a `<map>.size()` for the integer keys and another for the string keys, and one for the boolean keys.

Channel:
  - Use the `.length` property which will compose into `<queue>.size()`

Var:
  - Use the `.length` property which will return either 1 if the `var` is not a composite type or the `.length` property of the encapsulated composite type