// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"testing"

	"github.com/facebookincubator/symphony/graph/resolverutil"

	"github.com/facebookincubator/symphony/graph/graphql/models"

	"github.com/facebookincubator/symphony/graph/ent/propertytype"

	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/stretchr/testify/require"
)

var (
	parentsAHeader = []string{"Parent Equipment (3) A", "Position (3) A", "Parent Equipment (2) A", "Position (2) A", "Parent Equipment A", "Equipment Position A"}
	parentsBHeader = []string{"Parent Equipment (3) B", "Position (3) B", "Parent Equipment (2) B", "Position (2) B", "Parent Equipment B", "Equipment Position B"}
)

func TestLinkTitleInputValidation(t *testing.T) {
	r, err := newImporterTestResolver(t)
	require.NoError(t, err)
	importer := r.importer
	defer r.drv.Close()

	ctx := newImportContext(viewertest.NewContext(r.client))
	prepareBasicData(ctx, t, *r)

	err = importer.inputValidationsLinks(ctx, NewImportHeader([]string{"aa"}, ImportEntityLink))
	require.Error(t, err)
	err = importer.inputValidationsLinks(ctx, NewImportHeader(fixedFirstPortLink, ImportEntityLink))
	require.Error(t, err)
	minimalRow := append(append(append(append(fixedFirstPortLink, parentsAHeader...), fixedSecondPortLink...), parentsBHeader...), "Service Names")
	err = importer.inputValidationsLinks(ctx, NewImportHeader(minimalRow, ImportEntityLink))
	require.NoError(t, err)
}

func TestGeneralLinksEditImport(t *testing.T) {
	r, err := newImporterTestResolver(t)
	require.NoError(t, err)
	importer := r.importer
	defer r.drv.Close()

	ctx := newImportContext(viewertest.NewContext(r.client))
	ids := preparePortTypeData(ctx, t, *r)

	def1 := r.client.EquipmentPortDefinition.GetX(ctx, ids.portDef1)
	equip2 := r.client.Equipment.GetX(ctx, ids.equipParent2ID)
	etyp1 := equip2.QueryType().OnlyX(ctx)

	def2 := r.client.EquipmentPortDefinition.GetX(ctx, ids.portDef2)
	childEquip := r.client.Equipment.GetX(ctx, ids.equipChildID)
	etyp2 := childEquip.QueryType().OnlyX(ctx)
	var (
		row1 = []string{ids.linkID, def1.Name, equip2.Name, etyp1.Name, locationL, locationM, locationS, "", "", "", "", "", "", def2.Name, childEquip.Name, etyp2.Name, locationL, locationM, locationS, "", "", "", "", parentEquip, posName, "", "44", "2019-01-01", "FALSE"}
	)
	firstPortHeader := append(append(fixedFirstPortLink, locTypeNameL, locTypeNameM, locTypeNameS), parentsAHeader...)
	secondPortHeader := append(append(fixedSecondPortLink, locTypeNameL, locTypeNameM, locTypeNameS), parentsBHeader...)

	header := append(append(firstPortHeader, secondPortHeader...), "Service Names", propNameInt, propNameDate, propNameBool)

	fl := NewImportHeader(header, ImportEntityLink)
	err = importer.inputValidationsLinks(ctx, fl)
	require.NoError(t, err)

	r1 := NewImportRecord(row1, fl)

	link, err := importer.validateLineForExistingLink(ctx, ids.linkID, r1)
	require.NoError(t, err)
	ports := link.QueryPorts().AllX(ctx)
	for _, port := range ports {
		propertyTypes, err := importer.validatePropertiesForPortType(ctx, r1, port.QueryDefinition().QueryEquipmentPortType().OnlyX(ctx), ImportEntityLink)
		require.NoError(t, err)
		if port.ID == ids.childPortInst1 {
			require.Len(t, propertyTypes, 2)
			require.NotEqual(t, propertyTypes[0].PropertyTypeID, propertyTypes[1].PropertyTypeID)
			for _, value := range propertyTypes {
				ptyp := etyp2.QueryPortDefinitions().QueryEquipmentPortType().QueryLinkPropertyTypes().Where(propertytype.ID(value.PropertyTypeID)).OnlyX(ctx)
				switch ptyp.Name {
				case propNameDate:
					require.Equal(t, *value.StringValue, "2019-01-01")
					require.Equal(t, ptyp.Type, models.PropertyKindDate.String())
				case propNameBool:
					require.Equal(t, *value.BooleanValue, false)
					require.Equal(t, ptyp.Type, models.PropertyKindBool.String())
				default:
					require.Fail(t, "property type name should be one of the two")
				}
			}
		} else if port.ID == ids.parentPortInst2 {
			require.Len(t, propertyTypes, 1)
			val := propertyTypes[0]
			ptyp := etyp1.QueryPortDefinitions().QueryEquipmentPortType().QueryLinkPropertyTypes().Where(propertytype.ID(val.PropertyTypeID)).OnlyX(ctx)
			require.Equal(t, ptyp.Name, propNameInt)
			require.Equal(t, *val.IntValue, 44)
			require.Equal(t, ptyp.Type, models.PropertyKindInt.String())
		}
	}
}

