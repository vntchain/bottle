#include "vntlib.h"

// typedef struct
// {
//     string ccc;
//     struct
//     {
//         string aaa;
//     } bbb;
// } S1;

typedef struct
{
    string ccc;
    mapping(string, string) aaa;
} S1;

KEY mapping(string, S1) b;

constructor test() {}
