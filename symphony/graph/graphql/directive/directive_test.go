// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package directive

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/log/logtest"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testResolver struct {
	mock.Mock
}

func (tr *testResolver) resolve(ctx context.Context) (interface{}, error) {
	args := tr.Called(ctx)
	return args.Get(0), args.Error(1)
}

func TestDirectiveLength(t *testing.T) {
	var (
		tests = []struct {
			name    string
			input   interface{}
			err     error
			min     int
			max     *int
			wantErr bool
		}{
			{
				name:  "Valid",
				input: "test",
				min:   1,
			},
			{
				name:    "TooShort",
				input:   "foo",
				min:     5,
				wantErr: true,
			},
			{
				name:    "TooLong",
				input:   []int{1, 2, 3},
				max:     pointer.ToInt(2),
				wantErr: true,
			},
			{
				name:  "Unlimited",
				input: "hello world",
			},
			{
				name:    "Unresolved",
				input:   "test",
				err:     errors.New("bad resolver"),
				wantErr: true,
			},
			{
				name:    "NoLength",
				input:   42,
				min:     10,
				wantErr: true,
			},
		}
		d      = New(logtest.NewTestLogger(t))
		length = func(min int, max *int) func(interface{}, graphql.Resolver) (interface{}, error) {
			return func(in interface{}, next graphql.Resolver) (interface{}, error) {
				return d.Length(context.Background(), in, next, min, max)
			}
		}
	)
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var tr testResolver
			tr.On("resolve", mock.Anything).
				Return(tc.input, tc.err).
				Once()
			defer tr.AssertExpectations(t)

			output, err := length(tc.min, tc.max)(tc.input, tr.resolve)
			if !tc.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, tc.input, output)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestDirectiveRange(t *testing.T) {
	var (
		tests = []struct {
			name    string
			input   interface{}
			err     error
			min     *float64
			max     *float64
			wantErr bool
		}{
			{
				name:  "Valid",
				input: 42,
				min:   pointer.ToFloat64(10),
				max:   pointer.ToFloat64(50),
			},
			{
				name:    "TooSmall",
				input:   -5,
				min:     pointer.ToFloat64(0),
				wantErr: true,
			},
			{
				name:    "TooBig",
				input:   100,
				min:     pointer.ToFloat64(0),
				max:     pointer.ToFloat64(99),
				wantErr: true,
			},
			{
				name:    "NotInt",
				input:   "55",
				wantErr: true,
			},
		}
		d       = New(logtest.NewTestLogger(t))
		rangefn = func(min, max *float64) func(interface{}, graphql.Resolver) (interface{}, error) {
			return func(in interface{}, next graphql.Resolver) (interface{}, error) {
				return d.Range(context.Background(), in, next, min, max)
			}
		}
	)
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var tr testResolver
			tr.On("resolve", mock.Anything).
				Return(tc.input, tc.err).
				Once()
			defer tr.AssertExpectations(t)

			output, err := rangefn(tc.min, tc.max)(tc.input, tr.resolve)
			if !tc.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, tc.input, output)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestDirectiveUniqueField(t *testing.T) {
	var (
		d          = New(logtest.NewTestLogger(t))
		uniqueName = func(in interface{}, next graphql.Resolver) (interface{}, error) {
			return d.UniqueField(context.Background(), in, next, "property type", "Name")
		}
	)
	t.Run("SliceOfStructPointers", func(t *testing.T) {
		inputs := []*models.PropertyTypeInput{
			nil,
			{Name: "foo"},
			{Name: "bar"},
			{Name: "foo"},
			{Name: "baz"},
		}
		var tr testResolver
		tr.On("resolve", mock.Anything).
			Return(inputs, nil).
			Once()
		defer tr.AssertExpectations(t)

		outputs, err := uniqueName(inputs, tr.resolve)
		assert.Nil(t, outputs)
		assert.Error(t, err)
	})
	t.Run("SliceOfStructs", func(t *testing.T) {
		inputs := []struct {
			Name string
		}{
			{"foo"},
			{"bar"},
			{"baz"},
		}
		var tr testResolver
		tr.On("resolve", mock.Anything).
			Return(inputs, nil).
			Once()
		defer tr.AssertExpectations(t)

		outputs, err := uniqueName(inputs, tr.resolve)
		assert.Equal(t, inputs, outputs)
		assert.NoError(t, err)
	})
	t.Run("NoName", func(t *testing.T) {
		inputs := []struct {
			LastName string
		}{
			{"bar"},
			{"baz"},
		}
		var tr testResolver
		tr.On("resolve", mock.Anything).
			Return(inputs, nil).
			Once()
		defer tr.AssertExpectations(t)

		outputs, err := uniqueName(inputs, tr.resolve)
		assert.Equal(t, inputs, outputs)
		assert.NoError(t, err)
	})
	t.Run("NoStringName", func(t *testing.T) {
		inputs := []struct {
			Name int
		}{
			{1},
			{1},
		}
		var tr testResolver
		tr.On("resolve", mock.Anything).
			Return(inputs, nil).
			Once()
		defer tr.AssertExpectations(t)

		outputs, err := uniqueName(inputs, tr.resolve)
		assert.Equal(t, inputs, outputs)
		assert.NoError(t, err)
	})
	t.Run("NonSlice", func(t *testing.T) {
		var (
			tr testResolver
			nr = 42
		)
		tr.On("resolve", mock.Anything).
			Return(nr, nil).
			Once()
		defer tr.AssertExpectations(t)

		r, err := uniqueName(nr, tr.resolve)
		assert.Equal(t, nr, r)
		assert.NoError(t, err)
	})
	t.Run("IntPointerSlice", func(t *testing.T) {
		gen := func(v int) *int { return &v }
		var (
			tr testResolver
			nr = []*int{gen(1), gen(1)}
		)
		tr.On("resolve", mock.Anything).
			Return(nr, nil).
			Once()
		defer tr.AssertExpectations(t)

		r, err := uniqueName(nr, tr.resolve)
		assert.Equal(t, nr, r)
		assert.NoError(t, err)
	})
	t.Run("FuncSlice", func(t *testing.T) {
		inputs := []func(){func() {}}
		var tr testResolver
		tr.On("resolve", mock.Anything).
			Return(inputs, nil).
			Once()
		defer tr.AssertExpectations(t)

		outputs, err := uniqueName(inputs, tr.resolve)
		assert.Equal(t, inputs, outputs)
		assert.NoError(t, err)
	})
	t.Run("ResolveError", func(t *testing.T) {
		var tr testResolver
		tr.On("resolve", mock.Anything).
			Return(nil, io.EOF).
			Once()
		defer tr.AssertExpectations(t)

		r, err := uniqueName(42, tr.resolve)
		assert.Nil(t, r)
		assert.EqualError(t, err, io.EOF.Error())
	})
}
