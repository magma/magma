// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"context"
	"strconv"
	"testing"

	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/stretchr/testify/require"
)

const (
	equipmentTypeName  = "equipmentType"
	equipmentType2Name = "equipmentType2"
	equipmentType3Name = "equipmentType3"
	propName1          = "prop1"
	propName2          = "prop2"
	propName3          = "prop3"
	propName4          = "prop4"
	propName5          = "prop5"
	propName6          = "prop6"
	propName7          = "prop7"
	propName8          = "prop8"
	propDefValue       = "defaultVal"
	locTypeNameL       = "locTypeLarge"
	locTypeNameM       = "locTypeMedium"
	locTypeNameS       = "locTypeSmall"
)

type locTypeIDs struct {
	locTypeIDL int
	locTypeIDM int
	locTypeIDS int
}
type ids struct {
	locTypeIDL   int
	locTypeIDM   int
	locTypeIDS   int
	equipTypeID  int
	equipTypeID2 int
	equipTypeID3 int
}

func prepareEquipmentTypeData(ctx context.Context, t *testing.T, r TestImporterResolver) ids {
	lids := prepareBasicData(ctx, t, r)
	mr := r.importer.r.Mutation()

	strDefVal := propDefValue
	propDefInput1 := models.PropertyTypeInput{
		Name:        propName1,
		Type:        "string",
		StringValue: &strDefVal,
	}
	propDefInput2 := models.PropertyTypeInput{
		Name: propName2,
		Type: "int",
	}
	propDefInput3 := models.PropertyTypeInput{
		Name: propName3,
		Type: "date",
	}
	propDefInput4 := models.PropertyTypeInput{
		Name: propName4,
		Type: "bool",
	}
	propDefInput5 := models.PropertyTypeInput{
		Name: propName5,
		Type: "range",
	}
	propDefInput6 := models.PropertyTypeInput{
		Name: propName6,
		Type: "gps_location",
	}

	equipmentType1, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:       equipmentTypeName,
		Properties: []*models.PropertyTypeInput{&propDefInput1, &propDefInput2},
	})
	require.NoError(t, err)
	equipmentType2, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:       equipmentType2Name,
		Properties: []*models.PropertyTypeInput{&propDefInput3, &propDefInput4},
	})
	require.NoError(t, err)
	equipmentType3, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:       equipmentType3Name,
		Properties: []*models.PropertyTypeInput{&propDefInput5, &propDefInput6},
	})
	require.NoError(t, err)
	return ids{
		equipTypeID:  equipmentType1.ID,
		equipTypeID2: equipmentType2.ID,
		equipTypeID3: equipmentType3.ID,
		locTypeIDS:   lids.locTypeIDS,
		locTypeIDM:   lids.locTypeIDM,
		locTypeIDL:   lids.locTypeIDL,
	}
}

func prepareBasicData(ctx context.Context, t *testing.T, r TestImporterResolver) locTypeIDs {
	mr := r.importer.r.Mutation()

	l, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: locTypeNameL})
	require.NoError(t, err)
	m, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: locTypeNameM})
	require.NoError(t, err)
	s, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: locTypeNameS})
	require.NoError(t, err)

	_, err = r.importer.r.Mutation().EditLocationTypesIndex(ctx, []*models.LocationTypeIndex{
		{
			LocationTypeID: l.ID,
			Index:          0,
		},
		{
			LocationTypeID: m.ID,
			Index:          1,
		},
		{
			LocationTypeID: s.ID,
			Index:          2,
		},
	})
	require.NoError(t, err)

	return locTypeIDs{
		locTypeIDS: s.ID,
		locTypeIDM: m.ID,
		locTypeIDL: l.ID,
	}
}

