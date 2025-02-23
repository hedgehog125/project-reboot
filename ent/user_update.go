// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/hedgehog125/project-reboot/ent/loginattempt"
	"github.com/hedgehog125/project-reboot/ent/predicate"
	"github.com/hedgehog125/project-reboot/ent/user"
)

// UserUpdate is the builder for updating User entities.
type UserUpdate struct {
	config
	hooks    []Hook
	mutation *UserMutation
}

// Where appends a list predicates to the UserUpdate builder.
func (uu *UserUpdate) Where(ps ...predicate.User) *UserUpdate {
	uu.mutation.Where(ps...)
	return uu
}

// SetUsername sets the "username" field.
func (uu *UserUpdate) SetUsername(s string) *UserUpdate {
	uu.mutation.SetUsername(s)
	return uu
}

// SetNillableUsername sets the "username" field if the given value is not nil.
func (uu *UserUpdate) SetNillableUsername(s *string) *UserUpdate {
	if s != nil {
		uu.SetUsername(*s)
	}
	return uu
}

// SetAlertDiscordId sets the "alertDiscordId" field.
func (uu *UserUpdate) SetAlertDiscordId(s string) *UserUpdate {
	uu.mutation.SetAlertDiscordId(s)
	return uu
}

// SetNillableAlertDiscordId sets the "alertDiscordId" field if the given value is not nil.
func (uu *UserUpdate) SetNillableAlertDiscordId(s *string) *UserUpdate {
	if s != nil {
		uu.SetAlertDiscordId(*s)
	}
	return uu
}

// SetAlertEmail sets the "alertEmail" field.
func (uu *UserUpdate) SetAlertEmail(s string) *UserUpdate {
	uu.mutation.SetAlertEmail(s)
	return uu
}

// SetNillableAlertEmail sets the "alertEmail" field if the given value is not nil.
func (uu *UserUpdate) SetNillableAlertEmail(s *string) *UserUpdate {
	if s != nil {
		uu.SetAlertEmail(*s)
	}
	return uu
}

// SetContent sets the "content" field.
func (uu *UserUpdate) SetContent(b []byte) *UserUpdate {
	uu.mutation.SetContent(b)
	return uu
}

// SetFileName sets the "fileName" field.
func (uu *UserUpdate) SetFileName(s string) *UserUpdate {
	uu.mutation.SetFileName(s)
	return uu
}

// SetNillableFileName sets the "fileName" field if the given value is not nil.
func (uu *UserUpdate) SetNillableFileName(s *string) *UserUpdate {
	if s != nil {
		uu.SetFileName(*s)
	}
	return uu
}

// SetMime sets the "mime" field.
func (uu *UserUpdate) SetMime(s string) *UserUpdate {
	uu.mutation.SetMime(s)
	return uu
}

// SetNillableMime sets the "mime" field if the given value is not nil.
func (uu *UserUpdate) SetNillableMime(s *string) *UserUpdate {
	if s != nil {
		uu.SetMime(*s)
	}
	return uu
}

// SetNonce sets the "nonce" field.
func (uu *UserUpdate) SetNonce(b []byte) *UserUpdate {
	uu.mutation.SetNonce(b)
	return uu
}

// SetKeySalt sets the "keySalt" field.
func (uu *UserUpdate) SetKeySalt(b []byte) *UserUpdate {
	uu.mutation.SetKeySalt(b)
	return uu
}

// SetPasswordHash sets the "passwordHash" field.
func (uu *UserUpdate) SetPasswordHash(b []byte) *UserUpdate {
	uu.mutation.SetPasswordHash(b)
	return uu
}

// SetPasswordSalt sets the "passwordSalt" field.
func (uu *UserUpdate) SetPasswordSalt(b []byte) *UserUpdate {
	uu.mutation.SetPasswordSalt(b)
	return uu
}

// SetHashTime sets the "hashTime" field.
func (uu *UserUpdate) SetHashTime(u uint32) *UserUpdate {
	uu.mutation.ResetHashTime()
	uu.mutation.SetHashTime(u)
	return uu
}

// SetNillableHashTime sets the "hashTime" field if the given value is not nil.
func (uu *UserUpdate) SetNillableHashTime(u *uint32) *UserUpdate {
	if u != nil {
		uu.SetHashTime(*u)
	}
	return uu
}

