// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AlekSi/pointer"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/resolverutil"
	"github.com/pkg/errors"
)

type equipmentPortResolver struct{}

func (r equipmentPortResolver) ServiceEndpoints(ctx context.Context, obj *ent.EquipmentPort) ([]*ent.ServiceEndpoint, error) {
	return obj.QueryEndpoints().All(ctx)
}

func (equipmentPortResolver) Definition(ctx context.Context, obj *ent.EquipmentPort) (*ent.EquipmentPortDefinition, error) {
	return obj.QueryDefinition().Only(ctx)
}

func (equipmentPortResolver) ParentEquipment(ctx context.Context, obj *ent.EquipmentPort) (*ent.Equipment, error) {
	return obj.QueryParent().Only(ctx)
}

func (equipmentPortResolver) Link(ctx context.Context, obj *ent.EquipmentPort) (*ent.Link, error) {
	l, err := obj.QueryLink().Only(ctx)
	return l, ent.MaskNotFound(err)
}

func (equipmentPortResolver) Properties(ctx context.Context, obj *ent.EquipmentPort) ([]*ent.Property, error) {
	return obj.QueryProperties().All(ctx)
}

type equipmentPositionResolver struct{}

func (equipmentPositionResolver) Definition(ctx context.Context, obj *ent.EquipmentPosition) (*ent.EquipmentPositionDefinition, error) {
	return obj.QueryDefinition().Only(ctx)
}

func (equipmentPositionResolver) ParentEquipment(ctx context.Context, obj *ent.EquipmentPosition) (*ent.Equipment, error) {
	return obj.QueryParent().Only(ctx)
}

func (equipmentPositionResolver) AttachedEquipment(ctx context.Context, obj *ent.EquipmentPosition) (*ent.Equipment, error) {
	e, err := obj.QueryAttachment().Only(ctx)
	return e, ent.MaskNotFound(err)
}

type equipmentPortDefinitionResolver struct{}

func (equipmentPortDefinitionResolver) PortType(ctx context.Context, obj *ent.EquipmentPortDefinition) (*ent.EquipmentPortType, error) {
	l, err := obj.QueryEquipmentPortType().Only(ctx)
	return l, ent.MaskNotFound(err)
}

type equipmentPortTypeResolver struct{}

func (equipmentPortTypeResolver) PropertyTypes(ctx context.Context, obj *ent.EquipmentPortType) ([]*ent.PropertyType, error) {
	return obj.QueryPropertyTypes().All(ctx)
}

func (equipmentPortTypeResolver) LinkPropertyTypes(ctx context.Context, obj *ent.EquipmentPortType) ([]*ent.PropertyType, error) {
	return obj.QueryLinkPropertyTypes().All(ctx)
}

func (equipmentPortTypeResolver) NumberOfPortDefinitions(ctx context.Context, obj *ent.EquipmentPortType) (int, error) {
	return obj.QueryPortDefinitions().Count(ctx)
}

type equipmentTypeResolver struct{}

func (equipmentTypeResolver) Category(ctx context.Context, obj *ent.EquipmentType) (*string, error) {
	c, err := obj.QueryCategory().Only(ctx)
	if c != nil {
		return &c.Name, err
	}
	return nil, ent.MaskNotFound(err)
}

func (equipmentTypeResolver) PositionDefinitions(ctx context.Context, obj *ent.EquipmentType) ([]*ent.EquipmentPositionDefinition, error) {
	return obj.QueryPositionDefinitions().All(ctx)
}

func (equipmentTypeResolver) PortDefinitions(ctx context.Context, obj *ent.EquipmentType) ([]*ent.EquipmentPortDefinition, error) {
	return obj.QueryPortDefinitions().All(ctx)
}

func (equipmentTypeResolver) PropertyTypes(ctx context.Context, obj *ent.EquipmentType) ([]*ent.PropertyType, error) {
	return obj.QueryPropertyTypes().All(ctx)
}

func (equipmentTypeResolver) Equipments(ctx context.Context, obj *ent.EquipmentType) ([]*ent.Equipment, error) {
	return obj.QueryEquipment().All(ctx)
}

func (equipmentTypeResolver) NumberOfEquipment(ctx context.Context, obj *ent.EquipmentType) (int, error) {
	return obj.QueryEquipment().Count(ctx)
}

type equipmentResolver struct{ resolver }

func (r equipmentResolver) DescendentsIncludingSelf(ctx context.Context, obj *ent.Equipment) ([]*ent.Equipment, error) {
	equip := *obj
	var err error
	var ret []*ent.Equipment
	children, err := equip.QueryPositions().QueryAttachment().All(ctx)
	if err == nil {
		for _, child := range children {
			grandChildren, err := r.DescendentsIncludingSelf(ctx, child)
			if err == nil {
				ret = append(ret, grandChildren...)
			}
		}
	}
	ret = append(ret, obj)
	return ret, err
}

func (equipmentResolver) ParentLocation(ctx context.Context, obj *ent.Equipment) (*ent.Location, error) {
	l, err := obj.QueryLocation().Only(ctx)
	return l, ent.MaskNotFound(err)
}

func (equipmentResolver) ParentPosition(ctx context.Context, obj *ent.Equipment) (*ent.EquipmentPosition, error) {
	p, err := obj.QueryParentPosition().Only(ctx)
	return p, ent.MaskNotFound(err)
}

func (equipmentResolver) EquipmentType(ctx context.Context, obj *ent.Equipment) (*ent.EquipmentType, error) {
	return obj.QueryType().Only(ctx)
}

