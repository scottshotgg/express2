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
struct Poodle {};
struct Bulldog {
  std::string kind = "bulldog";
};
struct Husky {
  std::string kind = "husky";
};

// Interfaces:
typedef struct {
  void *self;
  std::string (*Kind)(void *self);
} Dog;

// Prototypes:
std::string impl_Dog_Husky_Kind(void *self);
std::string impl_Dog_Bulldog_Kind(void *self);
std::string Kind(Poodle p);
std::string Kind(Bulldog b);
std::string Kind(Husky h);
std::string impl_Dog_Poodle_Kind(void *self);

// Functions:
std::string impl_Dog_Husky_Kind(void *self) { return Kind(*(Husky *)self); }

std::string impl_Dog_Bulldog_Kind(void *self) { return Kind(*(Bulldog *)self); }

std::string Kind(Poodle p) {
  defer onReturn, onExit;
  return "poodle";
}

std::string Kind(Bulldog b) {
  defer onReturn, onExit;
  return b.kind;
}

std::string Kind(Husky h) {
  defer onReturn, onExit;
  return h.kind;
}

std::string impl_Dog_Poodle_Kind(void *self) { return Kind(*(Poodle *)self); }

// Main:
// generated: false
int main() {
  defer onReturn, onExit;
  Poodle __temp_0_Poodle = {};
  Dog d = {
      .self = &__temp_0_Poodle,
      .Kind = impl_Dog_Poodle_Kind,
  };
  std::cout << "I should be a Poodle:"
            << " " << d.Kind(d.self) << std::endl;
  Husky __temp_1_Husky = {};
  d = {
      .self = &__temp_1_Husky,
      .Kind = impl_Dog_Husky_Kind,
  };
  std::cout << "I should be a Husky:"
            << " " << d.Kind(d.self) << std::endl;
  Bulldog __temp_2_Bulldog = {};
  d = {
      .self = &__temp_2_Bulldog,
      .Kind = impl_Dog_Bulldog_Kind,
  };
  std::cout << "I should be a Bulldog:"
            << " " << d.Kind(d.self) << std::endl;
  Bulldog __temp_3_Bulldog = {
      .kind = "bulldog_2",
  };
  d = {
      .self = &__temp_3_Bulldog,
      .Kind = impl_Dog_Bulldog_Kind,
  };
  std::cout << "I should be a SECOND Bulldog:"
            << " " << d.Kind(d.self) << std::endl;
  Husky __temp_4_Husky = {
      .kind = "special_husky",
  };
  Dog dd = {
      .self = &__temp_4_Husky,
      .Kind = impl_Dog_Husky_Kind,
  };
  std::cout << "I should be a special Husky:"
            << " " << dd.Kind(dd.self) << std::endl;
  Dog ddd;
  Husky h = {
      .kind = "husky 2",
  };
  ddd = {
      .self = &h,
      .Kind = impl_Dog_Husky_Kind,
  };
  Bulldog b = {
      .kind = "poodle?",
  };
  std::cout << Kind(h) << std::endl;
  std::cout << Kind(b) << std::endl;
  std::cout << "I should be another Husky:"
            << " " << ddd.Kind(ddd.self) << std::endl;
}