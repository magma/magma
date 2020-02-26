/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package indexer_test

import (
	"github.com/stretchr/testify/assert"
	"magma/orc8r/cloud/go/services/state/indexer"
	"testing"
)

var emptyMatch []string

func TestMatchAll_Match(t *testing.T) {
	m := indexer.NewMatchAll()
	yesMatch := []string{"", "f", "fo", "foo", "foobar", "foobarbaz", "foobar !@#$%^&*()_+"}
	testMatchImpl(t, m, yesMatch, emptyMatch)
}

func TestMatchExact_Match(t *testing.T) {
	m := indexer.NewMatchExact("foobar")
	yesMatch := []string{"foobar"}
	noMatch := []string{"", "f", "fo", "foo", "foobarbaz", "foobar !@#$%^&*()_+"}
	testMatchImpl(t, m, yesMatch, noMatch)
}

func TestMatchPrefix_Match(t *testing.T) {
	m := indexer.NewMatchPrefix("foo")
	yesMatch := []string{"foo", "foobar", "foobarbaz", "foobar !@#$%^&*()_+"}
	noMatch := []string{"", "f", "fo"}
	testMatchImpl(t, m, yesMatch, noMatch)
}

func testMatchImpl(t *testing.T, m indexer.KeyMatcher, yesMatch, noMatch []string) {
	for _, str := range yesMatch {
		assert.True(t, m.Match(str))
	}
	for _, str := range noMatch {
		assert.False(t, m.Match(str))
	}
}
