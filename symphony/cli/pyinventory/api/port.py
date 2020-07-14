#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import Dict, Iterator

from pysymphony import SymphonyClient

from .._utils import get_graphql_property_inputs
from ..common.cache import PORT_TYPES
from ..common.constant import PAGINATION_STEP
from ..common.data_class import (
    Equipment,
    EquipmentPort,
    EquipmentPortDefinition,
    Link,
    PropertyValue,
)
from ..common.data_enum import Entity
from ..exceptions import EntityNotFoundError, EquipmentPortIsNotUniqueException
from ..graphql.input.edit_equipment_port import EditEquipmentPortInput
from ..graphql.input.link_side import LinkSide
from ..graphql.mutation.edit_equipment_port import EditEquipmentPortMutation
from ..graphql.mutation.edit_link import EditLinkInput, EditLinkMutation
from ..graphql.query.equipment_ports import EquipmentPortsQuery
from ..graphql.query.ports import PortsQuery


def get_port(
    client: SymphonyClient, equipment: Equipment, port_name: str
) -> EquipmentPort:
    """This function returns port in equipment based on its name.

        :param equipment: Existing equipment object
        :type equipment: :class:`~pyinventory.common.data_class.Equipment`
        :param port_name: Existing port name
        :type port_name: str

        :raises:
            * EquipmentPortIsNotUniqueException: There is more than one port with this name
            * :class:`~pyinventory.exceptions.EntityNotFoundError`: Equipment does not exist or port was not found

        :return: EquipmentPort object
        :rtype: :class:`~pyinventory.common.data_class.EquipmentPort`

        **Example**

        .. code-block:: python

            location = client.get_location(location_hirerchy=[("Country", "LS_IND_Prod_Copy")])
            equipment = client.get_equipment(name="indProdCpy1_AIO", location=location)
            port = client.get_port(equipment=equipment, port_name="Z AIO - Port 1")
    """
    equipment_with_ports = EquipmentPortsQuery.execute(client, id=equipment.id)

    if not equipment_with_ports:
        raise EntityNotFoundError(entity=Entity.Equipment, entity_id=equipment.id)

    ports = [
        port for port in equipment_with_ports.ports if port.definition.name == port_name
    ]

    if len(ports) > 1:
        raise EquipmentPortIsNotUniqueException(equipment.name, port_name)
    if len(ports) == 0:
        raise EntityNotFoundError(entity=Entity.EquipmentPort, entity_name=port_name)

    port_type_name = None
    port_type = ports[0].definition.portType
    if port_type is not None:
        port_type_name = port_type.name
    link = ports[0].link

    return EquipmentPort(
        id=ports[0].id,
        properties=ports[0].properties,
        definition=EquipmentPortDefinition(
            id=ports[0].definition.id,
            name=ports[0].definition.name,
            port_type_name=port_type_name,
        ),
        link=Link(
            id=link.id,
            properties=link.properties,
            service_ids=[s.id for s in link.services],
        )
        if link
        else None,
    )


def get_ports(client: SymphonyClient) -> Iterator[EquipmentPort]:
    """This function returns all existing ports

        :return: EquipmentPorts Iterator
        :rtype: Iterator[ :class:`~pyinventory.common.data_class.EquipmentPort` ]

        **Example**

        .. code-block:: python

            all_ports = client.get_ports()
    """

    def generate_pages(
        client: SymphonyClient,
    ) -> Iterator[PortsQuery.PortsQueryData.EquipmentPortConnection]:
        ports = PortsQuery.execute(client, first=PAGINATION_STEP)
        if ports:
            yield ports
        while ports is not None and ports.pageInfo.hasNextPage:
            ports = PortsQuery.execute(
                client, after=ports.pageInfo.endCursor, first=PAGINATION_STEP
            )
            if ports is not None:
                yield ports

    for page in generate_pages(client):
        for edge in page.edges:
            node = edge.node
            if node is not None:
                port_type = None
                if node.definition.portType is not None:
                    port_type = node.definition.portType
                link = None
                if node.link is not None:
                    link = node.link
                yield EquipmentPort(
                    id=node.id,
                    properties=node.properties,
                    definition=EquipmentPortDefinition(
                        id=node.definition.id,
                        name=node.definition.name,
                        port_type_name=port_type.name if port_type else None,
                    ),
                    link=Link(
                        id=link.id,
                        properties=link.properties,
                        service_ids=[s.id for s in link.services],
                    )
                    if link
                    else None,
                )


