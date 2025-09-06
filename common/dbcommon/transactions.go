package dbcommon

import (
	"context"
	"time"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ent"
)

// TODO: what happens if expired contexts are passed?

func WithReadTx[T any](
	ctx context.Context, db common.DatabaseService,
	fn func(tx *ent.Tx, ctx context.Context) (T, error),
) (T, error) {
	return withRetryingTx(ctx, db, db.ReadTx, fn)
}
func WithWriteTx(
	ctx context.Context, db common.DatabaseService,
	fn func(tx *ent.Tx, ctx context.Context) error,
) error {
	_, stdErr := withRetryingTx(ctx, db, db.WriteTx, func(tx *ent.Tx, ctx context.Context) (struct{}, error) {
		return struct{}{}, fn(tx, ctx)
	})
	return stdErr
}
func WithReadWriteTx[T any](
	ctx context.Context, db common.DatabaseService,
	fn func(tx *ent.Tx, ctx context.Context) (T, error),
) (T, error) {
	return withRetryingTx(ctx, db, db.WriteTx, fn)
}

func withRetryingTx[T any](
	ctx context.Context, db common.DatabaseService,
	txCallback func(ctx context.Context) (*ent.Tx, error),
	fn func(tx *ent.Tx, ctx context.Context) (T, error),
) (T, error) {
	var returnValue T
	stdErr := common.WithRetries(ctx, common.GetLogger(ctx, db), func() error {
		return withTx(ctx, db, txCallback, func(tx *ent.Tx, ctx context.Context) error {
			var stdErr error
			returnValue, stdErr = fn(tx, ctx)
			return stdErr
		})
	})
	return returnValue, stdErr
}
func withTx(
	ctx context.Context, db common.DatabaseService,
	txCallback func(ctx context.Context) (*ent.Tx, error),
	fn func(tx *ent.Tx, ctx context.Context) error,
) error {
	if ent.TxFromContext(ctx) != nil {
		return ErrWrapperWithTx.Wrap(ErrUnexpectedTransaction)
	}
	tx, stdErr := txCallback(ctx)
	if stdErr != nil {
		return ErrWrapperWithTx.Wrap(
			ErrWrapperStartTx.Wrap(stdErr),
		).ConfigureRetries(-1, 5*time.Millisecond, 1.5)
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
			// TODO: handle "transaction already committed or rolled back" errors?
			common.GetLogger(ctx, db).Error("error rolling back transaction", "error", rollbackErr, "originalError", stdErr)
		}
		return ErrWrapperWithTx.Wrap(ErrWrapperCallback.Wrap(stdErr))
	}
	stdErr = tx.Commit()
	if stdErr != nil {
		return ErrWrapperWithTx.Wrap(
			ErrWrapperCommitTx.Wrap(stdErr),
		).ConfigureRetries(-1, 5*time.Millisecond, 1.5)
	}
	return nil
}
