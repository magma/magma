#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import List, Optional

from pysymphony import SymphonyClient

from ..common.constant import SUPERUSER_ROLE, USER_ROLE
from ..common.data_enum import Entity
from ..exceptions import EntityNotFoundError
from ..graphql.enum.user_role import UserRole
from ..graphql.enum.user_status import UserStatus
from ..graphql.fragment.user import UserFragment
from ..graphql.input.edit_user import EditUserInput
from ..graphql.mutation.edit_user import EditUserMutation
from ..graphql.query.user import UserQuery
from ..graphql.query.users import UsersQuery


def get_user(client: SymphonyClient, email: str) -> UserFragment:
    """Returns `pyinventory.graphql.fragment.user.UserFragment` object by its email

        Args:
            email (str): the email address the user registered with

        Returns:
            `pyinventory.graphql.fragment.user.UserFragment` object

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: the user was not found
            FailedOperationException: internal inventory error

        Example:
            ```
            user = client.get_user(email="user@test.com")
            ```
    """
    user = UserQuery.execute(client, email)
    if user is None:
        raise EntityNotFoundError(entity=Entity.User, entity_name=email)
    return UserFragment(
        id=user.id,
        authID=user.authID,
        email=user.email,
        status=user.status,
        role=user.role,
    )


def add_user(client: SymphonyClient, email: str, password: str) -> UserFragment:
    """Adds new user to inventory with its email and password

        Args:
            email (str): the email address of the user
            password (str): password the user would connect with

        Returns:
            `pyinventory.graphql.fragment.user.UserFragment` object

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: the user was not created properly
            FailedOperationException: internal inventory error
            AssertionError: The user was not created for some known reason
            HTTPError: Error with connection

        Example:
            ```
            user = client.add_user(email="user@test.com", password="P0ssW!rd0f43")
            ```
    """
    resp = client.post(
        "/user/async/",
        {"email": email, "password": password, "role": USER_ROLE, "networkIDs": []},
    )

    if not resp.ok:
        error_message = resp.json().get("error", None)
        if error_message is not None:
            raise AssertionError(error_message)
        raise

    return get_user(client=client, email=email)


def edit_user(
    client: SymphonyClient,
    user: UserFragment,
    new_password: Optional[str] = None,
    new_role: Optional[UserRole] = None,
) -> None:
    """Edit user password and role

        Args:
            user ( `pyinventory.graphql.fragment.user.UserFragment` ): user to edit
            new_password (Optional[str]): new password the user would connect with
            new_role ( `pyinventory.graphql.enum.user_role.UserRole` ): user new role

        Raises:
            FailedOperationException: internal inventory error
            AssertionError: The user was not edited for some known reason
            HTTPError: Error with connection

        Example:
            ```
            user = client.add_user(email="user@test.com", password="P0ssW!rd0f43")
            client.edit_user(user=user, new_password="New_Password4Ever", new_role=UserRole.ADMIN)
            ```
    """
    params = {}
    if new_password is not None:
        params.update({"password": new_password})
    if new_role is not None:
        params.update(
            {"role": USER_ROLE if new_role == UserRole.USER else SUPERUSER_ROLE}
        )
    resp = client.put(f"/user/set/{user.email}", params)

    if not resp.ok:
        error_message = resp.json().get("error", None)
        if error_message is not None:
            raise AssertionError(error_message)
        raise

    if new_role is not None:
        EditUserMutation.execute(client, input=EditUserInput(id=user.id, role=new_role))


def deactivate_user(client: SymphonyClient, user: UserFragment) -> None:
    """Deactivate the user which would prevent the user from login in to symphony
       Users in symphony are never deleted. Only de-activated.


        Args:
            user ( `pyinventory.graphql.fragment.user.UserFragment` ): user to deactivate


        Raises:
            FailedOperationException: internal inventory error

        Example:
            ```
            user = client.get_user(email="user@test.com")
            client.deactivate_user(user=user)
            ```
    """
    EditUserMutation.execute(
        client, input=EditUserInput(id=user.id, status=UserStatus.DEACTIVATED)
    )


def activate_user(client: SymphonyClient, user: UserFragment) -> None:
    """Activate the user which would allow the user to login again to symphony

        Args:
            user ( `pyinventory.graphql.fragment.user.UserFragment` ): user to activate

        Raises:
            FailedOperationException: internal inventory error

        Example:
            ```
            user = client.get_user(email="user@test.com")
            client.activate_user(user=user)
            ```
    """
    EditUserMutation.execute(
        client, input=EditUserInput(id=user.id, status=UserStatus.ACTIVE)
    )


def get_users(client: SymphonyClient) -> List[UserFragment]:
    """Get the list of users in the system (both active and deactivate)

        Returns:
            List[ `pyinventory.graphql.fragment.user.UserFragment` ]

        Raises:
            FailedOperationException: internal inventory error

        Example:
            ```
            users = client.get_users()
            for user in users:
                print(user.email)
            ```
    """
    result = UsersQuery.execute(client)
    if result is None:
        return []
    users = []
    for edge in result.edges:
        node = edge.node
        if node is not None:
            users.append(
                UserFragment(
                    id=node.id,
                    authID=node.authID,
                    email=node.email,
                    status=node.status,
                    role=node.role,
                )
            )
    return users


def get_active_users(client: SymphonyClient) -> List[UserFragment]:
    """Get the list of the active users in the system

        Returns:
            List[ `pyinventory.graphql.fragment.user.UserFragment` ]

        Raises:
            FailedOperationException: internal inventory error

        Example:
            ```
            users = client.get_active_users()
            for user in users:
                print(user.email)
            ```
    """
    users = get_users(client=client)
    return [user for user in users if user.status == UserStatus.ACTIVE]
