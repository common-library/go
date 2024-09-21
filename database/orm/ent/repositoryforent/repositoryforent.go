// Code generated by ent, DO NOT EDIT.

package repositoryforent

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the repositoryforent type in the database.
	Label = "repository_for_ent"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// EdgeUserForEnts holds the string denoting the user_for_ents edge name in mutations.
	EdgeUserForEnts = "user_for_ents"
	// EdgeIssueForEnts holds the string denoting the issue_for_ents edge name in mutations.
	EdgeIssueForEnts = "issue_for_ents"
	// Table holds the table name of the repositoryforent in the database.
	Table = "repository_for_ents"
	// UserForEntsTable is the table that holds the user_for_ents relation/edge. The primary key declared below.
	UserForEntsTable = "repository_for_ent_user_for_ents"
	// UserForEntsInverseTable is the table name for the UserForEnt entity.
	// It exists in this package in order to avoid circular dependency with the "userforent" package.
	UserForEntsInverseTable = "user_for_ents"
	// IssueForEntsTable is the table that holds the issue_for_ents relation/edge.
	IssueForEntsTable = "issue_for_ents"
	// IssueForEntsInverseTable is the table name for the IssueForEnt entity.
	// It exists in this package in order to avoid circular dependency with the "issueforent" package.
	IssueForEntsInverseTable = "issue_for_ents"
	// IssueForEntsColumn is the table column denoting the issue_for_ents relation/edge.
	IssueForEntsColumn = "repository_for_ent_issue_for_ents"
)

// Columns holds all SQL columns for repositoryforent fields.
var Columns = []string{
	FieldID,
}

var (
	// UserForEntsPrimaryKey and UserForEntsColumn2 are the table columns denoting the
	// primary key for the user_for_ents relation (M2M).
	UserForEntsPrimaryKey = []string{"repository_for_ent_id", "user_for_ent_id"}
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

// OrderOption defines the ordering options for the RepositoryForEnt queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByUserForEntsCount orders the results by user_for_ents count.
func ByUserForEntsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newUserForEntsStep(), opts...)
	}
}

// ByUserForEnts orders the results by user_for_ents terms.
func ByUserForEnts(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newUserForEntsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByIssueForEntsCount orders the results by issue_for_ents count.
func ByIssueForEntsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newIssueForEntsStep(), opts...)
	}
}

// ByIssueForEnts orders the results by issue_for_ents terms.
func ByIssueForEnts(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newIssueForEntsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newUserForEntsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(UserForEntsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, UserForEntsTable, UserForEntsPrimaryKey...),
	)
}
func newIssueForEntsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(IssueForEntsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, IssueForEntsTable, IssueForEntsColumn),
	)
}
