// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"encoding/json"
	"errors"
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
)

type equipmentPortResolver struct{}

func (equipmentPortResolver) ServiceEndpoints(ctx context.Context, ep *ent.EquipmentPort) ([]*ent.ServiceEndpoint, error) {
	if endpoints, err := ep.Edges.EndpointsOrErr(); !ent.IsNotLoaded(err) {
		return endpoints, err
	}
	return ep.QueryEndpoints().All(ctx)
}

func (equipmentPortResolver) Definition(ctx context.Context, ep *ent.EquipmentPort) (*ent.EquipmentPortDefinition, error) {
	if definition, err := ep.Edges.DefinitionOrErr(); !ent.IsNotLoaded(err) {
		return definition, err
	}
	return ep.QueryDefinition().Only(ctx)
}

func (equipmentPortResolver) ParentEquipment(ctx context.Context, ep *ent.EquipmentPort) (*ent.Equipment, error) {
	if e, err := ep.Edges.ParentOrErr(); !ent.IsNotLoaded(err) {
		return e, err
	}
	return ep.QueryParent().Only(ctx)
}

func (equipmentPortResolver) Link(ctx context.Context, ep *ent.EquipmentPort) (*ent.Link, error) {
	l, err := ep.Edges.LinkOrErr()
	if ent.IsNotLoaded(err) {
		l, err = ep.QueryLink().Only(ctx)
	}
	return l, ent.MaskNotFound(err)
}

func (equipmentPortResolver) Properties(ctx context.Context, ep *ent.EquipmentPort) ([]*ent.Property, error) {
	if properties, err := ep.Edges.PropertiesOrErr(); !ent.IsNotLoaded(err) {
		return properties, err
	}
	return ep.QueryProperties().All(ctx)
}

type equipmentPositionResolver struct{}

func (equipmentPositionResolver) Definition(ctx context.Context, ep *ent.EquipmentPosition) (*ent.EquipmentPositionDefinition, error) {
	if definition, err := ep.Edges.DefinitionOrErr(); !ent.IsNotLoaded(err) {
		return definition, err
	}
	return ep.QueryDefinition().Only(ctx)
}

func (equipmentPositionResolver) ParentEquipment(ctx context.Context, ep *ent.EquipmentPosition) (*ent.Equipment, error) {
	if parent, err := ep.Edges.ParentOrErr(); !ent.IsNotLoaded(err) {
		return parent, err
	}
	return ep.QueryParent().Only(ctx)
}

func (equipmentPositionResolver) AttachedEquipment(ctx context.Context, ep *ent.EquipmentPosition) (*ent.Equipment, error) {
	e, err := ep.Edges.AttachmentOrErr()
	if ent.IsNotLoaded(err) {
		e, err = ep.QueryAttachment().Only(ctx)
	}
	return e, ent.MaskNotFound(err)
}

type equipmentPortDefinitionResolver struct{}

func (equipmentPortDefinitionResolver) PortType(ctx context.Context, epd *ent.EquipmentPortDefinition) (*ent.EquipmentPortType, error) {
	l, err := epd.Edges.EquipmentPortTypeOrErr()
	if ent.IsNotLoaded(err) {
		l, err = epd.QueryEquipmentPortType().Only(ctx)
	}
	return l, ent.MaskNotFound(err)
}

type equipmentPortTypeResolver struct{}

func (equipmentPortTypeResolver) PropertyTypes(ctx context.Context, ept *ent.EquipmentPortType) ([]*ent.PropertyType, error) {
	if types, err := ept.Edges.PropertyTypesOrErr(); !ent.IsNotLoaded(err) {
		return types, err
	}
	return ept.QueryPropertyTypes().All(ctx)
}

func (equipmentPortTypeResolver) LinkPropertyTypes(ctx context.Context, ept *ent.EquipmentPortType) ([]*ent.PropertyType, error) {
	if types, err := ept.Edges.LinkPropertyTypesOrErr(); !ent.IsNotLoaded(err) {
		return types, err
	}
	return ept.QueryLinkPropertyTypes().All(ctx)
}

func (equipmentPortTypeResolver) NumberOfPortDefinitions(ctx context.Context, ept *ent.EquipmentPortType) (int, error) {
	if pds, err := ept.Edges.PortDefinitionsOrErr(); !ent.IsNotLoaded(err) {
		return len(pds), nil
	}
	return ept.QueryPortDefinitions().Count(ctx)
}

type equipmentTypeResolver struct{}

func (equipmentTypeResolver) Category(ctx context.Context, typ *ent.EquipmentType) (*string, error) {
	c, err := typ.Edges.CategoryOrErr()
	if ent.IsNotLoaded(err) {
		c, err = typ.QueryCategory().Only(ctx)
	}
	if err == nil {
		return &c.Name, nil
	}
	return nil, ent.MaskNotFound(err)
}

