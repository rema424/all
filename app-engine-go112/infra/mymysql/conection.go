package mymysql

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// GlobalDB ...
var GlobalDB = newDB()

// Queryer ...
type Queryer interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type queryer struct {
	Queryer
}

func newDB() Queryer {
	var (
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		user     = os.Getenv("DB_USER")
		dbName   = os.Getenv("DB_NAME")     // NOTE: dbName may be empty
		password = os.Getenv("DB_PASSWORD") // NOTE: password may be empty
	)

	cfg := mysql.Config{
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

	return &queryer{
		Queryer: dbx,
	}
}

// TxFunc ...
type TxFunc func(context.Context) error

// RunInTransaction ...
func RunInTransaction(ctx context.Context, fn TxFunc) error {
	// if _, ok :=
	return runTransaction(ctx, fn)
}

func runTransaction(ctx context.Context, fn TxFunc) error {
	return nil
}

func (q *queryer) isInTransaction(ctx context.Context) bool {
	if val, ok := ctx.Value("isInTransaction").(bool); ok {
		return val
	}
	return false
}
