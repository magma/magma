// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"context"
	"strconv"

	"github.com/AlekSi/pointer"

	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/customer"
	"github.com/facebookincubator/symphony/pkg/ent/equipment"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentposition"
	"github.com/facebookincubator/symphony/pkg/ent/equipmenttype"
	"github.com/facebookincubator/symphony/pkg/ent/location"
	"github.com/facebookincubator/symphony/pkg/ent/locationtype"
	"github.com/facebookincubator/symphony/pkg/ent/service"
	"github.com/facebookincubator/symphony/pkg/ent/servicetype"
	"github.com/facebookincubator/symphony/pkg/viewer"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func pointerToServiceStatus(status models.ServiceStatus) *models.ServiceStatus {
	return &status
}

func (m *importer) getOrCreateEquipmentType(ctx context.Context, name string, positionsCount int, positionPrefix string, portsCount int, props []*models.PropertyTypeInput) *ent.EquipmentType {
	log := m.logger.For(ctx)
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

func (m *importer) queryLocationForTypeAndParent(ctx context.Context, name string, locType *ent.LocationType, parentID *int) (*ent.Location, error) {
	rq := locType.QueryLocations().Where(location.Name(name))
	if parentID != nil {
		rq = rq.Where(location.HasParentWith(location.ID(*parentID)))
	} else {
		rq = rq.Where(location.Not(location.HasParent()))
	}
	l, err := rq.Only(ctx)
	if l != nil {
		return l, nil
	}
	return nil, err
}

func (m *importer) getOrCreateLocation(
	ctx context.Context, name string, latitude, longitude float64,
	locType *ent.LocationType, parentID *int, props []*models.PropertyInput,
	externalID *string,
) (*ent.Location, bool, error) {
	log := m.logger.For(ctx)
	l, err := m.queryLocationForTypeAndParent(ctx, name, locType, parentID)
	if ent.MaskNotFound(err) != nil {
		return nil, false, err
	}
	if l != nil {
		return l, false, nil
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
		return nil, false, err
	}
	return l, true, nil
}

func (m *importer) getEquipmentIfExist(
	ctx context.Context, name string, equipType *ent.EquipmentType,
	loc *ent.Location, position *ent.EquipmentPosition,
) (*ent.Equipment, error) {
	log := m.logger.For(ctx)
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
			zap.Int("type", equipType.ID),
		)
		return equip, nil
	}
	return nil, nil
}

func (m *importer) getOrCreateEquipment(
	ctx context.Context, mr generated.MutationResolver, name string,
	equipType *ent.EquipmentType, externalID *string, loc *ent.Location,
	position *ent.EquipmentPosition, props []*models.PropertyInput,
) (*ent.Equipment, bool, error) {
	log := m.logger.For(ctx)
	eq, err := m.getEquipmentIfExist(ctx, name, equipType, loc, position)
	if err != nil || eq != nil {
		return eq, false, err
	}

	var locID *int
	if loc != nil {
		locID = &loc.ID
	}

	var parentEquipmentID, positionDefinitionID *int
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
	log.Debug("Creating new equipment",
		zap.String("equip.Name", equip.Name),
		zap.Int("equip.ID", equip.ID),
	)
	return equip, true, nil
}

func (m *importer) getServiceIfExist(ctx context.Context, name string, serviceType *ent.ServiceType) (*ent.Service, error) {
	log := m.logger.For(ctx)
	client := m.ClientFrom(ctx)
	rq := client.ServiceType.Query().
		Where(servicetype.ID(serviceType.ID)).
		QueryServices().
		Where(
			service.Name(name),
		)
	svc, err := rq.First(ctx)
	if ent.MaskNotFound(err) != nil {
		return nil, err
	}
	if svc != nil {
		log.Debug("service exists",
			zap.String("name", name),
			zap.Int("type", serviceType.ID),
		)
		return svc, nil
	}
	return nil, nil
}

func (m *importer) getOrCreateService(
	ctx context.Context, mr generated.MutationResolver, name string,
	serviceType *ent.ServiceType, props []*models.PropertyInput,
	customerID *int, externalID *string, status models.ServiceStatus,
) (*ent.Service, bool, error) {
	log := m.logger.For(ctx)
	svc, err := m.getServiceIfExist(ctx, name, serviceType)
	if err != nil || svc != nil {
		return svc, false, err
	}

	svc, err = mr.AddService(ctx, models.ServiceCreateData{
		Name:          name,
		ServiceTypeID: serviceType.ID,
		Properties:    props,
		Status:        pointerToServiceStatus(status),
		CustomerID:    customerID,
		ExternalID:    externalID,
	})
	if err != nil {
		log.Error("add svc", zap.String("name", name), zap.Error(err))
		return nil, false, err
	}
	log.Debug("Creating new svc", zap.String("svc.Name", svc.Name), zap.Int("svc.ID", svc.ID))

	return svc, true, nil
}

func (m *importer) getCustomerIfExist(ctx context.Context, name string) (*ent.Customer, error) {
	log := m.logger.For(ctx)
	client := m.ClientFrom(ctx)
	c, err := client.Customer.Query().Where(customer.Name(name)).First(ctx)
	if c != nil {
		log.Debug("customer exists",
			zap.String("name", name),
		)
		return c, nil
	}
	if !ent.IsNotFound(err) {
		return nil, err
	}
	return nil, nil
}

func (m *importer) getOrCreateCustomer(ctx context.Context, mr generated.MutationResolver, name string, externalID string) (*ent.Customer, error) {
	log := m.logger.For(ctx)
	_, err := m.getCustomerIfExist(ctx, name)
	if err != nil {
		return nil, err
	}

	exID := pointer.ToStringOrNil(externalID)
	c, err := mr.AddCustomer(ctx, models.AddCustomerInput{
		Name:       name,
		ExternalID: exID,
	})
	if err != nil {
		return nil, err
	}
	log.Debug("Creating new customer",
		zap.String("customer.Name", c.Name),
		zap.Int("customer.ID", c.ID),
	)

	return c, nil
}

func (m *importer) getOrCreateEquipmentLocationByFullPath(ctx context.Context, line, firstLine []string, includePropTypes bool) (int, error) {
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
			m.logger.For(ctx).Debug("didn't find parent- creating a new location", zap.String("name", name))
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
				Parent: func() *int {
					if parent != nil {
						return &parent.ID
					}
					return nil
				}(),
			})
			if err != nil {
				return 0, errors.WithMessage(err, "cannot add location")
			}
			resLocation = l
		}
		parent = resLocation
	}
	if resLocation != nil {
		return resLocation.ID, nil
	}
	return 0, nil
}

func (m *importer) getLocationIDByName(ctx context.Context, name string) (int, error) {
	return m.ClientFrom(ctx).Location.Query().
		Where(location.Name(name)).
		OnlyID(ctx)
}

func (m *importer) CloneContext(ctx context.Context) context.Context {
	return viewer.NewContext(ent.NewContext(context.Background(), m.ClientFrom(ctx)), viewer.FromContext(ctx))
}

func (m *importer) validateServiceExistsAndUnique(ctx context.Context, serviceNamesMap map[string]bool, serviceName string) (int, error) {
	client := m.ClientFrom(ctx)
	if _, ok := serviceNamesMap[serviceName]; ok {
		return 0, errors.Errorf("property can't be the endpoint of the same service more than once - service name=%q", serviceName)
	}
	serviceNamesMap[serviceName] = true
	s, err := client.Service.Query().Where(service.Name(serviceName)).Only(ctx)
	if err != nil {
		return 0, errors.Wrapf(err, "can't query service name=%q", serviceName)
	}
	return s.ID, nil
}
