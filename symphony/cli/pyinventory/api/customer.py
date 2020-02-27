#!/usr/bin/env python3
# pyre-strict

from typing import List, Optional

from ..client import SymphonyClient
from ..consts import Customer
from ..graphql.add_customer_input import AddCustomerInput
from ..graphql.add_customer_mutation import AddCustomerMutation
from ..graphql.customers_query import CustomersQuery
from ..graphql.remove_customer_mutation import RemoveCustomerMutation


def add_customer(
    client: SymphonyClient, name: str, external_id: Optional[str]
) -> Customer:
    customer_input = AddCustomerInput(name=name, externalId=external_id)
    result = AddCustomerMutation.execute(client, input=customer_input).addCustomer
    return Customer(name=result.name, id=result.id, externalId=result.externalId)


def get_all_customers(client: SymphonyClient) -> List[Customer]:
    customer_edges = CustomersQuery.execute(client).customers.edges

    customers = [edge.node for edge in customer_edges]

    return [
        Customer(name=customer.name, id=customer.id, externalId=customer.externalId)
        for customer in customers
    ]


def delete_customer(client: SymphonyClient, customer: Customer) -> None:
    RemoveCustomerMutation.execute(client, id=customer.id)
