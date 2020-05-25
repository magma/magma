#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from datetime import date
from numbers import Number
from typing import List, NamedTuple, Optional, Sequence, Tuple, Type, TypeVar, Union

from pysymphony.graphql.enum.image_entity import ImageEntity

from ..graphql.enum.property_kind import PropertyKind
from ..graphql.fragment.equipment_port_definition import EquipmentPortDefinitionFragment
from ..graphql.fragment.equipment_position_definition import (
    EquipmentPositionDefinitionFragment,
)
from ..graphql.fragment.property import PropertyFragment
from ..graphql.fragment.property_type import PropertyTypeFragment


ReturnType = TypeVar("ReturnType")
PropertyValue = Union[date, float, int, str, bool, Tuple[float, float]]
PropertyValueType = Union[
    Type[date], Type[float], Type[int], Type[str], Type[bool], Type[Tuple[float, float]]
]


class PropertyDefinition(NamedTuple):
    """
    Attributes:
        property_name (str): type name
        property_kind ( `pyinventory.graphql.enum.property_kind.PropertyKind` ): property kind
        default_value (Optional[PropertyValue]): default property value
        is_fixed (bool): fixed value flag
        external_id (str): property type external ID
        is_mandatory (bool): mandatory value flag
        is_deleted (bool): is delete flag
    """

    property_name: str
    property_kind: PropertyKind
    default_value: Optional[PropertyValue] = None
    is_fixed: Optional[bool] = False
    external_id: Optional[str] = None
    is_mandatory: Optional[bool] = False
    is_deleted: Optional[bool] = False


class DataTypeName(NamedTuple):
    """
    Attributes:
        data_type (PropertyValueType): data type
        graphql_field_name (Tuple[str, ...]): graphql field name, in case of `gps_location` it is Tuple[`latitudeValue`, `longitudeValue`]
    """

    data_type: PropertyValueType
    graphql_field_name: Tuple[str, ...]


TYPE_AND_FIELD_NAME = {
    "date": DataTypeName(data_type=date, graphql_field_name=("stringValue",)),
    "float": DataTypeName(data_type=float, graphql_field_name=("floatValue",)),
    "int": DataTypeName(data_type=int, graphql_field_name=("intValue",)),
    "email": DataTypeName(data_type=str, graphql_field_name=("stringValue",)),
    "string": DataTypeName(data_type=str, graphql_field_name=("stringValue",)),
    "bool": DataTypeName(data_type=bool, graphql_field_name=("booleanValue",)),
    "gps_location": DataTypeName(
        data_type=tuple, graphql_field_name=("latitudeValue", "longitudeValue")
    ),
}


class LocationType(NamedTuple):
    """
    Attributes:
        name (str): name
        id (str):  ID
        property_types (Sequence[ `pyinventory.graphql.fragment.property_type.PropertyTypeFragment` ]): property types sequence
    """

    name: str
    id: str
    property_types: Sequence[PropertyTypeFragment]


class Location(NamedTuple):
    """
    Attributes:
        name (str): name
        id (str): ID
        latitude (Number): latitude
        longitude (Number): longitude
        externalId (Optional[str]): external ID
        locationTypeName (str): location type name
    """

    name: str
    id: str
    latitude: Number
    longitude: Number
    externalId: Optional[str]
    locationTypeName: str
    properties: Sequence[PropertyFragment]


class EquipmentType(NamedTuple):
    """
    Attributes:
        name (str): name
        category (Optional[str]): category
        id (str): ID
        property_types (Sequence[PropertyTypeFragment]):  property types sequence
        position_definitions (Sequence[EquipmentPositionDefinitionFragment]): position definitions sequence
        port_definitions (Sequence[EquipmentPortDefinitionFragment]): port definition sequence
    """

    name: str
    category: Optional[str]
    id: str
    property_types: Sequence[PropertyTypeFragment]
    position_definitions: Sequence[EquipmentPositionDefinitionFragment]
    port_definitions: Sequence[EquipmentPortDefinitionFragment]


class EquipmentPortType(NamedTuple):
    """
    Attributes:
        id (str): ID
        name (str): name
        property_types (Sequence[PropertyTypeFragment]): property types sequence
        link_property_types (Sequence[PropertyTypeFragment]): link property types sequence
    """

    id: str
    name: str
    property_types: Sequence[PropertyTypeFragment]
    link_property_types: Sequence[PropertyTypeFragment]


class Equipment(NamedTuple):
    """
    Attributes:
        id (str): ID
        external_id (Optional[str]): external ID
        name (str): name
        equipment_type_name (str): equipment type name
    """

    id: str
    external_id: Optional[str]
    name: str
    equipment_type_name: str


class Link(NamedTuple):
    """
    Attributes:
        id (str): link ID
        properties (Sequence[PropertyFragment]): properties sequence
        service_ids (List[str]): service IDs
    """

    id: str
    properties: Sequence[PropertyFragment]
    service_ids: List[str]


class EquipmentPortDefinition(NamedTuple):
    """
    Attributes:
        id (str): ID
        name (str): name
        port_type_name (Optional[str]): port type name
    """

    id: str
    name: str
    port_type_name: Optional[str] = None


class EquipmentPort(NamedTuple):
    """
    Attributes:
        id (str): equipment port ID
        properties (Sequence[PropertyFragment]): properties sequence
        definition ( `pyinventory.common.data_class.EquipmentPortDefinition` ): port definition
        link (Optional[ `pyinventory.common.data_class.Link` ]): link
    """

    id: str
    properties: Sequence[PropertyFragment]
    definition: EquipmentPortDefinition
    link: Optional[Link]


class ServiceType(NamedTuple):
    """
    Attributes:
        name (str): name
        id (str): ID
        has_customer (bool): customer existence flag
        property_types (Sequence[PropertyTypeFragment]): property types sequence
    """

    name: str
    id: str
    has_customer: bool
    property_types: Sequence[PropertyTypeFragment]


class Customer(NamedTuple):
    """
    Attributes:
        name (str): name
        id (str): ID
        externalId (Optional[str]): external ID
    """

    name: str
    id: str
    externalId: Optional[str]


class ServiceEndpoint(NamedTuple):
    """
    Attributes:
        id (str): ID
        port (Optional[EquipmentPort]): port
        type (str): type
    """

    id: str
    port: Optional[EquipmentPort]
    type: str


class Service(NamedTuple):
    """
    Attributes:
        name (str): name
        id (str): ID
        external_id (Optional[str]): external ID
        customer (Optional[Customer]): customer
        properties (Sequence[PropertyFragment]): properties sequence
    """

    id: str
    name: str
    external_id: Optional[str]
    customer: Optional[Customer]
    properties: Sequence[PropertyFragment]


class Document(NamedTuple):
    """
    Attributes:
        id (str): ID
        name (str): name
        parentId (str): parent ID
        parentEntity (ImageEntity): parent entity
        category (Optional[str]): category
    """

    id: str
    name: str
    parentId: str
    parentEntity: ImageEntity
    category: Optional[str]
