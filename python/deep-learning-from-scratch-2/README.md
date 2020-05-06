## セットアップ

```sh
git clone https://github.com/oreilly-japan/deep-learning-from-scratch-2

virtualenv -p python3 venv

. venv/bin/activate

pip install numpy
```

## numpy

```sh
python

import numpy as np

x = np.array([1, 2, 3])
x.__class__
x.shape
x.ndim

```