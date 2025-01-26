// Code generated by ent, DO NOT EDIT.

package user

import (
	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldUsername holds the string denoting the username field in the database.
	FieldUsername = "username"
	// FieldContent holds the string denoting the content field in the database.
	FieldContent = "content"
	// FieldFileName holds the string denoting the filename field in the database.
	FieldFileName = "file_name"
	// FieldMime holds the string denoting the mime field in the database.
	FieldMime = "mime"
	// FieldNonce holds the string denoting the nonce field in the database.
	FieldNonce = "nonce"
	// FieldKeySalt holds the string denoting the keysalt field in the database.
	FieldKeySalt = "key_salt"
	// FieldPasswordHash holds the string denoting the passwordhash field in the database.
	FieldPasswordHash = "password_hash"
	// FieldPasswordSalt holds the string denoting the passwordsalt field in the database.
	FieldPasswordSalt = "password_salt"
	// FieldHashTime holds the string denoting the hashtime field in the database.
	FieldHashTime = "hash_time"
	// FieldHashMemory holds the string denoting the hashmemory field in the database.
	FieldHashMemory = "hash_memory"
	// FieldHashKeyLen holds the string denoting the hashkeylen field in the database.
	FieldHashKeyLen = "hash_key_len"
	// Table holds the table name of the user in the database.
	Table = "users"
)

// Columns holds all SQL columns for user fields.
var Columns = []string{
	FieldID,
	FieldUsername,
	FieldContent,
	FieldFileName,
	FieldMime,
	FieldNonce,
	FieldKeySalt,
	FieldPasswordHash,
	FieldPasswordSalt,
	FieldHashTime,
	FieldHashMemory,
	FieldHashKeyLen,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the User queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByUsername orders the results by the username field.
func ByUsername(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUsername, opts...).ToFunc()
}

// ByFileName orders the results by the fileName field.
func ByFileName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFileName, opts...).ToFunc()
}

// ByMime orders the results by the mime field.
func ByMime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMime, opts...).ToFunc()
}

// ByHashTime orders the results by the hashTime field.
func ByHashTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldHashTime, opts...).ToFunc()
}

// ByHashMemory orders the results by the hashMemory field.
func ByHashMemory(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldHashMemory, opts...).ToFunc()
}

// ByHashKeyLen orders the results by the hashKeyLen field.
func ByHashKeyLen(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldHashKeyLen, opts...).ToFunc()
}
