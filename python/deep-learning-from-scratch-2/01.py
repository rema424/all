import numpy as np
import matplotlib.pyplot as plt
from common import functions

x = np.array([1, 2, 3])
print(x)
print(x.__class__)
print(x.shape)
print(x.ndim)
print('==================================')

W = np.array([
    [1, 2, 3],
    [4, 5, 6]
])
print(W)
print(W.shape)
print(W.ndim)
print('==================================')

W = np.array([
    [1, 2, 3],
    [4, 5, 6]
])
X = np.array([
    [0, 1, 2],
    [3, 4, 5]
])
print(W)
print(X)
print(X + W)
print(X * W)
print('==================================')

A = np.array([
    [1, 2],
    [3, 4],
])
print(A)
print(A * 10)
print('==================================')

A = np.array([
    [1, 2],
    [3, 4],
])
b = np.array([10, 20])
print(A)
print(b)
print(A * b)
print('==================================')
a = np.array([1, 2, 3])
b = np.array([4, 5, 6])
print(a)
print(b)
print(np.dot(a, b))
print('==================================')
A = np.array([[1, 2], [3, 4]])
B = np.array([[5, 6], [7, 8]])
print(A)
print(B)
print(np.dot(A, B))
print('==================================')
W1 = np.random.randn(2, 4)
b1 = np.random.randn(4)
x = np.random.randn(10, 2)
h = np.dot(x, W1) + b1
# plt.hist(W1)
# plt.show()
print(x)
print(W1)
print(b1)
print(h)
print('==================================')
x = np.random.randn(10, 2)
W1 = np.random.randn(2, 4)
b1 = np.random.randn(4)
W2 = np.random.randn(4, 3)
b2 = np.random.randn(3)
h = np.dot(x, W1) + b1
a = functions.sigmoid(h)
s = np.dot(a, W2) + b2
print(x)
print(W1)
print(b1)
print(W2)
print(b2)
print(h)
print(a)
print(s)
print('==================================')
print()
print()
print()
