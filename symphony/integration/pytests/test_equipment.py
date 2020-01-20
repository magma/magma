#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


from utils.base_test import BaseTest


class TestEquipment(BaseTest):
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
                "Tp-Link T1600G", "Router", [("IP", "string", None, True)], {}, []
            )
        )
        self.location = self.client.add_location(
            [("City", "Lima")],
            {"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
            10,
            20,
        )

    def tearDown(self):
        for equipment_type in self.equipment_types_created:
            self.client.delete_equipment_type_with_equipments(equipment_type)
        for location_type in self.location_types_created:
            self.client.delete_location_type_with_locations(location_type)

    def test_equipment_created(self):

        equipment = self.client.add_equipment(
            "TPLinkRouter", "Tp-Link T1600G", self.location, {"IP": "127.0.0.1"}
        )
        fetched_equipment = self.client.get_equipment("TPLinkRouter", self.location)
        self.assertEqual(equipment, fetched_equipment)

    def test_get_or_create_equipment(self):
        equipment = self.client.get_or_create_equipment(
            "TPLinkRouter", "Tp-Link T1600G", self.location, {"IP": "127.0.0.1"}
        )
        equipment2 = self.client.get_or_create_equipment(
            "TPLinkRouter", "Tp-Link T1600G", self.location, {"IP": "127.0.0.1"}
        )
        self.assertEqual(equipment, equipment2)
