#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import Iterator, Tuple

from pysymphony import SymphonyClient

from ..common.constant import PAGINATION_STEP
from ..common.data_class import Equipment, Link
from ..common.data_enum import Entity
from ..exceptions import (
    EntityNotFoundError,
    LinkNotFoundException,
    PortAlreadyOccupiedException,
)
from ..graphql.input.add_link import AddLinkInput
from ..graphql.input.link_side import LinkSide
from ..graphql.mutation.add_link import AddLinkMutation
from ..graphql.query.equipment_ports import EquipmentPortsQuery
from ..graphql.query.links import LinksQuery
from .port import get_port


def get_all_links_and_port_names_of_equipment(
    client: SymphonyClient, equipment: Equipment
) -> Iterator[Tuple[Link, str]]:
    """Returns all links and port names in equipment.

        :param equipment: Equipment object to be copied, could be retrieved from

            * :meth:`~pyinventory.api.equipment.get_equipment`
            * :meth:`~pyinventory.api.equipment.get_equipment_in_position`
            * :meth:`~pyinventory.api.equipment.add_equipment`
            * :meth:`~pyinventory.api.equipment.add_equipment_to_position`

        :type equipment: :class:`~pyinventory.common.data_class.Equipment`

        :raises:
            * `pyinventory.exceptions.EntityNotFoundError`: Link not found
            * FailedOperationException: Internal inventory error

        :return: Iterator of Tuple[Link, str]

            * Link - :class:`~pyinventory.common.data_class.Link`
            * str - port name

        :rtype: Iterator[ Tuple[ :class:`~pyinventory.common.data_class.Link`, str ] ]

        **Example**

        .. code-block:: python

            location = client.get_location({("Country", "LS_IND_Prod_Copy")})
            equipment = client.get_equipment(name="indProdCpy1_AIO", location=location1)
            client.get_all_links_and_port_names_of_equipment(equipment=equipment)
    """

    equipment_data = EquipmentPortsQuery.execute(client, id=equipment.id)
    if equipment_data is None:
        raise EntityNotFoundError(entity=Entity.Equipment, entity_id=equipment.id)
    ports = equipment_data.ports
    for port in ports:
        port_link = port.link
        if port_link is not None:
            yield (
                Link(
                    id=port_link.id,
                    properties=port_link.properties,
                    service_ids=[s.id for s in port_link.services if port_link.services]
                    if port_link.services is not None
                    else [],
                ),
                port.definition.name,
            )


def add_link(
    client: SymphonyClient,
    equipment_a: Equipment,
    port_name_a: str,
    equipment_b: Equipment,
    port_name_b: str,
) -> Link:
    """Connects a link between two ports of two equipments.

        :param equipment_a: Equipment object to connect, could be retrieved from

            * :meth:`~pyinventory.api.equipment.get_equipment`
            * :meth:`~pyinventory.api.equipment.get_equipment_in_position`
            * :meth:`~pyinventory.api.equipment.add_equipment`
            * :meth:`~pyinventory.api.equipment.add_equipment_to_position`

        :type equipment_a: :class:`~pyinventory.common.data_class.Equipment`
        :param port_name_a: Name of the port in equipment A
        :type port_name_a: str
        :param equipment_b: Equipment object to connect, could be retrieved from

            * :meth:`~pyinventory.api.equipment.get_equipment`
            * :meth:`~pyinventory.api.equipment.get_equipment_in_position`
            * :meth:`~pyinventory.api.equipment.add_equipment`
            * :meth:`~pyinventory.api.equipment.add_equipment_to_position`

        :type equipment_b: :class:`~pyinventory.common.data_class.Equipment`
        :param port_name_b: Name of the port in equipment B
        :type port_name_b: str

        :raises:
            * AssertionError: `port_name` in any of the equipments does not exist,
              or match more than one port, or is already occupied by link
            * FailedOperationException: Internal inventory error

        :return: Link object
        :rtype: :class:`~pyinventory.common.data_class.Link`

        **Example**

        .. code-block:: python

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
    link = AddLinkMutation.execute(client, add_link_input)

    return Link(
        id=link.id,
        properties=link.properties,
        service_ids=[s.id for s in link.services],
    )


def get_links(client: SymphonyClient) -> Iterator[Link]:
    """This function returns all existing links

        :return: Links Iterator
        :rtype: Iterator[ :class:`~pyinventory.common.data_class.Link` ]

        **Example**

        .. code-block:: python

            all_links = client.get_links()
    """
    links = LinksQuery.execute(client, first=PAGINATION_STEP)
    edges = links.edges if links else []
    while links is not None and links.pageInfo.hasNextPage:
        links = LinksQuery.execute(
            client, after=links.pageInfo.endCursor, first=PAGINATION_STEP
        )
        if links is not None:
            edges.extend(links.edges)

    for edge in edges:
        node = edge.node
        if node is not None:
            yield Link(
                id=node.id,
                properties=node.properties,
                service_ids=[s.id for s in node.services],
            )


def get_link_in_port_of_equipment(
    client: SymphonyClient, equipment: Equipment, port_name: str
) -> Link:
    """Returns link in specific port by name in equipment.

        :param equipment: Equipment object that has link, could be retrieved from

            * :meth:`~pyinventory.api.equipment.get_equipment`
            * :meth:`~pyinventory.api.equipment.get_equipment_in_position`
            * :meth:`~pyinventory.api.equipment.add_equipment`
            * :meth:`~pyinventory.api.equipment.add_equipment_to_position`

        :type equipment: :class:`~pyinventory.common.data_class.Equipment`
        :param port_name: Name of the port in equipment
        :type port_name: str

        :raises:
            * LinkNotFoundException: Link not found
            * FailedOperationException: Internal inventory error

        :return: Link object
        :rtype: :class:`~pyinventory.common.data_class.Link`

        **Example**

        .. code-block:: python

            location = client.get_location({("Country", "LS_IND_Prod_Copy")})
            equipment = client.get_equipment(name="indProdCpy1_AIO", location=location)
            client.get_link_in_port_of_equipment(
                equipment=equipment,
                port_name="Port A"
            )
    """
    port = get_port(client, equipment, port_name)
    link = port.link
    if link is not None:
        return Link(
            id=link.id, properties=link.properties, service_ids=link.service_ids
        )
    raise LinkNotFoundException(equipment.name, port_name)
