// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package directive

import (
	"context"
	"reflect"
	"strings"

	"github.com/facebookincubator/symphony/cloud/log"
	"github.com/facebookincubator/symphony/graph/graphql/generated"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/gqlerror"
	"go.uber.org/zap"
)

type directive struct {
	logger log.Logger
}

// New creates a graphql directive root.
func New(logger log.Logger) generated.DirectiveRoot {
	d := &directive{logger}
	return generated.DirectiveRoot{
		Length:      d.Length,
		UniqueField: d.UniqueField,
	}
}

func (d *directive) Length(ctx context.Context, obj interface{}, next graphql.Resolver, min int, max *int) (res interface{}, err error) {
	if res, err = next(ctx); err != nil {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			d.logger.For(ctx).Error("panic recovery", zap.Reflect("error", err))
		}
	}()
	length := reflect.Indirect(reflect.ValueOf(res)).Len()
	if length < min {
		return nil, gqlerror.Errorf("Field length %d below allowed minimum %d", length, min)
	}
	if max != nil && length > *max {
		return nil, gqlerror.Errorf("Field length %d above allowed maximum %d", length, *max)
	}
	return
}

func (d *directive) UniqueField(ctx context.Context, obj interface{}, next graphql.Resolver, typ, field string) (interface{}, error) {
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
