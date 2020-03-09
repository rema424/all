import time


def make_generator():
    print('yield 1')
    yield 1
    print('yield 2')
    yield 2
    print('yield 3')
    yield 3


for v in make_generator():
    time.sleep(v)

gen = make_generator()
time.sleep(gen.__next__())
time.sleep(gen.__next__())
time.sleep(gen.__next__())
