// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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
	create_time        *time.Time
	update_time        *time.Time
	name               *string
	external_id        *string
	latitude           *float64
	longitude          *float64
	site_survey_needed *bool
	_type              map[string]struct{}
	parent             map[string]struct{}
	children           map[string]struct{}
	files              map[string]struct{}
	hyperlinks         map[string]struct{}
	equipment          map[string]struct{}
	properties         map[string]struct{}
	survey             map[string]struct{}
	wifi_scan          map[string]struct{}
	cell_scan          map[string]struct{}
	work_orders        map[string]struct{}
	floor_plans        map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (lc *LocationCreate) SetCreateTime(t time.Time) *LocationCreate {
	lc.create_time = &t
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
	lc.update_time = &t
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
	lc.name = &s
	return lc
}

// SetExternalID sets the external_id field.
func (lc *LocationCreate) SetExternalID(s string) *LocationCreate {
	lc.external_id = &s
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
	lc.latitude = &f
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
	lc.longitude = &f
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
	lc.site_survey_needed = &b
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
func (lc *LocationCreate) SetTypeID(id string) *LocationCreate {
	if lc._type == nil {
		lc._type = make(map[string]struct{})
	}
	lc._type[id] = struct{}{}
	return lc
}

// SetType sets the type edge to LocationType.
func (lc *LocationCreate) SetType(l *LocationType) *LocationCreate {
	return lc.SetTypeID(l.ID)
}

// SetParentID sets the parent edge to Location by id.
func (lc *LocationCreate) SetParentID(id string) *LocationCreate {
	if lc.parent == nil {
		lc.parent = make(map[string]struct{})
	}
	lc.parent[id] = struct{}{}
	return lc
}

// SetNillableParentID sets the parent edge to Location by id if the given value is not nil.
func (lc *LocationCreate) SetNillableParentID(id *string) *LocationCreate {
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
func (lc *LocationCreate) AddChildIDs(ids ...string) *LocationCreate {
	if lc.children == nil {
		lc.children = make(map[string]struct{})
	}
	for i := range ids {
		lc.children[ids[i]] = struct{}{}
	}
	return lc
}

// AddChildren adds the children edges to Location.
func (lc *LocationCreate) AddChildren(l ...*Location) *LocationCreate {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return lc.AddChildIDs(ids...)
}

// AddFileIDs adds the files edge to File by ids.
func (lc *LocationCreate) AddFileIDs(ids ...string) *LocationCreate {
	if lc.files == nil {
		lc.files = make(map[string]struct{})
	}
	for i := range ids {
		lc.files[ids[i]] = struct{}{}
	}
	return lc
}

// AddFiles adds the files edges to File.
func (lc *LocationCreate) AddFiles(f ...*File) *LocationCreate {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return lc.AddFileIDs(ids...)
}

// AddHyperlinkIDs adds the hyperlinks edge to Hyperlink by ids.
func (lc *LocationCreate) AddHyperlinkIDs(ids ...string) *LocationCreate {
	if lc.hyperlinks == nil {
		lc.hyperlinks = make(map[string]struct{})
	}
	for i := range ids {
		lc.hyperlinks[ids[i]] = struct{}{}
	}
	return lc
}

// AddHyperlinks adds the hyperlinks edges to Hyperlink.
func (lc *LocationCreate) AddHyperlinks(h ...*Hyperlink) *LocationCreate {
	ids := make([]string, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return lc.AddHyperlinkIDs(ids...)
}

// AddEquipmentIDs adds the equipment edge to Equipment by ids.
func (lc *LocationCreate) AddEquipmentIDs(ids ...string) *LocationCreate {
	if lc.equipment == nil {
		lc.equipment = make(map[string]struct{})
	}
	for i := range ids {
		lc.equipment[ids[i]] = struct{}{}
	}
	return lc
}

// AddEquipment adds the equipment edges to Equipment.
func (lc *LocationCreate) AddEquipment(e ...*Equipment) *LocationCreate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return lc.AddEquipmentIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (lc *LocationCreate) AddPropertyIDs(ids ...string) *LocationCreate {
	if lc.properties == nil {
		lc.properties = make(map[string]struct{})
	}
	for i := range ids {
		lc.properties[ids[i]] = struct{}{}
	}
	return lc
}

// AddProperties adds the properties edges to Property.
func (lc *LocationCreate) AddProperties(p ...*Property) *LocationCreate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return lc.AddPropertyIDs(ids...)
}

// AddSurveyIDs adds the survey edge to Survey by ids.
func (lc *LocationCreate) AddSurveyIDs(ids ...string) *LocationCreate {
	if lc.survey == nil {
		lc.survey = make(map[string]struct{})
	}
	for i := range ids {
		lc.survey[ids[i]] = struct{}{}
	}
	return lc
}

// AddSurvey adds the survey edges to Survey.
func (lc *LocationCreate) AddSurvey(s ...*Survey) *LocationCreate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lc.AddSurveyIDs(ids...)
}

// AddWifiScanIDs adds the wifi_scan edge to SurveyWiFiScan by ids.
func (lc *LocationCreate) AddWifiScanIDs(ids ...string) *LocationCreate {
	if lc.wifi_scan == nil {
		lc.wifi_scan = make(map[string]struct{})
	}
	for i := range ids {
		lc.wifi_scan[ids[i]] = struct{}{}
	}
	return lc
}

// AddWifiScan adds the wifi_scan edges to SurveyWiFiScan.
func (lc *LocationCreate) AddWifiScan(s ...*SurveyWiFiScan) *LocationCreate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lc.AddWifiScanIDs(ids...)
}

