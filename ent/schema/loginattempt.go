package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/hedgehog125/project-reboot/intertypes"
)

// LoginAttempt holds the schema definition for the LoginAttempt entity.
type LoginAttempt struct {
	ent.Schema
}

// Fields of the LoginAttempt.
func (LoginAttempt) Fields() []ent.Field { // TODO: auto archive
	return []ent.Field{
		field.Time("time").Default(time.Now),
		field.String("username"),
		field.String("code"), // The randomly generated authorisation code that will become valid after enough time
		field.Time("codeValidFrom"),
		field.JSON("info", &intertypes.LoginAttemptInfo{}),
	}
}

// Edges of the LoginAttempt.
func (LoginAttempt) Edges() []ent.Edge {
	return nil
}
