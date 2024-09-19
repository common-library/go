// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/common-library/go/database/orm/ent/issueforent"
	"github.com/common-library/go/database/orm/ent/repositoryforent"
	"github.com/common-library/go/database/orm/ent/userforent"
)

// RepositoryForEntCreate is the builder for creating a RepositoryForEnt entity.
type RepositoryForEntCreate struct {
	config
	mutation *RepositoryForEntMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// AddUserForEntIDs adds the "user_for_ents" edge to the UserForEnt entity by IDs.
func (rfec *RepositoryForEntCreate) AddUserForEntIDs(ids ...int) *RepositoryForEntCreate {
	rfec.mutation.AddUserForEntIDs(ids...)
	return rfec
}

// AddUserForEnts adds the "user_for_ents" edges to the UserForEnt entity.
func (rfec *RepositoryForEntCreate) AddUserForEnts(u ...*UserForEnt) *RepositoryForEntCreate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return rfec.AddUserForEntIDs(ids...)
}

// AddIssueForEntIDs adds the "issue_for_ents" edge to the IssueForEnt entity by IDs.
func (rfec *RepositoryForEntCreate) AddIssueForEntIDs(ids ...int) *RepositoryForEntCreate {
	rfec.mutation.AddIssueForEntIDs(ids...)
	return rfec
}

// AddIssueForEnts adds the "issue_for_ents" edges to the IssueForEnt entity.
func (rfec *RepositoryForEntCreate) AddIssueForEnts(i ...*IssueForEnt) *RepositoryForEntCreate {
	ids := make([]int, len(i))
	for j := range i {
		ids[j] = i[j].ID
	}
	return rfec.AddIssueForEntIDs(ids...)
}

// Mutation returns the RepositoryForEntMutation object of the builder.
func (rfec *RepositoryForEntCreate) Mutation() *RepositoryForEntMutation {
	return rfec.mutation
}

