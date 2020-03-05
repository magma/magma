#!/usr/bin/env python3
# pyre-strict
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import List, Tuple

from gql.gql.client import OperationException
from gql.gql.reporter import FailedOperationException

from ..client import SymphonyClient
from ..consts import Entity, Equipment, Link
from ..exceptions import (
    EntityNotFoundError,
    LinkNotFoundException,
    PortAlreadyOccupiedException,
)
from ..graphql.add_link_input import AddLinkInput
from ..graphql.add_link_mutation import AddLinkMutation
from ..graphql.equipment_ports_query import EquipmentPortsQuery
from ..graphql.link_side_input import LinkSide
from .port import get_port


ADD_LINK_MUTATION_NAME = "addLink"


def get_all_links_and_port_names_of_equipment(
    client: SymphonyClient, equipment: Equipment
) -> List[Tuple[Link, str]]:
    equipment_with_ports = EquipmentPortsQuery.execute(
        client, id=equipment.id
    ).equipment
    if not equipment_with_ports:
        raise EntityNotFoundError(entity=Entity.Equipment, entity_id=equipment.id)
    result = []
    for port in equipment_with_ports.ports:
        link = port.link
        if link is not None:
            result.append((Link(link.id), port.definition.name))
    return result


def add_link(
    client: SymphonyClient,
    equipment_a: Equipment,
    port_name_a: str,
    equipment_b: Equipment,
    port_name_b: str,
) -> Link:
    """Connects a link between two ports of two equipments.

        Args:
            equipment_a (pyinventory.consts.Equipment object): could be retrieved from the following apis:

                * `pyinventory.api.equipment.get_equipment`

                * `pyinventory.api.equipment.get_equipment_in_position`

                * `pyinventory.api.equipment.add_equipment`

                * `pyinventory.api.equipment.add_equipment_to_position`

            port_name_a (str): The name of port in equipment type
            equipment_b (pyinventory.consts.Equipment object): could be retrieved from the following apis:

                * `pyinventory.api.equipment.get_equipment`

                * `pyinventory.api.equipment.get_equipment_in_position`

                * `pyinventory.api.equipment.add_equipment`

                * `pyinventory.api.equipment.add_equipment_to_position`
            
            port_name_b (str): The name of port in equipment type

        Returns: pyinventory.consts.Link object with id field

        Raises: AssertionError if portName in any of the equipment does not exist, match more than one port
                                    or is already occupied by link
                FailedOperationException for internal inventory error

        Example:
        ```
        client.addLink(VSATEquipment, "Port A", MWEquipment, "Port B")
        ```
    """

    port_a = get_port(client, equipment_a, port_name_a)
    if port_a.link is not None:
        raise PortAlreadyOccupiedException(equipment_a.name, port_a.definition.name)
    port_b = get_port(client, equipment_b, port_name_b)
    if port_b.link is not None:
        raise PortAlreadyOccupiedException(equipment_b.name, port_b.definition.name)

    add_link_input = AddLinkInput(
        sides=[
            LinkSide(equipment=equipment_a.id, port=port_a.definition.id),
            LinkSide(equipment=equipment_b.id, port=port_b.definition.id),
        ],
        properties=[],
        serviceIds=[],
    )
    try:
        link = AddLinkMutation.execute(client, add_link_input).__dict__[
            ADD_LINK_MUTATION_NAME
        ]
        client.reporter.log_successful_operation(
            ADD_LINK_MUTATION_NAME, add_link_input.__dict__
        )
    except OperationException as e:
        raise FailedOperationException(
            client.reporter,
            e.err_msg,
            e.err_id,
            ADD_LINK_MUTATION_NAME,
            add_link_input.__dict__,
        )

    return Link(id=link.id)


def get_link_in_port_of_equipment(
    client: SymphonyClient, equipment: Equipment, port_name: str
) -> Link:
    port = get_port(client, equipment, port_name)
    link = port.link
    if link is not None:
        return Link(id=link.id)
    raise LinkNotFoundException(equipment.name, port_name)
