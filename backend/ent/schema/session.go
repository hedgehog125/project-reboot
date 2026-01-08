package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Session holds the schema definition for the Session entity.
type Session struct {
	ent.Schema
}

// Fields of the Session.
func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.Time("createdAt"),
		field.Bytes("code").
			Unique().
			MinLen(128), // The randomly generated authorisation code that will become valid after enough time
		field.Time("validFrom"),
		field.Time("validUntil"),
		field.String("userAgent"),
		field.String("ip"),
		field.UUID("userID", uuid.UUID{}),
	}
}

// Edges of the Session.
func (Session) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("sessions").
			Field("userID").Unique().Required(),
		edge.To("loginAlerts", LoginAlert.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (Session) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("code", "userID"),
	}
}
