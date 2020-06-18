#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import unittest

from pyinventory.common.cache import clear_types
from pysymphony import SymphonyClient

from . import truncate_client
from .grpc.rpc_pb2_grpc import TenantServiceStub


class BaseTest(unittest.TestCase):
    def __init__(
        self, test_name: str, client: SymphonyClient, stub: TenantServiceStub
    ) -> None:
        super().__init__(test_name)
        self.client = client
        self.stub = stub

    def setUp(self) -> None:
        truncate_client(self.stub)
        clear_types()
