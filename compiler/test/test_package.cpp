// Namespace:
// none

// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/var.cpp"
#include <map>
#include <string>

#include <libmill.h>

// Types:
// none

// Structs:

class myPackage {
  public:
typedef int myInt;
struct AnotherOne {
  bool abc = true;
};
struct myStruct {
  int i = 10;
  std::string something = "something";
  AnotherOne ayy = {};
};
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
int main() {
  defer onReturn, onExit;
  object o = {};
  enum {
    some,
    one = some + 2,
    here,
  };
  if (69 > one + 20) {
    int x = 7;
  } else if (some) {
    var y = "1000000" + true;
  } else {
    go([=](...) { something(); }());
  }
  go([=](...) { something(); }());
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
  {
    int j = 0;
    auto SOMETHING = {1, 2, 3};
    while (j < std::size(SOMETHING)) {
      i = j;
    }
    (j)++;
  }
  int k = 10;
  {
    int j = 0;
    auto SOMETHING = {1, 2, 3};
    while (j < std::size(SOMETHING)) {
      {
        int j = 0;
        auto SOMETHING = {4, 5, 6};
        while (j < std::size(SOMETHING)) {
          i = j;
        }
        (j)++;
      }
    }
    (j)++;
  }
  something();
}
int another(int i, std::string s) {
  defer onReturn, onExit;
  return 6666666;
}
}; // namespace myPackage

// Prototypes:
// none

// Functions:// none
// Main:
// generated: false

namespace main_lskjflksdfj {
  int main() {
    std::cout << "hey its me" << std::endl;

    return 0;
  }
}

int main() {
  main_lskjflksdfj::main();

  myPackage m;
  std::cout << m.another(0, "sd") << std::endl;
 }

// package main will essentially just be a namespace that is called from a main function
// all other packages will be translated into classes with their appropriate public/private
// stuff