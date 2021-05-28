// Namespace:
// none

// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/var.cpp"
#include <array>
#include <map>
#include <string>
#include <vector>

#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/libmill/libmill.h"

// Types:
// none

// Structs:
struct Token {
  std::string name = "";
  var value;
};

// Prototypes:
void delayedPrintln();

// Functions:
void delayedPrintln() {
  defer onReturn, onExit;
  msleep(now(), 1000)("hi");
}

// Main:
// generated: false
int main() {
  defer onReturn, onExit;
  onReturn.deferStack.push([=](...) {
    std::cout << "\n--- ENDING ---" << std::endl("\n--- STARTING ---\n");
  });
  go([=](...) { delayedPrintln(); }());
  Token ident = {
      .name = "ident",
  };
  std::map<string, Token> tokens = {
      {"ident", ident},
  };
  std::vector<var> vv = {666, "something_here", false, 73.986622195};
  {
    int x = 0;
    {
      while (x < std::size(vv)) {
        tokens["ident"].value = vv[x];
        std::cout << "vv[x]:" << vv[x]
                  << std::endl("tokens[\"ident\"].value:",
                               tokens["ident"].value, "\n");
        x++;
      }
    }
  }
  std::vector<int> i = {1, 2, 3, 4, 5, 6, 7, 8, 9};
  std::map<var, var> m = {};
  {
    int x = 0;
    {
      while (x < std::size(i)) {
        m[x] = i[x];
        x++;
      }
    }
  }
  std::cout << "---" << std::endl;
  {
    for (auto const &set : m) {
      auto x = set.first;
      m["result"] = m["result"] + m[x] * m[x];
      std::cout << "square: " << m[x] * m[x] << "\nresult: " << m["result"]
                << "\n"
                << std::endl;
    }
  }
}