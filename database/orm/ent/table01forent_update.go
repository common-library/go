// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/common-library/go/database/orm/ent/predicate"
	"github.com/common-library/go/database/orm/ent/table01forent"
)

// Table01ForEntUpdate is the builder for updating Table01ForEnt entities.
type Table01ForEntUpdate struct {
	config
	hooks    []Hook
	mutation *Table01ForEntMutation
}

// Where appends a list predicates to the Table01ForEntUpdate builder.
func (teu *Table01ForEntUpdate) Where(ps ...predicate.Table01ForEnt) *Table01ForEntUpdate {
	teu.mutation.Where(ps...)
	return teu
}

// SetCommonField01 sets the "common_field01" field.
func (teu *Table01ForEntUpdate) SetCommonField01(i int) *Table01ForEntUpdate {
	teu.mutation.ResetCommonField01()
	teu.mutation.SetCommonField01(i)
	return teu
}

// SetNillableCommonField01 sets the "common_field01" field if the given value is not nil.
func (teu *Table01ForEntUpdate) SetNillableCommonField01(i *int) *Table01ForEntUpdate {
	if i != nil {
		teu.SetCommonField01(*i)
	}
	return teu
}

// AddCommonField01 adds i to the "common_field01" field.
func (teu *Table01ForEntUpdate) AddCommonField01(i int) *Table01ForEntUpdate {
	teu.mutation.AddCommonField01(i)
	return teu
}

// ClearCommonField01 clears the value of the "common_field01" field.
func (teu *Table01ForEntUpdate) ClearCommonField01() *Table01ForEntUpdate {
	teu.mutation.ClearCommonField01()
	return teu
}

// SetUpdateTime sets the "update_time" field.
func (teu *Table01ForEntUpdate) SetUpdateTime(t time.Time) *Table01ForEntUpdate {
	teu.mutation.SetUpdateTime(t)
	return teu
}

// SetField01 sets the "field01" field.
func (teu *Table01ForEntUpdate) SetField01(s string) *Table01ForEntUpdate {
	teu.mutation.SetField01(s)
	return teu
}

// SetNillableField01 sets the "field01" field if the given value is not nil.
func (teu *Table01ForEntUpdate) SetNillableField01(s *string) *Table01ForEntUpdate {
	if s != nil {
		teu.SetField01(*s)
	}
	return teu
}

// SetField02 sets the "field02" field.
func (teu *Table01ForEntUpdate) SetField02(i int) *Table01ForEntUpdate {
	teu.mutation.ResetField02()
	teu.mutation.SetField02(i)
	return teu
}

// SetNillableField02 sets the "field02" field if the given value is not nil.
func (teu *Table01ForEntUpdate) SetNillableField02(i *int) *Table01ForEntUpdate {
	if i != nil {
		teu.SetField02(*i)
	}
	return teu
}

// AddField02 adds i to the "field02" field.
func (teu *Table01ForEntUpdate) AddField02(i int) *Table01ForEntUpdate {
	teu.mutation.AddField02(i)
	return teu
}

// SetField03 sets the "field03" field.
func (teu *Table01ForEntUpdate) SetField03(b bool) *Table01ForEntUpdate {
	teu.mutation.SetField03(b)
	return teu
}

// SetNillableField03 sets the "field03" field if the given value is not nil.
func (teu *Table01ForEntUpdate) SetNillableField03(b *bool) *Table01ForEntUpdate {
	if b != nil {
		teu.SetField03(*b)
	}
	return teu
}

// SetField04 sets the "field04" field.
func (teu *Table01ForEntUpdate) SetField04(t table01forent.Field04) *Table01ForEntUpdate {
	teu.mutation.SetField04(t)
	return teu
}

// SetNillableField04 sets the "field04" field if the given value is not nil.
func (teu *Table01ForEntUpdate) SetNillableField04(t *table01forent.Field04) *Table01ForEntUpdate {
	if t != nil {
		teu.SetField04(*t)
	}
	return teu
}

// SetField05 sets the "field05" field.
func (teu *Table01ForEntUpdate) SetField05(s string) *Table01ForEntUpdate {
	teu.mutation.SetField05(s)
	return teu
}

// SetNillableField05 sets the "field05" field if the given value is not nil.
func (teu *Table01ForEntUpdate) SetNillableField05(s *string) *Table01ForEntUpdate {
	if s != nil {
		teu.SetField05(*s)
	}
	return teu
}

