// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/common-library/go/database/orm/ent/issueforent"
	"github.com/common-library/go/database/orm/ent/predicate"
	"github.com/common-library/go/database/orm/ent/repositoryforent"
	"github.com/common-library/go/database/orm/ent/userforent"
)

// RepositoryForEntQuery is the builder for querying RepositoryForEnt entities.
type RepositoryForEntQuery struct {
	config
	ctx              *QueryContext
	order            []repositoryforent.OrderOption
	inters           []Interceptor
	predicates       []predicate.RepositoryForEnt
	withUserForEnts  *UserForEntQuery
	withIssueForEnts *IssueForEntQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the RepositoryForEntQuery builder.
func (rfeq *RepositoryForEntQuery) Where(ps ...predicate.RepositoryForEnt) *RepositoryForEntQuery {
	rfeq.predicates = append(rfeq.predicates, ps...)
	return rfeq
}

// Limit the number of records to be returned by this query.
func (rfeq *RepositoryForEntQuery) Limit(limit int) *RepositoryForEntQuery {
	rfeq.ctx.Limit = &limit
	return rfeq
}

// Offset to start from.
func (rfeq *RepositoryForEntQuery) Offset(offset int) *RepositoryForEntQuery {
	rfeq.ctx.Offset = &offset
	return rfeq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (rfeq *RepositoryForEntQuery) Unique(unique bool) *RepositoryForEntQuery {
	rfeq.ctx.Unique = &unique
	return rfeq
}

// Order specifies how the records should be ordered.
func (rfeq *RepositoryForEntQuery) Order(o ...repositoryforent.OrderOption) *RepositoryForEntQuery {
	rfeq.order = append(rfeq.order, o...)
	return rfeq
}

// QueryUserForEnts chains the current query on the "user_for_ents" edge.
func (rfeq *RepositoryForEntQuery) QueryUserForEnts() *UserForEntQuery {
	query := (&UserForEntClient{config: rfeq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rfeq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rfeq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(repositoryforent.Table, repositoryforent.FieldID, selector),
			sqlgraph.To(userforent.Table, userforent.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, repositoryforent.UserForEntsTable, repositoryforent.UserForEntsPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(rfeq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryIssueForEnts chains the current query on the "issue_for_ents" edge.
func (rfeq *RepositoryForEntQuery) QueryIssueForEnts() *IssueForEntQuery {
	query := (&IssueForEntClient{config: rfeq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rfeq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rfeq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(repositoryforent.Table, repositoryforent.FieldID, selector),
			sqlgraph.To(issueforent.Table, issueforent.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, repositoryforent.IssueForEntsTable, repositoryforent.IssueForEntsColumn),
		)
		fromU = sqlgraph.SetNeighbors(rfeq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first RepositoryForEnt entity from the query.
// Returns a *NotFoundError when no RepositoryForEnt was found.
func (rfeq *RepositoryForEntQuery) First(ctx context.Context) (*RepositoryForEnt, error) {
	nodes, err := rfeq.Limit(1).All(setContextOp(ctx, rfeq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{repositoryforent.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (rfeq *RepositoryForEntQuery) FirstX(ctx context.Context) *RepositoryForEnt {
	node, err := rfeq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first RepositoryForEnt ID from the query.
// Returns a *NotFoundError when no RepositoryForEnt ID was found.
func (rfeq *RepositoryForEntQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = rfeq.Limit(1).IDs(setContextOp(ctx, rfeq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{repositoryforent.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (rfeq *RepositoryForEntQuery) FirstIDX(ctx context.Context) int {
	id, err := rfeq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single RepositoryForEnt entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one RepositoryForEnt entity is found.
// Returns a *NotFoundError when no RepositoryForEnt entities are found.
func (rfeq *RepositoryForEntQuery) Only(ctx context.Context) (*RepositoryForEnt, error) {
	nodes, err := rfeq.Limit(2).All(setContextOp(ctx, rfeq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{repositoryforent.Label}
	default:
		return nil, &NotSingularError{repositoryforent.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (rfeq *RepositoryForEntQuery) OnlyX(ctx context.Context) *RepositoryForEnt {
	node, err := rfeq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only RepositoryForEnt ID in the query.
// Returns a *NotSingularError when more than one RepositoryForEnt ID is found.
// Returns a *NotFoundError when no entities are found.
func (rfeq *RepositoryForEntQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = rfeq.Limit(2).IDs(setContextOp(ctx, rfeq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{repositoryforent.Label}
	default:
		err = &NotSingularError{repositoryforent.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (rfeq *RepositoryForEntQuery) OnlyIDX(ctx context.Context) int {
	id, err := rfeq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of RepositoryForEnts.
func (rfeq *RepositoryForEntQuery) All(ctx context.Context) ([]*RepositoryForEnt, error) {
	ctx = setContextOp(ctx, rfeq.ctx, ent.OpQueryAll)
	if err := rfeq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*RepositoryForEnt, *RepositoryForEntQuery]()
	return withInterceptors[[]*RepositoryForEnt](ctx, rfeq, qr, rfeq.inters)
}

// AllX is like All, but panics if an error occurs.
func (rfeq *RepositoryForEntQuery) AllX(ctx context.Context) []*RepositoryForEnt {
	nodes, err := rfeq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of RepositoryForEnt IDs.
func (rfeq *RepositoryForEntQuery) IDs(ctx context.Context) (ids []int, err error) {
	if rfeq.ctx.Unique == nil && rfeq.path != nil {
		rfeq.Unique(true)
	}
	ctx = setContextOp(ctx, rfeq.ctx, ent.OpQueryIDs)
	if err = rfeq.Select(repositoryforent.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (rfeq *RepositoryForEntQuery) IDsX(ctx context.Context) []int {
	ids, err := rfeq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (rfeq *RepositoryForEntQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, rfeq.ctx, ent.OpQueryCount)
	if err := rfeq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, rfeq, querierCount[*RepositoryForEntQuery](), rfeq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (rfeq *RepositoryForEntQuery) CountX(ctx context.Context) int {
	count, err := rfeq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (rfeq *RepositoryForEntQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, rfeq.ctx, ent.OpQueryExist)
	switch _, err := rfeq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (rfeq *RepositoryForEntQuery) ExistX(ctx context.Context) bool {
	exist, err := rfeq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the RepositoryForEntQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (rfeq *RepositoryForEntQuery) Clone() *RepositoryForEntQuery {
	if rfeq == nil {
		return nil
	}
	return &RepositoryForEntQuery{
		config:           rfeq.config,
		ctx:              rfeq.ctx.Clone(),
		order:            append([]repositoryforent.OrderOption{}, rfeq.order...),
		inters:           append([]Interceptor{}, rfeq.inters...),
		predicates:       append([]predicate.RepositoryForEnt{}, rfeq.predicates...),
		withUserForEnts:  rfeq.withUserForEnts.Clone(),
		withIssueForEnts: rfeq.withIssueForEnts.Clone(),
		// clone intermediate query.
		sql:  rfeq.sql.Clone(),
		path: rfeq.path,
	}
}

// WithUserForEnts tells the query-builder to eager-load the nodes that are connected to
// the "user_for_ents" edge. The optional arguments are used to configure the query builder of the edge.
func (rfeq *RepositoryForEntQuery) WithUserForEnts(opts ...func(*UserForEntQuery)) *RepositoryForEntQuery {
	query := (&UserForEntClient{config: rfeq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	rfeq.withUserForEnts = query
	return rfeq
}

// WithIssueForEnts tells the query-builder to eager-load the nodes that are connected to
// the "issue_for_ents" edge. The optional arguments are used to configure the query builder of the edge.
func (rfeq *RepositoryForEntQuery) WithIssueForEnts(opts ...func(*IssueForEntQuery)) *RepositoryForEntQuery {
	query := (&IssueForEntClient{config: rfeq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	rfeq.withIssueForEnts = query
	return rfeq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
func (rfeq *RepositoryForEntQuery) GroupBy(field string, fields ...string) *RepositoryForEntGroupBy {
	rfeq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &RepositoryForEntGroupBy{build: rfeq}
	grbuild.flds = &rfeq.ctx.Fields
	grbuild.label = repositoryforent.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
func (rfeq *RepositoryForEntQuery) Select(fields ...string) *RepositoryForEntSelect {
	rfeq.ctx.Fields = append(rfeq.ctx.Fields, fields...)
	sbuild := &RepositoryForEntSelect{RepositoryForEntQuery: rfeq}
	sbuild.label = repositoryforent.Label
	sbuild.flds, sbuild.scan = &rfeq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a RepositoryForEntSelect configured with the given aggregations.
func (rfeq *RepositoryForEntQuery) Aggregate(fns ...AggregateFunc) *RepositoryForEntSelect {
	return rfeq.Select().Aggregate(fns...)
}

func (rfeq *RepositoryForEntQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range rfeq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, rfeq); err != nil {
				return err
			}
		}
	}
	for _, f := range rfeq.ctx.Fields {
		if !repositoryforent.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if rfeq.path != nil {
		prev, err := rfeq.path(ctx)
		if err != nil {
			return err
		}
		rfeq.sql = prev
	}
	return nil
}

func (rfeq *RepositoryForEntQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*RepositoryForEnt, error) {
	var (
		nodes       = []*RepositoryForEnt{}
		_spec       = rfeq.querySpec()
		loadedTypes = [2]bool{
			rfeq.withUserForEnts != nil,
			rfeq.withIssueForEnts != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*RepositoryForEnt).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &RepositoryForEnt{config: rfeq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, rfeq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := rfeq.withUserForEnts; query != nil {
		if err := rfeq.loadUserForEnts(ctx, query, nodes,
			func(n *RepositoryForEnt) { n.Edges.UserForEnts = []*UserForEnt{} },
			func(n *RepositoryForEnt, e *UserForEnt) { n.Edges.UserForEnts = append(n.Edges.UserForEnts, e) }); err != nil {
			return nil, err
		}
	}
	if query := rfeq.withIssueForEnts; query != nil {
		if err := rfeq.loadIssueForEnts(ctx, query, nodes,
			func(n *RepositoryForEnt) { n.Edges.IssueForEnts = []*IssueForEnt{} },
			func(n *RepositoryForEnt, e *IssueForEnt) { n.Edges.IssueForEnts = append(n.Edges.IssueForEnts, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (rfeq *RepositoryForEntQuery) loadUserForEnts(ctx context.Context, query *UserForEntQuery, nodes []*RepositoryForEnt, init func(*RepositoryForEnt), assign func(*RepositoryForEnt, *UserForEnt)) error {
	edgeIDs := make([]driver.Value, len(nodes))
	byID := make(map[int]*RepositoryForEnt)
	nids := make(map[int]map[*RepositoryForEnt]struct{})
	for i, node := range nodes {
		edgeIDs[i] = node.ID
		byID[node.ID] = node
		if init != nil {
			init(node)
		}
	}
	query.Where(func(s *sql.Selector) {
		joinT := sql.Table(repositoryforent.UserForEntsTable)
		s.Join(joinT).On(s.C(userforent.FieldID), joinT.C(repositoryforent.UserForEntsPrimaryKey[1]))
		s.Where(sql.InValues(joinT.C(repositoryforent.UserForEntsPrimaryKey[0]), edgeIDs...))
		columns := s.SelectedColumns()
		s.Select(joinT.C(repositoryforent.UserForEntsPrimaryKey[0]))
		s.AppendSelect(columns...)
		s.SetDistinct(false)
	})
	if err := query.prepareQuery(ctx); err != nil {
		return err
	}
	qr := QuerierFunc(func(ctx context.Context, q Query) (Value, error) {
		return query.sqlAll(ctx, func(_ context.Context, spec *sqlgraph.QuerySpec) {
			assign := spec.Assign
			values := spec.ScanValues
			spec.ScanValues = func(columns []string) ([]any, error) {
				values, err := values(columns[1:])
				if err != nil {
					return nil, err
				}
				return append([]any{new(sql.NullInt64)}, values...), nil
			}
			spec.Assign = func(columns []string, values []any) error {
				outValue := int(values[0].(*sql.NullInt64).Int64)
				inValue := int(values[1].(*sql.NullInt64).Int64)
				if nids[inValue] == nil {
					nids[inValue] = map[*RepositoryForEnt]struct{}{byID[outValue]: {}}
					return assign(columns[1:], values[1:])
				}
				nids[inValue][byID[outValue]] = struct{}{}
				return nil
			}
		})
	})
	neighbors, err := withInterceptors[[]*UserForEnt](ctx, query, qr, query.inters)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected "user_for_ents" node returned %v`, n.ID)
		}
		for kn := range nodes {
			assign(kn, n)
		}
	}
	return nil
}
func (rfeq *RepositoryForEntQuery) loadIssueForEnts(ctx context.Context, query *IssueForEntQuery, nodes []*RepositoryForEnt, init func(*RepositoryForEnt), assign func(*RepositoryForEnt, *IssueForEnt)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*RepositoryForEnt)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.IssueForEnt(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(repositoryforent.IssueForEntsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.repository_for_ent_issue_for_ents
		if fk == nil {
			return fmt.Errorf(`foreign-key "repository_for_ent_issue_for_ents" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "repository_for_ent_issue_for_ents" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (rfeq *RepositoryForEntQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := rfeq.querySpec()
	_spec.Node.Columns = rfeq.ctx.Fields
	if len(rfeq.ctx.Fields) > 0 {
		_spec.Unique = rfeq.ctx.Unique != nil && *rfeq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, rfeq.driver, _spec)
}

func (rfeq *RepositoryForEntQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(repositoryforent.Table, repositoryforent.Columns, sqlgraph.NewFieldSpec(repositoryforent.FieldID, field.TypeInt))
	_spec.From = rfeq.sql
	if unique := rfeq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if rfeq.path != nil {
		_spec.Unique = true
	}
	if fields := rfeq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, repositoryforent.FieldID)
		for i := range fields {
			if fields[i] != repositoryforent.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := rfeq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := rfeq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := rfeq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := rfeq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (rfeq *RepositoryForEntQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(rfeq.driver.Dialect())
	t1 := builder.Table(repositoryforent.Table)
	columns := rfeq.ctx.Fields
	if len(columns) == 0 {
		columns = repositoryforent.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if rfeq.sql != nil {
		selector = rfeq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if rfeq.ctx.Unique != nil && *rfeq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range rfeq.predicates {
		p(selector)
	}
	for _, p := range rfeq.order {
		p(selector)
	}
	if offset := rfeq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := rfeq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// RepositoryForEntGroupBy is the group-by builder for RepositoryForEnt entities.
type RepositoryForEntGroupBy struct {
	selector
	build *RepositoryForEntQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (rfegb *RepositoryForEntGroupBy) Aggregate(fns ...AggregateFunc) *RepositoryForEntGroupBy {
	rfegb.fns = append(rfegb.fns, fns...)
	return rfegb
}

// Scan applies the selector query and scans the result into the given value.
func (rfegb *RepositoryForEntGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, rfegb.build.ctx, ent.OpQueryGroupBy)
	if err := rfegb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*RepositoryForEntQuery, *RepositoryForEntGroupBy](ctx, rfegb.build, rfegb, rfegb.build.inters, v)
}

func (rfegb *RepositoryForEntGroupBy) sqlScan(ctx context.Context, root *RepositoryForEntQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(rfegb.fns))
	for _, fn := range rfegb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*rfegb.flds)+len(rfegb.fns))
		for _, f := range *rfegb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*rfegb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := rfegb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// RepositoryForEntSelect is the builder for selecting fields of RepositoryForEnt entities.
type RepositoryForEntSelect struct {
	*RepositoryForEntQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (rfes *RepositoryForEntSelect) Aggregate(fns ...AggregateFunc) *RepositoryForEntSelect {
	rfes.fns = append(rfes.fns, fns...)
	return rfes
}

// Scan applies the selector query and scans the result into the given value.
func (rfes *RepositoryForEntSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, rfes.ctx, ent.OpQuerySelect)
	if err := rfes.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*RepositoryForEntQuery, *RepositoryForEntSelect](ctx, rfes.RepositoryForEntQuery, rfes, rfes.inters, v)
}

func (rfes *RepositoryForEntSelect) sqlScan(ctx context.Context, root *RepositoryForEntQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(rfes.fns))
	for _, fn := range rfes.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*rfes.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := rfes.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
