package component

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Global variables
var (
	GDBx *sqlx.DB
	gDB  *DB
)

type querier interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// DB ...
type DB struct {
	querier // 無名フィールド。クエリを発行する実体を入れる。具体的には sqlx.DB または sqlx.Tx
	// ctx       context.Context
	// WarnTime  time.Duration
	// WarnRows  int
	// logHidden bool
}

func newDB(q querier) *DB {
	return &DB{q}
}

// OpenDB ...
func OpenDB() {
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

	GDBx, gDB = dbx, newDB(dbx)

	log.Println("db connected successfully")
}

// CloseDB ...
func CloseDB() {
	if GDBx == nil {
		log.Fatalln("failed to close dbx - err: not found dbx instance")
	}

	if err := GDBx.Close(); err != nil {
		log.Fatalln("failed to close dbx - err:", err.Error())
	}

	gDB = nil
	log.Println("db closed successfully")
}

// GetDB ...
func GetDB() *DB {
	if gDB == nil {
		log.Fatalln("failed to get db - err: not found db instance")
	}
	return gDB
}

// Get ...
func (db *DB) Get(dest interface{}, query string, args ...interface{}) error {
	start := time.Now()
	err := db.querier.Get(dest, query, args...)
	db.log(start, err, query, args, countRows(dest))
	return err
}

// NamedExec ...
func (db *DB) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return db.querier.NamedExec(query, arg)
}

// TxFunc ...
type TxFunc func(db *DB) error

// RunInTx ...
func (db *DB) RunInTx(fn TxFunc) error {
	// トランザクション開始済み（querierの実体が*sqlx.Tx）ならそのままクエリ実行
	if _, ok := db.querier.(*sqlx.Tx); ok {
		return fn(db)
	}

	// トランザクション開始前（querierの実体が*sqlx.DB）ならトランザクションを開始
	dbx, ok := db.querier.(*sqlx.DB)
	if !ok {
		return fmt.Errorf("invalid DB")
	}

	tx, err := dbx.Beginx()
	if err != nil {
		return err
	}

	return db.execTxFn(tx, fn)
}

func (db *DB) execTxFn(tx *sqlx.Tx, fn TxFunc) (err error) {
	defer func() {
		// recover() で panic を捕捉して制御を分ける
		if r := recover(); r != nil {
			// panic が発生したらロールバックを実行
			if rberr := tx.Rollback(); rberr != nil {
				// ロールバックエラーがある場合はログ出力
				if rberr == sql.ErrTxDone {
					log.Printf("debug: rollback error %s", rberr)
				} else {
					log.Printf("error: rollback error %s", rberr)
				}
			}

			// panic の内容をエラーとして返却
			if rerr, ok := r.(error); ok {
				err = rerr
			} else {
				err = fmt.Errorf("%v", r)
			}
		} else if err != nil {
			// panic が発生せず、error が発生した場合はロールバックを実行
			if rberr := tx.Rollback(); rberr != nil {
				// ロールバックエラーがある場合はログ出力
				if rberr == sql.ErrTxDone {
					log.Printf("debug: rollback error %s", rberr)
				} else {
					log.Printf("error: rollback error %s", rberr)
				}
			}
		} else {
			// panic も error も発生していない場合はコミットを実行
			if cmerr := tx.Commit(); cmerr != nil {
				if cmerr != sql.ErrTxDone {
					err = cmerr
				}
			}
		}
	}()

	err = fn(newDB(tx))

	return
}

func (db *DB) execTxFnLess(tx *sqlx.Tx, fn TxFunc) (err error) {
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = fn(newDB(tx))

	return
}

func (db *DB) log(start time.Time, err error, query string, args []interface{}, count int) {

}

func countRows(result interface{}) int {
	return 0
}
