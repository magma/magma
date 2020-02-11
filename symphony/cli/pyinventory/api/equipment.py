#!/usr/bin/env python3
# pyre-strict

from typing import Dict, List, Optional, Tuple

from gql.gql.client import OperationException
from tqdm import tqdm

from .._utils import PropertyValue, _get_graphql_properties, _get_property_value
from ..consts import Equipment, Location
from ..exceptions import (
    EquipmentIsNotUniqueException,
    EquipmentNotFoundException,
    EquipmentPositionIsNotUniqueException,
    EquipmentPositionNotFoundException,
)
from ..graphql.add_equipment_mutation import AddEquipmentInput, AddEquipmentMutation
from ..graphql.equipment_positions_query import EquipmentPositionsQuery
from ..graphql.equipment_search_query import EquipmentSearchQuery
from ..graphql.equipment_type_and_properties_query import (
    EquipmentTypeAndPropertiesQuery,
)
from ..graphql.location_equipments_query import LocationEquipmentsQuery
from ..graphql.remove_equipment_mutation import RemoveEquipmentMutation
from ..graphql_client import GraphqlClient
from ..reporter import FailedOperationException


ADD_EQUIPMENT_MUTATION_NAME = "addEquipment"
ADD_EQUIPMENT_TO_POSITION_MUTATION_NAME = "addEquipmentToPosition"
NUM_EQUIPMENTS_TO_SEARCH = 10


def _get_equipment_if_exists(
    client: GraphqlClient, name: str, location: Location
) -> Optional[Equipment]:

    equipments = LocationEquipmentsQuery.execute(
        client, id=location.id
    ).location.equipments

    equipments = [equipment for equipment in equipments if equipment.name == name]
    if len(equipments) > 1:
        raise EquipmentIsNotUniqueException(name)

    if len(equipments) == 0:
        return None
    return Equipment(equipments[0].name, equipments[0].id)


def get_equipment(client: GraphqlClient, name: str, location: Location) -> Equipment:
    """Get the equipment in a given location by name

        Args:
            name (str): equipment name
            location (client.Location object): retrieved from getLocation or
                                                addLocation api.

        Raises: AssertionException if location contains more than one equipments
                        with the same name or if equipment with the name is
                        not found FailedOperationException for internal
                        inventory error

        Returns: client.Equipment object (with name and id fields)
                 You can use the id to access the equipment from the UI:
                 https://{}.purpleheadband.cloud/inventory/inventory?equipment={}
    """

    equipment = _get_equipment_if_exists(client, name, location)
    if equipment is None:
        raise EquipmentNotFoundException(equipment_name=name)
    return equipment


def _get_equipment_in_position_if_exists(
    client: GraphqlClient, parent_equipment: Equipment, position_name: str
) -> Optional[Equipment]:
    _, equipment = _find_position_definition_id(client, parent_equipment, position_name)
    return equipment


def get_equipment_in_position(
    client: GraphqlClient, parent_equipment: Equipment, position_name: str
) -> Equipment:
    """Get the equipment attached in a given positionName of a given parentEquipment

        Args:
            parent_equipment (client.Equipment object): could be retrieved from
            the following apis:
                * getEquipment
                * getEquipmentInPosition
                * addEquipment
                * addEquipmentToPosition
            position_name (str): the name of the position in the equipment type.

        Raises: AssertionException if parent equipment has more than one
                    position with the given name, or none with this name or
                    if the position is not occupied._findPositionDefinitionId
                FailedOperationException for internal inventory error

        Returns: client.Equipment object (with name and id fields)
                 You can use the id to access the equipment from the UI:
                 https://{}.purpleheadband.cloud/inventory/inventory?equipment={}
    """

    equipment = _get_equipment_in_position_if_exists(
        client, parent_equipment, position_name
    )
    if equipment is None:
        raise EquipmentNotFoundException(
            parent_equipment_name=parent_equipment.name,
            parent_position_name=position_name,
        )
    return equipment


