package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	redisPool = newRedisPool()
	e         = createMux()
)

func createMux() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())

	return e
}

func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err)
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func init() {
	// curl -X POST -d 'userId=user:1' -d 'articleTitle=タイトル1' -d localhost:3333/articles
	e.POST("/articles", postArticle)
}

func main() {
	port := flag.String("port", "3333", "アプリケーションのアドレス")
	flag.Parse()

	// http.Handle("/", e)
	// log.Printf("server started at %s\n", *port)
	// log.Fatal(http.ListenAndServe("localhost:"+*port, nil))
	e.Logger.Fatal(e.Start(":" + *port))
}

func postArticle(c echo.Context) error {
	in := struct {
		UserID string `form:"userId"` // user:1111
		Title  string `form:"title"`  // タイトル
	}{}

	if err := c.Bind(&in); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	conn := redisPool.Get()
	if conn == nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	defer conn.Close()

	// 記事ID文字列
	articleID, err := redis.Int(conn.Do("incr", "articleID:"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	fmt.Println(articleID)

	// 投票セット
	voted := fmt.Sprintf("voted:%d", articleID)
	conn.Do("sadd", voted, in.UserID)
	conn.Do("expire", voted, oneWeekSeconds)

	// 記事ハッシュ
	now := time.Now().Unix()
	article := fmt.Sprintf("article:%d", articleID)
	conn.Do("HMSET", article,
		"title", in.Title,
		"poster", in.UserID,
		"time", now,
		"votes", 1,
	)

	// スコアソートセット
	conn.Do("zadd", now+voteScore, article)

	// 投稿時間ソートセット
	conn.Do("zadd", now, article)

	return c.JSON(200, in)
}

const (
	oneWeekSeconds = 7 * 86400
	voteScore      = 432
)

func articleVote(c echo.Context) error {
	conn := redisPool.Get()
	if conn == nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	defer conn.Close()

	articleID := c.FormValue("article-id")

	r, err := redis.Int(conn.Do("ZSCORE", "time:", "article:"+articleID))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	fmt.Println(r)

	cutoff := time.Now().Unix() - oneWeekSeconds
	if int64(r) < cutoff {
		return c.String(http.StatusOK, "投票期限が過ぎています")
	}

	return nil
}
