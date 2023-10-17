// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include <iostream>
#include <stdio.h>
#include <string>

// Namespaces:

// Types:
// none

// Structs:
struct Husky {
  std::string kind = "husky";
};

// Interfaces:

// Prototypes:
std::string Kind(Husky *h);

// Functions:
std::string Kind(Husky *h) {
  defer onReturn, onExit;
  return h->kind;
}

// Main:
// generated: false
int main() {
  defer onReturn, onExit;
  Husky h = {
      .kind = "husky 2",
  };
  int i = 7;
  int ii = &i;
  std::cout << Kind(&h) << std::endl;
}