package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

// Template holds the schema definition for the Template entity.
type Template struct {
	ent.Schema
}

// Fields of the Template.
func (Template) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("name").Unique(),
		field.String("type"),
		field.String("template").
			SchemaType(map[string]string{
				dialect.Postgres: "text",
			}),
		field.String("created_by"),
		field.Time("created_at").
			Default(time.Now),
		field.String("updated_by").
			Optional().
			Nillable(),
		field.Time("updated_at").
			Optional().
			Nillable(),
		field.String("deleted_by").
			Optional().
			Nillable(),
		field.Time("deleted_at").
			Optional().
			Nillable(),
	}
}

// Edges of the Template.
func (Template) Edges() []ent.Edge {
	return nil
}