// Save creates the RepositoryForEnt in the database.
func (rfec *RepositoryForEntCreate) Save(ctx context.Context) (*RepositoryForEnt, error) {
	return withHooks(ctx, rfec.sqlSave, rfec.mutation, rfec.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (rfec *RepositoryForEntCreate) SaveX(ctx context.Context) *RepositoryForEnt {
	v, err := rfec.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rfec *RepositoryForEntCreate) Exec(ctx context.Context) error {
	_, err := rfec.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rfec *RepositoryForEntCreate) ExecX(ctx context.Context) {
	if err := rfec.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (rfec *RepositoryForEntCreate) check() error {
	return nil
}

func (rfec *RepositoryForEntCreate) sqlSave(ctx context.Context) (*RepositoryForEnt, error) {
	if err := rfec.check(); err != nil {
		return nil, err
	}
	_node, _spec := rfec.createSpec()
	if err := sqlgraph.CreateNode(ctx, rfec.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	rfec.mutation.id = &_node.ID
	rfec.mutation.done = true
	return _node, nil
}

func (rfec *RepositoryForEntCreate) createSpec() (*RepositoryForEnt, *sqlgraph.CreateSpec) {
	var (
		_node = &RepositoryForEnt{config: rfec.config}
		_spec = sqlgraph.NewCreateSpec(repositoryforent.Table, sqlgraph.NewFieldSpec(repositoryforent.FieldID, field.TypeInt))
	)
	_spec.OnConflict = rfec.conflict
	if nodes := rfec.mutation.UserForEntsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   repositoryforent.UserForEntsTable,
			Columns: repositoryforent.UserForEntsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(userforent.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rfec.mutation.IssueForEntsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   repositoryforent.IssueForEntsTable,
			Columns: []string{repositoryforent.IssueForEntsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(issueforent.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.RepositoryForEnt.Create().
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (rfec *RepositoryForEntCreate) OnConflict(opts ...sql.ConflictOption) *RepositoryForEntUpsertOne {
	rfec.conflict = opts
	return &RepositoryForEntUpsertOne{
		create: rfec,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.RepositoryForEnt.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (rfec *RepositoryForEntCreate) OnConflictColumns(columns ...string) *RepositoryForEntUpsertOne {
	rfec.conflict = append(rfec.conflict, sql.ConflictColumns(columns...))
	return &RepositoryForEntUpsertOne{
		create: rfec,
	}
}

type (
	// RepositoryForEntUpsertOne is the builder for "upsert"-ing
	//  one RepositoryForEnt node.
	RepositoryForEntUpsertOne struct {
		create *RepositoryForEntCreate
	}

	// RepositoryForEntUpsert is the "OnConflict" setter.
	RepositoryForEntUpsert struct {
		*sql.UpdateSet
	}
)

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.RepositoryForEnt.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *RepositoryForEntUpsertOne) UpdateNewValues() *RepositoryForEntUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.RepositoryForEnt.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *RepositoryForEntUpsertOne) Ignore() *RepositoryForEntUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *RepositoryForEntUpsertOne) DoNothing() *RepositoryForEntUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the RepositoryForEntCreate.OnConflict
// documentation for more info.
func (u *RepositoryForEntUpsertOne) Update(set func(*RepositoryForEntUpsert)) *RepositoryForEntUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&RepositoryForEntUpsert{UpdateSet: update})
	}))
	return u
}

// Exec executes the query.
func (u *RepositoryForEntUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for RepositoryForEntCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *RepositoryForEntUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *RepositoryForEntUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *RepositoryForEntUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// RepositoryForEntCreateBulk is the builder for creating many RepositoryForEnt entities in bulk.
type RepositoryForEntCreateBulk struct {
	config
	err      error
	builders []*RepositoryForEntCreate
	conflict []sql.ConflictOption
}

// Save creates the RepositoryForEnt entities in the database.
func (rfecb *RepositoryForEntCreateBulk) Save(ctx context.Context) ([]*RepositoryForEnt, error) {
	if rfecb.err != nil {
		return nil, rfecb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(rfecb.builders))
	nodes := make([]*RepositoryForEnt, len(rfecb.builders))
	mutators := make([]Mutator, len(rfecb.builders))
	for i := range rfecb.builders {
		func(i int, root context.Context) {
			builder := rfecb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*RepositoryForEntMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, rfecb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = rfecb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, rfecb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, rfecb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (rfecb *RepositoryForEntCreateBulk) SaveX(ctx context.Context) []*RepositoryForEnt {
	v, err := rfecb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rfecb *RepositoryForEntCreateBulk) Exec(ctx context.Context) error {
	_, err := rfecb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rfecb *RepositoryForEntCreateBulk) ExecX(ctx context.Context) {
	if err := rfecb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.RepositoryForEnt.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (rfecb *RepositoryForEntCreateBulk) OnConflict(opts ...sql.ConflictOption) *RepositoryForEntUpsertBulk {
	rfecb.conflict = opts
	return &RepositoryForEntUpsertBulk{
		create: rfecb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.RepositoryForEnt.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (rfecb *RepositoryForEntCreateBulk) OnConflictColumns(columns ...string) *RepositoryForEntUpsertBulk {
	rfecb.conflict = append(rfecb.conflict, sql.ConflictColumns(columns...))
	return &RepositoryForEntUpsertBulk{
		create: rfecb,
	}
}

// RepositoryForEntUpsertBulk is the builder for "upsert"-ing
// a bulk of RepositoryForEnt nodes.
type RepositoryForEntUpsertBulk struct {
	create *RepositoryForEntCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.RepositoryForEnt.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *RepositoryForEntUpsertBulk) UpdateNewValues() *RepositoryForEntUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.RepositoryForEnt.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *RepositoryForEntUpsertBulk) Ignore() *RepositoryForEntUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *RepositoryForEntUpsertBulk) DoNothing() *RepositoryForEntUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the RepositoryForEntCreateBulk.OnConflict
// documentation for more info.
func (u *RepositoryForEntUpsertBulk) Update(set func(*RepositoryForEntUpsert)) *RepositoryForEntUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&RepositoryForEntUpsert{UpdateSet: update})
	}))
	return u
}

// Exec executes the query.
func (u *RepositoryForEntUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the RepositoryForEntCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for RepositoryForEntCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *RepositoryForEntUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
