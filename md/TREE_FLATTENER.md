# Tree Flattening

Take for example, a simple for in-range for loop:
```
for i in [ 1, 2, 3 ] {
  ... statements ...
}
```
<br>
At first, this seems simple to translate, but looking into it more you start to notice that this is really a multi-operation line:
1. induction variable declaration; `int i = 0;`
2. anonymous iterable declaration; `auto ITER = <ITER>;`
3. loop from 0 to the length of the ITER variable

Even in the last statement, you can see this as a compound operation:
1. calculate the length of the ITER; `LEN = len(ITER)`
2. if the induction variable is greater than or equal to length of the ITER then break; `i < LEN`
3. increment the induction variable; `i++`

All of this needs to be wrapped in a block as to not shadow other variables.

In the end, the for in-range loop can be visualized as:

```
{
  int i = 0
  auto random_name = { 1, 2, 3 }

  while (i < len(random_name)) {
    ... statements ...

    i++
  }
}
```
<br>
Considering this breakdown: you could go even further and always infinite while loop with a check inside. I am not sure what this would do to the LLVM optimizer so it is left there.