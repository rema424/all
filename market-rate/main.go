package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql" // Using MySQL driver
	"github.com/jmoiron/sqlx"
	"gonum.org/v1/gonum/stat"
)

var db *sqlx.DB

type Price struct {
	Word   string  `db:"word"`
	Price  float64 `db:"price"`
	Title  string  `db:"title"`
	Url    string  `db:"url"`
	Site   string  `db:"site"`
	ItemID int     `db:"item_id"`
}

type Rate struct {
	ID         int     `db:"id"`
	PartsName  string  `db:"parts_name"`
	Memo       string  `db:"memo"`
	NarrowRate float64 `db:"narrow_rate"`
	SampleSize int     `db:"sample_size"`
	Skew       float64 `db:"skew"`
	Stddev     float64 `db:"stddev"`
	Min        float64 `db:"min"`
	Max        float64 `db:"max"`
	Lower      float64 `db:"lower"`
	Upper      float64 `db:"upper"`
	Median     float64 `db:"median"`
	Mean       float64 `db:"mean"`
	Ratio      float64 `db:"ratio"`
	ItemID     int     `db:"item_id"`
}

func (r *Rate) adaptWheel() {
	r.Min = r.Min / 3 * 4
	r.Max = r.Max / 3 * 4
	r.Lower = r.Lower / 3 * 4
	r.Upper = r.Upper / 3 * 4
	r.Median = r.Median / 3 * 4
	r.Mean = r.Mean / 3 * 4
	r.Stddev = r.Stddev / 3 * 4
}

