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

// MyDB ...
type MyDB struct {
	querier   // 無名フィールド。クエリを発行する実体を入れる。具体的には sqlx.DB または sqlx.Tx
	ctx       context.Context
	WarnTime  time.Duration
	WarnRows  int
	logHidden bool
}

// NewMyDB ...
func NewMyDB(ctx context.Context, q querier) *MyDB {
	return &MyDB{
		querier:  q,
		ctx:      ctx,
		WarnTime: 150 * time.Millisecond,
		WarnRows: 1000,
	}
}

// Get ...
func (db *MyDB) Get(dest interface{}, query string, args ...interface{}) error {
	start := time.Now()
	err := db.querier.Get(dest, query, args...)
	db.log(start, err, query, args, countRows(dest))
	return err
}

// TxFunc ...
type TxFunc func(db *MyDB) error

// BeginTx ...
func (db *MyDB) BeginTx(fn TxFunc) error {
	// トランザクション開始済み（querierの実体が*sqlx.Tx）ならそのままクエリ実行
	if _, ok := db.querier.(*sqlx.Tx); ok {
		return fn(db)
	}

	// トランザクション開始前（querierの実体が*sqlx.DB）ならトランザクションを開始
	sqlxDB, ok := db.querier.(*sqlx.DB)
	if !ok {
		return fmt.Errorf("invalid MyDB")
	}

	tx, err := sqlxDB.Beginx()
	if err != nil {
		return err
	}

	return db.execTxFn(tx, fn)
}

func (db *MyDB) execTxFn(tx *sqlx.Tx, fn TxFunc) (err error) {
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

	txMyDB := &MyDB{
		querier:  tx, // *sqlx.DB ではなく *sqlx.Tx を利用してクエリを実行する
		ctx:      db.ctx,
		WarnTime: db.WarnTime,
		WarnRows: db.WarnRows,
	}

	err = fn(txMyDB)

	return
}

func (db *MyDB) execTxFnLess(tx *sqlx.Tx, fn TxFunc) (err error) {
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	txMyDB := &MyDB{
		querier: tx, // *sqlx.DB ではなく *sqlx.Tx を利用してクエリを実行する
		ctx:     db.ctx,
	}

	err = fn(txMyDB)

	return
}

func (db *MyDB) log(start time.Time, err error, query string, args []interface{}, count int) {

}

func countRows(result interface{}) int {
	return 0
}
