#include <stdio.h>

// Concrete student declaration
struct Student
{
  char name[8] = "student";
};

// Concrete Researcher declaration
struct Researcher
{
  char name[11] = "researcher";
};

// Interface Person declaration
typedef struct
{
  // Self reference
  void *self;

  // vtable
  char *(*const name_fn)(void *self);
} Person;

// Researcher implementation
char *name_fn_researcher(void *self)
{
  Researcher *r = (Researcher *)self;
  return r->name;
}

// Student implementation
char *name_fn_student(void *self)
{
  Student *s = (Student *)self;
  return s->name;
}

int main()
{
  // Concrete
  Student s;

  // Interface
  Person p = {
      // Assign self as well as the function pointer
      .self = &s,
      .name_fn = name_fn_researcher,
  };

  printf("%s", p.name_fn(p.self));
}
