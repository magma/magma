// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package directive_test

import (
	"context"
	"io"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/facebookincubator/symphony/graph/graphql/directive"
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

func TestDirectiveNumberValue(t *testing.T) {
	var (
		tests = []struct {
			name         string
			input        interface{}
			err          error
			multipleOf   *float64
			max          *float64
			min          *float64
			exclusiveMax *float64
			exclusiveMin *float64
			oneOf        []float64
			equals       *float64
			wantErr      bool
		}{
			{
				name:       "Valid",
				input:      uint(42),
				multipleOf: pointer.ToFloat64(3),
				max:        pointer.ToFloat64(50),
				min:        pointer.ToFloat64(40),
				oneOf:      []float64{10, 42, 66, 70},
				equals:     pointer.ToFloat64(42),
			},
			{
				name:    "NextError",
				err:     io.ErrUnexpectedEOF,
				wantErr: true,
			},
			{
				name:    "NilInput",
				input:   pointer.ToIntOrNil(0),
				wantErr: false,
			},
			{
				name:       "NotMultipleOf",
				input:      7,
				multipleOf: pointer.ToFloat64(2),
				wantErr:    true,
			},
			{
				name:       "ZeroMultipleOf",
				input:      float32(1),
				multipleOf: pointer.ToFloat64(0),
				wantErr:    true,
			},
			{
				name:       "NegativeMultipleOf",
				input:      float32(2),
				multipleOf: pointer.ToFloat64(-5),
				wantErr:    true,
			},
			{
				name:    "AboveMaximum",
				input:   51,
				max:     pointer.ToFloat64(50),
				wantErr: true,
			},
			{
				name:    "BelowMinimum",
				input:   -7,
				min:     pointer.ToFloat64(0),
				wantErr: true,
			},
			{
				name:         "AboveExclusiveMaximum",
				input:        50,
				exclusiveMax: pointer.ToFloat64(50),
				wantErr:      true,
			},
			{
				name:         "BelowExclusiveMinimum",
				input:        float64(0),
				exclusiveMin: pointer.ToFloat64(0),
				wantErr:      true,
			},
			{
				name:    "NotOneOf",
				input:   15,
				oneOf:   []float64{10, 16, 18},
				wantErr: true,
			},
			{
				name:    "NotEquals",
				input:   21,
				equals:  pointer.ToFloat64(20),
				wantErr: true,
			},
		}
		d               = directive.New(logtest.NewTestLogger(t))
		numberValueFunc = func(
			multipleOf, max, min, exclusiveMax, exclusiveMin *float64,
			oneOf []float64, equals *float64,
		) func(interface{}, graphql.Resolver) (interface{}, error) {
			return func(in interface{}, next graphql.Resolver) (interface{}, error) {
				return d.NumberValue(
					context.Background(), in, next,
					multipleOf, max, min, exclusiveMax,
					exclusiveMin, oneOf, equals,
				)
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

			output, err := numberValueFunc(
				tc.multipleOf, tc.max, tc.min,
				tc.exclusiveMax, tc.exclusiveMin,
				tc.oneOf, tc.equals,
			)(tc.input, tr.resolve)
			if !tc.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, tc.input, output)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestDirectiveStringValue(t *testing.T) {
	var (
		tests = []struct {
			name       string
			input      interface{}
			err        error
			minLength  *int
			maxLength  *int
			startsWith *string
			endsWith   *string
			includes   *string
			regex      *string
			oneOf      []string
			equals     *string
			wantErr    bool
		}{
			{
				name:       "Valid",
				input:      pointer.ToString("foobarbaz"),
				minLength:  pointer.ToInt(3),
				maxLength:  pointer.ToInt(10),
				startsWith: pointer.ToString("foo"),
				endsWith:   pointer.ToString("baz"),
				includes:   pointer.ToString("bar"),
				regex:      pointer.ToString(`^f.+z$`),
				oneOf:      []string{"foo", "bar", "foobarbaz"},
				equals:     pointer.ToString("foobarbaz"),
			},
			{
				name:    "NextError",
				err:     io.ErrUnexpectedEOF,
				wantErr: true,
			},
			{
				name:    "NilInput",
				input:   pointer.ToStringOrNil(""),
				wantErr: false,
			},
			{
				name:      "TooLong",
				input:     "hello",
				maxLength: pointer.ToInt(4),
				wantErr:   true,
			},
			{
				name:      "NegativeMaxLength",
				input:     "hello",
				maxLength: pointer.ToInt(-5),
				wantErr:   true,
			},
			{
				name:      "TooShort",
				input:     "hello",
				minLength: pointer.ToInt(10),
				wantErr:   true,
			},
			{
				name:      "NegativeMinLength",
				input:     "hello",
				maxLength: pointer.ToInt(-100),
				wantErr:   true,
			},
			{
				name:       "NoPrefix",
				input:      "world",
				startsWith: pointer.ToString("he"),
				wantErr:    true,
			},
			{
				name:     "NoSuffix",
				input:    "world",
				endsWith: pointer.ToString("lld"),
				wantErr:  true,
			},
			{
				name:     "NoContains",
				input:    "world",
				includes: pointer.ToString("rrl"),
				wantErr:  true,
			},
			{
				name:    "NoRegexMatch",
				input:   "world",
				regex:   pointer.ToString("^wworld$"),
				wantErr: true,
			},
			{
				name:    "NoOneOf",
				input:   "bar",
				oneOf:   []string{"foo", "baz"},
				wantErr: true,
			},
			{
				name:    "NoEquals",
				input:   "bar",
				equals:  pointer.ToString("baz"),
				wantErr: true,
			},
		}
		d               = directive.New(logtest.NewTestLogger(t))
		stringValueFunc = func(
			maxLength, minLength *int, startsWith, endsWith, includes,
			regex *string, oneOf []string, equals *string,
		) func(interface{}, graphql.Resolver) (interface{}, error) {
			return func(in interface{}, next graphql.Resolver) (interface{}, error) {
				return d.StringValue(
					context.Background(), in, next, maxLength,
					minLength, startsWith, endsWith, includes,
					regex, oneOf, equals,
				)
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

			output, err := stringValueFunc(
				tc.maxLength, tc.minLength, tc.startsWith,
				tc.endsWith, tc.includes, tc.regex,
				tc.oneOf, tc.equals,
			)(tc.input, tr.resolve)
			if !tc.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, tc.input, output)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestDirectiveList(t *testing.T) {
	var (
		tests = []struct {
			name        string
			input       interface{}
			err         error
			maxItems    *int
			minItems    *int
			uniqueItems *bool
			wantErr     bool
		}{
			{
				name:        "Valid",
				input:       []int{1, 2, 3},
				maxItems:    pointer.ToInt(5),
				minItems:    pointer.ToInt(3),
				uniqueItems: pointer.ToBool(true),
			},
			{
				name:    "NextError",
				input:   []int{1, 2, 3},
				err:     io.ErrUnexpectedEOF,
				wantErr: true,
			},
			{
				name:    "NilInput",
				input:   []uint8(nil),
				wantErr: false,
			},
			{
				name:     "AboveMaxItems",
				input:    []int8{1, 2, 3, 4, 5},
				maxItems: pointer.ToInt(4),
				wantErr:  true,
			},
			{
				name:     "NegativeMaxItems",
				input:    []uint{42},
				maxItems: pointer.ToInt(-6),
				wantErr:  true,
			},
			{
				name:     "BelowMinItems",
				input:    []string{"foo"},
				minItems: pointer.ToInt(2),
				wantErr:  true,
			},
			{
				name:     "NegativeMinItems",
				input:    []byte{8},
				minItems: pointer.ToInt(-42),
				wantErr:  true,
			},
			{
				name:        "NonUniqueStrings",
				input:       []string{"foo", "bar", "baz", "baz"},
				uniqueItems: pointer.ToBool(true),
				wantErr:     true,
			},
			{
				name:        "NonUniqueUint64s",
				input:       []uint64{100, 200, 300, 100, 500},
				uniqueItems: pointer.ToBool(true),
				wantErr:     true,
			},
			{
				name: "UniqueStructs",
				input: []struct {
					name string
					age  uint
				}{
					{
						name: "foo",
						age:  18,
					},
					{
						name: "bar",
						age:  19,
					},
					{
						name: "foo",
						age:  19,
					},
				},
				uniqueItems: pointer.ToBool(true),
			},
			{
				name: "NonUniqueStructs",
				input: []*struct {
					name string
					age  uint
				}{
					{
						name: "foo",
						age:  18,
					},
					{
						name: "bar",
						age:  19,
					},
					{
						name: "foo",
						age:  18,
					},
				},
				uniqueItems: pointer.ToBool(true),
				wantErr:     true,
			},
		}
		d        = directive.New(logtest.NewTestLogger(t))
		listFunc = func(maxItems, minItems *int, uniqueItems *bool) func(interface{}, graphql.Resolver) (interface{}, error) {
			return func(in interface{}, next graphql.Resolver) (interface{}, error) {
				return d.List(context.Background(), in, next, maxItems, minItems, uniqueItems)
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

			output, err := listFunc(tc.maxItems, tc.minItems, tc.uniqueItems)(tc.input, tr.resolve)
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
		d          = directive.New(logtest.NewTestLogger(t))
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

func TestDirectiveDeprecatedInputField(t *testing.T) {
	var (
		d                    = directive.New(logtest.NewTestLogger(t))
		deprecatedInputField = func(in interface{}, next graphql.Resolver) (interface{}, error) {
			return d.DeprecatedInput(context.Background(), in, next, "AddInput.input", "Don't use both", pointer.ToString("input2"))
		}
	)
	t.Run("OnlyDeprecatedField", func(t *testing.T) {
		input := map[string]interface{}{
			"input": "Valid",
		}
		var tr testResolver
		tr.On("resolve", mock.Anything).
			Return(input["input"], nil).
			Once()
		defer tr.AssertExpectations(t)

		outputs, err := deprecatedInputField(input, tr.resolve)
		assert.Equal(t, "Valid", outputs)
		assert.NoError(t, err)
	})
	t.Run("OnlyNewField", func(t *testing.T) {
		input := map[string]interface{}{
			"input2": "Valid",
		}
		var tr testResolver
		tr.On("resolve", mock.Anything).
			Return(input["input"], nil).
			Once()
		defer tr.AssertExpectations(t)

		outputs, err := deprecatedInputField(input, tr.resolve)
		assert.Equal(t, nil, outputs)
		assert.NoError(t, err)
	})
	t.Run("BothFields", func(t *testing.T) {
		input := map[string]interface{}{
			"input":  "Valid",
			"input2": "Invalid",
		}
		var tr testResolver
		tr.On("resolve", mock.Anything).
			Return(input["input"], nil).
			Once()
		defer tr.AssertExpectations(t)

		outputs, err := deprecatedInputField(input, tr.resolve)
		assert.Nil(t, outputs)
		assert.Error(t, err)
	})
	t.Run("BothFieldWithDeprecatedEmpty", func(t *testing.T) {
		input := map[string]interface{}{
			"input":  "",
			"input2": "Valid",
		}
		var tr testResolver
		tr.On("resolve", mock.Anything).
			Return(input["input"], nil).
			Once()
		defer tr.AssertExpectations(t)

		outputs, err := deprecatedInputField(input, tr.resolve)
		assert.Nil(t, outputs)
		assert.Error(t, err)
	})
	t.Run("BothFieldWithNewEmpty", func(t *testing.T) {
		input := map[string]interface{}{
			"input":  "Valid",
			"input2": "",
		}
		var tr testResolver
		tr.On("resolve", mock.Anything).
			Return(input["input"], nil).
			Once()
		defer tr.AssertExpectations(t)

		outputs, err := deprecatedInputField(input, tr.resolve)
		assert.Nil(t, outputs)
		assert.Error(t, err)
	})
}
