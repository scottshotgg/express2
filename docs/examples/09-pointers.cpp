// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include <iostream>
#include <stdio.h>

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
  int value = 42;
  std::cout << "Original value:" << " " << value << std::endl;
  int *ptr = &value;
  std::cout << "Value via pointer:" << " " << *ptr << std::endl;
  *ptr = 100;
  std::cout << "After modifying via pointer:" << " " << value << std::endl;
  int x = 5;
  int y = 10;
  int *px = &x;
  int *py = &y;
  int temp = *px;
  *px = *py;
  *py = temp;
  std::cout << "After swap: x =" << " " << x << " " << "y =" << " " << y
            << std::endl;
}