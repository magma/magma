#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import Dict, List, Optional, Tuple

from gql.gql.client import OperationException
from gql.gql.reporter import FailedOperationException

from .._utils import (
    format_properties,
    get_port_definition_input,
    get_position_definition_input,
    get_property_type_input,
)
from ..client import SymphonyClient
from ..consts import Entity, Equipment, EquipmentPortType, EquipmentType, PropertyValue
from ..exceptions import EntityNotFoundError, EquipmentTypeNotFoundException
from ..graphql.add_equipment_type_input import AddEquipmentTypeInput
from ..graphql.add_equipment_type_mutation import AddEquipmentTypeMutation
from ..graphql.edit_equipment_type_input import EditEquipmentTypeInput
from ..graphql.edit_equipment_type_mutation import EditEquipmentTypeMutation
from ..graphql.equipment_port_input import EquipmentPortInput
from ..graphql.equipment_port_types import EquipmentPortTypesQuery
from ..graphql.equipment_position_input import EquipmentPositionInput
from ..graphql.equipment_type_equipments_query import EquipmentTypeEquipmentQuery
from ..graphql.equipment_types_query import EquipmentTypesQuery
from ..graphql.property_type_input import PropertyTypeInput
from ..graphql.remove_equipment_type_mutation import RemoveEquipmentTypeMutation
from .equipment import delete_equipment


ADD_EQUIPMENT_TYPE_MUTATION_NAME = "addEquipmentType"
EDIT_EQUIPMENT_TYPE_MUTATION_NAME = "editEquipmentType"


def _populate_equipment_types(client: SymphonyClient) -> None:
    edges = EquipmentTypesQuery.execute(client).equipmentTypes.edges

    for edge in edges:
        node = edge.node
        if node:
            client.equipmentTypes[node.name] = EquipmentType(
                name=node.name,
                category=node.category,
                id=node.id,
                property_types=node.propertyTypes,
                position_definitions=node.positionDefinitions,
                port_definitions=node.portDefinitions,
            )


def _populate_equipment_port_types(client: SymphonyClient) -> None:
    edges = EquipmentPortTypesQuery.execute(client).equipmentPortTypes.edges

    for edge in edges:
        node = edge.node
        if node:
            client.portTypes[node.name] = EquipmentPortType(
                id=node.id,
                name=node.name,
                property_types=node.propertyTypes,
                link_property_types=node.linkPropertyTypes,
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
    ).__dict__[ADD_EQUIPMENT_TYPE_MUTATION_NAME]


def get_or_create_equipment_type(
    client: SymphonyClient,
    name: str,
    category: str,
    properties: List[Tuple[str, str, Optional[PropertyValue], Optional[bool]]],
    ports_dict: Dict[str, str],
    position_list: List[str],
) -> EquipmentType:
    """This function checks equipment type existence, 
        in case it is not found, creates one.

        Args:
            name (str): equipment name
            category (str): category name
            properties (List[Tuple[str, str, Optional[PropertyValue], Optional[bool]]]): 
            - str - type name
            - str - enum["string", "int", "bool", "float", "date", "enum", "range", 
            "email", "gps_location", "equipment", "location", "service", "datetime_local"]
            - PropertyValue - default property value
            - bool - fixed value flag

            ports_dict (Dict[str, str]): dict of property name to property value
            - str - port name
            - str - port type name
            
            position_list (List[str]): list of positions names

        Returns:
            pyinventory.consts.EquipmentType object
        
        Raises:
            FailedOperationException: internal inventory error

        Example:
            ```
            e_type = client.get_or_create_equipment_type(
                name="Tp-Link T1600G",
                category="Router",
                properties=[("IP", "string", None, True)],
                ports_dict={"Port 1": "eth port", "port 2": "eth port"},
                position_list=[],
            )
            ```
    """
    if name in client.equipmentTypes:
        return client.equipmentTypes[name]
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
    ).__dict__[EDIT_EQUIPMENT_TYPE_MUTATION_NAME]


