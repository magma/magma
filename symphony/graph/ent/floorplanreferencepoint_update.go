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
	"github.com/facebookincubator/symphony/graph/ent/floorplanreferencepoint"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// FloorPlanReferencePointUpdate is the builder for updating FloorPlanReferencePoint entities.
type FloorPlanReferencePointUpdate struct {
	config

	update_time  *time.Time
	x            *int
	addx         *int
	y            *int
	addy         *int
	latitude     *float64
	addlatitude  *float64
	longitude    *float64
	addlongitude *float64
	predicates   []predicate.FloorPlanReferencePoint
}

// Where adds a new predicate for the builder.
func (fprpu *FloorPlanReferencePointUpdate) Where(ps ...predicate.FloorPlanReferencePoint) *FloorPlanReferencePointUpdate {
	fprpu.predicates = append(fprpu.predicates, ps...)
	return fprpu
}

// SetX sets the x field.
func (fprpu *FloorPlanReferencePointUpdate) SetX(i int) *FloorPlanReferencePointUpdate {
	fprpu.x = &i
	fprpu.addx = nil
	return fprpu
}

// AddX adds i to x.
func (fprpu *FloorPlanReferencePointUpdate) AddX(i int) *FloorPlanReferencePointUpdate {
	if fprpu.addx == nil {
		fprpu.addx = &i
	} else {
		*fprpu.addx += i
	}
	return fprpu
}

// SetY sets the y field.
func (fprpu *FloorPlanReferencePointUpdate) SetY(i int) *FloorPlanReferencePointUpdate {
	fprpu.y = &i
	fprpu.addy = nil
	return fprpu
}

// AddY adds i to y.
func (fprpu *FloorPlanReferencePointUpdate) AddY(i int) *FloorPlanReferencePointUpdate {
	if fprpu.addy == nil {
		fprpu.addy = &i
	} else {
		*fprpu.addy += i
	}
	return fprpu
}

// SetLatitude sets the latitude field.
func (fprpu *FloorPlanReferencePointUpdate) SetLatitude(f float64) *FloorPlanReferencePointUpdate {
	fprpu.latitude = &f
	fprpu.addlatitude = nil
	return fprpu
}

// AddLatitude adds f to latitude.
func (fprpu *FloorPlanReferencePointUpdate) AddLatitude(f float64) *FloorPlanReferencePointUpdate {
	if fprpu.addlatitude == nil {
		fprpu.addlatitude = &f
	} else {
		*fprpu.addlatitude += f
	}
	return fprpu
}

// SetLongitude sets the longitude field.
func (fprpu *FloorPlanReferencePointUpdate) SetLongitude(f float64) *FloorPlanReferencePointUpdate {
	fprpu.longitude = &f
	fprpu.addlongitude = nil
	return fprpu
}

