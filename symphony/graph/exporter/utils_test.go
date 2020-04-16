// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exporter

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/stretchr/testify/require"
)

const strVal = "defVal"

func TestLocationHierarchy(t *testing.T) {
	r := newExporterTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr := r.Mutation()
	locTypeL, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: "example_type_large"})
	require.NoError(t, err)
	locTypeM, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: "example_type_medium"})
	require.NoError(t, err)
	locTypeS, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: "example_type_small"})
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
	client := ent.FromContext(ctx)
	locTypeHierarchy, err := locationTypeHierarchy(ctx, client)
	require.NoError(t, err)

	require.Equal(t, locTypeHierarchy[0], locTypeL.Name)
	require.Equal(t, locTypeHierarchy[1], locTypeM.Name)
	require.Equal(t, locTypeHierarchy[2], locTypeS.Name)

	gpLocation, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "grand_parent_loc",
		Type: locTypeL.ID,
	})
	require.NoError(t, err)
	pLocation, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name:   "parent_loc",
		Type:   locTypeM.ID,
		Parent: &gpLocation.ID,
	})
	require.NoError(t, err)
	clocation, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name:   "child_loc",
		Type:   locTypeS.ID,
		Parent: &pLocation.ID,
	})
	require.NoError(t, err)

	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "equipment_type",
	})
	require.NoError(t, err)
	equipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "child_equipment",
		Type:     equipmentType.ID,
		Location: &clocation.ID,
	})
	require.NoError(t, err)

	locHierarchy, err := locationHierarchyForEquipment(ctx, equipment, locTypeHierarchy)
	require.NoError(t, err)

	require.Equal(t, locHierarchy[0], gpLocation.Name)
	require.Equal(t, locHierarchy[1], pLocation.Name)
	require.Equal(t, locHierarchy[2], clocation.Name)
}

func TestParentHierarchy(t *testing.T) {
	r := newExporterTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr := r.Mutation()
	mapType := "map"
	mapZoomLvl := 7
	locType, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name:         "example_loc_type",
		MapType:      &mapType,
		MapZoomLevel: &mapZoomLvl,
	})
	require.NoError(t, err)

	location, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "example_loc_inst",
		Type: locType.ID,
	})
	require.NoError(t, err)

	position1 := models.EquipmentPositionInput{
		Name: "Position 1",
	}
	position2 := models.EquipmentPositionInput{
		Name: "Position 2",
	}

	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:      "equipment_type",
		Positions: []*models.EquipmentPositionInput{&position1, &position2},
	})
	require.NoError(t, err)
	posDef1 := equipmentType.QueryPositionDefinitions().Where(equipmentpositiondefinition.Name("Position 1")).OnlyX(ctx)
	posDef2 := equipmentType.QueryPositionDefinitions().Where(equipmentpositiondefinition.Name("Position 2")).OnlyX(ctx)

	grandParentEquipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "grand_parent_equipment",
		Type:     equipmentType.ID,
		Location: &location.ID,
	})
	require.NoError(t, err)

	parentEquipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:               "parent_equipment",
		Type:               equipmentType.ID,
		Parent:             &grandParentEquipment.ID,
		PositionDefinition: &posDef1.ID,
	})
	require.NoError(t, err)

	childEquipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:               "child_equipment",
		Type:               equipmentType.ID,
		Parent:             &parentEquipment.ID,
		PositionDefinition: &posDef2.ID,
	})
	require.NoError(t, err)

	hierarchy := parentHierarchyWithAllPositions(ctx, *childEquipment)
	require.NoError(t, err)

	require.Equal(t, hierarchy[0], "")
	require.Equal(t, hierarchy[1], "")
	require.Equal(t, hierarchy[2], grandParentEquipment.Name)
	require.Equal(t, hierarchy[3], posDef1.Name)
	require.Equal(t, hierarchy[4], parentEquipment.Name)
	require.Equal(t, hierarchy[5], posDef2.Name)

}

