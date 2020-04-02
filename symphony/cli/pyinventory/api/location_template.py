#!/usr/bin/env python3

from typing import Dict, List, Tuple

from ..client import SymphonyClient
from ..consts import Entity, Equipment, Location
from ..exceptions import EntityNotFoundError
from ..graphql.equipment_positions_query import EquipmentPositionsQuery
from ..graphql.location_equipments_query import LocationEquipmentsQuery
from .equipment import copy_equipment, copy_equipment_in_position
from .link import add_link, get_all_links_and_port_names_of_equipment


def _get_one_level_attachments_of_equipment(
    client: SymphonyClient, equipment: Equipment
) -> List[Tuple[str, Equipment]]:
    equipment_with_positions = EquipmentPositionsQuery.execute(
        client, id=equipment.id
    ).equipment
    if not equipment_with_positions:
        raise EntityNotFoundError(entity=Entity.Equipment, entity_id=equipment.id)
    attachments = []
    for position in equipment_with_positions.positions:
        attached_equipment = position.attachedEquipment
        if attached_equipment is not None:
            attachments.append(
                (
                    position.definition.id,
                    Equipment(
                        id=attached_equipment.id,
                        external_id=attached_equipment.externalId,
                        name=attached_equipment.name,
                        equipment_type_name=attached_equipment.equipmentType.name,
                    ),
                )
            )

    return attachments


def copy_equipment_with_all_attachments(
    client: SymphonyClient, equipment: Equipment, dest_location: Location
) -> Dict[Equipment, Equipment]:
    """Copy the equipment to the new location with all its attachments

        Args:
            equipment (pyinventory.consts.Equipment object): could be retrieved from the following apis:
                * getEquipment
                * getEquipmentInPosition
                * addEquipment
                * addEquipmentToPosition
            dest_location (pyinventory.consts.Location object): retrieved from getLocation or addLocation api.

        Raises: FailedOperationException for internal inventory error

        Returns: dict of source equipment (pyinventory.consts.Equipment) to new equipment (pyinventory.consts.Equipment)
                The dict includes the equipment given as parameter and also all the equipments
                attached to it
    """

    result = {}

    new_equipment = copy_equipment(client, equipment, dest_location)
    equipments = [(equipment, new_equipment)]

    while len(equipments) != 0:
        old_equipment, new_equipment = equipments.pop()
        result[old_equipment] = new_equipment
        attachments = _get_one_level_attachments_of_equipment(client, old_equipment)
        for position_name, child_equipment in attachments:
            new_child_equipment = copy_equipment_in_position(
                client, child_equipment, new_equipment, position_name
            )
            equipments.append((child_equipment, new_child_equipment))
    return result


def apply_location_template_to_location(
    client: SymphonyClient, template_location: Location, location: Location
) -> None:

    location_with_equipments = LocationEquipmentsQuery.execute(
        client, id=template_location.id
    ).location
    if not location_with_equipments:
        raise EntityNotFoundError(
            entity=Entity.Location, entity_id=template_location.id
        )
    equipments = [
        Equipment(
            id=equipment.id,
            external_id=equipment.externalId,
            name=equipment.name,
            equipment_type_name=equipment.equipmentType.name,
        )
        for equipment in location_with_equipments.equipments
    ]
    equipments_to_new_equipments = {}
    for equipment in equipments:
        # return back all and gather link ids
        equipments_to_new_equipments.update(
            copy_equipment_with_all_attachments(client, equipment, location)
        )

    equipments = equipments_to_new_equipments.keys()

    link_to_equipment_and_port = {}
    connected_links = []

    for equipment in equipments:
        links_and_ports = get_all_links_and_port_names_of_equipment(client, equipment)
        for link, port_name in links_and_ports:
            if link not in link_to_equipment_and_port:
                link_to_equipment_and_port[link] = (port_name, equipment)
            else:
                other_port_name, other_equipment = link_to_equipment_and_port.pop(link)
                connected_links.append(
                    (equipment, port_name, other_equipment, other_port_name)
                )

    assert (
        len(link_to_equipment_and_port) == 0
    ), "Some equipments in location are connected to equipments outside the location"

    for equipment, port_name, other_equipment, other_port_name in connected_links:
        new_equipment = equipments_to_new_equipments[equipment]
        new_other_equipment = equipments_to_new_equipments[other_equipment]
        add_link(client, new_equipment, port_name, new_other_equipment, other_port_name)
