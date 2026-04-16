// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include <iostream>
#include <stdio.h>
#include <string>

// Namespaces:
// none

// Types:
// none

// Structs:

// Prototypes:
void greet(std::string name);
int add(int a, int b);
int factorial(int n);

// Functions:
void greet(std::string name) {
  defer onReturn, onExit;
  std::cout << "Hello," << " " << name << std::endl;
}

int add(int a, int b) {
  defer onReturn, onExit;
  return a + b;
}

int factorial(int n) {
  defer onReturn, onExit;
  if (n <= 1) {
    return 1;
  }
  return n * factorial(n - 1);
}

// Main:
// generated: false
int main() {
  defer onReturn, onExit;
  greet("Alice");
  greet("Bob");
  int sum = add(5, 10);
  std::cout << "5 + 10 =" << " " << sum << std::endl;
  std::cout << "Factorial of 5:" << " " << factorial(5) << std::endl;
}