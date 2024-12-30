// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/hedgehog125/project-reboot/ent/loginattempt"
	"github.com/hedgehog125/project-reboot/ent/predicate"
)

// LoginAttemptQuery is the builder for querying LoginAttempt entities.
type LoginAttemptQuery struct {
	config
	ctx        *QueryContext
	order      []loginattempt.OrderOption
	inters     []Interceptor
	predicates []predicate.LoginAttempt
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the LoginAttemptQuery builder.
func (laq *LoginAttemptQuery) Where(ps ...predicate.LoginAttempt) *LoginAttemptQuery {
	laq.predicates = append(laq.predicates, ps...)
	return laq
}

// Limit the number of records to be returned by this query.
func (laq *LoginAttemptQuery) Limit(limit int) *LoginAttemptQuery {
	laq.ctx.Limit = &limit
	return laq
}

// Offset to start from.
func (laq *LoginAttemptQuery) Offset(offset int) *LoginAttemptQuery {
	laq.ctx.Offset = &offset
	return laq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (laq *LoginAttemptQuery) Unique(unique bool) *LoginAttemptQuery {
	laq.ctx.Unique = &unique
	return laq
}

// Order specifies how the records should be ordered.
func (laq *LoginAttemptQuery) Order(o ...loginattempt.OrderOption) *LoginAttemptQuery {
	laq.order = append(laq.order, o...)
	return laq
}

// First returns the first LoginAttempt entity from the query.
// Returns a *NotFoundError when no LoginAttempt was found.
func (laq *LoginAttemptQuery) First(ctx context.Context) (*LoginAttempt, error) {
	nodes, err := laq.Limit(1).All(setContextOp(ctx, laq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{loginattempt.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (laq *LoginAttemptQuery) FirstX(ctx context.Context) *LoginAttempt {
	node, err := laq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first LoginAttempt ID from the query.
// Returns a *NotFoundError when no LoginAttempt ID was found.
func (laq *LoginAttemptQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = laq.Limit(1).IDs(setContextOp(ctx, laq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{loginattempt.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (laq *LoginAttemptQuery) FirstIDX(ctx context.Context) int {
	id, err := laq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single LoginAttempt entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one LoginAttempt entity is found.
// Returns a *NotFoundError when no LoginAttempt entities are found.
func (laq *LoginAttemptQuery) Only(ctx context.Context) (*LoginAttempt, error) {
	nodes, err := laq.Limit(2).All(setContextOp(ctx, laq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{loginattempt.Label}
	default:
		return nil, &NotSingularError{loginattempt.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (laq *LoginAttemptQuery) OnlyX(ctx context.Context) *LoginAttempt {
	node, err := laq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only LoginAttempt ID in the query.
// Returns a *NotSingularError when more than one LoginAttempt ID is found.
// Returns a *NotFoundError when no entities are found.
func (laq *LoginAttemptQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = laq.Limit(2).IDs(setContextOp(ctx, laq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{loginattempt.Label}
	default:
		err = &NotSingularError{loginattempt.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (laq *LoginAttemptQuery) OnlyIDX(ctx context.Context) int {
	id, err := laq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of LoginAttempts.
func (laq *LoginAttemptQuery) All(ctx context.Context) ([]*LoginAttempt, error) {
	ctx = setContextOp(ctx, laq.ctx, ent.OpQueryAll)
	if err := laq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*LoginAttempt, *LoginAttemptQuery]()
	return withInterceptors[[]*LoginAttempt](ctx, laq, qr, laq.inters)
}

// AllX is like All, but panics if an error occurs.
func (laq *LoginAttemptQuery) AllX(ctx context.Context) []*LoginAttempt {
	nodes, err := laq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of LoginAttempt IDs.
func (laq *LoginAttemptQuery) IDs(ctx context.Context) (ids []int, err error) {
	if laq.ctx.Unique == nil && laq.path != nil {
		laq.Unique(true)
	}
	ctx = setContextOp(ctx, laq.ctx, ent.OpQueryIDs)
	if err = laq.Select(loginattempt.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (laq *LoginAttemptQuery) IDsX(ctx context.Context) []int {
	ids, err := laq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (laq *LoginAttemptQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, laq.ctx, ent.OpQueryCount)
	if err := laq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, laq, querierCount[*LoginAttemptQuery](), laq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (laq *LoginAttemptQuery) CountX(ctx context.Context) int {
	count, err := laq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (laq *LoginAttemptQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, laq.ctx, ent.OpQueryExist)
	switch _, err := laq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (laq *LoginAttemptQuery) ExistX(ctx context.Context) bool {
	exist, err := laq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the LoginAttemptQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (laq *LoginAttemptQuery) Clone() *LoginAttemptQuery {
	if laq == nil {
		return nil
	}
	return &LoginAttemptQuery{
		config:     laq.config,
		ctx:        laq.ctx.Clone(),
		order:      append([]loginattempt.OrderOption{}, laq.order...),
		inters:     append([]Interceptor{}, laq.inters...),
		predicates: append([]predicate.LoginAttempt{}, laq.predicates...),
		// clone intermediate query.
		sql:  laq.sql.Clone(),
		path: laq.path,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Time time.Time `json:"time,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.LoginAttempt.Query().
//		GroupBy(loginattempt.FieldTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (laq *LoginAttemptQuery) GroupBy(field string, fields ...string) *LoginAttemptGroupBy {
	laq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &LoginAttemptGroupBy{build: laq}
	grbuild.flds = &laq.ctx.Fields
	grbuild.label = loginattempt.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Time time.Time `json:"time,omitempty"`
//	}
//
//	client.LoginAttempt.Query().
//		Select(loginattempt.FieldTime).
//		Scan(ctx, &v)
func (laq *LoginAttemptQuery) Select(fields ...string) *LoginAttemptSelect {
	laq.ctx.Fields = append(laq.ctx.Fields, fields...)
	sbuild := &LoginAttemptSelect{LoginAttemptQuery: laq}
	sbuild.label = loginattempt.Label
	sbuild.flds, sbuild.scan = &laq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a LoginAttemptSelect configured with the given aggregations.
func (laq *LoginAttemptQuery) Aggregate(fns ...AggregateFunc) *LoginAttemptSelect {
	return laq.Select().Aggregate(fns...)
}

func (laq *LoginAttemptQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range laq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, laq); err != nil {
				return err
			}
		}
	}
	for _, f := range laq.ctx.Fields {
		if !loginattempt.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if laq.path != nil {
		prev, err := laq.path(ctx)
		if err != nil {
			return err
		}
		laq.sql = prev
	}
	return nil
}

func (laq *LoginAttemptQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*LoginAttempt, error) {
	var (
		nodes = []*LoginAttempt{}
		_spec = laq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*LoginAttempt).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &LoginAttempt{config: laq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, laq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (laq *LoginAttemptQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := laq.querySpec()
	_spec.Node.Columns = laq.ctx.Fields
	if len(laq.ctx.Fields) > 0 {
		_spec.Unique = laq.ctx.Unique != nil && *laq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, laq.driver, _spec)
}

func (laq *LoginAttemptQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(loginattempt.Table, loginattempt.Columns, sqlgraph.NewFieldSpec(loginattempt.FieldID, field.TypeInt))
	_spec.From = laq.sql
	if unique := laq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if laq.path != nil {
		_spec.Unique = true
	}
	if fields := laq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, loginattempt.FieldID)
		for i := range fields {
			if fields[i] != loginattempt.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := laq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := laq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := laq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := laq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (laq *LoginAttemptQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(laq.driver.Dialect())
	t1 := builder.Table(loginattempt.Table)
	columns := laq.ctx.Fields
	if len(columns) == 0 {
		columns = loginattempt.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if laq.sql != nil {
		selector = laq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if laq.ctx.Unique != nil && *laq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range laq.predicates {
		p(selector)
	}
	for _, p := range laq.order {
		p(selector)
	}
	if offset := laq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := laq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// LoginAttemptGroupBy is the group-by builder for LoginAttempt entities.
type LoginAttemptGroupBy struct {
	selector
	build *LoginAttemptQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (lagb *LoginAttemptGroupBy) Aggregate(fns ...AggregateFunc) *LoginAttemptGroupBy {
	lagb.fns = append(lagb.fns, fns...)
	return lagb
}

// Scan applies the selector query and scans the result into the given value.
func (lagb *LoginAttemptGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, lagb.build.ctx, ent.OpQueryGroupBy)
	if err := lagb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*LoginAttemptQuery, *LoginAttemptGroupBy](ctx, lagb.build, lagb, lagb.build.inters, v)
}

func (lagb *LoginAttemptGroupBy) sqlScan(ctx context.Context, root *LoginAttemptQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(lagb.fns))
	for _, fn := range lagb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*lagb.flds)+len(lagb.fns))
		for _, f := range *lagb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*lagb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := lagb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// LoginAttemptSelect is the builder for selecting fields of LoginAttempt entities.
type LoginAttemptSelect struct {
	*LoginAttemptQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (las *LoginAttemptSelect) Aggregate(fns ...AggregateFunc) *LoginAttemptSelect {
	las.fns = append(las.fns, fns...)
	return las
}

// Scan applies the selector query and scans the result into the given value.
func (las *LoginAttemptSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, las.ctx, ent.OpQuerySelect)
	if err := las.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*LoginAttemptQuery, *LoginAttemptSelect](ctx, las.LoginAttemptQuery, las, las.inters, v)
}

func (las *LoginAttemptSelect) sqlScan(ctx context.Context, root *LoginAttemptQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(las.fns))
	for _, fn := range las.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*las.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := las.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
