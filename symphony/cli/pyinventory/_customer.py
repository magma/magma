#!/usr/bin/env python3
# pyre-strict

from typing import List, Optional

from .consts import Customer
from .graphql.add_customer_mutation import AddCustomerInput, AddCustomerMutation
from .graphql.customers_query import CustomersQuery
from .graphql.remove_customer_mutation import RemoveCustomerMutation
from .graphql_client import GraphqlClient


def add_customer(
    client: GraphqlClient, name: str, external_id: Optional[str]
) -> Customer:
    customer_input = AddCustomerInput(name=name, externalId=external_id)
    result = AddCustomerMutation.execute(client, input=customer_input).addCustomer
    return Customer(name=result.name, id=result.id, externalId=result.externalId)


def get_all_customers(client: GraphqlClient) -> List[Customer]:
    customer_edges = CustomersQuery.execute(client).customers.edges

    customers = [edge.node for edge in customer_edges]

    return [
        Customer(name=customer.name, id=customer.id, externalId=customer.externalId)
        for customer in customers
    ]


def delete_customer(client: GraphqlClient, customer: Customer) -> None:
    RemoveCustomerMutation.execute(client, id=customer.id)
