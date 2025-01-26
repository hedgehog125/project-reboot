package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username"),

		field.Bytes("content"),
		field.String("mime"),
		field.Bytes("nonce"),
		field.Bytes("keySalt"),
		field.Bytes("passwordHash"),
		field.Bytes("passwordSalt"),
		field.Bytes("hashTime"),
		field.Bytes("hashMemory"),
		field.Bytes("hashKeyLen"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
