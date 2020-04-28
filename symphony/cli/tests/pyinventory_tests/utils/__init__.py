#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
import os
import sys
import time

import requests
from pyinventory import InventoryClient
from pyinventory.common.endpoint import LOCALHOST_INVENTORY_SERVER

from .constant import PLATFORM_SERVER_HEALTH_CHECK_URL, TestMode


if True:
    sys.path.insert(0, os.path.join(os.path.dirname(__file__), "../grpc"))
    from ..grpc.rpc_pb2_grpc import (
        TenantServiceStub,
        google_dot_protobuf_dot_wrappers__pb2,
    )


TEST_MODE: TestMode = TestMode.DEV
TENANT = "fb-test"


def wait_for_platform() -> None:
    if TEST_MODE == TestMode.REMOTE:
        return
    platform_server_health_check = PLATFORM_SERVER_HEALTH_CHECK_URL
    if TEST_MODE == TestMode.LOCAL:
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
    if TEST_MODE == TestMode.LOCAL:
        return InventoryClient(email, password, tenant=TENANT, is_local_host=True)
    elif TEST_MODE == TestMode.REMOTE:
        return InventoryClient(email, password, tenant=f"{TENANT}.staging")
    else:
        return InventoryClient(email, password, is_dev_mode=True)


def get_grpc_server_address() -> str:
    if TEST_MODE == TestMode.LOCAL:
        return LOCALHOST_INVENTORY_SERVER.format(TENANT) + ":8083"
    elif TEST_MODE == TestMode.REMOTE:
        raise NotImplementedError("T64902729")
    else:
        return "graph:443"


def truncate_client(stub: TenantServiceStub) -> None:
    stub.Truncate(google_dot_protobuf_dot_wrappers__pb2.StringValue(value=TENANT))
