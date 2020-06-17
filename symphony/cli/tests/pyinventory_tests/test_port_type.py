#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


from pyinventory import InventoryClient
from pyinventory.api.equipment_type import _populate_equipment_port_types
from pyinventory.api.port_type import (
    add_equipment_port_type,
    edit_equipment_port_type,
    get_equipment_port_type,
)
from pyinventory.consts import PropertyDefinition
from pyinventory.graphql.property_kind_enum import PropertyKind

from .grpc.rpc_pb2_grpc import TenantServiceStub
from .utils.base_test import BaseTest


class TestEquipmentPortType(BaseTest):
    def __init__(
        self, testName: str, client: InventoryClient, stub: TenantServiceStub
    ) -> None:
        super().__init__(testName, client, stub)

    def setUp(self) -> None:
        super().setUp()
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

    def test_equipment_port_type_populated(self) -> None:
        self.assertEqual(len(self.client.portTypes), 1)
        self.client.portTypes = {}
        _populate_equipment_port_types(client=self.client)
        self.assertEqual(len(self.client.portTypes), 1)

    def test_equipment_port_type_created(self) -> None:
        fetched_port_type = get_equipment_port_type(
            client=self.client, equipment_port_type_id=self.port_type1.id
        )
        self.assertEqual(self.port_type1.id, fetched_port_type.id)

    def test_equipment_port_type_edited(self) -> None:
        edited_port_type = edit_equipment_port_type(
            client=self.client,
            port_type=self.port_type1,
            new_name="new port type 1",
            new_properties={"port property": "new port property value"},
            new_link_properties={"link property": "new link property value"},
        )
        fetched_port_type = get_equipment_port_type(
            client=self.client, equipment_port_type_id=self.port_type1.id
        )
        self.assertEqual(fetched_port_type.name, edited_port_type.name)
        self.assertEqual(len(fetched_port_type.property_types), 1)
        self.assertEqual(
            fetched_port_type.property_types[0].stringValue, "new port property value"
        )
        self.assertEqual(len(fetched_port_type.link_property_types), 1)
        self.assertEqual(
            fetched_port_type.link_property_types[0].stringValue,
            "new link property value",
        )
