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

// Structs:
struct Person
{
  std::string name = "7";
  int age;
};

// Prototypes:
std::string Name(Person p);

// Functions:
std::string Name(Person p)
{
  defer onReturn, onExit;
  return p.name;
}

typedef struct
{
  std::string (*const Name)(void *self);
} Nameable;

typedef struct
{
  void *self;
  const Nameable *tc;

} Human;

#define impl_show(T, Namee, name_f)                                            \
  Nameable Namee(T x)                                                          \
  {                                                                            \
    std::string (*const namer_)(T e) = (name_f);                               \
    (void)Name;                                                                \
    static Nameable const tc = {.Name = (std::string(const)(void *))(name_f)}; \
    return (Human){.tc = &tc, .self = x};                                      \
  }

impl_show(Person *, prep_Person_show, Name);

// Main:
// generated: false
int main()
{
  defer onReturn, onExit;
  Person p = {
      .name = "scott",
  };

  std::cout
      << Name(p) << std::endl;
}

// Nameable prep_person_name(Person *x)
// {
//   char *(*const name_)(Person * e) = (name_f);
//   (void)name_;
//   static Nameable const tc = {.Name = (char *(*const)(void *))(name_f)};
//   return (Human){.tc = &tc, .self = x};
// }
