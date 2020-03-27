#!/usr/bin/env python3

from datetime import date, datetime
from enum import Enum
from numbers import Number
from typing import (
    Any,
    Dict,
    List,
    NamedTuple,
    Optional,
    Sequence,
    Tuple,
    Type,
    TypeVar,
    Union,
)

from .graphql.equipment_port_definition_fragment import EquipmentPortDefinitionFragment
from .graphql.equipment_position_definition_fragment import (
    EquipmentPositionDefinitionFragment,
)
from .graphql.image_entity_enum import ImageEntity
from .graphql.property_fragment import PropertyFragment
from .graphql.property_kind_enum import PropertyKind
from .graphql.property_type_fragment import PropertyTypeFragment
from .graphql.user_role_enum import UserRole
from .graphql.user_status_enum import UserStatus


__version__ = "2.5.0"

INVENTORY_ENDPOINT = "https://{}.thesymphony.cloud"
LOCALHOST_INVENTORY_ENDPOINT = "https://{}.localtest.me"
INVENTORY_GRAPHQL_ENDPOINT = "/graph/query"
INVENTORY_STORE_PUT_ENDPOINT = "/store/put"
INVENTORY_LOGIN_ENDPOINT = "/user/login"
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
    property_types: Sequence[PropertyTypeFragment]


class Location(NamedTuple):
    name: str
    id: str
    latitude: Number
    longitude: Number
    externalId: Optional[str]
    locationTypeName: str


class EquipmentType(NamedTuple):
    name: str
    category: Optional[str]
    id: str
    property_types: Sequence[PropertyTypeFragment]
    position_definitions: Sequence[EquipmentPositionDefinitionFragment]
    port_definitions: Sequence[EquipmentPortDefinitionFragment]


class EquipmentPortType(NamedTuple):
    """
    Attributes:
        id (str): equipment port type ID
        name (str): equipment port type name
        property_types (List[Dict[str, PropertyValue]]): list of equipment port type propertyTypes to their default values
        link_property_types (List[Dict[str, PropertyValue]]): list of equipment port type linkPropertyTypes to their default values
    """

    id: str
    name: str
    property_types: Sequence[PropertyTypeFragment]
    link_property_types: Sequence[PropertyTypeFragment]


class Equipment(NamedTuple):
    """
    Attributes:
        name (str): equipment name
        id (str): equipment ID
        equipment_type_name (str): equipment type name
    """

    name: str
    id: str
    equipment_type_name: str


class Link(NamedTuple):
    """
    Attributes:
        id (str): link ID
        service_ids (List[str]): service IDs 
    """

    id: str
    properties: Sequence[PropertyFragment]
    service_ids: List[str]


class EquipmentPortDefinition(NamedTuple):
    """
    Attributes:
        id (str): equipment port definition ID
        name (str): equipment port definition name
        port_type_name (Optional[str]): equipment port definition port type name
    """

    id: str
    name: str
    port_type_name: Optional[str] = None


class EquipmentPort(NamedTuple):
    """
    Attributes:
        id (str): equipment port ID
        properties (List[Dict[str, PropertyValue]]): list of equipment port properties
        definition (pyinventory.consts.EquipmentPortDefinition): port definition
        link (Optional[pyinventory.consts.Link]): link
    """

    id: str
    properties: Sequence[PropertyFragment]
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
    property_types: Sequence[PropertyTypeFragment]


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


class User(NamedTuple):
    id: str
    auth_id: str
    email: str
    status: UserStatus
    role: UserRole


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
    Property = "Property"
    User = "User"
