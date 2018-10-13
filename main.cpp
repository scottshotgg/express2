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
  int i = 0;
  bool b = false;
  float f = 0;
  char c = 0;
  std::string s = "";
  object o = {};
  o["me"] = "s";
  o["thing"] = 8;
  o["me"] = "9";
  var json = {};
  json["a"] = 1;
  json["b"] = 2;
  json["c"] = 3;
  var eight = 8;
  var v = 0;
  std::string scott = "";
  {
    std::string scott = "scott";
    a = 9;
  }
  scott = "me";
}