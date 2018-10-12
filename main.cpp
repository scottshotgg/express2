#include "lib/var.cpp"
#include <functional>
#include <string>
void something() {
  std::string stuff = "woah";
  var thing = "yeah";
}
int main() {
  float a = 6.6;
  a = 7;
  var v = 0;
  int i = 0;
  bool b = false;
  float f = 0;
  char c = 0;
  std::string s = "";
  std::string scott = "";
  {
    std::string scott = "scott";
    a = 9;
  };
  scott = "me";
}