// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package directive

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/pkg/log"

	"github.com/99designs/gqlgen/graphql"
	"github.com/AlekSi/pointer"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"go.uber.org/zap"
)

type directive struct {
	logger  log.Logger
	regexps sync.Map
}

// New creates a graphql directive root.
func New(logger log.Logger) generated.DirectiveRoot {
	d := &directive{logger: logger}
	return generated.DirectiveRoot{
		DeprecatedInput: d.DeprecatedInput,
		NumberValue:     d.NumberValue,
		StringValue:     d.StringValue,
		List:            d.List,
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

type number float64

func (n number) String() string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}

type numbers []float64

func (nrs numbers) String() string {
	var b strings.Builder
	b.WriteByte('[')
	for i, n := range nrs {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(number(n).String())
	}
	b.WriteByte(']')
	return b.String()
}

func (d *directive) NumberValue(
	ctx context.Context, _ interface{}, next graphql.Resolver,
	multipleOf, max, min, exclusiveMax, exclusiveMin *float64,
	oneOf []float64, equals *float64,
) (res interface{}, err error) {
	if res, err = next(ctx); err != nil || d.isNil(res) {
		return
	}
	defer d.recoverFromPanic(ctx, &err)

	var value float64
	switch rv := reflect.Indirect(reflect.ValueOf(res)); rv.Kind() {
	case reflect.Float32, reflect.Float64:
		value = rv.Float()
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		value = float64(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		value = float64(rv.Uint())
	default:
		return nil, fmt.Errorf("@numberValue directive on kind %s", rv.Kind())
	}

	if multipleOf != nil {
		switch {
		case *multipleOf <= 0:
			return nil, fmt.Errorf("@numberValue directive multipleOf %s must be positive", number(*multipleOf))
		case math.Mod(value, *multipleOf) != 0:
			return nil, gqlerror.Errorf("Numeric value %s is not a multiple of %s", number(value), number(*multipleOf))
		}
	}
	if max != nil && value > *max {
		return nil, gqlerror.Errorf("Numeric value %s above allowed maximum %s", number(value), number(*max))
	}
	if min != nil && value < *min {
		return nil, gqlerror.Errorf("Numeric value %s below allowed minimum %s", number(value), number(*min))
	}
	if exclusiveMax != nil && value >= *exclusiveMax {
		return nil, gqlerror.Errorf("Numeric value %s above allowed exclusive maximum %s", number(value), number(*exclusiveMax))
	}
	if exclusiveMin != nil && value <= *exclusiveMin {
		return nil, gqlerror.Errorf("Numeric value %s below allowed exclusive minimum %s", number(value), number(*exclusiveMin))
	}
	if len(oneOf) > 0 && !func() bool {
		for _, one := range oneOf {
			if value == one {
				return true
			}
		}
		return false
	}() {
		return nil, gqlerror.Errorf("Field value %s must be one of %s", number(value), numbers(oneOf))
	}
	if equals != nil && value != *equals {
		return nil, gqlerror.Errorf("Numeric value %s must be equal to %s", number(value), number(*equals))
	}
	return
}

func (d *directive) StringValue(
	ctx context.Context, _ interface{}, next graphql.Resolver,
	maxLength, minLength *int, startsWith, endsWith, includes,
	regex *string, oneOf []string, equals *string,
) (res interface{}, err error) {
	if res, err = next(ctx); err != nil || d.isNil(res) {
		return
	}
	defer d.recoverFromPanic(ctx, &err)

	rv := reflect.Indirect(reflect.ValueOf(res))
	if rv.Kind() != reflect.String {
		return nil, fmt.Errorf("@stringValue directive on kind %s", rv.Kind())
	}
	value := rv.String()

	if maxLength != nil {
		switch length := len(value); {
		case *maxLength < 0:
			return nil, fmt.Errorf("@stringValue directive maxLength %d cannot be negative", *maxLength)
		case length > *maxLength:
			return nil, gqlerror.Errorf("String %q length %d above allowed maximum %d", value, length, *maxLength)
		}
	}
	if minLength != nil {
		switch length := len(value); {
		case *minLength < 0:
			return nil, fmt.Errorf("@stringValue directive minLength %d cannot be negative", *minLength)
		case length < *minLength:
			return nil, gqlerror.Errorf("String %q length %d below allowed minimum %d", value, length, *minLength)
		}
	}
	if startsWith != nil && !strings.HasPrefix(value, *startsWith) {
		return nil, gqlerror.Errorf("String %q must start with %q", value, *startsWith)
	}
	if endsWith != nil && !strings.HasSuffix(value, *endsWith) {
		return nil, gqlerror.Errorf("String %q must end with %q", value, *endsWith)
	}
	if includes != nil && !strings.Contains(value, *includes) {
		return nil, gqlerror.Errorf("String %q must include %q", value, *includes)
	}
	if regex != nil && !d.regexp(*regex).MatchString(value) {
		return nil, gqlerror.Errorf("String %q must match regular expression %q", value, *regex)
	}
	if len(oneOf) > 0 && !func() bool {
		for _, one := range oneOf {
			if value == one {
				return true
			}
		}
		return false
	}() {
		return nil, gqlerror.Errorf("String value %q must be one of %v", value, oneOf)
	}
	if equals != nil && value != *equals {
		return nil, gqlerror.Errorf("String value %q must be equal to %s", value, *equals)
	}
	return
}

func (d *directive) regexp(expr string) *regexp.Regexp {
	if re, ok := d.regexps.Load(expr); ok {
		return re.(*regexp.Regexp)
	}
	re := regexp.MustCompile(expr)
	d.regexps.Store(expr, re)
	return re
}

func (d *directive) List(
	ctx context.Context, _ interface{}, next graphql.Resolver,
	maxItems, minItems *int, uniqueItems *bool,
) (res interface{}, err error) {
	if res, err = next(ctx); err != nil || d.isNil(res) {
		return
	}
	defer d.recoverFromPanic(ctx, &err)

	value := reflect.Indirect(reflect.ValueOf(res))
	switch value.Kind() {
	case reflect.Array, reflect.Slice:
	default:
		return nil, fmt.Errorf("@list directive on kind %s", value.Kind())
	}
	length := value.Len()

	if maxItems != nil {
		switch {
		case *maxItems < 0:
			return nil, fmt.Errorf("@list directive maxItems %d cannot be negative", *maxItems)
		case length > *maxItems:
			return nil, gqlerror.Errorf("List length %d above allowed maximum %d", length, *maxItems)
		}
	}
	if minItems != nil {
		switch {
		case *minItems < 0:
			return nil, fmt.Errorf("@list directive minItems %d cannot be negative", *minItems)
		case length < *minItems:
			return nil, gqlerror.Errorf("List length %d below allowed minimum %d", length, *minItems)
		}
	}
	if pointer.GetBool(uniqueItems) {
		for i := 0; i < length-1; i++ {
			for j := i + 1; j < length; j++ {
				if reflect.DeepEqual(
					value.Index(i).Interface(),
					value.Index(j).Interface(),
				) {
					return nil, gqlerror.Errorf("List %v items not unique as item[%d] = item[%d]", res, i, j)
				}
			}
		}
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

func (d *directive) isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch v := reflect.ValueOf(i); v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice,
		reflect.Map, reflect.Chan, reflect.Func,
		reflect.UnsafePointer:
		return v.IsNil()
	default:
		return false
	}
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
