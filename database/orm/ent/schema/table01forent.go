package schema

import (
	"errors"
	"math/rand/v2"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

type Table01ForEnt struct {
	ent.Schema
}

func (Table01ForEnt) Mixin() []ent.Mixin {
	return []ent.Mixin{
		CommonMixin{},

		mixin.Time{},
	}
}

func (Table01ForEnt) Fields() []ent.Field {
	return []ent.Field{
		field.String("field01").Unique(),
		field.Int("field02"),
		field.Bool("field03").Default(false),
		field.Enum("field04").Values("value01", "value02").Default("value01"),
		field.String("field05").Optional(),
		field.String("field06").Optional().Nillable(),
		field.Int64("field07").DefaultFunc(func() int64 {
			return rand.Int64()
		}).UpdateDefault(rand.Int64),
		field.Float("field08").Default(1).SchemaType(map[string]string{
			dialect.MySQL:    "decimal(6,2)",
			dialect.Postgres: "numeric",
		}),
		field.String("field09").Default("Field").Validate(func(s string) error {
			if strings.ToLower(s) == s {
				return errors.New("invalid")
			} else {
				return nil
			}
		}),

		field.String("field10").Optional().Sensitive(),
		field.String("field11").Optional().Comment("comment"),
		field.String("field12").Optional().Comment("comment").Annotations(entsql.WithComments(true)),
		field.String("field13").Optional().Deprecated("deprecated"),
		field.String("field14").Optional().StorageKey("storage_key"),
		field.String("field15").Optional().StructTag(`gqlgen:"gql_name"`),
	}
}

func (Table01ForEnt) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("field01", "field07").Unique(),
		index.Fields("field02", "field03"),
	}
}

func (Table01ForEnt) Edges() []ent.Edge {
	return nil
}