def add_equipment_type(
    client: SymphonyClient,
    name: str,
    category: str,
    properties: List[Tuple[str, str, Optional[PropertyValue], Optional[bool]]],
    ports_dict: Dict[str, str],
    position_list: List[str],
) -> EquipmentType:
    """This function creates new equipment type.

        Args:
            name (str): equipment type name
            category (str): category name
            properties (List[Tuple[str, str, Optional[PropertyValue], Optional[bool]]]): 
            - str - type name
            - str - enum["string", "int", "bool", "float", "date", "enum", "range", 
            "email", "gps_location", "equipment", "location", "service", "datetime_local"]
            - PropertyValue - default property value
            - bool - fixed value flag

            ports_dict (Dict[str, str]): dict of property name to property value
            - str - port name
            - str - port type name
            
            position_list (List[str]): list of positions names

        Returns:
            pyinventory.consts.EquipmentType object
        
        Raises:
            FailedOperationException: internal inventory error

        Example:
            ```
            e_type = client.add_equipment_type(
                name="Tp-Link T1600G",
                category="Router",
                properties=[("IP", "string", None, True)],
                ports_dict={"Port 1": "eth port", "port 2": "eth port"},
                position_list=[],
            )
            ```
    """
    new_property_types = format_properties(properties)

    port_definitions = [
        EquipmentPortInput(name=name, portTypeID=client.portTypes[_type].id)
        for name, _type in ports_dict.items()
    ]
    position_definitions = [
        EquipmentPositionInput(name=position) for position in position_list
    ]

    add_equipment_type_variables = {
        "name": name,
        "category": category,
        "positionDefinitions": position_definitions,
        "portDefinitions": port_definitions,
        "properties": new_property_types,
    }
    try:
        equipment_type = _add_equipment_type(
            client,
            name,
            category,
            new_property_types,
            position_definitions,
            port_definitions,
        )
        client.reporter.log_successful_operation(
            ADD_EQUIPMENT_TYPE_MUTATION_NAME, add_equipment_type_variables
        )
    except OperationException as e:
        raise FailedOperationException(
            client.reporter,
            e.err_msg,
            e.err_id,
            ADD_EQUIPMENT_TYPE_MUTATION_NAME,
            add_equipment_type_variables,
        )

    equipment_type = EquipmentType(
        name=equipment_type.name,
        category=equipment_type.category,
        id=equipment_type.id,
        property_types=equipment_type.propertyTypes,
        position_definitions=equipment_type.positionDefinitions,
        port_definitions=equipment_type.portDefinitions,
    )
    client.equipmentTypes[equipment_type.name] = equipment_type
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
            new_ports_dict (Dict[str, str]): ports dictionary, where key is a port name and value is port type name

        Returns:
            EquipmentType object

        Raises:
            FailedOperationException for internal inventory error

        Example:
            edited_equipment = client.edit_equipment_type("Card", [], {"Port 5": "Z Cards Only (LS - DND)"})
    """
    if name not in client.equipmentTypes:
        raise EquipmentTypeNotFoundException
    equipment_type = client.equipmentTypes[name]
    position_definitions = [
        get_position_definition_input(position_definition, is_new=False)
        for position_definition in equipment_type.position_definitions
    ] + [EquipmentPositionInput(name=position) for position in new_positions_list]
    port_definitions = [
        get_port_definition_input(port_definition, is_new=False)
        for port_definition in equipment_type.port_definitions
    ] + [
        EquipmentPortInput(name=name, portTypeID=client.portTypes[_type].id)
        for name, _type in new_ports_dict.items()
    ]

    edit_equipment_type_variables = {
        "name": name,
        "category": equipment_type.category,
        "positionDefinitions": position_definitions,
        "portDefinitions": port_definitions,
        "properties": equipment_type.property_types,
    }
    try:
        equipment_type = _edit_equipment_type(
            client,
            equipment_type.id,
            equipment_type.name,
            equipment_type.category,
            [
                get_property_type_input(property_type, is_new=False)
                for property_type in equipment_type.property_types
            ],
            position_definitions,
            port_definitions,
        )
        client.reporter.log_successful_operation(
            EDIT_EQUIPMENT_TYPE_MUTATION_NAME, edit_equipment_type_variables
        )
    except OperationException as e:
        raise FailedOperationException(
            client.reporter,
            e.err_msg,
            e.err_id,
            EDIT_EQUIPMENT_TYPE_MUTATION_NAME,
            edit_equipment_type_variables,
        )

    equipment_type = EquipmentType(
        name=equipment_type.name,
        category=equipment_type.category,
        id=equipment_type.id,
        property_types=equipment_type.propertyTypes,
        position_definitions=equipment_type.positionDefinitions,
        port_definitions=equipment_type.portDefinitions,
    )
    client.equipmentTypes[equipment_type.name] = equipment_type
    return equipment_type


def copy_equipment_type(
    client: SymphonyClient, curr_equipment_type_name: str, new_equipment_type_name: str
) -> EquipmentType:
    """Copy existing equipment type.

        Args:
            curr_equipment_type_name (str): existing equipment type name
            new_equipment_type_name (str): new equipment type name

        Returns:
            pyinventory.consts.EquipmentType object
        
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
    if curr_equipment_type_name not in client.equipmentTypes:
        raise Exception(
            "Equipment type " + curr_equipment_type_name + " does not exist"
        )

    equipment_type = client.equipmentTypes[curr_equipment_type_name]

    new_property_types = [
        get_property_type_input(property_type)
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

    equipment_type = _add_equipment_type(
        client,
        new_equipment_type_name,
        equipment_type.category,
        new_property_types,
        new_position_definitions,
        new_port_definitions,
    )

    new_equipment_type = EquipmentType(
        name=equipment_type.name,
        category=equipment_type.category,
        id=equipment_type.id,
        property_types=equipment_type.propertyTypes,
        position_definitions=equipment_type.positionDefinitions,
        port_definitions=equipment_type.portDefinitions,
    )

    client.equipmentTypes[new_equipment_type_name] = new_equipment_type
    return new_equipment_type


def delete_equipment_type_with_equipments(
    client: SymphonyClient, equipment_type: EquipmentType
) -> None:
    equipment_type_with_equipments = EquipmentTypeEquipmentQuery.execute(
        client, id=equipment_type.id
    ).equipmentType
    if not equipment_type_with_equipments:
        raise EntityNotFoundError(
            entity=Entity.EquipmentType, entity_id=equipment_type.id
        )
    for equipment in equipment_type_with_equipments.equipments:
        delete_equipment(
            client,
            Equipment(
                id=equipment.id,
                name=equipment.name,
                equipment_type_name=equipment.equipmentType.name,
            ),
        )

    RemoveEquipmentTypeMutation.execute(client, id=equipment_type.id)
