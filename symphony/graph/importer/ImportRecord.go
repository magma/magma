// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent/service"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type ImportRecord struct {
	line  []string
	title ImportHeader
}

// PortData is the data structure for PortData function
type PortData struct {
	ID                string
	Name              string
	TypeName          string
	EquipmentID       *string
	EquipmentName     string
	EquipmentTypeName string
}

func NewImportRecord(line []string, title ImportHeader) ImportRecord {
	return ImportRecord{
		line:  line,
		title: title,
	}
}

func (l ImportRecord) ZapField() zap.Field {
	return zap.Strings("line", l.line)
}

func (l ImportRecord) Len() int {
	return len(l.line)
}

func (l ImportRecord) Header() ImportHeader {
	return l.title
}

// GetPropertyInput returns a PropertyInput model from a proptypeName
func (l ImportRecord) GetPropertyInput(client *ent.Client, ctx context.Context, typ interface{}, proptypeName string) (*models.PropertyInput, error) {
	var pTyp *ent.PropertyType
	var err error
	switch l.entity() {
	case ImportEntityEquipment:
		typ := typ.(*ent.EquipmentType)
		pTyp, err = typ.QueryPropertyTypes().Where(propertytype.Name(proptypeName)).Only(ctx)
	case ImportEntityPort:
		typ := typ.(*ent.EquipmentPortType)
		pTyp, err = typ.QueryPropertyTypes().Where(propertytype.Name(proptypeName)).Only(ctx)
	case ImportEntityLink:
		typ := typ.(*ent.EquipmentPortType)
		pTyp, err = typ.QueryLinkPropertyTypes().Where(propertytype.Name(proptypeName)).Only(ctx)
	case ImportEntityService:
		typ := typ.(*ent.ServiceType)
		pTyp, err = typ.QueryPropertyTypes().Where(propertytype.Name(proptypeName)).Only(ctx)
	default:
		return nil, errors.Wrapf(err, "entity is not supported %s", l.entity())
	}
	if err != nil {
		return nil, errors.Wrapf(err, "property type does not exist %q", proptypeName)
	}

	idx := l.title.Find(proptypeName)
	if idx == -1 {
		return nil, nil
	}
	value := l.line[idx]
	if pTyp.Type == "service" && value != "" {
		if value, err = client.Service.Query().Where(service.Name(value)).OnlyID(ctx); err != nil {
			return nil, errors.Wrapf(err, "service name does not exist %q", l.line[idx])
		}
	}
	return getPropInput(*pTyp, value)
}

func (l ImportRecord) entity() ImportEntity {
	return l.Header().entity
}

func (l ImportRecord) ID() string {
	return l.line[0]
}

func (l ImportRecord) Name() string {
	return l.line[l.Header().NameIdx()]
}

func (l ImportRecord) TypeName() string {
	return l.line[2]
}

func (l ImportRecord) PortEquipmentName() string {
	return l.line[l.Header().PortEquipmentNameIdx()]
}

func (l ImportRecord) PortEquipmentTypeName() string {
	return l.line[l.Header().PortEquipmentTypeNameIdx()]
}

func (l ImportRecord) ExternalID() string {
	return l.line[l.title.ExternalIDIdx()]
}

func (l ImportRecord) ThirdParent() string {
	return l.line[l.title.ThirdParentIdx()]
}

func (l ImportRecord) ThirdPosition() string {
	return l.line[l.title.ThirdPositionIdx()]
}

func (l ImportRecord) SecondParent() string {
	return l.line[l.title.SecondParentIdx()]
}

func (l ImportRecord) SecondPosition() string {
	return l.line[l.title.SecondPositionIdx()]
}

func (l ImportRecord) DirectParent() string {
	return l.line[l.title.DirectParentIdx()]
}

func (l ImportRecord) Position() string {
	return l.line[l.title.PositionIdx()]
}

func (l ImportRecord) LocationsRangeArr() []string {
	s, e := l.title.LocationsRangeIdx()
	return l.line[s:e]
}

func (l ImportRecord) PropertiesMap() map[string]string {
	valueSlice := l.line[l.title.PropertyStartIdx():]
	typeSlice := l.title.line[l.title.PropertyStartIdx():]
	ret := make(map[string]string, len(valueSlice))
	for i, typ := range typeSlice {
		ret[typ] = valueSlice[i]
	}
	return ret
}

// ServiceExternalID is the external id of the service (used in other systems)
func (l ImportRecord) ServiceExternalID() string {
	return l.line[l.title.ServiceExternalIDIdx()]
}

// CustomerName is name of customer that uses the services
func (l ImportRecord) CustomerName() string {
	return l.line[l.title.CustomerNameIdx()]
}

// CustomerExternalID is the external id of customer that uses the services
func (l ImportRecord) CustomerExternalID() string {
	return l.line[l.title.CustomerExternalIDIdx()]
}

// Status is the status of the service (can be of types enum ServiceType in graphql)
func (l ImportRecord) Status() string {
	return l.line[l.title.StatusIdx()]
}

// PortData returns the relevant info for the port from the CSV
func (l ImportRecord) PortData(side *string) (*PortData, error) {
	if l.entity() == ImportEntityPort {
		return &PortData{
			ID:                l.ID(),
			Name:              l.Name(),
			TypeName:          l.TypeName(),
			EquipmentName:     l.PortEquipmentName(),
			EquipmentTypeName: l.PortEquipmentTypeName(),
		}, nil
	}
	return nil, errors.New("unsupported entity for link port Data")
}

// ConsumerPortsServices is the list of services where the port is their consumer endpoint
func (l ImportRecord) ConsumerPortsServices() string {
	return l.line[l.title.ConsumerPortsServicesIdx()]
}

// ProviderPortsServices is the list of services where the port is their provider endpoint
func (l ImportRecord) ProviderPortsServices() string {
	return l.line[l.title.ProviderPortsServicesIdx()]
}

func (l ImportRecord) ServiceNames() string {
	return l.line[l.title.ServiceNamesIdx()]
}

func (l ImportRecord) LinkGetTwoPortsSlices() [][]string {
	if l.entity() == ImportEntityLink {
		idxA, idxB := l.Header().LinkGetTwoPortsRange()
		return [][]string{l.line[idxA[0]:idxA[1]], l.line[idxB[0]:idxB[1]]}
	}
	return nil
}
