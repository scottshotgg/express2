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
  std::string breed = "bulldog";
};
struct Husky {
  std::string breed = "husky";
};

// Interfaces:
typedef struct {
  void *self;
  std::string (*Breed)(void *self);
} Dog;

// Prototypes:
std::string impl_Dog_Poodle_Breed(void *self);
std::string impl_Dog_Husky_Breed(void *self);
std::string impl_Dog_Bulldog_Breed(void *self);
std::string Breed(Poodle p);
std::string Breed(Bulldog b);
std::string Breed(Husky *h);

// Functions:
std::string impl_Dog_Poodle_Breed(void *self) { return Breed(*(Poodle *)self); }

std::string impl_Dog_Husky_Breed(void *self) { return Breed((Husky *)self); }

std::string impl_Dog_Bulldog_Breed(void *self) {
  return Breed(*(Bulldog *)self);
}

std::string Breed(Poodle p) {
  defer onReturn, onExit;
  return "poodle";
}

std::string Breed(Bulldog b) {
  defer onReturn, onExit;
  return b.breed;
}

std::string Breed(Husky *h) {
  defer onReturn, onExit;
  return h->breed;
}

// Main:
// generated: false
int main() {
  defer onReturn, onExit;
  Poodle p = {};
  Dog[{}] = {p};
  Dog d = {
      .self = &p,
      .Breed = impl_Dog_Poodle_Breed,
  };
  std::cout << "I should be a Poodle:"
            << " " << d.Breed(d.self) << std::endl;
  Husky __temp_0_Husky = {};
  d = {
      .self = &__temp_0_Husky,
      .Breed = impl_Dog_Husky_Breed,
  };
  std::cout << "I should be a Husky:"
            << " " << d.Breed(d.self) << std::endl;
  Bulldog __temp_1_Bulldog = {};
  d = {
      .self = &__temp_1_Bulldog,
      .Breed = impl_Dog_Bulldog_Breed,
  };
  std::cout << "I should be a Bulldog:"
            << " " << d.Breed(d.self) << std::endl;
  Bulldog __temp_2_Bulldog = {
      .breed = "bulldog_2",
  };
  d = {
      .self = &__temp_2_Bulldog,
      .Breed = impl_Dog_Bulldog_Breed,
  };
  std::cout << "I should be a SECOND Bulldog:"
            << " " << d.Breed(d.self) << std::endl;
}