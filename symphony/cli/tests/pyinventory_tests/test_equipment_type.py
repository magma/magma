#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


from pyinventory import InventoryClient
from pyinventory.api.equipment_type import (
    add_equipment_type,
    get_or_create_equipment_type,
)

from .grpc.rpc_pb2_grpc import TenantServiceStub
from .utils.base_test import BaseTest


class TestEquipmentType(BaseTest):
    def __init__(
        self, testName: str, client: InventoryClient, stub: TenantServiceStub
    ) -> None:
        super().__init__(testName, client, stub)

    def setUp(self) -> None:
        super().setUp()
        self.equipment_type = add_equipment_type(
            client=self.client,
            name="Tp-Link T1600G",
            category="Router",
            properties=[("IP", "string", None, True)],
            ports_dict={},
            position_list=[],
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
