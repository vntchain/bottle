#include "vntlib.h"

constructor test() {}

MUTABLE
int32 test_export_1() {
  SendFromContract(AddressFrom("0xaaa"), U256From("1000000"));
}

UNMUTABLE
int32 test_export_2() {
  SendFromContract(AddressFrom("0xaaa"), U256From("1000000"));
}

int test_export_5() {}