func (equipmentResolver) Positions(ctx context.Context, obj *ent.Equipment) ([]*ent.EquipmentPosition, error) {
	return obj.QueryPositions().All(ctx)
}

func (equipmentResolver) Ports(ctx context.Context, obj *ent.Equipment, availableOnly *bool) ([]*ent.EquipmentPort, error) {
	q := obj.QueryPorts()
	if pointer.GetBool(availableOnly) {
		q.Where(equipmentport.Not(equipmentport.HasLink()))
	}
	return q.All(ctx)
}

func (equipmentResolver) Properties(ctx context.Context, obj *ent.Equipment) ([]*ent.Property, error) {
	return obj.QueryProperties().All(ctx)
}

func (equipmentResolver) FutureState(_ context.Context, obj *ent.Equipment) (*models.FutureState, error) {
	fs := models.FutureState(obj.FutureState)
	return &fs, nil
}

func (equipmentResolver) WorkOrder(ctx context.Context, obj *ent.Equipment) (*ent.WorkOrder, error) {
	wo, err := obj.QueryWorkOrder().Only(ctx)
	return wo, ent.MaskNotFound(err)
}

func (equipmentResolver) Images(ctx context.Context, obj *ent.Equipment) ([]*ent.File, error) {
	return obj.QueryFiles().Where(file.Type(models.FileTypeImage.String())).All(ctx)
}

func (equipmentResolver) Files(ctx context.Context, obj *ent.Equipment) ([]*ent.File, error) {
	return obj.QueryFiles().Where(file.Type(models.FileTypeFile.String())).All(ctx)
}

func (equipmentResolver) Hyperlinks(ctx context.Context, obj *ent.Equipment) ([]*ent.Hyperlink, error) {
	return obj.QueryHyperlinks().All(ctx)
}

func (equipmentResolver) PositionHierarchy(ctx context.Context, eq *ent.Equipment) ([]*ent.EquipmentPosition, error) {
	var positions []*ent.EquipmentPosition
	ppos, err := eq.QueryParentPosition().Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, errors.Wrap(err, "querying parent position")
	}
	for ppos != nil {
		positions = append([]*ent.EquipmentPosition{ppos}, positions...)
		p, err := ppos.QueryParent().QueryParentPosition().Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return nil, errors.Wrap(err, "querying parent position")
		}

		ppos = p
	}
	return positions, nil
}

func (r equipmentResolver) LocationHierarchy(ctx context.Context, eq *ent.Equipment) ([]*ent.Location, error) {
	ph, err := r.PositionHierarchy(ctx, eq)
	if err != nil {
		return nil, err
	}
	var locs []*ent.Location
	var pl *ent.Location
	if len(ph) > 0 {
		pl = ph[0].QueryParent().QueryLocation().OnlyX(ctx)
	} else {
		pl = eq.QueryLocation().OnlyX(ctx)
	}
	for pl != nil {
		locs = append([]*ent.Location{pl}, locs...)
		l, err := pl.QueryParent().Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return nil, err
		}
		pl = l
	}
	return locs, nil
}

func (r equipmentResolver) Device(ctx context.Context, eq *ent.Equipment) (*models.Device, error) {
	if eq.DeviceID == "" {
		return nil, nil
	}
	if r.orc8rClient == nil {
		return nil, errors.New("unsupported orc8r field")
	}
	parts := strings.Split(eq.DeviceID, ".")
	if len(parts) < 2 {
		return nil, errors.New("invalid equipment device id")
	}

	uri := fmt.Sprintf("/magma/v1/networks/%s/gateways/%s/status", parts[1], parts[0])
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create http request: %w", err)
	}
	rsp, err := r.orc8rClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("performing orc8r status request: %w", err)
	}
	defer rsp.Body.Close()

	var result struct {
		CheckinTime int64 `json:"checkin_time"`
	}
	if err := json.NewDecoder(rsp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding orc8r response: %w", err)
	}

	dev := models.Device{ID: eq.DeviceID}
	if result.CheckinTime != 0 {
		up := checkinTimeIsUp(result.CheckinTime, 3*time.Minute)
		dev.Up = &up
	}
	return &dev, nil
}

// checkinTime is milliseconds since epoch
func checkinTimeIsUp(checkinTime int64, buffer time.Duration) bool {
	return checkinTime/1000+int64(buffer.Seconds()) > time.Now().Unix()
}

func (r equipmentResolver) Services(ctx context.Context, obj *ent.Equipment) ([]*ent.Service, error) {
	eqPred := resolverutil.BuildGeneralEquipmentAncestorFilter(equipment.ID(obj.ID), 1, 4)
	endpoints, err := r.ClientFrom(ctx).ServiceEndpoint.Query().Where(
		serviceendpoint.HasPortWith(equipmentport.HasParentWith(eqPred))).All(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "querying service endpoints")
	}
	eids := make([]string, len(endpoints))
	for _, ep := range endpoints {
		eids = append(eids, ep.ID)
	}

	services, err := r.ClientFrom(ctx).Service.Query().Where(service.HasEndpointsWith(serviceendpoint.IDIn(eids...))).All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "querying services where equipment port is an endpoint")
	}

	ids := make([]string, len(services))
	for i, svc := range services {
		ids[i] = svc.ID
	}

	linkServices, err := r.ClientFrom(ctx).Service.Query().Where(
		service.HasLinksWith(link.HasPortsWith(equipmentport.HasParentWith(
			resolverutil.BuildGeneralEquipmentAncestorFilter(equipment.ID(obj.ID), 1, 3)))),
		service.Not(service.IDIn(ids...))).All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "querying services where equipment connected to link of service")
	}

	services = append(services, linkServices...)

	return services, nil
}
