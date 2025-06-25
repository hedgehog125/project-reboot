package dbcommon

import (
	"context"
	"fmt"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ent"
)

func WithTx(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) *common.Error) *common.Error {
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
	commErr := fn(tx)
	shouldRecover = false
	if commErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			panic(fmt.Sprintf("error while rolling back transaction:\n%v\nerror that caused rollback:\n%v", commErr, rollbackErr))
		}
		return ErrWrapperCallback.Wrap(commErr).AddCategory(ErrTypeWithTx)
	}
	stdErr = tx.Commit()
	if stdErr != nil {
		return ErrWrapperCommitTx.Wrap(stdErr).AddCategory(ErrTypeWithTx)
	}
	return nil
}
