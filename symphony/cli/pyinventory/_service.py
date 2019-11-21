#!/usr/bin/env python3
# pyre-strict

from typing import Dict, List, Optional, Tuple, Union

from dacite import from_dict

from ._utils import PropertyValue, _get_properties_to_add, _make_property_types
from .consts import Customer, Equipment, Link, Service, ServiceType
from .graphql.add_service_mutation import AddServiceMutation, ServiceCreateData
from .graphql.add_service_type_mutation import (
    AddServiceTypeMutation,
    PropertyKind,
    ServiceTypeCreateData,
)
from .graphql.remove_service_mutation import RemoveServiceMutation
from .graphql.remove_service_type_mutation import RemoveServiceTypeMutation
from .graphql.service_details_query import ServiceDetailsQuery
from .graphql.service_type_services_query import ServiceTypeServicesQuery
from .graphql.service_types_query import ServiceTypesQuery
from .graphql_client import GraphqlClient


def _populate_service_types(client: GraphqlClient) -> None:
    edges = ServiceTypesQuery.execute(client).serviceTypes.edges
    for edge in edges:
        node = edge.node
        client.serviceTypes[node.name] = ServiceType(
            name=node.name,
            id=node.id,
            hasCustomer=node.hasCustomer,
            propertyTypes=[p.to_dict() for p in node.propertyTypes],
        )


def add_service_type(
    client: GraphqlClient,
    name: str,
    hasCustomer: bool,
    properties: List[Tuple[str, str, PropertyValue, bool]],
) -> ServiceType:
    property_types = _make_property_types(properties)

    def property_type_to_kind(
        key: str, value: Union[str, int, float, bool]
    ) -> Union[str, int, float, bool, PropertyKind]:
        return value if key != "type" else PropertyKind(value)

    new_property_types = [
        {k: property_type_to_kind(k, v) for k, v in property_type.items()}
        for property_type in property_types
    ]
    result = AddServiceTypeMutation.execute(
        client,
        data=ServiceTypeCreateData(
            name=name,
            hasCustomer=hasCustomer,
            properties=[
                from_dict(data_class=ServiceTypeCreateData.PropertyTypeInput, data=p)
                for p in new_property_types
            ],
        ),
    ).addServiceType

    service_type = ServiceType(
        name=result.name,
        id=result.id,
        hasCustomer=result.hasCustomer,
        propertyTypes=[p.to_dict() for p in result.propertyTypes],
    )
    client.serviceTypes[name] = service_type
    return service_type


def add_service(
    client: GraphqlClient,
    name: str,
    external_id: str,
    service_type: str,
    customer: Optional[Customer],
    properties_dict: Dict[str, PropertyValue],
    termination_points: List[Equipment],
    links: List[Link],
) -> Service:
    property_types = client.serviceTypes[service_type].propertyTypes
    properties = _get_properties_to_add(property_types, properties_dict)
    service_create_data = ServiceCreateData(
        name=name,
        externalId=external_id,
        serviceTypeId=client.serviceTypes[service_type].id,
        customerId=customer.id if customer is not None else None,
        properties=properties,
        upstreamServiceIds=[],
        terminationPointIds=[e.id for e in termination_points],
        linkIds=[l.id for l in links],
    )
    result = AddServiceMutation.execute(client, data=service_create_data).addService
    return Service(
        name=result.name,
        id=result.id,
        externalId=result.externalId,
        customer=Customer(
            name=result.customer.name,
            id=result.customer.id,
            externalId=result.customer.externalId,
        )
        if result.customer is not None
        else None,
        terminationPoints=[
            Equipment(name=e.name, id=e.id) for e in result.terminationPoints
        ],
        links=[Link(id=l.id) for l in result.links],
    )


def get_service(client: GraphqlClient, id: str) -> Service:
    result = ServiceDetailsQuery.execute(client, id=id).service
    return Service(
        name=result.name,
        id=result.id,
        externalId=result.externalId,
        customer=Customer(
            name=result.customer.name,
            id=result.customer.id,
            externalId=result.customer.externalId,
        )
        if result.customer is not None
        else None,
        terminationPoints=[
            Equipment(name=e.name, id=e.id) for e in result.terminationPoints
        ],
        links=[Link(id=l.id) for l in result.links],
    )


def delete_service_type_with_services(
    client: GraphqlClient, service_type: ServiceType
) -> None:
    services = ServiceTypeServicesQuery.execute(
        client, id=service_type.id
    ).serviceType.services
    for service in services:
        RemoveServiceMutation.execute(client, id=service.id)
    RemoveServiceTypeMutation.execute(client, id=service_type.id)
