// Code generated by ent, DO NOT EDIT.

package userforent

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the userforent type in the database.
	Label = "user_for_ent"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// EdgeRepositoryForEnts holds the string denoting the repository_for_ents edge name in mutations.
	EdgeRepositoryForEnts = "repository_for_ents"
	// Table holds the table name of the userforent in the database.
	Table = "user_for_ents"
	// RepositoryForEntsTable is the table that holds the repository_for_ents relation/edge. The primary key declared below.
	RepositoryForEntsTable = "repository_for_ent_user_for_ents"
	// RepositoryForEntsInverseTable is the table name for the RepositoryForEnt entity.
	// It exists in this package in order to avoid circular dependency with the "repositoryforent" package.
	RepositoryForEntsInverseTable = "repository_for_ents"
)

// Columns holds all SQL columns for userforent fields.
var Columns = []string{
	FieldID,
}

var (
	// RepositoryForEntsPrimaryKey and RepositoryForEntsColumn2 are the table columns denoting the
	// primary key for the repository_for_ents relation (M2M).
	RepositoryForEntsPrimaryKey = []string{"repository_for_ent_id", "user_for_ent_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the UserForEnt queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByRepositoryForEntsCount orders the results by repository_for_ents count.
func ByRepositoryForEntsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newRepositoryForEntsStep(), opts...)
	}
}

// ByRepositoryForEnts orders the results by repository_for_ents terms.
func ByRepositoryForEnts(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newRepositoryForEntsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newRepositoryForEntsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(RepositoryForEntsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, RepositoryForEntsTable, RepositoryForEntsPrimaryKey...),
	)
}
