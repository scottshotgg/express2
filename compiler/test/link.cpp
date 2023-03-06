// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include <iostream>
#include <stdio.h>
#include <string>

#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/libmill/libmill.h"
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
int Atoi(std::string s);
int convert(std::string s);

// Functions:
int Atoi(std::string s) {
  defer onReturn, onExit;
  return atoi(s.c_str());
}

int convert(std::string s) {
  defer onReturn, onExit;
  return Atoi(s);
}

// Main:
// generated: false
int main() {
  defer onReturn, onExit;
go(coroutine [=](...){{{
while(true){printf("hi\n");__time::Sleep(100);
}
}
}
}());
__time::Sleep(1000);
onReturn.deferStack.push([=](...) { printf("hey its me\n"); });
printf("Atoi: %d\n", convert("97"));
}