// AddHashTime adds u to the "hashTime" field.
func (uu *UserUpdate) AddHashTime(u int32) *UserUpdate {
	uu.mutation.AddHashTime(u)
	return uu
}

// SetHashMemory sets the "hashMemory" field.
func (uu *UserUpdate) SetHashMemory(u uint32) *UserUpdate {
	uu.mutation.ResetHashMemory()
	uu.mutation.SetHashMemory(u)
	return uu
}

// SetNillableHashMemory sets the "hashMemory" field if the given value is not nil.
func (uu *UserUpdate) SetNillableHashMemory(u *uint32) *UserUpdate {
	if u != nil {
		uu.SetHashMemory(*u)
	}
	return uu
}

// AddHashMemory adds u to the "hashMemory" field.
func (uu *UserUpdate) AddHashMemory(u int32) *UserUpdate {
	uu.mutation.AddHashMemory(u)
	return uu
}

// SetHashKeyLen sets the "hashKeyLen" field.
func (uu *UserUpdate) SetHashKeyLen(u uint32) *UserUpdate {
	uu.mutation.ResetHashKeyLen()
	uu.mutation.SetHashKeyLen(u)
	return uu
}

// SetNillableHashKeyLen sets the "hashKeyLen" field if the given value is not nil.
func (uu *UserUpdate) SetNillableHashKeyLen(u *uint32) *UserUpdate {
	if u != nil {
		uu.SetHashKeyLen(*u)
	}
	return uu
}

// AddHashKeyLen adds u to the "hashKeyLen" field.
func (uu *UserUpdate) AddHashKeyLen(u int32) *UserUpdate {
	uu.mutation.AddHashKeyLen(u)
	return uu
}

// AddLoginAttemptIDs adds the "loginAttempts" edge to the LoginAttempt entity by IDs.
func (uu *UserUpdate) AddLoginAttemptIDs(ids ...int) *UserUpdate {
	uu.mutation.AddLoginAttemptIDs(ids...)
	return uu
}

// AddLoginAttempts adds the "loginAttempts" edges to the LoginAttempt entity.
func (uu *UserUpdate) AddLoginAttempts(l ...*LoginAttempt) *UserUpdate {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return uu.AddLoginAttemptIDs(ids...)
}

// Mutation returns the UserMutation object of the builder.
func (uu *UserUpdate) Mutation() *UserMutation {
	return uu.mutation
}

// ClearLoginAttempts clears all "loginAttempts" edges to the LoginAttempt entity.
func (uu *UserUpdate) ClearLoginAttempts() *UserUpdate {
	uu.mutation.ClearLoginAttempts()
	return uu
}

// RemoveLoginAttemptIDs removes the "loginAttempts" edge to LoginAttempt entities by IDs.
func (uu *UserUpdate) RemoveLoginAttemptIDs(ids ...int) *UserUpdate {
	uu.mutation.RemoveLoginAttemptIDs(ids...)
	return uu
}

