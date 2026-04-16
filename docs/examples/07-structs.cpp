// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include <iostream>
#include <stdio.h>
#include <string>

// Namespaces:
// none

// Types:
// none

// Structs:
struct Person {
  std::string name = "";
  int age = 0;
};

// Prototypes:
void greetPerson(Person p);

// Functions:
void greetPerson(Person p) {
  defer onReturn, onExit;
  std::cout << "Hello," << " " << p.name << std::endl;
  std::cout << "Age:" << " " << p.age << std::endl;
}

// Main:
// generated: false
int main() {
  defer onReturn, onExit;
  Person alice = {
      .name = "Alice",
      .age = 30,
  };
  Person bob = {
      .name = "Bob",
      .age = 25,
  };
  std::cout << "Alice's name:" << " " << alice.name << std::endl;
  std::cout << "Bob's age:" << " " << bob.age << std::endl;
  alice.age = 31;
  std::cout << "Alice's new age:" << " " << alice.age << std::endl;
  greetPerson(alice);
}