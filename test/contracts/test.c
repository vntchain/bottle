#include "../../vntlib.h"

KEY int64 testa;
KEY mapping(string, string) a;
bool case1() {
  a.key = "testkey1";
  KEY string *tmpmap = &a.value; //临时变量
  tmpmap = "testvalue";
  return false;
}

string *get_a_value(int xx) {
  if (xx == 1 + 1) {
    return &a.value + 1;
  } else {
    return get_a_value(xx);
  }
}

UNMUTABLE
bool case2() {
  a.key = "testkey1";
  KEY string *tmpmap = get_a_value(1); //临时变量
  tmpmap = "testvalue";
  return false;
}
