# coding: utf-8


class SGD:
    '''
    確率的勾配降下法（Stochastic Gradient Descent）
    '''

    def __init__(self, learning_rate=0.01):
        self.learning_rate = learning_rate

    def update(self, params, grads):
        for i in range(len(params)):
            params[i] -= self.learning_rate * grads[i]
