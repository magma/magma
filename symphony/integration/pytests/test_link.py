#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


from pyinventory.exceptions import PortAlreadyOccupiedException
from utils.base_test import BaseTest


class TestLink(BaseTest):
    def setUp(self):
        super().setUp()
        self.location_types_created = []
        self.location_types_created.append(
            self.client.add_location_type(
                "City",
                [("Mayor", "string", None, True), ("Contact", "email", None, True)],
            )
        )
        self.equipment_types_created = []
        self.equipment_types_created.append(
            self.client.add_equipment_type(
                "Tp-Link T1600G",
                "Router",
                [("IP", "string", None, True)],
                {"Port 1": "Eth", "Port 2": "Eth"},
                [],
            )
        )

    def tearDown(self):
        for equipment_type in self.equipment_types_created:
            self.client.delete_equipment_type_with_equipments(equipment_type)
        for location_type in self.location_types_created:
            self.client.delete_location_type_with_locations(location_type)

    def test_cannot_create_link_if_port_occupied(self):
        location = self.client.add_location(
            [("City", "Lima")],
            {"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
            10,
            20,
        )
        router1 = self.client.add_equipment(
            "TPLinkRouter1", "Tp-Link T1600G", location, {"IP": "192.688.0.1"}
        )
        router2 = self.client.add_equipment(
            "TPLinkRouter2", "Tp-Link T1600G", location, {"IP": "192.688.0.2"}
        )
        router3 = self.client.add_equipment(
            "TPLinkRouter3", "Tp-Link T1600G", location, {"IP": "192.688.0.2"}
        )
        link = self.client.add_link(router1, "Port 1", router2, "Port 1")
        fetched_link1 = self.client.get_link_in_port_of_equipment(router1, "Port 1")
        fetched_link2 = self.client.get_link_in_port_of_equipment(router2, "Port 1")
        self.assertEqual(link, fetched_link1)
        self.assertEqual(link, fetched_link2)

        self.assertRaises(
            PortAlreadyOccupiedException,
            self.client.add_link,
            router2,
            "Port 1",
            router3,
            "Port 1",
        )

        link = self.client.add_link(router2, "Port 2", router3, "Port 1")
        fetched_link1 = self.client.get_link_in_port_of_equipment(router2, "Port 2")
        fetched_link2 = self.client.get_link_in_port_of_equipment(router3, "Port 1")
        self.assertEqual(link, fetched_link1)
        self.assertEqual(link, fetched_link2)
