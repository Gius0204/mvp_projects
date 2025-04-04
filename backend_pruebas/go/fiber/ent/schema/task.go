package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Task holds the schema definition for the Task entity.
type Task struct {
	ent.Schema
}

// Fields of the Task.
func (Task) Fields() []ent.Field {
	return []ent.Field{	
		field.Int("parent_id").
			Optional().
			Nillable(),

		field.Float("progress").
			Default(0.0),

		field.Int("estimated_time"),

		field.Time("created_at").
			Default(time.Now),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Task.
func (Task) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("parent", Task.Type).
			Ref("subtasks"). // <- antes era "children"
			Unique().
			Field("parent_id"),

		edge.To("subtasks", Task.Type), // <- antes era "children"
	}
}