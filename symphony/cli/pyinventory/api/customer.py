#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import Iterator, Optional

from pysymphony import SymphonyClient

from ..common.data_class import Customer
from ..graphql.input.add_customer import AddCustomerInput
from ..graphql.mutation.add_customer import AddCustomerMutation
from ..graphql.mutation.remove_customer import RemoveCustomerMutation
from ..graphql.query.customers import CustomersQuery


def add_customer(
    client: SymphonyClient, name: str, external_id: Optional[str]
) -> Customer:
    """This function adds Customer.

        :param name: Customer name
        :type name: str
        :param external_id: Customer external ID
        :type external_id: str, optional

        :return: Customer object
        :rtype: :class:`~pyinventory.common.data_class.Customer`

        **Example 1**

        .. code-block:: python

            new_customers = client.add_customer(name="new_customer")

        **Example 2**

        .. code-block:: python

            new_customers = client.add_customer(
                name="new_customer",
                external_id="12345678"
            )
    """
    customer_input = AddCustomerInput(name=name, externalId=external_id)
    result = AddCustomerMutation.execute(client, input=customer_input)
    return Customer(name=result.name, id=result.id, external_id=result.externalId)


def get_all_customers(client: SymphonyClient) -> Iterator[Customer]:

    """This function returns all Customers.

        :return: Customers Iterator
        :rtype: Iterator[ :class:`~pyinventory.common.data_class.Customer` ]

        **Example**

        .. code-block:: python

            customers = client.get_all_customers()
    """
    customers = CustomersQuery.execute(client)
    if not customers:
        return
    for customer in customers.edges:
        node = customer.node
        if node:
            yield Customer(name=node.name, id=node.id, external_id=node.externalId)


def delete_customer(client: SymphonyClient, customer: Customer) -> None:
    """This function delete Customer.

        :param name: Customer name
        :type name: :class:`~pyinventory.common.data_class.Customer`
        :rtype: None

        **Example**

        .. code-block:: python

            client.delete_customer(customer)
    """
    RemoveCustomerMutation.execute(client, id=customer.id)
