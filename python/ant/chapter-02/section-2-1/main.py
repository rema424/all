# coding: utf-8
from timeit import timeit


def fact(n: int) -> int:
    if n == 1:
        return 1
    return n * fact(n - 1)


def fib(n: int) -> int:
    if n <= 1:
        return n
    return fib(n-1) + fib(n-2)


def fib2(n: int) -> int:
    if n <= 1:
        return n

    memo = [0] * (n+1)

    def _fib(n: int) -> int:
        if n <= 1:
            return n
        if memo[n] != 0:
            return memo[n]
        memo[n] = _fib(n-1) + _fib(n-2)
        return memo[n]
    return _fib(n)


if __name__ == "__main__":
    print("="*30)
    print(fact(4))
    print(fact(3))
    print(fact(2))
    print("="*30)
    # print([fib(n) for n in range(33)])
    print("="*30)
    print([fib2(n) for n in range(33)])
    print("="*30)
