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

// Interfaces:

// Prototypes:
std::string Name(Person p);
std::string Name();

// Functions:
std::string Name(Person p) {
  defer onReturn, onExit;
  return p.name;
}

std::string Name() {
  defer onReturn, onExit;
  return "blah";
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