#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import Dict, List, Mapping, Optional, Tuple

from gql.gql.client import OperationException
from gql.gql.reporter import FailedOperationException
from tqdm import tqdm

from .._utils import PropertyValue, _get_property_value, get_graphql_property_inputs
from ..client import SymphonyClient
from ..common.constant import EQUIPMENTS_TO_SEARCH
from ..common.data_class import Equipment, EquipmentType, Location
from ..common.data_enum import Entity
from ..common.mutation_name import (
    ADD_EQUIPMENT,
    ADD_EQUIPMENT_TO_POSITION,
    EDIT_EQUIPMENT,
)
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
from ..graphql.equipment_filter_input import EquipmentFilterInput
from ..graphql.equipment_filter_type_enum import EquipmentFilterType
from ..graphql.equipment_positions_query import EquipmentPositionsQuery
from ..graphql.equipment_search_query import EquipmentSearchQuery
from ..graphql.equipment_type_and_properties_query import (
    EquipmentTypeAndPropertiesQuery,
)
from ..graphql.equipment_type_equipments_query import EquipmentTypeEquipmentQuery
from ..graphql.filter_operator_enum import FilterOperator
from ..graphql.location_equipments_query import LocationEquipmentsQuery
from ..graphql.property_kind_enum import PropertyKind
from ..graphql.remove_equipment_mutation import RemoveEquipmentMutation


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
        external_id=equipments[0].externalId,
        name=equipments[0].name,
        equipment_type_name=equipments[0].equipmentType.name,
    )


def get_equipment(client: SymphonyClient, name: str, location: Location) -> Equipment:
    """Get equipment by name in a given location.

        Args:
            name (str): equipment name
            location ( `pyinventory.common.data_class.Location`): location object could be retrieved from
            - `pyinventory.api.location.get_location`
            - `pyinventory.api.location.add_location`

        Returns:
            `pyinventory.common.data_class.Equipment` object:
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


def get_equipment_by_external_id(client: SymphonyClient, external_id: str) -> Equipment:
    """Get equipment by external ID.

        Args:
            external_id (str): equipment external ID

        Returns:
            `pyinventory.common.data_class.Equipment` object:
                You can use the ID to access the equipment from the UI:
                https://{}.thesymphony.cloud/inventory/inventory?equipment={}

        Raises:
            EquipmentIsNotUniqueException: location contains
                more than one equipment with the same external ID
            `pyinventory.exceptions.EntityNotFoundError`: the equipment was not found
            FailedOperationException: internal inventory error

        Example:
            ```
            equipment = client.get_equipment_by_external_id(external_id="123456")
            ```
    """
    equipment_filter = EquipmentFilterInput(
        filterType=EquipmentFilterType.EQUIP_INST_EXTERNAL_ID,
        operator=FilterOperator.IS,
        stringValue=external_id,
        idSet=[],
        stringSet=[],
    )

    equipments = EquipmentSearchQuery.execute(
        client, filters=[equipment_filter], limit=5
    ).equipmentSearch

    if not equipments or equipments.count == 0:
        raise EntityNotFoundError(
            entity=Entity.Equipment, msg=f"external_id={external_id}"
        )

    if equipments.count > 1:
        raise EquipmentIsNotUniqueException(external_id)

    return Equipment(
        id=equipments.equipment[0].id,
        external_id=equipments.equipment[0].externalId,
        name=equipments.equipment[0].name,
        equipment_type_name=equipments.equipment[0].equipmentType.name,
    )


def get_equipment_properties(
    client: SymphonyClient, equipment: Equipment
) -> Dict[str, PropertyValue]:
    """Get specific equipment properties.

        Args:
            equipment ( `pyinventory.common.data_class.Equipment` ): equipment object

        Returns:
            Dict[str, PropertyValue]: dictionary of property name to property value
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
            List[ `pyinventory.common.data_class.Equipment` ]: List of found equipments

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: equipment type with this ID does not exist

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
                external_id=equipment.externalId,
                name=equipment.name,
                equipment_type_name=equipment.equipmentType.name,
            )
        )

    return result


