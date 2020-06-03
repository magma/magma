#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from pyinventory.api.equipment_type import add_equipment_type
from pyinventory.api.service import add_service
from pyinventory.api.service_type import (
    _populate_service_types,
    add_service_type,
    delete_service_type,
    delete_service_type_with_services,
    edit_service_type,
    get_service_type,
)
from pyinventory.common.cache import SERVICE_TYPES
from pyinventory.common.data_class import PropertyDefinition, ServiceEndpointDefinition
from pyinventory.graphql.enum.property_kind import PropertyKind
from pysymphony import SymphonyClient

from ..utils.base_test import BaseTest
from ..utils.grpc.rpc_pb2_grpc import TenantServiceStub


class TestServiceType(BaseTest):
    def __init__(
        self, testName: str, client: SymphonyClient, stub: TenantServiceStub
    ) -> None:
        super().__init__(testName, client, stub)

    def setUp(self) -> None:
        super().setUp()
        self.service_type = add_service_type(
            client=self.client,
            name="Internet Access",
            has_customer=True,
            properties=[
                PropertyDefinition(
                    property_name="Service Package",
                    property_kind=PropertyKind.string,
                    default_value=None,
                    is_fixed=False,
                )
            ],
            endpoint_definitions=[],
        )
        self.service = add_service(
            client=self.client,
            name="Room 201 Internet Access",
            external_id="S3232",
            service_type=self.service_type.name,
            customer=None,
            properties_dict={"Service Package": "Public 5G"},
        )

    def test_service_type_populated(self) -> None:
        self.assertEqual(len(SERVICE_TYPES), 1)
        SERVICE_TYPES.clear()
        _populate_service_types(client=self.client)
        self.assertEqual(len(SERVICE_TYPES), 1)

    def test_service_type_created(self) -> None:
        fetched_service_type = get_service_type(
            client=self.client, service_type_id=self.service_type.id
        )
        self.assertEqual(fetched_service_type, self.service_type)

    def test_service_type_edited(self) -> None:
        equipment_type = add_equipment_type(
            client=self.client,
            name="Tp-Link T1600G",
            category="Router",
            properties=[
                PropertyDefinition(
                    property_name="IP",
                    property_kind=PropertyKind.string,
                    default_value=None,
                    is_fixed=False,
                )
            ],
            ports_dict={},
            position_list=[],
        )
        new_name = "New Service Package"
        new_properties = {"Service Package": "Public 5G"}
        endpoint_definitions = SERVICE_TYPES[
            self.service_type.name
        ].endpoint_definitions
        self.assertFalse(endpoint_definitions)
        edited_service_type = edit_service_type(
            client=self.client,
            service_type=self.service_type,
            new_name=new_name,
            new_properties=new_properties,
            new_endpoints=[
                ServiceEndpointDefinition(
                    id=None,
                    name="EndpointDefinition",
                    role="CPE",
                    endpoint_definition_index=0,
                    equipment_type_id=equipment_type.id,
                )
            ],
        )
        endpoint_definitions = SERVICE_TYPES[
            edited_service_type.name
        ].endpoint_definitions
        self.assertEqual(len(endpoint_definitions), 1)
        self.assertEqual(edited_service_type.name, new_name)
        self.assertEqual(len(edited_service_type.property_types), 1)
        self.assertEqual(edited_service_type.property_types[0].stringValue, "Public 5G")

    def test_service_type_delete(self) -> None:
        delete_service_type(client=self.client, service_type=self.service_type)
        self.assertEqual(len(SERVICE_TYPES), 0)

    def test_service_type_delete_with_services(self) -> None:
        delete_service_type_with_services(
            client=self.client, service_type=self.service_type
        )
        self.assertEqual(len(SERVICE_TYPES), 0)