// RemoveLoginAttempts removes "loginAttempts" edges to LoginAttempt entities.
func (uu *UserUpdate) RemoveLoginAttempts(l ...*LoginAttempt) *UserUpdate {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return uu.RemoveLoginAttemptIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (uu *UserUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, uu.sqlSave, uu.mutation, uu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (uu *UserUpdate) SaveX(ctx context.Context) int {
	affected, err := uu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (uu *UserUpdate) Exec(ctx context.Context) error {
	_, err := uu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uu *UserUpdate) ExecX(ctx context.Context) {
	if err := uu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (uu *UserUpdate) check() error {
	if v, ok := uu.mutation.Username(); ok {
		if err := user.UsernameValidator(v); err != nil {
			return &ValidationError{Name: "username", err: fmt.Errorf(`ent: validator failed for field "User.username": %w`, err)}
		}
	}
	if v, ok := uu.mutation.Content(); ok {
		if err := user.ContentValidator(v); err != nil {
			return &ValidationError{Name: "content", err: fmt.Errorf(`ent: validator failed for field "User.content": %w`, err)}
		}
	}
	if v, ok := uu.mutation.FileName(); ok {
		if err := user.FileNameValidator(v); err != nil {
			return &ValidationError{Name: "fileName", err: fmt.Errorf(`ent: validator failed for field "User.fileName": %w`, err)}
		}
	}
	if v, ok := uu.mutation.Mime(); ok {
		if err := user.MimeValidator(v); err != nil {
			return &ValidationError{Name: "mime", err: fmt.Errorf(`ent: validator failed for field "User.mime": %w`, err)}
		}
	}
	if v, ok := uu.mutation.Nonce(); ok {
		if err := user.NonceValidator(v); err != nil {
			return &ValidationError{Name: "nonce", err: fmt.Errorf(`ent: validator failed for field "User.nonce": %w`, err)}
		}
	}
	if v, ok := uu.mutation.KeySalt(); ok {
		if err := user.KeySaltValidator(v); err != nil {
			return &ValidationError{Name: "keySalt", err: fmt.Errorf(`ent: validator failed for field "User.keySalt": %w`, err)}
		}
	}
	if v, ok := uu.mutation.PasswordHash(); ok {
		if err := user.PasswordHashValidator(v); err != nil {
			return &ValidationError{Name: "passwordHash", err: fmt.Errorf(`ent: validator failed for field "User.passwordHash": %w`, err)}
		}
	}
	if v, ok := uu.mutation.PasswordSalt(); ok {
		if err := user.PasswordSaltValidator(v); err != nil {
			return &ValidationError{Name: "passwordSalt", err: fmt.Errorf(`ent: validator failed for field "User.passwordSalt": %w`, err)}
		}
	}
	return nil
}

func (uu *UserUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := uu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(user.Table, user.Columns, sqlgraph.NewFieldSpec(user.FieldID, field.TypeInt))
	if ps := uu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := uu.mutation.Username(); ok {
		_spec.SetField(user.FieldUsername, field.TypeString, value)
	}
	if value, ok := uu.mutation.AlertDiscordId(); ok {
		_spec.SetField(user.FieldAlertDiscordId, field.TypeString, value)
	}
	if value, ok := uu.mutation.AlertEmail(); ok {
		_spec.SetField(user.FieldAlertEmail, field.TypeString, value)
	}
	if value, ok := uu.mutation.Content(); ok {
		_spec.SetField(user.FieldContent, field.TypeBytes, value)
	}
	if value, ok := uu.mutation.FileName(); ok {
		_spec.SetField(user.FieldFileName, field.TypeString, value)
	}
	if value, ok := uu.mutation.Mime(); ok {
		_spec.SetField(user.FieldMime, field.TypeString, value)
	}
	if value, ok := uu.mutation.Nonce(); ok {
		_spec.SetField(user.FieldNonce, field.TypeBytes, value)
	}
	if value, ok := uu.mutation.KeySalt(); ok {
		_spec.SetField(user.FieldKeySalt, field.TypeBytes, value)
	}
	if value, ok := uu.mutation.PasswordHash(); ok {
		_spec.SetField(user.FieldPasswordHash, field.TypeBytes, value)
	}
	if value, ok := uu.mutation.PasswordSalt(); ok {
		_spec.SetField(user.FieldPasswordSalt, field.TypeBytes, value)
	}
	if value, ok := uu.mutation.HashTime(); ok {
		_spec.SetField(user.FieldHashTime, field.TypeUint32, value)
	}
	if value, ok := uu.mutation.AddedHashTime(); ok {
		_spec.AddField(user.FieldHashTime, field.TypeUint32, value)
	}
	if value, ok := uu.mutation.HashMemory(); ok {
		_spec.SetField(user.FieldHashMemory, field.TypeUint32, value)
	}
	if value, ok := uu.mutation.AddedHashMemory(); ok {
		_spec.AddField(user.FieldHashMemory, field.TypeUint32, value)
	}
	if value, ok := uu.mutation.HashKeyLen(); ok {
		_spec.SetField(user.FieldHashKeyLen, field.TypeUint32, value)
	}
	if value, ok := uu.mutation.AddedHashKeyLen(); ok {
		_spec.AddField(user.FieldHashKeyLen, field.TypeUint32, value)
	}
	if uu.mutation.LoginAttemptsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.LoginAttemptsTable,
			Columns: []string{user.LoginAttemptsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(loginattempt.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.RemovedLoginAttemptsIDs(); len(nodes) > 0 && !uu.mutation.LoginAttemptsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.LoginAttemptsTable,
			Columns: []string{user.LoginAttemptsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(loginattempt.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.LoginAttemptsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.LoginAttemptsTable,
			Columns: []string{user.LoginAttemptsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(loginattempt.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, uu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{user.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	uu.mutation.done = true
	return n, nil
}

// UserUpdateOne is the builder for updating a single User entity.
type UserUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *UserMutation
}

// SetUsername sets the "username" field.
func (uuo *UserUpdateOne) SetUsername(s string) *UserUpdateOne {
	uuo.mutation.SetUsername(s)
	return uuo
}

// SetNillableUsername sets the "username" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableUsername(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetUsername(*s)
	}
	return uuo
}

// SetAlertDiscordId sets the "alertDiscordId" field.
func (uuo *UserUpdateOne) SetAlertDiscordId(s string) *UserUpdateOne {
	uuo.mutation.SetAlertDiscordId(s)
	return uuo
}

// SetNillableAlertDiscordId sets the "alertDiscordId" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableAlertDiscordId(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetAlertDiscordId(*s)
	}
	return uuo
}

// SetAlertEmail sets the "alertEmail" field.
func (uuo *UserUpdateOne) SetAlertEmail(s string) *UserUpdateOne {
	uuo.mutation.SetAlertEmail(s)
	return uuo
}

// SetNillableAlertEmail sets the "alertEmail" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableAlertEmail(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetAlertEmail(*s)
	}
	return uuo
}

// SetContent sets the "content" field.
func (uuo *UserUpdateOne) SetContent(b []byte) *UserUpdateOne {
	uuo.mutation.SetContent(b)
	return uuo
}

// SetFileName sets the "fileName" field.
func (uuo *UserUpdateOne) SetFileName(s string) *UserUpdateOne {
	uuo.mutation.SetFileName(s)
	return uuo
}

// SetNillableFileName sets the "fileName" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableFileName(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetFileName(*s)
	}
	return uuo
}

// SetMime sets the "mime" field.
func (uuo *UserUpdateOne) SetMime(s string) *UserUpdateOne {
	uuo.mutation.SetMime(s)
	return uuo
}

// SetNillableMime sets the "mime" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableMime(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetMime(*s)
	}
	return uuo
}

