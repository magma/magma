#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import Dict, List, Optional

from pysymphony import SymphonyClient

from .._utils import get_port_definition_input, get_position_definition_input
from ..common.cache import EQUIPMENT_TYPES, PORT_TYPES
from ..common.data_class import (
    Equipment,
    EquipmentPortType,
    EquipmentType,
    PropertyDefinition,
)
from ..common.data_enum import Entity
from ..common.data_format import (
    format_to_property_definitions,
    format_to_property_type_input,
    format_to_property_type_inputs,
)
from ..exceptions import EntityNotFoundError
from ..graphql.input.add_equipment_type import AddEquipmentTypeInput
from ..graphql.input.edit_equipment_type import EditEquipmentTypeInput
from ..graphql.input.equipment_port import EquipmentPortInput
from ..graphql.input.equipment_position import EquipmentPositionInput
from ..graphql.input.property_type import PropertyTypeInput
from ..graphql.mutation.add_equipment_type import AddEquipmentTypeMutation
from ..graphql.mutation.edit_equipment_type import EditEquipmentTypeMutation
from ..graphql.mutation.remove_equipment_type import RemoveEquipmentTypeMutation
from ..graphql.query.equipment_port_types import EquipmentPortTypesQuery
from ..graphql.query.equipment_type_equipments import EquipmentTypeEquipmentQuery
from ..graphql.query.equipment_types import EquipmentTypesQuery
from .equipment import delete_equipment
from .property_type import (
    edit_property_type,
    get_property_type,
    get_property_type_by_external_id,
)


def _populate_equipment_types(client: SymphonyClient) -> None:
    edges = EquipmentTypesQuery.execute(client).edges

    for edge in edges:
        node = edge.node
        if node:
            EQUIPMENT_TYPES[node.name] = EquipmentType(
                name=node.name,
                category=node.category,
                id=node.id,
                property_types=format_to_property_definitions(node.propertyTypes),
                position_definitions=node.positionDefinitions,
                port_definitions=node.portDefinitions,
            )


def _populate_equipment_port_types(client: SymphonyClient) -> None:
    edges = EquipmentPortTypesQuery.execute(client).edges

    for edge in edges:
        node = edge.node
        if node:
            PORT_TYPES[node.name] = EquipmentPortType(
                id=node.id,
                name=node.name,
                property_types=format_to_property_definitions(node.propertyTypes),
                link_property_types=format_to_property_definitions(
                    node.linkPropertyTypes
                ),
            )


def _add_equipment_type(
    client: SymphonyClient,
    name: str,
    category: Optional[str],
    properties: List[PropertyTypeInput],
    position_definitions: List[EquipmentPositionInput],
    port_definitions: List[EquipmentPortInput],
) -> AddEquipmentTypeMutation.AddEquipmentTypeMutationData.EquipmentType:
    return AddEquipmentTypeMutation.execute(
        client,
        AddEquipmentTypeInput(
            name=name,
            category=category,
            positions=position_definitions,
            ports=port_definitions,
            properties=properties,
        ),
    )


def get_or_create_equipment_type(
    client: SymphonyClient,
    name: str,
    category: str,
    properties: List[PropertyDefinition],
    ports_dict: Dict[str, str],
    position_list: List[str],
) -> EquipmentType:
    """This function checks equipment type existence,
        in case it is not found, creates one.

        Args:
            name (str): equipment name
            category (str): category name
            properties (List[ `pyinventory.common.data_class.PropertyDefinition` ]): list of property definitions
            ports_dict (Dict[str, str]): dict of property name to property value
            - str - port name
            - str - port type name

            position_list (List[str]): list of positions names

        Returns:
            `pyinventory.common.data_class.EquipmentType` object

        Raises:
            FailedOperationException: internal inventory error

        Example:
            ```
            e_type = client.get_or_create_equipment_type(
                name="Tp-Link T1600G",
                category="Router",
                properties=[
                    PropertyDefinition(
                        property_name="IP",
                        property_kind=PropertyKind.string,
                        default_raw_value=None,
                        is_fixed=True
                    )
                ],
                ports_dict={"Port 1": "eth port", "port 2": "eth port"},
                position_list=[],
            )
            ```
    """
    if name in EQUIPMENT_TYPES:
        return EQUIPMENT_TYPES[name]
    return add_equipment_type(
        client, name, category, properties, ports_dict, position_list
    )