// AddCellScanIDs adds the cell_scan edge to SurveyCellScan by ids.
func (lc *LocationCreate) AddCellScanIDs(ids ...string) *LocationCreate {
	if lc.cell_scan == nil {
		lc.cell_scan = make(map[string]struct{})
	}
	for i := range ids {
		lc.cell_scan[ids[i]] = struct{}{}
	}
	return lc
}

// AddCellScan adds the cell_scan edges to SurveyCellScan.
func (lc *LocationCreate) AddCellScan(s ...*SurveyCellScan) *LocationCreate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lc.AddCellScanIDs(ids...)
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (lc *LocationCreate) AddWorkOrderIDs(ids ...string) *LocationCreate {
	if lc.work_orders == nil {
		lc.work_orders = make(map[string]struct{})
	}
	for i := range ids {
		lc.work_orders[ids[i]] = struct{}{}
	}
	return lc
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (lc *LocationCreate) AddWorkOrders(w ...*WorkOrder) *LocationCreate {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return lc.AddWorkOrderIDs(ids...)
}

// AddFloorPlanIDs adds the floor_plans edge to FloorPlan by ids.
func (lc *LocationCreate) AddFloorPlanIDs(ids ...string) *LocationCreate {
	if lc.floor_plans == nil {
		lc.floor_plans = make(map[string]struct{})
	}
	for i := range ids {
		lc.floor_plans[ids[i]] = struct{}{}
	}
	return lc
}

// AddFloorPlans adds the floor_plans edges to FloorPlan.
func (lc *LocationCreate) AddFloorPlans(f ...*FloorPlan) *LocationCreate {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return lc.AddFloorPlanIDs(ids...)
}

// Save creates the Location in the database.
func (lc *LocationCreate) Save(ctx context.Context) (*Location, error) {
	if lc.create_time == nil {
		v := location.DefaultCreateTime()
		lc.create_time = &v
	}
	if lc.update_time == nil {
		v := location.DefaultUpdateTime()
		lc.update_time = &v
	}
	if lc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if err := location.NameValidator(*lc.name); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
	}
	if lc.latitude == nil {
		v := location.DefaultLatitude
		lc.latitude = &v
	}
	if err := location.LatitudeValidator(*lc.latitude); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"latitude\": %v", err)
	}
	if lc.longitude == nil {
		v := location.DefaultLongitude
		lc.longitude = &v
	}
	if err := location.LongitudeValidator(*lc.longitude); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"longitude\": %v", err)
	}
	if lc.site_survey_needed == nil {
		v := location.DefaultSiteSurveyNeeded
		lc.site_survey_needed = &v
	}
	if len(lc._type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if lc._type == nil {
		return nil, errors.New("ent: missing required edge \"type\"")
	}
	if len(lc.parent) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"parent\"")
	}
	return lc.sqlSave(ctx)
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
		l    = &Location{config: lc.config}
		spec = &sqlgraph.CreateSpec{
			Table: location.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: location.FieldID,
			},
		}
	)
	if value := lc.create_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: location.FieldCreateTime,
		})
		l.CreateTime = *value
	}
	if value := lc.update_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: location.FieldUpdateTime,
		})
		l.UpdateTime = *value
	}
	if value := lc.name; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: location.FieldName,
		})
		l.Name = *value
	}
	if value := lc.external_id; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: location.FieldExternalID,
		})
		l.ExternalID = *value
	}
	if value := lc.latitude; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: location.FieldLatitude,
		})
		l.Latitude = *value
	}
	if value := lc.longitude; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  *value,
			Column: location.FieldLongitude,
		})
		l.Longitude = *value
	}
	if value := lc.site_survey_needed; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  *value,
			Column: location.FieldSiteSurveyNeeded,
		})
		l.SiteSurveyNeeded = *value
	}
	if nodes := lc._type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   location.TypeTable,
			Columns: []string{location.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: locationtype.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := lc.parent; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   location.ParentTable,
			Columns: []string{location.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := lc.children; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.ChildrenTable,
			Columns: []string{location.ChildrenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := lc.files; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.FilesTable,
			Columns: []string{location.FilesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: file.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := lc.hyperlinks; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.HyperlinksTable,
			Columns: []string{location.HyperlinksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: hyperlink.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := lc.equipment; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.EquipmentTable,
			Columns: []string{location.EquipmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipment.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := lc.properties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   location.PropertiesTable,
			Columns: []string{location.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: property.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := lc.survey; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.SurveyTable,
			Columns: []string{location.SurveyColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: survey.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := lc.wifi_scan; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.WifiScanTable,
			Columns: []string{location.WifiScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveywifiscan.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := lc.cell_scan; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.CellScanTable,
			Columns: []string{location.CellScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveycellscan.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := lc.work_orders; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.WorkOrdersTable,
			Columns: []string{location.WorkOrdersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workorder.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := lc.floor_plans; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   location.FloorPlansTable,
			Columns: []string{location.FloorPlansColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: floorplan.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges = append(spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, lc.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := spec.ID.Value.(int64)
	l.ID = strconv.FormatInt(id, 10)
	return l, nil
}
