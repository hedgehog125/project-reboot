package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// LoginAlert holds the schema definition for the LoginAlert entity.
type LoginAlert struct {
	ent.Schema
}

// Fields of the SuccessfulLoginAlerts.
func (LoginAlert) Fields() []ent.Field {
	return []ent.Field{
		field.Time("sentAt"),
		field.String("versionedMessengerType").MinLen(1).MaxLen(128),
		field.Bool("confirmed"),
		field.Int("sessionID"),
	}
}

// Edges of the SuccessfulLoginAlerts.
func (LoginAlert) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("session", Session.Type).Ref("loginAlerts").
			Field("sessionID").Unique().Required(),
	}
}
