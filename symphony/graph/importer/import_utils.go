// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"context"
	"strconv"

	"github.com/AlekSi/pointer"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/customer"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentposition"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func pointerToServiceStatus(status models.ServiceStatus) *models.ServiceStatus {
	return &status
}

func (m *importer) getOrCreateEquipmentType(ctx context.Context, name string, positionsCount int, positionPrefix string, portsCount int, props []*models.PropertyTypeInput) *ent.EquipmentType {
	log := m.log.For(ctx)
	client := m.ClientFrom(ctx)

	equipmentType, err := client.EquipmentType.Query().Where(equipmenttype.Name(name)).Only(ctx)
	if equipmentType != nil {
		return equipmentType
	}
	if !ent.IsNotFound(err) {
		panic(err)
	}
	var proprArr []*ent.PropertyType
	for _, input := range props {
		propEnt := client.PropertyType.
			Create().
			SetName(input.Name).
			SetType(input.Type.String()).
			SetNillableStringVal(input.StringValue).
			SetNillableIntVal(input.IntValue).
			SetNillableBoolVal(input.BooleanValue).
			SetNillableFloatVal(input.FloatValue).
			SetNillableLatitudeVal(input.LatitudeValue).
			SetNillableLongitudeVal(input.LongitudeValue).
			SetNillableIsInstanceProperty(input.IsInstanceProperty).
			SetNillableEditable(input.IsEditable).
			SaveX(ctx)
		proprArr = append(proprArr, propEnt)
	}
	wq := client.EquipmentType.Create().
		AddPropertyTypes(proprArr...).
		SetName(name)
	for i := 1; i <= positionsCount; i++ {
		p := client.EquipmentPositionDefinition.Create().
			SetName(positionPrefix + strconv.Itoa(i)).
			SaveX(ctx)
		wq.AddPositionDefinitions(p)
	}
	for i := 1; i <= portsCount; i++ {
		p := client.EquipmentPortDefinition.Create().
			SetName(strconv.Itoa(i)).
			SaveX(ctx)
		wq.AddPortDefinitions(p)
	}
	log.Debug("Creating new equipment type", zap.String("name", name))
	return wq.SaveX(ctx)
}

func (m *importer) queryLocationForTypeAndParent(ctx context.Context, name string, locType *ent.LocationType, parentID *string) (*ent.Location, error) {
	rq := locType.QueryLocations().Where(location.Name(name))
	if parentID != nil {
		rq = rq.Where(location.HasParentWith(location.ID(*parentID)))
	}
	l, err := rq.Only(ctx)
	if l != nil {
		return l, nil
	}
	return nil, err
}

func (m *importer) getOrCreateLocation(ctx context.Context, name string, latitude float64, longitude float64, locType *ent.LocationType, parentID *string, props []*models.PropertyInput, externalID *string) (*ent.Location, bool) {
	log := m.log.For(ctx)
	l, err := m.queryLocationForTypeAndParent(ctx, name, locType, parentID)
	if l != nil {
		return l, false
	}
	if !ent.IsNotFound(err) {
		log.Panic("query location failed", zap.String("name", name), zap.Error(err))
	}
	log.Debug("Creating new location", zap.String("name", name))
	l, err = m.r.Mutation().AddLocation(ctx, models.AddLocationInput{
		Name:       name,
		Type:       locType.ID,
		Parent:     parentID,
		Latitude:   &latitude,
		Longitude:  &longitude,
		Properties: props,
		ExternalID: externalID,
	})
	if err != nil {
		log.Panic(err.Error(), zap.Error(err))
	}
	return l, true
}

