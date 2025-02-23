// Code generated by ent, DO NOT EDIT.

package loginattempt

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the loginattempt type in the database.
	Label = "login_attempt"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldTime holds the string denoting the time field in the database.
	FieldTime = "time"
	// FieldCode holds the string denoting the code field in the database.
	FieldCode = "code"
	// FieldCodeValidFrom holds the string denoting the codevalidfrom field in the database.
	FieldCodeValidFrom = "code_valid_from"
	// FieldInfo holds the string denoting the info field in the database.
	FieldInfo = "info"
	// EdgeUser holds the string denoting the user edge name in mutations.
	EdgeUser = "user"
	// Table holds the table name of the loginattempt in the database.
	Table = "login_attempts"
	// UserTable is the table that holds the user relation/edge.
	UserTable = "login_attempts"
	// UserInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	UserInverseTable = "users"
	// UserColumn is the table column denoting the user relation/edge.
	UserColumn = "user_login_attempts"
)

// Columns holds all SQL columns for loginattempt fields.
var Columns = []string{
	FieldID,
	FieldTime,
	FieldCode,
	FieldCodeValidFrom,
	FieldInfo,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "login_attempts"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"user_login_attempts",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultTime holds the default value on creation for the "time" field.
	DefaultTime func() time.Time
	// CodeValidator is a validator for the "code" field. It is called by the builders before save.
	CodeValidator func([]byte) error
)

// OrderOption defines the ordering options for the LoginAttempt queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByTime orders the results by the time field.
func ByTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTime, opts...).ToFunc()
}

// ByCodeValidFrom orders the results by the codeValidFrom field.
func ByCodeValidFrom(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCodeValidFrom, opts...).ToFunc()
}

// ByUserField orders the results by user field.
func ByUserField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newUserStep(), sql.OrderByField(field, opts...))
	}
}
func newUserStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(UserInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, UserTable, UserColumn),
	)
}
