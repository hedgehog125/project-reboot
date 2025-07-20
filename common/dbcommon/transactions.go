package dbcommon

import (
	"context"
	"fmt"

	"github.com/hedgehog125/project-reboot/ent"
)

func WithTx(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) error) error {
	tx, stdErr := client.Tx(ctx)
	if stdErr != nil {
		return ErrWrapperStartTx.Wrap(stdErr).AddCategory(ErrTypeWithTx)
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
	stdErr = fn(tx)
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
		return ErrWrapperCommitTx.Wrap(stdErr).AddCategory(ErrTypeWithTx)
	}
	return nil
}
