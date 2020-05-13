// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exporter

import (
	"bytes"
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/facebookincubator/ent/dialect"
	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/authz"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/enttest"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"
	"github.com/facebookincubator/symphony/graph/ent/migrate"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpointdefinition"
	"github.com/facebookincubator/symphony/graph/event"
	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/graphql/resolver"
	"github.com/facebookincubator/symphony/graph/importer"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/facebookincubator/symphony/pkg/log/logtest"
	"github.com/facebookincubator/symphony/pkg/testdb"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/require"
)

var debug = flag.Bool("debug", false, "run database driver on debug mode")

const (
	equipmentTypeName          = "equipmentType"
	equipmentType2Name         = "equipmentType2"
	parentEquip                = "parentEquipmentName"
	currEquip                  = "currEquipmentName"
	currEquip2                 = "currEquipmentName2"
	positionName               = "Position"
	portName1                  = "port1"
	portName2                  = "port2"
	portName3                  = "port3"
	propNameStr                = "propNameStr"
	propNameDate               = "propNameDate"
	propNameBool               = "propNameBool"
	propNameInt                = "propNameInt"
	externalIDL                = "11"
	externalIDM                = "22"
	lat                        = 32.109
	long                       = 34.855
	newPropNameStr             = "newPropNameStr"
	propDefValue               = "defaultVal"
	propDefValue2              = "defaultVal2"
	propDevValInt              = 15
	propInstanceValue          = "newVal"
	locTypeNameL               = "locTypeLarge"
	locTypeNameM               = "locTypeMedium"
	locTypeNameS               = "locTypeSmall"
	grandParentLocation        = "grandParentLocation"
	parentLocation             = "parentLocation"
	childLocation              = "childLocation"
	firstServiceName           = "S1"
	secondServiceName          = "S2"
	MethodAdd           method = "ADD"
	MethodEdit          method = "EDIT"
)

type method string

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

type TestExporterResolver struct {
	generated.ResolverRoot
	drv      dialect.Driver
	client   *ent.Client
	exporter exporter
}

func newExporterTestResolver(t *testing.T) *TestExporterResolver {
	db, name, err := testdb.Open()
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	return newResolver(t, sql.OpenDB(name, db))
}

