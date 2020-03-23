#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import unittest
from typing import cast

from pyinventory import InventoryClient

from . import init_client, wait_for_platform
from .constant import TEST_USER_EMAIL


class BaseTest(unittest.TestCase):
    client: InventoryClient = cast(InventoryClient, None)

    @classmethod
    def setUpClass(cls) -> None:
        wait_for_platform()
        cls.client = init_client(TEST_USER_EMAIL, TEST_USER_EMAIL)

    @classmethod
    def tearDownClass(cls) -> None:
        cls.client.session.close()