func (equipmentTypeResolver) PositionDefinitions(ctx context.Context, typ *ent.EquipmentType) ([]*ent.EquipmentPositionDefinition, error) {
	if pds, err := typ.Edges.PositionDefinitionsOrErr(); !ent.IsNotLoaded(err) {
		return pds, err
	}
	return typ.QueryPositionDefinitions().All(ctx)
}

func (equipmentTypeResolver) PortDefinitions(ctx context.Context, typ *ent.EquipmentType) ([]*ent.EquipmentPortDefinition, error) {
	if pds, err := typ.Edges.PortDefinitionsOrErr(); !ent.IsNotLoaded(err) {
		return pds, err
	}
	return typ.QueryPortDefinitions().All(ctx)
}

func (equipmentTypeResolver) PropertyTypes(ctx context.Context, typ *ent.EquipmentType) ([]*ent.PropertyType, error) {
	if pts, err := typ.Edges.PropertyTypesOrErr(); !ent.IsNotLoaded(err) {
		return pts, err
	}
	return typ.QueryPropertyTypes().All(ctx)
}

func (equipmentTypeResolver) Equipments(ctx context.Context, typ *ent.EquipmentType) ([]*ent.Equipment, error) {
	if es, err := typ.Edges.EquipmentOrErr(); !ent.IsNotLoaded(err) {
		return es, err
	}
	return typ.QueryEquipment().All(ctx)
}

func (equipmentTypeResolver) NumberOfEquipment(ctx context.Context, typ *ent.EquipmentType) (int, error) {
	if es, err := typ.Edges.EquipmentOrErr(); !ent.IsNotLoaded(err) {
		return len(es), err
	}
	return typ.QueryEquipment().Count(ctx)
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

func (equipmentResolver) ParentLocation(ctx context.Context, e *ent.Equipment) (*ent.Location, error) {
	l, err := e.Edges.LocationOrErr()
	if ent.IsNotLoaded(err) {
		l, err = e.QueryLocation().Only(ctx)
	}
	return l, ent.MaskNotFound(err)
}

func (equipmentResolver) ParentPosition(ctx context.Context, e *ent.Equipment) (*ent.EquipmentPosition, error) {
	p, err := e.Edges.ParentPositionOrErr()
	if ent.IsNotLoaded(err) {
		p, err = e.QueryParentPosition().Only(ctx)
	}
	return p, ent.MaskNotFound(err)
}

func (equipmentResolver) EquipmentType(ctx context.Context, e *ent.Equipment) (*ent.EquipmentType, error) {
	if typ, err := e.Edges.TypeOrErr(); !ent.IsNotLoaded(err) {
		return typ, err
	}
	return e.QueryType().Only(ctx)
}

func (equipmentResolver) Positions(ctx context.Context, e *ent.Equipment) ([]*ent.EquipmentPosition, error) {
	if positions, err := e.Edges.PositionsOrErr(); !ent.IsNotLoaded(err) {
		return positions, err
	}
	return e.QueryPositions().All(ctx)
}

func (equipmentResolver) Ports(ctx context.Context, e *ent.Equipment, availableOnly *bool) ([]*ent.EquipmentPort, error) {
	if !pointer.GetBool(availableOnly) {
		if ports, err := e.Edges.PortsOrErr(); !ent.IsNotLoaded(err) {
			return ports, err
		}
	}
	query := e.QueryPorts()
	if pointer.GetBool(availableOnly) {
		query.Where(equipmentport.Not(equipmentport.HasLink()))
	}
	return query.All(ctx)
}

func (equipmentResolver) Properties(ctx context.Context, e *ent.Equipment) ([]*ent.Property, error) {
	if properties, err := e.Edges.PropertiesOrErr(); !ent.IsNotLoaded(err) {
		return properties, err
	}
	return e.QueryProperties().All(ctx)
}

func (equipmentResolver) FutureState(_ context.Context, e *ent.Equipment) (*models.FutureState, error) {
	state := models.FutureState(e.FutureState)
	return &state, nil
}

func (equipmentResolver) WorkOrder(ctx context.Context, e *ent.Equipment) (*ent.WorkOrder, error) {
	wo, err := e.Edges.WorkOrderOrErr()
	if ent.IsNotLoaded(err) {
		wo, err = e.QueryWorkOrder().Only(ctx)
	}
	return wo, ent.MaskNotFound(err)
}

func (equipmentResolver) filesOfType(ctx context.Context, e *ent.Equipment, typ string) ([]*ent.File, error) {
	fds, err := e.Edges.FilesOrErr()
	if ent.IsNotLoaded(err) {
		return e.QueryFiles().
			Where(file.Type(typ)).
			All(ctx)
	}
	files := make([]*ent.File, 0, len(fds))
	for _, f := range fds {
		if f.Type == typ {
			files = append(files, f)
		}
	}
	return files, nil
}

func (r equipmentResolver) Images(ctx context.Context, e *ent.Equipment) ([]*ent.File, error) {
	return r.filesOfType(ctx, e, models.FileTypeImage.String())
}

func (r equipmentResolver) Files(ctx context.Context, e *ent.Equipment) ([]*ent.File, error) {
	return r.filesOfType(ctx, e, models.FileTypeFile.String())
}

func (equipmentResolver) Hyperlinks(ctx context.Context, e *ent.Equipment) ([]*ent.Hyperlink, error) {
	if hls, err := e.Edges.HyperlinksOrErr(); !ent.IsNotLoaded(err) {
		return hls, err
	}
	return e.QueryHyperlinks().All(ctx)
}

func (equipmentResolver) PositionHierarchy(ctx context.Context, e *ent.Equipment) ([]*ent.EquipmentPosition, error) {
	var positions []*ent.EquipmentPosition
	ppos, err := e.QueryParentPosition().Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, fmt.Errorf("querying parent position: %w", err)
	}
	for ppos != nil {
		positions = append([]*ent.EquipmentPosition{ppos}, positions...)
		p, err := ppos.QueryParent().QueryParentPosition().Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return nil, fmt.Errorf("querying parent position: %w", err)
		}

		ppos = p
	}
	return positions, nil
}

