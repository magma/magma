// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/floorplan"
	"github.com/facebookincubator/symphony/graph/ent/hyperlink"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/survey"
	"github.com/facebookincubator/symphony/graph/ent/surveycellscan"
	"github.com/facebookincubator/symphony/graph/ent/surveywifiscan"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// LocationCreate is the builder for creating a Location entity.
type LocationCreate struct {
	config
	mutation *LocationMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (lc *LocationCreate) SetCreateTime(t time.Time) *LocationCreate {
	lc.mutation.SetCreateTime(t)
	return lc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (lc *LocationCreate) SetNillableCreateTime(t *time.Time) *LocationCreate {
	if t != nil {
		lc.SetCreateTime(*t)
	}
	return lc
}

// SetUpdateTime sets the update_time field.
func (lc *LocationCreate) SetUpdateTime(t time.Time) *LocationCreate {
	lc.mutation.SetUpdateTime(t)
	return lc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (lc *LocationCreate) SetNillableUpdateTime(t *time.Time) *LocationCreate {
	if t != nil {
		lc.SetUpdateTime(*t)
	}
	return lc
}

// SetName sets the name field.
func (lc *LocationCreate) SetName(s string) *LocationCreate {
	lc.mutation.SetName(s)
	return lc
}

// SetExternalID sets the external_id field.
func (lc *LocationCreate) SetExternalID(s string) *LocationCreate {
	lc.mutation.SetExternalID(s)
	return lc
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (lc *LocationCreate) SetNillableExternalID(s *string) *LocationCreate {
	if s != nil {
		lc.SetExternalID(*s)
	}
	return lc
}

// SetLatitude sets the latitude field.
func (lc *LocationCreate) SetLatitude(f float64) *LocationCreate {
	lc.mutation.SetLatitude(f)
	return lc
}

// SetNillableLatitude sets the latitude field if the given value is not nil.
func (lc *LocationCreate) SetNillableLatitude(f *float64) *LocationCreate {
	if f != nil {
		lc.SetLatitude(*f)
	}
	return lc
}

// SetLongitude sets the longitude field.
func (lc *LocationCreate) SetLongitude(f float64) *LocationCreate {
	lc.mutation.SetLongitude(f)
	return lc
}

// SetNillableLongitude sets the longitude field if the given value is not nil.
func (lc *LocationCreate) SetNillableLongitude(f *float64) *LocationCreate {
	if f != nil {
		lc.SetLongitude(*f)
	}
	return lc
}

// SetSiteSurveyNeeded sets the site_survey_needed field.
func (lc *LocationCreate) SetSiteSurveyNeeded(b bool) *LocationCreate {
	lc.mutation.SetSiteSurveyNeeded(b)
	return lc
}

// SetNillableSiteSurveyNeeded sets the site_survey_needed field if the given value is not nil.
func (lc *LocationCreate) SetNillableSiteSurveyNeeded(b *bool) *LocationCreate {
	if b != nil {
		lc.SetSiteSurveyNeeded(*b)
	}
	return lc
}

// SetTypeID sets the type edge to LocationType by id.
func (lc *LocationCreate) SetTypeID(id int) *LocationCreate {
	lc.mutation.SetTypeID(id)
	return lc
}

// SetType sets the type edge to LocationType.
func (lc *LocationCreate) SetType(l *LocationType) *LocationCreate {
	return lc.SetTypeID(l.ID)
}

// SetParentID sets the parent edge to Location by id.
func (lc *LocationCreate) SetParentID(id int) *LocationCreate {
	lc.mutation.SetParentID(id)
	return lc
}

// SetNillableParentID sets the parent edge to Location by id if the given value is not nil.
func (lc *LocationCreate) SetNillableParentID(id *int) *LocationCreate {
	if id != nil {
		lc = lc.SetParentID(*id)
	}
	return lc
}

// SetParent sets the parent edge to Location.
func (lc *LocationCreate) SetParent(l *Location) *LocationCreate {
	return lc.SetParentID(l.ID)
}

// AddChildIDs adds the children edge to Location by ids.
func (lc *LocationCreate) AddChildIDs(ids ...int) *LocationCreate {
	lc.mutation.AddChildIDs(ids...)
	return lc
}

// AddChildren adds the children edges to Location.
func (lc *LocationCreate) AddChildren(l ...*Location) *LocationCreate {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return lc.AddChildIDs(ids...)
}

// AddFileIDs adds the files edge to File by ids.
func (lc *LocationCreate) AddFileIDs(ids ...int) *LocationCreate {
	lc.mutation.AddFileIDs(ids...)
	return lc
}

// AddFiles adds the files edges to File.
func (lc *LocationCreate) AddFiles(f ...*File) *LocationCreate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return lc.AddFileIDs(ids...)
}

// AddHyperlinkIDs adds the hyperlinks edge to Hyperlink by ids.
func (lc *LocationCreate) AddHyperlinkIDs(ids ...int) *LocationCreate {
	lc.mutation.AddHyperlinkIDs(ids...)
	return lc
}

// AddHyperlinks adds the hyperlinks edges to Hyperlink.
func (lc *LocationCreate) AddHyperlinks(h ...*Hyperlink) *LocationCreate {
	ids := make([]int, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return lc.AddHyperlinkIDs(ids...)
}

// AddEquipmentIDs adds the equipment edge to Equipment by ids.
func (lc *LocationCreate) AddEquipmentIDs(ids ...int) *LocationCreate {
	lc.mutation.AddEquipmentIDs(ids...)
	return lc
}

// AddEquipment adds the equipment edges to Equipment.
func (lc *LocationCreate) AddEquipment(e ...*Equipment) *LocationCreate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return lc.AddEquipmentIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (lc *LocationCreate) AddPropertyIDs(ids ...int) *LocationCreate {
	lc.mutation.AddPropertyIDs(ids...)
	return lc
}

// AddProperties adds the properties edges to Property.
func (lc *LocationCreate) AddProperties(p ...*Property) *LocationCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return lc.AddPropertyIDs(ids...)
}

// AddSurveyIDs adds the survey edge to Survey by ids.
func (lc *LocationCreate) AddSurveyIDs(ids ...int) *LocationCreate {
	lc.mutation.AddSurveyIDs(ids...)
	return lc
}

// AddSurvey adds the survey edges to Survey.
func (lc *LocationCreate) AddSurvey(s ...*Survey) *LocationCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lc.AddSurveyIDs(ids...)
}

// AddWifiScanIDs adds the wifi_scan edge to SurveyWiFiScan by ids.
func (lc *LocationCreate) AddWifiScanIDs(ids ...int) *LocationCreate {
	lc.mutation.AddWifiScanIDs(ids...)
	return lc
}

// AddWifiScan adds the wifi_scan edges to SurveyWiFiScan.
func (lc *LocationCreate) AddWifiScan(s ...*SurveyWiFiScan) *LocationCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lc.AddWifiScanIDs(ids...)
}

// AddCellScanIDs adds the cell_scan edge to SurveyCellScan by ids.
func (lc *LocationCreate) AddCellScanIDs(ids ...int) *LocationCreate {
	lc.mutation.AddCellScanIDs(ids...)
	return lc
}

// AddCellScan adds the cell_scan edges to SurveyCellScan.
func (lc *LocationCreate) AddCellScan(s ...*SurveyCellScan) *LocationCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lc.AddCellScanIDs(ids...)
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (lc *LocationCreate) AddWorkOrderIDs(ids ...int) *LocationCreate {
	lc.mutation.AddWorkOrderIDs(ids...)
	return lc
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (lc *LocationCreate) AddWorkOrders(w ...*WorkOrder) *LocationCreate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return lc.AddWorkOrderIDs(ids...)
}

// AddFloorPlanIDs adds the floor_plans edge to FloorPlan by ids.
func (lc *LocationCreate) AddFloorPlanIDs(ids ...int) *LocationCreate {
	lc.mutation.AddFloorPlanIDs(ids...)
	return lc
}

// AddFloorPlans adds the floor_plans edges to FloorPlan.
func (lc *LocationCreate) AddFloorPlans(f ...*FloorPlan) *LocationCreate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return lc.AddFloorPlanIDs(ids...)
}

