package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// TwoFactorAction holds the schema definition for the TwoFactorAction entity.
type TwoFactorAction struct {
	ent.Schema
}

// Fields of the TwoFactorAction.
func (TwoFactorAction) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("type").MinLen(1).MaxLen(128),
		field.Int("version"),
		field.JSON("data", ""),
		field.Time("expiresAt"),
		field.String("code").MaxLen(9).MaxLen(9),
	}
}

// Edges of the TwoFactorAction.
func (TwoFactorAction) Edges() []ent.Edge {
	return nil
}
