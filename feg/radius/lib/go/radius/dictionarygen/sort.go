/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dictionarygen

import "sort"

func sortExternalAttributes(e []*externalAttribute) {
	sort.Stable(sortableExternalAttributes(e))
}

type sortableExternalAttributes []*externalAttribute

func (s sortableExternalAttributes) Len() int           { return len(s) }
func (s sortableExternalAttributes) Less(i, j int) bool { return s[i].Attribute < s[j].Attribute }
func (s sortableExternalAttributes) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
