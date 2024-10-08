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
	"github.com/common-library/go/database/orm/ent/predicate"
	"github.com/common-library/go/database/orm/ent/repositoryforent"
	"github.com/common-library/go/database/orm/ent/userforent"
)

// UserForEntQuery is the builder for querying UserForEnt entities.
type UserForEntQuery struct {
	config
	ctx                   *QueryContext
	order                 []userforent.OrderOption
	inters                []Interceptor
	predicates            []predicate.UserForEnt
	withRepositoryForEnts *RepositoryForEntQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the UserForEntQuery builder.
func (ufeq *UserForEntQuery) Where(ps ...predicate.UserForEnt) *UserForEntQuery {
	ufeq.predicates = append(ufeq.predicates, ps...)
	return ufeq
}

// Limit the number of records to be returned by this query.
func (ufeq *UserForEntQuery) Limit(limit int) *UserForEntQuery {
	ufeq.ctx.Limit = &limit
	return ufeq
}

// Offset to start from.
func (ufeq *UserForEntQuery) Offset(offset int) *UserForEntQuery {
	ufeq.ctx.Offset = &offset
	return ufeq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (ufeq *UserForEntQuery) Unique(unique bool) *UserForEntQuery {
	ufeq.ctx.Unique = &unique
	return ufeq
}

// Order specifies how the records should be ordered.
func (ufeq *UserForEntQuery) Order(o ...userforent.OrderOption) *UserForEntQuery {
	ufeq.order = append(ufeq.order, o...)
	return ufeq
}

// QueryRepositoryForEnts chains the current query on the "repository_for_ents" edge.
func (ufeq *UserForEntQuery) QueryRepositoryForEnts() *RepositoryForEntQuery {
	query := (&RepositoryForEntClient{config: ufeq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := ufeq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := ufeq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(userforent.Table, userforent.FieldID, selector),
			sqlgraph.To(repositoryforent.Table, repositoryforent.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, userforent.RepositoryForEntsTable, userforent.RepositoryForEntsPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(ufeq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first UserForEnt entity from the query.
// Returns a *NotFoundError when no UserForEnt was found.
func (ufeq *UserForEntQuery) First(ctx context.Context) (*UserForEnt, error) {
	nodes, err := ufeq.Limit(1).All(setContextOp(ctx, ufeq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{userforent.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (ufeq *UserForEntQuery) FirstX(ctx context.Context) *UserForEnt {
	node, err := ufeq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first UserForEnt ID from the query.
// Returns a *NotFoundError when no UserForEnt ID was found.
func (ufeq *UserForEntQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = ufeq.Limit(1).IDs(setContextOp(ctx, ufeq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{userforent.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (ufeq *UserForEntQuery) FirstIDX(ctx context.Context) int {
	id, err := ufeq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single UserForEnt entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one UserForEnt entity is found.
// Returns a *NotFoundError when no UserForEnt entities are found.
func (ufeq *UserForEntQuery) Only(ctx context.Context) (*UserForEnt, error) {
	nodes, err := ufeq.Limit(2).All(setContextOp(ctx, ufeq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{userforent.Label}
	default:
		return nil, &NotSingularError{userforent.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (ufeq *UserForEntQuery) OnlyX(ctx context.Context) *UserForEnt {
	node, err := ufeq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only UserForEnt ID in the query.
// Returns a *NotSingularError when more than one UserForEnt ID is found.
// Returns a *NotFoundError when no entities are found.
func (ufeq *UserForEntQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = ufeq.Limit(2).IDs(setContextOp(ctx, ufeq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{userforent.Label}
	default:
		err = &NotSingularError{userforent.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (ufeq *UserForEntQuery) OnlyIDX(ctx context.Context) int {
	id, err := ufeq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of UserForEnts.
func (ufeq *UserForEntQuery) All(ctx context.Context) ([]*UserForEnt, error) {
	ctx = setContextOp(ctx, ufeq.ctx, ent.OpQueryAll)
	if err := ufeq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*UserForEnt, *UserForEntQuery]()
	return withInterceptors[[]*UserForEnt](ctx, ufeq, qr, ufeq.inters)
}

// AllX is like All, but panics if an error occurs.
func (ufeq *UserForEntQuery) AllX(ctx context.Context) []*UserForEnt {
	nodes, err := ufeq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of UserForEnt IDs.
func (ufeq *UserForEntQuery) IDs(ctx context.Context) (ids []int, err error) {
	if ufeq.ctx.Unique == nil && ufeq.path != nil {
		ufeq.Unique(true)
	}
	ctx = setContextOp(ctx, ufeq.ctx, ent.OpQueryIDs)
	if err = ufeq.Select(userforent.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (ufeq *UserForEntQuery) IDsX(ctx context.Context) []int {
	ids, err := ufeq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (ufeq *UserForEntQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, ufeq.ctx, ent.OpQueryCount)
	if err := ufeq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, ufeq, querierCount[*UserForEntQuery](), ufeq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (ufeq *UserForEntQuery) CountX(ctx context.Context) int {
	count, err := ufeq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (ufeq *UserForEntQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, ufeq.ctx, ent.OpQueryExist)
	switch _, err := ufeq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (ufeq *UserForEntQuery) ExistX(ctx context.Context) bool {
	exist, err := ufeq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the UserForEntQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (ufeq *UserForEntQuery) Clone() *UserForEntQuery {
	if ufeq == nil {
		return nil
	}
	return &UserForEntQuery{
		config:                ufeq.config,
		ctx:                   ufeq.ctx.Clone(),
		order:                 append([]userforent.OrderOption{}, ufeq.order...),
		inters:                append([]Interceptor{}, ufeq.inters...),
		predicates:            append([]predicate.UserForEnt{}, ufeq.predicates...),
		withRepositoryForEnts: ufeq.withRepositoryForEnts.Clone(),
		// clone intermediate query.
		sql:  ufeq.sql.Clone(),
		path: ufeq.path,
	}
}

// WithRepositoryForEnts tells the query-builder to eager-load the nodes that are connected to
// the "repository_for_ents" edge. The optional arguments are used to configure the query builder of the edge.
func (ufeq *UserForEntQuery) WithRepositoryForEnts(opts ...func(*RepositoryForEntQuery)) *UserForEntQuery {
	query := (&RepositoryForEntClient{config: ufeq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	ufeq.withRepositoryForEnts = query
	return ufeq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
func (ufeq *UserForEntQuery) GroupBy(field string, fields ...string) *UserForEntGroupBy {
	ufeq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &UserForEntGroupBy{build: ufeq}
	grbuild.flds = &ufeq.ctx.Fields
	grbuild.label = userforent.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
func (ufeq *UserForEntQuery) Select(fields ...string) *UserForEntSelect {
	ufeq.ctx.Fields = append(ufeq.ctx.Fields, fields...)
	sbuild := &UserForEntSelect{UserForEntQuery: ufeq}
	sbuild.label = userforent.Label
	sbuild.flds, sbuild.scan = &ufeq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a UserForEntSelect configured with the given aggregations.
func (ufeq *UserForEntQuery) Aggregate(fns ...AggregateFunc) *UserForEntSelect {
	return ufeq.Select().Aggregate(fns...)
}

func (ufeq *UserForEntQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range ufeq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, ufeq); err != nil {
				return err
			}
		}
	}
	for _, f := range ufeq.ctx.Fields {
		if !userforent.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if ufeq.path != nil {
		prev, err := ufeq.path(ctx)
		if err != nil {
			return err
		}
		ufeq.sql = prev
	}
	return nil
}

func (ufeq *UserForEntQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*UserForEnt, error) {
	var (
		nodes       = []*UserForEnt{}
		_spec       = ufeq.querySpec()
		loadedTypes = [1]bool{
			ufeq.withRepositoryForEnts != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*UserForEnt).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &UserForEnt{config: ufeq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, ufeq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := ufeq.withRepositoryForEnts; query != nil {
		if err := ufeq.loadRepositoryForEnts(ctx, query, nodes,
			func(n *UserForEnt) { n.Edges.RepositoryForEnts = []*RepositoryForEnt{} },
			func(n *UserForEnt, e *RepositoryForEnt) {
				n.Edges.RepositoryForEnts = append(n.Edges.RepositoryForEnts, e)
			}); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (ufeq *UserForEntQuery) loadRepositoryForEnts(ctx context.Context, query *RepositoryForEntQuery, nodes []*UserForEnt, init func(*UserForEnt), assign func(*UserForEnt, *RepositoryForEnt)) error {
	edgeIDs := make([]driver.Value, len(nodes))
	byID := make(map[int]*UserForEnt)
	nids := make(map[int]map[*UserForEnt]struct{})
	for i, node := range nodes {
		edgeIDs[i] = node.ID
		byID[node.ID] = node
		if init != nil {
			init(node)
		}
	}
	query.Where(func(s *sql.Selector) {
		joinT := sql.Table(userforent.RepositoryForEntsTable)
		s.Join(joinT).On(s.C(repositoryforent.FieldID), joinT.C(userforent.RepositoryForEntsPrimaryKey[0]))
		s.Where(sql.InValues(joinT.C(userforent.RepositoryForEntsPrimaryKey[1]), edgeIDs...))
		columns := s.SelectedColumns()
		s.Select(joinT.C(userforent.RepositoryForEntsPrimaryKey[1]))
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
					nids[inValue] = map[*UserForEnt]struct{}{byID[outValue]: {}}
					return assign(columns[1:], values[1:])
				}
				nids[inValue][byID[outValue]] = struct{}{}
				return nil
			}
		})
	})
	neighbors, err := withInterceptors[[]*RepositoryForEnt](ctx, query, qr, query.inters)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected "repository_for_ents" node returned %v`, n.ID)
		}
		for kn := range nodes {
			assign(kn, n)
		}
	}
	return nil
}

func (ufeq *UserForEntQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := ufeq.querySpec()
	_spec.Node.Columns = ufeq.ctx.Fields
	if len(ufeq.ctx.Fields) > 0 {
		_spec.Unique = ufeq.ctx.Unique != nil && *ufeq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, ufeq.driver, _spec)
}

func (ufeq *UserForEntQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(userforent.Table, userforent.Columns, sqlgraph.NewFieldSpec(userforent.FieldID, field.TypeInt))
	_spec.From = ufeq.sql
	if unique := ufeq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if ufeq.path != nil {
		_spec.Unique = true
	}
	if fields := ufeq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, userforent.FieldID)
		for i := range fields {
			if fields[i] != userforent.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := ufeq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := ufeq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := ufeq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := ufeq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (ufeq *UserForEntQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(ufeq.driver.Dialect())
	t1 := builder.Table(userforent.Table)
	columns := ufeq.ctx.Fields
	if len(columns) == 0 {
		columns = userforent.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if ufeq.sql != nil {
		selector = ufeq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if ufeq.ctx.Unique != nil && *ufeq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range ufeq.predicates {
		p(selector)
	}
	for _, p := range ufeq.order {
		p(selector)
	}
	if offset := ufeq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := ufeq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// UserForEntGroupBy is the group-by builder for UserForEnt entities.
type UserForEntGroupBy struct {
	selector
	build *UserForEntQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (ufegb *UserForEntGroupBy) Aggregate(fns ...AggregateFunc) *UserForEntGroupBy {
	ufegb.fns = append(ufegb.fns, fns...)
	return ufegb
}

// Scan applies the selector query and scans the result into the given value.
func (ufegb *UserForEntGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ufegb.build.ctx, ent.OpQueryGroupBy)
	if err := ufegb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*UserForEntQuery, *UserForEntGroupBy](ctx, ufegb.build, ufegb, ufegb.build.inters, v)
}

func (ufegb *UserForEntGroupBy) sqlScan(ctx context.Context, root *UserForEntQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(ufegb.fns))
	for _, fn := range ufegb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*ufegb.flds)+len(ufegb.fns))
		for _, f := range *ufegb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*ufegb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ufegb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// UserForEntSelect is the builder for selecting fields of UserForEnt entities.
type UserForEntSelect struct {
	*UserForEntQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ufes *UserForEntSelect) Aggregate(fns ...AggregateFunc) *UserForEntSelect {
	ufes.fns = append(ufes.fns, fns...)
	return ufes
}

// Scan applies the selector query and scans the result into the given value.
func (ufes *UserForEntSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ufes.ctx, ent.OpQuerySelect)
	if err := ufes.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*UserForEntQuery, *UserForEntSelect](ctx, ufes.UserForEntQuery, ufes, ufes.inters, v)
}

func (ufes *UserForEntSelect) sqlScan(ctx context.Context, root *UserForEntQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ufes.fns))
	for _, fn := range ufes.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ufes.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ufes.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
