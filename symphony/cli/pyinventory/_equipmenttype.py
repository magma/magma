#!/usr/bin/env python3
# pyre-strict

from typing import Any, Dict, List, Tuple, Union

from dacite import from_dict
from gql.gql.client import OperationException

from ._utils import PropertyValue, _make_property_types
from .api.equipment import delete_equipment
from .consts import Equipment, EquipmentType
from .exceptions import EquipmentTypeNotFoundException
from .graphql.add_equipment_type_mutation import (
    AddEquipmentTypeInput,
    AddEquipmentTypeMutation,
    PropertyKind,
)
from .graphql.edit_equipment_type_mutation import (
    EditEquipmentTypeInput,
    EditEquipmentTypeMutation,
)
from .graphql.equipment_type_equipments_query import EquipmentTypeEquipmentQuery
from .graphql.equipment_types_query import EquipmentTypesQuery
from .graphql.remove_equipment_type_mutation import RemoveEquipmentTypeMutation
from .graphql_client import GraphqlClient
from .reporter import FailedOperationException


ADD_EQUIPMENT_TYPE_MUTATION_NAME = "addEquipmentType"
EDIT_EQUIPMENT_TYPE_MUTATION_NAME = "editEquipmentType"


def _populate_equipment_types(client: GraphqlClient) -> None:
    edges = EquipmentTypesQuery.execute(client).equipmentTypes.edges

    for edge in edges:
        node = edge.node
        client.equipmentTypes[node.name] = EquipmentType(
            name=node.name,
            category=node.category,
            id=node.id,
            propertyTypes=list(map(lambda p: p.to_dict(), node.propertyTypes)),
            positionDefinitions=list(
                map(lambda p: p.to_dict(), node.positionDefinitions)
            ),
            portDefinitions=list(map(lambda p: p.to_dict(), node.portDefinitions)),
        )


def _add_equipment_type(
    client: GraphqlClient,
    name: str,
    category: str,
    properties: List[Dict[str, Any]],
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
                    data_class=AddEquipmentTypeInput.EquipmentPositionInput, data=pos
                )
                for pos in position_definitions
            ],
            ports=[
                from_dict(
                    data_class=AddEquipmentTypeInput.EquipmentPortInput, data=port
                )
                for port in port_definitions
            ],
            properties=[
                from_dict(data_class=AddEquipmentTypeInput.PropertyTypeInput, data=prop)
                for prop in properties
            ],
        ),
    ).__dict__[ADD_EQUIPMENT_TYPE_MUTATION_NAME]


def get_or_create_equipment_type(
    client: GraphqlClient,
    name: str,
    category: str,
    properties: List[Tuple[str, str, PropertyValue, bool]],
    ports_dict: Dict[str, str],
    position_list: List[str],
) -> EquipmentType:
    if name in client.equipmentTypes:
        return client.equipmentTypes[name]
    return add_equipment_type(
        client, name, category, properties, ports_dict, position_list
    )


def _edit_equipment_type(
    client: GraphqlClient,
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
                    data_class=EditEquipmentTypeInput.EquipmentPositionInput, data=pos
                )
                for pos in position_definitions
            ],
            ports=[
                from_dict(
                    data_class=EditEquipmentTypeInput.EquipmentPortInput, data=port
                )
                for port in port_definitions
            ],
            properties=[
                from_dict(
                    data_class=EditEquipmentTypeInput.PropertyTypeInput, data=prop
                )
                for prop in properties
            ],
        ),
    ).__dict__[EDIT_EQUIPMENT_TYPE_MUTATION_NAME]


def add_equipment_type(
    client: GraphqlClient,
    name: str,
    category: str,
    properties: List[Tuple[str, str, PropertyValue, bool]],
    ports_dict: Dict[str, str],
    position_list: List[str],
) -> EquipmentType:
    property_types = _make_property_types(properties)

    def property_type_to_kind(
        key: str, value: Union[str, int, float, bool]
    ) -> Union[str, int, float, bool, PropertyKind]:
        return value if key != "type" else PropertyKind(value)

    new_property_types = [
        {k: property_type_to_kind(k, v) for k, v in property_type.items()}
        for property_type in property_types
    ]

    port_definitions = [
        {"name": name, "type": type} for name, type in ports_dict.items()
    ]
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
        propertyTypes=list(map(lambda p: p.to_dict(), equipment_type.propertyTypes)),
        positionDefinitions=list(
            map(lambda p: p.to_dict(), equipment_type.positionDefinitions)
        ),
        portDefinitions=list(
            map(lambda p: p.to_dict(), equipment_type.portDefinitions)
        ),
    )
    client.equipmentTypes[equipment_type.name] = equipment_type
    return equipment_type


def edit_equipment_type(
    client: GraphqlClient,
    name: str,
    new_positions_list: List[str],
    new_ports_dict: Dict[str, str],
) -> EquipmentType:
    if name not in client.equipmentTypes:
        raise EquipmentTypeNotFoundException
    equipment_type = client.equipmentTypes[name]
    position_definitions = equipment_type.positionDefinitions + [
        {"name": position} for position in new_positions_list
    ]
    port_definitions = equipment_type.portDefinitions + [
        {"name": name, "type": type} for name, type in new_ports_dict.items()
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
        propertyTypes=list(map(lambda p: p.to_dict(), equipment_type.propertyTypes)),
        positionDefinitions=list(
            map(lambda p: p.to_dict(), equipment_type.positionDefinitions)
        ),
        portDefinitions=list(
            map(lambda p: p.to_dict(), equipment_type.portDefinitions)
        ),
    )
    client.equipmentTypes[equipment_type.name] = equipment_type
    return equipment_type


def copy_equipment_type(
    client: GraphqlClient, curr_equipment_type_name: str, new_equipment_type_name: str
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
        propertyTypes=list(map(lambda p: p.to_dict(), equipment_type.propertyTypes)),
        positionDefinitions=list(
            map(lambda p: p.to_dict(), equipment_type.positionDefinitions)
        ),
        portDefinitions=list(
            map(lambda p: p.to_dict(), equipment_type.portDefinitions)
        ),
    )

    client.equipmentTypes[new_equipment_type_name] = new_equipment_type
    return new_equipment_type


def delete_equipment_type_with_equipments(
    client: GraphqlClient, equipment_type: EquipmentType
) -> None:
    equipments = EquipmentTypeEquipmentQuery.execute(
        client, id=equipment_type.id
    ).equipmentType.equipments
    for equipment in equipments:
        delete_equipment(client, Equipment(id=equipment.id, name=equipment.name))

    RemoveEquipmentTypeMutation.execute(client, id=equipment_type.id)
