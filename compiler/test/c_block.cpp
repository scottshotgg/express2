// Includes:

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include <iostream>
#include <stdio.h>
#include <string>

// Namespaces:
namespace __strings {
// Includes:
// none

// Imports:
// none

// Namespaces:
namespace strings {}

// Types:
// none

// Structs:

// Prototypes:
int Atoi(std::string s);
std::string Itoa(int i);

// Functions:
int Atoi(std::string s) {
  defer onReturn, onExit;
  return atoi(s.c_str());
}

std::string Itoa(int i) {
  defer onReturn, onExit;
  std::string s = "";
  s.append(1, char(i));
  return s;
}

// Main:
// generated: false

} // namespace __strings

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
  int i = __strings::Atoi("97");
  std::string s = __strings::Itoa(i);
  printf("i: %d\n", i);
  printf("s: %s\n", s.c_str());
}