// This file is used for the over-arching struct inhereted by all objects

#include <map>
#include "var.cpp"

struct __OBJECT__ {
  private:
    // Integrated map for the object
    map<var, var> m = {};

  public:
    // A few public functions
    var Get(var key) {
      return m[key];
    };

    void Set(var key, var value) {
      m[key] = value;
    };
};