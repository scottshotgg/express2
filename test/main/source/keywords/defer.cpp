// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include <iostream>
#include <stdio.h>

// Namespaces:
// none

// Types:
// none

// Structs:

// Prototypes:
void deferred();

// Functions:
void deferred() {
  defer onReturn, onExit;
  std::cout << "deferred call" << std::endl;
}

// Main:
// generated: false
int main() {
  defer onReturn, onExit;
  onReturn.deferStack.push([=](...) { deferred(); });
  onReturn.deferStack.push(
      [=](...) { std::cout << "inline deferred" << std::endl; });
  std::cout << "main body" << std::endl;
}