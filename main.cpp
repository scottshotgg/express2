#include "lib/var.cpp"
int main() {
  var something = "crazy";
  var k = {};
  ;
  for (int i = 0; i < 10; i++) {
    something = i;
    var j = {false, 2.3456789, "1.123456789", 8};
    k = {};
    k["a"] = i;
    k["b"] = j;
    k["c"] = 3;
  }
}