#!/usr/bin/env python3

from typing import Dict, List, Mapping, Optional, Tuple

from gql.gql.client import OperationException
from gql.gql.reporter import FailedOperationException
from tqdm import tqdm

from .._utils import PropertyValue, _get_property_value, get_graphql_property_inputs
from ..client import SymphonyClient
from ..consts import Entity, Equipment, EquipmentType, Location
from ..exceptions import (
    EntityNotFoundError,
    EquipmentIsNotUniqueException,
    EquipmentNotFoundException,
    EquipmentPositionIsNotUniqueException,
    EquipmentPositionNotFoundException,
)
from ..graphql.add_equipment_input import AddEquipmentInput
from ..graphql.add_equipment_mutation import AddEquipmentMutation
from ..graphql.edit_equipment_input import EditEquipmentInput
from ..graphql.edit_equipment_mutation import EditEquipmentMutation
from ..graphql.equipment_positions_query import EquipmentPositionsQuery
from ..graphql.equipment_search_query import EquipmentSearchQuery
from ..graphql.equipment_type_and_properties_query import (
    EquipmentTypeAndPropertiesQuery,
)
from ..graphql.equipment_type_equipments_query import EquipmentTypeEquipmentQuery
from ..graphql.location_equipments_query import LocationEquipmentsQuery
from ..graphql.property_kind_enum import PropertyKind
from ..graphql.remove_equipment_mutation import RemoveEquipmentMutation


ADD_EQUIPMENT_MUTATION_NAME = "addEquipment"
ADD_EQUIPMENT_TO_POSITION_MUTATION_NAME = "addEquipmentToPosition"
EDIT_EQUIPMENT_MUTATION_NAME = "editEquipment"
NUM_EQUIPMENTS_TO_SEARCH = 10


def _get_equipment_if_exists(
    client: SymphonyClient, name: str, location: Location
) -> Optional[Equipment]:

    location_with_equipments = LocationEquipmentsQuery.execute(
        client, id=location.id
    ).location
    if location_with_equipments is None:
        raise EntityNotFoundError(entity=Entity.Location, entity_id=location.id)
    equipments = [
        equipment
        for equipment in location_with_equipments.equipments
        if equipment.name == name
    ]
    if len(equipments) > 1:
        raise EquipmentIsNotUniqueException(name)

    if len(equipments) == 0:
        return None
    return Equipment(
        id=equipments[0].id,
        name=equipments[0].name,
        equipment_type_name=equipments[0].equipmentType.name,
    )


def get_equipment(client: SymphonyClient, name: str, location: Location) -> Equipment:
    """Get equipment by name in a given location.

        Args:
            name (str): equipment name
            location (pyinventory.consts.Location object): location object could be retrieved from 
            - `pyinventory.api.location.get_location`
            - `pyinventory.api.location.add_location`

        Returns:
            pyinventory.consts.Equipment object: 
                You can use the ID to access the equipment from the UI:
                https://{}.thesymphony.cloud/inventory/inventory?equipment={}

        Raises:
            EquipmentIsNotUniqueException: location contains 
                more than one equipment with the same name
            EquipmentNotFoundException: the equipment was not found
            FailedOperationException: internal inventory error

        Example:
            ```
            location = client.get_location([("Country", "LS_IND_Prod_Copy")])
            equipment = client.get_equipment("indProdCpy1_AIO", location)
            ```
    """

    equipment = _get_equipment_if_exists(client, name, location)
    if equipment is None:
        raise EquipmentNotFoundException(equipment_name=name)
    return equipment


def get_equipment_properties(
    client: SymphonyClient, equipment: Equipment
) -> Dict[str, PropertyValue]:
    """Get specific equipment properties.

        Args:
            equipment (pyinventory.consts.Equipment object): equipment object

        Returns:
            Dict[str, PropertyValue]: dict of property name to property value
            - str - property name
            - PropertyValue - new value of the same type for this property

        Example:
            ```
            location = client.get_location({("Country", "LS_IND_Prod_Copy")})
            equipment = client.get_equipment("indProdCpy1_AIO", location) 
            properties = client.get_equipment_properties(equipment=equipment)
            ```
    """
    equipment_type, properties_dict = _get_equipment_type_and_properties_dict(
        client, equipment
    )
    return properties_dict


