// Includes:

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/var.cpp"
#include <array>
#include <iostream>
#include <map>
#include <stdio.h>
#include <string>
#include <unistd.h>
#include <vector>

// Namespaces:
namespace __time {
// Includes:
// none

// Imports:
// none

// Namespaces:
// none

// Types:
// none

// Structs:

// Prototypes:
int Now();

// Functions:
int Now() {
  defer onReturn, onExit;
  return time(NULL);
}

// Main:
// generated: false

} // namespace __time

namespace __stack {
// Includes:
// none

// Imports:
// none

// Namespaces:
// none

// Types:
// none

// Structs:

// Prototypes:
// none

// Functions:// none
// Main:
// generated: false

} // namespace __stack

// Types:
// none

// Structs:
struct Token {
  std::string name = "";
  var value;
};
struct Stack {
  std::vector<var> value;
};
Stack s;
std::string num = "current_num";
std::string res = "result";

// Prototypes:
void Push(var v);
var Pop();
std::string to_string(var v);
void printResults(std::map<var, var> m, int x);

// Functions:
void Push(var v) {
  defer onReturn, onExit;
  s.value.push_back(v);
}

var Pop() {
  defer onReturn, onExit;
  var v = s.value[s.value.size() - 1];
  s.value.pop_back();
  return v;
}

std::string to_string(var v) {
  defer onReturn, onExit;
  return v.to_string();
}

void printResults(std::map<var, var> m, int x) {
  defer onReturn, onExit;
  std::cout << "square:"
            << " " << x * x << " "
            << "\nresult:"
            << " " << m[res] << " "
            << "\n"
            << std::endl;
}

// Main:
// generated: false
int main() {
  defer onReturn, onExit;
  std::cout << "size before:"
            << " " << s.value.size() << std::endl;
  Push(7);
  Push("hello");
  Push(false);
  Push(777.333);
  Push("7lskdjfsdlkjfldkj");
  std::cout << "size after:"
            << " " << s.value.size() << std::endl;
  std::cout << Pop() << std::endl;
  sleep(1);
  int t = __time::Now();
  std::cout << "Current time is:"
            << " " << t << std::endl;
  sleep(1);
  onReturn.deferStack.push(
      [=](...) { std::cout << "\n--- ENDING ---" << std::endl; });
  std::cout << "\n--- STARTING ---\n" << std::endl;
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
        std::cout << "vv[x]:"
                  << " " << vv[x] << std::endl;
        std::cout << "tokens[\"ident\"].value:"
                  << " " << tokens["ident"].value << " "
                  << "\n"
                  << std::endl;
        x++;
      }
    }
  }
  std::cout << "---\n" << std::endl;
  std::vector<int> i = {1, 2, 3, 4, 5, 6, 7, 8, 9};
  std::map<var, var> m;
  {
    int x = 0;
    {
      while (x < std::size(i)) {
        m[x] = x * x;
        m[res] = m[res] + m[x];
        m[x] = to_string(x * x);
        printResults(m, x);
        x++;
      }
    }
  }
  m[res] = m[res].to_string();
  m["done"] = true;
  std::cout << "m:"
            << " " << m << std::endl;
  sleep(1);
  std::cout << "Program took:"
            << " " << __time::Now() - t << " "
            << "seconds" << std::endl;
}