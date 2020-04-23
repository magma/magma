#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from enum import Enum


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
