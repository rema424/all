# coding: utf-8
import numpy as np


def sigmoid(x):
    return 1 / (1 + np.exp(-x))


def softmax_1(x):
    exp_x = np.exp(x)
    sum_exp_x = np.sum(exp_x)
    return exp_x / sum_exp_x


def softmax_2(x):
    # print('x.ndim', x.ndim)
    # print('x', x)
    # オーバーフロー対策として最大値で減産する
    if x.ndim == 1:
        max = np.max(x)
        x = x - max
        exp_x = np.exp(x)
        result = exp_x / np.sum(exp_x)
        # print('np.max(x)', max)
        # print('x - np.max(x)', x)
        # print('exp_x', exp_x)
        # print('result', result)
        return result
    elif x.ndim == 2:
        max = x.max(axis=1, keepdims=True)
        x = x - max
        exp_x = np.exp(x)
        result = exp_x / exp_x.sum(axis=1, keepdims=True)
        # print('np.max(x)', max)
        # print('x - np.max(x)', x)
        # print('exp_x', exp_x)
        # print('result', result)
        return result

def softmax_3(x):
    if x.ndim == 2:
        x = x - x.max(axis=1, keepdims=True)
        x = np.exp(x)
        x /= x.sum(axis=1, keepdims=True)
    elif x.ndim == 1:
        x = x - np.max(x)
        x = np.exp(x) / np.sum(np.exp(x))

    return x


def cross_entropy_error(y, t):
    # yはニューラルネットワークの出力
    # tは教師データ
    # if y.ndim == 1:
    #     t = t.reshape(1, t.size)
    #     y = y.reshape(1, y.size)

    # batch_size = y.shape[0]
    # delta = 1e-7
    # return -np.sum(np.log(y[np.arange(batch_size), t] + delta)) / batch_size
    if y.ndim == 1:
        t = t.reshape(1, t.size)
        y = y.reshape(1, y.size)

    # 教師データがone-hot-vectorの場合、正解ラベルのインデックスに変換
    if t.size == y.size:
        t = t.argmax(axis=1)

    batch_size = y.shape[0]

    return -np.sum(np.log(y[np.arange(batch_size), t] + 1e-7)) / batch_size


# a = np.array([0.3, 2.9, 4.0])
# y = softmax_2(a)
# print(a)
# print(y)
# print(np.sum(y))
# a = np.array([
#     [0.3, 2.9, 4.0],
#     [0.7, 2.3, 5.4]
# ])
# y = softmax_2(a)
# print(a)
# print(y)
# print(y.sum(axis=1, keepdims=True))
