// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/common-library/go/database/orm/ent/predicate"
	"github.com/common-library/go/database/orm/ent/table01forent"
)

// Table01ForEntDelete is the builder for deleting a Table01ForEnt entity.
type Table01ForEntDelete struct {
	config
	hooks    []Hook
	mutation *Table01ForEntMutation
}

// Where appends a list predicates to the Table01ForEntDelete builder.
func (ted *Table01ForEntDelete) Where(ps ...predicate.Table01ForEnt) *Table01ForEntDelete {
	ted.mutation.Where(ps...)
	return ted
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ted *Table01ForEntDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, ted.sqlExec, ted.mutation, ted.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (ted *Table01ForEntDelete) ExecX(ctx context.Context) int {
	n, err := ted.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ted *Table01ForEntDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(table01forent.Table, sqlgraph.NewFieldSpec(table01forent.FieldID, field.TypeInt))
	if ps := ted.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, ted.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	ted.mutation.done = true
	return affected, err
}

// Table01ForEntDeleteOne is the builder for deleting a single Table01ForEnt entity.
type Table01ForEntDeleteOne struct {
	ted *Table01ForEntDelete
}

// Where appends a list predicates to the Table01ForEntDelete builder.
func (tedo *Table01ForEntDeleteOne) Where(ps ...predicate.Table01ForEnt) *Table01ForEntDeleteOne {
	tedo.ted.mutation.Where(ps...)
	return tedo
}

// Exec executes the deletion query.
func (tedo *Table01ForEntDeleteOne) Exec(ctx context.Context) error {
	n, err := tedo.ted.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{table01forent.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (tedo *Table01ForEntDeleteOne) ExecX(ctx context.Context) {
	if err := tedo.Exec(ctx); err != nil {
		panic(err)
	}
}
