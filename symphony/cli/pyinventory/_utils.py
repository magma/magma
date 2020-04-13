#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import warnings
from datetime import datetime
from typing import Callable, Dict, List, Mapping, Optional, Sequence, Tuple, Union, cast

from dacite import Config, from_dict

from .consts import (
    TYPE_AND_FIELD_NAME,
    DataTypeName,
    Entity,
    PropertyDefinition,
    PropertyValue,
    ReturnType,
)
from .exceptions import EntityNotFoundError
from .graphql.equipment_port_definition_fragment import EquipmentPortDefinitionFragment
from .graphql.equipment_port_input import EquipmentPortInput
from .graphql.equipment_position_definition_fragment import (
    EquipmentPositionDefinitionFragment,
)
from .graphql.equipment_position_input import EquipmentPositionInput
from .graphql.property_fragment import PropertyFragment
from .graphql.property_input import PropertyInput
from .graphql.property_kind_enum import PropertyKind
from .graphql.property_type_fragment import PropertyTypeFragment
from .graphql.property_type_input import PropertyTypeInput


def format_to_type_and_field_name(type_key: str) -> Optional[DataTypeName]:
    formated = TYPE_AND_FIELD_NAME.get(type_key, None)
    return formated


def get_graphql_input_field(
    property_type_name: str, type_key: str, value: PropertyValue
) -> Dict[str, PropertyValue]:
    formated_type = format_to_type_and_field_name(type_key)
    if formated_type is None:
        raise Exception(
            f"property type {property_type_name} has not supported type {type_key}"
        )
    if type_key == "string":
        assert isinstance(value, str) or isinstance(
            value, bytes
        ), f"property {property_type_name} is not of type {type_key}"
    elif type_key == "gps_location":
        assert isinstance(
            value, tuple
        ), f"property {property_type_name} is not of type {type_key}"
        gps_value = value
        assert (
            len(gps_value) == 2
            and isinstance(gps_value[0], float)
            and isinstance(gps_value[1], float)
        ), f"property {property_type_name} is not of type tuple(float, float)"
        return {
            formated_type.graphql_field_name[0]: gps_value[0],
            formated_type.graphql_field_name[1]: gps_value[1],
        }
    else:
        assert isinstance(
            value, formated_type.data_type
        ), f"property {property_type_name} is not of type {type_key}"

    return {formated_type.graphql_field_name[0]: cast(PropertyValue, value)}


def get_graphql_property_type_inputs(
    property_types: Sequence[PropertyTypeFragment],
    properties_dict: Dict[str, PropertyValue],
) -> List[PropertyTypeInput]:
    """This function gets existing property types and dictionary, where key - are type names, and keys - new values
    formats data, validates existence of keys from properties_dict in property_types and returns list of PropertyTypeInput
 
        Args:
            property_types (List[pyinventory.graphql.property_type_fragment.PropertyTypeFragment]): list of existing property types
            properties_dict (Dict[str, pyinventory.consts.PropertyValue]): dictionary of properties, where
                str: name of existing property
                PropertyValue: new value of existing type for this property
       
        Returns:
            List[pyinventory.graphql.property_type_input.PropertyTypeInput]
 
        Raises:
            EntityNotFoundError if there any unknown property name in properties_dict keys
    """
    properties = []
    property_type_names = {}

    for property_type in property_types:
        property_type_names[property_type.name] = property_type

    for name, value in properties_dict.items():
        if name not in property_type_names:
            raise EntityNotFoundError(entity=Entity.PropertyType, entity_name=name)
        assert property_type_names[
            name
        ].isInstanceProperty, f"property {name} is not instance property"
        result = {
            "id": property_type_names[name].id,
            "name": name,
            "type": PropertyKind(property_type_names[name].type),
        }
        result.update(
            get_graphql_input_field(
                property_type_name=name,
                type_key=property_type_names[name].type.value,
                value=value,
            )
        )
        properties.append(
            from_dict(
                data_class=PropertyTypeInput, data=result, config=Config(strict=True)
            )
        )

    return properties


def get_graphql_property_inputs(
    property_types: Sequence[PropertyTypeFragment],
    properties_dict: Mapping[str, PropertyValue],
) -> List[PropertyInput]:
    """This function gets existing property types and dictionary, where key - are type names, and keys - new values
    formats data, validates existence of keys from properties_dict in property_types and returns list of PropertyInput
 
        Args:
            property_types (Sequence[pyinventory.graphql.property_type_fragment.PropertyTypeFragment]): list of existing property types
            properties_dict (Mapping[str, pyinventory.consts.PropertyValue]): dictionary of properties, where
                str: name of existing property
                PropertyValue: new value of existing type for this property
       
        Returns:
            List[pyinventory.graphql.property_input.PropertyInput]
 
        Raises:
            EntityNotFoundError if there any unknown property name in properties_dict keys
       
        Example:
        ```
            property_types = client.locationTypes[location_type].property_types
            properties = get_graphql_property_inputs(property_types, properties_dict)
        ```
    """
    properties = []
    property_type_names = {}

    for property_type in property_types:
        property_type_names[property_type.name] = property_type

    for name, value in properties_dict.items():
        if name not in property_type_names:
            raise EntityNotFoundError(entity=Entity.PropertyType, entity_name=name)
        assert property_type_names[
            name
        ].isInstanceProperty, f"property {name} is not instance property"
        result = {"propertyTypeID": property_type_names[name].id}
        result.update(
            get_graphql_input_field(
                property_type_name=name,
                type_key=property_type_names[name].type.value,
                value=value,
            )
        )
        properties.append(
            from_dict(data_class=PropertyInput, data=result, config=Config(strict=True))
        )

    return properties


