// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/hedgehog125/project-reboot/ent/loginattempt"
	"github.com/hedgehog125/project-reboot/intertypes"
)

// LoginAttempt is the model entity for the LoginAttempt schema.
type LoginAttempt struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Time holds the value of the "time" field.
	Time time.Time `json:"time,omitempty"`
	// Username holds the value of the "username" field.
	Username time.Time `json:"username,omitempty"`
	// Code holds the value of the "code" field.
	Code string `json:"code,omitempty"`
	// CodeValidFrom holds the value of the "codeValidFrom" field.
	CodeValidFrom time.Time `json:"codeValidFrom,omitempty"`
	// Info holds the value of the "info" field.
	Info         *intertypes.LoginAttemptInfo `json:"info,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*LoginAttempt) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case loginattempt.FieldInfo:
			values[i] = new([]byte)
		case loginattempt.FieldID:
			values[i] = new(sql.NullInt64)
		case loginattempt.FieldCode:
			values[i] = new(sql.NullString)
		case loginattempt.FieldTime, loginattempt.FieldUsername, loginattempt.FieldCodeValidFrom:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the LoginAttempt fields.
func (la *LoginAttempt) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case loginattempt.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			la.ID = int(value.Int64)
		case loginattempt.FieldTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field time", values[i])
			} else if value.Valid {
				la.Time = value.Time
			}
		case loginattempt.FieldUsername:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field username", values[i])
			} else if value.Valid {
				la.Username = value.Time
			}
		case loginattempt.FieldCode:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field code", values[i])
			} else if value.Valid {
				la.Code = value.String
			}
		case loginattempt.FieldCodeValidFrom:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field codeValidFrom", values[i])
			} else if value.Valid {
				la.CodeValidFrom = value.Time
			}
		case loginattempt.FieldInfo:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field info", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &la.Info); err != nil {
					return fmt.Errorf("unmarshal field info: %w", err)
				}
			}
		default:
			la.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the LoginAttempt.
// This includes values selected through modifiers, order, etc.
func (la *LoginAttempt) Value(name string) (ent.Value, error) {
	return la.selectValues.Get(name)
}

// Update returns a builder for updating this LoginAttempt.
// Note that you need to call LoginAttempt.Unwrap() before calling this method if this LoginAttempt
// was returned from a transaction, and the transaction was committed or rolled back.
func (la *LoginAttempt) Update() *LoginAttemptUpdateOne {
	return NewLoginAttemptClient(la.config).UpdateOne(la)
}

// Unwrap unwraps the LoginAttempt entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (la *LoginAttempt) Unwrap() *LoginAttempt {
	_tx, ok := la.config.driver.(*txDriver)
	if !ok {
		panic("ent: LoginAttempt is not a transactional entity")
	}
	la.config.driver = _tx.drv
	return la
}

// String implements the fmt.Stringer.
func (la *LoginAttempt) String() string {
	var builder strings.Builder
	builder.WriteString("LoginAttempt(")
	builder.WriteString(fmt.Sprintf("id=%v, ", la.ID))
	builder.WriteString("time=")
	builder.WriteString(la.Time.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("username=")
	builder.WriteString(la.Username.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("code=")
	builder.WriteString(la.Code)
	builder.WriteString(", ")
	builder.WriteString("codeValidFrom=")
	builder.WriteString(la.CodeValidFrom.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("info=")
	builder.WriteString(fmt.Sprintf("%v", la.Info))
	builder.WriteByte(')')
	return builder.String()
}

// LoginAttempts is a parsable slice of LoginAttempt.
type LoginAttempts []*LoginAttempt
