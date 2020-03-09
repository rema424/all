```sh
docker run --rm -v ~/dev/all/c/example:/9cc -w /9cc compilerbook gcc -o test2 test2.s
docker run --rm -v ~/dev/all/c/example:/9cc -w /9cc compilerbook ./test2; echo $?;
```
