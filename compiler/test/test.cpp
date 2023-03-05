// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/var.cpp"
#include <iostream>
#include <map>
#include <stdio.h>
#include <string>

// Namespaces:
// none

// Types:
typedef int myInt;

// Structs:
struct AnotherOne {
  bool abc = true;
};
struct myStruct {
  int i = 10;
  std::string something = "something";
  AnotherOne ayy;
};

// Prototypes:
void another(int i, std::string s);
void something();

// Functions:
void another(int i, std::string s) {
  defer onReturn, onExit;
  int j = 6666666;
  printf("something");
}

void something() {
  defer onReturn, onExit;
  myStruct s = {
      .i = 100 * 7 / 3,
      .something = "else",
      .ayy =
          {
              .abc = true,
          },
  };
  int i = 10;
  another(i, "s");
}

// Main:
// generated: false
int main() {
  defer onReturn, onExit;
  enum {
    some,
    one,
    here,
  };
  if (69 > one + 20) {
    int x = 7;
  } else if (some) {
    var y = "1000000" + true;
  } else {
  }
  onReturn.deferStack.push([=](...) { something(); });
  std::string thing = "thing";
  std::string nothing = "nothing";
  std::map<var, var> m = {
      {thing, "thing"},
      {"not_a_thing", nothing},
      {6, true},
      {false, "thing"},
  };
  var m_thing = m[thing];
  int a = 0;
  int *b = &a;
  int c = *b;
  int i = 800008;
  int k = 10;
  another(k, thing);
}