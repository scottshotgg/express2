// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include <iostream>
#include <stdio.h>

// Namespaces:
namespace __time {
// Includes:
// none

// Imports:
// none

// Namespaces:
namespace time {}

// Types:
// none

// Structs:

// Prototypes:
void Sleep(int i);
int Now();

// Functions:
void Sleep(int i) {
  defer onReturn, onExit;
  msleep(Now() + i);
}

int Now() {
  defer onReturn, onExit;
  return now();
}

// Main:
// generated: false

} // namespace __time

// Types:
// none

// Structs:

// Prototypes:
// none

// Functions:// none
// Main:
// generated: false
int main() { defer onReturn, onExit; }