#!/usr/bin/env python3
# pyre-strict

import warnings
from datetime import datetime
from typing import Callable, Dict, List, Optional, Tuple, Union, cast

from .consts import (
    TYPE_AND_FIELD_NAME,
    DataTypeName,
    PropertyDefinition,
    PropertyValue,
    ReturnType,
)
from .exceptions import EntityNotFoundError
from .graphql.property_input import PropertyInput
from .graphql.property_kind_enum import PropertyKind
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
    property_types: List[Dict[str, PropertyTypeInput]],
    properties_dict: Dict[str, PropertyValue],
) -> List[PropertyTypeInput]:
    """This function gets existing property types and dictionary, where key - are type names, and keys - new values
    formats data, validates existence of keys from properties_dict in property_types and returns list of PropertyTypeInput
 
        Args:
            property_types (List[Dict[str, pyinventory.graphql.property_type_input.PropertyTypeInput]]): list of existing property types
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
        property_type_names[property_type["name"]] = property_type

    for name, value in properties_dict.items():
        if name not in property_type_names:
            raise EntityNotFoundError(entity="PropertyType", entity_name=name)
        assert property_type_names[name][
            "isInstanceProperty"
        ], f"property {name} is not instance property"
        result = {
            "id": property_type_names[name]["id"],
            "name": name,
            "type": PropertyKind(property_type_names[name]["type"]),
        }
        result.update(
            get_graphql_input_field(
                property_type_name=name,
                type_key=property_type_names[name]["type"],
                value=value,
            )
        )
        properties.append(result)

    return properties


def get_graphql_property_inputs(
    property_types: List[Dict[str, PropertyTypeInput]],
    properties_dict: Dict[str, PropertyValue],
) -> List[PropertyInput]:
    """This function gets existing property types and dictionary, where key - are type names, and keys - new values
    formats data, validates existence of keys from properties_dict in property_types and returns list of PropertyInput
 
        Args:
            property_types (List[Dict[str, pyinventory.graphql.property_type_input.PropertyTypeInput]]): list of existing property types
            properties_dict (Dict[str, pyinventory.consts.PropertyValue]): dictionary of properties, where
                str: name of existing property
                PropertyValue: new value of existing type for this property
       
        Returns:
            List[pyinventory.graphql.property_input.PropertyInput]
 
        Raises:
            EntityNotFoundError if there any unknown property name in properties_dict keys
       
        Example:
        ```
            property_types = client.locationTypes[location_type].propertyTypes
            properties = get_graphql_property_inputs(property_types, properties_dict)
        ```
    """
    properties = []
    property_type_names = {}

    for property_type in property_types:
        property_type_names[property_type["name"]] = property_type

    for name, value in properties_dict.items():
        if name not in property_type_names:
            raise EntityNotFoundError(entity="PropertyType", entity_name=name)
        assert property_type_names[name][
            "isInstanceProperty"
        ], f"property {name} is not instance property"
        result = {"propertyTypeID": property_type_names[name]["id"]}
        result.update(
            get_graphql_input_field(
                property_type_name=name,
                type_key=property_type_names[name]["type"],
                value=value,
            )
        )
        properties.append(result)

    return properties


def _get_property_value(
    property_type: str, property: Dict[str, PropertyValue]
) -> Tuple[PropertyValue, ...]:
    formated_name = format_to_type_and_field_name(property_type)
    if formated_name is None:
        raise AssertionError(f"Unknown property type - {property_type}")

    str_fields = formated_name.graphql_field_name
    values = []
    for str_field in str_fields:
        if property_type == "date":
            date_data = property[str_field]
            values.append(datetime.strptime(cast(str, date_data), "%Y-%m-%d").date())
        else:
            values.append(property[str_field])
    return tuple(value for value in values)


def _get_property_default_value(
    name: str, type: str, value: PropertyValue
) -> Dict[str, PropertyValue]:
    if value is None:
        return {}
    return get_graphql_input_field(property_type_name=name, type_key=type, value=value)


def _make_property_types(
    properties: List[PropertyDefinition]
) -> List[Dict[str, PropertyValue]]:
    property_types = [
        {
            "name": arg.property_name,
            "type": arg.property_type,
            "index": i,
            **_get_property_default_value(
                arg.property_name, arg.property_type, arg.default_value
            ),
            "isInstanceProperty": arg.is_fixed,
        }
        for i, arg in enumerate(properties)
    ]
    return property_types


def property_type_to_kind(
    key: str, value: PropertyValue
) -> Union[PropertyValue, PropertyKind]:
    return value if key != "type" else PropertyKind(value)


def format_properties(
    properties: List[PropertyDefinition]
) -> List[Dict[str, Union[PropertyValue, PropertyKind]]]:
    property_types = _make_property_types(properties)
    return [
        {k: property_type_to_kind(k, v) for k, v in property_type.items()}
        for property_type in property_types
    ]


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
