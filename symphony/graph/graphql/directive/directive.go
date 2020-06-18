// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package directive

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"go.uber.org/zap"
)

type directive struct {
	logger log.Logger
}

// New creates a graphql directive root.
func New(logger log.Logger) generated.DirectiveRoot {
	d := &directive{logger}
	return generated.DirectiveRoot{
		DeprecatedInput: d.DeprecatedInput,
		Length:          d.Length,
		Range:           d.Range,
		UniqueField:     d.UniqueField,
	}
}

func (d *directive) recoverFromPanic(ctx context.Context, err *error) {
	if r := recover(); r != nil {
		d.logger.For(ctx).
			Error("panic recovery",
				zap.Stack("stacktrace"),
				zap.Reflect("error", r),
			)
		*err = fmt.Errorf("recovered from panic %q: %w", r, *err)
	}
}

func (d *directive) Length(ctx context.Context, _ interface{}, next graphql.Resolver, min int, max *int) (res interface{}, err error) {
	if res, err = next(ctx); err != nil {
		return
	}
	defer d.recoverFromPanic(ctx, &err)
	length := reflect.Indirect(reflect.ValueOf(res)).Len()
	if length < min {
		return nil, gqlerror.Errorf("Field length %d below allowed minimum %d", length, min)
	}
	if max != nil && length > *max {
		return nil, gqlerror.Errorf("Field length %d above allowed maximum %d", length, *max)
	}
	return
}

func (d *directive) Range(ctx context.Context, _ interface{}, next graphql.Resolver, min, max *float64) (res interface{}, err error) {
	if res, err = next(ctx); err != nil {
		return
	}
	defer d.recoverFromPanic(ctx, &err)
	var value float64
	switch rv := reflect.Indirect(reflect.ValueOf(res)); rv.Kind() {
	case reflect.Float32, reflect.Float64:
		value = rv.Float()
	default:
		value = float64(rv.Int())
	}
	if min != nil && value < *min {
		return nil, gqlerror.Errorf("Field value %f below allowed minimum %f", value, *min)
	}
	if max != nil && value > *max {
		return nil, gqlerror.Errorf("Field value %f above allowed maximum %f", value, *max)
	}
	return
}

func (d *directive) UniqueField(ctx context.Context, _ interface{}, next graphql.Resolver, typ, field string) (interface{}, error) {
	objs, err := next(ctx)
	if err != nil {
		return objs, err
	}
	if err := d.uniqueField(objs, typ, field); err != nil {
		d.logger.For(ctx).Warn("unique name violation", zap.Error(err))
		return nil, err
	}
	return objs, nil
}

func (d *directive) uniqueField(objs interface{}, objType, field string) *gqlerror.Error {
	val := reflect.ValueOf(objs)
	if val.Kind() != reflect.Slice {
		return nil
	}

	switch typ := val.Type().Elem(); typ.Kind() {
	case reflect.Ptr:
		typ = typ.Elem()
		if typ.Kind() != reflect.Struct {
			return nil
		}
		fallthrough
	case reflect.Struct:
		field, found := typ.FieldByName(field)
		if !found || field.Type.Kind() != reflect.String {
			return nil
		}
	default:
		return nil
	}

	length := val.Len()
	values := make(map[string]struct{}, length)
	for i := 0; i < length; i++ {
		v := val.Index(i)
		if v = reflect.Indirect(v); !v.IsValid() {
			continue
		}
		value := v.FieldByName(field).Interface().(string)
		if _, ok := values[value]; ok {
			return gqlerror.Errorf("Duplicate %s %s: %s", objType, strings.ToLower(field), value)
		}
		values[value] = struct{}{}
	}
	return nil
}

func (d *directive) isNil(i interface{}) (result bool) {
	if i == nil {
		return true
	}
	defer func() {
		if err := recover(); err != nil {
			result = false
		}
	}()
	return reflect.ValueOf(i).IsNil()
}

func (d *directive) DeprecatedInput(ctx context.Context, obj interface{}, next graphql.Resolver, name, duplicateError string, newField *string) (res interface{}, err error) {
	value, err := next(ctx)
	if err != nil || d.isNil(value) {
		return value, err
	}
	tags := []tag.Mutator{
		tag.Upsert(Field, name),
	}

	_ = stats.RecordWithTags(
		ctx,
		tags,
		ServerDeprecatedInputs.M(1),
	)
	if newField == nil {
		return next(ctx)
	}
	m, ok := obj.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected obj type %T", obj)
	}
	if newValue, ok := m[*newField]; ok && newValue != nil {
		return nil, gqlerror.Errorf(duplicateError)
	}
	return value, nil
}
