package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// LoginAlerts holds the schema definition for the LoginAlerts entity.
type LoginAlerts struct {
	ent.Schema
}

// Fields of the SuccessfulLoginAlerts.
func (LoginAlerts) Fields() []ent.Field {
	return []ent.Field{
		field.Time("time"),
		field.String("messengerType").MinLen(1).MaxLen(128),
		field.Bool("confirmed"),
		field.Int("sessionID"),
	}
}

// Edges of the SuccessfulLoginAlerts.
func (LoginAlerts) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("session", Session.Type).Ref("loginAlerts").
			Field("sessionID").Unique().Required(),
	}
}
