package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// PeriodicTask holds the schema definition for the PeriodicTask entity.
type PeriodicTask struct {
	ent.Schema
}

// Fields of the PeriodicJob.
func (PeriodicTask) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MinLen(1).MaxLen(128),
		// Note: there's no version because we should just be able to upgrade to new versions when the server restarts since there's no request body
		field.Time("lastRan").Optional(),
	}
}

// Edges of the PeriodicJob.
func (PeriodicTask) Edges() []ent.Edge {
	return nil
}

func (PeriodicTask) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
	}
}
