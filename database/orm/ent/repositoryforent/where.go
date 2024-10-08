// Code generated by ent, DO NOT EDIT.

package repositoryforent

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/common-library/go/database/orm/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.RepositoryForEnt {
	return predicate.RepositoryForEnt(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.RepositoryForEnt {
	return predicate.RepositoryForEnt(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.RepositoryForEnt {
	return predicate.RepositoryForEnt(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.RepositoryForEnt {
	return predicate.RepositoryForEnt(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.RepositoryForEnt {
	return predicate.RepositoryForEnt(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.RepositoryForEnt {
	return predicate.RepositoryForEnt(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.RepositoryForEnt {
	return predicate.RepositoryForEnt(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.RepositoryForEnt {
	return predicate.RepositoryForEnt(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.RepositoryForEnt {
	return predicate.RepositoryForEnt(sql.FieldLTE(FieldID, id))
}

// HasUserForEnts applies the HasEdge predicate on the "user_for_ents" edge.
func HasUserForEnts() predicate.RepositoryForEnt {
	return predicate.RepositoryForEnt(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, UserForEntsTable, UserForEntsPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasUserForEntsWith applies the HasEdge predicate on the "user_for_ents" edge with a given conditions (other predicates).
func HasUserForEntsWith(preds ...predicate.UserForEnt) predicate.RepositoryForEnt {
	return predicate.RepositoryForEnt(func(s *sql.Selector) {
		step := newUserForEntsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasIssueForEnts applies the HasEdge predicate on the "issue_for_ents" edge.
func HasIssueForEnts() predicate.RepositoryForEnt {
	return predicate.RepositoryForEnt(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, IssueForEntsTable, IssueForEntsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasIssueForEntsWith applies the HasEdge predicate on the "issue_for_ents" edge with a given conditions (other predicates).
func HasIssueForEntsWith(preds ...predicate.IssueForEnt) predicate.RepositoryForEnt {
	return predicate.RepositoryForEnt(func(s *sql.Selector) {
		step := newIssueForEntsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.RepositoryForEnt) predicate.RepositoryForEnt {
	return predicate.RepositoryForEnt(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.RepositoryForEnt) predicate.RepositoryForEnt {
	return predicate.RepositoryForEnt(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.RepositoryForEnt) predicate.RepositoryForEnt {
	return predicate.RepositoryForEnt(sql.NotPredicates(p))
}
