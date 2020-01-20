#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


import os
import time
import unittest

import requests
from pyinventory import InventoryClient


class BaseTest(unittest.TestCase):

    TEST_USER_EMAIL = "fbuser@fb.com"
    PLATFORM_SERVER_HEALTH_CHECK_URL = os.getenv(
        "PLATFORM_SERVER_HEALTH_CHECK_URL", "http://platform-server/healthz"
    )

    def setUp(self):
        self._waitForPlatform()
        self.client = InventoryClient(
            self.TEST_USER_EMAIL, self.TEST_USER_EMAIL, is_dev_mode=True
        )

    def _waitForPlatform(self):
        deadline = time.monotonic() + 60
        while time.monotonic() < deadline:
            try:
                response = requests.get(
                    self.PLATFORM_SERVER_HEALTH_CHECK_URL, timeout=0.5
                )
                if response.status_code == 200:
                    return
            except Exception:
                time.sleep(0.5)
        raise Exception("Failed to wait for platform")