func TestPropertiesForCSV(t *testing.T) {
	r := newExporterTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	client := ent.FromContext(ctx)

	mr := r.Mutation()
	locType, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: "example_loc_type"})
	require.NoError(t, err)

	location, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "example_loc_inst",
		Type: locType.ID,
	})
	require.NoError(t, err)

	propInput1 := models.PropertyTypeInput{
		Name: "Property type1",
		Type: "int",
	}
	propInput2 := models.PropertyTypeInput{
		Name: "Property type2",
		Type: "string",
	}
	propInput3 := models.PropertyTypeInput{
		Name: "Property type3",
		Type: "gps_location",
	}
	propInput4 := models.PropertyTypeInput{
		Name: "Property type4",
		Type: "range",
	}
	propInput5 := models.PropertyTypeInput{
		Name: "Property type5",
		Type: "bool",
	}
	propInput6 := models.PropertyTypeInput{
		Name: "Property type6",
		Type: "node",
	}
	propInput7 := models.PropertyTypeInput{
		Name: "Property type7",
		Type: "node",
	}

	propInput8 := models.PropertyTypeInput{
		Name: "Property type8",
		Type: "node",
	}

	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "equipment_type",
		Properties: []*models.PropertyTypeInput{
			&propInput1, &propInput2, &propInput3, &propInput4, &propInput5, &propInput6, &propInput7, &propInput8,
		},
	})
	require.NoError(t, err)
	propType1 := equipmentType.QueryPropertyTypes().Where(propertytype.Name("Property type1")).OnlyX(ctx)
	propType2 := equipmentType.QueryPropertyTypes().Where(propertytype.Name("Property type2")).OnlyX(ctx)
	propType3 := equipmentType.QueryPropertyTypes().Where(propertytype.Name("Property type3")).OnlyX(ctx)
	propType4 := equipmentType.QueryPropertyTypes().Where(propertytype.Name("Property type4")).OnlyX(ctx)
	propType5 := equipmentType.QueryPropertyTypes().Where(propertytype.Name("Property type5")).OnlyX(ctx)
	propType6 := equipmentType.QueryPropertyTypes().Where(propertytype.Name("Property type6")).OnlyX(ctx)
	propType7 := equipmentType.QueryPropertyTypes().Where(propertytype.Name("Property type7")).OnlyX(ctx)
	propType8 := equipmentType.QueryPropertyTypes().Where(propertytype.Name("Property type8")).OnlyX(ctx)

	intVal := 40
	strVal := strVal
	prop1 := models.PropertyInput{
		PropertyTypeID: propType1.ID,
		IntValue:       &intVal,
	}
	prop2 := models.PropertyInput{
		PropertyTypeID: propType2.ID,
		StringValue:    &strVal,
	}
	latVal := 40.32
	longVal := 40.34
	prop3 := models.PropertyInput{
		PropertyTypeID: propType3.ID,
		LatitudeValue:  &latVal,
		LongitudeValue: &longVal,
	}
	coords := fmt.Sprintf("%f", latVal) + ", " + fmt.Sprintf("%f", longVal)

	fromVal := 10.0
	toVal := 20.0
	prop4 := models.PropertyInput{
		PropertyTypeID: propType4.ID,
		RangeFromValue: &fromVal,
		RangeToValue:   &toVal,
	}
	rangeVal := fmt.Sprintf("%.3f", fromVal) + " - " + fmt.Sprintf("%.3f", toVal)

	boolVal := true
	prop5 := models.PropertyInput{
		PropertyTypeID: propType5.ID,
		BooleanValue:   &boolVal,
	}

	propEquipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "prop_equipment_type",
	})
	require.NoError(t, err)
	propEquipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "prop_equipment",
		Type:     propEquipmentType.ID,
		Location: &location.ID,
	})
	require.NoError(t, err)
	prop6 := models.PropertyInput{
		PropertyTypeID: propType6.ID,
		NodeIDValue:    &propEquipment.ID,
	}

	propLocationType, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: "prop_loc_type"})
	require.NoError(t, err)
	propLocation, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "prop_loc_inst",
		Type: propLocationType.ID,
	})
	require.NoError(t, err)
	prop7 := models.PropertyInput{
		PropertyTypeID: propType7.ID,
		NodeIDValue:    &propLocation.ID,
	}

	propServiceType, err := mr.AddServiceType(ctx, models.ServiceTypeCreateData{Name: "prop_service_type", HasCustomer: false})
	require.NoError(t, err)
	propService, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "prop_service_inst",
		ServiceTypeID: propServiceType.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)
	prop8 := models.PropertyInput{
		PropertyTypeID: propType8.ID,
		NodeIDValue:    &propService.ID,
	}

	equipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:       "child_equipment",
		Type:       equipmentType.ID,
		Location:   &location.ID,
		Properties: []*models.PropertyInput{&prop1, &prop2, &prop3, &prop4, &prop5, &prop6, &prop7, &prop8},
	})
	require.NoError(t, err)

	propertyTypes, err := propertyTypesSlice(ctx, []int{equipment.ID}, client, models.PropertyEntityEquipment)
	require.NoError(t, err)

	props, err := propertiesSlice(ctx, equipment, propertyTypes, models.PropertyEntityEquipment)
	require.NoError(t, err)
	require.Contains(t, props, strVal)
	require.Contains(t, props, strconv.Itoa(intVal))
	require.Contains(t, props, coords)
	require.Contains(t, props, rangeVal)
	require.Contains(t, props, strconv.FormatBool(boolVal))
}

