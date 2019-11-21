// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

type ImportHeader struct {
	line     []string
	prnt3Idx int
}

func NewImportHeader(line []string) ImportHeader {
	prnt3Idx := findIndex(line, "Parent Equipment (3)")
	return ImportHeader{
		line:     line,
		prnt3Idx: prnt3Idx,
	}
}

func (l ImportHeader) Find(s string) int {
	return findIndex(l.line, s)
}

func (l ImportHeader) ThirdParentIdx() int {
	return l.prnt3Idx
}

func (l ImportHeader) SecondParentIdx() int {
	return l.prnt3Idx + 1
}

func (l ImportHeader) DirectParentIdx() int {
	return l.prnt3Idx + 2
}

func (l ImportHeader) LocationTypesRangeArr() []string {
	s, e := l.LocationsRangeIdx()
	return l.line[s:e]
}

func (l ImportHeader) LocationsRangeIdx() (int, int) {
	return 3, l.prnt3Idx
}

func (l ImportHeader) PositionIdx() int {
	return l.prnt3Idx + 3
}

func (l ImportHeader) PropertyStartIdx() int {
	return l.PositionIdx() + 1
}