func TestTitleLocationTypeInputValidation(t *testing.T) {
	r := newImporterTestResolver(t)
	importer := r.importer
	defer r.drv.Close()

	ctx := newImportContext(viewertest.NewContext(r.client))
	ic := getImportContext(ctx)
	var (
		equipDataHeader = [...]string{"Equipment ID", "Equipment Name", "Equipment Type", "External ID"}
		parentsHeader   = [...]string{"Parent Equipment (3)", "Position (3)", "Parent Equipment (2)", "Position (2)", "Parent Equipment", "Equipment Position"}
	)
	locTypeIDS := prepareBasicData(ctx, t, *r)
	header, _ := NewImportHeader([]string{"aa"}, ImportEntityEquipment)
	err := importer.inputValidations(ctx, header)
	require.Error(t, err)

	header, _ = NewImportHeader(equipDataHeader[:], ImportEntityEquipment)
	err = importer.inputValidations(ctx, header)
	require.Error(t, err)

	locationTypeNotInOrder := append(append(equipDataHeader[:], []string{locTypeNameS, locTypeNameM, locTypeNameL}...), parentsHeader[:]...)

	header, _ = NewImportHeader(locationTypeNotInOrder, ImportEntityEquipment)
	err = importer.inputValidations(ctx, header)
	require.Error(t, err)

	locationTypeInOrder := append(append(equipDataHeader[:], []string{locTypeNameL, locTypeNameM, locTypeNameS}...), parentsHeader[:]...)

	header, _ = NewImportHeader(locationTypeInOrder, ImportEntityEquipment)
	err = importer.inputValidations(ctx, header)
	require.NoError(t, err)
	require.EqualValues(t, ic.indexToLocationTypeID, map[int]int{
		4: locTypeIDS.locTypeIDL,
		5: locTypeIDS.locTypeIDM,
		6: locTypeIDS.locTypeIDS,
	})
}

func TestTitleEquipmentTypeInputValidation(t *testing.T) {
	r := newImporterTestResolver(t)
	importer := r.importer
	defer r.drv.Close()

	ctx := newImportContext(viewertest.NewContext(r.client))
	ic := getImportContext(ctx)
	var (
		equipDataHeader = [...]string{"Equipment ID", "Equipment Name", "Equipment Type", "External ID"}
		parentsHeader   = [...]string{"Parent Equipment (3)", "Position (3)", "Parent Equipment (2)", "Position (2)", "Parent Equipment", "Equipment Position"}
	)
	locationTypeInOrder := append(append(equipDataHeader[:], []string{locTypeNameL, locTypeNameM, locTypeNameS}...), parentsHeader[:]...)
	titleWithProperties := append(locationTypeInOrder, propName1, propName2, propName3, propName4)

	ids := prepareEquipmentTypeData(ctx, t, *r)
	/*
		populating:
		equipmentTypeNameToID
		propNameToIndex
		equipmentTypeIDToProperties
	*/
	header, _ := NewImportHeader(titleWithProperties, ImportEntityEquipment)
	err := importer.populateEquipmentTypeNameToIDMap(ctx, header, true)
	require.NoError(t, err)
	require.EqualValues(t, ic.equipmentTypeNameToID, map[string]int{
		equipmentTypeName:  ids.equipTypeID,
		equipmentType2Name: ids.equipTypeID2,
		equipmentType3Name: ids.equipTypeID3,
	})
	require.EqualValues(t, ic.propNameToIndex, map[string]int{
		propName1: 13,
		propName2: 14,
		propName3: 15,
		propName4: 16,
	})
	require.EqualValues(t, ic.equipmentTypeIDToProperties[ic.equipmentTypeNameToID[equipmentTypeName]], []string{
		propName1,
		propName2,
	})
	require.EqualValues(t, ic.equipmentTypeIDToProperties[ic.equipmentTypeNameToID[equipmentType2Name]], []string{
		propName3,
		propName4,
	})
	require.EqualValues(t, ic.equipmentTypeIDToProperties[ic.equipmentTypeNameToID[equipmentType3Name]], []string{
		propName5,
		propName6,
	})
}

