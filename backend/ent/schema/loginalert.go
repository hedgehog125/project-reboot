package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// LoginAlert holds the schema definition for the LoginAlert entity.
type LoginAlert struct {
	ent.Schema
}

// Fields of the SuccessfulLoginAlerts.
func (LoginAlert) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.Nil).Default(uuid.New),
		field.Time("sentAt"),
		field.String("versionedMessengerType").MinLen(1).MaxLen(128),
		field.Bool("confirmed"),
		field.UUID("sessionID", uuid.Nil),
	}
}

// Edges of the SuccessfulLoginAlerts.
func (LoginAlert) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("session", Session.Type).Ref("loginAlerts").
			Field("sessionID").Unique().Required(),
	}
}
