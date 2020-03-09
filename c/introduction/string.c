#include <stdio.h>

void main()
{
  char month[] = "January";
  printf("This month is %s.\n", month);

  char name[20];
  printf("What's your name?: ");
  scanf("%19s", name);
  printf("Your name is %s.\n", name);

  // gets(name);
  // puts(name);

  // puts(gets(name));

  char meibo[4][20];

  int i;
  for (i = 0; i < 4; i++)
  {
    printf("%d人目の名前を入力してください:", i + 1);
    scanf("%19s", meibo[i]);
  }
  puts("\n");

  for (i = 0; i < 4; i++)
  {
    printf("%d人目は[%s]さんです\n", i + 1, meibo[i]);
  }
}