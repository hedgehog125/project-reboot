package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Session holds the schema definition for the Session entity.
type Session struct {
	ent.Schema
}

// Fields of the Session.
func (Session) Fields() []ent.Field { // TODO: auto delete once used? Or also after a certain amount of time
	return []ent.Field{
		field.Time("time").Default(time.Now),     // TODO: will this be an issue with testing?
		field.Bytes("code").Unique().MinLen(128), // The randomly generated authorisation code that will become valid after enough time
		field.Time("codeValidFrom"),
		field.String("userAgent"),
		field.String("ip"),
	}
}

// Edges of the Session.
func (Session) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("sessions").Unique(),
	}
}
