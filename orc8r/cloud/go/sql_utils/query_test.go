/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package sql_utils_test

import (
	"testing"

	"magma/orc8r/cloud/go/sql_utils"

	"github.com/stretchr/testify/assert"
)

func TestGetPlaceholderArgListWithSuffix(t *testing.T) {
	testCases := []struct {
		startIdx int
		numArgs  int
		suffix   string
		expected string
	}{
		{1, 3, "hello", "($1, $2, $3, hello)"},
		{1, 3, "", "($1, $2, $3)"},
		{1, 1, "world", "($1, world)"},
		{1, 1, "", "($1)"},
		{1, 0, "", "()"},
		{1, 0, "foo", "(foo)"},
		{5, 3, "bar", "($5, $6, $7, bar)"},
		{5, 3, "", "($5, $6, $7)"},
		{5, 1, "baz", "($5, baz)"},
		{5, 1, "", "($5)"},
		{5, 0, "qux", "(qux)"},
		{5, 0, "", "()"},
	}

	for _, testCase := range testCases {
		actual := sql_utils.GetPlaceholderArgListWithSuffix(testCase.startIdx, testCase.numArgs, testCase.suffix)
		assert.Equal(t, testCase.expected, actual)

		if testCase.suffix == "" {
			assert.Equal(t, actual, sql_utils.GetPlaceholderArgList(testCase.startIdx, testCase.numArgs))
		}
	}
}

func TestGetUpdateClauseString(t *testing.T) {
	testCases := []struct {
		startIdx int
		argNames []string
		expected string
	}{
		{1, []string{"foo", "bar", "baz"}, "foo = $1, bar = $2, baz = $3"},
		{1, []string{"foo"}, "foo = $1"},
		{1, []string{}, ""},
		{5, []string{"foo", "bar", "baz"}, "foo = $5, bar = $6, baz = $7"},
		{5, []string{"foo"}, "foo = $5"},
		{5, []string{}, ""},
	}

	for _, testCase := range testCases {
		actual := sql_utils.GetUpdateClauseString(testCase.startIdx, testCase.argNames...)
		assert.Equal(t, testCase.expected, actual)
	}
}

func TestGetInsertArgListString(t *testing.T) {
	testCases := []struct {
		args     []string
		expected string
	}{
		{[]string{"foo", "bar", "baz"}, "(foo, bar, baz)"},
		{[]string{"foo"}, "(foo)"},
		{[]string{}, "()"},
	}

	for _, testCase := range testCases {
		actual := sql_utils.GetInsertArgListString(testCase.args...)
		assert.Equal(t, testCase.expected, actual)
	}
}
