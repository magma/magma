#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import List, Sequence

from .._utils import format_property_definitions, get_property_type_input
from ..client import SymphonyClient
from ..common.data_class import PropertyDefinition
from ..common.data_enum import Entity
from ..exceptions import EntityNotFoundError
from ..graphql.property_type_fragment import PropertyTypeFragment
from ..graphql.property_type_input import PropertyTypeInput


def get_property_types(
    client: SymphonyClient, entity_type: Entity, entity_name: str
) -> Sequence[PropertyTypeFragment]:
    """Get property types on specific entity. `entity_type` - ["LocationType", "EquipmentType", "ServiceType", "EquipmentPortType"]

        Args:
            entity_type ( `pyinventory.common.data_enum.Entity` ): existing entity type
            entity_name (str): existing entity name

        Returns:
            Sequence[ `pyinventory.graphql.property_type_fragment.PropertyTypeFragment` ]

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: if entity type does not found or does not have property types

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
    # pyre-fixme[16]: `None` has no attribute `get`.
    existing_entity_type = existing_entity_types.get(entity_name, None)

    if existing_entity_type is None:
        raise EntityNotFoundError(entity=entity_type, entity_name=entity_name)

    return existing_entity_type.property_types


def get_property_type(
    client: SymphonyClient, entity_type: Entity, entity_name: str, property_type_id: str
) -> PropertyTypeFragment:
    """Get property type on specific entity. `entity_type` - ["LocationType", "EquipmentType", "ServiceType", "EquipmentPortType"]

        Args:
            entity_type ( `pyinventory.common.data_enum.Entity` ): existing entity type
            entity_name (str): existing entity name
            property_type_id (str): property type ID

        Returns:
            `pyinventory.graphql.property_type_fragment.PropertyTypeFragment` object

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: if property type with id=`property_type_id` does not found

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


def get_property_type_id(
    client: SymphonyClient,
    entity_type: Entity,
    entity_name: str,
    property_type_name: str,
) -> str:
    """Get property type ID on specific entity. `entity_type` - ["LocationType", "EquipmentType", "ServiceType", "EquipmentPortType"]

        Args:
            entity_type ( `pyinventory.common.data_enum.Entity` ): existing entity type
            entity_name (str): existing entity name
            property_type_name (str): property type ID

        Returns:
            property type ID (str): property type ID

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: if property type with id=`property_type_id` does not found

        Example:
            ```
            property_type = client.get_property_type_id(
                entity_type=Entity.EquipmentType,
                entity_name="Card",
                property_type_name="IP",
            )
            ```
    """
    property_types = get_property_types(
        client=client, entity_type=entity_type, entity_name=entity_name
    )
    for property_type in property_types:
        if property_type.name == property_type_name:
            return property_type.id

    raise EntityNotFoundError(
        entity=Entity.PropertyType, entity_name=property_type_name
    )


def get_property_type_by_external_id(
    client: SymphonyClient,
    entity_type: Entity,
    entity_name: str,
    property_type_external_id: str,
) -> PropertyTypeFragment:
    """Get property type by external ID on specific entity. `entity_type` - ["LocationType", "EquipmentType", "ServiceType", "EquipmentPortType"]

        Args:
            entity_type ( `pyinventory.common.data_enum.Entity` ): existing entity type
            entity_name (str): existing entity name
            property_type_external_id (str): property type external ID

        Returns:
            `pyinventory.graphql.property_type_fragment.PropertyTypeFragment` object

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: property type with external_id=`property_type_external_id` is not found

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


def edit_property_type(
    client: SymphonyClient,
    entity_type: Entity,
    entity_name: str,
    property_type_id: str,
    new_property_definition: PropertyDefinition,
) -> List[PropertyTypeInput]:
    """Edit specific property type on specific entity. `entity_type` - ["LocationType", "EquipmentType", "ServiceType", "EquipmentPortType"]

        Args:
            entity_type ( `pyinventory.common.data_enum.Entity` ): existing entity type
            entity_name (str): existing entity name
            property_type_id (str): existing property type ID
            new_property_definition ( `pyinventory.common.data_class.PropertyDefinition` ): new property definition

        Returns:
            List[ `pyinventory.graphql.property_type_input.PropertyTypeInput` ]

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: property type with external_id=`property_type_external_id` is not found

        Example:
            ```
            property_types = client.edit_property_type(
                entity_type=Entity.EquipmentType,
                entity_name="Card",
                property_type_id="12345",
                property_definition=PropertyDefinition(
                    property_name="new_name",
                    property_kind=PropertyKind.string,
                    default_value=None,
                    is_fixed=False,
                    external_id="ex_12345",
                ),
            )
            ```
    """
    property_types = get_property_types(
        client=client, entity_type=entity_type, entity_name=entity_name
    )
    edited_property_types = []

    for property_type in property_types:
        property_type_input = get_property_type_input(property_type, is_new=False)
        if property_type_input.id == property_type_id:
            formated_property_definitions = format_property_definitions(
                [new_property_definition]
            )
            formated_property_definitions[0].id = property_type_input.id
            property_type_input = formated_property_definitions[0]

        edited_property_types.append(property_type_input)

    return edited_property_types
