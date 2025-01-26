// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/hedgehog125/project-reboot/ent/user"
)

// User is the model entity for the User schema.
type User struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Username holds the value of the "username" field.
	Username string `json:"username,omitempty"`
	// Content holds the value of the "content" field.
	Content []byte `json:"content,omitempty"`
	// FileName holds the value of the "fileName" field.
	FileName string `json:"fileName,omitempty"`
	// Mime holds the value of the "mime" field.
	Mime string `json:"mime,omitempty"`
	// Nonce holds the value of the "nonce" field.
	Nonce []byte `json:"nonce,omitempty"`
	// KeySalt holds the value of the "keySalt" field.
	KeySalt []byte `json:"keySalt,omitempty"`
	// PasswordHash holds the value of the "passwordHash" field.
	PasswordHash []byte `json:"passwordHash,omitempty"`
	// PasswordSalt holds the value of the "passwordSalt" field.
	PasswordSalt []byte `json:"passwordSalt,omitempty"`
	// HashTime holds the value of the "hashTime" field.
	HashTime uint32 `json:"hashTime,omitempty"`
	// HashMemory holds the value of the "hashMemory" field.
	HashMemory uint32 `json:"hashMemory,omitempty"`
	// HashKeyLen holds the value of the "hashKeyLen" field.
	HashKeyLen   uint32 `json:"hashKeyLen,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*User) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case user.FieldContent, user.FieldNonce, user.FieldKeySalt, user.FieldPasswordHash, user.FieldPasswordSalt:
			values[i] = new([]byte)
		case user.FieldID, user.FieldHashTime, user.FieldHashMemory, user.FieldHashKeyLen:
			values[i] = new(sql.NullInt64)
		case user.FieldUsername, user.FieldFileName, user.FieldMime:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the User fields.
func (u *User) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case user.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			u.ID = int(value.Int64)
		case user.FieldUsername:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field username", values[i])
			} else if value.Valid {
				u.Username = value.String
			}
		case user.FieldContent:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field content", values[i])
			} else if value != nil {
				u.Content = *value
			}
		case user.FieldFileName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field fileName", values[i])
			} else if value.Valid {
				u.FileName = value.String
			}
		case user.FieldMime:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field mime", values[i])
			} else if value.Valid {
				u.Mime = value.String
			}
		case user.FieldNonce:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field nonce", values[i])
			} else if value != nil {
				u.Nonce = *value
			}
		case user.FieldKeySalt:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field keySalt", values[i])
			} else if value != nil {
				u.KeySalt = *value
			}
		case user.FieldPasswordHash:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field passwordHash", values[i])
			} else if value != nil {
				u.PasswordHash = *value
			}
		case user.FieldPasswordSalt:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field passwordSalt", values[i])
			} else if value != nil {
				u.PasswordSalt = *value
			}
		case user.FieldHashTime:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field hashTime", values[i])
			} else if value.Valid {
				u.HashTime = uint32(value.Int64)
			}
		case user.FieldHashMemory:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field hashMemory", values[i])
			} else if value.Valid {
				u.HashMemory = uint32(value.Int64)
			}
		case user.FieldHashKeyLen:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field hashKeyLen", values[i])
			} else if value.Valid {
				u.HashKeyLen = uint32(value.Int64)
			}
		default:
			u.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the User.
// This includes values selected through modifiers, order, etc.
func (u *User) Value(name string) (ent.Value, error) {
	return u.selectValues.Get(name)
}

// Update returns a builder for updating this User.
// Note that you need to call User.Unwrap() before calling this method if this User
// was returned from a transaction, and the transaction was committed or rolled back.
func (u *User) Update() *UserUpdateOne {
	return NewUserClient(u.config).UpdateOne(u)
}

// Unwrap unwraps the User entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (u *User) Unwrap() *User {
	_tx, ok := u.config.driver.(*txDriver)
	if !ok {
		panic("ent: User is not a transactional entity")
	}
	u.config.driver = _tx.drv
	return u
}

// String implements the fmt.Stringer.
func (u *User) String() string {
	var builder strings.Builder
	builder.WriteString("User(")
	builder.WriteString(fmt.Sprintf("id=%v, ", u.ID))
	builder.WriteString("username=")
	builder.WriteString(u.Username)
	builder.WriteString(", ")
	builder.WriteString("content=")
	builder.WriteString(fmt.Sprintf("%v", u.Content))
	builder.WriteString(", ")
	builder.WriteString("fileName=")
	builder.WriteString(u.FileName)
	builder.WriteString(", ")
	builder.WriteString("mime=")
	builder.WriteString(u.Mime)
	builder.WriteString(", ")
	builder.WriteString("nonce=")
	builder.WriteString(fmt.Sprintf("%v", u.Nonce))
	builder.WriteString(", ")
	builder.WriteString("keySalt=")
	builder.WriteString(fmt.Sprintf("%v", u.KeySalt))
	builder.WriteString(", ")
	builder.WriteString("passwordHash=")
	builder.WriteString(fmt.Sprintf("%v", u.PasswordHash))
	builder.WriteString(", ")
	builder.WriteString("passwordSalt=")
	builder.WriteString(fmt.Sprintf("%v", u.PasswordSalt))
	builder.WriteString(", ")
	builder.WriteString("hashTime=")
	builder.WriteString(fmt.Sprintf("%v", u.HashTime))
	builder.WriteString(", ")
	builder.WriteString("hashMemory=")
	builder.WriteString(fmt.Sprintf("%v", u.HashMemory))
	builder.WriteString(", ")
	builder.WriteString("hashKeyLen=")
	builder.WriteString(fmt.Sprintf("%v", u.HashKeyLen))
	builder.WriteByte(')')
	return builder.String()
}

// Users is a parsable slice of User.
type Users []*User
