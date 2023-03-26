#include <string>

typedef struct
{
  void *self;
  std::string (*const Breed)(void *self);
} Dog;

struct Bulldog
{
  std::string breed = "bulldog";
};

struct Husky
{
  std::string breed = "husky";
};

std::string Breed(Bulldog *b)
{
  return b->breed;
}

std::string Breed(Husky *h)
{
  return h->breed;
}

std::string impl_dog_bulldog_Breed(void *self)
{
  return Breed((Bulldog *)self);
}

std::string impl_dog_husky_Breed(void *self)
{
  return Breed((Husky *)self);
}

Dog impl_dog_husky(Husky *h)
{
  return Dog{
      .self = h,
      .Breed = impl_dog_husky_Breed,
  };
}

Dog impl_dog(void *self, std::string breed(void *self))
{
  return Dog{
      .self = self,
      .Breed = breed,
  };
}

/*
  Steps:
    1. Allow transpilation of interface; i.e, Dog above
    //2. Check that struct implements interface type only when
    //   used as an interface
    3. Create struct wrappers when used as interface
    3. Wrap struct functions when used
    4. [ interface . call() ] => [ interface . call(self) ]
*/

int main()
{
  Husky h;

  Dog d = impl_dog(
      &h,
      impl_dog_husky_Breed);

  printf("%s", d.Breed(d.self).c_str());
}