def edit_port_properties(
    client: SymphonyClient,
    equipment: Equipment,
    port_name: str,
    new_properties: Dict[str, PropertyValue],
) -> EquipmentPort:
    """This function returns edited port in equipment based on its name.

        :param equipment: Existing equipment object
        :type equipment: :class:`~pyinventory.common.data_class.Equipment`
        :param port_name: Equipment type name
        :type port_name: str
        :param new_properties: Dictionary of property name to property value

            * str - property name
            * PropertyValue - new value of the same type for this property

        :type new_properties: Dict[str, PropertyValue]

        :raises:
            * :class:`~pyinventory.exceptions.EntityNotFoundError`: :attr:`~pyinventory.common.data_class.EquipmentPortDefinition.port_type_name` is None,
              there are no properties or there any unknown property name in `new_properties` keys
            * FailedOperationException: internal inventory error

        :return: EquipmentPort object
        :rtype: :class:`~pyinventory.common.data_class.EquipmentPort`

        **Example**

        .. code-block:: python

            location = client.get_location(location_hirerchy=[("Country", "LS_IND_Prod_Copy")])
            equipment = client.get_equipment(name="indProdCpy1_AIO", location=location)
            edited_port = client.edit_port_properties(
                equipment=equipment,
                port_name="Z AIO - Port 1",
                new_properties={"Port Property 2": "test_it"},
            )
    """
    port = get_port(client, equipment, port_name)

    new_property_inputs = []
    if new_properties:
        port_type_name = port.definition.port_type_name
        if port_type_name is None:
            raise EntityNotFoundError(
                entity=Entity.Property,
                msg=f"Not possible to edit properties in '{port.definition.name}' port with undefined PortType",
            )
        property_types = PORT_TYPES[port_type_name].property_types
        new_property_inputs = get_graphql_property_inputs(
            property_types, new_properties
        )

    result = EditEquipmentPortMutation.execute(
        client,
        EditEquipmentPortInput(
            side=LinkSide(equipment=equipment.id, port=port.definition.id),
            properties=new_property_inputs,
        ),
    )
    port_type = result.definition.portType
    link = result.link
    return EquipmentPort(
        id=result.id,
        properties=result.properties,
        definition=EquipmentPortDefinition(
            id=result.definition.id,
            name=result.definition.name,
            port_type_name=port_type.name if port_type else None,
        ),
        link=Link(
            id=link.id,
            properties=link.properties,
            service_ids=[s.id for s in link.services],
        )
        if link
        else None,
    )


def edit_link_properties(
    client: SymphonyClient,
    equipment: Equipment,
    port_name: str,
    new_link_properties: Dict[str, PropertyValue],
) -> EquipmentPort:
    """This function returns edited port in equipment based on its name.

        :param equipment: List of property definitions
        :type equipment: :class:`~pyinventory.common.data_class.Equipment`
        :param port_name: Equipment type name
        :type port_name: str
        :param new_link_properties: Dictionary of property name to property value

            * str - property name
            * PropertyValue - new value of the same type for this property

        :type new_link_properties: Dict[str, PropertyValue], optional

        :raises:
            * :class:`~pyinventory.exceptions.EntityNotFoundError`: :attr:`~pyinventory.common.data_class.EquipmentPortDefinition.port_type_name` is None,
              there are no properties or there any unknown property name in `new_link_properties` keys
            * FailedOperationException: internal inventory error

        :return: EquipmentPort object
        :rtype: :class:`~pyinventory.common.data_class.EquipmentPort`

        **Example**

        .. code-block:: python

            location = client.get_location(location_hirerchy=[("Country", "LS_IND_Prod_Copy")])
            equipment = client.get_equipment(name="indProdCpy1_AIO", location=location)
            edited_port = client.edit_link_properties(
                equipment=equipment,
                port_name="Z AIO - Port 1",
                new_link_properties={"Link Property 1": 98765},
            )
    """
    port = get_port(client, equipment, port_name)
    link = port.link
    if link is None:
        raise EntityNotFoundError(entity=Entity.Link, entity_name=port_name)

    definition_port_type_name = ""
    if port.definition.port_type_name is None:
        raise EntityNotFoundError(
            entity=Entity.Property,
            msg=f"Not possible to edit link properties in '{port.definition.name}' port with undefined PortType",
        )
    else:
        definition_port_type_name = port.definition.port_type_name
    new_link_property_inputs = []
    if new_link_properties and definition_port_type_name:
        link_property_types = PORT_TYPES[definition_port_type_name].link_property_types
        new_link_property_inputs = get_graphql_property_inputs(
            link_property_types, new_link_properties
        )

    result = EditLinkMutation.execute(
        client,
        EditLinkInput(
            id=link.id, properties=new_link_property_inputs, serviceIds=link.service_ids
        ),
    )
    return EquipmentPort(
        id=port.id,
        properties=port.properties,
        definition=EquipmentPortDefinition(
            id=port.definition.id,
            name=port.definition.name,
            port_type_name=port.definition.port_type_name,
        ),
        link=Link(
            id=result.id,
            properties=result.properties,
            service_ids=[s.id for s in result.services],
        )
        if result
        else None,
    )