// ClearField05 clears the value of the "field05" field.
func (teu *Table01ForEntUpdate) ClearField05() *Table01ForEntUpdate {
	teu.mutation.ClearField05()
	return teu
}

// SetField06 sets the "field06" field.
func (teu *Table01ForEntUpdate) SetField06(s string) *Table01ForEntUpdate {
	teu.mutation.SetField06(s)
	return teu
}

// SetNillableField06 sets the "field06" field if the given value is not nil.
func (teu *Table01ForEntUpdate) SetNillableField06(s *string) *Table01ForEntUpdate {
	if s != nil {
		teu.SetField06(*s)
	}
	return teu
}

// ClearField06 clears the value of the "field06" field.
func (teu *Table01ForEntUpdate) ClearField06() *Table01ForEntUpdate {
	teu.mutation.ClearField06()
	return teu
}

// SetField07 sets the "field07" field.
func (teu *Table01ForEntUpdate) SetField07(i int64) *Table01ForEntUpdate {
	teu.mutation.ResetField07()
	teu.mutation.SetField07(i)
	return teu
}

// AddField07 adds i to the "field07" field.
func (teu *Table01ForEntUpdate) AddField07(i int64) *Table01ForEntUpdate {
	teu.mutation.AddField07(i)
	return teu
}

// SetField08 sets the "field08" field.
func (teu *Table01ForEntUpdate) SetField08(f float64) *Table01ForEntUpdate {
	teu.mutation.ResetField08()
	teu.mutation.SetField08(f)
	return teu
}

// SetNillableField08 sets the "field08" field if the given value is not nil.
func (teu *Table01ForEntUpdate) SetNillableField08(f *float64) *Table01ForEntUpdate {
	if f != nil {
		teu.SetField08(*f)
	}
	return teu
}

// AddField08 adds f to the "field08" field.
func (teu *Table01ForEntUpdate) AddField08(f float64) *Table01ForEntUpdate {
	teu.mutation.AddField08(f)
	return teu
}

// SetField09 sets the "field09" field.
func (teu *Table01ForEntUpdate) SetField09(s string) *Table01ForEntUpdate {
	teu.mutation.SetField09(s)
	return teu
}

// SetNillableField09 sets the "field09" field if the given value is not nil.
func (teu *Table01ForEntUpdate) SetNillableField09(s *string) *Table01ForEntUpdate {
	if s != nil {
		teu.SetField09(*s)
	}
	return teu
}

// SetField10 sets the "field10" field.
func (teu *Table01ForEntUpdate) SetField10(s string) *Table01ForEntUpdate {
	teu.mutation.SetField10(s)
	return teu
}

// SetNillableField10 sets the "field10" field if the given value is not nil.
func (teu *Table01ForEntUpdate) SetNillableField10(s *string) *Table01ForEntUpdate {
	if s != nil {
		teu.SetField10(*s)
	}
	return teu
}

// ClearField10 clears the value of the "field10" field.
func (teu *Table01ForEntUpdate) ClearField10() *Table01ForEntUpdate {
	teu.mutation.ClearField10()
	return teu
}

// SetField11 sets the "field11" field.
func (teu *Table01ForEntUpdate) SetField11(s string) *Table01ForEntUpdate {
	teu.mutation.SetField11(s)
	return teu
}

// SetNillableField11 sets the "field11" field if the given value is not nil.
func (teu *Table01ForEntUpdate) SetNillableField11(s *string) *Table01ForEntUpdate {
	if s != nil {
		teu.SetField11(*s)
	}
	return teu
}

// ClearField11 clears the value of the "field11" field.
func (teu *Table01ForEntUpdate) ClearField11() *Table01ForEntUpdate {
	teu.mutation.ClearField11()
	return teu
}

// SetField12 sets the "field12" field.
func (teu *Table01ForEntUpdate) SetField12(s string) *Table01ForEntUpdate {
	teu.mutation.SetField12(s)
	return teu
}

// SetNillableField12 sets the "field12" field if the given value is not nil.
func (teu *Table01ForEntUpdate) SetNillableField12(s *string) *Table01ForEntUpdate {
	if s != nil {
		teu.SetField12(*s)
	}
	return teu
}

// ClearField12 clears the value of the "field12" field.
func (teu *Table01ForEntUpdate) ClearField12() *Table01ForEntUpdate {
	teu.mutation.ClearField12()
	return teu
}

