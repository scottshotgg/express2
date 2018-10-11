#include <string>
int main() {
  float a = 6.6;
  a = 7;
  float f = 0;
  std::string s = "";
  {
    std::string s = "scott";
    a = 9;
  };
  s = "me";
}