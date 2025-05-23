// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/hedgehog125/project-reboot/ent/session"
	"github.com/hedgehog125/project-reboot/ent/user"
)

// Session is the model entity for the Session schema.
type Session struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Time holds the value of the "time" field.
	Time time.Time `json:"time,omitempty"`
	// Code holds the value of the "code" field.
	Code []byte `json:"code,omitempty"`
	// CodeValidFrom holds the value of the "codeValidFrom" field.
	CodeValidFrom time.Time `json:"codeValidFrom,omitempty"`
	// UserAgent holds the value of the "userAgent" field.
	UserAgent string `json:"userAgent,omitempty"`
	// IP holds the value of the "ip" field.
	IP string `json:"ip,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the SessionQuery when eager-loading is set.
	Edges         SessionEdges `json:"edges"`
	user_sessions *int
	selectValues  sql.SelectValues
}

// SessionEdges holds the relations/edges for other nodes in the graph.
type SessionEdges struct {
	// User holds the value of the user edge.
	User *User `json:"user,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// UserOrErr returns the User value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e SessionEdges) UserOrErr() (*User, error) {
	if e.User != nil {
		return e.User, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: user.Label}
	}
	return nil, &NotLoadedError{edge: "user"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Session) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case session.FieldCode:
			values[i] = new([]byte)
		case session.FieldID:
			values[i] = new(sql.NullInt64)
		case session.FieldUserAgent, session.FieldIP:
			values[i] = new(sql.NullString)
		case session.FieldTime, session.FieldCodeValidFrom:
			values[i] = new(sql.NullTime)
		case session.ForeignKeys[0]: // user_sessions
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Session fields.
func (s *Session) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case session.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			s.ID = int(value.Int64)
		case session.FieldTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field time", values[i])
			} else if value.Valid {
				s.Time = value.Time
			}
		case session.FieldCode:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field code", values[i])
			} else if value != nil {
				s.Code = *value
			}
		case session.FieldCodeValidFrom:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field codeValidFrom", values[i])
			} else if value.Valid {
				s.CodeValidFrom = value.Time
			}
		case session.FieldUserAgent:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field userAgent", values[i])
			} else if value.Valid {
				s.UserAgent = value.String
			}
		case session.FieldIP:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field ip", values[i])
			} else if value.Valid {
				s.IP = value.String
			}
		case session.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field user_sessions", value)
			} else if value.Valid {
				s.user_sessions = new(int)
				*s.user_sessions = int(value.Int64)
			}
		default:
			s.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Session.
// This includes values selected through modifiers, order, etc.
func (s *Session) Value(name string) (ent.Value, error) {
	return s.selectValues.Get(name)
}

// QueryUser queries the "user" edge of the Session entity.
func (s *Session) QueryUser() *UserQuery {
	return NewSessionClient(s.config).QueryUser(s)
}

// Update returns a builder for updating this Session.
// Note that you need to call Session.Unwrap() before calling this method if this Session
// was returned from a transaction, and the transaction was committed or rolled back.
func (s *Session) Update() *SessionUpdateOne {
	return NewSessionClient(s.config).UpdateOne(s)
}

// Unwrap unwraps the Session entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (s *Session) Unwrap() *Session {
	_tx, ok := s.config.driver.(*txDriver)
	if !ok {
		panic("ent: Session is not a transactional entity")
	}
	s.config.driver = _tx.drv
	return s
}

// String implements the fmt.Stringer.
func (s *Session) String() string {
	var builder strings.Builder
	builder.WriteString("Session(")
	builder.WriteString(fmt.Sprintf("id=%v, ", s.ID))
	builder.WriteString("time=")
	builder.WriteString(s.Time.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("code=")
	builder.WriteString(fmt.Sprintf("%v", s.Code))
	builder.WriteString(", ")
	builder.WriteString("codeValidFrom=")
	builder.WriteString(s.CodeValidFrom.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("userAgent=")
	builder.WriteString(s.UserAgent)
	builder.WriteString(", ")
	builder.WriteString("ip=")
	builder.WriteString(s.IP)
	builder.WriteByte(')')
	return builder.String()
}

// Sessions is a parsable slice of Session.
type Sessions []*Session
