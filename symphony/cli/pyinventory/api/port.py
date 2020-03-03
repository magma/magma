#!/usr/bin/env python3
# pyre-strict
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from dataclasses import asdict

from ..client import SymphonyClient
from ..consts import Entity, Equipment, EquipmentPort
from ..exceptions import EntityNotFoundError, EquipmentPortIsNotUniqueException
from ..graphql.equipment_ports_query import EquipmentPortsQuery


EDIT_EQUIPMENT_PORT_MUTATION_NAME = "editEquipmentPort"


def get_port(
    client: SymphonyClient, equipment: Equipment, port_name: str
) -> EquipmentPort:
    """This function returns port in equipment based on its name.

        Args:
            equipment (pyinventory.consts.Equipment object): existing equipment object
            port_name (str): existing port name

        Returns:
            pyinventory.consts.EquipmentPort object

        Raises:
            EquipmentPortIsNotUniqueException: there is more than one port with this name
            EntityNotFoundError: the port was not found

        Example:
            ```
            location = client.get_location([("Country", "LS_IND_Prod_Copy")])
            equipment = client.get_equipment("indProdCpy1_AIO", location)
            port = client.get_port(equipment, "Z AIO - Port 1") 
            ```
    """
    equipment_ports = EquipmentPortsQuery.execute(
        client, id=equipment.id
    ).equipment.ports

    ports = [port for port in equipment_ports if port.definition.name == port_name]
    if len(ports) > 1:
        raise EquipmentPortIsNotUniqueException(equipment.name, port_name)
    if len(ports) == 0:
        raise EntityNotFoundError(entity=Entity.EquipmentPort, entity_name=port_name)
    return EquipmentPort(
        id=ports[0].id,
        properties=[asdict(p) for p in ports[0].properties],
        definition=ports[0].definition,
        link=ports[0].link,
    )
