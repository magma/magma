#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


from pyinventory.api.equipment import (
    add_equipment,
    get_equipment,
    get_equipment_properties,
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
from pyinventory.api.port import edit_port_properties, get_port
from pyinventory.api.port_type import (
    add_equipment_port_type,
    delete_equipment_port_type,
)
from pyinventory.consts import PropertyDefinition
from pyinventory.graphql.property_kind_enum import PropertyKind

from .utils.base_test import BaseTest


class TestEquipment(BaseTest):
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
                ports_dict={"tp_link_port": "port type 1"},
                position_list=[],
            )
        )
        self.location = add_location(
            client=self.client,
            location_hirerchy=[("City", "Lima")],
            properties_dict={"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
            lat=10,
            long=20,
        )
        self.equipment = add_equipment(
            client=self.client,
            name="TPLinkRouter",
            equipment_type="Tp-Link T1600G",
            location=self.location,
            properties_dict={"IP": "127.0.0.1"},
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

    def test_equipment_created(self) -> None:

        fetched_equipment = get_equipment(
            client=self.client, name="TPLinkRouter", location=self.location
        )
        self.assertEqual(self.equipment, fetched_equipment)

    def test_get_or_create_equipment(self) -> None:
        equipment2 = get_or_create_equipment(
            client=self.client,
            name="TPLinkRouter",
            equipment_type="Tp-Link T1600G",
            location=self.location,
            properties_dict={"IP": "127.0.0.1"},
        )
        self.assertEqual(self.equipment, equipment2)

    def test_equipment_properties(self) -> None:
        properties = get_equipment_properties(
            client=self.client, equipment=self.equipment
        )
        self.assertTrue("IP" in properties)
        self.assertEquals("127.0.0.1", properties["IP"])

    def test_equipment_get_port(self) -> None:
        fetched_port = get_port(
            client=self.client, equipment=self.equipment, port_name="tp_link_port"
        )
        self.assertEqual(self.port_type1.name, fetched_port.definition.port_type_name)

    def test_equipment_edit_port_properties(self) -> None:
        edit_port_properties(
            client=self.client,
            equipment=self.equipment,
            port_name="tp_link_port",
            new_properties={"port property": "test_port_property"},
        )
        fetched_port = get_port(
            client=self.client, equipment=self.equipment, port_name="tp_link_port"
        )
        port_properties = fetched_port.properties
        self.assertEqual(len(port_properties), 1)

        property_type = port_properties[0].propertyType
        self.assertEqual(property_type.name, "port property")
        self.assertEqual(port_properties[0].stringValue, "test_port_property")
