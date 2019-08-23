// Namespace:
// none

// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/var.cpp"
#include <map>
#include <string>

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
  enum {
    Male,
    Female,
    Helicopter,
  };
  struct Person {
    std::string Name = "";
    int Age = 0;
    int Gender = Male;
    std::map<var, var> Characteristics = {
        {444, 222},
    };
  };
  FILE *f = fopen("something", "w+");
  std::map<var, var> chars = {
      {"IsProAF", true},
      {69, "truth"},
  };
  Person test = {
      .Name = "scott",
      .Age = 24,
      .Gender = Helicopter,
      .Characteristics = chars,
  };
  std::string output = test.Name + " is a bawss";
  fputs(output.c_str(), f);
  fclose(f);
}