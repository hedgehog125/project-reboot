// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/hedgehog125/project-reboot/ent/predicate"
	"github.com/hedgehog125/project-reboot/ent/twofactoraction"
)

// TwoFactorActionDelete is the builder for deleting a TwoFactorAction entity.
type TwoFactorActionDelete struct {
	config
	hooks    []Hook
	mutation *TwoFactorActionMutation
}

// Where appends a list predicates to the TwoFactorActionDelete builder.
func (tfad *TwoFactorActionDelete) Where(ps ...predicate.TwoFactorAction) *TwoFactorActionDelete {
	tfad.mutation.Where(ps...)
	return tfad
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (tfad *TwoFactorActionDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, tfad.sqlExec, tfad.mutation, tfad.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (tfad *TwoFactorActionDelete) ExecX(ctx context.Context) int {
	n, err := tfad.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (tfad *TwoFactorActionDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(twofactoraction.Table, sqlgraph.NewFieldSpec(twofactoraction.FieldID, field.TypeUUID))
	if ps := tfad.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, tfad.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	tfad.mutation.done = true
	return affected, err
}

// TwoFactorActionDeleteOne is the builder for deleting a single TwoFactorAction entity.
type TwoFactorActionDeleteOne struct {
	tfad *TwoFactorActionDelete
}

// Where appends a list predicates to the TwoFactorActionDelete builder.
func (tfado *TwoFactorActionDeleteOne) Where(ps ...predicate.TwoFactorAction) *TwoFactorActionDeleteOne {
	tfado.tfad.mutation.Where(ps...)
	return tfado
}

// Exec executes the deletion query.
func (tfado *TwoFactorActionDeleteOne) Exec(ctx context.Context) error {
	n, err := tfado.tfad.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{twofactoraction.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (tfado *TwoFactorActionDeleteOne) ExecX(ctx context.Context) {
	if err := tfado.Exec(ctx); err != nil {
		panic(err)
	}
}
