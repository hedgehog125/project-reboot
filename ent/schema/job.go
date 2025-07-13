package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Job holds the schema definition for the Job entity.
type Job struct {
	ent.Schema
}

// Fields of the Job.
func (Job) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.Time("created").Default(time.Now),
		field.Time("due").Default(time.Now),
		field.String("type").MinLen(1).MaxLen(128),
		field.Int("version"),
		field.JSON("data", ""),
		field.Enum("status").Values("pending", "running", "failed").Default("pending"), // Completed jobs are deleted
		field.Int("retries").Default(0),
	}
}

// Edges of the Job.
func (Job) Edges() []ent.Edge {
	return nil
}

func (Job) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status", "due"),
	}
}
