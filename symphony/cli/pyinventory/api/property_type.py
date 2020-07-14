#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import Dict, List, Optional, Sequence, Union

from pysymphony import SymphonyClient

from ..common.cache import (
    EQUIPMENT_TYPES,
    LOCATION_TYPES,
    PORT_TYPES,
    SERVICE_TYPES,
    Cache,
)
from ..common.data_class import (
    EquipmentPortType,
    EquipmentType,
    LocationType,
    PropertyDefinition,
    ServiceType,
)
from ..common.data_enum import Entity
from ..common.data_format import format_to_property_type_input
from ..exceptions import EntityNotFoundError
from ..graphql.input.property_type import PropertyTypeInput


def get_property_types(
    client: SymphonyClient, entity_type: Entity, entity_name: str
) -> Sequence[PropertyDefinition]:
    """Get property types on specific entity.

        :param entity_type: Existing entity type
        :type entity_type: :class:`~pyinventory.common.data_enum.Entity`
        :param entity_name: Existing entity name
        :type entity_name: str

        :raises:
            :class:`~pyinventory.exceptions.EntityNotFoundError`: `entity_type` does not found or does not have property types

        :return: Sequence[PropertyDefinition]
        :rtype: Sequence[ :class:`~pyinventory.common.data_class.PropertyDefinition` ]

        **Example**

        .. code-block:: python

            property_type = client.get_property_types(
                entity_type=Entity.EquipmentType,
                entity_name="Card",
            )
    """

    caches: Dict[
        Entity,
        Union[
            Cache[LocationType],
            Cache[EquipmentType],
            Cache[ServiceType],
            Cache[EquipmentPortType],
        ],
    ] = {
        Entity.LocationType: LOCATION_TYPES,
        Entity.EquipmentType: EQUIPMENT_TYPES,
        Entity.ServiceType: SERVICE_TYPES,
        Entity.EquipmentPortType: PORT_TYPES,
    }

    existing_entity_types: Optional[
        Union[
            Cache[LocationType],
            Cache[EquipmentType],
            Cache[ServiceType],
            Cache[EquipmentPortType],
        ],
    ] = caches.get(entity_type, None)

    if existing_entity_types is None:
        raise EntityNotFoundError(entity=entity_type)

    existing_entity_type = None
    if entity_name in existing_entity_types:
        existing_entity_type = existing_entity_types[entity_name]

    if existing_entity_type is None:
        raise EntityNotFoundError(entity=entity_type, entity_name=entity_name)

    return existing_entity_type.property_types


def get_property_type(
    client: SymphonyClient, entity_type: Entity, entity_name: str, property_type_id: str
) -> PropertyDefinition:
    """Get property type on specific entity.

        :param entity_type: Existing entity type
        :type entity_type: :class:`~pyinventory.common.data_enum.Entity`
        :param entity_name: Existing entity name
        :type entity_name: str
        :param property_type_id: Existing property type ID
        :type property_type_id: str

        :raises:
            :class:`~pyinventory.exceptions.EntityNotFoundError`: Property type with id=`property_type_id` does not found

        :return: PropertyDefinition object
        :rtype: :class:`~pyinventory.common.data_class.PropertyDefinition`

        **Example**

        .. code-block:: python

            property_type = client.get_property_type(
                entity_type=Entity.EquipmentType,
                entity_name="Card",
                property_type_id="12345",
            )
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
    """Get property type ID on specific entity.

        :param entity_type: Existing entity type
        :type entity_type: :class:`~pyinventory.common.data_enum.Entity`
        :param entity_name: Existing entity name
        :type entity_name: str
        :param property_type_name: Existing property type name
        :type property_type_name: str

        :raises:
            :class:`~pyinventory.exceptions.EntityNotFoundError`: Property type with name=`property_type_name` does not found

        :return: property type ID
        :rtype: str

        **Example**

        .. code-block:: python

            property_type = client.get_property_type_id(
                entity_type=Entity.EquipmentType,
                entity_name="Card",
                property_type_name="IP",
            )
    """
    property_types = get_property_types(
        client=client, entity_type=entity_type, entity_name=entity_name
    )
    for property_type in property_types:
        if (
            property_type.property_name == property_type_name
            and property_type.id is not None
        ):
            return property_type.id

    raise EntityNotFoundError(
        entity=Entity.PropertyType, entity_name=property_type_name
    )


def get_property_type_by_external_id(
    client: SymphonyClient,
    entity_type: Entity,
    entity_name: str,
    property_type_external_id: str,
) -> PropertyDefinition:
    """Get property type by external ID on specific entity.

        :param entity_type: Existing entity type
        :type entity_type: :class:`~pyinventory.common.data_enum.Entity`
        :param entity_name: Existing entity name
        :type entity_name: str
        :param property_type_external_id: Existing property type external ID
        :type property_type_external_id: str

        :raises:
            :class:`~pyinventory.exceptions.EntityNotFoundError`: Property type with external_id=`property_type_external_id` does not found

        :return: PropertyDefinition object
        :rtype: :class:`~pyinventory.common.data_class.PropertyDefinition`

        **Example**

        .. code-block:: python

            property_type = client.get_property_type_by_external_id(
                entity_type=Entity.EquipmentType,
                entity_name="Card",
                property_type_external_id="12345",
            )
    """
    property_types = get_property_types(
        client=client, entity_type=entity_type, entity_name=entity_name
    )
    for property_type in property_types:
        if property_type.external_id == property_type_external_id:
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

        :param entity_type: Existing entity type
        :type entity_type: :class:`~pyinventory.common.data_enum.Entity`
        :param entity_name: Existing entity name
        :type entity_name: str
        :param property_type_id: Existing property type ID
        :type property_type_id: str
        :param new_property_definition: New property definition
        :type new_property_definition: :class:`~pyinventory.common.data_class.PropertyDefinition`

        :class:`~pyinventory.exceptions.EntityNotFoundError`: Property type with id=`property_type_id` does not found

        :return: PropertyTypeInputs List
        :rtype: List[ :class:`~pyinventory.graphql.input.property_type.PropertyTypeInput` ]

        **Example**

        .. code-block:: python

            property_types = client.edit_property_type(
                entity_type=Entity.EquipmentType,
                entity_name="Card",
                property_type_id="12345",
                property_definition=PropertyDefinition(
                    property_name="new_name",
                    property_kind=PropertyKind.string,
                    default_raw_value=None,
                    is_fixed=False,
                    external_id="ex_12345",
                ),
            )
    """
    property_types = get_property_types(
        client=client, entity_type=entity_type, entity_name=entity_name
    )
    edited_property_types = []

    for property_type in property_types:
        property_type_input = format_to_property_type_input(property_type, is_new=False)
        if property_type_input.id == property_type_id:
            formated_property_definitions = format_to_property_type_input(
                new_property_definition, is_new=False
            )
            formated_property_definitions.id = property_type_input.id
            property_type_input = formated_property_definitions

        edited_property_types.append(property_type_input)

    return edited_property_types
