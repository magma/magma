// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jobs

import (
	"context"
	"net/http"
	"strconv"

	"github.com/AlekSi/pointer"
	"github.com/pkg/errors"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/service"
	"github.com/facebookincubator/symphony/pkg/ent/servicetype"
	"github.com/facebookincubator/symphony/pkg/viewer"
)

type serviceEquipmentListData struct {
	EquipmentList []*ent.Equipment
}

const MaxEndpoints = 5

// syncServices job syncs the services according to changes
func (m *jobs) syncServices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	v := viewer.FromContext(ctx)
	sc := getServicesContext(ctx)
	log := m.logger.For(ctx)
	client := ent.FromContext(ctx)
	log.Info("services sync run. tenant: " + v.Tenant())

	services, err := client.Service.Query().Where(
		service.HasTypeWith(servicetype.DiscoveryMethodEQ(servicetype.DiscoveryMethodINVENTORY))).
		All(ctx)

	if err != nil {
		errorReturn(w, "can't get services", log, err)
		return
	}

	log.Info("service sync - looking for outdated services to delete")
	for _, srvc := range services {
		typ, err := srvc.QueryType().Only(ctx)
		if err != nil {
			log.Warn("[SKIP] can't get service type" + srvc.Name + ". error: " + err.Error())
			continue
		}
		if typ.IsDeleted {
			log.Debug("Deleting 'isDeleted' marked Service" + strconv.Itoa(srvc.ID))
			err = deleteService(ctx, srvc)
			if err != nil {
				log.Warn("[SKIP] can't delete service" + srvc.Name + ". error: " + err.Error())
				continue
			}
			sc.deleted++
		}
		err = m.validateEndpointsExistAndLinked(ctx, srvc)
		if err != nil {
			log.Warn("[SKIP] error while validating existing service" + srvc.Name + ". error: " + err.Error())
			continue
		}
	}
	log.Info("done deleting services, deleted instances: " + strconv.Itoa(sc.deleted))

	log.Info("service Sync - Add new services")
	serviceTypes, err := client.ServiceType.Query().
		Where(
			servicetype.DiscoveryMethodEQ(servicetype.DiscoveryMethodINVENTORY),
			servicetype.IsDeleted(false),
		).
		All(ctx)

	if err != nil {
		errorReturn(w, "can't get service types", log, err)
		return
	}
	for _, sType := range serviceTypes {
		log.Info("going over type: " + sType.Name)
		servicesDataListToAdd, err := m.getServicesDetailsList(ctx, sType)
		if err != nil {
			log.Warn("[SKIP] can't get service details for service type" + sType.Name + ". error: " + err.Error())
			continue
		}
		err = m.createServicesFromList(ctx, servicesDataListToAdd, sType)
		if err != nil {
			log.Warn("[SKIP] can't create services for type" + sType.Name + ". error: " + err.Error())
			continue
		}
	}
}

func (m *jobs) generateName(ctx context.Context, s serviceEquipmentListData, id int) (*string, error) {
	locStart, err := m.r.Equipment().FirstLocation(ctx, s.EquipmentList[0])
	if err != nil {
		return nil, errors.Wrapf(err, "can't query location on equipment %v", s.EquipmentList[0].ID)
	}
	locEnd, err := m.r.Equipment().FirstLocation(ctx, s.EquipmentList[len(s.EquipmentList)-1])
	if err != nil {
		return nil, errors.Wrapf(err, "can't query location on equipment %v", s.EquipmentList[len(s.EquipmentList)-1].ID)
	}

	idAsString := strconv.Itoa(id)
	return pointer.ToString(locStart.Name + "_" + locEnd.Name + "_" + idAsString[len(idAsString)-4:]), nil
}
