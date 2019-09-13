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

	_ "github.com/go-sql-driver/mysql" // Using MySQL driver
	"github.com/jmoiron/sqlx"
	"gonum.org/v1/gonum/stat"
)

var db *sqlx.DB

type Price struct {
	Word   string  `db:""`
	Price  float64 `db:""`
	Title  string  `db:""`
	Url    string  `db:""`
	Wite   string  `db:""`
	Method string  `db:""`
}

type Rate struct {
	ID         int     `db:"id"`
	PartsName  string  `db:"parts_name"`
	Method     string  `db:"method"`
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
}

func main() {
	db = DB()
	defer db.Close()

	fmt.Println("【method】")
	fmt.Println("[1] 通常のスクレイプ")
	fmt.Println("[2] 2000円未満を弾いたスクレイプ")
	fmt.Print("1 or 2 : ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	method := scanner.Text()

	if method != "1" && method != "2" {
		panic("不正な値が入力されました。")
	}

	file, err := os.Open("items.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// ヘッダー行を捨てる
	_, _ = reader.Read()

	query := `INSERT INTO prices
		(word, price, title, url, site, method)
	VALUES
		(?, ?, ?, ?, ?, ?)
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

		if _, err := db.Exec(query, word, price, title, url, site, method, price, title); err != nil {
			log.Println(err)
		}
	}

	// パーツ一覧抽出
	wordsExtractQuery := `SELECT distinct(word) from prices where method = ?;`
	var words []string
	if err := db.Select(&words, wordsExtractQuery, method); err != nil {
		panic(err)
	}
	fmt.Println(words)

	// パーツごとにループ処理して計算
	priceExtractQuery := `SELECT price from prices where word = ? and method = ?;`
	rateCreateQuery := `INSERT IGNORE INTO rates
		(parts_name, method, narrow_rate, sample_size, skew, stddev, min, lower, median, mean, upper, max, ratio)
	VALUES
		(:parts_name, :method, :narrow_rate, :sample_size, :skew, :stddev, :min, :lower, :median, :mean, :upper, :max, :ratio);`
	var prices []float64
	narrowRates := []float64{0.00, 0.05, 0.10, 0.15, 0.20}
	for _, word := range words {
		fmt.Println(word)
		if err := db.Select(&prices, priceExtractQuery, word, method); err != nil {
			panic(err)
		}
		fmt.Println(prices)
		sort.Float64s(prices)
		fmt.Println(prices)
		// sampleSize := len(prices)
		// 通常
		for _, narrowRate := range narrowRates {
			rate := calc(prices, narrowRate)
			rate.PartsName = word
			rate.Method = method
			if _, err := db.NamedQuery(rateCreateQuery, rate); err != nil {
				panic(err)
			}
		}
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

func calc(prices []float64, narrowRate float64) Rate {
	size := len(prices)
	cut := int(math.Floor(float64(size) * narrowRate))
	newPrices := prices[cut : len(prices)-cut]
	newSize := len(newPrices)

	// 歪度
	skew := stat.Skew(newPrices, nil)
	// 標準偏差
	stddev := stat.StdDev(newPrices, nil)
	// 中央値
	median := stat.Quantile(0.5, stat.Empirical, newPrices, nil)
	// 平均値
	mean := stat.Mean(newPrices, nil)
	// 平均値 - 標準偏差
	lower := mean - stddev
	// 平均値 + 標準偏差
	upper := mean + stddev
	// 最小値・最大値
	min, max := minmax(newPrices)
	// ratio
	ratio := float64(count(newPrices, lower, upper)) / float64(newSize)

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