def add_equipment(
    client: GraphqlClient,
    name: str,
    equipment_type: str,
    location: Location,
    properties_dict: Dict[str, PropertyValue],
) -> Equipment:
    """Create a new equipment inside a given location. The equipment will be of the given equipment type
        , with the given name and with the given properties.
        If equipment with his name in this location already exists the existing equipment is returned

        Args:
            name (str): name of the new equipment
            equipment_type (str): name of the equipment type
            location (client.Location object): retrieved from getLocation or addLocation api.
            properties_dict: dict of property name to property value. the property value should match
                            the property type. Otherwise exception is raised

        Returns: client.Equipment object (with name and id fields)
                 You can use the id to access the equipment from the UI:
                 https://{}.purpleheadband.cloud/inventory/inventory?equipment={}

        Raises: AssertionException if location contains more than one equipments with the
                                    same name or if property value in propertiesDict does not match
                                    the property type
                FailedOperationException for internal inventory error

        Example:
            from datetime import date
            equipment = client.addEquipment(
                "Router X123",
                "Router",
                location,
                {
                    'Date Property ': date.today(),
                    'Lat/Lng Property: ': (-1.23,9.232),
                    'E-mail Property ': "user@fb.com",
                    'Number Property ': 11,
                    'String Property ': "aa",
                    'Float Property': 1.23
                })
    """

    property_types = client.equipmentTypes[equipment_type].propertyTypes
    properties = _get_graphql_properties(property_types, properties_dict)

    add_equipment_input = AddEquipmentInput(
        name=name,
        type=client.equipmentTypes[equipment_type].id,
        location=location.id,
        properties=properties,
    )

    try:
        equipment = AddEquipmentMutation.execute(client, add_equipment_input).__dict__[
            ADD_EQUIPMENT_MUTATION_NAME
        ]
        client.reporter.log_successful_operation(
            ADD_EQUIPMENT_MUTATION_NAME, add_equipment_input.__dict__
        )
    except OperationException as e:
        raise FailedOperationException(
            client.reporter,
            e.err_msg,
            e.err_id,
            ADD_EQUIPMENT_MUTATION_NAME,
            add_equipment_input.__dict__,
        )

    return Equipment(equipment.name, equipment.id)


def _find_position_definition_id(
    client: GraphqlClient, equipment: Equipment, position_name: str
) -> Tuple[str, Optional[Equipment]]:

    equipment_data = EquipmentPositionsQuery.execute(client, id=equipment.id).equipment

    positions = equipment_data.equipmentType.positionDefinitions
    existing_positions = equipment_data.positions

    positions = [position for position in positions if position.name == position_name]
    if len(positions) > 1:
        raise EquipmentPositionIsNotUniqueException(equipment.name, position_name)
    if len(positions) == 0:
        raise EquipmentPositionNotFoundException(equipment.name, position_name)
    position = positions[0]

    installed_positions = [
        existing_position
        for existing_position in existing_positions
        if existing_position.definition.name == position_name
    ]
    if len(installed_positions) > 1:
        raise EquipmentIsNotUniqueException(
            parent_equipment_name=equipment.name, parent_position_name=position_name
        )
    if (
        len(installed_positions) == 1
        and installed_positions[0].attachedEquipment is not None
    ):
        return (
            position.id,
            Equipment(
                id=installed_positions[0].attachedEquipment.id,
                name=installed_positions[0].attachedEquipment.name,
            ),
        )
    return position.id, None


def add_equipment_to_position(
    client: GraphqlClient,
    name: str,
    equipment_type: str,
    existing_equipment: Equipment,
    position_name: str,
    properties_dict: Dict[str, PropertyValue],
) -> Equipment:
    """Create a new equipment inside a given positionName of the given existingEquipment.
        The equipment will be of the given equipment type, with the given name and with the given properties.
        If equipment with his name in this position already exists the existing equipment is returned

        Args:
            name (str): name of the new equipment
            equipment_type (str): name of the equipment type
            existing_equipment (client.Equipment object): could be retrieved
            from the following apis:
                * getEquipment
                * getEquipmentInPosition
                * addEquipment
                * addEquipmentToPosition
            position_name (str): the name of the position in the equipment type.
            properties_dict: dict of property name to property value. the property value should match
                            the property type. Otherwise exception is raised

        Returns: client.Equipment object (with name and id fields)
                 You can use the id to access the equipment from the UI:
                 https://{}.purpleheadband.cloud/inventory/inventory?equipment={}

        Raises: AssertionException if parent equipment has more than one position with the given name
                            or if property value in propertiesDict does not match the property type
                FailedOperationException for internal inventory error

        Example:
            from datetime import date
            equipment = client.addEquipmentToPosition(
                "Card Y123",
                "Card",
                equipment,
                "Pos 1",
                {
                    'Date Property ': date.today(),
                    'Lat/Lng Property: ': (-1.23,9.232),
                    'E-mail Property ': "user@fb.com",
                    'Number Property ': 11,
                    'String Property ': "aa",
                    'Float Property': 1.23
                })
    """

    position_definition_id, _ = _find_position_definition_id(
        client, existing_equipment, position_name
    )
    property_types = client.equipmentTypes[equipment_type].propertyTypes
    properties = _get_graphql_properties(property_types, properties_dict)

    add_equipment_input = AddEquipmentInput(
        name=name,
        type=client.equipmentTypes[equipment_type].id,
        parent=existing_equipment.id,
        positionDefinition=position_definition_id,
        properties=properties,
    )

    try:
        equipment = AddEquipmentMutation.execute(client, add_equipment_input).__dict__[
            ADD_EQUIPMENT_MUTATION_NAME
        ]
        client.reporter.log_successful_operation(
            ADD_EQUIPMENT_TO_POSITION_MUTATION_NAME, add_equipment_input.__dict__
        )
    except OperationException as e:
        raise FailedOperationException(
            client.reporter,
            e.err_msg,
            e.err_id,
            ADD_EQUIPMENT_TO_POSITION_MUTATION_NAME,
            add_equipment_input.__dict__,
        )

    return Equipment(equipment.name, equipment.id)