// AddLongitude adds f to longitude.
func (fprpu *FloorPlanReferencePointUpdate) AddLongitude(f float64) *FloorPlanReferencePointUpdate {
	if fprpu.addlongitude == nil {
		fprpu.addlongitude = &f
	} else {
		*fprpu.addlongitude += f
	}
	return fprpu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (fprpu *FloorPlanReferencePointUpdate) Save(ctx context.Context) (int, error) {
	if fprpu.update_time == nil {
		v := floorplanreferencepoint.UpdateDefaultUpdateTime()
		fprpu.update_time = &v
	}
	return fprpu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (fprpu *FloorPlanReferencePointUpdate) SaveX(ctx context.Context) int {
	affected, err := fprpu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (fprpu *FloorPlanReferencePointUpdate) Exec(ctx context.Context) error {
	_, err := fprpu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fprpu *FloorPlanReferencePointUpdate) ExecX(ctx context.Context) {
	if err := fprpu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fprpu *FloorPlanReferencePointUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(fprpu.driver.Dialect())
		selector = builder.Select(floorplanreferencepoint.FieldID).From(builder.Table(floorplanreferencepoint.Table))
	)
	for _, p := range fprpu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = fprpu.driver.Query(ctx, query, args, rows); err != nil {
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

	tx, err := fprpu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(floorplanreferencepoint.Table)
	)
	updater = updater.Where(sql.InInts(floorplanreferencepoint.FieldID, ids...))
	if value := fprpu.update_time; value != nil {
		updater.Set(floorplanreferencepoint.FieldUpdateTime, *value)
	}
	if value := fprpu.x; value != nil {
		updater.Set(floorplanreferencepoint.FieldX, *value)
	}
	if value := fprpu.addx; value != nil {
		updater.Add(floorplanreferencepoint.FieldX, *value)
	}
	if value := fprpu.y; value != nil {
		updater.Set(floorplanreferencepoint.FieldY, *value)
	}
	if value := fprpu.addy; value != nil {
		updater.Add(floorplanreferencepoint.FieldY, *value)
	}
	if value := fprpu.latitude; value != nil {
		updater.Set(floorplanreferencepoint.FieldLatitude, *value)
	}
	if value := fprpu.addlatitude; value != nil {
		updater.Add(floorplanreferencepoint.FieldLatitude, *value)
	}
	if value := fprpu.longitude; value != nil {
		updater.Set(floorplanreferencepoint.FieldLongitude, *value)
	}
	if value := fprpu.addlongitude; value != nil {
		updater.Add(floorplanreferencepoint.FieldLongitude, *value)
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

// FloorPlanReferencePointUpdateOne is the builder for updating a single FloorPlanReferencePoint entity.
type FloorPlanReferencePointUpdateOne struct {
	config
	id string

	update_time  *time.Time
	x            *int
	addx         *int
	y            *int
	addy         *int
	latitude     *float64
	addlatitude  *float64
	longitude    *float64
	addlongitude *float64
}

// SetX sets the x field.
func (fprpuo *FloorPlanReferencePointUpdateOne) SetX(i int) *FloorPlanReferencePointUpdateOne {
	fprpuo.x = &i
	fprpuo.addx = nil
	return fprpuo
}

// AddX adds i to x.
func (fprpuo *FloorPlanReferencePointUpdateOne) AddX(i int) *FloorPlanReferencePointUpdateOne {
	if fprpuo.addx == nil {
		fprpuo.addx = &i
	} else {
		*fprpuo.addx += i
	}
	return fprpuo
}

// SetY sets the y field.
func (fprpuo *FloorPlanReferencePointUpdateOne) SetY(i int) *FloorPlanReferencePointUpdateOne {
	fprpuo.y = &i
	fprpuo.addy = nil
	return fprpuo
}

// AddY adds i to y.
func (fprpuo *FloorPlanReferencePointUpdateOne) AddY(i int) *FloorPlanReferencePointUpdateOne {
	if fprpuo.addy == nil {
		fprpuo.addy = &i
	} else {
		*fprpuo.addy += i
	}
	return fprpuo
}

// SetLatitude sets the latitude field.
func (fprpuo *FloorPlanReferencePointUpdateOne) SetLatitude(f float64) *FloorPlanReferencePointUpdateOne {
	fprpuo.latitude = &f
	fprpuo.addlatitude = nil
	return fprpuo
}

// AddLatitude adds f to latitude.
func (fprpuo *FloorPlanReferencePointUpdateOne) AddLatitude(f float64) *FloorPlanReferencePointUpdateOne {
	if fprpuo.addlatitude == nil {
		fprpuo.addlatitude = &f
	} else {
		*fprpuo.addlatitude += f
	}
	return fprpuo
}

// SetLongitude sets the longitude field.
func (fprpuo *FloorPlanReferencePointUpdateOne) SetLongitude(f float64) *FloorPlanReferencePointUpdateOne {
	fprpuo.longitude = &f
	fprpuo.addlongitude = nil
	return fprpuo
}

// AddLongitude adds f to longitude.
func (fprpuo *FloorPlanReferencePointUpdateOne) AddLongitude(f float64) *FloorPlanReferencePointUpdateOne {
	if fprpuo.addlongitude == nil {
		fprpuo.addlongitude = &f
	} else {
		*fprpuo.addlongitude += f
	}
	return fprpuo
}

// Save executes the query and returns the updated entity.
func (fprpuo *FloorPlanReferencePointUpdateOne) Save(ctx context.Context) (*FloorPlanReferencePoint, error) {
	if fprpuo.update_time == nil {
		v := floorplanreferencepoint.UpdateDefaultUpdateTime()
		fprpuo.update_time = &v
	}
	return fprpuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (fprpuo *FloorPlanReferencePointUpdateOne) SaveX(ctx context.Context) *FloorPlanReferencePoint {
	fprp, err := fprpuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return fprp
}

// Exec executes the query on the entity.
func (fprpuo *FloorPlanReferencePointUpdateOne) Exec(ctx context.Context) error {
	_, err := fprpuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fprpuo *FloorPlanReferencePointUpdateOne) ExecX(ctx context.Context) {
	if err := fprpuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fprpuo *FloorPlanReferencePointUpdateOne) sqlSave(ctx context.Context) (fprp *FloorPlanReferencePoint, err error) {
	var (
		builder  = sql.Dialect(fprpuo.driver.Dialect())
		selector = builder.Select(floorplanreferencepoint.Columns...).From(builder.Table(floorplanreferencepoint.Table))
	)
	floorplanreferencepoint.ID(fprpuo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = fprpuo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		fprp = &FloorPlanReferencePoint{config: fprpuo.config}
		if err := fprp.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into FloorPlanReferencePoint: %v", err)
		}
		id = fprp.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("FloorPlanReferencePoint with id: %v", fprpuo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one FloorPlanReferencePoint with the same id: %v", fprpuo.id)
	}

	tx, err := fprpuo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(floorplanreferencepoint.Table)
	)
	updater = updater.Where(sql.InInts(floorplanreferencepoint.FieldID, ids...))
	if value := fprpuo.update_time; value != nil {
		updater.Set(floorplanreferencepoint.FieldUpdateTime, *value)
		fprp.UpdateTime = *value
	}
	if value := fprpuo.x; value != nil {
		updater.Set(floorplanreferencepoint.FieldX, *value)
		fprp.X = *value
	}
	if value := fprpuo.addx; value != nil {
		updater.Add(floorplanreferencepoint.FieldX, *value)
		fprp.X += *value
	}
	if value := fprpuo.y; value != nil {
		updater.Set(floorplanreferencepoint.FieldY, *value)
		fprp.Y = *value
	}
	if value := fprpuo.addy; value != nil {
		updater.Add(floorplanreferencepoint.FieldY, *value)
		fprp.Y += *value
	}
	if value := fprpuo.latitude; value != nil {
		updater.Set(floorplanreferencepoint.FieldLatitude, *value)
		fprp.Latitude = *value
	}
	if value := fprpuo.addlatitude; value != nil {
		updater.Add(floorplanreferencepoint.FieldLatitude, *value)
		fprp.Latitude += *value
	}
	if value := fprpuo.longitude; value != nil {
		updater.Set(floorplanreferencepoint.FieldLongitude, *value)
		fprp.Longitude = *value
	}
	if value := fprpuo.addlongitude; value != nil {
		updater.Add(floorplanreferencepoint.FieldLongitude, *value)
		fprp.Longitude += *value
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
	return fprp, nil
}
