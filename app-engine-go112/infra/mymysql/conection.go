package mymysql

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

var gDB = NewDB()

// Queryer ...
type Queryer interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type queryer struct {
	Queryer
}

// // DB ...
// type DB struct {
// 	*sqlx.DB
// 	txMap map[string]*sqlx.Tx
// }

// NewDB ...
func NewDB() Queryer {
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

	return &queryer{dbx}
}

func newTx(tx *sqlx.Tx) Queryer {
	return &queryer{tx}
}

// TxFunc ...
type TxFunc func(ctx context.Context, q Queryer) error

// RunInTransaction ...
func (q *queryer) RunInTransaction(ctx context.Context, fn TxFunc) error {
	if val, ok := ctx.Value("transaction").(*queryer); ok {
		if tx, ok := val.Queryer.(*sqlx.Tx); ok {
			return fn(ctx, tx)
		}
	}

	var tx *sqlx.Tx
	if val, ok := q.Queryer.(*sqlx.Tx); ok {
		tx = val
	} else if val, ok := q.Queryer.(*sqlx.DB); ok {
		var err error
		tx, err = val.Beginx()
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("invalid DB")
	}

	return runTransaction(ctx, tx, fn)
}

func runTransaction(ctx context.Context, tx *sqlx.Tx, fn TxFunc) (err error) {
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = fn(newTx(tx))
	return
}

func (q *queryer) isInTransaction(ctx context.Context) bool {
	if val, ok := ctx.Value("isInTransaction").(bool); ok {
		return val
	}
	return false
}
