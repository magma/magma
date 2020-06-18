#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from datetime import date
from numbers import Number
from typing import List, NamedTuple, Optional, Sequence, Tuple, Type, TypeVar, Union

from pysymphony.graphql.enum.image_entity import ImageEntity

from ..graphql.enum.property_kind import PropertyKind
from ..graphql.enum.user_role import UserRole
from ..graphql.enum.user_status import UserStatus
from ..graphql.fragment.equipment_port_definition import EquipmentPortDefinitionFragment
from ..graphql.fragment.equipment_position_definition import (
    EquipmentPositionDefinitionFragment,
)
from ..graphql.fragment.property import PropertyFragment


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
        default_raw_value (str): default property value as a string
            ```
            string         - "string"
            int            - "123"
            bool           - "true" / "True" / "TRUE"
            float           - "0.123456"
            date           - "24/10/2020"
            range          - "0.123456 - 0.2345" / "1 - 2"
            email          - "email@some.domain"
            gps_location   - "0.1234, 0.2345"
            ```
        id (Optional[str]):  ID
        is_fixed (bool): fixed value flag
        external_id (str): property type external ID
        is_mandatory (bool): mandatory value flag
        is_deleted (bool): is delete flag
    """

    property_name: str
    property_kind: PropertyKind
    default_raw_value: Optional[str]
    id: Optional[str] = None
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
        property_types (Sequence[ `pyinventory.common.data_class.PropertyDefinition` ]): property types sequence
    """

    name: str
    id: str
    property_types: Sequence[PropertyDefinition]


class Location(NamedTuple):
    """
    Attributes:
        name (str): name
        id (str): ID
        latitude (Number): latitude
        longitude (Number): longitude
        external_id (Optional[str]): external ID
        location_type_name (str): location type name
        properties (Sequence[ `pyinventory.graphql.fragment.property.PropertyFragment` ])
    """

    name: str
    id: str
    latitude: Number
    longitude: Number
    external_id: Optional[str]
    location_type_name: str
    properties: Sequence[PropertyFragment]


class EquipmentType(NamedTuple):
    """
    Attributes:
        name (str): name
        category (Optional[str]): category
        id (str): ID
        property_types (Sequence[ `pyinventory.common.data_class.PropertyDefinition` ]):  property types sequence
        position_definitions (Sequence[ `pyinventory.graphql.fragment.equipment_position_definition.EquipmentPositionDefinitionFragment` ]): position definitions sequence
        port_definitions (Sequence[ `pyinventory.graphql.fragment.equipment_port_definition.EquipmentPortDefinitionFragment` ]): port definition sequence
    """

    name: str
    category: Optional[str]
    id: str
    property_types: Sequence[PropertyDefinition]
    position_definitions: Sequence[EquipmentPositionDefinitionFragment]
    port_definitions: Sequence[EquipmentPortDefinitionFragment]


class EquipmentPortType(NamedTuple):
    """
    Attributes:
        id (str): ID
        name (str): name
        property_types (Sequence[ `pyinventory.common.data_class.PropertyDefinition` ]): property types sequence
        link_property_types (Sequence[ `pyinventory.common.data_class.PropertyDefinition` ]): link property types sequence
    """

    id: str
    name: str
    property_types: Sequence[PropertyDefinition]
    link_property_types: Sequence[PropertyDefinition]


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
        properties (Sequence[ `pyinventory.graphql.fragment.property.PropertyFragment` ]): properties sequence
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


class Customer(NamedTuple):
    """
    Attributes:
        id (str): ID
        name (str): name
        external_id (Optional[str]): external ID
    """

    id: str
    name: str
    external_id: Optional[str]


class Service(NamedTuple):
    """
    Attributes:
        name (str): name
        id (str): ID
        external_id (Optional[str]): external ID
        service_type_name (str): existing service tyoe name
        customer (Optional[ `pyinventory.common.data_class.Customer` ]): customer
        properties (Sequence[ `pyinventory.graphql.fragment.property..PropertyFragment` ]): properties sequence
    """

    id: str
    name: str
    external_id: Optional[str]
    service_type_name: str
    customer: Optional[Customer]
    properties: Sequence[PropertyFragment]


class ServiceEndpointDefinition(NamedTuple):
    """
    Attributes:
        id (str): ID
        name (str): name
        endpoint_definition_index (int): index
        role (str): role
        equipment_type_id (str): equipment type ID
    """

    id: Optional[str]
    name: str
    endpoint_definition_index: int
    role: Optional[str]
    equipment_type_id: str


class ServiceEndpoint(NamedTuple):
    """
    Attributes:
        id (str): ID
        equipment_id (str): existing equipment ID
        service_id (str): existing service ID
        definition_id (str): existing service endpoint definition ID
    """

    id: str
    equipment_id: str
    service_id: str
    definition_id: str


class ServiceType(NamedTuple):
    """
    Attributes:
        name (str): name
        id (str): ID
        has_customer (bool): customer existence flag
        property_types (Sequence[ `pyinventory.common.data_class.PropertyDefinition` ]): property types sequence
        endpoint_definitions (List[ `pyinventory.common.data_class.ServiceEndpointDefinition` ]): service endpoint definitions list
    """

    name: str
    id: str
    has_customer: bool
    property_types: Sequence[PropertyDefinition]
    endpoint_definitions: List[ServiceEndpointDefinition]


class Document(NamedTuple):
    """
    Attributes:
        id (str): ID
        name (str): name
        parent_id (str): parent ID
        parent_entity (ImageEntity): parent entity
        category (Optional[str]): category
    """

    id: str
    name: str
    parent_id: str
    parent_entity: ImageEntity
    category: Optional[str]


class User(NamedTuple):
    """
    Attributes:
        id (str): ID
        auth_id (str): auth ID
        email (str): email
        status (UserStatus): status
        role (UserRole): role
    """

    id: str
    auth_id: str
    email: str
    status: UserStatus
    role: UserRole
