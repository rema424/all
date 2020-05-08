# coding: utf-8
import sys
sys.path.append('deep-learning-from-scratch-2/')
# sys.path.append('..')

from dataset import spiral
import matplotlib.pyplot as plt

x, t = spiral.load_data(1)
print('x', x.shape, x)
print('t', t.shape, t)

# データ点のプロット
N = 100
CLS_NUM = 3
markers = ['o', 'x', '^']
for i in range(CLS_NUM):
    plt.scatter(x[i*N:(i+1)*N, 0], x[i*N:(i+1)*N, 1], s=40, marker=markers[i])
plt.show()
