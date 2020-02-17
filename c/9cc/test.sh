#!/bin/bash
try() {
  expected="$1"
  input="$2"

  ./9cc "$input" > tmp.s
  docker run --rm -v ~/dev/all/c/9cc:/9cc -w /9cc compilerbook gcc -o tmp tmp.s
  docker run --rm -v ~/dev/all/c/9cc:/9cc -w /9cc compilerbook ./tmp
  actual="$?"

  if [ "$actual" = "$expected" ]; then
    echo "$input => $actual"
  else
    echo "$input => $expected expected, but got $actual"
    exit 1
  fi
}

try 0 0
try 42 42

echo OK
