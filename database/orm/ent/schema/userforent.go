package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
)

type UserForEnt struct {
	ent.Schema
}

func (UserForEnt) Fields() []ent.Field {
	return nil
}

func (UserForEnt) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("repository_for_ents", RepositoryForEnt.Type).Ref("user_for_ents"),
	}
}
