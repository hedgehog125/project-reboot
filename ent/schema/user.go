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
		field.String("username").Unique(),
		field.String("alertDiscordId").Optional(),
		field.String("alertEmail").Optional(),

		field.Bytes("content"),
		field.String("fileName"),
		field.String("mime"),
		field.Bytes("nonce"),
		field.Bytes("keySalt"),
		field.Bytes("passwordHash"),
		field.Bytes("passwordSalt"),
		field.Uint32("hashTime"),
		field.Uint32("hashMemory"),
		field.Uint32("hashKeyLen"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
