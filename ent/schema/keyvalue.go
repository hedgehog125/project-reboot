package schema

import (
	"encoding/json"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// KeyValue holds the schema definition for the KeyValue entity.
type KeyValue struct {
	ent.Schema
}

// Fields of the KeyValue.
func (KeyValue) Fields() []ent.Field {
	return []ent.Field{
		field.String("key").MinLen(1).MaxLen(128),
		field.JSON("value", json.RawMessage{}),
	}
}

// Edges of the KeyValue.
func (KeyValue) Edges() []ent.Edge {
	return nil
}

func (KeyValue) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("key").Unique(),
	}
}
