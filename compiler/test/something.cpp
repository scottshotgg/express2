// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include <iostream>
#include <stdio.h>

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
  enum {
    monday,
    tuesday,
    wednesday,
    thursday,
  };
  C C = {
      .int i = 6,
      .printf("%d", i),
  };
  char c = '6';
  return 7;
}

// Main:
// generated: false
int main() {
  defer onReturn, onExit;
  printf("what: %d", what());
}