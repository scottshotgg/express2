// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include <array>
#include <iostream>
#include <stdio.h>
#include <vector>

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
  std::vector<int> numbers = {1, 2, 3, 4, 5};
  std::cout << "First element:" << " " << numbers[0] << std::endl;
  std::cout << "Third element:" << " " << numbers[2] << std::endl;
  std::cout << "For-in loop (indices):" << std::endl;
  {
    int i = 0;
    while (i < std::size(numbers)) {
      std::cout << "Index:" << " " << i << std::endl;
      (i)++;
    }
  }
  std::cout << "For-of loop (values):" << std::endl;
  {
    int _idx_0 = 0;
    while (_idx_0 < std::size(numbers)) {
      auto value = numbers[_idx_0];
      std::cout << "Value:" << " " << value << std::endl;
      (_idx_0)++;
    }
  }
}