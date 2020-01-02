// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

type ImportHeader struct {
	line     []string
	prnt3Idx int
	entity   ImportEntity
}

// NewImportHeader creates a new header to be used for import
func NewImportHeader(line []string, entity ImportEntity) ImportHeader {
	prnt3Idx := findIndex(line, "Parent Equipment (3)")
	return ImportHeader{
		line:     line,
		prnt3Idx: prnt3Idx,
		entity:   entity,
	}
}

func (l ImportHeader) Find(s string) int {
	return findIndex(l.line, s)
}

func (l ImportHeader) ExternalIDIdx() int {
	return findIndex(l.line, "External ID")
}

func (l ImportHeader) ThirdParentIdx() int {
	return l.prnt3Idx
}

func (l ImportHeader) ThirdPositionIdx() int {
	if l.entity == ImportEntityEquipment {
		return l.prnt3Idx + 1
	}
	return -1
}

func (l ImportHeader) SecondParentIdx() int {
	if l.entity == ImportEntityEquipment {
		return l.prnt3Idx + 2
	} else if l.entity == ImportEntityPort {
		return l.prnt3Idx + 1
	}
	return -1
}

func (l ImportHeader) SecondPositionIdx() int {
	if l.entity == ImportEntityEquipment {
		return l.prnt3Idx + 3
	}
	return -1
}

func (l ImportHeader) DirectParentIdx() int {
	if l.entity == ImportEntityEquipment {
		return l.prnt3Idx + 4
	} else if l.entity == ImportEntityPort {
		return l.prnt3Idx + 2
	}
	return -1
}

func (l ImportHeader) PositionIdx() int {
	if l.entity == ImportEntityEquipment {
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
	if l.entity == ImportEntityEquipment {
		return l.ExternalIDIdx() + 1, l.prnt3Idx
	} else if l.entity == ImportEntityPort {
		return 5, l.prnt3Idx
	}
	return -1, -1
}

func (l ImportHeader) PropertyStartIdx() int {
	switch l.entity {
	case ImportEntityEquipment:
		return l.PositionIdx() + 1
	case ImportEntityPort:
		return l.PositionIdx() + 5
	case ImportEntityService:
		return l.StatusIdx() + 1
	case ImportEntityLink:
		return minimalLinksLineLength()
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

func (l ImportHeader) PortAIDIdx() int {
	if l.entity == ImportEntityLink {
		return findIndex(l.line, "Port A ID")
	}
	return -1
}

func (l ImportHeader) PortBIDIdx() int {
	if l.entity == ImportEntityLink {
		return findIndex(l.line, "Port B ID")
	}
	return -1
}
