/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnionFind(t *testing.T) {
	uf := newUnionFind([]string{"1", "2", "3", "4", "5"})

	// 1 | 2 | 3 | 4 | 5
	actual := uf.getComponents()
	assert.Equal(
		t,
		[][]string{
			{"1"},
			{"2"},
			{"3"},
			{"4"},
			{"5"},
		},
		actual,
	)

	uf.union("1", "2")
	uf.union("3", "4")
	// 1 -> 2 | 3 -> 4 | 5
	actual = uf.getComponents()
	assert.Equal(
		t,
		[][]string{
			{"5"},
			{"1", "2"},
			{"3", "4"},
		},
		actual,
	)

	uf.union("4", "2")
	// 1 -> (2, 3 -> 4) | 5
	actual = uf.getComponents()
	assert.Equal(
		t,
		[][]string{
			{"5"},
			{"1", "2", "3", "4"},
		},
		actual,
	)

	// paths were compressed from last call to getComponents
	uf.union("4", "5")
	// 1 -> (2, 3, 4, 5)
	actual = uf.getComponents()
	assert.Equal(
		t,
		[][]string{
			{"1", "2", "3", "4", "5"},
		},
		actual,
	)

	uf = newUnionFind([]string{"1", "2", "3", "4", "5"})
	uf.union("1", "2")
	uf.union("2", "3")
	uf.union("4", "5")
	// 1 -> (2, 3) | 4 -> 5
	uf.union("5", "3")
	// 1 -> (2, 3, 4 -> 5)
	// 4 -> (1 -> [2, 3], 5)
	actual = uf.getComponents()
	assert.Equal(
		t,
		[][]string{
			{"1", "2", "3", "4", "5"},
		},
		actual,
	)
}
