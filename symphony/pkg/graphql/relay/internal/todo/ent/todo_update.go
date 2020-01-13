// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/pkg/graphql/relay/internal/todo/ent/predicate"
	"github.com/facebookincubator/symphony/pkg/graphql/relay/internal/todo/ent/todo"
)

// TodoUpdate is the builder for updating Todo entities.
type TodoUpdate struct {
	config
	text       *string
	predicates []predicate.Todo
}

// Where adds a new predicate for the builder.
func (tu *TodoUpdate) Where(ps ...predicate.Todo) *TodoUpdate {
	tu.predicates = append(tu.predicates, ps...)
	return tu
}

// SetText sets the text field.
func (tu *TodoUpdate) SetText(s string) *TodoUpdate {
	tu.text = &s
	return tu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (tu *TodoUpdate) Save(ctx context.Context) (int, error) {
	if tu.text != nil {
		if err := todo.TextValidator(*tu.text); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"text\": %v", err)
		}
	}
	return tu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (tu *TodoUpdate) SaveX(ctx context.Context) int {
	affected, err := tu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tu *TodoUpdate) Exec(ctx context.Context) error {
	_, err := tu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tu *TodoUpdate) ExecX(ctx context.Context) {
	if err := tu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tu *TodoUpdate) sqlSave(ctx context.Context) (n int, err error) {
	spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   todo.Table,
			Columns: todo.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: todo.FieldID,
			},
		},
	}
	if ps := tu.predicates; len(ps) > 0 {
		spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := tu.text; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: todo.FieldText,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, tu.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// TodoUpdateOne is the builder for updating a single Todo entity.
type TodoUpdateOne struct {
	config
	id   string
	text *string
}

// SetText sets the text field.
func (tuo *TodoUpdateOne) SetText(s string) *TodoUpdateOne {
	tuo.text = &s
	return tuo
}

// Save executes the query and returns the updated entity.
func (tuo *TodoUpdateOne) Save(ctx context.Context) (*Todo, error) {
	if tuo.text != nil {
		if err := todo.TextValidator(*tuo.text); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"text\": %v", err)
		}
	}
	return tuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (tuo *TodoUpdateOne) SaveX(ctx context.Context) *Todo {
	t, err := tuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return t
}

// Exec executes the query on the entity.
func (tuo *TodoUpdateOne) Exec(ctx context.Context) error {
	_, err := tuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tuo *TodoUpdateOne) ExecX(ctx context.Context) {
	if err := tuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tuo *TodoUpdateOne) sqlSave(ctx context.Context) (t *Todo, err error) {
	spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   todo.Table,
			Columns: todo.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  tuo.id,
				Type:   field.TypeString,
				Column: todo.FieldID,
			},
		},
	}
	if value := tuo.text; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: todo.FieldText,
		})
	}
	t = &Todo{config: tuo.config}
	spec.Assign = t.assignValues
	spec.ScanValues = t.scanValues()
	if err = sqlgraph.UpdateNode(ctx, tuo.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return t, nil
}
