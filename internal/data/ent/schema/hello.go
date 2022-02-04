package schema

import "entgo.io/ent"

// Hello holds the schema definition for the Hello entity.
type Hello struct {
	ent.Schema
}

// Fields of the Hello.
func (Hello) Fields() []ent.Field {
	return nil
}

// Edges of the Hello.
func (Hello) Edges() []ent.Edge {
	return nil
}
