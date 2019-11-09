package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// Get ...
func Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return getQuerent(ctx).Executor.Get(dest, query, args...)
}

// Select ...
func Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return getQuerent(ctx).Executor.Select(dest, query, args...)
}

// Exec ...
func Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return getQuerent(ctx).Executor.Exec(query, args...)
}

// NamedExec ...
func NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return getQuerent(ctx).Executor.NamedExec(query, arg)
}

// Query ...
func Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return getQuerent(ctx).Executor.Query(query, args...)
}

// TxFunc ...
type TxFunc func(context.Context) (interface{}, error)

// RunInTx ...
func RunInTx(ctx context.Context, txFn TxFunc) (interface{}, error) {
	tx, ok := getQuerent(ctx).Executor.(*sqlx.Tx)
	if !ok || tx == nil {
		var err error

		// トランザクションを取得
		tx, err = globalDB.Beginx()
		if err != nil {
			return nil, err
		}

		// コンテキストにトランザクションを格納
		ctx, err = setQuerent(ctx, newQuerent(tx))
		if err != nil {
			return nil, err
		}
	}

	// トランザクション内（ctx）で一連の SQL を実行
	v, err := txFn(ctx)

	// ロールバック・コミットの実行
	if pnc := recover(); pnc != nil {
		// panic が発生したらロールバックを実行
		if rlbkErr := tx.Rollback(); rlbkErr != nil {
			// ロールバックエラーが発生したら通知し、特別な措置を講じて回復する
			log.Printf("failed to rollback - err: %s\n", rlbkErr.Error())
		}
		// panic の内容を返却
		if pncErr, ok := pnc.(error); ok {
			return nil, pncErr
		}
		return nil, fmt.Errorf("%v", pnc)
	} else if err != nil {
		// SQL 実行エラーが発生したらロールバックを実行
		if rlbkErr := tx.Rollback(); rlbkErr != nil {
			// ロールバックエラーが発生したら通知し、特別な措置を講じて回復する
			log.Printf("failed to rollback - err: %s\n", rlbkErr.Error())
		}
		// SQL 実行エラーを返却
		return nil, err
	} else {
		// panic も SQL 実行エラーも発生しない場合はコミットを実行
		if cmtErr := tx.Commit(); err != nil {
			// コミットエラーを返却
			return nil, cmtErr
		}
	}

	// コミットまで成功したら値を返却
	return v, nil
}
