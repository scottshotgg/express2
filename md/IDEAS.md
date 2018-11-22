# Ideas about the language

- static and dynamic typing
- type inference using `let`
- immutability for static variables by default
- `const` -> value can't be changed

const c = 7

let something = 9
<br>
int something = 9
<br>
var something = 9
<br>
let something = [ 8.9, 9.9, 3.2, 3.3333 ] -> float[4]
<br>
-> static length, static type
<br>
let something = [ 8.9, 9.9, 3.2, 3.3333, .. ] -> float[]
<br>
-> dynamic length, static type
<br>
let something = [ "woah dude", 9, false, { a: "yeah man" }, 6.66 ]

var[] something = [ "woah dude", 9, false, { a: "yeah man" }, 6.66, .. ]

struct thing = {
  int a
  float b
}

let z = thing{
  a: 16,
  b: 1.1111
}
