package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type CommonMixin struct {
	mixin.Schema
}

func (CommonMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Int("common_field01").Optional(),
	}
}

func (CommonMixin) Edges() []ent.Edge {
	return nil
}
