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

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/floorplan"
	"github.com/facebookincubator/symphony/graph/ent/floorplanreferencepoint"
	"github.com/facebookincubator/symphony/graph/ent/floorplanscale"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// FloorPlanUpdate is the builder for updating FloorPlan entities.
type FloorPlanUpdate struct {
	config

	update_time           *time.Time
	name                  *string
	location              map[string]struct{}
	reference_point       map[string]struct{}
	scale                 map[string]struct{}
	image                 map[string]struct{}
	clearedLocation       bool
	clearedReferencePoint bool
	clearedScale          bool
	clearedImage          bool
	predicates            []predicate.FloorPlan
}

// Where adds a new predicate for the builder.
func (fpu *FloorPlanUpdate) Where(ps ...predicate.FloorPlan) *FloorPlanUpdate {
	fpu.predicates = append(fpu.predicates, ps...)
	return fpu
}

// SetName sets the name field.
func (fpu *FloorPlanUpdate) SetName(s string) *FloorPlanUpdate {
	fpu.name = &s
	return fpu
}

// SetLocationID sets the location edge to Location by id.
func (fpu *FloorPlanUpdate) SetLocationID(id string) *FloorPlanUpdate {
	if fpu.location == nil {
		fpu.location = make(map[string]struct{})
	}
	fpu.location[id] = struct{}{}
	return fpu
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (fpu *FloorPlanUpdate) SetNillableLocationID(id *string) *FloorPlanUpdate {
	if id != nil {
		fpu = fpu.SetLocationID(*id)
	}
	return fpu
}

// SetLocation sets the location edge to Location.
func (fpu *FloorPlanUpdate) SetLocation(l *Location) *FloorPlanUpdate {
	return fpu.SetLocationID(l.ID)
}

// SetReferencePointID sets the reference_point edge to FloorPlanReferencePoint by id.
func (fpu *FloorPlanUpdate) SetReferencePointID(id string) *FloorPlanUpdate {
	if fpu.reference_point == nil {
		fpu.reference_point = make(map[string]struct{})
	}
	fpu.reference_point[id] = struct{}{}
	return fpu
}

// SetNillableReferencePointID sets the reference_point edge to FloorPlanReferencePoint by id if the given value is not nil.
func (fpu *FloorPlanUpdate) SetNillableReferencePointID(id *string) *FloorPlanUpdate {
	if id != nil {
		fpu = fpu.SetReferencePointID(*id)
	}
	return fpu
}

// SetReferencePoint sets the reference_point edge to FloorPlanReferencePoint.
func (fpu *FloorPlanUpdate) SetReferencePoint(f *FloorPlanReferencePoint) *FloorPlanUpdate {
	return fpu.SetReferencePointID(f.ID)
}

// SetScaleID sets the scale edge to FloorPlanScale by id.
func (fpu *FloorPlanUpdate) SetScaleID(id string) *FloorPlanUpdate {
	if fpu.scale == nil {
		fpu.scale = make(map[string]struct{})
	}
	fpu.scale[id] = struct{}{}
	return fpu
}

// SetNillableScaleID sets the scale edge to FloorPlanScale by id if the given value is not nil.
func (fpu *FloorPlanUpdate) SetNillableScaleID(id *string) *FloorPlanUpdate {
	if id != nil {
		fpu = fpu.SetScaleID(*id)
	}
	return fpu
}

// SetScale sets the scale edge to FloorPlanScale.
func (fpu *FloorPlanUpdate) SetScale(f *FloorPlanScale) *FloorPlanUpdate {
	return fpu.SetScaleID(f.ID)
}

// SetImageID sets the image edge to File by id.
func (fpu *FloorPlanUpdate) SetImageID(id string) *FloorPlanUpdate {
	if fpu.image == nil {
		fpu.image = make(map[string]struct{})
	}
	fpu.image[id] = struct{}{}
	return fpu
}

// SetNillableImageID sets the image edge to File by id if the given value is not nil.
func (fpu *FloorPlanUpdate) SetNillableImageID(id *string) *FloorPlanUpdate {
	if id != nil {
		fpu = fpu.SetImageID(*id)
	}
	return fpu
}

// SetImage sets the image edge to File.
func (fpu *FloorPlanUpdate) SetImage(f *File) *FloorPlanUpdate {
	return fpu.SetImageID(f.ID)
}

// ClearLocation clears the location edge to Location.
func (fpu *FloorPlanUpdate) ClearLocation() *FloorPlanUpdate {
	fpu.clearedLocation = true
	return fpu
}

// ClearReferencePoint clears the reference_point edge to FloorPlanReferencePoint.
func (fpu *FloorPlanUpdate) ClearReferencePoint() *FloorPlanUpdate {
	fpu.clearedReferencePoint = true
	return fpu
}

// ClearScale clears the scale edge to FloorPlanScale.
func (fpu *FloorPlanUpdate) ClearScale() *FloorPlanUpdate {
	fpu.clearedScale = true
	return fpu
}

// ClearImage clears the image edge to File.
func (fpu *FloorPlanUpdate) ClearImage() *FloorPlanUpdate {
	fpu.clearedImage = true
	return fpu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (fpu *FloorPlanUpdate) Save(ctx context.Context) (int, error) {
	if fpu.update_time == nil {
		v := floorplan.UpdateDefaultUpdateTime()
		fpu.update_time = &v
	}
	if len(fpu.location) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	if len(fpu.reference_point) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"reference_point\"")
	}
	if len(fpu.scale) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"scale\"")
	}
	if len(fpu.image) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"image\"")
	}
	return fpu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (fpu *FloorPlanUpdate) SaveX(ctx context.Context) int {
	affected, err := fpu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (fpu *FloorPlanUpdate) Exec(ctx context.Context) error {
	_, err := fpu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fpu *FloorPlanUpdate) ExecX(ctx context.Context) {
	if err := fpu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fpu *FloorPlanUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(fpu.driver.Dialect())
		selector = builder.Select(floorplan.FieldID).From(builder.Table(floorplan.Table))
	)
	for _, p := range fpu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = fpu.driver.Query(ctx, query, args, rows); err != nil {
		return 0, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("ent: failed reading id: %v", err)
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return 0, nil
	}

	tx, err := fpu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(floorplan.Table)
	)
	updater = updater.Where(sql.InInts(floorplan.FieldID, ids...))
	if value := fpu.update_time; value != nil {
		updater.Set(floorplan.FieldUpdateTime, *value)
	}
	if value := fpu.name; value != nil {
		updater.Set(floorplan.FieldName, *value)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if fpu.clearedLocation {
		query, args := builder.Update(floorplan.LocationTable).
			SetNull(floorplan.LocationColumn).
			Where(sql.InInts(location.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(fpu.location) > 0 {
		for eid := range fpu.location {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(floorplan.LocationTable).
				Set(floorplan.LocationColumn, eid).
				Where(sql.InInts(floorplan.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if fpu.clearedReferencePoint {
		query, args := builder.Update(floorplan.ReferencePointTable).
			SetNull(floorplan.ReferencePointColumn).
			Where(sql.InInts(floorplanreferencepoint.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(fpu.reference_point) > 0 {
		for eid := range fpu.reference_point {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(floorplan.ReferencePointTable).
				Set(floorplan.ReferencePointColumn, eid).
				Where(sql.InInts(floorplan.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if fpu.clearedScale {
		query, args := builder.Update(floorplan.ScaleTable).
			SetNull(floorplan.ScaleColumn).
			Where(sql.InInts(floorplanscale.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(fpu.scale) > 0 {
		for eid := range fpu.scale {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(floorplan.ScaleTable).
				Set(floorplan.ScaleColumn, eid).
				Where(sql.InInts(floorplan.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if fpu.clearedImage {
		query, args := builder.Update(floorplan.ImageTable).
			SetNull(floorplan.ImageColumn).
			Where(sql.InInts(file.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(fpu.image) > 0 {
		for eid := range fpu.image {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(floorplan.ImageTable).
				Set(floorplan.ImageColumn, eid).
				Where(sql.InInts(floorplan.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// FloorPlanUpdateOne is the builder for updating a single FloorPlan entity.
type FloorPlanUpdateOne struct {
	config
	id string

	update_time           *time.Time
	name                  *string
	location              map[string]struct{}
	reference_point       map[string]struct{}
	scale                 map[string]struct{}
	image                 map[string]struct{}
	clearedLocation       bool
	clearedReferencePoint bool
	clearedScale          bool
	clearedImage          bool
}

// SetName sets the name field.
func (fpuo *FloorPlanUpdateOne) SetName(s string) *FloorPlanUpdateOne {
	fpuo.name = &s
	return fpuo
}

// SetLocationID sets the location edge to Location by id.
func (fpuo *FloorPlanUpdateOne) SetLocationID(id string) *FloorPlanUpdateOne {
	if fpuo.location == nil {
		fpuo.location = make(map[string]struct{})
	}
	fpuo.location[id] = struct{}{}
	return fpuo
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (fpuo *FloorPlanUpdateOne) SetNillableLocationID(id *string) *FloorPlanUpdateOne {
	if id != nil {
		fpuo = fpuo.SetLocationID(*id)
	}
	return fpuo
}

// SetLocation sets the location edge to Location.
func (fpuo *FloorPlanUpdateOne) SetLocation(l *Location) *FloorPlanUpdateOne {
	return fpuo.SetLocationID(l.ID)
}

// SetReferencePointID sets the reference_point edge to FloorPlanReferencePoint by id.
func (fpuo *FloorPlanUpdateOne) SetReferencePointID(id string) *FloorPlanUpdateOne {
	if fpuo.reference_point == nil {
		fpuo.reference_point = make(map[string]struct{})
	}
	fpuo.reference_point[id] = struct{}{}
	return fpuo
}

// SetNillableReferencePointID sets the reference_point edge to FloorPlanReferencePoint by id if the given value is not nil.
func (fpuo *FloorPlanUpdateOne) SetNillableReferencePointID(id *string) *FloorPlanUpdateOne {
	if id != nil {
		fpuo = fpuo.SetReferencePointID(*id)
	}
	return fpuo
}

// SetReferencePoint sets the reference_point edge to FloorPlanReferencePoint.
func (fpuo *FloorPlanUpdateOne) SetReferencePoint(f *FloorPlanReferencePoint) *FloorPlanUpdateOne {
	return fpuo.SetReferencePointID(f.ID)
}

// SetScaleID sets the scale edge to FloorPlanScale by id.
func (fpuo *FloorPlanUpdateOne) SetScaleID(id string) *FloorPlanUpdateOne {
	if fpuo.scale == nil {
		fpuo.scale = make(map[string]struct{})
	}
	fpuo.scale[id] = struct{}{}
	return fpuo
}

// SetNillableScaleID sets the scale edge to FloorPlanScale by id if the given value is not nil.
func (fpuo *FloorPlanUpdateOne) SetNillableScaleID(id *string) *FloorPlanUpdateOne {
	if id != nil {
		fpuo = fpuo.SetScaleID(*id)
	}
	return fpuo
}

// SetScale sets the scale edge to FloorPlanScale.
func (fpuo *FloorPlanUpdateOne) SetScale(f *FloorPlanScale) *FloorPlanUpdateOne {
	return fpuo.SetScaleID(f.ID)
}

// SetImageID sets the image edge to File by id.
func (fpuo *FloorPlanUpdateOne) SetImageID(id string) *FloorPlanUpdateOne {
	if fpuo.image == nil {
		fpuo.image = make(map[string]struct{})
	}
	fpuo.image[id] = struct{}{}
	return fpuo
}

// SetNillableImageID sets the image edge to File by id if the given value is not nil.
func (fpuo *FloorPlanUpdateOne) SetNillableImageID(id *string) *FloorPlanUpdateOne {
	if id != nil {
		fpuo = fpuo.SetImageID(*id)
	}
	return fpuo
}

// SetImage sets the image edge to File.
func (fpuo *FloorPlanUpdateOne) SetImage(f *File) *FloorPlanUpdateOne {
	return fpuo.SetImageID(f.ID)
}

// ClearLocation clears the location edge to Location.
func (fpuo *FloorPlanUpdateOne) ClearLocation() *FloorPlanUpdateOne {
	fpuo.clearedLocation = true
	return fpuo
}

// ClearReferencePoint clears the reference_point edge to FloorPlanReferencePoint.
func (fpuo *FloorPlanUpdateOne) ClearReferencePoint() *FloorPlanUpdateOne {
	fpuo.clearedReferencePoint = true
	return fpuo
}

// ClearScale clears the scale edge to FloorPlanScale.
func (fpuo *FloorPlanUpdateOne) ClearScale() *FloorPlanUpdateOne {
	fpuo.clearedScale = true
	return fpuo
}

// ClearImage clears the image edge to File.
func (fpuo *FloorPlanUpdateOne) ClearImage() *FloorPlanUpdateOne {
	fpuo.clearedImage = true
	return fpuo
}

// Save executes the query and returns the updated entity.
func (fpuo *FloorPlanUpdateOne) Save(ctx context.Context) (*FloorPlan, error) {
	if fpuo.update_time == nil {
		v := floorplan.UpdateDefaultUpdateTime()
		fpuo.update_time = &v
	}
	if len(fpuo.location) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	if len(fpuo.reference_point) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"reference_point\"")
	}
	if len(fpuo.scale) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"scale\"")
	}
	if len(fpuo.image) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"image\"")
	}
	return fpuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (fpuo *FloorPlanUpdateOne) SaveX(ctx context.Context) *FloorPlan {
	fp, err := fpuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return fp
}

// Exec executes the query on the entity.
func (fpuo *FloorPlanUpdateOne) Exec(ctx context.Context) error {
	_, err := fpuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fpuo *FloorPlanUpdateOne) ExecX(ctx context.Context) {
	if err := fpuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fpuo *FloorPlanUpdateOne) sqlSave(ctx context.Context) (fp *FloorPlan, err error) {
	var (
		builder  = sql.Dialect(fpuo.driver.Dialect())
		selector = builder.Select(floorplan.Columns...).From(builder.Table(floorplan.Table))
	)
	floorplan.ID(fpuo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = fpuo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		fp = &FloorPlan{config: fpuo.config}
		if err := fp.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into FloorPlan: %v", err)
		}
		id = fp.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("FloorPlan with id: %v", fpuo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one FloorPlan with the same id: %v", fpuo.id)
	}

	tx, err := fpuo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(floorplan.Table)
	)
	updater = updater.Where(sql.InInts(floorplan.FieldID, ids...))
	if value := fpuo.update_time; value != nil {
		updater.Set(floorplan.FieldUpdateTime, *value)
		fp.UpdateTime = *value
	}
	if value := fpuo.name; value != nil {
		updater.Set(floorplan.FieldName, *value)
		fp.Name = *value
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if fpuo.clearedLocation {
		query, args := builder.Update(floorplan.LocationTable).
			SetNull(floorplan.LocationColumn).
			Where(sql.InInts(location.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(fpuo.location) > 0 {
		for eid := range fpuo.location {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(floorplan.LocationTable).
				Set(floorplan.LocationColumn, eid).
				Where(sql.InInts(floorplan.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if fpuo.clearedReferencePoint {
		query, args := builder.Update(floorplan.ReferencePointTable).
			SetNull(floorplan.ReferencePointColumn).
			Where(sql.InInts(floorplanreferencepoint.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(fpuo.reference_point) > 0 {
		for eid := range fpuo.reference_point {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(floorplan.ReferencePointTable).
				Set(floorplan.ReferencePointColumn, eid).
				Where(sql.InInts(floorplan.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if fpuo.clearedScale {
		query, args := builder.Update(floorplan.ScaleTable).
			SetNull(floorplan.ScaleColumn).
			Where(sql.InInts(floorplanscale.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(fpuo.scale) > 0 {
		for eid := range fpuo.scale {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(floorplan.ScaleTable).
				Set(floorplan.ScaleColumn, eid).
				Where(sql.InInts(floorplan.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if fpuo.clearedImage {
		query, args := builder.Update(floorplan.ImageTable).
			SetNull(floorplan.ImageColumn).
			Where(sql.InInts(file.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(fpuo.image) > 0 {
		for eid := range fpuo.image {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(floorplan.ImageTable).
				Set(floorplan.ImageColumn, eid).
				Where(sql.InInts(floorplan.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return fp, nil
}
