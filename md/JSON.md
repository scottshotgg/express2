# Native JSON in Express

Native JSON markdown usage in Express is a push to allow the programmer to natively express JSON values in Express. Like Go, this seems to conflict with the public and private capitalization at first, but consider the structure of a `map` in Express:

```cpp
map m = {
  "some" : "value",
  6 : false,
  hey : "its me",
  78 : 89,
  false : "what what",

  // Arrays can be used as a keys and values
  [ true, false, true, false ] : 0xA,
  'c' : [ 7, 8, 9 ],

  // A block can be used as a keys and values
  { 7 : 7 } : "something",
  55.55 : {
    "another" : "block",
    "of" : "stuff"
  }
}
```
<br>

Here we have no actual variables that would inheret the public or private usage since the values of the map are in fact just literal values, which are not subject to public/private behavior.