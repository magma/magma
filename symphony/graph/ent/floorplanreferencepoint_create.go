// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/floorplanreferencepoint"
)

// FloorPlanReferencePointCreate is the builder for creating a FloorPlanReferencePoint entity.
type FloorPlanReferencePointCreate struct {
	config
	create_time *time.Time
	update_time *time.Time
	x           *int
	y           *int
	latitude    *float64
	longitude   *float64
}

// SetCreateTime sets the create_time field.
func (fprpc *FloorPlanReferencePointCreate) SetCreateTime(t time.Time) *FloorPlanReferencePointCreate {
	fprpc.create_time = &t
	return fprpc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (fprpc *FloorPlanReferencePointCreate) SetNillableCreateTime(t *time.Time) *FloorPlanReferencePointCreate {
	if t != nil {
		fprpc.SetCreateTime(*t)
	}
	return fprpc
}

// SetUpdateTime sets the update_time field.
func (fprpc *FloorPlanReferencePointCreate) SetUpdateTime(t time.Time) *FloorPlanReferencePointCreate {
	fprpc.update_time = &t
	return fprpc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (fprpc *FloorPlanReferencePointCreate) SetNillableUpdateTime(t *time.Time) *FloorPlanReferencePointCreate {
	if t != nil {
		fprpc.SetUpdateTime(*t)
	}
	return fprpc
}

// SetX sets the x field.
func (fprpc *FloorPlanReferencePointCreate) SetX(i int) *FloorPlanReferencePointCreate {
	fprpc.x = &i
	return fprpc
}

// SetY sets the y field.
func (fprpc *FloorPlanReferencePointCreate) SetY(i int) *FloorPlanReferencePointCreate {
	fprpc.y = &i
	return fprpc
}

// SetLatitude sets the latitude field.
func (fprpc *FloorPlanReferencePointCreate) SetLatitude(f float64) *FloorPlanReferencePointCreate {
	fprpc.latitude = &f
	return fprpc
}

// SetLongitude sets the longitude field.
func (fprpc *FloorPlanReferencePointCreate) SetLongitude(f float64) *FloorPlanReferencePointCreate {
	fprpc.longitude = &f
	return fprpc
}

// Save creates the FloorPlanReferencePoint in the database.
func (fprpc *FloorPlanReferencePointCreate) Save(ctx context.Context) (*FloorPlanReferencePoint, error) {
	if fprpc.create_time == nil {
		v := floorplanreferencepoint.DefaultCreateTime()
		fprpc.create_time = &v
	}
	if fprpc.update_time == nil {
		v := floorplanreferencepoint.DefaultUpdateTime()
		fprpc.update_time = &v
	}
	if fprpc.x == nil {
		return nil, errors.New("ent: missing required field \"x\"")
	}
	if fprpc.y == nil {
		return nil, errors.New("ent: missing required field \"y\"")
	}
	if fprpc.latitude == nil {
		return nil, errors.New("ent: missing required field \"latitude\"")
	}
	if fprpc.longitude == nil {
		return nil, errors.New("ent: missing required field \"longitude\"")
	}
	return fprpc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (fprpc *FloorPlanReferencePointCreate) SaveX(ctx context.Context) *FloorPlanReferencePoint {
	v, err := fprpc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (fprpc *FloorPlanReferencePointCreate) sqlSave(ctx context.Context) (*FloorPlanReferencePoint, error) {
	var (
		builder = sql.Dialect(fprpc.driver.Dialect())
		fprp    = &FloorPlanReferencePoint{config: fprpc.config}
	)
	tx, err := fprpc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(floorplanreferencepoint.Table).Default()
	if value := fprpc.create_time; value != nil {
		insert.Set(floorplanreferencepoint.FieldCreateTime, *value)
		fprp.CreateTime = *value
	}
	if value := fprpc.update_time; value != nil {
		insert.Set(floorplanreferencepoint.FieldUpdateTime, *value)
		fprp.UpdateTime = *value
	}
	if value := fprpc.x; value != nil {
		insert.Set(floorplanreferencepoint.FieldX, *value)
		fprp.X = *value
	}
	if value := fprpc.y; value != nil {
		insert.Set(floorplanreferencepoint.FieldY, *value)
		fprp.Y = *value
	}
	if value := fprpc.latitude; value != nil {
		insert.Set(floorplanreferencepoint.FieldLatitude, *value)
		fprp.Latitude = *value
	}
	if value := fprpc.longitude; value != nil {
		insert.Set(floorplanreferencepoint.FieldLongitude, *value)
		fprp.Longitude = *value
	}
	id, err := insertLastID(ctx, tx, insert.Returning(floorplanreferencepoint.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	fprp.ID = strconv.FormatInt(id, 10)
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return fprp, nil
}
