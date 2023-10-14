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

// Interfaces:

// Prototypes:
int convert(std::string s);
int Atoi(std::string s);

// Functions:
int convert(std::string s) {
  defer onReturn, onExit;
  return Atoi(s);
}

int Atoi(std::string s) {
  defer onReturn, onExit;
  return atoi(s.c_str());
}

// Main:
// generated: false
int main() {
  defer onReturn, onExit;
  onReturn.deferStack.push([=](...) { printf("hey its me\n"); });
  printf("Atoi: %d\n", convert("97"));
}