// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/floorplan"
	"github.com/facebookincubator/symphony/graph/ent/hyperlink"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/survey"
	"github.com/facebookincubator/symphony/graph/ent/surveycellscan"
	"github.com/facebookincubator/symphony/graph/ent/surveywifiscan"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// LocationUpdate is the builder for updating Location entities.
type LocationUpdate struct {
	config
	hooks      []Hook
	mutation   *LocationMutation
	predicates []predicate.Location
}

// Where adds a new predicate for the builder.
func (lu *LocationUpdate) Where(ps ...predicate.Location) *LocationUpdate {
	lu.predicates = append(lu.predicates, ps...)
	return lu
}

// SetName sets the name field.
func (lu *LocationUpdate) SetName(s string) *LocationUpdate {
	lu.mutation.SetName(s)
	return lu
}

// SetExternalID sets the external_id field.
func (lu *LocationUpdate) SetExternalID(s string) *LocationUpdate {
	lu.mutation.SetExternalID(s)
	return lu
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (lu *LocationUpdate) SetNillableExternalID(s *string) *LocationUpdate {
	if s != nil {
		lu.SetExternalID(*s)
	}
	return lu
}

// ClearExternalID clears the value of external_id.
func (lu *LocationUpdate) ClearExternalID() *LocationUpdate {
	lu.mutation.ClearExternalID()
	return lu
}

// SetLatitude sets the latitude field.
func (lu *LocationUpdate) SetLatitude(f float64) *LocationUpdate {
	lu.mutation.ResetLatitude()
	lu.mutation.SetLatitude(f)
	return lu
}

// SetNillableLatitude sets the latitude field if the given value is not nil.
func (lu *LocationUpdate) SetNillableLatitude(f *float64) *LocationUpdate {
	if f != nil {
		lu.SetLatitude(*f)
	}
	return lu
}

// AddLatitude adds f to latitude.
func (lu *LocationUpdate) AddLatitude(f float64) *LocationUpdate {
	lu.mutation.AddLatitude(f)
	return lu
}

// SetLongitude sets the longitude field.
func (lu *LocationUpdate) SetLongitude(f float64) *LocationUpdate {
	lu.mutation.ResetLongitude()
	lu.mutation.SetLongitude(f)
	return lu
}

// SetNillableLongitude sets the longitude field if the given value is not nil.
func (lu *LocationUpdate) SetNillableLongitude(f *float64) *LocationUpdate {
	if f != nil {
		lu.SetLongitude(*f)
	}
	return lu
}

// AddLongitude adds f to longitude.
func (lu *LocationUpdate) AddLongitude(f float64) *LocationUpdate {
	lu.mutation.AddLongitude(f)
	return lu
}

// SetSiteSurveyNeeded sets the site_survey_needed field.
func (lu *LocationUpdate) SetSiteSurveyNeeded(b bool) *LocationUpdate {
	lu.mutation.SetSiteSurveyNeeded(b)
	return lu
}

// SetNillableSiteSurveyNeeded sets the site_survey_needed field if the given value is not nil.
func (lu *LocationUpdate) SetNillableSiteSurveyNeeded(b *bool) *LocationUpdate {
	if b != nil {
		lu.SetSiteSurveyNeeded(*b)
	}
	return lu
}

// ClearSiteSurveyNeeded clears the value of site_survey_needed.
func (lu *LocationUpdate) ClearSiteSurveyNeeded() *LocationUpdate {
	lu.mutation.ClearSiteSurveyNeeded()
	return lu
}

// SetTypeID sets the type edge to LocationType by id.
func (lu *LocationUpdate) SetTypeID(id int) *LocationUpdate {
	lu.mutation.SetTypeID(id)
	return lu
}

// SetType sets the type edge to LocationType.
func (lu *LocationUpdate) SetType(l *LocationType) *LocationUpdate {
	return lu.SetTypeID(l.ID)
}

// SetParentID sets the parent edge to Location by id.
func (lu *LocationUpdate) SetParentID(id int) *LocationUpdate {
	lu.mutation.SetParentID(id)
	return lu
}

// SetNillableParentID sets the parent edge to Location by id if the given value is not nil.
func (lu *LocationUpdate) SetNillableParentID(id *int) *LocationUpdate {
	if id != nil {
		lu = lu.SetParentID(*id)
	}
	return lu
}

// SetParent sets the parent edge to Location.
func (lu *LocationUpdate) SetParent(l *Location) *LocationUpdate {
	return lu.SetParentID(l.ID)
}

// AddChildIDs adds the children edge to Location by ids.
func (lu *LocationUpdate) AddChildIDs(ids ...int) *LocationUpdate {
	lu.mutation.AddChildIDs(ids...)
	return lu
}

// AddChildren adds the children edges to Location.
func (lu *LocationUpdate) AddChildren(l ...*Location) *LocationUpdate {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return lu.AddChildIDs(ids...)
}

// AddFileIDs adds the files edge to File by ids.
func (lu *LocationUpdate) AddFileIDs(ids ...int) *LocationUpdate {
	lu.mutation.AddFileIDs(ids...)
	return lu
}

// AddFiles adds the files edges to File.
func (lu *LocationUpdate) AddFiles(f ...*File) *LocationUpdate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return lu.AddFileIDs(ids...)
}

