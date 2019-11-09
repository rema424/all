package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	my "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type ctxValKey string

const (
	querentKey ctxValKey = "querent-key"
)

var (
	globalDB      *sqlx.DB = newDB()
	globalQuerent *Querent = newQuerent(globalDB)
)

// Querent ...
type Querent struct {
	Executor // *sqlx.DB or *sqlx.Tx
}

// Executor ...
type Executor interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func newDB() *sqlx.DB {
	var (
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		user     = os.Getenv("DB_USER")
		dbName   = os.Getenv("DB_NAME")     // NOTE: dbName may be empty
		password = os.Getenv("DB_PASSWORD") // NOTE: password may be empty
	)

	cfg := my.Config{
		User:                 user,
		Passwd:               password,
		Net:                  "tcp",
		Addr:                 host + ":" + port,
		DBName:               dbName,
		Collation:            "utf8mb4_bin",
		InterpolateParams:    true,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	dbx, err := sqlx.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatalln(err)
	}
	if err := dbx.Ping(); err != nil {
		log.Fatalln(err)
	}

	dbx.SetMaxOpenConns(30)
	dbx.SetMaxIdleConns(30)
	dbx.SetConnMaxLifetime(60 * time.Second)

	return dbx
}

func newQuerent(e Executor) *Querent {
	return &Querent{e}
}

func getQuerent(ctx context.Context) *Querent {
	val, ok := ctx.Value(querentKey).(*Querent)
	if ok && val != nil && val.Executor != nil {
		return val
	}
	return globalQuerent
}

func setQuerent(ctx context.Context, q *Querent) (context.Context, error) {
	if q == nil {
		return ctx, fmt.Errorf("receive invalid querent")
	}
	return context.WithValue(ctx, querentKey, q), nil
}

// Close ...
func Close() {
	if globalDB == nil {
		log.Println("failed to close db - err: db is nil")
		return
	}

	if err := globalDB.Close(); err != nil {
		log.Println("failed to close db - err:", err.Error())
	}
	log.Println("successfully closed db")
}
