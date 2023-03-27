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
struct Person {
  std::string name = "7";
  int age;
};

// Prototypes:
std::string Name();
std::string Name(Person p);

// Functions:
std::string Name() {
  defer onReturn, onExit;
  return "blah";
}

std::string Name(Person p) {
  defer onReturn, onExit;
  return p.name;
}

// Main:
// generated: false
int main() {
  defer onReturn, onExit;
  Person p = {
      .name = "scott",
  };
  std::cout << Name(p) << std::endl;
}