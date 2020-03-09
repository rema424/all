#include <stdio.h>

void main()
{
  long a;
  char b;
  float c;

  printf("a:%lu\nb:%lu\nc:%lu\ndouble:%lu\n", sizeof(a), sizeof(b), sizeof(c), sizeof(double));

  int d[20];
  int i;
  for (i = 0; i < 20; i++)
  {
    d[i] = i;
    printf("%d\n", i);
  }
  for (i = 0; i < sizeof(d) / sizeof(int); i++)
  {
    d[i] = i;
    printf("%d\n", i);
  }
}