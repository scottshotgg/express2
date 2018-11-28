# On length in Express:

Strings:
  - Use the `.length` property which will compose into `strlen(<string>)`'

Arrays:
  - Use the `.length` property

Lists:
  - Use the `.length` property which will not return the `capacity` property calculated by `<vector>.size()`

Structs:
  - Use the `.length` property which will return the number of attributes the struct has

Objects:
  - Use the `.length` property which will compose into a `<map>.size()` for the integer keys and another for the string keys, and one for the boolean keys.

Channels:
  - Use the `.length` property which will compose into `<queue>.size()`

Vars:
  - Use the `.length` property which will return either 1 if the `var` is not a composite type or the `.length` property of the encapsulated composite type