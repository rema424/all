```sh
cd /Applications
git clone https://github.com/luin/medis.git && cd medis
npm install
npm run build
npm start
```

```sh
redis-cli --raw

# 記事IDをインクリメント
incr article:

# 記事:1 に対する投票ユーザーをSETに追加
sadd voted:1 user:1
sadd voted:1 user:2
sadd voted:1 user:3
smembers voted:1

# 有効時間を設定
expire voted:1 15
ttl voted:1
smembers voted:1

# もう一度追加
sadd voted:1 user:1 user:2 user:3

# 記事のハッシュを作る
hmset article:1 title: "タイトルです" link "リンクです" poster "user:4" time "今です" votes 3
hgetall article:1
del article:1
hgetall article:1



zadd time: 1571882597 article:1

zscore time: article:1
```
