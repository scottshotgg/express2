// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include <array>
#include <iostream>
#include <stdio.h>
#include <string>
#include <vector>

// Namespaces:
// none

// Types:
typedef int Score;

// Structs:
struct Student {
  std::string name = "";
  Score grade = 0;
  bool passing = false;
};

// Prototypes:
void printStudent(Student s);
bool isPassing(Score grade);

// Functions:
void printStudent(Student s) {
  defer onReturn, onExit;
  std::cout << "Name:" << " " << s.name << " " << "Grade:" << " " << s.grade
            << " " << "Passing:" << " " << s.passing << std::endl;
}

bool isPassing(Score grade) {
  defer onReturn, onExit;
  if (grade >= 60) {
    return true;
  }
  return false;
}

// Main:
// generated: false
int main() {
  defer onReturn, onExit;
  Student alice = {
      .name = "Alice",
      .grade = 92,
      .passing = false,
  };
  Student bob = {
      .name = "Bob",
      .grade = 55,
      .passing = false,
  };
  Student carol = {
      .name = "Carol",
      .grade = 78,
      .passing = false,
  };
  alice.passing = isPassing(alice.grade);
  bob.passing = isPassing(bob.grade);
  carol.passing = isPassing(carol.grade);
  printStudent(alice);
  printStudent(bob);
  printStudent(carol);
  if (alice.passing) {
    std::cout << alice.name << " " << "passed!" << std::endl;
  }
  if (!bob.passing) {
    std::cout << bob.name << " " << "did not pass." << std::endl;
  }
  std::vector<int> scores = {92, 55, 78};
  int total = 0;
  {
    int i = 0;
    while (i < std::size(scores)) {
      total = total + scores[i];
      (i)++;
    }
  }
  std::cout << "Total score:" << " " << total << std::endl;
  int *gptr = &alice.grade;
  std::cout << "Alice grade via pointer:" << " " << *gptr << std::endl;
}