def get_equipments_by_location(
    client: SymphonyClient, location_id: str
) -> List[Equipment]:
    """Get equipments by ID of specific location.

        Args:
            location_id (str): location ID

        Returns:
            List[ `pyinventory.common.data_class.Equipment` ]: List of found equipments

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: location with this ID does not exist

        Example:
            ```
            equipments = client.get_equipments_by_location(location_id="60129542651")
            ```
    """
    location_details = LocationEquipmentsQuery.execute(client, id=location_id).location
    if location_details is None:
        raise EntityNotFoundError(entity=Entity.Location, entity_id=location_id)
    result = []
    for equipment in location_details.equipments:
        result.append(
            Equipment(
                id=equipment.id,
                external_id=equipment.externalId,
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
    """Get the equipment attached in a given `position_name` of a given `parent_equipment`

        Args:
            parent_equipment ( `pyinventory.common.data_class.Equipment` ): could be retrieved from
            - `pyinventory.api.equipment.get_equipment`
            - `pyinventory.api.equipment.get_equipment_in_position`
            - `pyinventory.api.equipment.add_equipment`
            - `pyinventory.api.equipment.add_equipment_to_position`

            position_name (str): position name

        Returns:
            `pyinventory.common.data_class.Equipment` object:
                You can use the ID to access the equipment from the UI:
                https://{}.thesymphony.cloud/inventory/inventory?equipment={}

        Raises:
            AssertionException: if parent equipment has more than one
                position with the given name, or none with this name or
                if the position is not occupied.
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
        client=client, parent_equipment=parent_equipment, position_name=position_name
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
    external_id: Optional[str] = None,
) -> Equipment:
    """Create a new equipment in a given `location`.
        The equipment will be of the given `equipment_type`,
        with the given `name` and with the given `properties`.
        If equipment with this name already exists in this location,
        then existing equipment is returned.

        Args:
            name (str): new equipment name
            equipment_type (str): equipment type name
            location ( `pyinventory.common.data_class.Location` ): location object could be retrieved from
            - `pyinventory.api.location.get_location`
            - `pyinventory.api.location.add_location`

            properties_dict (Mapping[str, PropertyValue]): dictionary of property name to property value
            - str - property name
            - PropertyValue - new value of the same type for this property

            external_id (Optional[str]): equipment external ID

        Returns:
            `pyinventory.common.data_class.Equipment`:
                You can use the ID to access the equipment from the UI:
                https://{}.thesymphony.cloud/inventory/inventory?equipment={}

        Raises:
            AssertionException: location contains more than one equipment with the
                same name or if property value in properties_dict does not match the property type
            FailedOperationException: internal inventory error

        Example:
            ```
            from datetime import date
            equipment = client.add_equipment(
                name="Router X123",
                equipment_type="Router",
                location=location,
                properties_dict={
                    "Date Property": date.today(),
                    "Lat/Lng Property": (-1.23,9.232),
                    "E-mail Property": "user@fb.com",
                    "Number Property": 11,
                    "String Property": "aa",
                    "Float Property": 1.23
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
        externalId=external_id,
    )

    try:
        equipment = AddEquipmentMutation.execute(client, add_equipment_input).__dict__[
            ADD_EQUIPMENT
        ]
        client.reporter.log_successful_operation(
            ADD_EQUIPMENT, add_equipment_input.__dict__
        )
    except OperationException as e:
        raise FailedOperationException(
            client.reporter,
            e.err_msg,
            e.err_id,
            ADD_EQUIPMENT,
            add_equipment_input.__dict__,
        )

    return Equipment(
        id=equipment.id,
        external_id=equipment.externalId,
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
            equipment ( `pyinventory.common.data_class.Equipment` ): equipment object
            new_name (Optional[str]): equipment new name
            new_properties (Optional[Dict[str, PropertyValue]]): dictionary of property name to property value
                str - property name
                PropertyValue - new value of the same type for this property

        Returns:
            `pyinventory.common.data_class.Equipment` object

        Raises:
            FailedOperationException: internal inventory error

        Example:
            ```
            location = client.get_location({("Country", "LS_IND_Prod_Copy")})
            equipment = client.get_equipment(name="indProdCpy1_AIO", location=location)
            edited_equipment = client.edit_equipment(
                equipment=equipment,
                new_name="new_name",
                new_properties={"Z AIO - Number": 123},
            )
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
            EDIT_EQUIPMENT
        ]
        client.reporter.log_successful_operation(
            EDIT_EQUIPMENT, edit_equipment_input.__dict__
        )

    except OperationException as e:
        raise FailedOperationException(
            client.reporter,
            e.err_msg,
            e.err_id,
            EDIT_EQUIPMENT,
            edit_equipment_input.__dict__,
        )
    return Equipment(
        id=result.id,
        external_id=result.externalId,
        name=result.name,
        equipment_type_name=result.equipmentType.name,
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
                    external_id=attached_equipment.externalId,
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
    external_id: Optional[str] = None,
) -> Equipment:
    """Create a new equipment inside a given `position_name` of the given `existing_equipment`.
        The equipment will be of the given `equipment_type`, with the given `name` and with the given `properties`.
        If equipment with this name already exists in this position, then existing equipment is returned.

        Args:
            name (str): new equipment name
            equipment_type (str): equipment type name
            existing_equipment ( `pyinventory.common.data_class.Equipment` ): could be retrieved from
            - `pyinventory.api.equipment.get_equipment`
            - `pyinventory.api.equipment.get_equipment_in_position`
            - `pyinventory.api.equipment.add_equipment`
            - `pyinventory.api.equipment.add_equipment_to_position`

            position_name (str): position name in the equipment type.
            properties_dict (Mapping[str, PropertyValue]): dictionary of property name to property value
            - str - property name
            - PropertyValue - new value of the same type for this property

            external_id (Optional[str]): equipment external ID

        Returns:
            `pyinventory.common.data_class.Equipment` object:
                You can use the ID to access the equipment from the UI:
                https://{}.thesymphony.cloud/inventory/inventory?equipment={}

        Raises:
            AssertionException: if parent equipment has more than one position with the given name
                            or if property value in `properties_dict` does not match the property type
            FailedOperationException: for internal inventory error
            `pyinventory.exceptions.EntityNotFoundError`: if `existing_equipment` does not exist

        Example:
            ```
            from datetime import date
            equipment = client.add_equipment_to_position(
                name="Card Y123",
                equipment_type="Card",
                existing_equipment=equipment,
                position_name="Pos 1",
                properties_dict={
                    "Date Property": date.today(),
                    "Lat/Lng Property": (-1.23,9.232),
                    "E-mail Property": "user@fb.com",
                    "Number Property": 11,
                    "String Property": "aa",
                    "Float Property": 1.23
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
        externalId=external_id,
    )

    try:
        equipment = AddEquipmentMutation.execute(client, add_equipment_input).__dict__[
            ADD_EQUIPMENT
        ]
        client.reporter.log_successful_operation(
            ADD_EQUIPMENT_TO_POSITION, add_equipment_input.__dict__
        )
    except OperationException as e:
        raise FailedOperationException(
            client.reporter,
            e.err_msg,
            e.err_id,
            ADD_EQUIPMENT_TO_POSITION,
            add_equipment_input.__dict__,
        )

    return Equipment(
        id=equipment.id,
        external_id=equipment.externalId,
        name=equipment.name,
        equipment_type_name=equipment.equipmentType.name,
    )


def delete_equipment(client: SymphonyClient, equipment: Equipment) -> None:
    """This function delete Equipment.

        Args:
            equipment ( `pyinventory.common.data_class.Equipment` ): equipment object

        Example:
            ```
            client.delete_equipment(equipment=equipment)
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
            Tuple[List[ `pyinventory.common.data_class.Equipment` ], int]

        Example:
            ```
            client.search_for_equipments(limit=10)
            ```
    """
    equipments = EquipmentSearchQuery.execute(
        client, filters=[], limit=limit
    ).equipmentSearch

    total_count = equipments.count
    equipments = [
        Equipment(
            id=equipment.id,
            external_id=equipment.externalId,
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
    equipments, total_count = search_for_equipments(
        client=client, limit=EQUIPMENTS_TO_SEARCH
    )

    for equipment in equipments:
        delete_equipment(client=client, equipment=equipment)

    if total_count == len(equipments):
        return

    with tqdm(total=total_count) as progress_bar:
        progress_bar.update(len(equipments))
        while len(equipments) != 0:
            equipments, _ = search_for_equipments(
                client=client, limit=EQUIPMENTS_TO_SEARCH
            )
            for equipment in equipments:
                delete_equipment(client=client, equipment=equipment)
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
        property_value = _get_property_value(
            property_type=property_type.type.value, property=property
        )
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
    new_external_id: Optional[str] = None,
) -> Equipment:
    """Copy equipment in position.

        Args:
            equipment ( `pyinventory.common.data_class.Equipment` ): equipment object to be copied
            dest_parent_equipment ( `pyinventory.common.data_class.Equipment` ): parent equipment, destination to copy to
            dest_position_name (str): destination position name
            new_external_id (Optional[str]): new external ID for equipment

        Returns:
            `pyinventory.common.data_class.Equipment` object

        Example:
            ```
            location = client.get_location({("Country", "LS_IND_Prod_Copy")})
            equipment_to_copy = client.get_equipment(name="indProdCpy1_AIO", location=location)
            parent_equipment = client.get_equipment(name="parent", location=location)
            copied_equipment = client.copy_equipment_in_position(
                equipment=equipment,
                dest_parent_equipment=parent_equipment,
                dest_position_name="destination position name",
            )
            ```
    """
    equipment_type, properties_dict = _get_equipment_type_and_properties_dict(
        client, equipment
    )
    return add_equipment_to_position(
        client=client,
        name=equipment.name,
        equipment_type=equipment_type,
        existing_equipment=dest_parent_equipment,
        position_name=dest_position_name,
        properties_dict=properties_dict,
        external_id=new_external_id,
    )


def copy_equipment(
    client: SymphonyClient,
    equipment: Equipment,
    dest_location: Location,
    new_external_id: Optional[str] = None,
) -> Equipment:
    """Copy equipment.

        Args:
            equipment ( `pyinventory.common.data_class.Equipment` ): equipment object to be copied
            dest_location ( `pyinventory.common.data_class.Location` ): destination locatoin to copy to
            new_external_id (Optional[str]): equipment external ID

        Returns:
            `pyinventory.common.data_class.Equipment` object

        Example:
            ```
            location = client.get_location({("Country", "LS_IND_Prod_Copy")})
            equipment = client.get_equipment(name="indProdCpy1_AIO", location=location)
            new_location = client.get_location({("Country", "LS_IND_Prod")})
            copied_equipment = client.copy_equipment(
                equipment=equipment,
                dest_location=new_location,
            )
            ```
    """
    equipment_type, properties_dict = _get_equipment_type_and_properties_dict(
        client=client, equipment=equipment
    )
    return add_equipment(
        client=client,
        name=equipment.name,
        equipment_type=equipment_type,
        location=dest_location,
        properties_dict=properties_dict,
        external_id=new_external_id,
    )


def get_equipment_type_of_equipment(
    client: SymphonyClient, equipment: Equipment
) -> EquipmentType:
    """This function returns equipment type object of equipment.

        Args:
            equipment ( `pyinventory.common.data_class.Equipment` ): equipment object

        Returns:
            `pyinventory.common.data_class.EquipmentType` object

        Example:
            ```
            location = client.get_location({("Country", "LS_IND_Prod_Copy")})
            equipment = client.get_equipment(name="indProdCpy1_AIO", location=location)
            equipment_type = client.get_equipment_type_of_equipment(equipment=equipment)
            ```
    """
    equipment_type, _ = _get_equipment_type_and_properties_dict(
        client=client, equipment=equipment
    )
    return client.equipmentTypes[equipment_type]


def get_or_create_equipment(
    client: SymphonyClient,
    name: str,
    equipment_type: str,
    location: Location,
    properties_dict: Mapping[str, PropertyValue],
    external_id: Optional[str] = None,
) -> Equipment:
    """This function checks equipment existence by name in specific location,
        in case it is not found by name, creates one.

        Args:
            name (str): equipment name
            equipment_type (str): equipment type name
            location ( `pyinventory.common.data_class.Location` ): location object could be retrieved from
            - `pyinventory.api.location.get_location`
            - `pyinventory.api.location.add_location`

            properties_dict (Mapping[str, PropertyValue]): dictionary of property name to property value
            - str - property name
            - PropertyValue - new value of the same type for this property

            external_id (Optional[str]): equipment external ID

        Returns:
            `pyinventory.common.data_class.Equipment` object

        Raises:
            AssertionException: location contains more than one equipment with the
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
                    "Date Property": date.today(),
                    "Lat/Lng Property": (-1.23,9.232),
                    "E-mail Property": "user@fb.com",
                    "Number Property": 11,
                    "String Property": "aa",
                    "Float Property": 1.23
                })
            ```
    """
    equipment = _get_equipment_if_exists(client, name, location)
    if equipment is not None:
        return equipment
    return add_equipment(
        client=client,
        name=name,
        equipment_type=equipment_type,
        location=location,
        properties_dict=properties_dict,
        external_id=external_id,
    )


def get_or_create_equipment_in_position(
    client: SymphonyClient,
    name: str,
    equipment_type: str,
    existing_equipment: Equipment,
    position_name: str,
    properties_dict: Mapping[str, PropertyValue],
    external_id: Optional[str] = None,
) -> Equipment:
    """This function checks equipment existence by name in specific location,
        in case it is not found by name, creates one.

        Args:
            name (str): equipment name
            equipment_type (str): equipment type name
            existing_equipment ( `pyinventory.common.data_class.Equipment` ): existing equipment
            position_name (str): position name
            properties_dict (Mapping[str, PropertyValue]): dictionary of property name to property value
            - str - property name
            - PropertyValue - new value of the same type for this property

            external_id (Optional[str]): equipment external ID

        Returns:
            `pyinventory.common.data_class.Equipment` object

        Raises:
            AssertionException: location contains more than one equipment with the
                same name or if property value in properties_dict does not match the property type
            FailedOperationException: internal inventory error

        Example:
            ```
            location = client.get_location({("Country", "LS_IND_Prod_Copy")})
            e_equipment = client.get_equipment(name="indProdCpy1_AIO", location=location)
            equipment_in_position = client.get_or_create_equipment_in_position(
                name="indProdCpy1_AIO",
                equipment_type="router",
                existing_equipment=e_equipment,
                position_name="some_position",
                properties_dict={
                    "Date Property": date.today(),
                    "Lat/Lng Property": (-1.23,9.232),
                    "E-mail Property": "user@fb.com",
                    "Number Property": 11,
                    "String Property": "aa",
                    "Float Property": 1.23
                })
            ```
    """
    equipment = _get_equipment_in_position_if_exists(
        client=client, parent_equipment=existing_equipment, position_name=position_name
    )
    if equipment is not None:
        return equipment

    return add_equipment_to_position(
        client=client,
        name=name,
        equipment_type=equipment_type,
        existing_equipment=existing_equipment,
        position_name=position_name,
        properties_dict=properties_dict,
        external_id=external_id,
    )
