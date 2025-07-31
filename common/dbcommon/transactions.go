package dbcommon

import (
	"context"
	"fmt"
	"time"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ent"
)

func WithReadTx(ctx context.Context, db common.DatabaseService, fn func(tx *ent.Tx, ctx context.Context) error) error {
	return withRetryingTx(ctx, db.ReadTx, fn)
}
func WithWriteTx(ctx context.Context, db common.DatabaseService, fn func(tx *ent.Tx, ctx context.Context) error) error {
	return withRetryingTx(ctx, db.WriteTx, fn)
}

func withRetryingTx(
	ctx context.Context,
	txCallback func(ctx context.Context) (*ent.Tx, error),
	fn func(tx *ent.Tx, ctx context.Context) error,
) error {
	return common.WithRetries(ctx, func() error {
		return withTx(ctx, txCallback, fn)
	})
}
func withTx(
	ctx context.Context,
	txCallback func(ctx context.Context) (*ent.Tx, error),
	fn func(tx *ent.Tx, ctx context.Context) error,
) error {
	tx, stdErr := txCallback(ctx)
	if stdErr != nil {
		return ErrWrapperStartTx.Wrap(stdErr).
			AddCategory(ErrTypeWithTx).
			ConfigureRetries(-1, 50*time.Millisecond, 2)
	}
	shouldRecover := true
	defer func() {
		if !shouldRecover {
			return
		}
		pErr := recover()
		if pErr != nil {
			tx.Rollback()
			panic(pErr)
		}
	}()
	stdErr = fn(tx, ent.NewTxContext(ctx, tx))
	shouldRecover = false
	if stdErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			// TODO: log instead
			panic(fmt.Sprintf("error while rolling back transaction:\n%v\nerror that caused rollback:\n%v", stdErr, rollbackErr))
		}
		return ErrWrapperCallback.Wrap(stdErr).AddCategory(ErrTypeWithTx)
	}
	stdErr = tx.Commit()
	if stdErr != nil {
		return ErrWrapperCommitTx.Wrap(stdErr).
			AddCategory(ErrTypeWithTx).
			ConfigureRetries(-1, 50*time.Millisecond, 2)

	}
	return nil
}
