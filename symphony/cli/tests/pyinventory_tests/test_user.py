#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


import random
import string

from gql.gql.transport.session import UserDeactivatedException
from pyinventory import InventoryClient
from pyinventory.api.user import (
    activate_user,
    add_user,
    deactivate_user,
    edit_user,
    get_active_users,
)
from pyinventory.graphql.enum.user_role import UserRole
from pyinventory.graphql.enum.user_status import UserStatus

from .grpc.rpc_pb2_grpc import TenantServiceStub
from .utils import init_client
from .utils.base_test import BaseTest


class TestUser(BaseTest):
    def __init__(
        self, testName: str, client: InventoryClient, stub: TenantServiceStub
    ) -> None:
        super().__init__(testName, client, stub)

    @staticmethod
    def random_string(stringLength: int = 10) -> str:
        letters = string.ascii_lowercase
        return "".join(random.choice(letters) for i in range(stringLength))

    def test_user_created(self) -> None:
        user_name = f"{self.random_string()}@fb.com"
        u = add_user(client=self.client, email=user_name, password=user_name)
        self.assertEqual(user_name, u.email)
        self.assertEqual(UserStatus.ACTIVE, u.status)
        active_users = get_active_users(client=self.client)
        self.assertEqual(2, len(active_users))
        client2 = init_client(email=user_name, password=user_name)
        active_users = get_active_users(client=client2)
        self.assertEqual(2, len(active_users))

    def test_user_edited(self) -> None:
        user_name = f"{self.random_string()}@fb.com"
        new_password = self.random_string()
        u = add_user(client=self.client, email=user_name, password=user_name)
        edit_user(
            client=self.client,
            user=u,
            new_password=new_password,
            new_role=UserRole.OWNER,
        )
        client2 = init_client(email=user_name, password=new_password)
        active_users = get_active_users(client=client2)
        self.assertEqual(2, len(active_users))

    def test_user_deactivated(self) -> None:
        user_name = f"{self.random_string()}@fb.com"
        u = add_user(client=self.client, email=user_name, password=user_name)
        deactivate_user(client=self.client, user=u)
        active_users = get_active_users(client=self.client)
        self.assertEqual(1, len(active_users))
        with self.assertRaises(UserDeactivatedException):
            init_client(email=user_name, password=user_name)

    def test_user_reactivated(self) -> None:
        user_name = f"{self.random_string()}@fb.com"
        u = add_user(client=self.client, email=user_name, password=user_name)
        deactivate_user(client=self.client, user=u)
        activate_user(client=self.client, user=u)
        active_users = get_active_users(client=self.client)
        self.assertEqual(2, len(active_users))