// Save creates the Location in the database.
func (lc *LocationCreate) Save(ctx context.Context) (*Location, error) {
	if _, ok := lc.mutation.CreateTime(); !ok {
		v := location.DefaultCreateTime()
		lc.mutation.SetCreateTime(v)
	}
	if _, ok := lc.mutation.UpdateTime(); !ok {
		v := location.DefaultUpdateTime()
		lc.mutation.SetUpdateTime(v)
	}
	if _, ok := lc.mutation.Name(); !ok {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if v, ok := lc.mutation.Name(); ok {
		if err := location.NameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if _, ok := lc.mutation.Latitude(); !ok {
		v := location.DefaultLatitude
		lc.mutation.SetLatitude(v)
	}
	if v, ok := lc.mutation.Latitude(); ok {
		if err := location.LatitudeValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"latitude\": %v", err)
		}
	}
	if _, ok := lc.mutation.Longitude(); !ok {
		v := location.DefaultLongitude
		lc.mutation.SetLongitude(v)
	}
	if v, ok := lc.mutation.Longitude(); ok {
		if err := location.LongitudeValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"longitude\": %v", err)
		}
	}
	if _, ok := lc.mutation.SiteSurveyNeeded(); !ok {
		v := location.DefaultSiteSurveyNeeded
		lc.mutation.SetSiteSurveyNeeded(v)
	}
	if _, ok := lc.mutation.TypeID(); !ok {
		return nil, errors.New("ent: missing required edge \"type\"")
	}
	var (
		err  error
		node *Location
	)
	if len(lc.hooks) == 0 {
		node, err = lc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*LocationMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			lc.mutation = mutation
			node, err = lc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(lc.hooks) - 1; i >= 0; i-- {
			mut = lc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, lc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (lc *LocationCreate) SaveX(ctx context.Context) *Location {
	v, err := lc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (lc *LocationCreate) sqlSave(ctx context.Context) (*Location, error) {
	var (
		l     = &Location{config: lc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: location.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: location.FieldID,
			},
		}
	)
	if value, ok := lc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: location.FieldCreateTime,
		})
		l.CreateTime = value
	}
	if value, ok := lc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: location.FieldUpdateTime,
		})
		l.UpdateTime = value
	}
	if value, ok := lc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: location.FieldName,
		})
		l.Name = value
	}
	if value, ok := lc.mutation.ExternalID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: location.FieldExternalID,
		})
		l.ExternalID = value
	}
	if value, ok := lc.mutation.Latitude(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: location.FieldLatitude,
		})
		l.Latitude = value
	}
	if value, ok := lc.mutation.Longitude(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: location.FieldLongitude,
		})
		l.Longitude = value
	}
	if value, ok := lc.mutation.SiteSurveyNeeded(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: location.FieldSiteSurveyNeeded,
		})
		l.SiteSurveyNeeded = value
	}
	if nodes := lc.mutation.TypeIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := lc.mutation.ParentIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := lc.mutation.ChildrenIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := lc.mutation.FilesIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := lc.mutation.HyperlinksIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := lc.mutation.EquipmentIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := lc.mutation.PropertiesIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := lc.mutation.SurveyIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := lc.mutation.WifiScanIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := lc.mutation.CellScanIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := lc.mutation.WorkOrdersIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := lc.mutation.FloorPlansIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, lc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	l.ID = int(id)
	return l, nil
}