func TestLocationHierarchy(t *testing.T) {
	r := newImporterTestResolver(t)
	importer := r.importer
	defer r.drv.Close()
	ctx := newImportContext(viewertest.NewContext(r.client))

	ids := prepareEquipmentTypeData(ctx, t, *r)

	var (
		equipDataHeader = [...]string{"Equipment ID", "Equipment Name", "Equipment Type", "External ID"}
		parentsHeader   = [...]string{"Parent Equipment (3)", "Position (3)", "Parent Equipment (2)", "Position (2)", "Parent Equipment", "Equipment Position"}
		test1           = []string{"", "", equipmentTypeName, "1", "locNameL", "", "", "", "", "", ""}
		test2           = []string{"", "", equipmentTypeName, "1", "locNameL", "locNameM", "locNameS", "", "", "", ""}
		test3           = []string{"", "", equipmentTypeName, "1", "", "locNameM", "", "", "", "", ""}
	)
	locationTypeInOrder := append(append(equipDataHeader[:], []string{locTypeNameL, locTypeNameM, locTypeNameS}...), parentsHeader[:]...)
	title, _ := NewImportHeader(locationTypeInOrder, ImportEntityEquipment)
	err := importer.inputValidations(ctx, title)
	require.NoError(t, err)

	rec, _ := NewImportRecord(test1, title)
	loc, err := importer.verifyOrCreateLocationHierarchy(ctx, rec, true, nil)
	require.NoError(t, err)
	require.Equal(t, loc.Name, "locNameL")
	require.Equal(t, loc.QueryType().OnlyXID(ctx), ids.locTypeIDL)
	require.False(t, loc.QueryChildren().ExistX(ctx))

	rec2, _ := NewImportRecord(test2, title)
	loc2, err := importer.verifyOrCreateLocationHierarchy(ctx, rec2, true, nil)
	require.NoError(t, err)
	require.Equal(t, loc2.Name, "locNameS")
	require.Equal(t, loc2.QueryType().OnlyXID(ctx), ids.locTypeIDS)
	require.Equal(t, loc2.QueryParent().OnlyX(ctx).Name, "locNameM")

	rec3, _ := NewImportRecord(test3, title)
	loc3, err := importer.verifyOrCreateLocationHierarchy(ctx, rec3, true, nil)
	require.NoError(t, err)
	require.Equal(t, loc3.Name, "locNameM")
	require.Equal(t, loc3.QueryType().OnlyXID(ctx), ids.locTypeIDM)
	require.False(t, loc3.QueryChildren().ExistX(ctx))
}

func TestPosition(t *testing.T) {
	r := newImporterTestResolver(t)
	importer := r.importer
	defer r.drv.Close()
	ctx := newImportContext(viewertest.NewContext(r.client))
	prepareEquipmentTypeData(ctx, t, *r)

	pos1 := models.EquipmentPositionInput{
		Name: "pos1",
	}
	var (
		equipDataHeader = [...]string{"Equipment ID", "Equipment Name", "Equipment Type", "External ID"}
		parentsHeader   = [...]string{"Parent Equipment (3)", "Position (3)", "Parent Equipment (2)", "Position (2)", "Parent Equipment", "Equipment Position"}
		locCreate       = []string{"", "", equipmentTypeName, "1", "locNameL", "locNameM", "", "", "", "", "", "", ""}
		test1           = []string{"", "test", "type1", "1", "locNameL", "locNameM", "", "", "", "", "", "", "pos1"}
		test2           = []string{"", "test", "type1", "1", "locNameL", "locNameM", "", "", "", "", "", "equip1", ""}
		test3           = []string{"", "test", "type1", "1", "locNameL", "locNameM", "", "", "", "equip1", "", "", "pos1"}
		test4           = []string{"", "test", "type1", "1", "locNameL", "locNameM", "", "", "", "", "", "equip1", "pos1"}
	)
	locationTypeInOrder := append(append(equipDataHeader[:], []string{locTypeNameL, locTypeNameM, locTypeNameS}...), parentsHeader[:]...)
	title, _ := NewImportHeader(locationTypeInOrder, ImportEntityEquipment)
	err := importer.inputValidations(ctx, title)
	require.NoError(t, err)

	rec, _ := NewImportRecord(locCreate, title)
	loc, err := importer.verifyOrCreateLocationHierarchy(ctx, rec, true, nil)
	require.NoError(t, err)
	equipmentType, err := importer.r.Mutation().AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:      "type1",
		Positions: []*models.EquipmentPositionInput{&pos1},
	})
	require.NoError(t, err)
	equip, err := importer.r.Mutation().AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "equip1",
		Type:     equipmentType.ID,
		Location: &loc.ID,
	})
	require.NoError(t, err)

	rec1, _ := NewImportRecord(test1, title)
	_, _, err = importer.getPositionDetailsIfExists(ctx, loc, rec1, true)
	require.Error(t, err)

	rec2, _ := NewImportRecord(test2, title)
	eq, def, err := importer.getPositionDetailsIfExists(ctx, loc, rec2, true)
	require.Nil(t, err)
	require.Nil(t, eq)
	require.Nil(t, def)

	rec3, _ := NewImportRecord(test3, title)
	_, _, err = importer.getPositionDetailsIfExists(ctx, loc, rec3, true)
	require.NoError(t, err)

	rec4, _ := NewImportRecord(test4, title)
	equipID, defID, err := importer.getPositionDetailsIfExists(ctx, loc, rec4, true)
	require.NoError(t, err)
	fetchedEquip, err := importer.r.Query().Equipment(ctx, *equipID)
	require.NoError(t, err)
	fetchedDefinition := equipmentType.QueryPositionDefinitions().OnlyX(ctx)
	require.Equal(t, fetchedEquip.Name, equip.Name)
	require.Equal(t, fetchedEquip.ID, equip.ID)
	require.Equal(t, fetchedDefinition.ID, *defID)
}

