#!/usr/bin/env python3
# pyre-strict

from datetime import date, datetime
from enum import Enum
from typing import Any, Dict, List, NamedTuple, Optional, Tuple, Type, TypeVar, Union

from .graphql.image_entity_enum import ImageEntity
from .graphql.property_kind_enum import PropertyKind


__version__ = "2.4.0"

INVENTORY_ENDPOINT = "https://{}.thesymphony.cloud"
LOCALHOST_INVENTORY_ENDPOINT = "https://{}.localtest.me"
INVENTORY_GRAPHQL_ENDPOINT = "/graph/query"
INVENTORY_STORE_PUT_ENDPOINT = "/store/put"
INVENTORY_STORE_DELETE_ENDPOINT = "/store/delete?key={}"


ReturnType = TypeVar("ReturnType")
PropertyValue = Union[date, float, int, str, bool, Tuple[float, float]]
PropertyValueType = Union[
    Type[date], Type[float], Type[int], Type[str], Type[bool], Type[Tuple[float, float]]
]


class PropertyDefinition(NamedTuple):
    """
    Attributes:
        property_name (str): type name
        property_type (str): enum["string", "int", "bool", "float", "date", "enum", "range", 
            "email", "gps_location", "equipment", "location", "service", "datetime_local"]
        default_value (PropertyValue): default property value
        is_fixed (bool): fixed value flag
    """

    property_name: str
    property_kind: PropertyKind
    default_value: Optional[PropertyValue]
    is_fixed: Optional[bool] = False


class DataTypeName(NamedTuple):
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
    name: str
    id: str
    propertyTypes: List[Dict[str, PropertyValue]]


class Location(NamedTuple):
    name: str
    id: str
    latitude: float
    longitude: float
    externalId: Optional[str]
    locationTypeName: str


class EquipmentType(NamedTuple):
    name: str
    category: Optional[str]
    id: str
    propertyTypes: List[Dict[str, PropertyValue]]
    positionDefinitions: List[Dict[str, Any]]
    portDefinitions: List[Dict[str, Any]]


class EquipmentPortType(NamedTuple):
    """
    Attributes:
        id (str): equipment port type ID
        name (str): equipment port type name
        properties (List[Dict[str, PropertyValue]]): list of equipment port type propertyTypes to their default values
        link_properties (List[Dict[str, PropertyValue]]): list of equipment port type linkPropertyTypes to their default values
    """

    id: str
    name: str
    properties: List[Dict[str, PropertyValue]]
    link_properties: List[Dict[str, PropertyValue]]


class Equipment(NamedTuple):
    """
    Attributes:
        name (str): equipment name
        id (str): equipment ID
    """

    name: str
    id: str


class Link(NamedTuple):
    id: str


class EquipmentPortDefinition(NamedTuple):
    """
    Attributes:
        id (str): equipment port definition ID
        name (str): equipment port definition name
    """

    id: str
    name: str


class EquipmentPort(NamedTuple):
    """
    Attributes:
        id (str): equipment port ID
        properties (List[Dict[str, PropertyValue]]): list of equipment port properties
        definition (pyinventory.Consts.EquipmentPortDefinition): port definition
        link (Optional[pyinventory.consts.Link]): link
    """

    id: str
    properties: List[Dict[str, PropertyValue]]
    definition: EquipmentPortDefinition
    link: Optional[Link]


class SiteSurvey(NamedTuple):
    name: str
    id: str
    completionTime: datetime
    sourceFileId: Optional[str]
    sourceFileName: Optional[str]
    sourceFileKey: Optional[str]
    forms: Dict[str, Dict[str, Any]]


class ServiceType(NamedTuple):
    name: str
    id: str
    hasCustomer: bool
    propertyTypes: List[Dict[str, PropertyValue]]


class Customer(NamedTuple):
    name: str
    id: str
    externalId: Optional[str]


class ServiceEndpoint(NamedTuple):
    id: str
    port: EquipmentPort
    role: str


class Service(NamedTuple):
    name: str
    id: str
    externalId: Optional[str]
    customer: Optional[Customer]
    endpoints: List[ServiceEndpoint]
    links: List[Link]


class Document(NamedTuple):
    id: str
    name: str
    parentId: str
    parentEntity: ImageEntity
    category: Optional[str]


class Entity(Enum):
    Location = "Location"
    LocationType = "LocationType"
    Equipment = "Equipment"
    EquipmentType = "EquipmentType"
    EquipmentPort = "EquipmentPort"
    EquipmentPortType = "EquipmentPortType"
    Link = "Link"
    Service = "Service"
    ServiceType = "ServiceType"
    ServiceEndpoint = "ServiceEndpoint"
    SiteSurvey = "SiteSurvey"
    Customer = "Customer"
    Document = "Document"
    PropertyType = "PropertyType"
