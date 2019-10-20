package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func connectDB() *sqlx.DB {
	var (
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		user     = os.Getenv("DB_USER")
		dbName   = os.Getenv("DB_NAME")     // NOTE: dbName may be empty
		password = os.Getenv("DB_PASSWORD") // NOTE: password may be empty
	)

	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "3306"
	}
	if user == "" {
		user = "workuser"
	}
	if dbName == "" {
		dbName = "mychat"
	}
	if password == "" {
		password = "Passw0rd!"
	}

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

	db, err := sqlx.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatalln(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	log.Println("db connected successfully")

	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(30)
	db.SetConnMaxLifetime(60 * time.Second)

	return db
}

func closeDB(db *sqlx.DB) {
	db.Close()
}

type querier interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// DBx ...
type DBx struct {
	ctx     context.Context
	querier // 無名フィールド。クエリを発行する実体を入れる。具体的には sqlx.DB または sqlx.Tx
}

// NewDBx ...
func NewDBx(ctx context.Context, q querier) *DBx {
	return &DBx{
		ctx:     ctx,
		querier: q,
	}
}

// Get ...
func (db *DBx) Get(dest interface{}, query string, args ...interface{}) error {
	return db.querier.Get(dest, query, args...)
}

// Select ...
func (db *DBx) Select(dest interface{}, query string, args ...interface{}) error {
	return db.querier.Select(dest, query, args...)
}

// Exec ...
func (db *DBx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.querier.Exec(query, args...)
}

// NamedExec ...
func (db *DBx) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return db.querier.NamedExec(query, arg)
}

// Query ...
func (db *DBx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.querier.Query(query, args...)
}

// TxFunc ...
type TxFunc func(db *DBx) error

// BeginTx ...
func (db *DBx) BeginTx(fn TxFunc) error {
	// トランザクション開始済み（querierの実体が*sqlx.Tx）ならそのままクエリ実行
	if _, ok := db.querier.(*sqlx.Tx); ok {
		return fn(db)
	}

	// トランザクション開始前（querierの実体が*sqlx.DB）ならトランザクションを開始
	sqlxDB, ok := db.querier.(*sqlx.DB)
	if !ok {
		return fmt.Errorf("invalid DB")
	}

	tx, err := sqlxDB.Beginx()
	if err != nil {
		return err
	}

	return db.execTxFn(tx, fn)
}

func (db *DBx) execTxFn(tx *sqlx.Tx, fn TxFunc) (err error) {
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	txDBx := &DBx{
		ctx:     db.ctx,
		querier: tx, // *sqlx.DB ではなく *sqlx.Tx を利用してクエリを実行する
	}

	err = fn(txDBx)

	return
}