func main() {
	db = DB()
	defer db.Close()

	// fmt.Println("【method】")
	// fmt.Println("[1] 通常のスクレイプ")
	// fmt.Println("[2] 2000円未満を弾いたスクレイプ")
	// fmt.Print("1 or 2 : ")

	// scanner := bufio.NewScanner(os.Stdin)
	// scanner.Scan()
	// method := scanner.Text()

	// if method != "1" && method != "2" {
	// 	panic("不正な値が入力されました。")
	// }

	file, err := os.Open("items.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// ヘッダー行を捨てる
	_, _ = reader.Read()

	query := `INSERT INTO prices
		(word, price, title, url, site)
	VALUES
		(?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE price = ?, title = ?;
	`
	var word, price, title, url, site string
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		price = line[1]
		site = line[2]
		title = line[3]
		url = line[4]
		word = line[5]

		if _, err := db.Exec(query, word, price, title, url, site, price, title); err != nil {
			panic(err)
		}
	}

	// csvにitem_idを紐づける
	bindItemIDQueryToPrices := `update prices, (select id, name, long_name from items) as tmp
	set prices.item_id = tmp.id
	where prices.word = tmp.name
		 or prices.word = tmp.long_name;
	`
	if _, err := db.Exec(bindItemIDQueryToPrices); err != nil {
		panic(err)
	}

	// パーツ一覧抽出
	wordsExtractQuery := `SELECT distinct(word) as word
	from prices
	having word not in (select distinct(parts_name) from rates);`

	var words []string
	if err := db.Select(&words, wordsExtractQuery); err != nil {
		panic(err)
	}
	fmt.Println(words)

	// パーツごとにループ処理して計算
	priceExtractQuery := `SELECT * from prices where word = ?;`

	rateCreateQuery := `INSERT IGNORE INTO rates
		(parts_name, memo, narrow_rate, sample_size, skew, stddev, min, lower, median, mean, upper, max, ratio)
	VALUES
		(:parts_name, :memo, :narrow_rate, :sample_size, :skew, :stddev, :min, :lower, :median, :mean, :upper, :max, :ratio);`

	var prices []*Price
	narrowRates := []float64{0.00, 0.05, 0.10, 0.15, 0.20}
	excludeWords := []string{"ジャンク", "欠品", "希少"}

	for _, word := range words {
		fmt.Println(word)
		if err := db.Select(&prices, priceExtractQuery, word); err != nil {
			panic(err)
		}

		// 通常
		for _, narrowRate := range narrowRates {
			rate := calc(prices, narrowRate, "フィルターなし")
			rate.PartsName = word
			if _, err := db.NamedQuery(rateCreateQuery, rate); err != nil {
				panic(err)
			}
		}

		// 単価1000円未満を除外
		pricesOver1000yen := make([]*Price, 0, len(prices))
		for _, price := range prices {
			if price.Price > 1000 {
				pricesOver1000yen = append(pricesOver1000yen, price)
			}
		}
		for _, narrowRate := range narrowRates {
			rate := calc(pricesOver1000yen, narrowRate, "1000円未満除外")
			rate.PartsName = word
			if _, err := db.NamedQuery(rateCreateQuery, rate); err != nil {
				panic(err)
			}
		}

		// 除外ワードでフィルター
		pricesFilteredByWords := make([]*Price, 0, len(pricesOver1000yen))
		for _, price := range pricesOver1000yen {
			var include bool

			for _, excludeWord := range excludeWords {
				if strings.Contains(price.Title, excludeWord) {
					include = true
					break
				}
			}

			if !include {
				pricesFilteredByWords = append(pricesFilteredByWords, price)
			}
		}
		for _, narrowRate := range narrowRates {
			rate := calc(pricesFilteredByWords, narrowRate, "禁止ワードを除外")
			rate.PartsName = word
			if _, err := db.NamedQuery(rateCreateQuery, rate); err != nil {
				panic(err)
			}
		}

		// メーカー名・パーツ名完全一致でフィルター
		pricesFilteredByPartsName := make([]*Price, 0, len(pricesFilteredByWords))
		for _, price := range pricesFilteredByWords {
			if strings.Contains(price.Title, word) {
				pricesFilteredByPartsName = append(pricesFilteredByPartsName, price)
			}
		}
		for _, narrowRate := range narrowRates {
			rate := calc(pricesFilteredByPartsName, narrowRate, "パーツ名完全一致")
			rate.PartsName = word
			if _, err := db.NamedQuery(rateCreateQuery, rate); err != nil {
				panic(err)
			}
		}
	}

	// ホイールの処理
	wheelExtractQuery := `SELECT distinct(prices.word) as word
	from prices
	inner join items on items.id = prices.item_id and items.is_wheel = 1
	having word not in (select distinct(parts_name) from wheel_rates);`

	wheelRateCreateQuery := `INSERT IGNORE INTO wheel_rates
		(parts_name, memo, narrow_rate, sample_size, skew, stddev, min, lower, median, mean, upper, max, ratio)
	VALUES
		(:parts_name, :memo, :narrow_rate, :sample_size, :skew, :stddev, :min, :lower, :median, :mean, :upper, :max, :ratio);`

	var wheels []string
	if err := db.Select(&wheels, wheelExtractQuery); err != nil {
		panic(err)
	}
	fmt.Println(wheels)

	for _, word := range wheels {
		fmt.Println(word)
		if err := db.Select(&prices, priceExtractQuery, word); err != nil {
			panic(err)
		}

		for _, price := range prices {
			price.Price = price.Price / 3 * 4
		}

		// 通常
		for _, narrowRate := range narrowRates {
			rate := calc(prices, narrowRate, "フィルターなし")
			rate.PartsName = word
			// rate.adaptWheel()
			if _, err := db.NamedQuery(wheelRateCreateQuery, rate); err != nil {
				panic(err)
			}
		}

		// 単価1000円未満を除外
		pricesOver1000yen := make([]*Price, 0, len(prices))
		for _, price := range prices {
			if price.Price > 80000 {
				pricesOver1000yen = append(pricesOver1000yen, price)
			}
		}
		for _, narrowRate := range narrowRates {
			rate := calc(pricesOver1000yen, narrowRate, "80000円未満除外")
			rate.PartsName = word
			// rate.adaptWheel()
			if _, err := db.NamedQuery(wheelRateCreateQuery, rate); err != nil {
				panic(err)
			}
		}

		// 除外ワードでフィルター
		pricesFilteredByWords := make([]*Price, 0, len(pricesOver1000yen))
		for _, price := range pricesOver1000yen {
			var include bool

			for _, excludeWord := range excludeWords {
				if strings.Contains(price.Title, excludeWord) {
					include = true
					break
				}
			}

			if !include {
				pricesFilteredByWords = append(pricesFilteredByWords, price)
			}
		}
		for _, narrowRate := range narrowRates {
			rate := calc(pricesFilteredByWords, narrowRate, "禁止ワードを除外")
			rate.PartsName = word
			// rate.adaptWheel()
			if _, err := db.NamedQuery(wheelRateCreateQuery, rate); err != nil {
				panic(err)
			}
		}

		// メーカー名・パーツ名完全一致でフィルター
		pricesFilteredByPartsName := make([]*Price, 0, len(pricesFilteredByWords))
		for _, price := range pricesFilteredByWords {
			if strings.Contains(price.Title, word) {
				pricesFilteredByPartsName = append(pricesFilteredByPartsName, price)
			}
		}
		for _, narrowRate := range narrowRates {
			rate := calc(pricesFilteredByPartsName, narrowRate, "パーツ名完全一致")
			rate.PartsName = word
			// rate.adaptWheel()
			if _, err := db.NamedQuery(wheelRateCreateQuery, rate); err != nil {
				panic(err)
			}
		}
	}

	// ratesにitem_idを紐づける
	bindItemIDQueryToRates := `update rates, (select id, name, long_name from items) as tmp
		set rates.item_id = tmp.id
		where rates.parts_name = tmp.name
			 or rates.parts_name = tmp.long_name;`
	if _, err := db.Exec(bindItemIDQueryToRates); err != nil {
		panic(err)
	}
	// wheel_ratesにitem_idを紐づける
	bindItemIDQueryToWheelRates := `update wheel_rates, (select id, name, long_name from items) as tmp
		set wheel_rates.item_id = tmp.id
		where wheel_rates.parts_name = tmp.name
			 or wheel_rates.parts_name = tmp.long_name;`
	if _, err := db.Exec(bindItemIDQueryToWheelRates); err != nil {
		panic(err)
	}
}

