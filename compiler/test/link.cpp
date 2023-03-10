// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/var.cpp"
#include <iostream>
#include <map>
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
go(coroutine [=](...){{printf("hi\n");__time::Sleep(100);
}
}());
std::map<var, var> m = {
    {"i", 9},
    {9, "blah"},
};
std::cout << "m.i:"
          << " " << m["i"] << std::endl;
std::cout << "m.9:"
          << " " << m[9] << std::endl;
struct Expiry {
int seconds = 9;
int nanos;
};
Expiry e = {
    .nanos = 1,
};
printf("Default seconds value: %d\n", e.seconds);
__time::Sleep(1000);
onReturn.deferStack.push([=](...) { printf("hey its me\n"); });
printf("Atoi: %d\n", convert("97"));
}