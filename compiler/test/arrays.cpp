// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/var.cpp"
#include <iostream>
#include <map>
#include <string>
#include <vector>

// Namespaces:

// Types:
// none

// Structs:
struct person
{
  std::string name;
};

// Interfaces:

// Prototypes:
// none

// Functions:// none
// Main:
// generated: false
int main()
{
  defer onReturn, onExit;
  std::vector<int> int_array = {6, 7, 8};
  std::vector<int> int_vector = {6, 7, 8};
  std::vector<bool> bool_array = {true, false, true};
  std::vector<bool> bool_vector = {true, false, true};
  std::vector<float> float_array = {6.6, 7.77, 8.888};
  std::vector<float> float_vector = {6.6, 7.77, 8.888};
  map<var, var> m = {
      {6, "six"},
  };
  std::vector<map<var, var>> map_array = {
      {{
          6,
          "six",
      }},
  };
}