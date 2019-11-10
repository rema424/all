package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// Get ...
func (a *Accessor) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	if err := a.validate(); err != nil {
		return err
	}
	return a.build(ctx).querent.Get(dest, query, args...)
}

// Select ...
func (a *Accessor) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	if err := a.validate(); err != nil {
		return err
	}
	return a.build(ctx).querent.Select(dest, query, args...)
}

// Exec ...
func (a *Accessor) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if err := a.validate(); err != nil {
		return nil, err
	}
	return a.build(ctx).querent.Exec(query, args...)
}

// NamedExec ...
func (a *Accessor) NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	if err := a.validate(); err != nil {
		return nil, err
	}
	return a.build(ctx).querent.NamedExec(query, arg)
}

// Query ...
func (a *Accessor) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if err := a.validate(); err != nil {
		return nil, err
	}
	return a.build(ctx).querent.Query(query, args...)
}

// TxFunc ...
type TxFunc func(context.Context) (interface{}, error)

// RunInTx ...
func (a *Accessor) RunInTx(ctx context.Context, txFn TxFunc) (interface{}, error) {
	fmt.Println("start infra mysql RunInTx")
	defer fmt.Println("finish infra mysql RunInTx")

	if err := a.validate(); err != nil {
		return nil, err
	}

	a = a.build(ctx)

	tx, ok := a.querent.(*sqlx.Tx)
	if !ok || tx == nil {
		var err error

		// トランザクションを取得
		dbx, ok := a.querent.(*sqlx.DB)
		if !ok || dbx == nil {
			return nil, fmt.Errorf("failed to begin transaction - invalid dbx")
		}

		tx, err = dbx.Beginx()
		if err != nil {
			return nil, err
		}

		// コンテキストにトランザクションを格納
		ctx, err = set(ctx, newAccessor(tx))
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