// SetNonce sets the "nonce" field.
func (uuo *UserUpdateOne) SetNonce(b []byte) *UserUpdateOne {
	uuo.mutation.SetNonce(b)
	return uuo
}

// SetKeySalt sets the "keySalt" field.
func (uuo *UserUpdateOne) SetKeySalt(b []byte) *UserUpdateOne {
	uuo.mutation.SetKeySalt(b)
	return uuo
}

// SetPasswordHash sets the "passwordHash" field.
func (uuo *UserUpdateOne) SetPasswordHash(b []byte) *UserUpdateOne {
	uuo.mutation.SetPasswordHash(b)
	return uuo
}

// SetPasswordSalt sets the "passwordSalt" field.
func (uuo *UserUpdateOne) SetPasswordSalt(b []byte) *UserUpdateOne {
	uuo.mutation.SetPasswordSalt(b)
	return uuo
}

// SetHashTime sets the "hashTime" field.
func (uuo *UserUpdateOne) SetHashTime(u uint32) *UserUpdateOne {
	uuo.mutation.ResetHashTime()
	uuo.mutation.SetHashTime(u)
	return uuo
}

// SetNillableHashTime sets the "hashTime" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableHashTime(u *uint32) *UserUpdateOne {
	if u != nil {
		uuo.SetHashTime(*u)
	}
	return uuo
}

// AddHashTime adds u to the "hashTime" field.
func (uuo *UserUpdateOne) AddHashTime(u int32) *UserUpdateOne {
	uuo.mutation.AddHashTime(u)
	return uuo
}

// SetHashMemory sets the "hashMemory" field.
func (uuo *UserUpdateOne) SetHashMemory(u uint32) *UserUpdateOne {
	uuo.mutation.ResetHashMemory()
	uuo.mutation.SetHashMemory(u)
	return uuo
}

// SetNillableHashMemory sets the "hashMemory" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableHashMemory(u *uint32) *UserUpdateOne {
	if u != nil {
		uuo.SetHashMemory(*u)
	}
	return uuo
}

// AddHashMemory adds u to the "hashMemory" field.
func (uuo *UserUpdateOne) AddHashMemory(u int32) *UserUpdateOne {
	uuo.mutation.AddHashMemory(u)
	return uuo
}

