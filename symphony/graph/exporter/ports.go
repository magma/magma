// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exporter

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/AlekSi/pointer"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/resolverutil"
	"github.com/facebookincubator/symphony/pkg/ctxgroup"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type portfilterInput struct {
	Name          models.EquipmentFilterType `json:"name"`
	Operator      models.FilterOperator      `jsons:"operator"`
	StringValue   string                     `json:"stringValue"`
	IDSet         []string                   `json:"idSet"`
	PropertyValue models.PropertyTypeInput   `json:"propertyValue"`
	BoolValue     bool                       `json:"boolValue"`
}

type portsRower struct {
	log log.Logger
}

func (er portsRower) rows(ctx context.Context, url *url.URL) ([][]string, error) {
	log := er.log.For(ctx)

	var (
		err            error
		filterInput    []*models.PortFilterInput
		portDataHeader = [...]string{bom + "Port ID", "Port Name", "Port Type", "Equipment Name", "Equipment Type"}
		parentsHeader  = [...]string{"Parent Equipment (3)", "Parent Equipment (2)", "Parent Equipment", "Equipment Position"}
		linkHeader     = [...]string{"Linked Port ID", "Linked Port Name", "Linked Equipment ID", "Linked Equipment"}
		serviceHeader  = [...]string{"Consumer Endpoint for These Services", "Provider Endpoint for These Services"}
	)
	filtersParam := url.Query().Get("filters")
	if filtersParam != "" {
		filterInput, err = paramToPortFilterInput(filtersParam)
		if err != nil {
			log.Error("cannot filter ports", zap.Error(err))
			return nil, errors.Wrap(err, "cannot filter ports")
		}
	}
	client := ent.FromContext(ctx)

	ports, err := resolverutil.PortSearch(ctx, client, filterInput, nil)
	if err != nil {
		log.Error("cannot query ports", zap.Error(err))
		return nil, errors.Wrap(err, "cannot query ports")
	}
	cg := ctxgroup.WithContext(ctx, ctxgroup.MaxConcurrency(32))

	portsList := ports.Ports
	allrows := make([][]string, len(portsList)+1)

	var orderedLocTypes, propertyTypes []string
	cg.Go(func(ctx context.Context) error {
		orderedLocTypes, err = locationTypeHierarchy(ctx, client)
		if err != nil {
			log.Error("cannot query location types", zap.Error(err))
			return errors.Wrap(err, "cannot query location types")
		}
		return nil
	})
	cg.Go(func(ctx context.Context) error {
		portIDs := make([]string, len(portsList))
		for i, p := range portsList {
			portIDs[i] = p.ID
		}
		propertyTypes, err = propertyTypesSlice(ctx, portIDs, client, models.PropertyEntityPort)
		if err != nil {
			log.Error("cannot query property types", zap.Error(err))
			return errors.Wrap(err, "cannot query property types")
		}
		return nil
	})
	if err := cg.Wait(); err != nil {
		return nil, err
	}

	title := append(portDataHeader[:], orderedLocTypes...)
	title = append(title, parentsHeader[:]...)
	title = append(title, linkHeader[:]...)
	title = append(title, serviceHeader[:]...)
	title = append(title, propertyTypes...)

	allrows[0] = title
	cg = ctxgroup.WithContext(ctx, ctxgroup.MaxConcurrency(32))
	for i, value := range portsList {
		value, i := value, i
		cg.Go(func(ctx context.Context) error {
			row, err := portToSlice(ctx, value, orderedLocTypes, propertyTypes)
			if err != nil {
				return err
			}
			allrows[i+1] = row
			return nil
		})
	}
	if err := cg.Wait(); err != nil {
		log.Error("error in wait", zap.Error(err))
		return nil, errors.WithMessage(err, "error in wait")
	}
	return allrows, nil
}