def delete_equipment(client: GraphqlClient, equipment: Equipment) -> None:
    RemoveEquipmentMutation.execute(client, id=equipment.id)


def search_for_equipments(
    client: GraphqlClient, limit: int
) -> Tuple[List[Equipment], int]:

    equipments = EquipmentSearchQuery.execute(
        client, filters=[], limit=limit
    ).equipmentSearch

    total_count = equipments.count
    equipments = [
        Equipment(id=equipment.id, name=equipment.name)
        for equipment in equipments.equipment
    ]
    return equipments, total_count


def delete_all_equipments(client: GraphqlClient) -> None:
    equipments, total_count = search_for_equipments(client, NUM_EQUIPMENTS_TO_SEARCH)

    for equipment in equipments:
        delete_equipment(client, equipment)

    if total_count == len(equipments):
        return

    with tqdm(total=total_count) as progress_bar:
        progress_bar.update(len(equipments))
        while len(equipments) != 0:
            equipments, _ = search_for_equipments(client, NUM_EQUIPMENTS_TO_SEARCH)
            for equipment in equipments:
                delete_equipment(client, equipment)
            progress_bar.update(len(equipments))


def _get_equipment_type_and_properties_dict(
    client: GraphqlClient, equipment: Equipment
) -> Tuple[str, Dict[str, PropertyValue]]:

    result = EquipmentTypeAndPropertiesQuery.execute(client, id=equipment.id).equipment

    equipment_type = result.equipmentType.name

    properties_dict = {}
    property_types = client.equipmentTypes[equipment_type].propertyTypes
    for property in result.properties:
        property_type_id = property.propertyType.id
        property_types_with_id = [
            property_type
            for property_type in property_types
            if property_type["id"] == property_type_id
        ]
        assert (
            len(property_types_with_id) == 1
        ), "Equipment type {} has two property types with same id {}".format(
            equipment_type, property_type_id
        )
        property_type = property_types_with_id[0]
        property_value = _get_property_value(property_type, property.to_dict())
        properties_dict[property_type["name"]] = property_value
    return equipment_type, properties_dict


def copy_equipment_in_position(
    client: GraphqlClient,
    equipment: Equipment,
    dest_parent_equipment: Equipment,
    dest_position_name: str,
) -> Equipment:
    equipment_type, properties_dict = _get_equipment_type_and_properties_dict(
        client, equipment
    )
    return add_equipment_to_position(
        client,
        equipment.name,
        equipment_type,
        dest_parent_equipment,
        dest_position_name,
        properties_dict,
    )


def copy_equipment(
    client: GraphqlClient, equipment: Equipment, dest_location: Location
) -> Equipment:
    equipment_type, properties_dict = _get_equipment_type_and_properties_dict(
        client, equipment
    )
    return add_equipment(
        client, equipment.name, equipment_type, dest_location, properties_dict
    )


def get_equipment_type_of_equipment(
    client: GraphqlClient, equipment: Equipment
) -> Equipment:
    equipment_type, _ = _get_equipment_type_and_properties_dict(client, equipment)
    return client.equipmentTypes[equipment_type]


def get_or_create_equipment(
    client: GraphqlClient,
    name: str,
    equipment_type: str,
    location: Location,
    properties_dict: Dict[str, PropertyValue],
) -> Equipment:
    equipment = _get_equipment_if_exists(client, name, location)
    if equipment is not None:
        return equipment
    return add_equipment(client, name, equipment_type, location, properties_dict)


def get_or_create_equipment_in_position(
    client: GraphqlClient,
    name: str,
    equipment_type: str,
    existing_equipment: Equipment,
    position_name: str,
    properties_dict: Dict[str, PropertyValue],
) -> Equipment:
    equipment = _get_equipment_in_position_if_exists(
        client, existing_equipment, position_name
    )
    if equipment is not None:
        return equipment

    return add_equipment_to_position(
        client, name, equipment_type, existing_equipment, position_name, properties_dict
    )