// AddHyperlinkIDs adds the hyperlinks edge to Hyperlink by ids.
func (lu *LocationUpdate) AddHyperlinkIDs(ids ...int) *LocationUpdate {
	lu.mutation.AddHyperlinkIDs(ids...)
	return lu
}

// AddHyperlinks adds the hyperlinks edges to Hyperlink.
func (lu *LocationUpdate) AddHyperlinks(h ...*Hyperlink) *LocationUpdate {
	ids := make([]int, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return lu.AddHyperlinkIDs(ids...)
}

// AddEquipmentIDs adds the equipment edge to Equipment by ids.
func (lu *LocationUpdate) AddEquipmentIDs(ids ...int) *LocationUpdate {
	lu.mutation.AddEquipmentIDs(ids...)
	return lu
}

// AddEquipment adds the equipment edges to Equipment.
func (lu *LocationUpdate) AddEquipment(e ...*Equipment) *LocationUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return lu.AddEquipmentIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (lu *LocationUpdate) AddPropertyIDs(ids ...int) *LocationUpdate {
	lu.mutation.AddPropertyIDs(ids...)
	return lu
}

// AddProperties adds the properties edges to Property.
func (lu *LocationUpdate) AddProperties(p ...*Property) *LocationUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return lu.AddPropertyIDs(ids...)
}

// AddSurveyIDs adds the survey edge to Survey by ids.
func (lu *LocationUpdate) AddSurveyIDs(ids ...int) *LocationUpdate {
	lu.mutation.AddSurveyIDs(ids...)
	return lu
}

// AddSurvey adds the survey edges to Survey.
func (lu *LocationUpdate) AddSurvey(s ...*Survey) *LocationUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lu.AddSurveyIDs(ids...)
}

// AddWifiScanIDs adds the wifi_scan edge to SurveyWiFiScan by ids.
func (lu *LocationUpdate) AddWifiScanIDs(ids ...int) *LocationUpdate {
	lu.mutation.AddWifiScanIDs(ids...)
	return lu
}

// AddWifiScan adds the wifi_scan edges to SurveyWiFiScan.
func (lu *LocationUpdate) AddWifiScan(s ...*SurveyWiFiScan) *LocationUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lu.AddWifiScanIDs(ids...)
}

// AddCellScanIDs adds the cell_scan edge to SurveyCellScan by ids.
func (lu *LocationUpdate) AddCellScanIDs(ids ...int) *LocationUpdate {
	lu.mutation.AddCellScanIDs(ids...)
	return lu
}

// AddCellScan adds the cell_scan edges to SurveyCellScan.
func (lu *LocationUpdate) AddCellScan(s ...*SurveyCellScan) *LocationUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lu.AddCellScanIDs(ids...)
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (lu *LocationUpdate) AddWorkOrderIDs(ids ...int) *LocationUpdate {
	lu.mutation.AddWorkOrderIDs(ids...)
	return lu
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (lu *LocationUpdate) AddWorkOrders(w ...*WorkOrder) *LocationUpdate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return lu.AddWorkOrderIDs(ids...)
}

// AddFloorPlanIDs adds the floor_plans edge to FloorPlan by ids.
func (lu *LocationUpdate) AddFloorPlanIDs(ids ...int) *LocationUpdate {
	lu.mutation.AddFloorPlanIDs(ids...)
	return lu
}

// AddFloorPlans adds the floor_plans edges to FloorPlan.
func (lu *LocationUpdate) AddFloorPlans(f ...*FloorPlan) *LocationUpdate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return lu.AddFloorPlanIDs(ids...)
}

// ClearType clears the type edge to LocationType.
func (lu *LocationUpdate) ClearType() *LocationUpdate {
	lu.mutation.ClearType()
	return lu
}

// ClearParent clears the parent edge to Location.
func (lu *LocationUpdate) ClearParent() *LocationUpdate {
	lu.mutation.ClearParent()
	return lu
}

// RemoveChildIDs removes the children edge to Location by ids.
func (lu *LocationUpdate) RemoveChildIDs(ids ...int) *LocationUpdate {
	lu.mutation.RemoveChildIDs(ids...)
	return lu
}

// RemoveChildren removes children edges to Location.
func (lu *LocationUpdate) RemoveChildren(l ...*Location) *LocationUpdate {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return lu.RemoveChildIDs(ids...)
}

