#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


from pyinventory.api.equipment import add_equipment
from pyinventory.api.equipment_type import (
    add_equipment_type,
    delete_equipment_type_with_equipments,
)
from pyinventory.api.link import add_link, get_link_in_port_of_equipment
from pyinventory.api.location import add_location
from pyinventory.api.location_type import (
    add_location_type,
    delete_location_type_with_locations,
)
from pyinventory.api.port_type import (
    add_equipment_port_type,
    delete_equipment_port_type,
)
from pyinventory.consts import PropertyDefinition
from pyinventory.exceptions import PortAlreadyOccupiedException
from pyinventory.graphql.property_kind_enum import PropertyKind

from .utils.base_test import BaseTest


class TestLink(BaseTest):
    def setUp(self) -> None:
        self.port_type1 = add_equipment_port_type(
            self.client,
            name="port type 1",
            properties=[
                PropertyDefinition(
                    property_name="port property",
                    property_kind=PropertyKind.string,
                    default_value="port property value",
                    is_fixed=False,
                )
            ],
            link_properties=[
                PropertyDefinition(
                    property_name="link property",
                    property_kind=PropertyKind.string,
                    default_value="link property value",
                    is_fixed=False,
                )
            ],
        )
        self.location_types_created = []
        self.location_types_created.append(
            add_location_type(
                client=self.client,
                name="City",
                properties=[
                    ("Mayor", "string", None, True),
                    ("Contact", "email", None, True),
                ],
            )
        )
        self.equipment_types_created = []
        self.equipment_types_created.append(
            add_equipment_type(
                client=self.client,
                name="Tp-Link T1600G",
                category="Router",
                properties=[("IP", "string", None, True)],
                ports_dict={"Port 1": "port type 1", "Port 2": "port type 1"},
                position_list=[],
            )
        )

    def tearDown(self) -> None:
        for equipment_type in self.equipment_types_created:
            delete_equipment_type_with_equipments(
                client=self.client, equipment_type=equipment_type
            )
        for location_type in self.location_types_created:
            delete_location_type_with_locations(
                client=self.client, location_type=location_type
            )
        delete_equipment_port_type(
            client=self.client, equipment_port_type_id=self.port_type1.id
        )

    def test_cannot_create_link_if_port_occupied(self) -> None:
        location = add_location(
            client=self.client,
            location_hirerchy=[("City", "Lima")],
            properties_dict={"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
            lat=10,
            long=20,
        )
        router1 = add_equipment(
            client=self.client,
            name="TPLinkRouter1",
            equipment_type="Tp-Link T1600G",
            location=location,
            properties_dict={"IP": "192.688.0.1"},
        )
        router2 = add_equipment(
            client=self.client,
            name="TPLinkRouter2",
            equipment_type="Tp-Link T1600G",
            location=location,
            properties_dict={"IP": "192.688.0.2"},
        )
        router3 = add_equipment(
            client=self.client,
            name="TPLinkRouter3",
            equipment_type="Tp-Link T1600G",
            location=location,
            properties_dict={"IP": "192.688.0.2"},
        )
        link = add_link(
            client=self.client,
            equipment_a=router1,
            port_name_a="Port 1",
            equipment_b=router2,
            port_name_b="Port 1",
        )
        fetched_link1 = get_link_in_port_of_equipment(
            client=self.client, equipment=router1, port_name="Port 1"
        )
        fetched_link2 = get_link_in_port_of_equipment(
            client=self.client, equipment=router2, port_name="Port 1"
        )
        self.assertEqual(link, fetched_link1)
        self.assertEqual(link, fetched_link2)

        self.assertRaises(
            PortAlreadyOccupiedException,
            add_link,
            self.client,
            router2,
            "Port 1",
            router3,
            "Port 1",
        )

        link = add_link(
            client=self.client,
            equipment_a=router2,
            port_name_a="Port 2",
            equipment_b=router3,
            port_name_b="Port 1",
        )
        fetched_link1 = get_link_in_port_of_equipment(
            client=self.client, equipment=router2, port_name="Port 2"
        )
        fetched_link2 = get_link_in_port_of_equipment(
            client=self.client, equipment=router3, port_name="Port 1"
        )
        self.assertEqual(link, fetched_link1)
        self.assertEqual(link, fetched_link2)