func (m *importer) getEquipmentIfExist(ctx context.Context, mr generated.MutationResolver, name string, equipType *ent.EquipmentType, externalID *string, loc *ent.Location, position *ent.EquipmentPosition, props []*models.PropertyInput) (*ent.Equipment, error) {
	log := m.log.For(ctx)
	client := m.ClientFrom(ctx)
	rq := client.EquipmentType.Query().
		Where(equipmenttype.ID(equipType.ID)).
		QueryEquipment().
		Where(
			equipment.Name(name),
		)
	if loc != nil {
		rq = rq.Where(equipment.HasLocationWith(location.ID(loc.ID)))
	}
	if position != nil {
		rq = rq.Where(
			equipment.HasParentPositionWith(equipmentposition.ID(position.ID)),
		)
	}
	equip, err := rq.First(ctx)
	if ent.MaskNotFound(err) != nil {
		return nil, err
	}
	if equip != nil {
		log.Debug("equipment exists",
			zap.String("name", name),
			zap.String("type", equipType.ID),
		)
		return equip, nil
	}
	return nil, nil
}

func (m *importer) getOrCreateEquipment(ctx context.Context, mr generated.MutationResolver, name string, equipType *ent.EquipmentType, externalID *string, loc *ent.Location, position *ent.EquipmentPosition, props []*models.PropertyInput) (*ent.Equipment, bool, error) {
	log := m.log.For(ctx)
	eq, err := m.getEquipmentIfExist(ctx, mr, name, equipType, externalID, loc, position, props)
	if err != nil || eq != nil {
		return eq, false, err
	}

	var locID *string
	if loc != nil {
		locID = &loc.ID
	}

	var parentEquipmentID, positionDefinitionID *string
	if position != nil {
		p := position.QueryParent().OnlyXID(ctx)
		d := position.QueryDefinition().OnlyXID(ctx)
		parentEquipmentID = &p
		positionDefinitionID = &d
	}
	equip, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:               name,
		Type:               equipType.ID,
		Location:           locID,
		Parent:             parentEquipmentID,
		PositionDefinition: positionDefinitionID,
		Properties:         props,
		ExternalID:         externalID,
	})
	if err != nil {
		log.Error("add equipment", zap.String("name", name), zap.Error(err))
		return nil, false, err
	}
	log.Debug("Creating new equipment", zap.String("equip.Name", equip.Name), zap.String("equip.ID", equip.ID))
	return equip, true, nil
}

func (m *importer) getServiceIfExist(ctx context.Context, mr generated.MutationResolver, name string, serviceType *ent.ServiceType, props []*models.PropertyInput, customerID *string, externalID *string, status models.ServiceStatus) (*ent.Service, error) {
	log := m.log.For(ctx)
	client := m.ClientFrom(ctx)
	rq := client.ServiceType.Query().
		Where(servicetype.ID(serviceType.ID)).
		QueryServices().
		Where(
			service.Name(name),
		)
	service, err := rq.First(ctx)
	if ent.MaskNotFound(err) != nil {
		return nil, err
	}
	if service != nil {
		log.Debug("service exists",
			zap.String("name", name),
			zap.String("type", serviceType.ID),
		)
		return service, nil
	}
	return nil, nil
}

func (m *importer) getOrCreateService(
	ctx context.Context, mr generated.MutationResolver, name string, serviceType *ent.ServiceType, props []*models.PropertyInput, customerID *string, externalID *string, status models.ServiceStatus) (*ent.Service, bool, error) {
	log := m.log.For(ctx)
	service, err := m.getServiceIfExist(ctx, mr, name, serviceType, props, customerID, externalID, status)

	if err != nil || service != nil {
		return service, false, err
	}

	service, err = mr.AddService(ctx, models.ServiceCreateData{
		Name:          name,
		ServiceTypeID: serviceType.ID,
		Properties:    props,
		Status:        pointerToServiceStatus(status),
		CustomerID:    customerID,
		ExternalID:    externalID,
	})
	if err != nil {
		log.Error("add service", zap.String("name", name), zap.Error(err))
		return nil, false, err
	}
	log.Debug("Creating new service", zap.String("service.Name", service.Name), zap.String("service.ID", service.ID))

	return service, true, nil
}

func (m *importer) getCustomerIfExist(ctx context.Context, name string) (*ent.Customer, error) {
	log := m.log.For(ctx)
	client := m.ClientFrom(ctx)
	customer, err := client.Customer.Query().Where(customer.Name(name)).First(ctx)
	if customer != nil {
		log.Debug("customer exists",
			zap.String("name", name),
		)
		return customer, nil
	}
	if !ent.IsNotFound(err) {
		return nil, err
	}
	return nil, nil
}

