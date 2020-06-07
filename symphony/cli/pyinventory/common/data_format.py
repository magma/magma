#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import List, Sequence

from ..graphql.fragment.property_type import PropertyTypeFragment
from ..graphql.input.property_type import PropertyTypeInput
from .data_class import PropertyDefinition


def format_to_property_definition(
    property_type_fragment: PropertyTypeFragment,
) -> PropertyDefinition:
    """This function gets `pyinventory.graphql.fragment.property_type.PropertyTypeFragment` object as argument and formats it to `pyinventory.common.data_class.PropertyDefinition` object

        Args:
            property_type_fragment ( `pyinventory.graphql.fragment.property_type.PropertyTypeFragment` ): existing property type fragment object

        Returns:
            `pyinventory.common.data_class.PropertyDefinition` object

        Example:
            ```
            property_definition = format_to_property_definition(
                property_type_fragment=property_type_fragment,
            )
            ```
    """
    return PropertyDefinition(
        id=property_type_fragment.id,
        property_name=property_type_fragment.name,
        property_kind=property_type_fragment.type,
        default_raw_value=property_type_fragment.rawValue,
        is_fixed=not property_type_fragment.isInstanceProperty,
        external_id=property_type_fragment.externalId,
        is_mandatory=property_type_fragment.isMandatory,
        is_deleted=property_type_fragment.isDeleted,
    )


def format_to_property_definitions(
    data: Sequence[PropertyTypeFragment],
) -> Sequence[PropertyDefinition]:
    """This function gets Sequence[ `pyinventory.graphql.fragment.property_type.PropertyTypeFragment` ] as argument and formats it to Sequence[ `pyinventory.common.data_class.PropertyDefinition` ]

        Args:
            data (Sequence[ `pyinventory.graphql.fragment.property_type.PropertyTypeFragment` ]): existing property type fragments sequence

        Returns:
            Sequence[ `pyinventory.common.data_class.PropertyDefinition` ]

        Example:
            ```
            property_definitions = format_to_property_definitions(
                data=property_type_fragments,
            )
            ```
    """
    return [
        format_to_property_definition(property_type_fragment)
        for property_type_fragment in data
    ]


def format_to_property_type_input(
    property_definition: PropertyDefinition, is_new: bool = True
) -> PropertyTypeInput:
    """This function gets `pyinventory.common.data_class.PropertyDefinition` object as argument and formats it to `pyinventory.graphql.input.property_type.PropertyTypeInput` object

        Args:
            property_definition ( `pyinventory.graphql.input.property_type.PropertyTypeInput` ): existing property definition object

        Returns:
            `pyinventory.graphql.input.property_type.PropertyTypeInput` object

        Example:
            ```
            property_type_input = format_to_property_type_input(
                property_definition=property_definition,
            )
            ```
    """
    string_value = None
    int_value = None
    boolean_value = None
    float_value = None
    latitude_value = None
    longitude_value = None
    range_from_value = None
    range_to_value = None

    kind = property_definition.property_kind.value

    if property_definition.default_raw_value is not None:
        default_raw_value = property_definition.default_raw_value
        if kind == "int":
            int_value = int(default_raw_value, base=10)
        elif kind == "bool":
            boolean_value = True if default_raw_value.lower() == "true" else False
        elif kind == "float":
            float_value = float(default_raw_value)
        elif kind == "range":
            string_range = default_raw_value.split(" - ")
            range_from_value = float(string_range[0])
            range_to_value = float(string_range[1])
        elif kind == "gps_location":
            string_coordinates = default_raw_value.split(", ")
            latitude_value = float(string_coordinates[0])
            longitude_value = float(string_coordinates[1])
        else:
            string_value = default_raw_value

    return PropertyTypeInput(
        id=property_definition.id if not is_new else None,
        name=property_definition.property_name,
        type=property_definition.property_kind,
        externalId=property_definition.external_id
        if not is_new and property_definition.external_id
        else None,
        stringValue=string_value,
        intValue=int_value,
        booleanValue=boolean_value,
        floatValue=float_value,
        latitudeValue=latitude_value,
        longitudeValue=longitude_value,
        rangeFromValue=range_from_value,
        rangeToValue=range_to_value,
        isInstanceProperty=not property_definition.is_fixed,
        isMandatory=property_definition.is_mandatory,
        isDeleted=property_definition.is_deleted,
    )


def format_to_property_type_inputs(
    data: Sequence[PropertyDefinition],
) -> List[PropertyTypeInput]:
    """This function gets Sequence[ `pyinventory.common.data_class.PropertyDefinition` ] as argument and formats it to Sequence[ `pyinventory.graphql.input.property_type.PropertyTypeInput` ]

        Args:
            data (Sequence[ `pyinventory.common.data_class.PropertyDefinition` ]): existing property definitions sequence

        Returns:
            Sequence[ `pyinventory.graphql.input.property_type.PropertyTypeInput` ]

        Example:
            ```
            property_type_inputs = format_to_property_type_inputs(
                data=property_type_definitions,
            )
            ```
    """
    return [
        format_to_property_type_input(property_definition)
        for property_definition in data
    ]
