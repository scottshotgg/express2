// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include <array>
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
  int temperature = 75;
  if (temperature > 90) {
    std::cout << "It is hot outside" << std::endl;
  } else if (temperature > 70) {
    std::cout << "The weather is pleasant" << std::endl;
  } else {
    std::cout << "It is cold outside" << std::endl;
  }
  int x = 10;
  if (x == 10) {
    std::cout << "x is exactly 10" << std::endl;
  }
  {
    int i = 0;
    auto arr_0 = {10, 20, 30};
    while (i < std::size(arr_0)) {
      std::cout << "Index:" << " " << i << std::endl;
      (i)++;
    }
  }
}