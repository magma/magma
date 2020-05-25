#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from pyinventory.api.service_type import (
    _populate_service_types,
    add_service_type,
    get_service_type,
)
from pyinventory.common.cache import SERVICE_TYPES
from pyinventory.common.data_class import PropertyDefinition
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
                    default_value="Public 5G",
                    is_fixed=False,
                ),
                PropertyDefinition(
                    property_name="Address Family",
                    property_kind=PropertyKind.string,
                    default_value=None,
                    is_fixed=False,
                ),
            ],
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
