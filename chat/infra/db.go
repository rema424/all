package infra

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
		host     = mustGetenv("DB_HOST")
		port     = mustGetenv("DB_PORT")
		user     = mustGetenv("DB_USER")
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

// MyDB ...
// type MyDB struct {
// 	*sqlx.DB
// }

type querier interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// ExDB ...
type ExDB struct {
	querier   // 無名フィールド。クエリを発行する実体を入れる。具体的には sqlx.DB または sqlx.Tx
	ctx       context.Context
	WarnTime  time.Duration
	WarnRows  int
	logHidden bool
}

// NewExDB ...
func NewExDB(ctx context.Context, db *sqlx.DB) *ExDB {
	return &ExDB{
		querier:  db,
		ctx:      ctx,
		WarnTime: 150 * time.Millisecond,
		WarnRows: 1000,
	}
}

// Get ...
func (db *ExDB) Get(dest interface{}, query string, args ...interface{}) error {
	start := time.Now()
	err := db.querier.Get(dest, query, args...)
	db.log(start, err, query, args, countRows(dest))
	return err
}

// TxFunc ...
type TxFunc func(db *ExDB) error

// ExecTx ...
func (db *ExDB) ExecTx(fn TxFunc) error {
	// トランザクション開始済み（querierの実体が*sqlx.Tx）ならそのままクエリ実行
	if _, ok := db.querier.(*sqlx.Tx); ok {
		return fn(db)
	}

	// トランザクション開始前（querierの実体が*sqlx.DB）ならトランザクションを開始
	sqlxDB, ok := db.querier.(*sqlx.DB)
	if !ok {
		return fmt.Errorf("invalid ExDB")
	}

	tx, err := sqlxDB.Begin()
	if err != nil {
		return err
	}

	newExDB := &ExDB{
		querier:   tx,
		ctx:       db.ctx,
		WarnTime:  db.WarnTime,
		WarnRows:  db.WarnRows,
		logHidden: db.logHidden,
	}

	return execTx
}

func (db *ExDB) log(start time.Time, err error, query string, args []interface{}, count int) {

}

func countRows(result interface{}) int {
	return 0
}
