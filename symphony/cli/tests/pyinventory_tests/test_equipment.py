#!/usr/bin/env python3
# pyre-strict
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


from pyinventory.api.equipment import (
    add_equipment,
    get_equipment,
    get_or_create_equipment,
)
from pyinventory.api.equipment_type import (
    add_equipment_type,
    delete_equipment_type_with_equipments,
)
from pyinventory.api.location import add_location
from pyinventory.api.location_type import (
    add_location_type,
    delete_location_type_with_locations,
)

from .utils.base_test import BaseTest


class TestEquipment(BaseTest):
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
                {},
                [],
            )
        )
        self.location = add_location(
            self.client,
            [("City", "Lima")],
            {"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
            10,
            20,
        )

    def tearDown(self) -> None:
        for equipment_type in self.equipment_types_created:
            delete_equipment_type_with_equipments(self.client, equipment_type)
        for location_type in self.location_types_created:
            delete_location_type_with_locations(self.client, location_type)

    def test_equipment_created(self) -> None:

        equipment = add_equipment(
            self.client,
            "TPLinkRouter",
            "Tp-Link T1600G",
            self.location,
            {"IP": "127.0.0.1"},
        )
        fetched_equipment = get_equipment(self.client, "TPLinkRouter", self.location)
        self.assertEqual(equipment, fetched_equipment)

    def test_get_or_create_equipment(self) -> None:
        equipment = get_or_create_equipment(
            self.client,
            "TPLinkRouter",
            "Tp-Link T1600G",
            self.location,
            {"IP": "127.0.0.1"},
        )
        equipment2 = get_or_create_equipment(
            self.client,
            "TPLinkRouter",
            "Tp-Link T1600G",
            self.location,
            {"IP": "127.0.0.1"},
        )
        self.assertEqual(equipment, equipment2)