func TestPropertyTypesForCSV(t *testing.T) {
	r := newExporterTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	client := ent.FromContext(ctx)

	mr := r.Mutation()
	locType, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: "example_loc_type"})
	require.NoError(t, err)

	location, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "example_loc_inst",
		Type: locType.ID,
	})
	require.NoError(t, err)

	strVal := strVal
	intVal := 40
	propInput1 := models.PropertyTypeInput{
		Name:     "Property type1",
		Type:     "int",
		IntValue: &intVal,
	}
	propInput2 := models.PropertyTypeInput{
		Name:        "Property type2",
		Type:        "string",
		StringValue: &strVal,
	}
	latVal := 40.32
	longVal := 40.34
	propInput3 := models.PropertyTypeInput{
		Name:           "Property type3",
		Type:           "gps_location",
		LatitudeValue:  &latVal,
		LongitudeValue: &longVal,
	}
	coords := fmt.Sprintf("%f", latVal) + ", " + fmt.Sprintf("%f", longVal)

	fromVal := 10.0
	toVal := 20.0
	propInput4 := models.PropertyTypeInput{
		Name:           "Property type4",
		Type:           "range",
		RangeFromValue: &fromVal,
		RangeToValue:   &toVal,
	}
	rangeVal := fmt.Sprintf("%.3f", fromVal) + " - " + fmt.Sprintf("%.3f", toVal)

	boolVal := true
	propInput5 := models.PropertyTypeInput{
		Name:         "Property type5",
		Type:         "bool",
		BooleanValue: &boolVal,
	}

	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "equipment_type",
		Properties: []*models.PropertyTypeInput{
			&propInput1, &propInput2, &propInput3, &propInput4, &propInput5,
		},
	})
	require.NoError(t, err)

	equipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "child_equipment",
		Type:     equipmentType.ID,
		Location: &location.ID,
	})
	require.NoError(t, err)

	propertyTypes, err := propertyTypesSlice(ctx, []int{equipment.ID}, client, models.PropertyEntityEquipment)
	require.NoError(t, err)

	props, err := propertiesSlice(ctx, equipment, propertyTypes, models.PropertyEntityEquipment)
	require.NoError(t, err)
	require.Contains(t, props, strVal)
	require.Contains(t, props, strconv.Itoa(intVal))
	require.Contains(t, props, coords)
	require.Contains(t, props, rangeVal)
	require.Contains(t, props, strconv.FormatBool(boolVal))
}

func TestSamePropertyTypesForCSV(t *testing.T) {
	r := newExporterTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	client := ent.FromContext(ctx)

	mr := r.Mutation()
	locType, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: "example_loc_type"})
	require.NoError(t, err)

	location, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "example_loc_inst",
		Type: locType.ID,
	})
	require.NoError(t, err)

	intVal := 40
	propInput1 := models.PropertyTypeInput{
		Name:     "Ptype1",
		Type:     models.PropertyKindInt,
		IntValue: &intVal,
	}

	equipmentTypeA, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:       "equipment_typeA",
		Properties: []*models.PropertyTypeInput{&propInput1},
	})
	require.NoError(t, err)
	pa := equipmentTypeA.QueryPropertyTypes().OnlyX(ctx)
	propInput2 := models.PropertyTypeInput{
		Name:     "Ptype2",
		Type:     models.PropertyKindInt,
		IntValue: &intVal,
	}
	equipmentTypeB, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:       "equipment_typeB",
		Properties: []*models.PropertyTypeInput{&propInput1, &propInput2},
	})
	require.NoError(t, err)
	equipmentTypeB.QueryPropertyTypes().Where(propertytype.Name("Ptype2")).OnlyX(ctx)

	equ, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "child_equipment",
		Type:     equipmentTypeA.ID,
		Location: &location.ID,
	})
	require.NoError(t, err)

	propertyTypes, err := propertyTypesSlice(ctx, []int{equ.ID}, client, models.PropertyEntityEquipment)
	require.Len(t, propertyTypes, 1)
	require.Contains(t, propertyTypes, pa.Name)
	require.NoError(t, err)
}
