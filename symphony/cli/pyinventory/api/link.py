#!/usr/bin/env python3
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
    """Returns all links and port names in equipment.

        Args:
            equipment ( `pyinventory.consts.Equipment` ): could be retrieved from
            - `pyinventory.api.equipment.get_equipment`
            - `pyinventory.api.equipment.get_equipment_in_position`
            - `pyinventory.api.equipment.add_equipment`
            - `pyinventory.api.equipment.add_equipment_to_position`

        Returns:
            List[Tuple[ `pyinventory.consts.Link` , str]]:

            - `pyinventory.consts.Link` - link object
            - str - port definition name

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: if link not found
            FailedOperationException: for internal inventory error

        Example:
            ```
            location = client.get_location({("Country", "LS_IND_Prod_Copy")})
            equipment = client.get_equipment(name="indProdCpy1_AIO", location=location1)
            client.get_all_links_and_port_names_of_equipment(equipment=equipment)
            ```
    """

    equipment_data = EquipmentPortsQuery.execute(client, id=equipment.id).equipment
    if equipment_data is None:
        raise EntityNotFoundError(entity=Entity.Equipment, entity_id=equipment.id)
    ports = equipment_data.ports
    result = []
    for port in ports:
        port_link = port.link
        if port_link is not None:
            link = Link(
                id=port_link.id,
                properties=port_link.properties,
                service_ids=[s.id for s in port_link.services if port_link.services]
                if port_link.services is not None
                else [],
            )
            result.append((link, port.definition.name))
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
            equipment_a ( `pyinventory.consts.Equipment` ): could be retrieved from
            - `pyinventory.api.equipment.get_equipment`
            - `pyinventory.api.equipment.get_equipment_in_position`
            - `pyinventory.api.equipment.add_equipment`
            - `pyinventory.api.equipment.add_equipment_to_position`

            port_name_a (str): The name of port in equipment type
            equipment_b ( `pyinventory.consts.Equipment` ): could be retrieved from the following apis:
            - `pyinventory.api.equipment.get_equipment`
            - `pyinventory.api.equipment.get_equipment_in_position`
            - `pyinventory.api.equipment.add_equipment`
            - `pyinventory.api.equipment.add_equipment_to_position`

            port_name_b (str): The name of port in equipment type

        Returns:
            `pyinventory.consts.Link` object

        Raises:
            AssertionError: if port_name in any of the equipment does not exist, or match more than one port
                                    or is already occupied by link
            FailedOperationException: for internal inventory error

        Example:
            ```
            location1 = client.get_location({("Country", "LS_IND_Prod_Copy")})
            equipment1 = client.get_equipment(name="indProdCpy1_AIO", location=location1)
            location2 = client.get_location({("Country", "LS_IND_Prod")})
            equipment2 = client.get_equipment(name="indProd1_AIO", location=location2)
            client.add_link(
                equipment_a=equipment1,
                port_name_a="Port A",
                equipment_b=equipment2,
                port_name_b="Port B"
            )
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

    return Link(
        id=link.id,
        properties=link.properties,
        service_ids=[s.id for s in link.services],
    )


def get_link_in_port_of_equipment(
    client: SymphonyClient, equipment: Equipment, port_name: str
) -> Link:
    """Returns link in specific port by name in equipment.

        Args:
            equipment ( `pyinventory.consts.Equipment` ): could be retrieved from
            - `pyinventory.api.equipment.get_equipment`
            - `pyinventory.api.equipment.get_equipment_in_position`
            - `pyinventory.api.equipment.add_equipment`
            - `pyinventory.api.equipment.add_equipment_to_position`

            port_name (str): The name of port in equipment type

        Returns:
            `pyinventory.consts.Link` object

        Raises:
            LinkNotFoundException: if link not found
            FailedOperationException: for internal inventory error

        Example:
            ```
            location = client.get_location({("Country", "LS_IND_Prod_Copy")})
            equipment = client.get_equipment(name="indProdCpy1_AIO", location=location)
            client.get_link_in_port_of_equipment(
                equipment=equipment,
                port_name="Port A"
            )
            ```
    """
    port = get_port(client, equipment, port_name)
    link = port.link
    if link is not None:
        return Link(
            id=link.id, properties=link.properties, service_ids=link.service_ids
        )
    raise LinkNotFoundException(equipment.name, port_name)
