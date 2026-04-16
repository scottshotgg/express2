// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include </home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/var.cpp>
#include <iostream>
#include <stdio.h>
#include <string>

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
  var dynamicVar = 42;
  std::cout << "Dynamic var (int):" << " " << dynamicVar << std::endl;
  dynamicVar = "Now I'm a string";
  std::cout << "Dynamic var (string):" << " " << dynamicVar << std::endl;
  dynamicVar = true;
  std::cout << "Dynamic var (bool):" << " " << dynamicVar << std::endl;
  int inferredInt = 100;
  std::string inferredString = "Hello";
  bool inferredBool = false;
  float inferredFloat = 3.14;
  std::cout << "Inferred int:" << " " << inferredInt << std::endl;
  std::cout << "Inferred string:" << " " << inferredString << std::endl;
  std::cout << "Inferred bool:" << " " << inferredBool << std::endl;
  std::cout << "Inferred float:" << " " << inferredFloat << std::endl;
  int explicitInt = 200;
  std::string explicitString = "Explicitly typed";
  bool explicitBool = true;
  float explicitFloat = 2.718;
  std::cout << "Explicit int:" << " " << explicitInt << std::endl;
  std::cout << "Explicit string:" << " " << explicitString << std::endl;
  std::cout << "Explicit bool:" << " " << explicitBool << std::endl;
  std::cout << "Explicit float:" << " " << explicitFloat << std::endl;
}