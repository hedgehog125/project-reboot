package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/hedgehog125/project-reboot/intertypes"
)

// LoginAttempt holds the schema definition for the LoginAttempt entity.
type LoginAttempt struct {
	ent.Schema
}

// Fields of the LoginAttempt.
func (LoginAttempt) Fields() []ent.Field { // TODO: auto delete once used? Or also after a certain amount of time
	return []ent.Field{
		field.Time("time").Default(time.Now),
		field.Bytes("code").Unique().MinLen(128), // The randomly generated authorisation code that will become valid after enough time
		field.Time("codeValidFrom"),
		field.JSON("info", &intertypes.LoginAttemptInfo{}),
	}
}

// Edges of the LoginAttempt.
func (LoginAttempt) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("loginAttempts").Unique(),
	}
}