func newResolver(t *testing.T, drv dialect.Driver) *TestExporterResolver {
	if *debug {
		drv = dialect.Debug(drv)
	}
	client := enttest.NewClient(t,
		enttest.WithOptions(ent.Driver(drv)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	logger := logtest.NewTestLogger(t)
	r := resolver.New(resolver.Config{
		Logger:     logger,
		Subscriber: event.NewNopSubscriber(),
	})
	e := exporter{logger, equipmentRower{logger}}
	return &TestExporterResolver{r, drv, client, e}
}

func prepareData(ctx context.Context, t *testing.T, r TestExporterResolver) {
	mr := r.Mutation()

	locTypeL, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: locTypeNameL})
	require.NoError(t, err)
	locTypeM, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: locTypeNameM})
	require.NoError(t, err)
	locTypeS, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: locTypeNameS, Properties: []*models.PropertyTypeInput{
		{
			Name:        propNameStr,
			Type:        models.PropertyKindString,
			StringValue: pointer.ToString("default"),
		},
		{
			Name: propNameBool,
			Type: models.PropertyKindBool,
		},
		{
			Name:        propNameDate,
			Type:        models.PropertyKindDate,
			StringValue: pointer.ToString("1988-03-29"),
		},
	}})
	require.NoError(t, err)

	_, err = mr.EditLocationTypesIndex(ctx, []*models.LocationTypeIndex{
		{
			LocationTypeID: locTypeL.ID,
			Index:          0,
		},
		{
			LocationTypeID: locTypeM.ID,
			Index:          1,
		},
		{
			LocationTypeID: locTypeS.ID,
			Index:          2,
		},
	})
	require.NoError(t, err)

	gpLocation, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name:       grandParentLocation,
		Type:       locTypeL.ID,
		ExternalID: pointer.ToString(externalIDL),
		Latitude:   pointer.ToFloat64(lat),
		Longitude:  pointer.ToFloat64(long),
	})

	require.NoError(t, err)
	pLocation, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name:       parentLocation,
		Type:       locTypeM.ID,
		Parent:     &gpLocation.ID,
		ExternalID: pointer.ToString(externalIDM),
		Latitude:   pointer.ToFloat64(lat),
		Longitude:  pointer.ToFloat64(long),
	})
	require.NoError(t, err)
	strProp := locTypeS.QueryPropertyTypes().Where(propertytype.Type(models.PropertyKindString.String())).OnlyX(ctx)
	boolProp := locTypeS.QueryPropertyTypes().Where(propertytype.Type(models.PropertyKindBool.String())).OnlyX(ctx)
	clocation, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name:   childLocation,
		Type:   locTypeS.ID,
		Parent: &pLocation.ID,
		Properties: []*models.PropertyInput{
			{
				PropertyTypeID: strProp.ID,
				StringValue:    pointer.ToString("override"),
			},
			{
				PropertyTypeID: boolProp.ID,
				BooleanValue:   pointer.ToBool(true),
			},
		},
	})
	require.NoError(t, err)
	position1 := models.EquipmentPositionInput{
		Name: positionName,
	}

	ptyp, _ := mr.AddEquipmentPortType(ctx, models.AddEquipmentPortTypeInput{
		Name: "portType1",
		Properties: []*models.PropertyTypeInput{
			{
				Name:        propStr,
				Type:        "string",
				StringValue: pointer.ToString("t1"),
			},
			{
				Name: propStr2,
				Type: "string",
			},
		},
		LinkProperties: []*models.PropertyTypeInput{
			{
				Name:        propNameStr,
				Type:        "string",
				StringValue: pointer.ToString("t1"),
			},
			{
				Name: propNameBool,
				Type: "bool",
			},
			{
				Name:     propNameInt,
				Type:     "int",
				IntValue: pointer.ToInt(100),
			},
		},
	})
	port1 := models.EquipmentPortInput{
		Name:       portName1,
		PortTypeID: &ptyp.ID,
	}
	strDefVal := propDefValue
	intDefVal := propDevValInt
	propDefInput1 := models.PropertyTypeInput{
		Name:        propNameStr,
		Type:        "string",
		StringValue: &strDefVal,
	}
	propDefInput2 := models.PropertyTypeInput{
		Name:     propNameInt,
		Type:     "int",
		IntValue: &intDefVal,
	}
	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:      equipmentTypeName,
		Positions: []*models.EquipmentPositionInput{&position1},
		Ports:     []*models.EquipmentPortInput{&port1},
	})
	require.NoError(t, err)

	port2 := models.EquipmentPortInput{
		Name: portName2,
	}
	equipmentType2, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:       equipmentType2Name,
		Properties: []*models.PropertyTypeInput{&propDefInput1, &propDefInput2},
		Ports:      []*models.EquipmentPortInput{&port2},
	})
	require.NoError(t, err)

	posDef1 := equipmentType.QueryPositionDefinitions().Where(equipmentpositiondefinition.Name(positionName)).OnlyX(ctx)
	propDef1 := equipmentType2.QueryPropertyTypes().Where(propertytype.Name(propNameStr)).OnlyX(ctx)

	parentEquipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     parentEquip,
		Type:     equipmentType.ID,
		Location: &clocation.ID,
	})
	require.NoError(t, err)

	strVal := propInstanceValue
	propInstance1 := models.PropertyInput{
		PropertyTypeID: propDef1.ID,
		StringValue:    &strVal,
	}
	childEquip, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:               currEquip,
		Type:               equipmentType2.ID,
		Parent:             &parentEquipment.ID,
		PositionDefinition: &posDef1.ID,
		ExternalID:         pointer.ToString(externalIDM),
		Properties:         []*models.PropertyInput{&propInstance1},
	})
	require.NoError(t, err)

	portDef1 := equipmentType.QueryPortDefinitions().Where(equipmentportdefinition.Name(portName1)).OnlyX(ctx)
	portDef2 := equipmentType2.QueryPortDefinitions().Where(equipmentportdefinition.Name(portName2)).OnlyX(ctx)
	_, _ = mr.AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{Equipment: parentEquipment.ID, Port: portDef1.ID},
			{Equipment: childEquip.ID, Port: portDef2.ID},
		},
	})

	val := propDefValue2
	propertyInput := models.PropertyTypeInput{Name: newPropNameStr, StringValue: &val, Type: models.PropertyKindString}
	_, err = r.Mutation().EditEquipmentType(ctx, models.EditEquipmentTypeInput{
		ID:         equipmentType2.ID,
		Name:       equipmentType2.Name,
		Properties: []*models.PropertyTypeInput{&propertyInput},
	})
	require.NoError(t, err)

	portID1, err := parentEquipment.QueryPorts().Where(equipmentport.HasDefinitionWith(equipmentportdefinition.ID(portDef1.ID))).OnlyID(ctx)
	require.NoError(t, err)
	portID2, err := childEquip.QueryPorts().Where(equipmentport.HasDefinitionWith(equipmentportdefinition.ID(portDef2.ID))).OnlyID(ctx)
	require.NoError(t, err)

	serviceType, _ := mr.AddServiceType(ctx,
		models.ServiceTypeCreateData{
			Name:        "L2 Service",
			HasCustomer: false,
			Endpoints: []*models.ServiceEndpointDefinitionInput{
				{
					Name:            "endpoint type1",
					Role:            pointer.ToString("CONSUMER"),
					Index:           0,
					EquipmentTypeID: equipmentType.ID,
				},
				{
					Index:           1,
					Name:            "endpoint type2",
					Role:            pointer.ToString("PROVIDER"),
					EquipmentTypeID: equipmentType2.ID,
				},
			},
		})
	s1, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          firstServiceName,
		ServiceTypeID: serviceType.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)
	s2, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          secondServiceName,
		ServiceTypeID: serviceType.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	ept0 := serviceType.QueryEndpointDefinitions().Where(serviceendpointdefinition.Index(0)).OnlyX(ctx)

	_, _ = mr.AddServiceEndpoint(ctx, models.AddServiceEndpointInput{
		ID:          s1.ID,
		EquipmentID: parentEquipment.ID,
		PortID:      pointer.ToInt(portID1),
		Definition:  ept0.ID,
	})

	_, _ = mr.AddServiceEndpoint(ctx, models.AddServiceEndpointInput{
		ID:          s2.ID,
		EquipmentID: parentEquipment.ID,
		PortID:      pointer.ToInt(portID1),
		Definition:  ept0.ID,
	})

	ept1 := serviceType.QueryEndpointDefinitions().Where(serviceendpointdefinition.Index(1)).OnlyX(ctx)

	_, _ = mr.AddServiceEndpoint(ctx, models.AddServiceEndpointInput{
		ID:          s1.ID,
		EquipmentID: childEquip.ID,
		PortID:      pointer.ToInt(portID2),
		Definition:  ept1.ID,
	})
	/*
		helper: data now is of type:
		loc(grandParent):
			loc(parent):
				loc(child):
						parentEquipment(equipmentType): with portType1 (has 2 string props)
						childEquipment(equipmentType2): (no props props)
						these ports are linked together
		services:
			firstService:
					endpoints: parentEquipment consumer, childEquipment provider
			secondService:
					endpoints: parentEquipment consumer
	*/
}

