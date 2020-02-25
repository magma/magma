#!/usr/bin/env python3
# pyre-strict

from typing import List, Tuple

from gql.gql.client import OperationException

from ..consts import Equipment, EquipmentPort, Link
from ..exceptions import (
    EquipmentPortIsNotUniqueException,
    EquipmentPortNotFoundException,
    LinkNotFoundException,
    PortAlreadyOccupiedException,
)
from ..graphql.add_link_input import AddLinkInput
from ..graphql.add_link_mutation import AddLinkMutation
from ..graphql.equipment_ports_query import EquipmentPortsQuery
from ..graphql.link_side_input import LinkSide
from ..graphql_client import GraphqlClient
from ..reporter import FailedOperationException


ADD_LINK_MUTATION_NAME = "addLink"


def get_all_links_and_port_names_of_equipment(
    client: GraphqlClient, equipment: Equipment
) -> List[Tuple[Link, str]]:
    ports = EquipmentPortsQuery.execute(client, id=equipment.id).equipment.ports
    return [
        (Link(port.link.id), port.definition.name)
        for port in ports
        if port.link is not None
    ]


def _find_port_info(
    client: GraphqlClient, equipment: Equipment, port_name: str
) -> EquipmentPortsQuery.EquipmentPortsQueryData.Node.EquipmentPort:
    ports = EquipmentPortsQuery.execute(client, id=equipment.id).equipment.ports

    ports = [port for port in ports if port.definition.name == port_name]
    if len(ports) > 1:
        raise EquipmentPortIsNotUniqueException(equipment.name, port_name)
    if len(ports) == 0:
        raise EquipmentPortNotFoundException(equipment.name, port_name)
    return ports[0]


def _find_port_definition_id(
    client: GraphqlClient, equipment: Equipment, port_name: str
) -> str:
    port = _find_port_info(client, equipment, port_name)
    if port.link is not None:
        raise PortAlreadyOccupiedException(equipment.name, port_name)

    return port.definition.id


def get_port(
    client: GraphqlClient, equipment: Equipment, port_name: str
) -> EquipmentPort:
    port = _find_port_info(client, equipment, port_name)
    return EquipmentPort(id=port.id)


def add_link(
    client: GraphqlClient,
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

    port_id_a = _find_port_definition_id(client, equipment_a, port_name_a)
    port_id_b = _find_port_definition_id(client, equipment_b, port_name_b)

    add_link_input = AddLinkInput(
        sides=[
            LinkSide(equipment=equipment_a.id, port=port_id_a),
            LinkSide(equipment=equipment_b.id, port=port_id_b),
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
    client: GraphqlClient, equipment: Equipment, port_name: str
) -> Link:
    port = _find_port_info(client, equipment, port_name)
    link = port.link
    if link is not None:
        return Link(id=link.id)
    raise LinkNotFoundException(equipment.name, port_name)