// DB gets a connection to the database.
func DB() *sqlx.DB {
	var dsn = mustGetenv("DSN")

	conn, err := sqlx.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Sprintf("DB: %v", err))
	}

	return conn
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("%s environment variable not set.", k)
	}
	return v
}

func minmax(prices []float64) (float64, float64) {
	switch len(prices) {
	case 0:
		return 0, 0
	case 1:
		return prices[0], prices[0]
	default:
		return prices[0], prices[len(prices)-1]
	}
}

func calc(prices []*Price, narrowRate float64, memo string) Rate {
	size := len(prices)
	cut := int(math.Floor(float64(size) * narrowRate))
	newPrices := prices[cut : len(prices)-cut]
	newSize := len(newPrices)
	nums := extractPrices(newPrices)
	sort.Float64s(nums)

	fmt.Println("nums", nums)

	// 歪度
	skew := stat.Skew(nums, nil)
	// 標準偏差
	stddev := stat.StdDev(nums, nil)
	// 中央値
	var median float64
	if newSize > 3 {
		median = stat.Quantile(0.5, stat.Empirical, nums, nil)
	} else if newSize == 2 {
		median = stat.Mean(nums, nil)
	} else {
		median = math.NaN()
	}
	// 平均値
	mean := stat.Mean(nums, nil)
	// 平均値 - 標準偏差
	lower := mean - stddev
	// 平均値 + 標準偏差
	upper := mean + stddev
	// 最小値・最大値
	min, max := minmax(nums)
	// ratio
	ratio := float64(count(nums, lower, upper)) / float64(newSize)

	fmt.Println("narrowRate:", narrowRate)
	fmt.Println("size:", newSize)
	fmt.Println(newPrices)
	fmt.Println("skew:", skew)
	fmt.Println("stddev:", stddev)
	fmt.Println("min:", min)
	fmt.Println("lower:", lower)
	fmt.Println("median:", median)
	fmt.Println("mean:", mean)
	fmt.Println("upper:", upper)
	fmt.Println("max:", max)
	fmt.Println("ratio:", ratio)

	return Rate{
		Memo:       memo,
		NarrowRate: narrowRate,
		SampleSize: newSize,
		Skew:       skew,
		Stddev:     stddev,
		Min:        min,
		Lower:      lower,
		Median:     median,
		Mean:       mean,
		Upper:      upper,
		Max:        max,
		Ratio:      ratio,
	}
}

func count(prices []float64, lower, upper float64) int {
	c := 0
	for _, price := range prices {
		if lower < price && price < upper {
			c++
		}
	}
	return c
}

func extractPrices(prices []*Price) []float64 {
	nums := make([]float64, len(prices))
	for i, price := range prices {
		nums[i] = price.Price
	}
	return nums
}

func scan() {
	f, err := os.Open("salary.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var xs []float64
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		var v float64
		txt := scan.Text()
		_, err = fmt.Sscanf(txt, "%f", &v)
		if err != nil {
			log.Fatalf(
				"could not convert to float64 %q: %v",
				txt, err,
			)
		}
		xs = append(xs, v)
	}

	// make sure scanning the file and extracting values
	// went fine, without any error.
	if err = scan.Err(); err != nil {
		log.Fatalf("error scanning file: %v", err)
	}

	fmt.Printf("data sample size: %v\n", len(xs))

	mean := stat.Mean(xs, nil)
	variance := stat.Variance(xs, nil)
	stddev := math.Sqrt(variance)

	sort.Float64s(xs)
	median := stat.Quantile(0.5, stat.Empirical, xs, nil)

	fmt.Printf("mean=     %v\n", mean)
	fmt.Printf("median=   %v\n", median)
	fmt.Printf("variance= %v\n", variance)
	fmt.Printf("std-dev=  %v\n", stddev)
}