def _get_property_value(
    property_type: str, property: PropertyFragment
) -> Tuple[PropertyValue, ...]:
    formated_name = format_to_type_and_field_name(property_type)
    if formated_name is None:
        raise AssertionError(f"Unknown property type - {property_type}")

    str_fields = formated_name.graphql_field_name
    values = []
    for str_field in str_fields:
        if property_type == "date":
            date_data = property.__dict__[str_field]
            values.append(datetime.strptime(cast(str, date_data), "%Y-%m-%d").date())
        else:
            values.append(property.__dict__[str_field])
    return tuple(value for value in values)


def _get_property_default_value(
    name: str, type: str, value: Optional[PropertyValue]
) -> Dict[str, PropertyValue]:
    if value is None:
        return {}
    return get_graphql_input_field(property_type_name=name, type_key=type, value=value)


def _make_property_types(
    properties: Sequence[Tuple[str, str, Optional[PropertyValue], Optional[bool]]]
) -> List[PropertyTypeInput]:
    property_types = [
        from_dict(
            data_class=PropertyTypeInput,
            data={
                "name": arg[0],
                "type": PropertyKind(arg[1]),
                "index": i,
                **_get_property_default_value(arg[0], arg[1], arg[2]),
                "isInstanceProperty": arg[3],
            },
            config=Config(strict=True),
        )
        for i, arg in enumerate(properties)
    ]
    return property_types


# TODO(T63055378): remove
def property_type_to_kind(
    key: str, value: PropertyValue
) -> Union[PropertyValue, PropertyKind]:
    return value if key != "type" else PropertyKind(value)


# TODO(T63055378): remove and change usage to format_property_definitions
def format_properties(
    properties: Sequence[Tuple[str, str, Optional[PropertyValue], Optional[bool]]]
) -> List[PropertyTypeInput]:
    property_types = _make_property_types(properties)
    return property_types


def get_property_type_input(
    property_type: PropertyTypeFragment, is_new: bool = True
) -> PropertyTypeInput:
    return PropertyTypeInput(
        name=property_type.name,
        type=property_type.type,
        id=property_type.id if not is_new else None,
        externalId=property_type.externalId if not is_new else None,
        index=property_type.index,
        category=property_type.category,
        stringValue=property_type.stringValue,
        intValue=property_type.intValue,
        booleanValue=property_type.booleanValue,
        floatValue=property_type.floatValue,
        latitudeValue=property_type.latitudeValue,
        longitudeValue=property_type.longitudeValue,
        rangeFromValue=property_type.rangeFromValue,
        rangeToValue=property_type.rangeToValue,
        isEditable=property_type.isEditable,
        isInstanceProperty=property_type.isInstanceProperty,
        isMandatory=property_type.isMandatory,
        isDeleted=property_type.isDeleted,
    )


def get_position_definition_input(
    position_definition: EquipmentPositionDefinitionFragment, is_new: bool = True
) -> EquipmentPositionInput:
    return EquipmentPositionInput(
        name=position_definition.name,
        id=position_definition.id if not is_new else None,
        index=position_definition.index,
        visibleLabel=position_definition.visibleLabel,
    )


def get_port_definition_input(
    port_definition: EquipmentPortDefinitionFragment, is_new: bool = True
) -> EquipmentPortInput:
    return EquipmentPortInput(
        name=port_definition.name,
        id=port_definition.id if not is_new else None,
        index=port_definition.index,
        visibleLabel=port_definition.visibleLabel,
    )


def format_property_definitions(
    properties: List[PropertyDefinition]
) -> List[PropertyTypeInput]:
    property_types = [
        from_dict(
            data_class=PropertyTypeInput,
            data={
                "name": prop.property_name,
                "type": PropertyKind(prop.property_kind),
                "index": i,
                "externalId": prop.external_id,
                **_get_property_default_value(
                    prop.property_name, prop.property_kind.value, prop.default_value
                ),
                "isInstanceProperty": not prop.is_fixed,
            },
            config=Config(strict=True),
        )
        for i, prop in enumerate(properties)
    ]
    return property_types


def deprecated(
    deprecated_in: str, deprecated_by: str
) -> Callable[[Callable[..., ReturnType]], Callable[..., ReturnType]]:
    def wrapped(func: Callable[..., ReturnType]) -> Callable[..., ReturnType]:
        def wrapper(*args: str, **kwargs: int) -> Callable[..., ReturnType]:
            func_name = func.__name__
            message = f"{func_name} is deprecated in {deprecated_in}. Use the {deprecated_by} function instead."
            warnings.warn(message, DeprecationWarning, stacklevel=2)
            return cast(Callable[..., ReturnType], func(*args, **kwargs))

        return cast(Callable[..., ReturnType], wrapper)

    return wrapped
