#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import Sequence

from ..client import SymphonyClient
from ..consts import Entity
from ..exceptions import EntityNotFoundError
from ..graphql.property_type_fragment import PropertyTypeFragment


def get_property_types(
    client: SymphonyClient, entity_type: Entity, entity_name: str
) -> Sequence[PropertyTypeFragment]:
    """Get property types on specific entity.
    `entity_type` - ["LocationType", "EquipmentType", "ServiceType", "EquipmentPortType"]

        Args:
            entity_type (pyinventory.consts.Entity): existing entity type
            entity_name (str): existing entity name

        Returns:
            Sequence [pyinventory.graphql.property_type_input.PropertyTypeFragment ]

        Raises:
            EntityNotFoundError: if entity type does not found or does not have property types

        Example:
            ```
            property_type = client.get_property_types(
                entity_type=Entity.EquipmentType,
                entity_name="Card",
            )
            ```
    """

    existing_entity_types = {
        Entity.LocationType: client.locationTypes,
        Entity.EquipmentType: client.equipmentTypes,
        Entity.ServiceType: client.serviceTypes,
        Entity.EquipmentPortType: client.portTypes,
    }.get(entity_type, None)

    if existing_entity_types is None:
        raise EntityNotFoundError(entity=entity_type)
    existing_entity_type = existing_entity_types.get(entity_name, None)

    if existing_entity_type is None:
        raise EntityNotFoundError(entity=entity_type, entity_name=entity_name)

    return existing_entity_type.property_types


def get_property_type(
    client: SymphonyClient, entity_type: Entity, entity_name: str, property_type_id: str
) -> PropertyTypeFragment:
    """Get property type on specific entity.
    `entity_type` - ["LocationType", "EquipmentType", "ServiceType", "EquipmentPortType"]

        Args:
            entity_type (pyinventory.consts.Entity): existing entity type
            entity_name (str): existing entity name
            property_type_id (str): property type ID

        Returns:
            pyinventory.graphql.property_type_fragment.PropertyTypeFragment object

        Raises:
            EntityNotFounError: if property type with id=`property_type_id` does not found

        Example:
            ```
            property_type = client.get_property_type(
                entity_type=Entity.EquipmentType,
                entity_name="Card",
                property_type_id="12345",
            )
            ```
    """
    property_types = get_property_types(
        client=client, entity_type=entity_type, entity_name=entity_name
    )
    for property_type in property_types:
        if property_type.id == property_type_id:
            return property_type

    raise EntityNotFoundError(entity=Entity.PropertyType, entity_id=property_type_id)


def get_property_type_by_external_id(
    client: SymphonyClient,
    entity_type: Entity,
    entity_name: str,
    property_type_external_id: str,
) -> PropertyTypeFragment:
    """Get property type by external ID on specific entity.
    `entity_type` - ["LocationType", "EquipmentType", "ServiceType", "EquipmentPortType"]

        Args:
            entity_type (pyinventory.consts.Entity): existing entity type
            entity_name (str): existing entity name
            property_type_external_id (str): property type external ID

        Returns:
            pyinventory.graphql.property_type_fragment.PropertyTypeFragment object

        Raises:
            EntityNotFounError: property type with external_id=`property_type_external_id` is not found

        Example:
            ```
            property_type = client.get_property_type_by_external_id(
                entity_type=Entity.EquipmentType,
                entity_name="Card",
                property_type_external_id="12345",
            )
            ```
    """
    property_types = get_property_types(
        client=client, entity_type=entity_type, entity_name=entity_name
    )
    for property_type in property_types:
        if property_type.externalId == property_type_external_id:
            return property_type

    raise EntityNotFoundError(
        entity=Entity.PropertyType, msg=f"<external_id: {property_type_external_id}>"
    )