def _edit_equipment_type(
    client: SymphonyClient,
    equipment_type_id: str,
    name: str,
    category: Optional[str],
    properties: List[PropertyTypeInput],
    position_definitions: List[EquipmentPositionInput],
    port_definitions: List[EquipmentPortInput],
) -> EditEquipmentTypeMutation.EditEquipmentTypeMutationData.EquipmentType:
    return EditEquipmentTypeMutation.execute(
        client,
        EditEquipmentTypeInput(
            id=equipment_type_id,
            name=name,
            category=category,
            positions=position_definitions,
            ports=port_definitions,
            properties=properties,
        ),
    )


def _update_equipment_type(
    client: SymphonyClient,
    equipment_type_id: str,
    name: str,
    category: Optional[str],
    properties: List[PropertyTypeInput],
    position_definitions: List[EquipmentPositionInput],
    port_definitions: List[EquipmentPortInput],
) -> EquipmentType:

    equipment_type_result = _edit_equipment_type(
        client=client,
        equipment_type_id=equipment_type_id,
        name=name,
        category=category,
        properties=properties,
        position_definitions=position_definitions,
        port_definitions=port_definitions,
    )
    equipment_type = EquipmentType(
        name=equipment_type_result.name,
        category=equipment_type_result.category,
        id=equipment_type_result.id,
        property_types=format_to_property_definitions(
            equipment_type_result.propertyTypes
        ),
        position_definitions=equipment_type_result.positionDefinitions,
        port_definitions=equipment_type_result.portDefinitions,
    )
    EQUIPMENT_TYPES[name] = equipment_type
    return equipment_type


def add_equipment_type(
    client: SymphonyClient,
    name: str,
    category: str,
    properties: List[PropertyDefinition],
    ports_dict: Dict[str, str],
    position_list: List[str],
) -> EquipmentType:
    """This function creates new equipment type.

        Args:
            name (str): equipment type name
            category (str): category name
            properties (List[ `pyinventory.common.data_class.PropertyDefinition` ]): list of property definitions
            ports_dict (Dict[str, str]): dictionary of port name to port type name
            - str - port name
            - str - port type name

            position_list (List[str]): list of positions names

        Returns:
            `pyinventory.common.data_class.EquipmentType` object

        Raises:
            FailedOperationException: internal inventory error

        Example:
            ```
            e_type = client.add_equipment_type(
                name="Tp-Link T1600G",
                category="Router",
                properties=[
                    PropertyDefinition(
                        property_name="IP",
                        property_kind=PropertyKind.string,
                        default_raw_value=None,
                        is_fixed=True
                    )
                ],
                ports_dict={"Port 1": "eth port", "port 2": "eth port"},
                position_list=[],
            )
            ```
    """
    new_property_types = format_to_property_type_inputs(data=properties)

    port_definitions = [
        EquipmentPortInput(name=name, portTypeID=PORT_TYPES[_type].id)
        for name, _type in ports_dict.items()
    ]
    position_definitions = [
        EquipmentPositionInput(name=position) for position in position_list
    ]
    equipment_type_result = _add_equipment_type(
        client,
        name,
        category,
        new_property_types,
        position_definitions,
        port_definitions,
    )
    equipment_type = EquipmentType(
        name=equipment_type_result.name,
        category=equipment_type_result.category,
        id=equipment_type_result.id,
        property_types=format_to_property_definitions(
            equipment_type_result.propertyTypes
        ),
        position_definitions=equipment_type_result.positionDefinitions,
        port_definitions=equipment_type_result.portDefinitions,
    )
    EQUIPMENT_TYPES[name] = equipment_type
    return equipment_type