// RemoveFileIDs removes the files edge to File by ids.
func (lu *LocationUpdate) RemoveFileIDs(ids ...int) *LocationUpdate {
	lu.mutation.RemoveFileIDs(ids...)
	return lu
}

// RemoveFiles removes files edges to File.
func (lu *LocationUpdate) RemoveFiles(f ...*File) *LocationUpdate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return lu.RemoveFileIDs(ids...)
}

// RemoveHyperlinkIDs removes the hyperlinks edge to Hyperlink by ids.
func (lu *LocationUpdate) RemoveHyperlinkIDs(ids ...int) *LocationUpdate {
	lu.mutation.RemoveHyperlinkIDs(ids...)
	return lu
}

// RemoveHyperlinks removes hyperlinks edges to Hyperlink.
func (lu *LocationUpdate) RemoveHyperlinks(h ...*Hyperlink) *LocationUpdate {
	ids := make([]int, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return lu.RemoveHyperlinkIDs(ids...)
}

// RemoveEquipmentIDs removes the equipment edge to Equipment by ids.
func (lu *LocationUpdate) RemoveEquipmentIDs(ids ...int) *LocationUpdate {
	lu.mutation.RemoveEquipmentIDs(ids...)
	return lu
}

// RemoveEquipment removes equipment edges to Equipment.
func (lu *LocationUpdate) RemoveEquipment(e ...*Equipment) *LocationUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return lu.RemoveEquipmentIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (lu *LocationUpdate) RemovePropertyIDs(ids ...int) *LocationUpdate {
	lu.mutation.RemovePropertyIDs(ids...)
	return lu
}

// RemoveProperties removes properties edges to Property.
func (lu *LocationUpdate) RemoveProperties(p ...*Property) *LocationUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return lu.RemovePropertyIDs(ids...)
}

// RemoveSurveyIDs removes the survey edge to Survey by ids.
func (lu *LocationUpdate) RemoveSurveyIDs(ids ...int) *LocationUpdate {
	lu.mutation.RemoveSurveyIDs(ids...)
	return lu
}

// RemoveSurvey removes survey edges to Survey.
func (lu *LocationUpdate) RemoveSurvey(s ...*Survey) *LocationUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lu.RemoveSurveyIDs(ids...)
}

// RemoveWifiScanIDs removes the wifi_scan edge to SurveyWiFiScan by ids.
func (lu *LocationUpdate) RemoveWifiScanIDs(ids ...int) *LocationUpdate {
	lu.mutation.RemoveWifiScanIDs(ids...)
	return lu
}

// RemoveWifiScan removes wifi_scan edges to SurveyWiFiScan.
func (lu *LocationUpdate) RemoveWifiScan(s ...*SurveyWiFiScan) *LocationUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lu.RemoveWifiScanIDs(ids...)
}

// RemoveCellScanIDs removes the cell_scan edge to SurveyCellScan by ids.
func (lu *LocationUpdate) RemoveCellScanIDs(ids ...int) *LocationUpdate {
	lu.mutation.RemoveCellScanIDs(ids...)
	return lu
}

// RemoveCellScan removes cell_scan edges to SurveyCellScan.
func (lu *LocationUpdate) RemoveCellScan(s ...*SurveyCellScan) *LocationUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lu.RemoveCellScanIDs(ids...)
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrder by ids.
func (lu *LocationUpdate) RemoveWorkOrderIDs(ids ...int) *LocationUpdate {
	lu.mutation.RemoveWorkOrderIDs(ids...)
	return lu
}

// RemoveWorkOrders removes work_orders edges to WorkOrder.
func (lu *LocationUpdate) RemoveWorkOrders(w ...*WorkOrder) *LocationUpdate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return lu.RemoveWorkOrderIDs(ids...)
}

// RemoveFloorPlanIDs removes the floor_plans edge to FloorPlan by ids.
func (lu *LocationUpdate) RemoveFloorPlanIDs(ids ...int) *LocationUpdate {
	lu.mutation.RemoveFloorPlanIDs(ids...)
	return lu
}

