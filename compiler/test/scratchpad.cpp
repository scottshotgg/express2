// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include <iostream>
#include <string>

// Namespaces:

// Types:
// none

// Structs:
struct Poodle {
  std::string kind = "123";
};

// Interfaces:

// Prototypes:
std::string Kind(Poodle p);

// Functions:
std::string Kind(Poodle p) {
  defer onReturn, onExit;
  return "poodle";
}

// Main:
// generated: false
int main() {
  defer onReturn, onExit;
  Poodle p = {};
  p.kind = "blah";
  Kind(p);
  p = Poodle{};
}