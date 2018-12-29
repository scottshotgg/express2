// Includes:
#include <array>
#include <string>

// Imports:
// none

// Types:
// none

// Prototypes:
void another(int i, std::string s);
void something();

// Functions:
void another(int i, std::string s) { int j = 6666666; }

void something() {
  typedef int myInt;
  int i = 10;
  another(10, "s");
}

// Misc:
typedef int myInt;
struct myStruct {
  int i = 10;
  std::string something = "something";
};

// Main:
// generated: false
int main() {
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