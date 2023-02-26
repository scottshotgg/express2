// Includes:

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include <iostream>

// Namespaces:
namespace __something {
// Includes:
// none

// Imports:
// none

// Namespaces:
namespace something {}

// Types:
// none

// Structs:

// Prototypes:
int what();

// Functions:
int what() {
  defer onReturn, onExit;
  return 7;
}

// Main:
// generated: false

} // namespace __something

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
  __something::what();
}