package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
)

type IssueForEnt struct {
	ent.Schema
}

func (IssueForEnt) Fields() []ent.Field {
	return nil
}

func (IssueForEnt) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("repository", RepositoryForEnt.Type).Ref("issue_for_ents").Unique(),
	}
}
