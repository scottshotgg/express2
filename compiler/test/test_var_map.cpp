// Includes:
#include <array>
#include <map>
#include <string>
#include "../../lib/var.cpp"
#include <iostream>

// Imports:
// none

// Types:
typedef int myInt;
// none

// Structs:
struct AnotherOne {
  bool abc = true;
};
struct myStruct {
  int i = 10;
  std::string something = "something";
  AnotherOne ayy = {};
};
// none

// Prototypes:
void something();
void another(int i, std::string s);

// Functions:
void something() {
  myStruct s = {
      .i = 100 * 7 / 3,
      .something = "else",
      .ayy =
          {
              .abc = true,
          },
  };
  int i = 10;
  another(i, "s");
}

void another(int i, std::string s) { int j = 6666666; }

// Main:
// generated: false
int main() {
  std::string thing = "thing";
  std::string nothing = "nothing";
  std::map<var, var> m = {
      {thing, {"here", {"i", "am"}}},
      {{"not_a_thing","nothi"}, 33378301},
      {{{false,"hi"},"something"}, 78301},
      {{{false,"hi"},"some"}, 72228301},
      {6, nothing},
      {true, 982.321},
  };

  std::cout << "Testing object keying:" << std::endl;
  std::cout << "true :" << m[true] << std::endl;
  std::cout << "thing :" << m["thing"] << std::endl;
  std::cout << "6: " << m[6] << std::endl;
  std::cout << "{{false,\"hi\"},\"some\"} :" << m[{{false,"hi"},"some"}] << std::endl;
  std::cout << std::endl;

  var h = 77;
  var g;

  std::cout << "Testing null:" << std::endl;
  std::cout << g << std::endl;
  std::cout << h << std::endl;
  std::cout << h["thing"] << std::endl;

  int i = 0;
  {
    int j = 0;
    auto RANDOM_NAME_LATER = {1, 2, 3};
    {
      while (j < std::size(RANDOM_NAME_LATER)) {
        i = j;
        (j)++;
      }
    }
  }
  int k = 10;
  {
    int j = 0;
    auto RANDOM_NAME_LATER = {1, 2, 3};
    {
      while (j < std::size(RANDOM_NAME_LATER)) {
        {
          int j = 0;
          auto RANDOM_NAME_LATER = {4, 5, 6};
          {
            while (j < std::size(RANDOM_NAME_LATER)) {
              i = j;
              (j)++;
            }
          }
        }
        (j)++;
      }
    }
  }
  something();
}