// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/lib/defer.cpp"
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/lib/var.cpp"
#include <array>
#include <map>
#include <string>

#include <libmill.h>

// Types:
typedef int myInt;

// Structs:
struct AnotherOne {
  bool abc = true;
};
struct myStruct {
  int i = 10;
  std::string something = "something";
  AnotherOne ayy = {};
};

// Prototypes:
void something();
void another(int i, std::string s);

// Functions:
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

void another(int i, std::string s) {
  defer onReturn, onExit;
  int j = 6666666;
}

// Main:
// generated: false
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
  };
  int a = 0;
  int *b = &a;
  int c = *b;
  int i = 800008;
  {
    int j = 0;
    auto RANDOM_NAME_LATER = {1, 2, 3};
    {
      while (j < std::size(RANDOM_NAME_LATER)) {
        i = j;
        (j)++;
      }
    }
  }
  int k = 10;
  {
    int j = 0;
    auto RANDOM_NAME_LATER = {1, 2, 3};
    {
      while (j < std::size(RANDOM_NAME_LATER)) {
        {
          int j = 0;
          auto RANDOM_NAME_LATER = {4, 5, 6};
          {
            while (j < std::size(RANDOM_NAME_LATER)) {
              i = j;
              (j)++;
            }
          }
        }
        (j)++;
      }
    }
  }
  something();
}