#include "vntlib.h"

KEY struct {
  int32 a1;
  int64 a2;
  uint64 a3;
  uint64 a4;
  uint256 a5;
  string a6;
  address a7;
  mapping(int32, int32) a8;
  array(int32) a9;
} a;
KEY struct {
  mapping(int32 *, int32) b1;
  array(int64 *) b2;
  string *b3;
  int b4;
  /*address *b5;*/
  // uint256 * b6;
} b;

constructor test() {}
