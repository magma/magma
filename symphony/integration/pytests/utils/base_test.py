#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


import time
import unittest

import requests
from pyinventory import InventoryClient
from utils.constant import PLATFORM_SERVER_HEALTH_CHECK_URL, TEST_USER_EMAIL


class BaseTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls._waitForPlatform()
        cls.client = InventoryClient(TEST_USER_EMAIL, TEST_USER_EMAIL, is_dev_mode=True)

    @classmethod
    def tearDownClass(cls):
        cls.client.session.close()

    @classmethod
    def _waitForPlatform(self):
        deadline = time.monotonic() + 60
        while time.monotonic() < deadline:
            try:
                response = requests.get(PLATFORM_SERVER_HEALTH_CHECK_URL, timeout=0.5)
                if response.status_code == 200:
                    return
            except Exception:
                time.sleep(0.5)
        raise Exception("Failed to wait for platform")