def edit_equipment_type(
    client: SymphonyClient,
    name: str,
    new_positions_list: List[str],
    new_ports_dict: Dict[str, str],
) -> EquipmentType:
    """Edit existing equipment type.

        Args:
            name (str): equipment type name
            new_positions_list (List[str]): new position list
            new_ports_dict (Dict[str, str]): dictionary of port name to port type name
            - str - port name
            - str - port type name

        Returns:
            `pyinventory.common.data_class.EquipmentType` object

        Raises:
            FailedOperationException: internal inventory error

        Example:
            ```
            edited_equipment = client.edit_equipment_type(
                name="Card",
                new_positions_list=[],
                new_ports_dict={"Port 5": "Z Cards Only (LS - DND)"}
            )
            ```
    """
    equipment_type = EQUIPMENT_TYPES[name]
    edited_property_types = [
        format_to_property_type_input(property_type)
        for property_type in equipment_type.property_types
    ]
    position_definitions = [
        get_position_definition_input(position_definition, is_new=False)
        for position_definition in equipment_type.position_definitions
    ] + [EquipmentPositionInput(name=position) for position in new_positions_list]
    port_definitions = [
        get_port_definition_input(port_definition, is_new=False)
        for port_definition in equipment_type.port_definitions
    ] + [
        EquipmentPortInput(name=name, portTypeID=PORT_TYPES[_type].id)
        for name, _type in new_ports_dict.items()
    ]

    return _update_equipment_type(
        client=client,
        equipment_type_id=equipment_type.id,
        name=equipment_type.name,
        category=equipment_type.category,
        properties=edited_property_types,
        position_definitions=position_definitions,
        port_definitions=port_definitions,
    )


def copy_equipment_type(
    client: SymphonyClient, curr_equipment_type_name: str, new_equipment_type_name: str
) -> EquipmentType:
    """Copy existing equipment type.

        Args:
            curr_equipment_type_name (str): existing equipment type name
            new_equipment_type_name (str): new equipment type name

        Returns:
            `pyinventory.common.data_class.EquipmentType` object

        Raises:
            FailedOperationException: internal inventory error

        Example:
            ```
            e_type = client.copy_equipment_type(
                curr_equipment_type_name="Card",
                new_equipment_type_name="External_Card",
            )
            ```
    """
    equipment_type = EQUIPMENT_TYPES[curr_equipment_type_name]

    new_property_types = [
        format_to_property_type_input(property_type)
        for property_type in equipment_type.property_types
    ]

    new_position_definitions = [
        get_position_definition_input(position_definition)
        for position_definition in equipment_type.position_definitions
    ]

    new_port_definitions = [
        get_port_definition_input(port_definition)
        for port_definition in equipment_type.port_definitions
    ]

    equipment_type_result = _add_equipment_type(
        client,
        new_equipment_type_name,
        equipment_type.category,
        new_property_types,
        new_position_definitions,
        new_port_definitions,
    )

    new_equipment_type = EquipmentType(
        name=equipment_type_result.name,
        category=equipment_type_result.category,
        id=equipment_type_result.id,
        property_types=format_to_property_definitions(
            equipment_type_result.propertyTypes
        ),
        position_definitions=equipment_type_result.positionDefinitions,
        port_definitions=equipment_type_result.portDefinitions,
    )

    EQUIPMENT_TYPES[new_equipment_type_name] = new_equipment_type
    return new_equipment_type


def get_equipment_type_property_type(
    client: SymphonyClient, equipment_type_name: str, property_type_id: str
) -> PropertyDefinition:
    """Get property type by ID on specific equipment type.

        Args:
            equipment_type_name (str): existing equipment type name
            property_type_id (str): property type ID

        Returns:
            `pyinventory.common.data_class.PropertyDefinition`  object

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: property type with id=`property_type_id` is not found

        Example:
            ```
            property_type = client.get_equipment_type_property_type(
                equipment_type_name="Card",
                property_type_id="12345",
            )
            ```
    """
    return get_property_type(
        client=client,
        entity_type=Entity.EquipmentType,
        entity_name=equipment_type_name,
        property_type_id=property_type_id,
    )


