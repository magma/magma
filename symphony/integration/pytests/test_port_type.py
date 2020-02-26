#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


from utils.base_test import BaseTest


class TestEquipmentPortType(BaseTest):
    def setUp(self):
        self.port_type1 = self.client.add_equipment_port_type(
            name="port type 1",
            properties=[("port property", "string", "port property value", True)],
            link_properties=[
                ("link port property", "string", "link port property value", True)
            ],
        )

    def tearDown(self):
        self.client.delete_equipment_port_type(self.port_type1.id)

    def test_equipment_port_type_created(self):
        fetched_port_type = self.client.get_equipment_port_type(self.port_type1.id)
        self.assertEqual(self.port_type1.id, fetched_port_type.id)

    def test_equipment_port_type_edited(self):
        edited_port_type = self.client.edit_equipment_port_type(
            port_type=self.port_type1,
            new_name="new port type 1",
            new_properties={"port property": "new port property value"},
            new_link_properties={"link port property": "new link port property value"},
        )
        fetched_port_type = self.client.get_equipment_port_type(self.port_type1.id)
        self.assertEqual(fetched_port_type.name, edited_port_type.name)
        self.assertEqual(len(fetched_port_type.properties), 1)
        self.assertEqual(
            fetched_port_type.properties[0]["stringValue"], "new port property value"
        )
        self.assertEqual(len(fetched_port_type.link_properties), 1)
        self.assertEqual(
            fetched_port_type.link_properties[0]["stringValue"],
            "new link port property value",
        )
