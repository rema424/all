#include <stdio.h>

void main()
{
  printf("Hello C world.\n");

  int n = 10;
  char ch = 'A';
  printf("変数nの値は%dです。\n", n);
  printf("変数chの値は%cです。\n", ch);

  double pi = 3.1415926;
  printf("変数nの値は%+010dです。\n", n);
  printf("円周率は%+010.5fです。\n", pi);
}