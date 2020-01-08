// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"testing"

	"github.com/facebookincubator/symphony/graph/graphql/models"

	"github.com/facebookincubator/symphony/graph/ent/propertytype"

	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/stretchr/testify/require"
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
	err = importer.inputValidationsLinks(ctx, NewImportHeader(fixedFirstLineLink[:5], ImportEntityLink))
	require.Error(t, err)
	err = importer.inputValidationsLinks(ctx, NewImportHeader(fixedFirstLineLink, ImportEntityLink))
	require.NoError(t, err)
}

func TestGeneralLinksImport(t *testing.T) {
	r, err := newImporterTestResolver(t)
	require.NoError(t, err)
	importer := r.importer
	defer r.drv.Close()

	ctx := newImportContext(viewertest.NewContext(r.client))
	ids := preparePortTypeData(ctx, t, *r)

	def1 := r.client.EquipmentPortDefinition.GetX(ctx, ids.portDef1)
	typ1 := def1.QueryEquipmentPortType().OnlyX(ctx)
	equip2 := r.client.Equipment.GetX(ctx, ids.equipParent2ID)
	etyp1 := equip2.QueryType().OnlyX(ctx)

	def2 := r.client.EquipmentPortDefinition.GetX(ctx, ids.portDef2)
	typ2 := def2.QueryEquipmentPortType().OnlyX(ctx)
	childEquip := r.client.Equipment.GetX(ctx, ids.equipChildID)
	etyp2 := childEquip.QueryType().OnlyX(ctx)
	var (
		row1 = []string{ids.linkID, ids.parentPortInst2, def1.Name, typ1.Name, equip2.ID, equip2.Name, etyp1.Name, ids.childPortInst1, def2.Name, typ2.Name, childEquip.ID, childEquip.Name, etyp2.Name, "", "44", "2019-01-01", "FALSE"}
	)

	header := append(fixedFirstLineLink, []string{propNameInt, propNameDate, propNameBool}...)
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
