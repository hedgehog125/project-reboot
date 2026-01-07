package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
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
		// Admins might be able to be locked in the future
		field.Bool("locked").Default(false),
		field.Time("lockedUntil").Nillable().Optional(),
		field.Time("sessionsValidFrom"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("stash", Stash.Type).
			Annotations(entsql.OnDelete(entsql.Restrict)).Unique(),
		edge.To("sessions", Session.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("logs", LogEntry.Type).
			Annotations(entsql.OnDelete(entsql.SetNull)),
	}
}