func TestValidatePropertiesForType(t *testing.T) {
	r := newImporterTestResolver(t)
	importer := r.importer
	q := r.importer.r.Query()
	defer r.drv.Close()
	ctx := newImportContext(viewertest.NewContext(r.client))
	data := prepareEquipmentTypeData(ctx, t, *r)

	var (
		equipDataHeader = [...]string{"Equipment ID", "Equipment Name", "Equipment Type", "External ID"}
		parentsHeader   = [...]string{"Parent Equipment (3)", "Position (3)", "Parent Equipment (2)", "Position (2)", "Parent Equipment", "Equipment Position"}
		row1            = []string{"", "e1", equipmentTypeName, "1id", "locNameL", "locNameM", "", "", "", "", "", "", "", "strVal", "54", "", "", "", ""}
		row2            = []string{"", "e2", equipmentType2Name, "1id", "locNameL", "locNameM", "", "", "", "", "", "", "", "", "", "29/03/88", "false", "", ""}
		row3            = []string{"", "e3", equipmentType3Name, "1id", "locNameL", "locNameM", "", "", "", "", "", "", "", "", "", "", "", "30.23-50", "45.8,88.9"}
	)

	locationTypeInOrder := append(append(equipDataHeader[:], []string{locTypeNameL, locTypeNameM, locTypeNameS}...), parentsHeader[:]...)
	finalFirstRow := append(locationTypeInOrder, propName1, propName2, propName3, propName4, propName5, propName6)
	fl, _ := NewImportHeader(locationTypeInOrder, ImportEntityEquipment)
	err := importer.inputValidations(ctx, fl)
	require.NoError(t, err)

	fl, _ = NewImportHeader(finalFirstRow, ImportEntityEquipment)
	err = importer.populateEquipmentTypeNameToIDMap(ctx, fl, true)
	r1, _ := NewImportRecord(row1, fl)
	require.NoError(t, err)
	etyp1, err := q.EquipmentType(ctx, data.equipTypeID)
	require.NoError(t, err)
	ptypes, err := importer.validatePropertiesForEquipmentType(ctx, r1, etyp1)
	require.NoError(t, err)
	require.Len(t, ptypes, 2)
	require.NotEqual(t, ptypes[0].PropertyTypeID, ptypes[1].PropertyTypeID)
	for _, value := range ptypes {
		ptyp := etyp1.QueryPropertyTypes().Where(propertytype.ID(value.PropertyTypeID)).OnlyX(ctx)
		switch ptyp.Name {
		case propName1:
			require.Equal(t, *value.StringValue, "strVal")
			require.Equal(t, ptyp.Type, "string")
		case propName2:
			require.Equal(t, *value.IntValue, 54)
			require.Equal(t, ptyp.Type, "int")
		default:
			require.Fail(t, "property type name should be one of the two")
		}
	}
	etyp2, err := q.EquipmentType(ctx, data.equipTypeID2)
	require.NoError(t, err)

	r2, _ := NewImportRecord(row2, fl)
	ptypes2, err := importer.validatePropertiesForEquipmentType(ctx, r2, etyp2)
	require.NoError(t, err)
	require.Len(t, ptypes2, 2)
	for _, value := range ptypes2 {
		ptyp := etyp2.QueryPropertyTypes().Where(propertytype.ID(value.PropertyTypeID)).OnlyX(ctx)
		switch ptyp.Name {
		case propName3:
			require.Equal(t, *value.StringValue, "29/03/88")
			require.Equal(t, ptyp.Type, "date")
		case propName4:
			require.Equal(t, *value.BooleanValue, false)
			require.Equal(t, ptyp.Type, "bool")
		default:
			require.Fail(t, "property type name should be one of the two")
		}
	}

	etyp3, err := q.EquipmentType(ctx, data.equipTypeID3)
	require.NoError(t, err)

	r3, _ := NewImportRecord(row3, fl)
	ptypes3, err := importer.validatePropertiesForEquipmentType(ctx, r3, etyp3)
	require.NoError(t, err)
	require.Len(t, ptypes3, 2)
	require.NotEqual(t, ptypes3[0].PropertyTypeID, ptypes3[1].PropertyTypeID)
	for _, value := range ptypes3 {
		ptyp := etyp3.QueryPropertyTypes().Where(propertytype.ID(value.PropertyTypeID)).OnlyX(ctx)
		switch ptyp.Name {
		case propName5:
			require.Equal(t, *value.RangeFromValue, 30.23)
			require.EqualValues(t, *value.RangeToValue, 50)
			require.Equal(t, ptyp.Type, "range")
		case ptyp.Name:
			require.Equal(t, *value.LatitudeValue, 45.8)
			require.Equal(t, *value.LongitudeValue, 88.9)
			require.Equal(t, ptyp.Type, "gps_location")
		default:
			require.Fail(t, "property type name should be one of the two")
		}
	}
}

