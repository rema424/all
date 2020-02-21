#include <stdio.h>

void main()
{
  int n;
  char ch;

  printf("数字を入力してください --> ");
  scanf("%d", &n);
  printf("文字を入力してください --> ");
  scanf("%c", &ch);

  printf("\n\n");
  printf("入力された数字は%dです。\n", n);
  printf("入力された文字は%cです。\n", ch);

  scanf("%x", &n);
  printf("%d\n", n);
}