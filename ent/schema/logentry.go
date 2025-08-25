package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// LogEntry holds the schema definition for the LogEntry entity.
type LogEntry struct {
	ent.Schema
}

// Fields of the LogEntry.
func (LogEntry) Fields() []ent.Field {
	return []ent.Field{field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.Time("time"),      // The entries should be batch created, so a default time wouldn't be accurate
		field.Bool("timeKnown"), // Some logs don't have a time, so an inaccurate time is added during processing
		field.Int("level"),
		field.String("message"),
		field.JSON("attributes", map[string]any{}),
		field.String("sourceFile"),
		field.String("sourceFunction"),
		field.Int("sourceLine"),
		field.String("publicMessage"),
	}
}

// Edges of the LogEntry.
func (LogEntry) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).Unique().Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
