// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include <iostream>
#include <stdio.h>
#include <string>

// Namespaces:

// Types:
// none
class Student
{
  virtual std::string Name() = 0;
};

// Structs:
class Person : Student
{
public:
  std::string name = "7";
  int age;
  std::string Name();
};

// Prototypes:
std::string Name(Person p);
std::string Name();

// Functions:
std::string Person::Name()
{
  defer onReturn, onExit;
  return this->name;
}

std::string Name()
{
  defer onReturn, onExit;
  return "blah";
}

// Main:
// generated: false
int main()
{
  defer onReturn, onExit;
  // Person p = {
  //     .name = "scott",
  // };

  Person p;

  Student &s = p;
  std::cout << Name(p) << std::endl;
}