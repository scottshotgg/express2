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
struct Blah {};
struct Bulldog {};
struct Husky {};

// Interfaces:
typedef struct {
  void *self;
  std::string (*const Breed)(void *self);
} Dog;

// Prototypes:
std::string Breed(Bulldog b);
std::string Breed(Husky h);

// Functions:
std::string Breed(Bulldog b) {
  defer onReturn, onExit;
  return "bulldog";
}

std::string Breed(Husky h) {
  defer onReturn, onExit;
  return "husky";
}

// Main:
// generated: false
int main() {
  defer onReturn, onExit;
  Husky h = {};
  Dog d = {
      .self = &h,
      .Breed = helper(Breed(h)),
  };
  std::cout << d.Breed() << std::endl;
}