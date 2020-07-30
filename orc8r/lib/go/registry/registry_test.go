/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package registry

import (
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestServiceRegistry_GetAnnotationFields(t *testing.T) {
	tests := []struct {
		name            string
		annotationValue string
		separator       string
		want            []string
	}{
		{
			name:            "empty",
			annotationValue: "",
			separator:       ",",
			want:            []string{},
		},
		{
			name:            "all whitespace",
			annotationValue: "  \n\n  ",
			separator:       ",",
			want:            []string{},
		},
		{
			name:            "single element",
			annotationValue: "42",
			separator:       ",",
			want:            []string{"42"},
		},
		{
			name:            "multiple elements",
			annotationValue: "42,foo",
			separator:       ",",
			want:            []string{"42", "foo"},
		},
		{
			name:            "multiple elements with whitespace",
			annotationValue: "  42 ,\n  foo  ",
			separator:       ",",
			want:            []string{"42", "foo"},
		},
		{
			name:            "trailing separator",
			annotationValue: "  a,       b, c,\n\nd,    e,\n\n  f,  \n  ",
			separator:       ",",
			want:            []string{"a", "b", "c", "d", "e", "f"},
		},
		{
			name:            "space separator",
			annotationValue: "  a,       b, c,\n\nd,    e,\n\n  f,  \n  ",
			separator:       " ",
			want:            []string{"a,", "b,", "c,d,", "e,", "f,"},
		},
		{
			name:            "newline separator",
			annotationValue: "  a,       b, c,\n\nd,    e,\n\n  f,  \n  ",
			separator:       "\n",
			want:            []string{"a,b,c,", "d,e,", "f,"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ServiceRegistry{
				ServiceLocations: map[string]ServiceLocation{
					"srv": {Annotations: map[string]string{"annotationName": tt.annotationValue}},
				},
			}
			got, err := r.GetAnnotationFields("srv", "annotationName", tt.separator)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
