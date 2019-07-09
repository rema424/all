```sh
brew tap heroku/brew && brew install heroku

heroku update

mkdir $GOPATH/src/github.com/rema424/all/heroku-go

cd $_

# heroku login
heroku login -i

# heroku apps:create heroku-go --buildpack heroku/go
heroku apps:create --buildpack heroku/go

heroku apps

heroku rename heroku-go-echo -a rocky-depths-41914

git init

heroku git:remote -a heroku-go-echo

heroku addons:add cleardb

go get -u github.com/kardianos/govendor

govendor init

govendor fetch +out

echo "vendor/*" >> .gitignore
echo '!vendor/vendor.json' >> .gitignore
```

heroku git サブモジュールは認証が必要です
https://codeday.me/jp/qa/20190322/458912.html
