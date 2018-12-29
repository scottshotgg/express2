>In Express, the empty struct from C++ is Expressed as an open block.

>In C++ where you would write:
```
(struct{})
```
>In Express you can leave the struct off and just write:
```
{}
```
<br>

>In Express when you write:
```cpp
type thing = struct {
  x
  y
  x
}

struct thing = {
  x
  y
  x
}
```

>Both are translated to C++ as:
```cpp
struct something {
  .x;
  .y;
  .z;
}
```
<br>

>Blanking (empty initialization) the struct:
```
something s;
```

>Will translate to:
```
something s = {};
```
<br>

>Whereas an initialization with a value such as:
```cpp
struct s = {
  x = 1
  y = 7
  z = 45
}
```

> Will become:
```cpp
struct s = {
  .x = 1,
  .y = 7,
  .z = 45,
}
```