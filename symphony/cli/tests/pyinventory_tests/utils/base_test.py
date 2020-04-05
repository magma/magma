#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import unittest

from pyinventory import InventoryClient

from ..grpc.rpc_pb2_grpc import TenantServiceStub
from . import truncate_client


class BaseTest(unittest.TestCase):
    def __init__(
        self, testName: str, client: InventoryClient, stub: TenantServiceStub
    ) -> None:
        super().__init__(testName)
        self.client = client
        self.stub = stub

    def setUp(self) -> None:
        truncate_client(self.stub)
        self.client._clear_types()