// RemoveFloorPlans removes floor_plans edges to FloorPlan.
func (lu *LocationUpdate) RemoveFloorPlans(f ...*FloorPlan) *LocationUpdate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return lu.RemoveFloorPlanIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (lu *LocationUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := lu.mutation.UpdateTime(); !ok {
		v := location.UpdateDefaultUpdateTime()
		lu.mutation.SetUpdateTime(v)
	}
	if v, ok := lu.mutation.Name(); ok {
		if err := location.NameValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if v, ok := lu.mutation.Latitude(); ok {
		if err := location.LatitudeValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"latitude\": %v", err)
		}
	}
	if v, ok := lu.mutation.Longitude(); ok {
		if err := location.LongitudeValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"longitude\": %v", err)
		}
	}

	if _, ok := lu.mutation.TypeID(); lu.mutation.TypeCleared() && !ok {
		return 0, errors.New("ent: clearing a unique edge \"type\"")
	}

	var (
		err      error
		affected int
	)
	if len(lu.hooks) == 0 {
		affected, err = lu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*LocationMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			lu.mutation = mutation
			affected, err = lu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(lu.hooks) - 1; i >= 0; i-- {
			mut = lu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, lu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (lu *LocationUpdate) SaveX(ctx context.Context) int {
	affected, err := lu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (lu *LocationUpdate) Exec(ctx context.Context) error {
	_, err := lu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lu *LocationUpdate) ExecX(ctx context.Context) {
	if err := lu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (lu *LocationUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   location.Table,
			Columns: location.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: location.FieldID,
			},
		},
	}
	if ps := lu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := lu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: location.FieldUpdateTime,
		})
	}
	if value, ok := lu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: location.FieldName,
		})
	}
	if value, ok := lu.mutation.ExternalID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: location.FieldExternalID,
		})
	}
	if lu.mutation.ExternalIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: location.FieldExternalID,
		})
	}
	if value, ok := lu.mutation.Latitude(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: location.FieldLatitude,
		})
	}
	if value, ok := lu.mutation.AddedLatitude(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: location.FieldLatitude,
		})
	}
	if value, ok := lu.mutation.Longitude(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: location.FieldLongitude,
		})
	}
	if value, ok := lu.mutation.AddedLongitude(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: location.FieldLongitude,
		})
	}
	if value, ok := lu.mutation.SiteSurveyNeeded(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: location.FieldSiteSurveyNeeded,
		})
	}
	if lu.mutation.SiteSurveyNeededCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Column: location.FieldSiteSurveyNeeded,
		})
	}
	if lu.mutation.TypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   location.TypeTable,
			Columns: []string{location.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: locationtype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.mutation.TypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   location.TypeTable,
			Columns: []string{location.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: locationtype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if lu.mutation.ParentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   location.ParentTable,
			Columns: []string{location.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: location.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.mutation.ParentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   location.ParentTable,
			Columns: []string{location.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: location.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := lu.mutation.RemovedChildrenIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.ChildrenTable,
			Columns: []string{location.ChildrenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: location.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.mutation.ChildrenIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.ChildrenTable,
			Columns: []string{location.ChildrenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: location.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := lu.mutation.RemovedFilesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.FilesTable,
			Columns: []string{location.FilesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: file.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.mutation.FilesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.FilesTable,
			Columns: []string{location.FilesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: file.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := lu.mutation.RemovedHyperlinksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.HyperlinksTable,
			Columns: []string{location.HyperlinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: hyperlink.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.mutation.HyperlinksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.HyperlinksTable,
			Columns: []string{location.HyperlinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: hyperlink.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := lu.mutation.RemovedEquipmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.EquipmentTable,
			Columns: []string{location.EquipmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.mutation.EquipmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.EquipmentTable,
			Columns: []string{location.EquipmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := lu.mutation.RemovedPropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.PropertiesTable,
			Columns: []string{location.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: property.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.mutation.PropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.PropertiesTable,
			Columns: []string{location.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: property.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := lu.mutation.RemovedSurveyIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.SurveyTable,
			Columns: []string{location.SurveyColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: survey.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.mutation.SurveyIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.SurveyTable,
			Columns: []string{location.SurveyColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: survey.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := lu.mutation.RemovedWifiScanIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.WifiScanTable,
			Columns: []string{location.WifiScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveywifiscan.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.mutation.WifiScanIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.WifiScanTable,
			Columns: []string{location.WifiScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveywifiscan.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := lu.mutation.RemovedCellScanIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.CellScanTable,
			Columns: []string{location.CellScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveycellscan.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.mutation.CellScanIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.CellScanTable,
			Columns: []string{location.CellScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveycellscan.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := lu.mutation.RemovedWorkOrdersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.WorkOrdersTable,
			Columns: []string{location.WorkOrdersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorder.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.mutation.WorkOrdersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.WorkOrdersTable,
			Columns: []string{location.WorkOrdersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorder.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := lu.mutation.RemovedFloorPlansIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.FloorPlansTable,
			Columns: []string{location.FloorPlansColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: floorplan.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.mutation.FloorPlansIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.FloorPlansTable,
			Columns: []string{location.FloorPlansColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: floorplan.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, lu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{location.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// LocationUpdateOne is the builder for updating a single Location entity.
type LocationUpdateOne struct {
	config
	hooks    []Hook
	mutation *LocationMutation
}

// SetName sets the name field.
func (luo *LocationUpdateOne) SetName(s string) *LocationUpdateOne {
	luo.mutation.SetName(s)
	return luo
}

// SetExternalID sets the external_id field.
func (luo *LocationUpdateOne) SetExternalID(s string) *LocationUpdateOne {
	luo.mutation.SetExternalID(s)
	return luo
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (luo *LocationUpdateOne) SetNillableExternalID(s *string) *LocationUpdateOne {
	if s != nil {
		luo.SetExternalID(*s)
	}
	return luo
}

// ClearExternalID clears the value of external_id.
func (luo *LocationUpdateOne) ClearExternalID() *LocationUpdateOne {
	luo.mutation.ClearExternalID()
	return luo
}

// SetLatitude sets the latitude field.
func (luo *LocationUpdateOne) SetLatitude(f float64) *LocationUpdateOne {
	luo.mutation.ResetLatitude()
	luo.mutation.SetLatitude(f)
	return luo
}

// SetNillableLatitude sets the latitude field if the given value is not nil.
func (luo *LocationUpdateOne) SetNillableLatitude(f *float64) *LocationUpdateOne {
	if f != nil {
		luo.SetLatitude(*f)
	}
	return luo
}

// AddLatitude adds f to latitude.
func (luo *LocationUpdateOne) AddLatitude(f float64) *LocationUpdateOne {
	luo.mutation.AddLatitude(f)
	return luo
}

// SetLongitude sets the longitude field.
func (luo *LocationUpdateOne) SetLongitude(f float64) *LocationUpdateOne {
	luo.mutation.ResetLongitude()
	luo.mutation.SetLongitude(f)
	return luo
}

// SetNillableLongitude sets the longitude field if the given value is not nil.
func (luo *LocationUpdateOne) SetNillableLongitude(f *float64) *LocationUpdateOne {
	if f != nil {
		luo.SetLongitude(*f)
	}
	return luo
}

// AddLongitude adds f to longitude.
func (luo *LocationUpdateOne) AddLongitude(f float64) *LocationUpdateOne {
	luo.mutation.AddLongitude(f)
	return luo
}

// SetSiteSurveyNeeded sets the site_survey_needed field.
func (luo *LocationUpdateOne) SetSiteSurveyNeeded(b bool) *LocationUpdateOne {
	luo.mutation.SetSiteSurveyNeeded(b)
	return luo
}

// SetNillableSiteSurveyNeeded sets the site_survey_needed field if the given value is not nil.
func (luo *LocationUpdateOne) SetNillableSiteSurveyNeeded(b *bool) *LocationUpdateOne {
	if b != nil {
		luo.SetSiteSurveyNeeded(*b)
	}
	return luo
}

// ClearSiteSurveyNeeded clears the value of site_survey_needed.
func (luo *LocationUpdateOne) ClearSiteSurveyNeeded() *LocationUpdateOne {
	luo.mutation.ClearSiteSurveyNeeded()
	return luo
}

// SetTypeID sets the type edge to LocationType by id.
func (luo *LocationUpdateOne) SetTypeID(id int) *LocationUpdateOne {
	luo.mutation.SetTypeID(id)
	return luo
}

// SetType sets the type edge to LocationType.
func (luo *LocationUpdateOne) SetType(l *LocationType) *LocationUpdateOne {
	return luo.SetTypeID(l.ID)
}

// SetParentID sets the parent edge to Location by id.
func (luo *LocationUpdateOne) SetParentID(id int) *LocationUpdateOne {
	luo.mutation.SetParentID(id)
	return luo
}

// SetNillableParentID sets the parent edge to Location by id if the given value is not nil.
func (luo *LocationUpdateOne) SetNillableParentID(id *int) *LocationUpdateOne {
	if id != nil {
		luo = luo.SetParentID(*id)
	}
	return luo
}

// SetParent sets the parent edge to Location.
func (luo *LocationUpdateOne) SetParent(l *Location) *LocationUpdateOne {
	return luo.SetParentID(l.ID)
}

// AddChildIDs adds the children edge to Location by ids.
func (luo *LocationUpdateOne) AddChildIDs(ids ...int) *LocationUpdateOne {
	luo.mutation.AddChildIDs(ids...)
	return luo
}

// AddChildren adds the children edges to Location.
func (luo *LocationUpdateOne) AddChildren(l ...*Location) *LocationUpdateOne {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return luo.AddChildIDs(ids...)
}

// AddFileIDs adds the files edge to File by ids.
func (luo *LocationUpdateOne) AddFileIDs(ids ...int) *LocationUpdateOne {
	luo.mutation.AddFileIDs(ids...)
	return luo
}

// AddFiles adds the files edges to File.
func (luo *LocationUpdateOne) AddFiles(f ...*File) *LocationUpdateOne {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return luo.AddFileIDs(ids...)
}

// AddHyperlinkIDs adds the hyperlinks edge to Hyperlink by ids.
func (luo *LocationUpdateOne) AddHyperlinkIDs(ids ...int) *LocationUpdateOne {
	luo.mutation.AddHyperlinkIDs(ids...)
	return luo
}

// AddHyperlinks adds the hyperlinks edges to Hyperlink.
func (luo *LocationUpdateOne) AddHyperlinks(h ...*Hyperlink) *LocationUpdateOne {
	ids := make([]int, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return luo.AddHyperlinkIDs(ids...)
}

// AddEquipmentIDs adds the equipment edge to Equipment by ids.
func (luo *LocationUpdateOne) AddEquipmentIDs(ids ...int) *LocationUpdateOne {
	luo.mutation.AddEquipmentIDs(ids...)
	return luo
}

// AddEquipment adds the equipment edges to Equipment.
func (luo *LocationUpdateOne) AddEquipment(e ...*Equipment) *LocationUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return luo.AddEquipmentIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (luo *LocationUpdateOne) AddPropertyIDs(ids ...int) *LocationUpdateOne {
	luo.mutation.AddPropertyIDs(ids...)
	return luo
}

// AddProperties adds the properties edges to Property.
func (luo *LocationUpdateOne) AddProperties(p ...*Property) *LocationUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return luo.AddPropertyIDs(ids...)
}

// AddSurveyIDs adds the survey edge to Survey by ids.
func (luo *LocationUpdateOne) AddSurveyIDs(ids ...int) *LocationUpdateOne {
	luo.mutation.AddSurveyIDs(ids...)
	return luo
}

// AddSurvey adds the survey edges to Survey.
func (luo *LocationUpdateOne) AddSurvey(s ...*Survey) *LocationUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return luo.AddSurveyIDs(ids...)
}

// AddWifiScanIDs adds the wifi_scan edge to SurveyWiFiScan by ids.
func (luo *LocationUpdateOne) AddWifiScanIDs(ids ...int) *LocationUpdateOne {
	luo.mutation.AddWifiScanIDs(ids...)
	return luo
}

// AddWifiScan adds the wifi_scan edges to SurveyWiFiScan.
func (luo *LocationUpdateOne) AddWifiScan(s ...*SurveyWiFiScan) *LocationUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return luo.AddWifiScanIDs(ids...)
}

// AddCellScanIDs adds the cell_scan edge to SurveyCellScan by ids.
func (luo *LocationUpdateOne) AddCellScanIDs(ids ...int) *LocationUpdateOne {
	luo.mutation.AddCellScanIDs(ids...)
	return luo
}

// AddCellScan adds the cell_scan edges to SurveyCellScan.
func (luo *LocationUpdateOne) AddCellScan(s ...*SurveyCellScan) *LocationUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return luo.AddCellScanIDs(ids...)
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (luo *LocationUpdateOne) AddWorkOrderIDs(ids ...int) *LocationUpdateOne {
	luo.mutation.AddWorkOrderIDs(ids...)
	return luo
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (luo *LocationUpdateOne) AddWorkOrders(w ...*WorkOrder) *LocationUpdateOne {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return luo.AddWorkOrderIDs(ids...)
}

// AddFloorPlanIDs adds the floor_plans edge to FloorPlan by ids.
func (luo *LocationUpdateOne) AddFloorPlanIDs(ids ...int) *LocationUpdateOne {
	luo.mutation.AddFloorPlanIDs(ids...)
	return luo
}

// AddFloorPlans adds the floor_plans edges to FloorPlan.
func (luo *LocationUpdateOne) AddFloorPlans(f ...*FloorPlan) *LocationUpdateOne {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return luo.AddFloorPlanIDs(ids...)
}

// ClearType clears the type edge to LocationType.
func (luo *LocationUpdateOne) ClearType() *LocationUpdateOne {
	luo.mutation.ClearType()
	return luo
}

// ClearParent clears the parent edge to Location.
func (luo *LocationUpdateOne) ClearParent() *LocationUpdateOne {
	luo.mutation.ClearParent()
	return luo
}

// RemoveChildIDs removes the children edge to Location by ids.
func (luo *LocationUpdateOne) RemoveChildIDs(ids ...int) *LocationUpdateOne {
	luo.mutation.RemoveChildIDs(ids...)
	return luo
}

// RemoveChildren removes children edges to Location.
func (luo *LocationUpdateOne) RemoveChildren(l ...*Location) *LocationUpdateOne {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return luo.RemoveChildIDs(ids...)
}

// RemoveFileIDs removes the files edge to File by ids.
func (luo *LocationUpdateOne) RemoveFileIDs(ids ...int) *LocationUpdateOne {
	luo.mutation.RemoveFileIDs(ids...)
	return luo
}

// RemoveFiles removes files edges to File.
func (luo *LocationUpdateOne) RemoveFiles(f ...*File) *LocationUpdateOne {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return luo.RemoveFileIDs(ids...)
}

// RemoveHyperlinkIDs removes the hyperlinks edge to Hyperlink by ids.
func (luo *LocationUpdateOne) RemoveHyperlinkIDs(ids ...int) *LocationUpdateOne {
	luo.mutation.RemoveHyperlinkIDs(ids...)
	return luo
}

// RemoveHyperlinks removes hyperlinks edges to Hyperlink.
func (luo *LocationUpdateOne) RemoveHyperlinks(h ...*Hyperlink) *LocationUpdateOne {
	ids := make([]int, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return luo.RemoveHyperlinkIDs(ids...)
}

// RemoveEquipmentIDs removes the equipment edge to Equipment by ids.
func (luo *LocationUpdateOne) RemoveEquipmentIDs(ids ...int) *LocationUpdateOne {
	luo.mutation.RemoveEquipmentIDs(ids...)
	return luo
}

// RemoveEquipment removes equipment edges to Equipment.
func (luo *LocationUpdateOne) RemoveEquipment(e ...*Equipment) *LocationUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return luo.RemoveEquipmentIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (luo *LocationUpdateOne) RemovePropertyIDs(ids ...int) *LocationUpdateOne {
	luo.mutation.RemovePropertyIDs(ids...)
	return luo
}

// RemoveProperties removes properties edges to Property.
func (luo *LocationUpdateOne) RemoveProperties(p ...*Property) *LocationUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return luo.RemovePropertyIDs(ids...)
}

// RemoveSurveyIDs removes the survey edge to Survey by ids.
func (luo *LocationUpdateOne) RemoveSurveyIDs(ids ...int) *LocationUpdateOne {
	luo.mutation.RemoveSurveyIDs(ids...)
	return luo
}

// RemoveSurvey removes survey edges to Survey.
func (luo *LocationUpdateOne) RemoveSurvey(s ...*Survey) *LocationUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return luo.RemoveSurveyIDs(ids...)
}

// RemoveWifiScanIDs removes the wifi_scan edge to SurveyWiFiScan by ids.
func (luo *LocationUpdateOne) RemoveWifiScanIDs(ids ...int) *LocationUpdateOne {
	luo.mutation.RemoveWifiScanIDs(ids...)
	return luo
}

// RemoveWifiScan removes wifi_scan edges to SurveyWiFiScan.
func (luo *LocationUpdateOne) RemoveWifiScan(s ...*SurveyWiFiScan) *LocationUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return luo.RemoveWifiScanIDs(ids...)
}

// RemoveCellScanIDs removes the cell_scan edge to SurveyCellScan by ids.
func (luo *LocationUpdateOne) RemoveCellScanIDs(ids ...int) *LocationUpdateOne {
	luo.mutation.RemoveCellScanIDs(ids...)
	return luo
}

// RemoveCellScan removes cell_scan edges to SurveyCellScan.
func (luo *LocationUpdateOne) RemoveCellScan(s ...*SurveyCellScan) *LocationUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return luo.RemoveCellScanIDs(ids...)
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrder by ids.
func (luo *LocationUpdateOne) RemoveWorkOrderIDs(ids ...int) *LocationUpdateOne {
	luo.mutation.RemoveWorkOrderIDs(ids...)
	return luo
}

// RemoveWorkOrders removes work_orders edges to WorkOrder.
func (luo *LocationUpdateOne) RemoveWorkOrders(w ...*WorkOrder) *LocationUpdateOne {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return luo.RemoveWorkOrderIDs(ids...)
}

// RemoveFloorPlanIDs removes the floor_plans edge to FloorPlan by ids.
func (luo *LocationUpdateOne) RemoveFloorPlanIDs(ids ...int) *LocationUpdateOne {
	luo.mutation.RemoveFloorPlanIDs(ids...)
	return luo
}

// RemoveFloorPlans removes floor_plans edges to FloorPlan.
func (luo *LocationUpdateOne) RemoveFloorPlans(f ...*FloorPlan) *LocationUpdateOne {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return luo.RemoveFloorPlanIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (luo *LocationUpdateOne) Save(ctx context.Context) (*Location, error) {
	if _, ok := luo.mutation.UpdateTime(); !ok {
		v := location.UpdateDefaultUpdateTime()
		luo.mutation.SetUpdateTime(v)
	}
	if v, ok := luo.mutation.Name(); ok {
		if err := location.NameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if v, ok := luo.mutation.Latitude(); ok {
		if err := location.LatitudeValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"latitude\": %v", err)
		}
	}
	if v, ok := luo.mutation.Longitude(); ok {
		if err := location.LongitudeValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"longitude\": %v", err)
		}
	}

	if _, ok := luo.mutation.TypeID(); luo.mutation.TypeCleared() && !ok {
		return nil, errors.New("ent: clearing a unique edge \"type\"")
	}

	var (
		err  error
		node *Location
	)
	if len(luo.hooks) == 0 {
		node, err = luo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*LocationMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			luo.mutation = mutation
			node, err = luo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(luo.hooks) - 1; i >= 0; i-- {
			mut = luo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, luo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (luo *LocationUpdateOne) SaveX(ctx context.Context) *Location {
	l, err := luo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return l
}

// Exec executes the query on the entity.
func (luo *LocationUpdateOne) Exec(ctx context.Context) error {
	_, err := luo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (luo *LocationUpdateOne) ExecX(ctx context.Context) {
	if err := luo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (luo *LocationUpdateOne) sqlSave(ctx context.Context) (l *Location, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   location.Table,
			Columns: location.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: location.FieldID,
			},
		},
	}
	id, ok := luo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing Location.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := luo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: location.FieldUpdateTime,
		})
	}
	if value, ok := luo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: location.FieldName,
		})
	}
	if value, ok := luo.mutation.ExternalID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: location.FieldExternalID,
		})
	}
	if luo.mutation.ExternalIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: location.FieldExternalID,
		})
	}
	if value, ok := luo.mutation.Latitude(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: location.FieldLatitude,
		})
	}
	if value, ok := luo.mutation.AddedLatitude(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: location.FieldLatitude,
		})
	}
	if value, ok := luo.mutation.Longitude(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: location.FieldLongitude,
		})
	}
	if value, ok := luo.mutation.AddedLongitude(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: location.FieldLongitude,
		})
	}
	if value, ok := luo.mutation.SiteSurveyNeeded(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: location.FieldSiteSurveyNeeded,
		})
	}
	if luo.mutation.SiteSurveyNeededCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Column: location.FieldSiteSurveyNeeded,
		})
	}
	if luo.mutation.TypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   location.TypeTable,
			Columns: []string{location.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: locationtype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.mutation.TypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   location.TypeTable,
			Columns: []string{location.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: locationtype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if luo.mutation.ParentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   location.ParentTable,
			Columns: []string{location.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: location.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.mutation.ParentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   location.ParentTable,
			Columns: []string{location.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: location.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := luo.mutation.RemovedChildrenIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.ChildrenTable,
			Columns: []string{location.ChildrenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: location.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.mutation.ChildrenIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.ChildrenTable,
			Columns: []string{location.ChildrenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: location.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := luo.mutation.RemovedFilesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.FilesTable,
			Columns: []string{location.FilesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: file.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.mutation.FilesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.FilesTable,
			Columns: []string{location.FilesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: file.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := luo.mutation.RemovedHyperlinksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.HyperlinksTable,
			Columns: []string{location.HyperlinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: hyperlink.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.mutation.HyperlinksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.HyperlinksTable,
			Columns: []string{location.HyperlinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: hyperlink.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := luo.mutation.RemovedEquipmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.EquipmentTable,
			Columns: []string{location.EquipmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.mutation.EquipmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.EquipmentTable,
			Columns: []string{location.EquipmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := luo.mutation.RemovedPropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.PropertiesTable,
			Columns: []string{location.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: property.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.mutation.PropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.PropertiesTable,
			Columns: []string{location.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: property.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := luo.mutation.RemovedSurveyIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.SurveyTable,
			Columns: []string{location.SurveyColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: survey.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.mutation.SurveyIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.SurveyTable,
			Columns: []string{location.SurveyColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: survey.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := luo.mutation.RemovedWifiScanIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.WifiScanTable,
			Columns: []string{location.WifiScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveywifiscan.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.mutation.WifiScanIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.WifiScanTable,
			Columns: []string{location.WifiScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveywifiscan.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := luo.mutation.RemovedCellScanIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.CellScanTable,
			Columns: []string{location.CellScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveycellscan.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.mutation.CellScanIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.CellScanTable,
			Columns: []string{location.CellScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveycellscan.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := luo.mutation.RemovedWorkOrdersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.WorkOrdersTable,
			Columns: []string{location.WorkOrdersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorder.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.mutation.WorkOrdersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.WorkOrdersTable,
			Columns: []string{location.WorkOrdersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorder.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := luo.mutation.RemovedFloorPlansIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.FloorPlansTable,
			Columns: []string{location.FloorPlansColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: floorplan.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.mutation.FloorPlansIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.FloorPlansTable,
			Columns: []string{location.FloorPlansColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: floorplan.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	l = &Location{config: luo.config}
	_spec.Assign = l.assignValues
	_spec.ScanValues = l.scanValues()
	if err = sqlgraph.UpdateNode(ctx, luo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{location.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return l, nil
}
