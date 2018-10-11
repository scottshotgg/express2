#include "lib/var.cpp"
#include <functional>
#include <string>
int main() {
  std::function<void()> something = []() {
    std::string stuff = "woah";
    var thing = "yeah";
  };
  float a = 6.6;
  a = 7;
  var thing = 0;
  float f = 0;
  std::string s = "";
  {
    std::string s = "scott";
    a = 9;
  };
  s = "me";
}