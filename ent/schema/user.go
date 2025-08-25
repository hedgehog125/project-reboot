package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").Unique().NotEmpty(),
		field.String("alertDiscordId").Default(""),
		field.String("alertEmail").Default(""),
		field.Bool("locked").Default(false),
		field.Time("lockedUntil").Nillable().Optional(),

		field.Bytes("content").NotEmpty(),
		field.String("fileName").NotEmpty(),
		field.String("mime").NotEmpty(),
		field.Bytes("nonce").NotEmpty(),
		field.Bytes("keySalt").NotEmpty(),
		field.Uint32("hashTime"),
		field.Uint32("hashMemory"),
		field.Uint8("hashThreads"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("sessions", Session.Type).Ref("user"),
		edge.From("logs", LogEntry.Type).Ref("user"),
	}
}
