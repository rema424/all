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
	// curl -X POST -d 'userId=1' -d 'title=タイトル1' -d localhost:3333/articles
	e.POST("/articles", postArticle)

	// curl -X POST -d 'userId=1' localhost:3333/articles/9/votes
	e.POST("/articles/:articleID/votes", articleVote)
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
	conn.Do("sadd", voted, "user:"+in.UserID)
	conn.Do("expire", voted, oneWeekSeconds)

	// 記事ハッシュ
	now := time.Now().Unix()
	article := fmt.Sprintf("article:%d", articleID)
	conn.Do("HMSET", article,
		"title", in.Title,
		"poster", in.UserID,
		"created", now,
		"votes", 1,
	)

	// スコアソートセット
	conn.Do("zadd", "score:", now+voteScore, article)

	// 投稿時間ソートセット
	conn.Do("zadd", "created:", now, article)

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

	articleID := c.Param("articleID")
	userID := c.FormValue("userId")

	if articleID == "" || userID == "" {
		return c.String(http.StatusInternalServerError, "パラメータが足りません。")
	}

	created, err := redis.Int(conn.Do("ZSCORE", "created:", "article:"+articleID))
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "なんかエラー")
	}
	fmt.Println(created)

	cutoff := time.Now().Unix() - oneWeekSeconds
	if int64(created) < cutoff {
		return c.String(400, "投票期限が過ぎています")
	}

	if _, err := conn.Do("sadd", "voted:"+articleID, "user:"+userID); err != nil {
		return c.String(400, "投票失敗")
	}
	conn.Do("zincrby", "score:", voteScore, "article:"+articleID)
	conn.Do("hincrby", "article:"+articleID, "votes", 1)

	return c.String(200, "成功です。")
}
