#include "lib/var.cpp"
int main() {
  int its = 6;
  object something = {};
  something["hey"] = "hey";
  something["its"] = its;
  something["me"] = {};
  something["me"]["woah"] = {};
  something["me"]["woah"]["another"] = "one";
  ;
  something["me"]["what"] = "yeah";
  something["me"]["yeah"] = {6, 6, 6};
  something["me"]["yeah"] = 87;
  ;
  ;
}