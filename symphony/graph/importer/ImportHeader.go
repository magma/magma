// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"fmt"

	"github.com/pkg/errors"
)

type ImportHeader struct {
	line     []string
	prnt3Idx int
	entity   ImportEntity
}

// NewImportHeader creates a new header to be used for import
func NewImportHeader(line []string, entity ImportEntity) (ImportHeader, error) {
	prnt3Idx := findStringContainsIndex(line, "Parent Equipment (3)")

	switch entity {
	case ImportEntityService, ImportEntityLocation:
	default:
		if prnt3Idx == -1 {
			return ImportHeader{}, errors.New("couldn't find Parent Equipment headers")
		}
	}

	return ImportHeader{
		line:     line,
		prnt3Idx: prnt3Idx,
		entity:   entity,
	}, nil
}

func (l ImportHeader) Find(s string) int {
	return findIndex(l.line, s)
}

func (l ImportHeader) NameIdx() int {
	if l.entity == ImportEntityPortInLink {
		return 0
	}
	return 1
}

func (l ImportHeader) PortEquipmentNameIdx() int {
	if l.entity == ImportEntityPort {
		return 3
	}
	if l.entity == ImportEntityPortInLink {
		return 1
	}
	return -1
}

func (l ImportHeader) PortEquipmentTypeNameIdx() int {
	if l.entity == ImportEntityPort {
		return 4
	}
	if l.entity == ImportEntityPortInLink {
		return 2
	}
	return -1
}

func (l ImportHeader) ExternalIDIdx() int {
	return findIndex(l.line, "External ID")
}

// LatitudeIdx returns the index of "latitude" column
func (l ImportHeader) LatitudeIdx() int {
	return findIndex(l.line, "Latitude")
}

// LongitudeIdx returns the index of "longitude" column
func (l ImportHeader) LongitudeIdx() int {
	return findIndex(l.line, "Longitude")
}

func (l ImportHeader) ThirdParentIdx() int {
	return l.prnt3Idx
}

func (l ImportHeader) ThirdPositionIdx() int {
	if l.entity == ImportEntityEquipment || l.entity == ImportEntityPortInLink {
		return l.prnt3Idx + 1
	}
	return -1
}

func (l ImportHeader) SecondParentIdx() int {
	if l.entity == ImportEntityEquipment || l.entity == ImportEntityPortInLink {
		return l.prnt3Idx + 2
	} else if l.entity == ImportEntityPort {
		return l.prnt3Idx + 1
	}
	return -1
}

func (l ImportHeader) SecondPositionIdx() int {
	if l.entity == ImportEntityEquipment || l.entity == ImportEntityPortInLink {
		return l.prnt3Idx + 3
	}
	return -1
}

func (l ImportHeader) DirectParentIdx() int {
	if l.entity == ImportEntityEquipment || l.entity == ImportEntityPortInLink {
		return l.prnt3Idx + 4
	} else if l.entity == ImportEntityPort {
		return l.prnt3Idx + 2
	}
	return -1
}

func (l ImportHeader) PositionIdx() int {
	if l.entity == ImportEntityEquipment || l.entity == ImportEntityPortInLink {
		return l.prnt3Idx + 5
	} else if l.entity == ImportEntityPort {
		return l.prnt3Idx + 3
	}
	return -1
}

func (l ImportHeader) LocationTypesRangeArr() []string {
	s, e := l.LocationsRangeIdx()
	return l.line[s:e]
}

func (l ImportHeader) LocationsRangeIdx() (int, int) {
	switch l.entity {
	case ImportEntityEquipment:
		return l.ExternalIDIdx() + 1, l.prnt3Idx
	case ImportEntityPort:
		return 5, l.prnt3Idx
	case ImportEntityPortInLink:
		return 3, l.prnt3Idx
	case ImportEntityLocation:
		return 1, l.ExternalIDIdx()
	}
	return -1, -1
}

// PropertyEndIdx is the index of last property on the file. currently it's always the last value
func (l ImportHeader) PropertyEndIdx() int {
	return len(l.line) - 1
}

func (l ImportHeader) PropertyStartIdx() int {
	switch l.entity {
	case ImportEntityEquipment:
		return l.PositionIdx() + 1
	case ImportEntityPort:
		return l.PositionIdx() + 6
	case ImportEntityService:
		return l.StatusIdx() + 1
	case ImportEntityLink:
		return l.LinkSecondPortStartIdx() * 2
	case ImportEntityLocation:
		return l.ExternalIDIdx() + 3
	}
	return -1
}

// ServiceExternalIDIdx is the index of the external id of the service (used in other systems) in the exported csv
func (l ImportHeader) ServiceExternalIDIdx() int {
	if l.entity == ImportEntityService {
		return 3
	}
	return -1
}

// CustomerNameIdx is the index of the name of customer that uses the services in the exported csv
func (l ImportHeader) CustomerNameIdx() int {
	if l.entity == ImportEntityService {
		return 4
	}
	return -1
}

// CustomerExternalIDIdx is the index of the external id of customer that uses the services in the exported csv
func (l ImportHeader) CustomerExternalIDIdx() int {
	if l.entity == ImportEntityService {
		return 5
	}
	return -1
}

// StatusIdx is the index of the status of the service (can be of types enum ServiceType in graphql) in the exported csv
func (l ImportHeader) StatusIdx() int {
	if l.entity == ImportEntityService {
		return 6
	}
	return -1
}

func (l ImportHeader) ServiceNamesIdx() int {
	return findIndex(l.line, "Service Names")
}

func (l ImportHeader) LinkGetTwoPortsRange() ([]int, []int) {
	if l.entity == ImportEntityLink {
		splitIdx := l.LinkSecondPortStartIdx()
		return []int{1, splitIdx}, []int{splitIdx, l.ServiceNamesIdx()}
	}
	return nil, nil
}

// LinkGetTwoPortsSlices get metric of two slices, one for each port
func (l ImportHeader) LinkGetTwoPortsSlices() ([][]string, error) {
	if l.entity == ImportEntityLink {
		idxA, idxB := l.LinkGetTwoPortsRange()
		if idxA[0] == -1 || idxA[1] == -1 {
			return nil, errors.New("one of the port headers is missing")
		}
		if idxB[0] == -1 || idxB[1] == -1 {
			return nil, errors.New("one of the port B headers, or 'Service Names' column is missing")
		}
		return [][]string{l.line[idxA[0]:idxA[1]], l.line[idxB[0]:idxB[1]]}, nil
	}
	return nil, fmt.Errorf("invalid entity %v", l.entity)
}

func (l ImportHeader) LinkSecondPortStartIdx() int {
	if l.entity == ImportEntityLink {
		return findIndex(l.line, "Port B Name")
	}
	return -1
}

func (l ImportHeader) LinkLocationsRangesIdx() ([]int, []int) {
	prnt3ForSecondPortIdx := findIndex(l.line, "Parent Equipment (3) B")
	if l.entity == ImportEntityLink {
		return []int{4, l.prnt3Idx}, []int{l.LinkSecondPortStartIdx() + 3, prnt3ForSecondPortIdx}
	}
	return nil, nil
}
