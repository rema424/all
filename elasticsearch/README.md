# Elasticsearch

## インストール

- [オフィシャル版 Elastic Homebrew タップをリリース](https://www.elastic.co/jp/blog/official-elastic-homebrew-tap-package-manager-macos)
- [elastic/go-elasticsearch](https://github.com/elastic/go-elasticsearch)
- [Elastic Stack and Product Documentation](https://www.elastic.co/guide/index.html)
- [Elastic Stack で作る BI 環境　誰でもできるデータ分析入門 記事一覧](https://thinkit.co.jp/series/7517)
- [初めての Elasticsearch with Docker](https://qiita.com/kiyokiyo_kzsby/items/344fb2e9aead158a5545)
- [Elasticsearch + Kibana を docker-compose でさくっと動かす](https://qiita.com/nobuman/items/6308ea3bfd0aa0c58fdb)
- []()
- []()
- []()
- []()
- []()

```sh
brew tap elastic/tap

# デフォルト版
brew install elasticsearch-full
#OSS配布
brew install elasticsearch-oss
# ==> Caveats
# Data:    /usr/local/var/lib/elasticsearch/elasticsearch_rema/
# Logs:    /usr/local/var/log/elasticsearch/elasticsearch_rema.log
# Plugins: /usr/local/var/elasticsearch/plugins/
# Config:  /usr/local/etc/elasticsearch/

# To have launchd start elastic/tap/elasticsearch-full now and restart at login:
#   brew services start elastic/tap/elasticsearch-full
# Or, if you don't want/need a background service you can just run:
#   elasticsearch

brew install kibana-full
# ==> Caveats
# Config: /usr/local/etc/kibana/
# If you wish to preserve your plugins upon upgrade, make a copy of
# /usr/local/opt/kibana-full/plugins before upgrading, and copy it into the
# new keg location after upgrading.

# To have launchd start elastic/tap/kibana-full now and restart at login:
#   brew services start elastic/tap/kibana-full
# Or, if you don't want/need a background service you can just run:
#   kibana

brew cask install homebrew/cask-versions/adoptopenjdk8
# ==> Running installer for adoptopenjdk8; your password may be necessary.
# ==> Package installers may write to any location; options such as --appdir are ignored.

brew install logstash-full
# logstash-full: Java 1.8 is required to install this formula.
# Install AdoptOpenJDK 8 with Homebrew Cask:
#   brew cask install homebrew/cask-versions/adoptopenjdk8
# ==> Caveats
# Please read the getting started guide located at:
#   https://www.elastic.co/guide/en/logstash/current/getting-started-with-logstash.html

# To have launchd start elastic/tap/logstash-full now and restart at login:
#   brew services start elastic/tap/logstash-full
# Or, if you don't want/need a background service you can just run:
#   logstash

ls $(brew --prefix)/Homebrew/Library/Taps/elastic/homebrew-tap/Formula/*.rb

ls /usr/local/etc/elasticsearch/
cat /usr/local/etc/elasticsearch/elasticsearch.yml

elasticsearch
open http://localhost:9200

kibana
open http://localhost:5601
```
