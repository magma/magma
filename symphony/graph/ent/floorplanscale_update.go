// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/floorplanscale"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// FloorPlanScaleUpdate is the builder for updating FloorPlanScale entities.
type FloorPlanScaleUpdate struct {
	config

	update_time           *time.Time
	reference_point1_x    *int
	addreference_point1_x *int
	reference_point1_y    *int
	addreference_point1_y *int
	reference_point2_x    *int
	addreference_point2_x *int
	reference_point2_y    *int
	addreference_point2_y *int
	scale_in_meters       *float64
	addscale_in_meters    *float64
	predicates            []predicate.FloorPlanScale
}

// Where adds a new predicate for the builder.
func (fpsu *FloorPlanScaleUpdate) Where(ps ...predicate.FloorPlanScale) *FloorPlanScaleUpdate {
	fpsu.predicates = append(fpsu.predicates, ps...)
	return fpsu
}

// SetReferencePoint1X sets the reference_point1_x field.
func (fpsu *FloorPlanScaleUpdate) SetReferencePoint1X(i int) *FloorPlanScaleUpdate {
	fpsu.reference_point1_x = &i
	fpsu.addreference_point1_x = nil
	return fpsu
}

// AddReferencePoint1X adds i to reference_point1_x.
func (fpsu *FloorPlanScaleUpdate) AddReferencePoint1X(i int) *FloorPlanScaleUpdate {
	if fpsu.addreference_point1_x == nil {
		fpsu.addreference_point1_x = &i
	} else {
		*fpsu.addreference_point1_x += i
	}
	return fpsu
}

// SetReferencePoint1Y sets the reference_point1_y field.
func (fpsu *FloorPlanScaleUpdate) SetReferencePoint1Y(i int) *FloorPlanScaleUpdate {
	fpsu.reference_point1_y = &i
	fpsu.addreference_point1_y = nil
	return fpsu
}

// AddReferencePoint1Y adds i to reference_point1_y.
func (fpsu *FloorPlanScaleUpdate) AddReferencePoint1Y(i int) *FloorPlanScaleUpdate {
	if fpsu.addreference_point1_y == nil {
		fpsu.addreference_point1_y = &i
	} else {
		*fpsu.addreference_point1_y += i
	}
	return fpsu
}

// SetReferencePoint2X sets the reference_point2_x field.
func (fpsu *FloorPlanScaleUpdate) SetReferencePoint2X(i int) *FloorPlanScaleUpdate {
	fpsu.reference_point2_x = &i
	fpsu.addreference_point2_x = nil
	return fpsu
}

// AddReferencePoint2X adds i to reference_point2_x.
func (fpsu *FloorPlanScaleUpdate) AddReferencePoint2X(i int) *FloorPlanScaleUpdate {
	if fpsu.addreference_point2_x == nil {
		fpsu.addreference_point2_x = &i
	} else {
		*fpsu.addreference_point2_x += i
	}
	return fpsu
}

// SetReferencePoint2Y sets the reference_point2_y field.
func (fpsu *FloorPlanScaleUpdate) SetReferencePoint2Y(i int) *FloorPlanScaleUpdate {
	fpsu.reference_point2_y = &i
	fpsu.addreference_point2_y = nil
	return fpsu
}

// AddReferencePoint2Y adds i to reference_point2_y.
func (fpsu *FloorPlanScaleUpdate) AddReferencePoint2Y(i int) *FloorPlanScaleUpdate {
	if fpsu.addreference_point2_y == nil {
		fpsu.addreference_point2_y = &i
	} else {
		*fpsu.addreference_point2_y += i
	}
	return fpsu
}

// SetScaleInMeters sets the scale_in_meters field.
func (fpsu *FloorPlanScaleUpdate) SetScaleInMeters(f float64) *FloorPlanScaleUpdate {
	fpsu.scale_in_meters = &f
	fpsu.addscale_in_meters = nil
	return fpsu
}

