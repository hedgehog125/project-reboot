// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// LoginAttemptsColumns holds the columns for the "login_attempts" table.
	LoginAttemptsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "time", Type: field.TypeTime},
		{Name: "code", Type: field.TypeString},
	}
	// LoginAttemptsTable holds the schema information for the "login_attempts" table.
	LoginAttemptsTable = &schema.Table{
		Name:       "login_attempts",
		Columns:    LoginAttemptsColumns,
		PrimaryKey: []*schema.Column{LoginAttemptsColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		LoginAttemptsTable,
	}
)

func init() {
}
