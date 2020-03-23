#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


from pyinventory.api.equipment_type import (
    add_equipment_type,
    delete_equipment_type_with_equipments,
    get_or_create_equipment_type,
)
from pyinventory.api.location import add_location
from pyinventory.api.location_type import (
    add_location_type,
    delete_location_type_with_locations,
)

from .utils.base_test import BaseTest


class TestEquipmentType(BaseTest):
    def setUp(self) -> None:
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
        self.equipment_type = add_equipment_type(
            client=self.client,
            name="Tp-Link T1600G",
            category="Router",
            properties=[("IP", "string", None, True)],
            ports_dict={},
            position_list=[],
        )
        self.equipment_types_created = []
        self.equipment_types_created.append(self.equipment_type)
        self.location = add_location(
            client=self.client,
            location_hirerchy=[("City", "Lima")],
            properties_dict={"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
            lat=10,
            long=20,
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

    def test_equipment_type_created(self) -> None:
        fetched_equipment_type = get_or_create_equipment_type(
            client=self.client,
            name="Tp-Link T1600G",
            category="Router",
            properties=[("IP", "string", None, True)],
            ports_dict={},
            position_list=[],
        )
        self.assertEqual(self.equipment_type, fetched_equipment_type)
