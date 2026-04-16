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
  int a = 10;
  int b = 3;
  int sum = a + b;
  int difference = a - b;
  int product = a * b;
  int quotient = a / b;
  int remainder = a % b;
  std::cout << "a =" << " " << a << std::endl;
  std::cout << "b =" << " " << b << std::endl;
  std::cout << "sum =" << " " << sum << std::endl;
  std::cout << "difference =" << " " << difference << std::endl;
  std::cout << "product =" << " " << product << std::endl;
  std::cout << "quotient =" << " " << quotient << std::endl;
  std::cout << "remainder =" << " " << remainder << std::endl;
  int counter = 0;
  (counter)++;
  (counter)++;
  (counter)--;
  std::cout << "counter =" << " " << counter << std::endl;
}