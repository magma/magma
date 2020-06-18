// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jobs

import (
	"context"
	"math/rand"
	"strconv"

	"github.com/prometheus/common/log"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/ent/equipment"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentport"
	"github.com/facebookincubator/symphony/pkg/ent/equipmenttype"
	"github.com/facebookincubator/symphony/pkg/ent/link"
	"github.com/pkg/errors"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/service"
	"github.com/facebookincubator/symphony/pkg/ent/serviceendpoint"
	"github.com/facebookincubator/symphony/pkg/ent/serviceendpointdefinition"
	"github.com/facebookincubator/symphony/pkg/ent/servicetype"
)

func getStructuredCurrentServicesForType(ctx context.Context, sType *ent.ServiceType) ([]serviceEquipmentListData, error) {
	client := ent.FromContext(ctx)
	var detailedServices []serviceEquipmentListData
	services, err := client.Service.Query().
		Where(
			service.HasEndpoints(),
			service.HasTypeWith(servicetype.ID(sType.ID)),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}
	for _, srvc := range services {
		endpoints, err := srvc.QueryType().
			QueryEndpointDefinitions().
			Order(ent.Asc(serviceendpointdefinition.FieldIndex)).
			QueryEndpoints().
			Where(serviceendpoint.HasServiceWith(service.ID(srvc.ID))).
			All(ctx)
		if err != nil {
			return nil, err
		}
		equipList := make([]*ent.Equipment, len(endpoints))
		for i, ep := range endpoints {
			equip, err := ep.QueryEquipment().Only(ctx)
			if err != nil {
				return nil, err
			}
			equipList[i] = equip
		}

		detailedServices = append(detailedServices, serviceEquipmentListData{
			EquipmentList: equipList,
		})
	}
	return detailedServices, nil
}

func (m *jobs) getServicesDetailsList(ctx context.Context, sType *ent.ServiceType) ([]serviceEquipmentListData, error) {
	var serviceDataList []serviceEquipmentListData
	endpointDefs, err := sType.QueryEndpointDefinitions().
		Order(ent.Asc(serviceendpointdefinition.FieldIndex)).All(ctx)
	if err != nil {
		log.Warn("[SKIP] can't get endpoints definitions for service type" + sType.Name + ". error: " + err.Error())
		return nil, err
	}

	if len(endpointDefs) < 2 || len(endpointDefs) > MaxEndpoints {
		log.Info("[SKIPPING SERVICE TYPE] either too many or not enough endpoint types " + sType.Name)
		return nil, err
	}

	e1s, err := endpointDefs[0].QueryEquipmentType().QueryEquipment().All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "can't get equipment list for first endpoint type")
	}
	for _, e1 := range e1s {
		e2s, err := getNextEquipmentInstances(ctx, e1, endpointDefs[1], []int{})

		if err != nil {
			return nil, err
		}
		for _, e2 := range e2s {
			if len(endpointDefs) == 2 {
				log.Info("adding service to 'toAdd' list: ", sType.Name, "equipment:", e1.ID, e2.ID)
				serviceDataList = append(serviceDataList, serviceEquipmentListData{
					[]*ent.Equipment{e1, e2},
				})
				continue
			}
			e3s, err := getNextEquipmentInstances(ctx, e2, endpointDefs[2], []int{e1.ID})
			if err != nil {
				return nil, err
			}
			for _, e3 := range e3s {
				if len(endpointDefs) == 3 {
					log.Info("adding service to 'toAdd' list: ", sType.Name, "equipment:", e1.ID, e2.ID, e3.ID)
					serviceDataList = append(serviceDataList, serviceEquipmentListData{
						[]*ent.Equipment{e1, e2, e3},
					})
					continue
				}
				e4s, err := getNextEquipmentInstances(ctx, e3, endpointDefs[3], []int{e1.ID, e2.ID})
				if err != nil {
					return nil, err
				}
				for _, e4 := range e4s {
					if len(endpointDefs) == 4 {
						log.Info("adding service to 'toAdd' list: ", sType.Name, "equipment:", e1.ID, e2.ID, e3.ID, e4.ID)
						serviceDataList = append(serviceDataList, serviceEquipmentListData{
							[]*ent.Equipment{e1, e2, e3, e4},
						})
						continue
					}
					e5s, err := getNextEquipmentInstances(ctx, e4, endpointDefs[4], []int{e1.ID, e2.ID, e3.ID})
					if err != nil {
						return nil, err
					}
					for _, e5 := range e5s {
						if len(endpointDefs) != MaxEndpoints {
							return nil, errors.Errorf("service types support up to 5 endpoint definitions")
						}
						log.Info("adding service to 'toAdd' list: ", sType.Name, "equipment:", e1.ID, e2.ID, e3.ID, e4.ID, e5.ID)
						serviceDataList = append(serviceDataList, serviceEquipmentListData{
							[]*ent.Equipment{e1, e2, e3, e4, e5},
						})
						continue
					}
				}
			}
		}
	}
	return serviceDataList, nil
}

