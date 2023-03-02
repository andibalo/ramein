package schema

import "entgo.io/ent"

// Template holds the schema definition for the Template entity.
type Template struct {
	ent.Schema
}

// Fields of the Template.
func (Template) Fields() []ent.Field {
	return nil
}

// Edges of the Template.
func (Template) Edges() []ent.Edge {
	return nil
}
