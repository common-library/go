package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
)

type RepositoryForEnt struct {
	ent.Schema
}

func (RepositoryForEnt) Fields() []ent.Field {
	return nil
}

func (RepositoryForEnt) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user_for_ents", UserForEnt.Type),
		edge.To("issue_for_ents", IssueForEnt.Type),
	}
}