func prepareHandlerAndExport(t *testing.T, r *TestExporterResolver, e http.Handler) (context.Context, *http.Response) {
	auth := authz.Handler(e, logtest.NewTestLogger(t))
	th := viewer.TenancyHandler(auth,
		viewer.NewFixedTenancy(r.client),
		logtest.NewTestLogger(t),
	)
	server := httptest.NewServer(th)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	viewertest.SetDefaultViewerHeaders(req)

	ctx := viewertest.NewContext(context.Background(), r.client)
	prepareData(ctx, t, *r)
	locs := r.client.Location.Query().AllX(ctx)
	require.Len(t, locs, 3)
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return ctx, res
}

func importLinksPortsFile(t *testing.T, client *ent.Client, r io.Reader, entity importer.ImportEntity, method method, skipLines, withVerify bool) {
	readr := csv.NewReader(r)
	var buf *bytes.Buffer
	var contentType, url string
	switch entity {
	case importer.ImportEntityLink:
		buf, contentType = writeModifiedLinksCSV(t, readr, method, skipLines, withVerify)
	case importer.ImportEntityPort:
		buf, contentType = writeModifiedPortsCSV(t, readr, skipLines, withVerify)
	}

	h, _ := importer.NewHandler(
		importer.Config{
			Logger:     logtest.NewTestLogger(t),
			Subscriber: event.NewNopSubscriber(),
		},
	)
	auth := authz.Handler(h, logtest.NewTestLogger(t))
	th := viewer.TenancyHandler(auth,
		viewer.NewFixedTenancy(client),
		logtest.NewTestLogger(t),
	)
	server := httptest.NewServer(th)
	defer server.Close()
	switch entity {
	case importer.ImportEntityLink:
		url = server.URL + "/export_links"
	case importer.ImportEntityPort:
		fmt.Println("server.URL", server.URL)
		url = server.URL + "/export_ports"
	}
	req, err := http.NewRequest(http.MethodPost, url, buf)
	require.Nil(t, err)

	viewertest.SetDefaultViewerHeaders(req)
	req.Header.Set("Content-Type", contentType)
	resp, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)
	resp.Body.Close()
}
