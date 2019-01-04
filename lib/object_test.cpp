#include "object.cpp"
#include <iostream>
#include <memory>

struct myObject : __OBJECT__ {
  int a = 2;
};

void test(void* obj) {
  myObject objD = (*(myObject *)obj);

  std::cout << "objD.a is: " << objD.a << std::endl;
}

int main() {
  myObject a = {};

  std::cout << a.a << std::endl;

  // void* objPtr = &a;

  // std::cout << objPtr << std::endl;

  // test(objPtr);

  // auto sp = std::shared_ptr<myObject>(&a);
}