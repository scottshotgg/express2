// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include <iostream>
#include <stdio.h>
#include <string>

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
int main() {
  defer onReturn, onExit;
  std::cout << "Hello, World!" << std::endl;
  std::cout << 42 << std::endl;
  int x = 10;
  std::string name = "Express";
  std::cout << "The value of x is:" << " " << x << " "
            << "and the language is:" << " " << name << std::endl;
}