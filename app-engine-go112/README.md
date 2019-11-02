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

# 【Go】やはりお前らの DDD レイヤードアーキテクチャは間違っている

## はじめに

タイトルは半分釣りです。ごめんなさい。

ただ聞いてください。残りの半分は本当に問題提起です。

実際の問題提起は「お前らの DDD レイヤードアーキテクチャ**についての記事**は間違っている」ですけど。

## 問題提起？

「ドメイン駆動設計」ですよね？？？

にも関わらず「DDD で実装してみました」系の記事でドメインに着目して解説してる記事少なくないですか？？？

ディレクトリ構成だったり依存の方向性だったり、技術に注目した記事ばかりじゃないですか？？？

それじゃ「ドメイン駆動」じゃなくて「技術駆動」じゃないですか？？？

何かしらの業務を仮定した上でそれをどのようにプログラムに落とし込んでいくか、その時ディレクトリ構成はどうなるか。これが本当に見つけたいことなのではないでしょうか。

これが問題提起です。

ということで、言い出しっぺの法則でこの記事の残りの部分ではシステム化対象とする業務を中心に添えつつ DDD のレイヤードアーキテクチャについて考えて行きたいと思います。

題材とするのは**小売業務**で、開発に用いる言語は Go です。

## 業務の仮定

システム化対象とする仮定の業務を次に示します。

> 電話「プルルルル。プルルルル。ガチャ」
>
> 店員「はい。〇〇ストアです。」
>
> 客「すみません。商品を購入したいのですが、○○ と □□ と △△ の在庫って今ありますか？」
>
> 店員「確認しますね。少々お待ちください。」
>
> （店員確認中...）
>
> 店員「お待たせしました。全て在庫ございますよ。ただ □□ と △△ は残りわずかで店内に並べてあるものが最後となっています。」
>
> 客「おお！ぜひ購入したいです！30 分後にお店に伺うのでお取り置きしてもらうことって可能ですか？？」
>
> 店員「承りますよ。それぞれ何個ずつお取り置きしますか？」
>
> 客「それぞれ 2 個ずつお願いします！」
>
> 店員「かしこまりました。ただいま在庫を確保しますので少々お待ちください。」
>
> （店員、陳列棚と倉庫から商品を確保してレジ裏に運び中...）
>
> 店員「お待たせしました。在庫の確保ができたのでお取り置きしておきます。お名前と電話番号を伺ってもよろしいでしょうか。」
>
> 客「はい！名前が ◇◇◇◇ で、電話番号が xxx-xxxx-xxxx です！」
>
> 店員「それではお待ちしております。なお、□□ と △△ は人気商品となっていますので、今から 1 時間以内に受け取りにいらっしゃらない場合はお取り置きをやめて店頭に再度陳列しますのでご了承ください。」
>
> 客「わかりました！それでは 30 分後に伺います！よろしくお願いします！」
>
> 〜 完 〜

以上が今回題材とする小売業務の内容です。

上記のシチュエーションではまだ売買は成立していないので、注文や購入ではなく**予約**になるかと思います。

## オブジェクトの抽出

それでは次にオブジェクトになりそうなものを抽出してみます。

- 顧客（customer）
- 店員（employee）
- 商品（item）
- 在庫（stock）
- 注文（order）
- 注文詳細（order_detail）
- 予約（reservation）

「商品」と「在庫」については補足が必要です。ここでは、商品は「概念としての商品」を指していて、在庫は「実体としての商品」を指しています。例を挙げると、たとえば Mac Book Pro の予約をする場合、購入側は概念（モデル）としての Mac Book Pro を指して予約し、販売側はシリアルナンバーのついた個別具体の Mac Book Pro の在庫を確保します。言い方を変えると、購入側は「このシリアルナンバーの Mac Book Pro を下さい」のように個別具体の商品を指して予約はしません。（なお、これは新品市場での話であって、中古市場だと個別具体を指して予約や注文をすることが一般的かと思います。）DDD においては、このような概念としてのオブジェクトを「値オブジェクト」、個別具体のオブジェクトを「エンティティ」と名付けて区別しています。RDB の設計においては、概念の方を「商品マスタ」、実体の方を「商品テーブル」と名付けて区別しているプロジェクトもあるようです。本記事においては概念の方を「商品」、実体の方を「在庫」と呼んで扱っていきます。

抽出したオブジェクトをコードにしてみます。

```go
type Customer struct {
	ID          int
	Name        string
	PhoneNumber int
}

type Employee struct {
	ID   int
	Name string
}

type Item struct {
	ID    int
	Name  string
	Price int
}

type Items []Item

type Order struct {
	ID           int
	OrderDetails OrderDetails
	Customer     Customer
	Employee     Employee
	CreatedAt    time.Time
}

type OrderDetail struct {
	OrderID  int
	Item     Item
	Quantity int
}

type OrderDetails []OrderDetail

type Stock struct {
	ID     int
	Status StockStatus
	Item   Item
}

type Stocks []Stock

type StockStatus int

const (
	StockStatusOnSale StockStatus = 1
	StockStatusReserved StockStatus = 2
	StockStatusSoldOut StockStatus = 3
)

type Reservation struct {
	ID       int
	OrderID  int
	Stocks   Stocks
	ExpireAt time.Time
}
```
