#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


from pyinventory import InventoryClient
from pyinventory._utils import get_property_type_input
from pyinventory.api.equipment_type import (
    add_equipment_type,
    edit_equipment_type_property_type,
    get_or_create_equipment_type,
)
from pyinventory.api.property_type import get_property_type_id, get_property_types
from pyinventory.consts import Entity, PropertyDefinition
from pyinventory.graphql.property_kind_enum import PropertyKind

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

    def test_equipment_type_add_external_id_to_property_type(self) -> None:
        equipment_type_name = self.equipment_type.name
        property_type_name = "IP"
        property_type_id = get_property_type_id(
            client=self.client,
            entity_type=Entity.EquipmentType,
            entity_name=equipment_type_name,
            property_type_name=property_type_name,
        )
        e_type = edit_equipment_type_property_type(
            client=self.client,
            equipment_type_name=equipment_type_name,
            property_type_id=property_type_id,
            new_property_definition=PropertyDefinition(
                property_name=property_type_name,
                property_kind=PropertyKind.string,
                default_value=None,
                is_fixed=False,
                external_id="12345",
            ),
        )
        property_types = get_property_types(
            client=self.client,
            entity_type=Entity.EquipmentType,
            entity_name=e_type.name,
        )
        fetched_property_type = None
        for property_type in property_types:
            if property_type.name == property_type_name:
                fetched_property_type = property_type
        self.assertIsNotNone(fetched_property_type)
        self.assertEqual(fetched_property_type.externalId, "12345")

    def test_equipment_type_property_type_name(self) -> None:
        equipment_type_name = self.equipment_type.name
        property_type_name = "IP"
        new_name = "new_IP"
        property_type_id = get_property_type_id(
            client=self.client,
            entity_type=Entity.EquipmentType,
            entity_name=equipment_type_name,
            property_type_name=property_type_name,
        )
        e_type = edit_equipment_type_property_type(
            client=self.client,
            equipment_type_name=equipment_type_name,
            property_type_id=property_type_id,
            new_property_definition=PropertyDefinition(
                property_name=new_name,
                property_kind=PropertyKind.string,
                default_value=None,
                is_fixed=False,
                external_id=None,
            ),
        )
        property_types = get_property_types(
            client=self.client,
            entity_type=Entity.EquipmentType,
            entity_name=e_type.name,
        )
        fetched_property_type = None
        for property_type in property_types:
            property_type_input = get_property_type_input(property_type)
            if property_type_input.name == property_type_name:
                fetched_property_type = property_type_input
        self.assertEqual(fetched_property_type, None)
        for property_type in property_types:
            property_type_input = get_property_type_input(property_type)
            if property_type_input.name == new_name:
                fetched_property_type = property_type_input
        self.assertIsNotNone(fetched_property_type)
        self.assertEqual(fetched_property_type.name, new_name)
