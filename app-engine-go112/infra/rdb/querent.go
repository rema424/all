package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	my "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Accessor ...
type Accessor struct {
	querent querent // *sqlx.DB or *sqlx.Tx
}

type querent interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// Config ...
type Config struct {
	User                 string
	Passwd               string
	Host                 string
	Port                 string
	Net                  string
	Addr                 string
	DBName               string
	Collation            string
	InterpolateParams    bool
	AllowNatevePasswords bool
	ParseTime            bool
	MaxOpenConns         int
	MaxIdleConns         int
	ConnMaxLifetime      time.Duration
}

func (c Config) build() Config {
	// var (
	// 	host     = os.Getenv("DB_HOST")
	// 	port     = os.Getenv("DB_PORT")
	// 	user     = os.Getenv("DB_USER")
	// 	dbName   = os.Getenv("DB_NAME")     // NOTE: dbName may be empty
	// 	password = os.Getenv("DB_PASSWORD") // NOTE: password may be empty
	// )
	if c.User == "" {
		c.User = "root"
	}
	if c.Net == "" {
		c.Net = "tcp"
	}
	if c.Host == "" {
		c.Host = "127.0.0.1"
	}
	if c.Port == "" {
		c.Port = "3306"
	}
	if c.Addr == "" {
		c.Addr = c.Host + ":" + c.Port
	}
	if c.Collation == "" {
		c.Collation = "utf8mb4_bin"
	}
	if c.MaxOpenConns < 0 {
		c.MaxOpenConns = 30
	}
	if c.MaxIdleConns < 0 {
		c.MaxIdleConns = 30
	}
	if c.ConnMaxLifetime < 0 {
		c.ConnMaxLifetime = 60 * time.Second
	}
	return c
}

func newDB(c Config) *sqlx.DB {
	c = c.build()

	mycfg := my.Config{
		User:                 c.User,
		Passwd:               c.Passwd,
		Net:                  c.Net,
		Addr:                 c.Addr,
		DBName:               c.DBName,
		Collation:            c.Collation,
		InterpolateParams:    c.InterpolateParams,
		AllowNativePasswords: c.AllowNatevePasswords,
		ParseTime:            c.ParseTime,
	}

	dbx, err := sqlx.Open("mysql", mycfg.FormatDSN())
	if err != nil {
		log.Fatalln(err)
	}
	if err := dbx.Ping(); err != nil {
		log.Fatalln(err)
	}
	log.Println("db opened successfully")

	dbx.SetMaxOpenConns(c.MaxOpenConns)
	dbx.SetMaxIdleConns(c.MaxIdleConns)
	dbx.SetConnMaxLifetime(c.ConnMaxLifetime)

	return dbx
}

func newAccessor(q querent) *Accessor {
	return &Accessor{q}
}

// Open ...
func Open(c Config) *Accessor {
	return newAccessor(newDB(c))
}

// Close ...
func (a *Accessor) Close() error {
	if err := a.validate(); err != nil {
		return fmt.Errorf("failed to close db - %s", err.Error())
	}

	dbx, ok := a.querent.(*sqlx.DB)
	if !ok || dbx == nil {
		return fmt.Errorf("failed to close db - invalid dbx")
	}

	if err := dbx.Close(); err != nil {
		return fmt.Errorf("failed to close db - %s", err.Error())
	}

	log.Println("db closed successfully")
	return nil
}

type ctxValKey string

const (
	accessorKey ctxValKey = "accessor-key"
)

func set(ctx context.Context, a *Accessor) (context.Context, error) {
	if err := a.validate(); err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, accessorKey, a), nil
}

func (a *Accessor) build(ctx context.Context) *Accessor {
	if val, ok := ctx.Value(accessorKey).(*Accessor); ok {
		return val
	}
	return a
}

func (a *Accessor) validate() error {
	if a == nil || a.querent == nil {
		return fmt.Errorf("invalid mysql accessor")
	}
	return nil
}
