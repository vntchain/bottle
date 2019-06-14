#include "vntlib.h"

KEY mapping(string, string) key1;
KEY array(string) key2;
KEY struct
{
  string a;
} key3;

typedef struct
{
  string a;
} s1;

KEY s1 key4;

struct s2
{
  string a;
};
KEY struct s2 key5;

KEY mapping(string, mapping(string, string)) key6;

KEY array(array(string)) key7;

KEY mapping(string, s1) key8;

KEY mapping(string, struct s2) key9;

KEY array(array(s1)) key10;

KEY array(array(struct s2)) key11;

typedef struct
{
  string a;
  mapping(string, string) b;
} s3;

KEY s3 key12;

KEY mapping(string, s3) key13;

struct s4
{
  string a;
  mapping(string, string) b;
};

KEY struct s4 key14;

KEY mapping(string, struct s4) key15;

constructor test()
{
}
