# coding: utf-8
import numpy as np


class Sigmoid:
    params = []

    def __init__(self):
        pass

    def forward(self, x):
        return 1 / (1 + np.exp(-x))


class Affine:
    params = []

    def __init__(self, W, b):
        self.params = [W, b]

    def forward(self, x):
        W, b = self.params
        return np.dot(x, W) + b


class TwoLayerNet:
    layers = []
    params = []

    def __init__(self, input_size, hidden_size, output_size):
        # input__sizeはサンプルデータの個数ではなく、特徴量（属性値）の個数である
        I, H, O = input_size, hidden_size, output_size

        # 重みとバイアスの初期化
        W1 = np.random.randn(I, H)  # I個のデータからH個の中間データを生成
        b1 = np.random.randn(H)     # H個の中間データを生成するのでH
        W2 = np.random.randn(H, O)  # H個のデータからO個の出力データを生成
        b2 = np.random.randn(O)     # O個のデータを生成するのでO

        # レイヤの生成
        self.layers = [
            Affine(W1, b1),
            Sigmoid(),
            Affine(W2, b2),
        ]

        # 全ての重みをリストにまとめる
        for l in self.layers:
            self.params += l.params

    def predict(self, x):
        print(x)
        for l in self.layers:
            x = l.forward(x)
            print(l.__class__)
            print(x)
        return x


# 2個の特徴を持つ10個のデータを用いて推論
x = np.random.randn(10, 2)
model = TwoLayerNet(2, 4, 3)
s = model.predict(x)

print('x', x)
print('s', s)
print('model.params')
for p in model.params:
    print(p)
