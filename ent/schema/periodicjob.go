package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// PeriodicJob holds the schema definition for the PeriodicJob entity.
type PeriodicJob struct {
	ent.Schema
}

// Fields of the PeriodicJob.
func (PeriodicJob) Fields() []ent.Field {
	return []ent.Field{
		field.String("type").MinLen(1).MaxLen(128),
		field.Int("version"),
		field.Time("lastScheduledNewJob").Optional(),
	}
}

// Edges of the PeriodicJob.
func (PeriodicJob) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("jobs", Job.Type).Annotations(entsql.OnDelete(entsql.SetNull)),
	}
}

func (PeriodicJob) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("type", "version").Unique(),
	}
}
