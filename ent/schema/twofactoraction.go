package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// TwoFactorAction holds the schema definition for the TwoFactorAction entity.
type TwoFactorAction struct {
	ent.Schema
}

// Fields of the TwoFactorAction.
func (TwoFactorAction) Fields() []ent.Field {
	return []ent.Field{
		field.String("type").MinLen(1).MaxLen(128),
		field.Int("version"),
		field.JSON("data", &map[string]any{}),
		field.Time("expiresAt"),
		field.String("code").MaxLen(6).MaxLen(6),
	}
}

// Edges of the TwoFactorAction.
func (TwoFactorAction) Edges() []ent.Edge {
	return nil
}
