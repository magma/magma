#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from dataclasses import asdict
from typing import Any, Dict, List, Optional, Tuple

from dacite import Config, from_dict
from gql.gql.client import OperationException
from gql.gql.reporter import FailedOperationException

from .._utils import format_properties
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
                propertyTypes=list(map(lambda p: asdict(p), node.propertyTypes)),
                positionDefinitions=list(
                    map(lambda p: asdict(p), node.positionDefinitions)
                ),
                portDefinitions=list(map(lambda p: asdict(p), node.portDefinitions)),
            )


def _populate_equipment_port_types(client: SymphonyClient) -> None:
    edges = EquipmentPortTypesQuery.execute(client).equipmentPortTypes.edges

    for edge in edges:
        node = edge.node
        if node:
            client.portTypes[node.name] = EquipmentPortType(
                id=node.id,
                name=node.name,
                properties=list(map(lambda p: asdict(p), node.propertyTypes)),
                link_properties=list(map(lambda p: asdict(p), node.linkPropertyTypes)),
            )


def _add_equipment_type(
    client: SymphonyClient,
    name: str,
    category: str,
    properties: List[PropertyTypeInput],
    position_definitions: List[Dict[str, Any]],
    port_definitions: List[Dict[str, Any]],
) -> AddEquipmentTypeMutation.AddEquipmentTypeMutationData.EquipmentType:
    return AddEquipmentTypeMutation.execute(
        client,
        AddEquipmentTypeInput(
            name=name,
            category=category,
            positions=[
                from_dict(
                    data_class=EquipmentPositionInput,
                    data=pos,
                    config=Config(strict=True),
                )
                for pos in position_definitions
            ],
            ports=[
                from_dict(
                    data_class=EquipmentPortInput, data=port, config=Config(strict=True)
                )
                for port in port_definitions
            ],
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
    if name in client.equipmentTypes:
        return client.equipmentTypes[name]
    return add_equipment_type(
        client, name, category, properties, ports_dict, position_list
    )


def _edit_equipment_type(
    client: SymphonyClient,
    equipment_type_id: str,
    name: str,
    category: str,
    properties: List[Dict[str, Any]],
    position_definitions: List[Dict[str, Any]],
    port_definitions: List[Dict[str, Any]],
) -> EditEquipmentTypeMutation.EditEquipmentTypeMutationData.EquipmentType:
    return EditEquipmentTypeMutation.execute(
        client,
        EditEquipmentTypeInput(
            id=equipment_type_id,
            name=name,
            category=category,
            positions=[
                from_dict(
                    data_class=EquipmentPositionInput,
                    data=pos,
                    config=Config(strict=True),
                )
                for pos in position_definitions
            ],
            ports=[
                from_dict(
                    data_class=EquipmentPortInput, data=port, config=Config(strict=True)
                )
                for port in port_definitions
            ],
            properties=[
                from_dict(
                    data_class=PropertyTypeInput, data=prop, config=Config(strict=True)
                )
                for prop in properties
            ],
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

    new_property_types = format_properties(properties)

    port_definitions = [{"name": name} for name, _ in ports_dict.items()]
    position_definitions = [{"name": position} for position in position_list]

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
        propertyTypes=list(map(lambda p: asdict(p), equipment_type.propertyTypes)),
        positionDefinitions=list(
            map(lambda p: asdict(p), equipment_type.positionDefinitions)
        ),
        portDefinitions=list(map(lambda p: asdict(p), equipment_type.portDefinitions)),
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
    position_definitions = equipment_type.positionDefinitions + [
        {"name": position} for position in new_positions_list
    ]
    port_definitions = equipment_type.portDefinitions + [
        {"name": name, "portTypeID": client.portTypes[_type].id}
        for name, _type in new_ports_dict.items()
    ]

    edit_equipment_type_variables = {
        "name": name,
        "category": equipment_type.category,
        "positionDefinitions": position_definitions,
        "portDefinitions": port_definitions,
        "properties": equipment_type.propertyTypes,
    }
    try:
        equipment_type = _edit_equipment_type(
            client,
            equipment_type.id,
            equipment_type.name,
            equipment_type.category,
            equipment_type.propertyTypes,
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
        propertyTypes=list(map(lambda p: asdict(p), equipment_type.propertyTypes)),
        positionDefinitions=list(
            map(lambda p: asdict(p), equipment_type.positionDefinitions)
        ),
        portDefinitions=list(map(lambda p: asdict(p), equipment_type.portDefinitions)),
    )
    client.equipmentTypes[equipment_type.name] = equipment_type
    return equipment_type


def copy_equipment_type(
    client: SymphonyClient, curr_equipment_type_name: str, new_equipment_type_name: str
) -> EquipmentType:
    if curr_equipment_type_name not in client.equipmentTypes:
        raise Exception(
            "Equipment type " + curr_equipment_type_name + " does not exist"
        )

    equipment_type = client.equipmentTypes[curr_equipment_type_name]

    new_property_types = [
        {key: value for (key, value) in property_type.items() if key != "id"}
        for property_type in equipment_type.propertyTypes
    ]

    new_position_definitions = [
        {key: value for (key, value) in position_definition.items() if key != "id"}
        for position_definition in equipment_type.positionDefinitions
    ]

    new_port_definitions = [
        {key: value for (key, value) in port_definition.items() if key != "id"}
        for port_definition in equipment_type.portDefinitions
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
        propertyTypes=list(map(lambda p: asdict(p), equipment_type.propertyTypes)),
        positionDefinitions=list(
            map(lambda p: asdict(p), equipment_type.positionDefinitions)
        ),
        portDefinitions=list(map(lambda p: asdict(p), equipment_type.portDefinitions)),
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
        delete_equipment(client, Equipment(id=equipment.id, name=equipment.name))

    RemoveEquipmentTypeMutation.execute(client, id=equipment_type.id)
