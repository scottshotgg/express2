// Namespace:
// none

// Includes:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/compiler/test/test_import.expr"

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/var.cpp"
#include <"/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/std.cpp">
#include <"/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/file.cpp">
#include <map>
#include <string>

// Types:
// none

// Structs:

// Prototypes:
// none

// Functions:// none
// Main:
// generated: false
int main() {
  defer onReturn, onExit;
  int i = 0;
  bool b = 2 + 2 == 3 + 3;
  struct something {
    int a = 5;
    std::string s = "0";
  };
  struct yo {
    something sth = {};
    bool ayy = false;
  };
  something s = {};
  yo e = {
      .sth =
          {
              .a = 9,
          },
  };
  s.a = 88;
  std::map<var, var> m = {
      {"me", "you"},
      {"num", 8},
  };
  var m2 = m["me"];
  int result = m["num"] + e.sth.a;
  Println(m["me"], e.sth.a, e.ayy, result);
  s.a = s.a + result;
  Print(s.a);
  Something();
}