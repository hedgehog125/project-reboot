// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/hedgehog125/project-reboot/ent/loginattempt"
	"github.com/hedgehog125/project-reboot/ent/predicate"
)

// LoginAttemptUpdate is the builder for updating LoginAttempt entities.
type LoginAttemptUpdate struct {
	config
	hooks    []Hook
	mutation *LoginAttemptMutation
}

// Where appends a list predicates to the LoginAttemptUpdate builder.
func (lau *LoginAttemptUpdate) Where(ps ...predicate.LoginAttempt) *LoginAttemptUpdate {
	lau.mutation.Where(ps...)
	return lau
}

// SetTime sets the "time" field.
func (lau *LoginAttemptUpdate) SetTime(t time.Time) *LoginAttemptUpdate {
	lau.mutation.SetTime(t)
	return lau
}

// SetNillableTime sets the "time" field if the given value is not nil.
func (lau *LoginAttemptUpdate) SetNillableTime(t *time.Time) *LoginAttemptUpdate {
	if t != nil {
		lau.SetTime(*t)
	}
	return lau
}

// SetCode sets the "code" field.
func (lau *LoginAttemptUpdate) SetCode(s string) *LoginAttemptUpdate {
	lau.mutation.SetCode(s)
	return lau
}

// SetNillableCode sets the "code" field if the given value is not nil.
func (lau *LoginAttemptUpdate) SetNillableCode(s *string) *LoginAttemptUpdate {
	if s != nil {
		lau.SetCode(*s)
	}
	return lau
}

// Mutation returns the LoginAttemptMutation object of the builder.
func (lau *LoginAttemptUpdate) Mutation() *LoginAttemptMutation {
	return lau.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (lau *LoginAttemptUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, lau.sqlSave, lau.mutation, lau.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (lau *LoginAttemptUpdate) SaveX(ctx context.Context) int {
	affected, err := lau.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (lau *LoginAttemptUpdate) Exec(ctx context.Context) error {
	_, err := lau.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lau *LoginAttemptUpdate) ExecX(ctx context.Context) {
	if err := lau.Exec(ctx); err != nil {
		panic(err)
	}
}

func (lau *LoginAttemptUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(loginattempt.Table, loginattempt.Columns, sqlgraph.NewFieldSpec(loginattempt.FieldID, field.TypeInt))
	if ps := lau.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := lau.mutation.Time(); ok {
		_spec.SetField(loginattempt.FieldTime, field.TypeTime, value)
	}
	if value, ok := lau.mutation.Code(); ok {
		_spec.SetField(loginattempt.FieldCode, field.TypeString, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, lau.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{loginattempt.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	lau.mutation.done = true
	return n, nil
}

// LoginAttemptUpdateOne is the builder for updating a single LoginAttempt entity.
type LoginAttemptUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *LoginAttemptMutation
}

// SetTime sets the "time" field.
func (lauo *LoginAttemptUpdateOne) SetTime(t time.Time) *LoginAttemptUpdateOne {
	lauo.mutation.SetTime(t)
	return lauo
}

// SetNillableTime sets the "time" field if the given value is not nil.
func (lauo *LoginAttemptUpdateOne) SetNillableTime(t *time.Time) *LoginAttemptUpdateOne {
	if t != nil {
		lauo.SetTime(*t)
	}
	return lauo
}

// SetCode sets the "code" field.
func (lauo *LoginAttemptUpdateOne) SetCode(s string) *LoginAttemptUpdateOne {
	lauo.mutation.SetCode(s)
	return lauo
}

// SetNillableCode sets the "code" field if the given value is not nil.
func (lauo *LoginAttemptUpdateOne) SetNillableCode(s *string) *LoginAttemptUpdateOne {
	if s != nil {
		lauo.SetCode(*s)
	}
	return lauo
}

// Mutation returns the LoginAttemptMutation object of the builder.
func (lauo *LoginAttemptUpdateOne) Mutation() *LoginAttemptMutation {
	return lauo.mutation
}

// Where appends a list predicates to the LoginAttemptUpdate builder.
func (lauo *LoginAttemptUpdateOne) Where(ps ...predicate.LoginAttempt) *LoginAttemptUpdateOne {
	lauo.mutation.Where(ps...)
	return lauo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (lauo *LoginAttemptUpdateOne) Select(field string, fields ...string) *LoginAttemptUpdateOne {
	lauo.fields = append([]string{field}, fields...)
	return lauo
}

// Save executes the query and returns the updated LoginAttempt entity.
func (lauo *LoginAttemptUpdateOne) Save(ctx context.Context) (*LoginAttempt, error) {
	return withHooks(ctx, lauo.sqlSave, lauo.mutation, lauo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (lauo *LoginAttemptUpdateOne) SaveX(ctx context.Context) *LoginAttempt {
	node, err := lauo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (lauo *LoginAttemptUpdateOne) Exec(ctx context.Context) error {
	_, err := lauo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lauo *LoginAttemptUpdateOne) ExecX(ctx context.Context) {
	if err := lauo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (lauo *LoginAttemptUpdateOne) sqlSave(ctx context.Context) (_node *LoginAttempt, err error) {
	_spec := sqlgraph.NewUpdateSpec(loginattempt.Table, loginattempt.Columns, sqlgraph.NewFieldSpec(loginattempt.FieldID, field.TypeInt))
	id, ok := lauo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "LoginAttempt.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := lauo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, loginattempt.FieldID)
		for _, f := range fields {
			if !loginattempt.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != loginattempt.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := lauo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := lauo.mutation.Time(); ok {
		_spec.SetField(loginattempt.FieldTime, field.TypeTime, value)
	}
	if value, ok := lauo.mutation.Code(); ok {
		_spec.SetField(loginattempt.FieldCode, field.TypeString, value)
	}
	_node = &LoginAttempt{config: lauo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, lauo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{loginattempt.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	lauo.mutation.done = true
	return _node, nil
}