// nolint: ineffassign
func portToSlice(ctx context.Context, port *ent.EquipmentPort, orderedLocTypes []string, propertyTypes []string) ([]string, error) {
	var (
		posName              string
		lParents, properties []string
		linkData             = make([]string, 4)
		eParents             = make([]string, maxEquipmentParents)
		serviceData          = make([]string, 2)
	)
	parentEquip, err := port.QueryParent().Only(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "querying equipment for port (id=%s)", port.ID)
	}
	portDefinition := port.QueryDefinition().OnlyX(ctx)
	g := ctxgroup.WithContext(ctx)

	g.Go(func(ctx context.Context) error {
		lParents, err = locationHierarchyForEquipment(ctx, parentEquip, orderedLocTypes)
		return err
	})
	g.Go(func(ctx context.Context) error {
		properties, err = propertiesSlice(ctx, port, propertyTypes, models.PropertyEntityPort)
		return err
	})
	g.Go(func(ctx context.Context) error {
		pos, err := parentEquip.QueryParentPosition().Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return err
		}
		err = nil
		if pos != nil {
			def, err := pos.QueryDefinition().Only(ctx)
			if err != nil {
				return err
			}
			posName = def.Name
			eParents = parentHierarchy(ctx, *parentEquip)
		}
		return nil
	})
	g.Go(func(ctx context.Context) error {
		link, err := port.QueryLink().Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return err
		}
		if ent.IsNotFound(err) {
			err = nil
			return nil
		}
		err = nil

		if link != nil {
			otherPort, err := link.QueryPorts().Where(equipmentport.Not(equipmentport.ID(port.ID))).Only(ctx)
			if err != nil {
				return err
			}
			otherEquip := otherPort.QueryParent().OnlyX(ctx)
			linkData = []string{otherPort.ID, otherPort.QueryDefinition().OnlyX(ctx).Name, otherEquip.ID, otherEquip.Name}
		}
		return nil
	})
	g.Go(func(ctx context.Context) error {
		consumerServicesStr, err := getServicesOfPortAsEndpoint(ctx, port, models.ServiceEndpointRoleConsumer)
		if err != nil {
			return err
		}
		providerServicesStr, err := getServicesOfPortAsEndpoint(ctx, port, models.ServiceEndpointRoleProvider)
		if err != nil {
			return err
		}
		serviceData = []string{consumerServicesStr, providerServicesStr}
		return nil
	})
	if err := g.Wait(); err != nil {
		return nil, err
	}

	portType := ""
	pt, err := portDefinition.QueryEquipmentPortType().Only(ctx)
	if err == nil {
		portType = pt.Name
	}
	row := []string{port.ID, portDefinition.Name, portType, parentEquip.Name, parentEquip.QueryType().OnlyX(ctx).Name}
	row = append(row, lParents...)
	row = append(row, eParents...)
	row = append(row, posName)
	row = append(row, linkData...)
	row = append(row, serviceData...)
	row = append(row, properties...)

	return row, nil
}

func getServicesOfPortAsEndpoint(ctx context.Context, port *ent.EquipmentPort, role models.ServiceEndpointRole) (string, error) {
	services, err := port.
		QueryEndpoints().
		Where(serviceendpoint.Role(role.String())).
		QueryService().
		All(ctx)
	if err != nil {
		return "", errors.Wrapf(err, "querying port for services (id=%s)", port.ID)
	}
	var servicesList []string
	for _, service := range services {
		servicesList = append(servicesList, service.Name)
	}
	return strings.Join(servicesList, ";"), nil
}

func paramToPortFilterInput(params string) ([]*models.PortFilterInput, error) {
	var ret []*models.PortFilterInput
	var inputs []portfilterInput
	err := json.Unmarshal([]byte(params), &inputs)
	if err != nil {
		return nil, err
	}

	for _, f := range inputs {
		upperName := strings.ToUpper(f.Name.String())
		upperOp := strings.ToUpper(f.Operator.String())
		StringVal := f.StringValue
		propVal := f.PropertyValue
		maxDepth := 5
		inp := models.PortFilterInput{
			FilterType:    models.PortFilterType(upperName),
			Operator:      models.FilterOperator(upperOp),
			StringValue:   &StringVal,
			PropertyValue: &propVal,
			BoolValue:     pointer.ToBool(f.BoolValue),
			IDSet:         f.IDSet,
			MaxDepth:      &maxDepth,
		}
		ret = append(ret, &inp)
	}
	return ret, nil
}
