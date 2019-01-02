# New Ideas

1) Statements (and, in general, abstract langage features) as first class values
<br>
`stmt s = for i of ["j", "k", "l"] { <i> = rand.Int() }`

2) Interpret expression operator
<br>
`<i>` which means interpret the expression provided into an ident

3) Use variable descent typing

4) Scope operator
<br>
`$<ident>` means that this variable is declared in the scope above instead of the current scope; can also be used with a block