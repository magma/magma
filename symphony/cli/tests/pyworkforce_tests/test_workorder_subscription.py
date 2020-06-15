#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


import time
from typing import Any, Dict, List

from pysymphony import SymphonyClient
from pyworkforce.api.workorder import add_workorder, add_workorder_type

from ..utils.base_test import BaseTest
from ..utils.constant import TEST_USER_EMAIL, TestMode
from ..utils.grpc.rpc_pb2_grpc import TenantServiceStub
from .subscription_client import SubscriptionClient


class TestWorkOrderSubscription(BaseTest):
    def __init__(
        self, testName: str, client: SymphonyClient, stub: TenantServiceStub
    ) -> None:
        super().__init__(testName, client, stub)

    def setUp(self) -> None:
        super().setUp()
        self.type = add_workorder_type(self.client, "Work Order Template")

    def test_subscribe_to_work_order_added(self) -> None:
        from ..utils import TEST_MODE

        url = "wss://fb-test.localtest.me/graph/query"
        if TEST_MODE == TestMode.DEV:
            url = "wss://fb-test.thesymphony.cloud/graph/query"

        sub_client = SubscriptionClient(url, TEST_USER_EMAIL, TEST_USER_EMAIL)
        workorders_added: List[Dict[str, str]] = []

        def callback(_id: str, data: Dict[str, Any]) -> None:
            workorders_added.append(data["payload"]["data"]["workOrderAdded"])

        query = """
            subscription {
                workOrderAdded {
                    id
                    name
                }
            }
        """

        sub_id = sub_client.subscribe(query, callback=callback)

        workorder = add_workorder(
            self.client, name="My Work Order", workorder_type=self.type
        )
        i = 0
        while len(workorders_added) == 0:
            time.sleep(1)
            i = i + 1
            if i == 3:
                break

        self.assertEqual(1, len(workorders_added))
        self.assertEqual(workorder.id, workorders_added[0]["id"])
        self.assertEqual(workorder.name, workorders_added[0]["name"])

        sub_client.stop_subscribe(sub_id)
        sub_client.close()