func (m *importer) getOrCreateCustomer(ctx context.Context, mr generated.MutationResolver, name string, externalID string) (*ent.Customer, error) {
	log := m.log.For(ctx)
	_, err := m.getCustomerIfExist(ctx, name)
	if err != nil {
		return nil, err
	}

	exID := pointer.ToStringOrNil(externalID)
	customer, err := mr.AddCustomer(ctx, models.AddCustomerInput{
		Name:       name,
		ExternalID: exID,
	})
	if err != nil {
		return nil, err
	}
	log.Debug("Creating new customer", zap.String("customer.Name", customer.Name),
		zap.String("customer.ID", customer.ID))

	return customer, nil
}

func (m *importer) getOrCreateEquipmentLocationByFullPath(ctx context.Context, line, firstLine []string, includePropTypes bool) (string, error) {
	var (
		lastLocationTypeIdx   = getLowestLocationHierarchyIdxForRow(ctx, line)
		indexToLocationTypeID = getImportContext(ctx).indexToLocationTypeID
		resLocation, parent   *ent.Location
	)
	for i, name := range line {
		if i > lastLocationTypeIdx {
			break
		}
		if name == "" {
			continue
		}
		q := m.ClientFrom(ctx).LocationType.Query().
			QueryLocations().
			Where(location.Name(name))
		if parent != nil {
			q = q.Where(location.HasParentWith(location.ID(parent.ID)))
		}
		resLocation = q.FirstX(ctx)
		if resLocation == nil {
			m.log.For(ctx).Debug("didn't find parent- creating a new location", zap.String("name", name))
			locationTypeID := indexToLocationTypeID[i]
			var pinputs []*models.PropertyInput

			if i == lastLocationTypeIdx && includePropTypes {
				locType := m.ClientFrom(ctx).LocationType.Query().Where(locationtype.ID(locationTypeID)).OnlyX(ctx)
				propTypes := locType.QueryPropertyTypes().AllX(ctx)
				for _, ptype := range propTypes {
					index := findIndex(firstLine, ptype.Name)
					if index != -1 {
						pinputs = append(pinputs, &models.PropertyInput{
							PropertyTypeID: ptype.ID,
							StringValue:    &line[index],
						})
					}
				}
			}

			l, err := m.r.Mutation().AddLocation(ctx, models.AddLocationInput{
				Name:       name,
				Type:       locationTypeID,
				Properties: pinputs,
				Parent: func() *string {
					if parent != nil {
						return &parent.ID
					}
					return nil
				}(),
			})
			if err != nil {
				return "", errors.WithMessage(err, "cannot add location")
			}
			resLocation = l
		}
		parent = resLocation
	}
	if resLocation != nil {
		return resLocation.ID, nil
	}
	return "", nil
}

func (m *importer) getLocationIDByName(ctx context.Context, loc string) (string, error) {
	return m.ClientFrom(ctx).Location.Query().
		Where(location.Name(loc)).
		OnlyID(ctx)
}

func (m *importer) CloneContext(ctx context.Context) context.Context {
	return viewer.NewContext(ent.NewContext(context.Background(), m.ClientFrom(ctx)), viewer.FromContext(ctx))
}

func (m *importer) validateServiceExistsAndUnique(ctx context.Context, serviceNamesMap map[string]bool, serviceName string) (string, error) {
	client := m.ClientFrom(ctx)
	if _, ok := serviceNamesMap[serviceName]; ok {
		return "", errors.Errorf("Property can't be the endpoint of the same service more than once - service name=%q", serviceName)
	}
	serviceNamesMap[serviceName] = true
	s, err := client.Service.Query().Where(service.Name(serviceName)).Only(ctx)
	if err != nil {
		return "", errors.Wrapf(err, "can't query service name=%q", serviceName)
	}
	return s.ID, nil
}
