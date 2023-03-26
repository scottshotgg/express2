typedef struct vtable_Fail
{
  int (*val)(struct Fail *this_ptr);
  bool (*error)(struct Fail *this_ptr);
} vtable_Fail;

typedef struct vtable_IntFoo
{
  int (*val)(struct IntFoo *this_ptr);
  bool (*error)(struct IntFoo *this_ptr);
} vtable_IntFoo;

int Fail_val(struct Fail *this_ptr)
{
  return 0;
}

void IntFoo_ctor(struct IntFoo *this_ptr, int value)
{
  this_ptr->v = value;
}

int IntFoo_val(struct IntFoo *this_ptr)
{
  return this_ptr->v;
}

struct Fail
{
  vtable_Fail *vtable;
};

struct IntFoo
{
  vtable_IntFoo *vtable;
  int v;
};

int main()
{
  Fail fail;
}