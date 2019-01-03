// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include <array>
#include <map>
#include <string>
#include <libmill.h>

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
void another(int i, std::string s);

// Functions:
coroutine void something() {
  defer onReturn, onExit;
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
  int waitTime = rand() % 1000;
  printf("waiting %d\n", waitTime);
  msleep(now() + waitTime);
  printf("something here %d\n\n", waitTime);
}

void another(int i, std::string s) {
  defer onReturn, onExit;
  int j = 6666666;
}

// Main:
// generated: false
int main() {
  defer onReturn, onExit;

  printf("hi\n");
  go(something());

  onReturn.deferStack.push([=](...){ something(); });

  enum {
    some,
    one = some + 2,
    here,
  };
  std::string thing = "thing";
  std::string nothing = "nothing";
  std::map<std::string, std::string> m = {
      {thing, "thing"},
      {"not_a_thing", nothing},
  };
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

  printf("i am here\n");

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

  printf("hi\n");

  msleep(now() + 1000);
}