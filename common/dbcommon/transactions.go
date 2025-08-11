package dbcommon

import (
	"context"
	"fmt"
	"time"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ent"
)

func WithReadTx[T any](
	ctx context.Context, db common.DatabaseService,
	fn func(tx *ent.Tx, ctx context.Context) (T, error),
) (T, error) {
	return withRetryingTx(ctx, db.ReadTx, fn)
}
func WithWriteTx(
	ctx context.Context, db common.DatabaseService,
	fn func(tx *ent.Tx, ctx context.Context) error,
) error {
	_, stdErr := withRetryingTx(ctx, db.WriteTx, func(tx *ent.Tx, ctx context.Context) (struct{}, error) {
		return struct{}{}, fn(tx, ctx)
	})
	return stdErr
}
func WithReadWriteTx[T any](
	ctx context.Context, db common.DatabaseService,
	fn func(tx *ent.Tx, ctx context.Context) (T, error),
) (T, error) {
	return withRetryingTx(ctx, db.WriteTx, fn)
}

func withRetryingTx[T any](
	ctx context.Context,
	txCallback func(ctx context.Context) (*ent.Tx, error),
	fn func(tx *ent.Tx, ctx context.Context) (T, error),
) (T, error) {
	var returnValue T
	stdErr := common.WithRetries(ctx, func() error {
		return withTx(ctx, txCallback, func(tx *ent.Tx, ctx context.Context) error {
			var stdErr error
			returnValue, stdErr = fn(tx, ctx)
			return stdErr
		})
	})
	return returnValue, stdErr
}
func withTx(
	ctx context.Context,
	txCallback func(ctx context.Context) (*ent.Tx, error),
	fn func(tx *ent.Tx, ctx context.Context) error,
) error {
	tx, stdErr := txCallback(ctx)
	if stdErr != nil {
		return ErrWrapperWithTx.Wrap(
			ErrWrapperStartTx.Wrap(stdErr),
		).ConfigureRetries(-1, 50*time.Millisecond, 2)
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
		return ErrWrapperWithTx.Wrap(ErrWrapperCallback.Wrap(stdErr))
	}
	stdErr = tx.Commit()
	if stdErr != nil {
		return ErrWrapperWithTx.Wrap(
			ErrWrapperCommitTx.Wrap(stdErr),
		).ConfigureRetries(-1, 50*time.Millisecond, 2)
	}
	return nil
}
