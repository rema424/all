```sh
brew install python

pip3 install scrapy

pip3 install flake8

pip3 install autopep8

pip3 install pipdeptree

pip3 show scrapy | grep Location | cut -d " " -f 2

vi settings.json
  # "python.autoComplete.extraPaths": [
  #   "/usr/local/lib/python3.7/site-packages"
  # ],
  #  "prettier.disableLanguages": [
  #   "vue", "python"
  # ],
  # "[python]": {
  #   "editor.defaultFormatter": "ms-python.python"
  # },
  # "python.linting.pylintEnabled": false,
  # "python.linting.flake8Enabled": true,
  # "python.linting.lintOnSave": true,
  # "python.formatting.provider": "autopep8",
  # "editor.formatOnSave": true,

scrapy startproject tutorial

cd tutorial

scrapy genspider quotes quotes.toscrape.com

vi tutorial/items.py

vi tutorial/spiders/quotes.py

vi tutorial/settings.py

scrapy runspider tutorial/spiders/quotes.py

scrapy runspider tutorial/spiders/quotes.py -t csv -o stdout:

scrapy runspider tutorial/spiders/quotes.py -t csv -o stdout: --nolog

scrapy runspider tutorial/spiders/quotes.py -t json -o stdout: --nolog

scrapy runspider tutorial/spiders/quotes.py -o out_$(date "+%Y-%m-%d_%H:%M").css

scrapy crawl quotes

scrapy crawl quotes -t csv -o stdout:

scrapy crawl quotes -t csv -o stdout: --nolog

scrapy crawl quotes -t json -o stdout: --nolog

scrapy crawl quotes -o out_$(date "+%Y-%m-%d_%H:%M").css

scrapy shell http://quotes.toscrape.com

ctrl + D

cd ..

scrapy startproject cartune cartune-test

cd cartune-test

scrapy genspider pickups cartune.me

vi items.py

vi pickups.py

scrapy crawl pickups -t json -o stdout: --nolog

vi settings.py
# FEED_EXPORT_ENCODING = 'utf-8'

scrapy crawl pickups -t json -o stdout: --nolog

scrapy genspider pickups_detail cartune.me

vi items.py
```
