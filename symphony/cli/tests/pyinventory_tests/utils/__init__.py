#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
import time

import requests
from pyinventory import InventoryClient

from .constant import PLATFORM_SERVER_HEALTH_CHECK_URL


RUN_LOCALLY = False


def wait_for_platform() -> None:
    platform_server_health_check = PLATFORM_SERVER_HEALTH_CHECK_URL
    if RUN_LOCALLY:
        platform_server_health_check = "http://fb-test.localtest.me/healthz"

    deadline = time.monotonic() + 60
    while time.monotonic() < deadline:
        try:
            response = requests.get(platform_server_health_check, timeout=0.5)
            if response.status_code == 200:
                return
            print(
                f"Response failed with status code {response.status_code}"
                f' and with message "{response.text}"'
            )
        except Exception as e:
            print(f"Request failed with exception {e}")
            time.sleep(0.5)
    raise Exception("Failed to wait for platform")


def init_client(email: str, password: str) -> InventoryClient:
    if RUN_LOCALLY:
        return InventoryClient(email, password, is_local_host=True)
    else:
        return InventoryClient(email, password, is_dev_mode=True)
