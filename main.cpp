#include "lib/var.cpp"
int main() {
  var something = "crazy";
  int crazy = 7;
  var json = {};
  json["a"] = 1;
  json["b"] = 2;
  json["c"] = 3;
  ;
  object a = {};
  a["a"] = {};
  a["a"]["hey"] = "its me";
  ;
  a["b"] = 0;
  a["c"] = 3;
  ;
  float float_array[] = {9.9, 9, 5.5};
  object array_of_objects = {};
  array_of_objects[0] = {};
  array_of_objects[1] = {};
  array_of_objects[2] = {};
  array_of_objects[3] = {};
  {};
  array_of_objects[0]["a"] = 8;
  {};
  array_of_objects[1]["b"] = 7;
  {};
  array_of_objects[2]["a"] = {};
  array_of_objects[2]["a"]["hey"] = "its me";
  ;
  array_of_objects[2]["b"] = 0;
  array_of_objects[2]["c"] = 3;
  {};
  array_of_objects[3]["a"] = {};
  array_of_objects[3]["a"]["hey"] = "its me";
  ;
  array_of_objects[3]["b"] = 0;
  array_of_objects[3]["c"] = 3;
  ;
}