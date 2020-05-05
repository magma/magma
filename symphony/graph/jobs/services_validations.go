// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jobs

import (
	"context"
	"strconv"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpointdefinition"
	"github.com/pkg/errors"
)

func (m *jobs) validateEndpointsExistAndLinked(ctx context.Context, srvc *ent.Service) error {
	log := m.logger.For(ctx)
	sc := getServicesContext(ctx)

	orderedDefQuery := srvc.QueryType().
		QueryEndpointDefinitions().
		Order(ent.Asc(serviceendpointdefinition.FieldIndex))
	endpoints, err := orderedDefQuery.
		Clone().
		QueryEndpoints().
		Where(serviceendpoint.HasServiceWith(service.ID(srvc.ID))).
		All(ctx)
	if err != nil {
		return errors.Wrap(err, "query service endpoints")
	}
	definitions, err := orderedDefQuery.All(ctx)
	if err != nil {
		return errors.Wrap(err, "query service endpoint definitions")
	}
	var prev *ent.Equipment

	for i, ep := range endpoints {
		// Verify all endpoint has an equipment
		curr, err := ep.QueryEquipment().Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				log.Debug("deleting Service with missing equipment" + strconv.Itoa(srvc.ID))
				err = deleteService(ctx, srvc)
				if err != nil {
					return errors.Wrap(err, "deleting service")
				}
				sc.deleted++
				return nil
			}
			return errors.Wrap(err, "fetching equipment from endpoint")
		}
		// Verify equipment matches equipmentType on ServiceEndpointDefinition
		eTypeID, err := curr.QueryType().OnlyID(ctx)
		if err != nil {
			return errors.Wrapf(err, "fetching equipment type from equipment %v", curr.ID)
		}

		eTypeFromEndpointDef, err := definitions[i].QueryEquipmentType().Only(ctx)
		if err != nil {
			return errors.Wrapf(err, "fetching equipment type from endpoint definition  %v", definitions[i].ID)
		}
		if eTypeFromEndpointDef.ID != eTypeID {
			log.Debug("deleting Service with mismatched endpoint & equipment: " + strconv.Itoa(srvc.ID))
			err = deleteService(ctx, srvc)
			if err != nil {
				return errors.Wrap(err, "deleting service")
			}
			sc.deleted++
			return nil
		}

		// Verify all equipment are linked by the index order
		if prev != nil {
			linked, err := prev.QueryPorts().Where(
				equipmentport.HasLinkWith(link.HasPortsWith(
					equipmentport.HasParentWith(
						equipment.ID(curr.ID),
						equipment.Not(equipment.ID(prev.ID)),
					)))).Exist(ctx)

			if ent.MaskNotFound(err) != nil {
				return errors.Wrap(err, "checking equipment are linked")
			}
			if !linked {
				log.Debug("deleting Service with unlinked equipment" + strconv.Itoa(srvc.ID))
				err = deleteService(ctx, srvc)
				if err != nil {
					return errors.Wrap(err, "deleting service")
				}
				sc.deleted++
				return nil
			}
		}
		prev = curr
	}
	return nil
}
