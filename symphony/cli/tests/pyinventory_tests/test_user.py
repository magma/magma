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
from pyinventory.graphql.user_role_enum import UserRole
from pyinventory.graphql.user_status_enum import UserStatus

from .utils import init_client
from .utils.base_test import BaseTest
from .utils.constant import TEST_USER_EMAIL


class TestUser(BaseTest):
    def tearDown(self) -> None:
        active_users = get_active_users(self.client)
        for user in active_users:
            if user.email != TEST_USER_EMAIL:
                deactivate_user(self.client, user)

    @staticmethod
    def random_string(stringLength: int = 10) -> str:
        letters = string.ascii_lowercase
        return "".join(random.choice(letters) for i in range(stringLength))

    def test_user_created(self) -> None:
        user_name = f"{self.random_string()}@fb.com"
        u = add_user(self.client, user_name, user_name)
        self.assertEqual(user_name, u.email)
        self.assertEqual(UserStatus.ACTIVE, u.status)
        active_users = get_active_users(self.client)
        self.assertEqual(2, len(active_users))
        client2 = init_client(user_name, user_name)
        active_users = get_active_users(client2)
        self.assertEqual(2, len(active_users))

    def test_user_edited(self) -> None:
        user_name = f"{self.random_string()}@fb.com"
        new_password = self.random_string()
        u = add_user(self.client, user_name, user_name)
        edit_user(self.client, u, new_password, UserRole.OWNER)
        client2 = init_client(user_name, new_password)
        active_users = get_active_users(client2)
        self.assertEqual(2, len(active_users))

    def test_user_deactivated(self) -> None:
        user_name = f"{self.random_string()}@fb.com"
        u = add_user(self.client, user_name, user_name)
        deactivate_user(self.client, u)
        active_users = get_active_users(self.client)
        self.assertEqual(1, len(active_users))
        with self.assertRaises(UserDeactivatedException):
            init_client(user_name, user_name)

    def test_user_reactivated(self) -> None:
        user_name = f"{self.random_string()}@fb.com"
        u = add_user(self.client, user_name, user_name)
        deactivate_user(self.client, u)
        activate_user(self.client, u)
        active_users = get_active_users(self.client)
        self.assertEqual(2, len(active_users))
