#!/usr/bin/env python3
# pyre-strict

from datetime import datetime
from enum import Enum
from typing import Any, Dict, List, NamedTuple, Optional

from .reporter import DummyReporter


__version__ = "2.3.3"

DUMMY_REPORTER = DummyReporter()

INVENTORY_ENDPOINT = "https://{}.thesymphony.cloud"
LOCALHOST_INVENTORY_ENDPOINT = "https://{}.localtest.me"
INVENTORY_GRAPHQL_ENDPOINT = "/graph/query"
INVENTORY_LOGIN_ENDPOINT = "/user/login"
INVENTORY_STORE_PUT_ENDPOINT = "/store/put"
INVENTORY_STORE_DELETE_ENDPOINT = "/store/delete?key={}"

PROPERTY_TYPE_VALUES = """stringValue
    intValue
    booleanValue
    floatValue
    latitudeValue
    longitudeValue
    isEditable
    isInstanceProperty"""

PROPERTY_VALUES = """stringValue
    intValue
    booleanValue
    floatValue
    latitudeValue
    longitudeValue"""


class ImageEntity(Enum):
    LOCATION = "LOCATION"
    WORK_ORDER = "WORK_ORDER"
    SITE_SURVEY = "SITE_SURVEY"
    EQUIPMENT = "EQUIPMENT"


class LocationType(NamedTuple):
    name: str
    id: str
    propertyTypes: List[Dict[str, Any]]


class Location(NamedTuple):
    name: str
    id: str
    externalId: Optional[str]


class EquipmentType(NamedTuple):
    name: str
    category: Optional[str]
    id: str
    propertyTypes: List[Dict[str, Any]]
    positionDefinitions: List[Dict[str, Any]]
    portDefinitions: List[Dict[str, Any]]


class Equipment(NamedTuple):
    name: str
    id: str


class Link(NamedTuple):
    id: str


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
    propertyTypes: List[Dict[str, Any]]


class Customer(NamedTuple):
    name: str
    id: str
    externalId: Optional[str]


class Service(NamedTuple):
    name: str
    id: str
    externalId: Optional[str]
    customer: Optional[Customer]
    terminationPoints: List[Equipment]
    links: List[Link]


class Document(NamedTuple):
    id: str
    name: str
    parentId: str
    parentEntity: ImageEntity