func TestValidateForExistingEquipment(t *testing.T) {
	r := newImporterTestResolver(t)
	importer := r.importer
	defer r.drv.Close()
	ctx := newImportContext(viewertest.NewContext(r.client))
	prepareEquipmentTypeData(ctx, t, *r)

	pos1 := models.EquipmentPositionInput{
		Name: "pos1",
	}
	pos2 := models.EquipmentPositionInput{
		Name: "pos2",
	}
	var (
		equipDataHeader = [...]string{"Equipment ID", "Equipment Name", "Equipment Type", "External ID"}
		parentsHeader   = [...]string{"Parent Equipment (3)", "Position (3)", "Parent Equipment (2)", "Position (2)", "Parent Equipment", "Equipment Position"}
		locCreate       = []string{"", "", equipmentTypeName, "1id", "locNameL", "locNameM", "", "", "", "", ""}
	)
	locationTypeInOrder := append(append(equipDataHeader[:], []string{locTypeNameL, locTypeNameM, locTypeNameS}...), parentsHeader[:]...)
	title, _ := NewImportHeader(locationTypeInOrder, ImportEntityEquipment)
	err := importer.inputValidations(ctx, title)
	require.NoError(t, err)

	rec, _ := NewImportRecord(locCreate, title)
	loc, err := importer.verifyOrCreateLocationHierarchy(ctx, rec, true, nil)
	require.NoError(t, err)
	equipmentType, err := importer.r.Mutation().AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:      "type1",
		Positions: []*models.EquipmentPositionInput{&pos1, &pos2},
	})
	require.NoError(t, err)
	parent, err := importer.r.Mutation().AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "parent",
		Type:     equipmentType.ID,
		Location: &loc.ID,
	})
	require.NoError(t, err)
	posDefs := equipmentType.QueryPositionDefinitions().AllX(ctx)
	child, err := importer.r.Mutation().AddEquipment(ctx, models.AddEquipmentInput{
		Name:               "child",
		Type:               equipmentType.ID,
		Parent:             &parent.ID,
		PositionDefinition: &posDefs[0].ID,
	})
	require.NoError(t, err)
	grandchild, err := importer.r.Mutation().AddEquipment(ctx, models.AddEquipmentInput{
		Name:               "grandchild",
		Type:               equipmentType.ID,
		Parent:             &child.ID,
		PositionDefinition: &posDefs[1].ID,
	})
	require.NoError(t, err)
	var (
		test1 = []string{strconv.Itoa(child.ID), "c_new_name", "type1", "1id", "locNameL", "locNameM", "", "", "", "", "", "parent", "pos1"}
		test2 = []string{strconv.Itoa(grandchild.ID), "gc_new_name", "type1", "1id", "locNameL", "locNameM", "", "", "", "parent", "pos1", "child", "pos2"}
	)

	rec1, _ := NewImportRecord(test1, title)
	_, err = importer.validateLineForExistingEquipment(ctx, child.ID, rec1)
	require.NoError(t, err)

	rec2, _ := NewImportRecord(test2, title)
	_, err = importer.validateLineForExistingEquipment(ctx, grandchild.ID, rec2)
	require.NoError(t, err)
}