// SetField13 sets the "field13" field.
func (teu *Table01ForEntUpdate) SetField13(s string) *Table01ForEntUpdate {
	teu.mutation.SetField13(s)
	return teu
}

// SetNillableField13 sets the "field13" field if the given value is not nil.
func (teu *Table01ForEntUpdate) SetNillableField13(s *string) *Table01ForEntUpdate {
	if s != nil {
		teu.SetField13(*s)
	}
	return teu
}

// ClearField13 clears the value of the "field13" field.
func (teu *Table01ForEntUpdate) ClearField13() *Table01ForEntUpdate {
	teu.mutation.ClearField13()
	return teu
}

// SetField14 sets the "field14" field.
func (teu *Table01ForEntUpdate) SetField14(s string) *Table01ForEntUpdate {
	teu.mutation.SetField14(s)
	return teu
}

// SetNillableField14 sets the "field14" field if the given value is not nil.
func (teu *Table01ForEntUpdate) SetNillableField14(s *string) *Table01ForEntUpdate {
	if s != nil {
		teu.SetField14(*s)
	}
	return teu
}

// ClearField14 clears the value of the "field14" field.
func (teu *Table01ForEntUpdate) ClearField14() *Table01ForEntUpdate {
	teu.mutation.ClearField14()
	return teu
}

// SetField15 sets the "field15" field.
func (teu *Table01ForEntUpdate) SetField15(s string) *Table01ForEntUpdate {
	teu.mutation.SetField15(s)
	return teu
}

// SetNillableField15 sets the "field15" field if the given value is not nil.
func (teu *Table01ForEntUpdate) SetNillableField15(s *string) *Table01ForEntUpdate {
	if s != nil {
		teu.SetField15(*s)
	}
	return teu
}

// ClearField15 clears the value of the "field15" field.
func (teu *Table01ForEntUpdate) ClearField15() *Table01ForEntUpdate {
	teu.mutation.ClearField15()
	return teu
}

