package mymysql

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var gDB = newDB()

type db struct {
	db    *sqlx.DB
	txMap map[string]*sqlx.Tx
}

func newDB() *db {
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

	return &db{
		db:    dbx,
		txMap: make(map[string]*sqlx.Tx),
	}
}

// TxFunc ...
type TxFunc func(context.Context) error

func (db *db) RunInTransaction(ctx context.Context, fn TxFunc) error {
	return db.runTransaction(ctx, fn)
}

func (db *db) runTransaction(ctx context.Context, fn TxFunc) error {
	return nil
}

func (db *db) isInTransaction(ctx context.Context) bool {
	if val, ok := ctx.Value("isInTransaction").(bool); ok {
		return val
	}
	return false
}
