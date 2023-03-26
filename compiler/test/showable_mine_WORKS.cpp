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
char *name_fn_researcher(Researcher *r)
{
  return r->name;
}

// Student implementation
char *name_fn_student(Student *s)
{
  return s->name;
}

// Wrapper function for interface usage
char *wrap_name_fn_researcher(void *self)
{
  return name_fn_researcher((Researcher *)self);
}

// Wrapper function for interface usage
char *wrap_name_fn_student(void *self)
{
  return name_fn_student((Student *)self);
}

int main()
{
  // Concrete
  Student s;

  // Interface
  Person p = {
      // Assign self as well as the function pointer
      .self = &s,
      .name_fn = wrap_name_fn_researcher,
  };

  printf("%s", p.name_fn(p.self));
}