def get_equipment_type_property_type_by_external_id(
    client: SymphonyClient, equipment_type_name: str, property_type_external_id: str
) -> PropertyDefinition:
    """Get property type by external ID on specific equipment type.

        Args:
            equipment_type_name (str): existing equipment type name
            property_type_external_id (str): property type external ID

        Returns:
            `pyinventory.common.data_class.PropertyDefinition`  object

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: property type with external_id=`property_type_external_id` is not found

        Example:
            ```
            property_type = client.get_equipment_type_property_type_by_external_id(
                equipment_type_name="Card",
                property_type_external_id="12345",
            )
            ```
    """
    return get_property_type_by_external_id(
        client=client,
        entity_type=Entity.EquipmentType,
        entity_name=equipment_type_name,
        property_type_external_id=property_type_external_id,
    )


def edit_equipment_type_property_type(
    client: SymphonyClient,
    equipment_type_name: str,
    property_type_id: str,
    new_property_definition: PropertyDefinition,
) -> EquipmentType:
    """Edit specific property type on specific equipment type.

        Args:
            equipment_type_name (str): existing equipment type name
            property_type_id (str): existing property type id
            new_property_definition ( `pyinventory.common.data_class.PropertyDefinition` ): new property definition

        Returns:
            `pyinventory.common.data_class.EquipmentType`

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: if property type name is not found
            FailedOperationException: internal inventory error

        Example:
            ```
            e_type = client.edit_equipment_type_property_type(
                equipment_type_name="Card",
                property_type_id="111669149698",
                new_property_definition=PropertyDefinition(
                    property_name=property_type_name,
                    property_kind=PropertyKind.string,
                    default_raw_value=None,
                    is_fixed=False,
                    external_id="12345",
                ),
            )
            ```
    """
    equipment_type = EQUIPMENT_TYPES[equipment_type_name]
    edited_property_types = edit_property_type(
        client=client,
        entity_type=Entity.EquipmentType,
        entity_name=equipment_type_name,
        property_type_id=property_type_id,
        new_property_definition=new_property_definition,
    )
    position_definitions = [
        get_position_definition_input(position_definition, is_new=False)
        for position_definition in equipment_type.position_definitions
        if equipment_type.position_definitions
    ]
    port_definitions = [
        get_port_definition_input(port_definition, is_new=False)
        for port_definition in equipment_type.port_definitions
        if equipment_type.port_definitions
    ]

    return _update_equipment_type(
        client=client,
        equipment_type_id=equipment_type.id,
        name=equipment_type.name,
        category=equipment_type.category,
        properties=edited_property_types,
        position_definitions=position_definitions,
        port_definitions=port_definitions,
    )


def delete_equipment_type_with_equipments(
    client: SymphonyClient, equipment_type: EquipmentType
) -> None:
    """Delete equipment type with existing equipments.

        Args:
            equipment_type ( `pyinventory.common.data_class.EquipmentType` ): equipment type object

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: if equipment_type does not exist

        Example:
            ```
            equipment_type = client.get_or_create_equipment_type(
                name="Tp-Link T1600G",
                category="Router",
                properties=[("IP", "string", None, True)],
                ports_dict={"Port 1": "eth port", "port 2": "eth port"},
                position_list=[],
            )
            client.delete_equipment_type_with_equipments(equipment_type=equipment_type)
            ```
    """
    equipment_type_with_equipments = EquipmentTypeEquipmentQuery.execute(
        client, id=equipment_type.id
    )
    if not equipment_type_with_equipments:
        raise EntityNotFoundError(
            entity=Entity.EquipmentType, entity_id=equipment_type.id
        )
    for equipment in equipment_type_with_equipments.equipments:
        delete_equipment(
            client,
            Equipment(
                id=equipment.id,
                external_id=equipment.externalId,
                name=equipment.name,
                equipment_type_name=equipment.equipmentType.name,
            ),
        )

    RemoveEquipmentTypeMutation.execute(client, id=equipment_type.id)