// SetHashKeyLen sets the "hashKeyLen" field.
func (uuo *UserUpdateOne) SetHashKeyLen(u uint32) *UserUpdateOne {
	uuo.mutation.ResetHashKeyLen()
	uuo.mutation.SetHashKeyLen(u)
	return uuo
}

// SetNillableHashKeyLen sets the "hashKeyLen" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableHashKeyLen(u *uint32) *UserUpdateOne {
	if u != nil {
		uuo.SetHashKeyLen(*u)
	}
	return uuo
}

// AddHashKeyLen adds u to the "hashKeyLen" field.
func (uuo *UserUpdateOne) AddHashKeyLen(u int32) *UserUpdateOne {
	uuo.mutation.AddHashKeyLen(u)
	return uuo
}

// AddLoginAttemptIDs adds the "loginAttempts" edge to the LoginAttempt entity by IDs.
func (uuo *UserUpdateOne) AddLoginAttemptIDs(ids ...int) *UserUpdateOne {
	uuo.mutation.AddLoginAttemptIDs(ids...)
	return uuo
}

// AddLoginAttempts adds the "loginAttempts" edges to the LoginAttempt entity.
func (uuo *UserUpdateOne) AddLoginAttempts(l ...*LoginAttempt) *UserUpdateOne {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return uuo.AddLoginAttemptIDs(ids...)
}

// Mutation returns the UserMutation object of the builder.
func (uuo *UserUpdateOne) Mutation() *UserMutation {
	return uuo.mutation
}

// ClearLoginAttempts clears all "loginAttempts" edges to the LoginAttempt entity.
func (uuo *UserUpdateOne) ClearLoginAttempts() *UserUpdateOne {
	uuo.mutation.ClearLoginAttempts()
	return uuo
}

// RemoveLoginAttemptIDs removes the "loginAttempts" edge to LoginAttempt entities by IDs.
func (uuo *UserUpdateOne) RemoveLoginAttemptIDs(ids ...int) *UserUpdateOne {
	uuo.mutation.RemoveLoginAttemptIDs(ids...)
	return uuo
}

// RemoveLoginAttempts removes "loginAttempts" edges to LoginAttempt entities.
func (uuo *UserUpdateOne) RemoveLoginAttempts(l ...*LoginAttempt) *UserUpdateOne {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return uuo.RemoveLoginAttemptIDs(ids...)
}

