#!/usr/bin/env python3
# pyre-strict
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
from pyinventory.exceptions import PortAlreadyOccupiedException

from .utils.base_test import BaseTest


class TestLink(BaseTest):
    def setUp(self) -> None:
        super().setUp()
        self.location_types_created = []
        self.location_types_created.append(
            add_location_type(
                self.client,
                "City",
                [("Mayor", "string", None, True), ("Contact", "email", None, True)],
            )
        )
        self.equipment_types_created = []
        self.equipment_types_created.append(
            add_equipment_type(
                self.client,
                "Tp-Link T1600G",
                "Router",
                [("IP", "string", None, True)],
                {"Port 1": "Eth", "Port 2": "Eth"},
                [],
            )
        )

    def tearDown(self) -> None:
        for equipment_type in self.equipment_types_created:
            delete_equipment_type_with_equipments(self.client, equipment_type)
        for location_type in self.location_types_created:
            delete_location_type_with_locations(self.client, location_type)

    def test_cannot_create_link_if_port_occupied(self) -> None:
        location = add_location(
            self.client,
            [("City", "Lima")],
            {"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
            10,
            20,
        )
        router1 = add_equipment(
            self.client,
            "TPLinkRouter1",
            "Tp-Link T1600G",
            location,
            {"IP": "192.688.0.1"},
        )
        router2 = add_equipment(
            self.client,
            "TPLinkRouter2",
            "Tp-Link T1600G",
            location,
            {"IP": "192.688.0.2"},
        )
        router3 = add_equipment(
            self.client,
            "TPLinkRouter3",
            "Tp-Link T1600G",
            location,
            {"IP": "192.688.0.2"},
        )
        link = add_link(self.client, router1, "Port 1", router2, "Port 1")
        fetched_link1 = get_link_in_port_of_equipment(self.client, router1, "Port 1")
        fetched_link2 = get_link_in_port_of_equipment(self.client, router2, "Port 1")
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

        link = add_link(self.client, router2, "Port 2", router3, "Port 1")
        fetched_link1 = get_link_in_port_of_equipment(self.client, router2, "Port 2")
        fetched_link2 = get_link_in_port_of_equipment(self.client, router3, "Port 1")
        self.assertEqual(link, fetched_link1)
        self.assertEqual(link, fetched_link2)
