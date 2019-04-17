#include "vntlib.h"

constructor test() {}

MUTABLE
int32 test_export_1() {}

UNMUTABLE
int32 test_export_2() { test_export_1(); }

int test_export_5() {}
