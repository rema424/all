# Chat

```sh
brew install redis

To have launchd start redis now and restart at login:
  brew services start redis
Or, if you don't want/need a background service you can just run:
  redis-server /usr/local/etc/redis.conf
```



```sh
redis-server

redis-cli

get foo

set foo bar

del foo

get foo

quit
```