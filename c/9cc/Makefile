build:
	docker run --rm -v ~/dev/all/c/9cc:/9cc -w /9cc compilerbook gcc -o 9cc 9cc.c

run: build
	docker run --rm -v ~/dev/all/c/9cc:/9cc -w /9cc compilerbook ./9cc ${arg}
	echo $?

builds:
	docker run --rm -v ~/dev/all/c/9cc:/9cc -w /9cc compilerbook gcc -o tmp tmp.s

runs: builds
	docker run --rm -v ~/dev/all/c/9cc:/9cc -w /9cc compilerbook ./tmp; echo $?;
