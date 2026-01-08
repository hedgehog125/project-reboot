package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Stash holds the schema definition for the Stash entity.
type Stash struct {
	ent.Schema
}

// Fields of the Stash.
func (Stash) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.Bytes("content").NotEmpty(),
		field.String("fileName").NotEmpty(),
		field.String("mime").NotEmpty(),
		field.Bytes("nonce").NotEmpty(),
		field.Bytes("keySalt").NotEmpty(),
		field.Uint32("hashTime"),
		field.Uint32("hashMemory"),
		field.Uint8("hashThreads"),
		field.UUID("userID", uuid.UUID{}),
	}
}

// Edges of the Stash.
func (Stash) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("stash").
			Field("userID").Unique().Required(),
	}
}