// Mutation returns the Table01ForEntMutation object of the builder.
func (teu *Table01ForEntUpdate) Mutation() *Table01ForEntMutation {
	return teu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (teu *Table01ForEntUpdate) Save(ctx context.Context) (int, error) {
	teu.defaults()
	return withHooks(ctx, teu.sqlSave, teu.mutation, teu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (teu *Table01ForEntUpdate) SaveX(ctx context.Context) int {
	affected, err := teu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (teu *Table01ForEntUpdate) Exec(ctx context.Context) error {
	_, err := teu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (teu *Table01ForEntUpdate) ExecX(ctx context.Context) {
	if err := teu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (teu *Table01ForEntUpdate) defaults() {
	if _, ok := teu.mutation.UpdateTime(); !ok {
		v := table01forent.UpdateDefaultUpdateTime()
		teu.mutation.SetUpdateTime(v)
	}
	if _, ok := teu.mutation.Field07(); !ok {
		v := table01forent.UpdateDefaultField07()
		teu.mutation.SetField07(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (teu *Table01ForEntUpdate) check() error {
	if v, ok := teu.mutation.Field04(); ok {
		if err := table01forent.Field04Validator(v); err != nil {
			return &ValidationError{Name: "field04", err: fmt.Errorf(`ent: validator failed for field "Table01ForEnt.field04": %w`, err)}
		}
	}
	if v, ok := teu.mutation.Field09(); ok {
		if err := table01forent.Field09Validator(v); err != nil {
			return &ValidationError{Name: "field09", err: fmt.Errorf(`ent: validator failed for field "Table01ForEnt.field09": %w`, err)}
		}
	}
	return nil
}

func (teu *Table01ForEntUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := teu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(table01forent.Table, table01forent.Columns, sqlgraph.NewFieldSpec(table01forent.FieldID, field.TypeInt))
	if ps := teu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := teu.mutation.CommonField01(); ok {
		_spec.SetField(table01forent.FieldCommonField01, field.TypeInt, value)
	}
	if value, ok := teu.mutation.AddedCommonField01(); ok {
		_spec.AddField(table01forent.FieldCommonField01, field.TypeInt, value)
	}
	if teu.mutation.CommonField01Cleared() {
		_spec.ClearField(table01forent.FieldCommonField01, field.TypeInt)
	}
	if value, ok := teu.mutation.UpdateTime(); ok {
		_spec.SetField(table01forent.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := teu.mutation.Field01(); ok {
		_spec.SetField(table01forent.FieldField01, field.TypeString, value)
	}
	if value, ok := teu.mutation.Field02(); ok {
		_spec.SetField(table01forent.FieldField02, field.TypeInt, value)
	}
	if value, ok := teu.mutation.AddedField02(); ok {
		_spec.AddField(table01forent.FieldField02, field.TypeInt, value)
	}
	if value, ok := teu.mutation.Field03(); ok {
		_spec.SetField(table01forent.FieldField03, field.TypeBool, value)
	}
	if value, ok := teu.mutation.Field04(); ok {
		_spec.SetField(table01forent.FieldField04, field.TypeEnum, value)
	}
	if value, ok := teu.mutation.Field05(); ok {
		_spec.SetField(table01forent.FieldField05, field.TypeString, value)
	}
	if teu.mutation.Field05Cleared() {
		_spec.ClearField(table01forent.FieldField05, field.TypeString)
	}
	if value, ok := teu.mutation.Field06(); ok {
		_spec.SetField(table01forent.FieldField06, field.TypeString, value)
	}
	if teu.mutation.Field06Cleared() {
		_spec.ClearField(table01forent.FieldField06, field.TypeString)
	}
	if value, ok := teu.mutation.Field07(); ok {
		_spec.SetField(table01forent.FieldField07, field.TypeInt64, value)
	}
	if value, ok := teu.mutation.AddedField07(); ok {
		_spec.AddField(table01forent.FieldField07, field.TypeInt64, value)
	}
	if value, ok := teu.mutation.Field08(); ok {
		_spec.SetField(table01forent.FieldField08, field.TypeFloat64, value)
	}
	if value, ok := teu.mutation.AddedField08(); ok {
		_spec.AddField(table01forent.FieldField08, field.TypeFloat64, value)
	}
	if value, ok := teu.mutation.Field09(); ok {
		_spec.SetField(table01forent.FieldField09, field.TypeString, value)
	}
	if value, ok := teu.mutation.Field10(); ok {
		_spec.SetField(table01forent.FieldField10, field.TypeString, value)
	}
	if teu.mutation.Field10Cleared() {
		_spec.ClearField(table01forent.FieldField10, field.TypeString)
	}
	if value, ok := teu.mutation.Field11(); ok {
		_spec.SetField(table01forent.FieldField11, field.TypeString, value)
	}
	if teu.mutation.Field11Cleared() {
		_spec.ClearField(table01forent.FieldField11, field.TypeString)
	}
	if value, ok := teu.mutation.Field12(); ok {
		_spec.SetField(table01forent.FieldField12, field.TypeString, value)
	}
	if teu.mutation.Field12Cleared() {
		_spec.ClearField(table01forent.FieldField12, field.TypeString)
	}
	if value, ok := teu.mutation.Field13(); ok {
		_spec.SetField(table01forent.FieldField13, field.TypeString, value)
	}
	if teu.mutation.Field13Cleared() {
		_spec.ClearField(table01forent.FieldField13, field.TypeString)
	}
	if value, ok := teu.mutation.Field14(); ok {
		_spec.SetField(table01forent.FieldField14, field.TypeString, value)
	}
	if teu.mutation.Field14Cleared() {
		_spec.ClearField(table01forent.FieldField14, field.TypeString)
	}
	if value, ok := teu.mutation.Field15(); ok {
		_spec.SetField(table01forent.FieldField15, field.TypeString, value)
	}
	if teu.mutation.Field15Cleared() {
		_spec.ClearField(table01forent.FieldField15, field.TypeString)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, teu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{table01forent.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	teu.mutation.done = true
	return n, nil
}

// Table01ForEntUpdateOne is the builder for updating a single Table01ForEnt entity.
type Table01ForEntUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *Table01ForEntMutation
}

// SetCommonField01 sets the "common_field01" field.
func (teuo *Table01ForEntUpdateOne) SetCommonField01(i int) *Table01ForEntUpdateOne {
	teuo.mutation.ResetCommonField01()
	teuo.mutation.SetCommonField01(i)
	return teuo
}

// SetNillableCommonField01 sets the "common_field01" field if the given value is not nil.
func (teuo *Table01ForEntUpdateOne) SetNillableCommonField01(i *int) *Table01ForEntUpdateOne {
	if i != nil {
		teuo.SetCommonField01(*i)
	}
	return teuo
}

// AddCommonField01 adds i to the "common_field01" field.
func (teuo *Table01ForEntUpdateOne) AddCommonField01(i int) *Table01ForEntUpdateOne {
	teuo.mutation.AddCommonField01(i)
	return teuo
}

// ClearCommonField01 clears the value of the "common_field01" field.
func (teuo *Table01ForEntUpdateOne) ClearCommonField01() *Table01ForEntUpdateOne {
	teuo.mutation.ClearCommonField01()
	return teuo
}

// SetUpdateTime sets the "update_time" field.
func (teuo *Table01ForEntUpdateOne) SetUpdateTime(t time.Time) *Table01ForEntUpdateOne {
	teuo.mutation.SetUpdateTime(t)
	return teuo
}

// SetField01 sets the "field01" field.
func (teuo *Table01ForEntUpdateOne) SetField01(s string) *Table01ForEntUpdateOne {
	teuo.mutation.SetField01(s)
	return teuo
}

// SetNillableField01 sets the "field01" field if the given value is not nil.
func (teuo *Table01ForEntUpdateOne) SetNillableField01(s *string) *Table01ForEntUpdateOne {
	if s != nil {
		teuo.SetField01(*s)
	}
	return teuo
}

// SetField02 sets the "field02" field.
func (teuo *Table01ForEntUpdateOne) SetField02(i int) *Table01ForEntUpdateOne {
	teuo.mutation.ResetField02()
	teuo.mutation.SetField02(i)
	return teuo
}

// SetNillableField02 sets the "field02" field if the given value is not nil.
func (teuo *Table01ForEntUpdateOne) SetNillableField02(i *int) *Table01ForEntUpdateOne {
	if i != nil {
		teuo.SetField02(*i)
	}
	return teuo
}

// AddField02 adds i to the "field02" field.
func (teuo *Table01ForEntUpdateOne) AddField02(i int) *Table01ForEntUpdateOne {
	teuo.mutation.AddField02(i)
	return teuo
}

// SetField03 sets the "field03" field.
func (teuo *Table01ForEntUpdateOne) SetField03(b bool) *Table01ForEntUpdateOne {
	teuo.mutation.SetField03(b)
	return teuo
}

// SetNillableField03 sets the "field03" field if the given value is not nil.
func (teuo *Table01ForEntUpdateOne) SetNillableField03(b *bool) *Table01ForEntUpdateOne {
	if b != nil {
		teuo.SetField03(*b)
	}
	return teuo
}

// SetField04 sets the "field04" field.
func (teuo *Table01ForEntUpdateOne) SetField04(t table01forent.Field04) *Table01ForEntUpdateOne {
	teuo.mutation.SetField04(t)
	return teuo
}

// SetNillableField04 sets the "field04" field if the given value is not nil.
func (teuo *Table01ForEntUpdateOne) SetNillableField04(t *table01forent.Field04) *Table01ForEntUpdateOne {
	if t != nil {
		teuo.SetField04(*t)
	}
	return teuo
}

// SetField05 sets the "field05" field.
func (teuo *Table01ForEntUpdateOne) SetField05(s string) *Table01ForEntUpdateOne {
	teuo.mutation.SetField05(s)
	return teuo
}

// SetNillableField05 sets the "field05" field if the given value is not nil.
func (teuo *Table01ForEntUpdateOne) SetNillableField05(s *string) *Table01ForEntUpdateOne {
	if s != nil {
		teuo.SetField05(*s)
	}
	return teuo
}

// ClearField05 clears the value of the "field05" field.
func (teuo *Table01ForEntUpdateOne) ClearField05() *Table01ForEntUpdateOne {
	teuo.mutation.ClearField05()
	return teuo
}

// SetField06 sets the "field06" field.
func (teuo *Table01ForEntUpdateOne) SetField06(s string) *Table01ForEntUpdateOne {
	teuo.mutation.SetField06(s)
	return teuo
}

// SetNillableField06 sets the "field06" field if the given value is not nil.
func (teuo *Table01ForEntUpdateOne) SetNillableField06(s *string) *Table01ForEntUpdateOne {
	if s != nil {
		teuo.SetField06(*s)
	}
	return teuo
}

// ClearField06 clears the value of the "field06" field.
func (teuo *Table01ForEntUpdateOne) ClearField06() *Table01ForEntUpdateOne {
	teuo.mutation.ClearField06()
	return teuo
}

// SetField07 sets the "field07" field.
func (teuo *Table01ForEntUpdateOne) SetField07(i int64) *Table01ForEntUpdateOne {
	teuo.mutation.ResetField07()
	teuo.mutation.SetField07(i)
	return teuo
}

// AddField07 adds i to the "field07" field.
func (teuo *Table01ForEntUpdateOne) AddField07(i int64) *Table01ForEntUpdateOne {
	teuo.mutation.AddField07(i)
	return teuo
}

// SetField08 sets the "field08" field.
func (teuo *Table01ForEntUpdateOne) SetField08(f float64) *Table01ForEntUpdateOne {
	teuo.mutation.ResetField08()
	teuo.mutation.SetField08(f)
	return teuo
}

// SetNillableField08 sets the "field08" field if the given value is not nil.
func (teuo *Table01ForEntUpdateOne) SetNillableField08(f *float64) *Table01ForEntUpdateOne {
	if f != nil {
		teuo.SetField08(*f)
	}
	return teuo
}

// AddField08 adds f to the "field08" field.
func (teuo *Table01ForEntUpdateOne) AddField08(f float64) *Table01ForEntUpdateOne {
	teuo.mutation.AddField08(f)
	return teuo
}

// SetField09 sets the "field09" field.
func (teuo *Table01ForEntUpdateOne) SetField09(s string) *Table01ForEntUpdateOne {
	teuo.mutation.SetField09(s)
	return teuo
}

// SetNillableField09 sets the "field09" field if the given value is not nil.
func (teuo *Table01ForEntUpdateOne) SetNillableField09(s *string) *Table01ForEntUpdateOne {
	if s != nil {
		teuo.SetField09(*s)
	}
	return teuo
}

// SetField10 sets the "field10" field.
func (teuo *Table01ForEntUpdateOne) SetField10(s string) *Table01ForEntUpdateOne {
	teuo.mutation.SetField10(s)
	return teuo
}

// SetNillableField10 sets the "field10" field if the given value is not nil.
func (teuo *Table01ForEntUpdateOne) SetNillableField10(s *string) *Table01ForEntUpdateOne {
	if s != nil {
		teuo.SetField10(*s)
	}
	return teuo
}

// ClearField10 clears the value of the "field10" field.
func (teuo *Table01ForEntUpdateOne) ClearField10() *Table01ForEntUpdateOne {
	teuo.mutation.ClearField10()
	return teuo
}

// SetField11 sets the "field11" field.
func (teuo *Table01ForEntUpdateOne) SetField11(s string) *Table01ForEntUpdateOne {
	teuo.mutation.SetField11(s)
	return teuo
}

// SetNillableField11 sets the "field11" field if the given value is not nil.
func (teuo *Table01ForEntUpdateOne) SetNillableField11(s *string) *Table01ForEntUpdateOne {
	if s != nil {
		teuo.SetField11(*s)
	}
	return teuo
}

// ClearField11 clears the value of the "field11" field.
func (teuo *Table01ForEntUpdateOne) ClearField11() *Table01ForEntUpdateOne {
	teuo.mutation.ClearField11()
	return teuo
}

// SetField12 sets the "field12" field.
func (teuo *Table01ForEntUpdateOne) SetField12(s string) *Table01ForEntUpdateOne {
	teuo.mutation.SetField12(s)
	return teuo
}

// SetNillableField12 sets the "field12" field if the given value is not nil.
func (teuo *Table01ForEntUpdateOne) SetNillableField12(s *string) *Table01ForEntUpdateOne {
	if s != nil {
		teuo.SetField12(*s)
	}
	return teuo
}

// ClearField12 clears the value of the "field12" field.
func (teuo *Table01ForEntUpdateOne) ClearField12() *Table01ForEntUpdateOne {
	teuo.mutation.ClearField12()
	return teuo
}

// SetField13 sets the "field13" field.
func (teuo *Table01ForEntUpdateOne) SetField13(s string) *Table01ForEntUpdateOne {
	teuo.mutation.SetField13(s)
	return teuo
}

// SetNillableField13 sets the "field13" field if the given value is not nil.
func (teuo *Table01ForEntUpdateOne) SetNillableField13(s *string) *Table01ForEntUpdateOne {
	if s != nil {
		teuo.SetField13(*s)
	}
	return teuo
}

// ClearField13 clears the value of the "field13" field.
func (teuo *Table01ForEntUpdateOne) ClearField13() *Table01ForEntUpdateOne {
	teuo.mutation.ClearField13()
	return teuo
}

// SetField14 sets the "field14" field.
func (teuo *Table01ForEntUpdateOne) SetField14(s string) *Table01ForEntUpdateOne {
	teuo.mutation.SetField14(s)
	return teuo
}

// SetNillableField14 sets the "field14" field if the given value is not nil.
func (teuo *Table01ForEntUpdateOne) SetNillableField14(s *string) *Table01ForEntUpdateOne {
	if s != nil {
		teuo.SetField14(*s)
	}
	return teuo
}

// ClearField14 clears the value of the "field14" field.
func (teuo *Table01ForEntUpdateOne) ClearField14() *Table01ForEntUpdateOne {
	teuo.mutation.ClearField14()
	return teuo
}

// SetField15 sets the "field15" field.
func (teuo *Table01ForEntUpdateOne) SetField15(s string) *Table01ForEntUpdateOne {
	teuo.mutation.SetField15(s)
	return teuo
}

// SetNillableField15 sets the "field15" field if the given value is not nil.
func (teuo *Table01ForEntUpdateOne) SetNillableField15(s *string) *Table01ForEntUpdateOne {
	if s != nil {
		teuo.SetField15(*s)
	}
	return teuo
}

// ClearField15 clears the value of the "field15" field.
func (teuo *Table01ForEntUpdateOne) ClearField15() *Table01ForEntUpdateOne {
	teuo.mutation.ClearField15()
	return teuo
}

// Mutation returns the Table01ForEntMutation object of the builder.
func (teuo *Table01ForEntUpdateOne) Mutation() *Table01ForEntMutation {
	return teuo.mutation
}

// Where appends a list predicates to the Table01ForEntUpdate builder.
func (teuo *Table01ForEntUpdateOne) Where(ps ...predicate.Table01ForEnt) *Table01ForEntUpdateOne {
	teuo.mutation.Where(ps...)
	return teuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (teuo *Table01ForEntUpdateOne) Select(field string, fields ...string) *Table01ForEntUpdateOne {
	teuo.fields = append([]string{field}, fields...)
	return teuo
}

// Save executes the query and returns the updated Table01ForEnt entity.
func (teuo *Table01ForEntUpdateOne) Save(ctx context.Context) (*Table01ForEnt, error) {
	teuo.defaults()
	return withHooks(ctx, teuo.sqlSave, teuo.mutation, teuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (teuo *Table01ForEntUpdateOne) SaveX(ctx context.Context) *Table01ForEnt {
	node, err := teuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (teuo *Table01ForEntUpdateOne) Exec(ctx context.Context) error {
	_, err := teuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (teuo *Table01ForEntUpdateOne) ExecX(ctx context.Context) {
	if err := teuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (teuo *Table01ForEntUpdateOne) defaults() {
	if _, ok := teuo.mutation.UpdateTime(); !ok {
		v := table01forent.UpdateDefaultUpdateTime()
		teuo.mutation.SetUpdateTime(v)
	}
	if _, ok := teuo.mutation.Field07(); !ok {
		v := table01forent.UpdateDefaultField07()
		teuo.mutation.SetField07(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (teuo *Table01ForEntUpdateOne) check() error {
	if v, ok := teuo.mutation.Field04(); ok {
		if err := table01forent.Field04Validator(v); err != nil {
			return &ValidationError{Name: "field04", err: fmt.Errorf(`ent: validator failed for field "Table01ForEnt.field04": %w`, err)}
		}
	}
	if v, ok := teuo.mutation.Field09(); ok {
		if err := table01forent.Field09Validator(v); err != nil {
			return &ValidationError{Name: "field09", err: fmt.Errorf(`ent: validator failed for field "Table01ForEnt.field09": %w`, err)}
		}
	}
	return nil
}

func (teuo *Table01ForEntUpdateOne) sqlSave(ctx context.Context) (_node *Table01ForEnt, err error) {
	if err := teuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(table01forent.Table, table01forent.Columns, sqlgraph.NewFieldSpec(table01forent.FieldID, field.TypeInt))
	id, ok := teuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Table01ForEnt.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := teuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, table01forent.FieldID)
		for _, f := range fields {
			if !table01forent.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != table01forent.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := teuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := teuo.mutation.CommonField01(); ok {
		_spec.SetField(table01forent.FieldCommonField01, field.TypeInt, value)
	}
	if value, ok := teuo.mutation.AddedCommonField01(); ok {
		_spec.AddField(table01forent.FieldCommonField01, field.TypeInt, value)
	}
	if teuo.mutation.CommonField01Cleared() {
		_spec.ClearField(table01forent.FieldCommonField01, field.TypeInt)
	}
	if value, ok := teuo.mutation.UpdateTime(); ok {
		_spec.SetField(table01forent.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := teuo.mutation.Field01(); ok {
		_spec.SetField(table01forent.FieldField01, field.TypeString, value)
	}
	if value, ok := teuo.mutation.Field02(); ok {
		_spec.SetField(table01forent.FieldField02, field.TypeInt, value)
	}
	if value, ok := teuo.mutation.AddedField02(); ok {
		_spec.AddField(table01forent.FieldField02, field.TypeInt, value)
	}
	if value, ok := teuo.mutation.Field03(); ok {
		_spec.SetField(table01forent.FieldField03, field.TypeBool, value)
	}
	if value, ok := teuo.mutation.Field04(); ok {
		_spec.SetField(table01forent.FieldField04, field.TypeEnum, value)
	}
	if value, ok := teuo.mutation.Field05(); ok {
		_spec.SetField(table01forent.FieldField05, field.TypeString, value)
	}
	if teuo.mutation.Field05Cleared() {
		_spec.ClearField(table01forent.FieldField05, field.TypeString)
	}
	if value, ok := teuo.mutation.Field06(); ok {
		_spec.SetField(table01forent.FieldField06, field.TypeString, value)
	}
	if teuo.mutation.Field06Cleared() {
		_spec.ClearField(table01forent.FieldField06, field.TypeString)
	}
	if value, ok := teuo.mutation.Field07(); ok {
		_spec.SetField(table01forent.FieldField07, field.TypeInt64, value)
	}
	if value, ok := teuo.mutation.AddedField07(); ok {
		_spec.AddField(table01forent.FieldField07, field.TypeInt64, value)
	}
	if value, ok := teuo.mutation.Field08(); ok {
		_spec.SetField(table01forent.FieldField08, field.TypeFloat64, value)
	}
	if value, ok := teuo.mutation.AddedField08(); ok {
		_spec.AddField(table01forent.FieldField08, field.TypeFloat64, value)
	}
	if value, ok := teuo.mutation.Field09(); ok {
		_spec.SetField(table01forent.FieldField09, field.TypeString, value)
	}
	if value, ok := teuo.mutation.Field10(); ok {
		_spec.SetField(table01forent.FieldField10, field.TypeString, value)
	}
	if teuo.mutation.Field10Cleared() {
		_spec.ClearField(table01forent.FieldField10, field.TypeString)
	}
	if value, ok := teuo.mutation.Field11(); ok {
		_spec.SetField(table01forent.FieldField11, field.TypeString, value)
	}
	if teuo.mutation.Field11Cleared() {
		_spec.ClearField(table01forent.FieldField11, field.TypeString)
	}
	if value, ok := teuo.mutation.Field12(); ok {
		_spec.SetField(table01forent.FieldField12, field.TypeString, value)
	}
	if teuo.mutation.Field12Cleared() {
		_spec.ClearField(table01forent.FieldField12, field.TypeString)
	}
	if value, ok := teuo.mutation.Field13(); ok {
		_spec.SetField(table01forent.FieldField13, field.TypeString, value)
	}
	if teuo.mutation.Field13Cleared() {
		_spec.ClearField(table01forent.FieldField13, field.TypeString)
	}
	if value, ok := teuo.mutation.Field14(); ok {
		_spec.SetField(table01forent.FieldField14, field.TypeString, value)
	}
	if teuo.mutation.Field14Cleared() {
		_spec.ClearField(table01forent.FieldField14, field.TypeString)
	}
	if value, ok := teuo.mutation.Field15(); ok {
		_spec.SetField(table01forent.FieldField15, field.TypeString, value)
	}
	if teuo.mutation.Field15Cleared() {
		_spec.ClearField(table01forent.FieldField15, field.TypeString)
	}
	_node = &Table01ForEnt{config: teuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, teuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{table01forent.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	teuo.mutation.done = true
	return _node, nil
}
