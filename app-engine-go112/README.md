# APP-ENGINE-GO112

## 準備

- Google Account の作成
- GCP プロジェクトの作成
- Cloud SDK のセットアップ

### Google Account の作成

https://accounts.google.com/signup

### GCP プロジェクトの作成

https://console.cloud.google.com/projectselector2/home/dashboard

### Cloud SDK のセットアップ

```sh
curl https://sdk.cloud.google.com | bash
exec -l $SHELL
gcloud init
```

### 課金プロジェクトのリンク

https://console.cloud.google.com/projectselector2/billing

### App Engine のセットアップ

```sh
gcloud components update
gcloud config list
gcloud projects describe [YOUR_PROJECT_ID]
gcloud app create --project=[YOUR_PROJECT_ID]
gcloud components install app-engine-go
```

### 各種操作

```sh
# ローカル実行
go run appengine/default/main.go

# デプロイ
gcloud app deploy appengine/default/app.yaml
```

### 逆引き DDD

- トランザクション（データの整合性）はどこで担保する？
- データストアを複数跨ぐ処理はどうする？（RDB とオブジェクトストレージ）
