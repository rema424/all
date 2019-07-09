```sh
# アプリの作成
mkdir $GOPATH/src/github.com/rema424/all/heroku-go
cd $_
go get -u github.com/kardianos/govendor
govendor init
# v4 は Go Modules でしか使えない?ため
govendor fetch github.com/labstack/echo@v3.3.10
govendor fetch github.com/labstack/echo/middleware@v3.3.10
govendor fetch +out
go run main.go

# git
git init
echo "vendor/*" >> .gitignore
echo '!vendor/vendor.json' >> .gitignore
echo "web: $(basename `pwd`)" > Procfile
git add .
git commit -m "Initial commit"

# heroku cli
brew tap heroku/brew && brew install heroku
heroku update
heroku login -i

# heroku app
heroku apps:create heroku-go-echo -b heroku/go
# heroku apps:create --buildpack heroku/go
heroku apps
heroku git:remote -a heroku-go-echo
git push heroku master
# heroku rename heroku-go-echo -a rocky-depths-41914
heroku open -a heroku-go-echo

heroku addons:add cleardb
heroku addons:docs cleardb
```

heroku git サブモジュールは認証が必要です
https://codeday.me/jp/qa/20190322/458912.html