func getNextEquipmentInstances(ctx context.Context, e *ent.Equipment, ed *ent.ServiceEndpointDefinition, prevEquipmentIDs []int) ([]*ent.Equipment, error) {
	typ2, err := ed.QueryEquipmentType().Only(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "can't get equipment type")
	}
	nextEquipmentInstancesQuery := e.QueryPorts().
		QueryLink().
		QueryPorts().
		QueryParent().
		Where(equipment.And(
			equipment.HasTypeWith(equipmenttype.ID(typ2.ID)),
			equipment.Not(equipment.ID(e.ID)),
		))
	if len(prevEquipmentIDs) > 0 {
		nextEquipmentInstancesQuery =
			nextEquipmentInstancesQuery.Where(equipment.IDNotIn(prevEquipmentIDs...))
	}
	nextEquipmentInstances, err := nextEquipmentInstancesQuery.All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "can't query equipment on link")
	}
	return nextEquipmentInstances, nil
}

// nolint: gosec
func (m *jobs) createServicesFromList(ctx context.Context, serviceDetails []serviceEquipmentListData, serviceType *ent.ServiceType) error {
	client := ent.FromContext(ctx)
	existingServicesStructuredList, err := getStructuredCurrentServicesForType(ctx, serviceType)
	if err != nil {
		return errors.Wrap(err, "error fetching current services")
	}

	for _, serviceData := range serviceDetails {
		// TODO search can be optimized
		if isServiceEquipmentListAlreadyExists(existingServicesStructuredList, serviceData) {
			log.Info("skipping new service: already exists. type:", serviceType.Name, " first equipment:", serviceData.EquipmentList[0].Name, serviceData.EquipmentList[0].ID)
			continue
		}
		log.Info("saving new service. type:", serviceType.Name, " first equipment:", serviceData.EquipmentList[0].Name, serviceData.EquipmentList[0].ID)

		links, err := getLinksFromEquipmentList(ctx, serviceData.EquipmentList)
		if err != nil || len(links) == 0 {
			return errors.Wrapf(err, "can't get links for service")
		}
		srvc, err := client.Service.Create().
			SetStatus(models.ServiceStatusPending.String()).
			AddLinks(links...).
			SetName(strconv.Itoa(rand.Int())[:5]).
			SetType(serviceType).
			Save(ctx)
		if err != nil {
			return err
		}

		name, err := m.generateName(ctx, serviceData, srvc.ID)
		if err != nil {
			return err
		}
		srvc, err = client.Service.UpdateOneID(srvc.ID).SetName(*name).Save(ctx)
		if err != nil {
			return err
		}
		endpointDefs, err := serviceType.
			QueryEndpointDefinitions().
			Order(ent.Asc(serviceendpointdefinition.FieldIndex)).
			All(ctx)
		if err != nil {
			return err
		}
		for i, equip := range serviceData.EquipmentList {
			_, err = client.ServiceEndpoint.Create().
				SetEquipment(equip).
				SetDefinition(endpointDefs[i]).
				SetService(srvc).
				Save(ctx)
			if err != nil {
				return errors.Wrapf(err, "saving service endpoint on service %v", srvc.Name)
			}
		}
	}
	return nil
}

func isServiceEquipmentListAlreadyExists(all []serviceEquipmentListData, serviceToAdd serviceEquipmentListData) bool {
	for _, curr := range all {
		if isSameEquipmentList(serviceToAdd, curr) {
			return true
		}
	}
	return false
}

func isSameEquipmentList(curr serviceEquipmentListData, serviceToAdd serviceEquipmentListData) bool {
	allEquips := false
	if len(curr.EquipmentList) == len(serviceToAdd.EquipmentList) {
		allEquips = true
		for i, equip := range serviceToAdd.EquipmentList {
			if curr.EquipmentList[i].ID != equip.ID {
				return false
			}
		}
		if allEquips {
			return true
		}
	}
	return false
}

func getLinksFromEquipmentList(ctx context.Context, equipmentList []*ent.Equipment) ([]*ent.Link, error) {
	var (
		links []*ent.Link
		prev  *ent.Equipment
	)
	for _, curr := range equipmentList {
		if prev == nil {
			prev = curr
			continue
		}
		l, err := prev.QueryPorts().
			QueryLink().
			Where(link.HasPortsWith(equipmentport.HasParentWith(equipment.ID(curr.ID)))).
			First(ctx)
		if err != nil {
			return nil, err
		}
		links = append(links, l)
		prev = curr
	}
	return links, nil
}