func (r equipmentResolver) LocationHierarchy(ctx context.Context, e *ent.Equipment) ([]*ent.Location, error) {
	positions, err := r.PositionHierarchy(ctx, e)
	if err != nil {
		return nil, err
	}
	var (
		locations []*ent.Location
		query     *ent.LocationQuery
	)
	if len(positions) > 0 {
		query = positions[0].QueryParent().QueryLocation()
	} else {
		query = e.QueryLocation()
	}
	for parent := query.WithParent().OnlyX(ctx); parent != nil; {
		locations = append([]*ent.Location{parent}, locations...)
		grandparent, err := parent.Edges.ParentOrErr()
		if ent.IsNotLoaded(err) {
			grandparent, err = parent.QueryParent().WithParent().Only(ctx)
		}
		if err != nil && !ent.IsNotFound(err) {
			return nil, err
		}
		parent = grandparent
	}
	return locations, nil
}

func (r equipmentResolver) Device(ctx context.Context, e *ent.Equipment) (*models.Device, error) {
	if e.DeviceID == "" {
		return nil, nil
	}
	if r.orc8r.client == nil {
		return nil, errors.New("unsupported field")
	}
	parts := strings.Split(e.DeviceID, ".")
	if len(parts) < 2 {
		return nil, errors.New("invalid equipment device id")
	}

	uri := fmt.Sprintf("/magma/v1/networks/%s/gateways/%s/status", parts[1], parts[0])
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create http request: %w", err)
	}
	rsp, err := r.orc8r.client.Do(req)
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

	dev := models.Device{ID: e.DeviceID}
	if result.CheckinTime != 0 {
		dev.Up = pointer.ToBool(
			checkinTimeIsUp(result.CheckinTime, 3*time.Minute),
		)
	}
	return &dev, nil
}

// checkinTime is milliseconds since epoch
func checkinTimeIsUp(checkinTime int64, duration time.Duration) bool {
	secs := int64(duration.Seconds()) + checkinTime/1000
	return time.Now().Before(time.Unix(secs, 0))
}

func (r equipmentResolver) Services(ctx context.Context, e *ent.Equipment) ([]*ent.Service, error) {
	eqPred := resolverutil.BuildGeneralEquipmentAncestorFilter(equipment.ID(e.ID), 1, 4)
	eids, err := r.ClientFrom(ctx).ServiceEndpoint.Query().
		Where(serviceendpoint.HasPortWith(
			equipmentport.HasParentWith(eqPred),
		)).
		IDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying service endpoint ids: %w", err)
	}

	services, err := r.ClientFrom(ctx).Service.Query().
		Where(service.HasEndpointsWith(
			serviceendpoint.IDIn(eids...),
		),
		).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying services where equipment port is an endpoint: %w", err)
	}
	ids := make([]int, len(services))
	for i, svc := range services {
		ids[i] = svc.ID
	}

	linkServices, err := r.ClientFrom(ctx).Service.Query().Where(
		service.HasLinksWith(link.HasPortsWith(equipmentport.HasParentWith(
			resolverutil.BuildGeneralEquipmentAncestorFilter(equipment.ID(e.ID), 1, 3)))),
		service.Not(service.IDIn(ids...))).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying services where equipment connected to link of service: %w", err)
	}
	return append(services, linkServices...), nil
}