// AddScaleInMeters adds f to scale_in_meters.
func (fpsu *FloorPlanScaleUpdate) AddScaleInMeters(f float64) *FloorPlanScaleUpdate {
	if fpsu.addscale_in_meters == nil {
		fpsu.addscale_in_meters = &f
	} else {
		*fpsu.addscale_in_meters += f
	}
	return fpsu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (fpsu *FloorPlanScaleUpdate) Save(ctx context.Context) (int, error) {
	if fpsu.update_time == nil {
		v := floorplanscale.UpdateDefaultUpdateTime()
		fpsu.update_time = &v
	}
	return fpsu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (fpsu *FloorPlanScaleUpdate) SaveX(ctx context.Context) int {
	affected, err := fpsu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (fpsu *FloorPlanScaleUpdate) Exec(ctx context.Context) error {
	_, err := fpsu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fpsu *FloorPlanScaleUpdate) ExecX(ctx context.Context) {
	if err := fpsu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fpsu *FloorPlanScaleUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(fpsu.driver.Dialect())
		selector = builder.Select(floorplanscale.FieldID).From(builder.Table(floorplanscale.Table))
	)
	for _, p := range fpsu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = fpsu.driver.Query(ctx, query, args, rows); err != nil {
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

	tx, err := fpsu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(floorplanscale.Table).Where(sql.InInts(floorplanscale.FieldID, ids...))
	)
	if value := fpsu.update_time; value != nil {
		updater.Set(floorplanscale.FieldUpdateTime, *value)
	}
	if value := fpsu.reference_point1_x; value != nil {
		updater.Set(floorplanscale.FieldReferencePoint1X, *value)
	}
	if value := fpsu.addreference_point1_x; value != nil {
		updater.Add(floorplanscale.FieldReferencePoint1X, *value)
	}
	if value := fpsu.reference_point1_y; value != nil {
		updater.Set(floorplanscale.FieldReferencePoint1Y, *value)
	}
	if value := fpsu.addreference_point1_y; value != nil {
		updater.Add(floorplanscale.FieldReferencePoint1Y, *value)
	}
	if value := fpsu.reference_point2_x; value != nil {
		updater.Set(floorplanscale.FieldReferencePoint2X, *value)
	}
	if value := fpsu.addreference_point2_x; value != nil {
		updater.Add(floorplanscale.FieldReferencePoint2X, *value)
	}
	if value := fpsu.reference_point2_y; value != nil {
		updater.Set(floorplanscale.FieldReferencePoint2Y, *value)
	}
	if value := fpsu.addreference_point2_y; value != nil {
		updater.Add(floorplanscale.FieldReferencePoint2Y, *value)
	}
	if value := fpsu.scale_in_meters; value != nil {
		updater.Set(floorplanscale.FieldScaleInMeters, *value)
	}
	if value := fpsu.addscale_in_meters; value != nil {
		updater.Add(floorplanscale.FieldScaleInMeters, *value)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// FloorPlanScaleUpdateOne is the builder for updating a single FloorPlanScale entity.
type FloorPlanScaleUpdateOne struct {
	config
	id string

	update_time           *time.Time
	reference_point1_x    *int
	addreference_point1_x *int
	reference_point1_y    *int
	addreference_point1_y *int
	reference_point2_x    *int
	addreference_point2_x *int
	reference_point2_y    *int
	addreference_point2_y *int
	scale_in_meters       *float64
	addscale_in_meters    *float64
}

// SetReferencePoint1X sets the reference_point1_x field.
func (fpsuo *FloorPlanScaleUpdateOne) SetReferencePoint1X(i int) *FloorPlanScaleUpdateOne {
	fpsuo.reference_point1_x = &i
	fpsuo.addreference_point1_x = nil
	return fpsuo
}

// AddReferencePoint1X adds i to reference_point1_x.
func (fpsuo *FloorPlanScaleUpdateOne) AddReferencePoint1X(i int) *FloorPlanScaleUpdateOne {
	if fpsuo.addreference_point1_x == nil {
		fpsuo.addreference_point1_x = &i
	} else {
		*fpsuo.addreference_point1_x += i
	}
	return fpsuo
}

// SetReferencePoint1Y sets the reference_point1_y field.
func (fpsuo *FloorPlanScaleUpdateOne) SetReferencePoint1Y(i int) *FloorPlanScaleUpdateOne {
	fpsuo.reference_point1_y = &i
	fpsuo.addreference_point1_y = nil
	return fpsuo
}

// AddReferencePoint1Y adds i to reference_point1_y.
func (fpsuo *FloorPlanScaleUpdateOne) AddReferencePoint1Y(i int) *FloorPlanScaleUpdateOne {
	if fpsuo.addreference_point1_y == nil {
		fpsuo.addreference_point1_y = &i
	} else {
		*fpsuo.addreference_point1_y += i
	}
	return fpsuo
}

// SetReferencePoint2X sets the reference_point2_x field.
func (fpsuo *FloorPlanScaleUpdateOne) SetReferencePoint2X(i int) *FloorPlanScaleUpdateOne {
	fpsuo.reference_point2_x = &i
	fpsuo.addreference_point2_x = nil
	return fpsuo
}

// AddReferencePoint2X adds i to reference_point2_x.
func (fpsuo *FloorPlanScaleUpdateOne) AddReferencePoint2X(i int) *FloorPlanScaleUpdateOne {
	if fpsuo.addreference_point2_x == nil {
		fpsuo.addreference_point2_x = &i
	} else {
		*fpsuo.addreference_point2_x += i
	}
	return fpsuo
}

// SetReferencePoint2Y sets the reference_point2_y field.
func (fpsuo *FloorPlanScaleUpdateOne) SetReferencePoint2Y(i int) *FloorPlanScaleUpdateOne {
	fpsuo.reference_point2_y = &i
	fpsuo.addreference_point2_y = nil
	return fpsuo
}

// AddReferencePoint2Y adds i to reference_point2_y.
func (fpsuo *FloorPlanScaleUpdateOne) AddReferencePoint2Y(i int) *FloorPlanScaleUpdateOne {
	if fpsuo.addreference_point2_y == nil {
		fpsuo.addreference_point2_y = &i
	} else {
		*fpsuo.addreference_point2_y += i
	}
	return fpsuo
}

// SetScaleInMeters sets the scale_in_meters field.
func (fpsuo *FloorPlanScaleUpdateOne) SetScaleInMeters(f float64) *FloorPlanScaleUpdateOne {
	fpsuo.scale_in_meters = &f
	fpsuo.addscale_in_meters = nil
	return fpsuo
}

// AddScaleInMeters adds f to scale_in_meters.
func (fpsuo *FloorPlanScaleUpdateOne) AddScaleInMeters(f float64) *FloorPlanScaleUpdateOne {
	if fpsuo.addscale_in_meters == nil {
		fpsuo.addscale_in_meters = &f
	} else {
		*fpsuo.addscale_in_meters += f
	}
	return fpsuo
}

// Save executes the query and returns the updated entity.
func (fpsuo *FloorPlanScaleUpdateOne) Save(ctx context.Context) (*FloorPlanScale, error) {
	if fpsuo.update_time == nil {
		v := floorplanscale.UpdateDefaultUpdateTime()
		fpsuo.update_time = &v
	}
	return fpsuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (fpsuo *FloorPlanScaleUpdateOne) SaveX(ctx context.Context) *FloorPlanScale {
	fps, err := fpsuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return fps
}

// Exec executes the query on the entity.
func (fpsuo *FloorPlanScaleUpdateOne) Exec(ctx context.Context) error {
	_, err := fpsuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fpsuo *FloorPlanScaleUpdateOne) ExecX(ctx context.Context) {
	if err := fpsuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fpsuo *FloorPlanScaleUpdateOne) sqlSave(ctx context.Context) (fps *FloorPlanScale, err error) {
	var (
		builder  = sql.Dialect(fpsuo.driver.Dialect())
		selector = builder.Select(floorplanscale.Columns...).From(builder.Table(floorplanscale.Table))
	)
	floorplanscale.ID(fpsuo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = fpsuo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		fps = &FloorPlanScale{config: fpsuo.config}
		if err := fps.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into FloorPlanScale: %v", err)
		}
		id = fps.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("FloorPlanScale with id: %v", fpsuo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one FloorPlanScale with the same id: %v", fpsuo.id)
	}

	tx, err := fpsuo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(floorplanscale.Table).Where(sql.InInts(floorplanscale.FieldID, ids...))
	)
	if value := fpsuo.update_time; value != nil {
		updater.Set(floorplanscale.FieldUpdateTime, *value)
		fps.UpdateTime = *value
	}
	if value := fpsuo.reference_point1_x; value != nil {
		updater.Set(floorplanscale.FieldReferencePoint1X, *value)
		fps.ReferencePoint1X = *value
	}
	if value := fpsuo.addreference_point1_x; value != nil {
		updater.Add(floorplanscale.FieldReferencePoint1X, *value)
		fps.ReferencePoint1X += *value
	}
	if value := fpsuo.reference_point1_y; value != nil {
		updater.Set(floorplanscale.FieldReferencePoint1Y, *value)
		fps.ReferencePoint1Y = *value
	}
	if value := fpsuo.addreference_point1_y; value != nil {
		updater.Add(floorplanscale.FieldReferencePoint1Y, *value)
		fps.ReferencePoint1Y += *value
	}
	if value := fpsuo.reference_point2_x; value != nil {
		updater.Set(floorplanscale.FieldReferencePoint2X, *value)
		fps.ReferencePoint2X = *value
	}
	if value := fpsuo.addreference_point2_x; value != nil {
		updater.Add(floorplanscale.FieldReferencePoint2X, *value)
		fps.ReferencePoint2X += *value
	}
	if value := fpsuo.reference_point2_y; value != nil {
		updater.Set(floorplanscale.FieldReferencePoint2Y, *value)
		fps.ReferencePoint2Y = *value
	}
	if value := fpsuo.addreference_point2_y; value != nil {
		updater.Add(floorplanscale.FieldReferencePoint2Y, *value)
		fps.ReferencePoint2Y += *value
	}
	if value := fpsuo.scale_in_meters; value != nil {
		updater.Set(floorplanscale.FieldScaleInMeters, *value)
		fps.ScaleInMeters = *value
	}
	if value := fpsuo.addscale_in_meters; value != nil {
		updater.Add(floorplanscale.FieldScaleInMeters, *value)
		fps.ScaleInMeters += *value
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return fps, nil
}
