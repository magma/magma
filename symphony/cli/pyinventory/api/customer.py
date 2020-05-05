#!/usr/bin/env python3

from typing import List, Optional

from pysymphony import SymphonyClient

from ..common.data_class import Customer
from ..graphql.add_customer_input import AddCustomerInput
from ..graphql.add_customer_mutation import AddCustomerMutation
from ..graphql.customers_query import CustomersQuery
from ..graphql.remove_customer_mutation import RemoveCustomerMutation


def add_customer(
    client: SymphonyClient, name: str, external_id: Optional[str]
) -> Customer:
    """This function adds Customer.

        Args:
            name (str): name for the Customer
            external_id (Optional[str]): external ID for the Customer

        Returns:
            `pyinventory.common.data_class.Customer` object

        Example:
            ```
            new_customers = client.add_customer(name="new_customer")
            ```
            or
            ```
            new_customers = client.add_customer(name="new_customer", external_id="12345678")
            ```
    """
    customer_input = AddCustomerInput(name=name, externalId=external_id)
    result = AddCustomerMutation.execute(client, input=customer_input)
    return Customer(name=result.name, id=result.id, externalId=result.externalId)


def get_all_customers(client: SymphonyClient) -> List[Customer]:

    """This function returns all Customers.

        Returns:
            List[ `pyinventory.common.data_class.Customer` ]

        Example:
            ```
            customers = client.get_all_customers()
            ```
    """
    customers = CustomersQuery.execute(client)
    if not customers:
        return []
    result = []
    for customer in customers.edges:
        node = customer.node
        if node:
            result.append(
                Customer(name=node.name, id=node.id, externalId=node.externalId)
            )
    return result


def delete_customer(client: SymphonyClient, customer: Customer) -> None:
    """This function delete Customer.

        Args:
            customer ( `pyinventory.common.data_class.Customer` ): customer object

        Example:
            ```
            client.delete_customer(customer)
            ```
    """
    RemoveCustomerMutation.execute(client, id=customer.id)
