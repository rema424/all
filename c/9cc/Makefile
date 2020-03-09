CFLAGS=-std=c11 -g -static

9cc: 9cc.c
	docker run --rm -v ~/dev/all/c/9cc:/9cc -w /9cc compilerbook gcc -o 9cc 9cc.c

run: 9cc
	docker run --rm -v ~/dev/all/c/9cc:/9cc -w /9cc compilerbook ./9cc ${arg}

test: 9cc
	./test.sh

clean: rm -f 9cc *.o *~ tmp*

.PHONY: test clean

builds:
	docker run --rm -v ~/dev/all/c/9cc:/9cc -w /9cc compilerbook gcc -o tmp tmp.s

runs: builds
	docker run --rm -v ~/dev/all/c/9cc:/9cc -w /9cc compilerbook ./tmp; echo $$?