def get_equipments_by_type(
    client: SymphonyClient, equipment_type_id: str
) -> List[Equipment]:
    """Get equipments by ID of specific type.

        Args:
            equipment_type_id (str): equipment type ID

        Returns:
            List[ pyinventory.consts.Equipment ]: List of found equipments

        Raises:
            EntityNotFoundError: equipment type with this ID does not exist

        Example:
            ```
            equipments = client.get_equipments_by_type(equipment_type_id="34359738369") 
            ```
    """
    equipment_type_with_equipments = EquipmentTypeEquipmentQuery.execute(
        client, id=equipment_type_id
    ).equipmentType
    if not equipment_type_with_equipments:
        raise EntityNotFoundError(
            entity=Entity.EquipmentType, entity_id=equipment_type_id
        )
    result = []
    for equipment in equipment_type_with_equipments.equipments:
        result.append(
            Equipment(
                id=equipment.id,
                name=equipment.name,
                equipment_type_name=equipment.equipmentType.name,
            )
        )

    return result


def _get_equipment_in_position_if_exists(
    client: SymphonyClient, parent_equipment: Equipment, position_name: str
) -> Optional[Equipment]:
    _, equipment = _find_position_definition_id(client, parent_equipment, position_name)
    return equipment


def get_equipment_in_position(
    client: SymphonyClient, parent_equipment: Equipment, position_name: str
) -> Equipment:
    """Get the equipment attached in a given positionName of a given parentEquipment

        Args:
            parent_equipment (pyinventory.consts.Equipment object): could be retrieved from
            - `pyinventory.api.equipment.get_equipment`
            - `pyinventory.api.equipment.get_equipment_in_position`
            - `pyinventory.api.equipment.add_equipment`
            - `pyinventory.api.equipment.add_equipment_to_position`

            position_name (str): position name

        Returns:
            pyinventory.consts.Equipment object: 
                You can use the ID to access the equipment from the UI:
                https://{}.thesymphony.cloud/inventory/inventory?equipment={}

        Raises:
            AssertionException: if parent equipment has more than one
                position with the given name, or none with this name or
                if the position is not occupied._findPositionDefinitionId
            FailedOperationException: for internal inventory error
            `pyinventory.exceptions.EntityNotFoundError`: if parent_equipment does not exist

        Example:
            ```
            location = client.get_location([("Country", "LS_IND_Prod_Copy")])
            p_equipment = client.get_equipment("indProdCpy1_AIO", location)
            equipment = client.get_equipment_in_position(p_equipment, "some_position")
            ```
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
    client: SymphonyClient,
    name: str,
    equipment_type: str,
    location: Location,
    properties_dict: Mapping[str, PropertyValue],
) -> Equipment:
    """Create a new equipment in a given location. 
        The equipment will be of the given equipment type, 
        with the given name and with the given properties.
        If equipment with his name in this location already exists, 
        the existing equipment is returned

        Args:
            name (str): new equipment name
            equipment_type (str): equipment type name
            location (pyinventory.consts.Location object): location object could be retrieved from 
            - `pyinventory.api.location.get_location`
            - `pyinventory.api.location.add_location`
            
            properties_dict (Mapping[str, PropertyValue]): dict of property name to property value
            - str - property name
            - PropertyValue - new value of the same type for this property

        Returns:
            pyinventory.consts.Equipment object: 
                You can use the ID to access the equipment from the UI:
                https://{}.thesymphony.cloud/inventory/inventory?equipment={}

        Raises:
            AssertionException: location contains more than one equipments with the
                same name or if property value in properties_dict does not match the property type
            FailedOperationException: internal inventory error

        Example:
            ```
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
            ```
    """

    property_types = client.equipmentTypes[equipment_type].property_types
    properties = get_graphql_property_inputs(property_types, properties_dict)

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

    return Equipment(
        id=equipment.id,
        name=equipment.name,
        equipment_type_name=equipment.equipmentType.name,
    )


def edit_equipment(
    client: SymphonyClient,
    equipment: Equipment,
    new_name: Optional[str] = None,
    new_properties: Optional[Dict[str, PropertyValue]] = None,
) -> Equipment:
    """Edit existing equipment.

        Args:
            equipment (pyinventory.consts.Equipment object): equipment object
            new_name (Optional[str]): equipment new name
            new_properties (Optional[Dict[str, pyinventory.consts.PropertyValue]]): Dict, where
                str - property name
                PropertyValue - new value of the same type for this property

        Returns:
            pyinventory.consts.Equipment object

        Raises:
            FailedOperationException: internal inventory error

        Example:
            ```
            location = client.get_location({("Country", "LS_IND_Prod_Copy")})
            equipment = client.get_equipment("indProdCpy1_AIO", location) 
            edited_equipment = client.edit_equipment(equipment=equipment, new_name="new_name", new_properties={"Z AIO - Number": 123})
            ```
    """
    properties = []
    property_types = client.equipmentTypes[equipment.equipment_type_name].property_types
    if new_properties:
        properties = get_graphql_property_inputs(property_types, new_properties)
    edit_equipment_input = EditEquipmentInput(
        id=equipment.id,
        name=new_name if new_name else equipment.name,
        properties=properties,
    )

    try:
        result = EditEquipmentMutation.execute(client, edit_equipment_input).__dict__[
            EDIT_EQUIPMENT_MUTATION_NAME
        ]
        client.reporter.log_successful_operation(
            EDIT_EQUIPMENT_MUTATION_NAME, edit_equipment_input.__dict__
        )

    except OperationException as e:
        raise FailedOperationException(
            client.reporter,
            e.err_msg,
            e.err_id,
            EDIT_EQUIPMENT_MUTATION_NAME,
            edit_equipment_input.__dict__,
        )
    return Equipment(
        id=result.id, name=result.name, equipment_type_name=result.equipmentType.name
    )


def _find_position_definition_id(
    client: SymphonyClient, equipment: Equipment, position_name: str
) -> Tuple[str, Optional[Equipment]]:

    equipment_data = EquipmentPositionsQuery.execute(client, id=equipment.id).equipment

    if not equipment_data:
        raise EntityNotFoundError(entity=Entity.Equipment, entity_id=equipment.id)

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
    if len(installed_positions) == 1:
        attached_equipment = installed_positions[0].attachedEquipment
        if attached_equipment is not None:
            return (
                position.id,
                Equipment(
                    id=attached_equipment.id,
                    name=attached_equipment.name,
                    equipment_type_name=attached_equipment.equipmentType.name,
                ),
            )
    return position.id, None


def add_equipment_to_position(
    client: SymphonyClient,
    name: str,
    equipment_type: str,
    existing_equipment: Equipment,
    position_name: str,
    properties_dict: Mapping[str, PropertyValue],
) -> Equipment:
    """Create a new equipment inside a given positionName of the given existingEquipment.
        The equipment will be of the given equipment type, with the given name and with the given properties.
        If equipment with his name in this position already exists the existing equipment is returned

        Args:
            name (str): new equipment name
            equipment_type (str): equipment type name
            existing_equipment (pyinventory.consts.Equipment object): could be retrieved from
            - `pyinventory.api.equipment.get_equipment`
            - `pyinventory.api.equipment.get_equipment_in_position`
            - `pyinventory.api.equipment.add_equipment`
            - `pyinventory.api.equipment.add_equipment_to_position`
            
            position_name (str): position name in the equipment type.            
            properties_dict (Mapping[str, PropertyValue]): dict of property name to property value
            - str - property name
            - PropertyValue - new value of the same type for this property

        Returns:
            pyinventory.consts.Equipment object: 
                You can use the ID to access the equipment from the UI:
                https://{}.thesymphony.cloud/inventory/inventory?equipment={}

        Raises:
            AssertionException: if parent equipment has more than one position with the given name
                            or if property value in propertiesDict does not match the property type
            FailedOperationException: for internal inventory error
            `pyinventory.exceptions.EntityNotFoundError`: if existing_equipment does not exist

        Example:
            ```
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
            ```
    """

    position_definition_id, _ = _find_position_definition_id(
        client, existing_equipment, position_name
    )
    property_types = client.equipmentTypes[equipment_type].property_types
    properties = get_graphql_property_inputs(property_types, properties_dict)

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

    return Equipment(
        id=equipment.id,
        name=equipment.name,
        equipment_type_name=equipment.equipmentType.name,
    )


def delete_equipment(client: SymphonyClient, equipment: Equipment) -> None:
    """This function delete Equipment.
        
        Args:
            equipment (pyinventory.consts.Equipment object): equipment object
        
        Example:
            ```
            client.delete_equipment(equipment) 
            ```
    """
    RemoveEquipmentMutation.execute(client, id=equipment.id)


def search_for_equipments(
    client: SymphonyClient, limit: int
) -> Tuple[List[Equipment], int]:
    """Search for equipments.

        Args:
            limit (int): search result limit

        Returns:
            Tuple[List[ `pyinventory.consts.Equipment` , int]

        Example:
            ```
            client.search_for_equipments(10)
            ```
    """
    equipments = EquipmentSearchQuery.execute(
        client, filters=[], limit=limit
    ).equipmentSearch

    total_count = equipments.count
    equipments = [
        Equipment(
            id=equipment.id,
            name=equipment.name,
            equipment_type_name=equipment.equipmentType.name,
        )
        for equipment in equipments.equipment
    ]
    return equipments, total_count


def delete_all_equipments(client: SymphonyClient) -> None:
    """This function delete all Equipments.
        
        Example:
            ```
            client.delete_all_equipment() 
            ```
    """
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
    client: SymphonyClient, equipment: Equipment
) -> Tuple[str, Dict[str, PropertyValue]]:

    result = EquipmentTypeAndPropertiesQuery.execute(client, id=equipment.id).equipment
    if result is None:
        raise EntityNotFoundError(entity=Entity.Equipment, entity_id=equipment.id)
    equipment_type = result.equipmentType.name

    properties_dict = {}
    property_types = client.equipmentTypes[equipment_type].property_types
    for property in result.properties:
        property_type_id = property.propertyType.id
        property_types_with_id = [
            property_type
            for property_type in property_types
            if property_type.id == property_type_id
        ]
        assert (
            len(property_types_with_id) == 1
        ), "Equipment type {} has two property types with same id {}".format(
            equipment_type, property_type_id
        )
        property_type = property_types_with_id[0]
        property_value = _get_property_value(property_type.type.value, property)
        if property_type.type == PropertyKind.gps_location:
            properties_dict[property_type.name] = (property_value[0], property_value[1])
        else:
            properties_dict[property_type.name] = property_value[0]
    return equipment_type, properties_dict


def copy_equipment_in_position(
    client: SymphonyClient,
    equipment: Equipment,
    dest_parent_equipment: Equipment,
    dest_position_name: str,
) -> Equipment:
    """Copy equipment in position.

        Args:
            equipment (pyinventory.consts.Equipment object): equipment object to be copied
            dest_parent_equipment (pyinventory.consts.Equipment object): parent equipment, destination to copy to
            dest_position_name (str): destination position name

        Returns:
            pyinventory.consts.Equipment object

        Example:
            ```
            location = client.get_location({("Country", "LS_IND_Prod_Copy")})
            equipment_to_copy = client.get_equipment("indProdCpy1_AIO", location) 
            parent_equipment = client.get_equipment("parent", location) 
            copied_equipment = client.copy_equipment_in_position(equipment=equipment, dest_parent_equipment=parent_equipment, dest_position_name="destination position name")
            ```
    """
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
    client: SymphonyClient, equipment: Equipment, dest_location: Location
) -> Equipment:
    """Copy equipment.

        Args:
            equipment (pyinventory.consts.Equipment object): equipment object to be copied
            dest_location (pyinventory.consts.Location): destination locatoin to copy to

        Returns:
            pyinventory.consts.Equipment object

        Example:
            ```
            location = client.get_location({("Country", "LS_IND_Prod_Copy")})
            equipment = client.get_equipment("indProdCpy1_AIO", location)
            new_location = client.get_location({("Country", "LS_IND_Prod")})
            copied_equipment = client.copy_equipment(equipment=equipment, dest_location=new_location)
            ```
    """
    equipment_type, properties_dict = _get_equipment_type_and_properties_dict(
        client, equipment
    )
    return add_equipment(
        client, equipment.name, equipment_type, dest_location, properties_dict
    )


def get_equipment_type_of_equipment(
    client: SymphonyClient, equipment: Equipment
) -> EquipmentType:
    """This function returns equipment type object of equipment.

        Args:
            equipment (pyinventory.consts.Equipment object): equipment object

        Returns:
            pyinventory.consts.EquipmentType object

        Example:
            ```
            location = client.get_location({("Country", "LS_IND_Prod_Copy")})
            equipment = client.get_equipment("indProdCpy1_AIO", location) 
            equipment_type = client.get_equipment_type_of_equipment(equipment=equipment)
            ```
    """
    equipment_type, _ = _get_equipment_type_and_properties_dict(client, equipment)
    return client.equipmentTypes[equipment_type]


def get_or_create_equipment(
    client: SymphonyClient,
    name: str,
    equipment_type: str,
    location: Location,
    properties_dict: Mapping[str, PropertyValue],
) -> Equipment:
    """This function checks if equipment existence in specific location by name, 
        in case it is not found by name, creates one.

        Args:
            name (str): equipment name
            equipment_type (str): equipment type name
            location (pyinventory.consts.Location object): location object could be retrieved from 
            - `pyinventory.api.location.get_location`
            - `pyinventory.api.location.add_location`            
            properties_dict (Mapping[str, PropertyValue]): dict of property name to property value
            - str - property name
            - PropertyValue - new value of the same type for this property

        Returns:
            pyinventory.consts.Equipment object
        
        Raises:
            AssertionException: location contains more than one equipments with the
                same name or if property value in properties_dict does not match the property type
            FailedOperationException: internal inventory error

        Example:
            ```
            location = client.get_location({("Country", "LS_IND_Prod_Copy")})
            equipment = client.get_or_create_equipment(
                name="indProdCpy1_AIO", 
                equipment_type="router", 
                location=location,
                properties_dict={
                    'Date Property ': date.today(),
                    'Lat/Lng Property: ': (-1.23,9.232),
                    'E-mail Property ': "user@fb.com",
                    'Number Property ': 11,
                    'String Property ': "aa",
                    'Float Property': 1.23
                })
            ```
    """
    equipment = _get_equipment_if_exists(client, name, location)
    if equipment is not None:
        return equipment
    return add_equipment(client, name, equipment_type, location, properties_dict)


def get_or_create_equipment_in_position(
    client: SymphonyClient,
    name: str,
    equipment_type: str,
    existing_equipment: Equipment,
    position_name: str,
    properties_dict: Mapping[str, PropertyValue],
) -> Equipment:
    """This function checks equipment existence in specific location by name, 
        in case it is not found by name, creates one.

        Args:
            name (str): equipment name
            equipment_type (str): equipment type name
            existing_equipment (pyinventory.consts.Equipment object): existing equipment
            position_name (str): position name
            properties_dict (Mapping[str, PropertyValue]): dict of property name to property value
            - str - property name
            - PropertyValue - new value of the same type for this property

        Returns:
            pyinventory.consts.Equipment object
        
        Raises:
            AssertionException: location contains more than one equipments with the
                same name or if property value in properties_dict does not match the property type
            FailedOperationException: internal inventory error

        Example:
            ```
            location = client.get_location({("Country", "LS_IND_Prod_Copy")})
            e_equipment = client.get_equipment("indProdCpy1_AIO", location) 
            equipment_in_position = client.get_or_create_equipment_in_position(
                name="indProdCpy1_AIO", 
                equipment_type="router", 
                existing_equipment=e_equipment,
                position_name="some_position",
                properties_dict={
                    'Date Property ': date.today(),
                    'Lat/Lng Property: ': (-1.23,9.232),
                    'E-mail Property ': "user@fb.com",
                    'Number Property ': 11,
                    'String Property ': "aa",
                    'Float Property': 1.23
                })
            ```
    """
    equipment = _get_equipment_in_position_if_exists(
        client, existing_equipment, position_name
    )
    if equipment is not None:
        return equipment

    return add_equipment_to_position(
        client, name, equipment_type, existing_equipment, position_name, properties_dict
    )