func TestGeneralLinksAddImport(t *testing.T) {
	r, err := newImporterTestResolver(t)
	require.NoError(t, err)
	importer := r.importer
	defer r.drv.Close()

	ctx := newImportContext(viewertest.NewContext(r.client))
	ids := preparePortTypeData(ctx, t, *r)
	def1 := r.client.EquipmentPortDefinition.GetX(ctx, ids.portDef1)
	equipParent := r.client.Equipment.GetX(ctx, ids.equipParentID)
	etyp1 := equipParent.QueryType().OnlyX(ctx)
	posDef := etyp1.QueryPositionDefinitions().OnlyX(ctx)

	def2 := r.client.EquipmentPortDefinition.GetX(ctx, ids.portDef2)
	childEquip := r.client.Equipment.GetX(ctx, ids.equipChild2ID)
	etyp2 := childEquip.QueryType().OnlyX(ctx)

	equipParent2 := r.client.Equipment.GetX(ctx, ids.equipParent2ID)
	var (
		row1 = []string{"", def1.Name, equipParent.Name, etyp1.Name, locationL, locationM, locationS, "", "", "", "", "", "", def2.Name, childEquip.Name, etyp2.Name, locationL, locationM, locationS, "", "", "", "", parentEquip2, posName, "", "44", "2019-01-01", "FALSE"}
	)
	firstPortHeader := append(append(fixedFirstPortLink, locTypeNameL, locTypeNameM, locTypeNameS), parentsAHeader...)
	secondPortHeader := append(append(fixedSecondPortLink, locTypeNameL, locTypeNameM, locTypeNameS), parentsBHeader...)

	header := append(append(firstPortHeader, secondPortHeader...), "Service Names", propNameInt, propNameDate, propNameBool)

	fl := NewImportHeader(header, ImportEntityLink)
	err = importer.inputValidationsLinks(ctx, fl)
	require.NoError(t, err)
	r1 := NewImportRecord(row1, fl)

	pr1, pr2, err := importer.getTwoPortRecords(r1)
	require.NoError(t, err)
	// port1 test
	require.Equal(t, row1[1:13], pr1.line)
	parentLoc, err := importer.verifyOrCreateLocationHierarchy(ctx, *pr1)
	require.NoError(t, err)
	require.Equal(t, locationS, parentLoc.Name)
	eqID, defID, err := importer.getPositionDetailsIfExists(ctx, parentLoc, *pr1, false)
	require.NoError(t, err)
	require.Nil(t, eqID)
	require.Nil(t, defID)
	eq, created, err := importer.getOrCreateEquipment(ctx, importer.r.Mutation(), pr1.PortEquipmentName(), etyp1, nil, parentLoc, nil, nil)
	require.NoError(t, err)
	require.False(t, created)
	require.Equal(t, equipParent.ID, eq.ID)

	portType, err := def1.QueryEquipmentPortType().Only(ctx)
	require.NoError(t, err)
	propertyInputs, err := importer.validatePropertiesForPortType(ctx, r1, portType, ImportEntityLink)
	require.NoError(t, err)
	require.Len(t, propertyInputs, 1)
	require.Equal(t, *propertyInputs[0].IntValue, 44)
	// port2 test
	require.Equal(t, row1[13:25], pr2.line)
	parentLoc, err = importer.verifyOrCreateLocationHierarchy(ctx, *pr2)
	require.NoError(t, err)
	require.Equal(t, locationS, parentLoc.Name)
	eqID, defID, err = importer.getPositionDetailsIfExists(ctx, parentLoc, *pr2, false)
	require.NoError(t, err)
	require.Equal(t, equipParent2.ID, *eqID)
	require.Equal(t, posDef.ID, *defID)
	pos, err := resolverutil.GetOrCreatePosition(ctx, importer.ClientFrom(ctx), eqID, defID, false)
	require.NoError(t, err)
	eq, created, err = importer.getOrCreateEquipment(ctx, importer.r.Mutation(), pr2.PortEquipmentName(), etyp2, nil, nil, pos, nil)
	require.NoError(t, err)
	require.False(t, created)
	require.Equal(t, ids.equipChild2ID, eq.ID)
	portType, err = def2.QueryEquipmentPortType().Only(ctx)
	require.NoError(t, err)
	propertyInputs, err = importer.validatePropertiesForPortType(ctx, r1, portType, ImportEntityLink)
	require.NoError(t, err)
	require.Len(t, propertyInputs, 2)
	require.NotEqual(t, propertyInputs[0].PropertyTypeID, propertyInputs[1].PropertyTypeID)
	for _, inp := range propertyInputs {
		ptype := importer.ClientFrom(ctx).PropertyType.GetX(ctx, inp.PropertyTypeID)
		if ptype.Type == models.PropertyKindDate.String() {
			require.Equal(t, "2019-01-01", *inp.StringValue)
		} else if ptype.Type == models.PropertyKindBool.String() {
			require.False(t, *inp.BooleanValue)
		}
	}
}
