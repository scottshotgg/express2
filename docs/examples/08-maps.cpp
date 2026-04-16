// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include <iostream>
#include <map>
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
  std::map<std::string, std::string> colors = {
      {"red", "FF0000"},
      {"green", "00FF00"},
      {"blue", "0000FF"},
  };
  std::cout << "Red hex:" << " " << colors["red"] << std::endl;
  std::cout << "Blue hex:" << " " << colors["blue"] << std::endl;
  colors["white"] = "FFFFFF";
  std::cout << "White hex:" << " " << colors["white"] << std::endl;
}