// Where appends a list predicates to the UserUpdate builder.
func (uuo *UserUpdateOne) Where(ps ...predicate.User) *UserUpdateOne {
	uuo.mutation.Where(ps...)
	return uuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (uuo *UserUpdateOne) Select(field string, fields ...string) *UserUpdateOne {
	uuo.fields = append([]string{field}, fields...)
	return uuo
}

// Save executes the query and returns the updated User entity.
func (uuo *UserUpdateOne) Save(ctx context.Context) (*User, error) {
	return withHooks(ctx, uuo.sqlSave, uuo.mutation, uuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (uuo *UserUpdateOne) SaveX(ctx context.Context) *User {
	node, err := uuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (uuo *UserUpdateOne) Exec(ctx context.Context) error {
	_, err := uuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uuo *UserUpdateOne) ExecX(ctx context.Context) {
	if err := uuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (uuo *UserUpdateOne) check() error {
	if v, ok := uuo.mutation.Username(); ok {
		if err := user.UsernameValidator(v); err != nil {
			return &ValidationError{Name: "username", err: fmt.Errorf(`ent: validator failed for field "User.username": %w`, err)}
		}
	}
	if v, ok := uuo.mutation.Content(); ok {
		if err := user.ContentValidator(v); err != nil {
			return &ValidationError{Name: "content", err: fmt.Errorf(`ent: validator failed for field "User.content": %w`, err)}
		}
	}
	if v, ok := uuo.mutation.FileName(); ok {
		if err := user.FileNameValidator(v); err != nil {
			return &ValidationError{Name: "fileName", err: fmt.Errorf(`ent: validator failed for field "User.fileName": %w`, err)}
		}
	}
	if v, ok := uuo.mutation.Mime(); ok {
		if err := user.MimeValidator(v); err != nil {
			return &ValidationError{Name: "mime", err: fmt.Errorf(`ent: validator failed for field "User.mime": %w`, err)}
		}
	}
	if v, ok := uuo.mutation.Nonce(); ok {
		if err := user.NonceValidator(v); err != nil {
			return &ValidationError{Name: "nonce", err: fmt.Errorf(`ent: validator failed for field "User.nonce": %w`, err)}
		}
	}
	if v, ok := uuo.mutation.KeySalt(); ok {
		if err := user.KeySaltValidator(v); err != nil {
			return &ValidationError{Name: "keySalt", err: fmt.Errorf(`ent: validator failed for field "User.keySalt": %w`, err)}
		}
	}
	if v, ok := uuo.mutation.PasswordHash(); ok {
		if err := user.PasswordHashValidator(v); err != nil {
			return &ValidationError{Name: "passwordHash", err: fmt.Errorf(`ent: validator failed for field "User.passwordHash": %w`, err)}
		}
	}
	if v, ok := uuo.mutation.PasswordSalt(); ok {
		if err := user.PasswordSaltValidator(v); err != nil {
			return &ValidationError{Name: "passwordSalt", err: fmt.Errorf(`ent: validator failed for field "User.passwordSalt": %w`, err)}
		}
	}
	return nil
}

func (uuo *UserUpdateOne) sqlSave(ctx context.Context) (_node *User, err error) {
	if err := uuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(user.Table, user.Columns, sqlgraph.NewFieldSpec(user.FieldID, field.TypeInt))
	id, ok := uuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "User.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := uuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, user.FieldID)
		for _, f := range fields {
			if !user.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != user.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := uuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := uuo.mutation.Username(); ok {
		_spec.SetField(user.FieldUsername, field.TypeString, value)
	}
	if value, ok := uuo.mutation.AlertDiscordId(); ok {
		_spec.SetField(user.FieldAlertDiscordId, field.TypeString, value)
	}
	if value, ok := uuo.mutation.AlertEmail(); ok {
		_spec.SetField(user.FieldAlertEmail, field.TypeString, value)
	}
	if value, ok := uuo.mutation.Content(); ok {
		_spec.SetField(user.FieldContent, field.TypeBytes, value)
	}
	if value, ok := uuo.mutation.FileName(); ok {
		_spec.SetField(user.FieldFileName, field.TypeString, value)
	}
	if value, ok := uuo.mutation.Mime(); ok {
		_spec.SetField(user.FieldMime, field.TypeString, value)
	}
	if value, ok := uuo.mutation.Nonce(); ok {
		_spec.SetField(user.FieldNonce, field.TypeBytes, value)
	}
	if value, ok := uuo.mutation.KeySalt(); ok {
		_spec.SetField(user.FieldKeySalt, field.TypeBytes, value)
	}
	if value, ok := uuo.mutation.PasswordHash(); ok {
		_spec.SetField(user.FieldPasswordHash, field.TypeBytes, value)
	}
	if value, ok := uuo.mutation.PasswordSalt(); ok {
		_spec.SetField(user.FieldPasswordSalt, field.TypeBytes, value)
	}
	if value, ok := uuo.mutation.HashTime(); ok {
		_spec.SetField(user.FieldHashTime, field.TypeUint32, value)
	}
	if value, ok := uuo.mutation.AddedHashTime(); ok {
		_spec.AddField(user.FieldHashTime, field.TypeUint32, value)
	}
	if value, ok := uuo.mutation.HashMemory(); ok {
		_spec.SetField(user.FieldHashMemory, field.TypeUint32, value)
	}
	if value, ok := uuo.mutation.AddedHashMemory(); ok {
		_spec.AddField(user.FieldHashMemory, field.TypeUint32, value)
	}
	if value, ok := uuo.mutation.HashKeyLen(); ok {
		_spec.SetField(user.FieldHashKeyLen, field.TypeUint32, value)
	}
	if value, ok := uuo.mutation.AddedHashKeyLen(); ok {
		_spec.AddField(user.FieldHashKeyLen, field.TypeUint32, value)
	}
	if uuo.mutation.LoginAttemptsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.LoginAttemptsTable,
			Columns: []string{user.LoginAttemptsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(loginattempt.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.RemovedLoginAttemptsIDs(); len(nodes) > 0 && !uuo.mutation.LoginAttemptsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.LoginAttemptsTable,
			Columns: []string{user.LoginAttemptsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(loginattempt.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.LoginAttemptsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.LoginAttemptsTable,
			Columns: []string{user.LoginAttemptsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(loginattempt.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &User{config: uuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, uuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{user.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	uuo.mutation.done = true
	return _node, nil
}
