#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import Dict, List, Optional, Tuple

from .._utils import PropertyValue, format_properties, get_graphql_property_inputs
from ..client import SymphonyClient
from ..common.data_enum import Entity
from ..consts import (
    Customer,
    EquipmentPort,
    EquipmentPortDefinition,
    Link,
    Service,
    ServiceEndpoint,
    ServiceType,
)
from ..exceptions import EntityNotFoundError
from ..graphql.add_service_endpoint_input import AddServiceEndpointInput
from ..graphql.add_service_endpoint_mutation import AddServiceEndpointMutation
from ..graphql.add_service_link_mutation import AddServiceLinkMutation
from ..graphql.add_service_mutation import AddServiceMutation
from ..graphql.add_service_type_mutation import AddServiceTypeMutation
from ..graphql.remove_service_mutation import RemoveServiceMutation
from ..graphql.remove_service_type_mutation import RemoveServiceTypeMutation
from ..graphql.service_create_data_input import ServiceCreateData
from ..graphql.service_details_query import ServiceDetailsQuery
from ..graphql.service_status_enum import ServiceStatus
from ..graphql.service_type_create_data_input import ServiceTypeCreateData
from ..graphql.service_type_services_query import ServiceTypeServicesQuery
from ..graphql.service_types_query import ServiceTypesQuery


def _populate_service_types(client: SymphonyClient) -> None:
    service_types = ServiceTypesQuery.execute(client).serviceTypes
    if not service_types:
        return
    edges = service_types.edges
    for edge in edges:
        node = edge.node
        if node is not None:
            client.serviceTypes[node.name] = ServiceType(
                name=node.name,
                id=node.id,
                hasCustomer=node.hasCustomer,
                property_types=node.propertyTypes,
            )


def add_service_type(
    client: SymphonyClient,
    name: str,
    hasCustomer: bool,
    properties: List[Tuple[str, str, Optional[PropertyValue], Optional[bool]]],
) -> ServiceType:

    new_property_types = format_properties(properties)
    result = AddServiceTypeMutation.execute(
        client,
        data=ServiceTypeCreateData(
            name=name, hasCustomer=hasCustomer, properties=new_property_types
        ),
    ).addServiceType

    service_type = ServiceType(
        name=result.name,
        id=result.id,
        hasCustomer=result.hasCustomer,
        property_types=result.propertyTypes,
    )
    client.serviceTypes[name] = service_type
    return service_type


def add_service(
    client: SymphonyClient,
    name: str,
    external_id: Optional[str],
    service_type: str,
    customer: Optional[Customer],
    properties_dict: Dict[str, PropertyValue],
    links: List[Link],
) -> Service:
    property_types = client.serviceTypes[service_type].property_types
    properties = get_graphql_property_inputs(property_types, properties_dict)
    service_create_data = ServiceCreateData(
        name=name,
        externalId=external_id,
        serviceTypeId=client.serviceTypes[service_type].id,
        status=ServiceStatus.PENDING,
        customerId=customer.id if customer is not None else None,
        properties=properties,
        upstreamServiceIds=[],
    )
    result = AddServiceMutation.execute(client, data=service_create_data).addService
    for l in links:
        result = AddServiceLinkMutation.execute(
            client, id=result.id, linkId=l.id
        ).addServiceLink
    returned_customer = result.customer
    endpoints = []
    for e in result.endpoints:
        port = e.port
        link = port.link if port is not None else None
        endpoints.append(
            ServiceEndpoint(
                id=e.id,
                port=EquipmentPort(
                    id=port.id,
                    properties=port.properties,
                    definition=EquipmentPortDefinition(
                        id=port.definition.id, name=port.definition.name
                    ),
                    link=Link(
                        link.id,
                        properties=link.properties,
                        service_ids=[s.id for s in link.services],
                    )
                    if link
                    else None,
                )
                if port
                else None,
                # TODO add service_endpoint_type api
                type="1",
            )
        )
    return Service(
        name=result.name,
        id=result.id,
        externalId=result.externalId,
        customer=Customer(
            name=returned_customer.name,
            id=returned_customer.id,
            externalId=returned_customer.externalId,
        )
        if returned_customer
        else None,
        endpoints=endpoints,
        links=[
            Link(
                id=l.id, properties=l.properties, service_ids=[s.id for s in l.services]
            )
            for l in result.links
        ],
    )


def add_service_endpoint(
    client: SymphonyClient, service: Service, port: EquipmentPort
) -> None:
    AddServiceEndpointMutation.execute(
        client,
        input=AddServiceEndpointInput(
            id=service.id, portId=port.id, definition="1", equipmentID="1"
        ),
    )


def get_service(client: SymphonyClient, id: str) -> Service:
    result = ServiceDetailsQuery.execute(client, id=id).service
    if result is None:
        raise EntityNotFoundError(entity=Entity.Service, entity_id=id)
    customer = result.customer
    endpoints = []
    for e in result.endpoints:
        port = e.port
        link = port.link if port is not None else None
        endpoints.append(
            ServiceEndpoint(
                id=e.id,
                port=EquipmentPort(
                    id=port.id,
                    properties=port.properties,
                    definition=EquipmentPortDefinition(
                        id=port.definition.id, name=port.definition.name
                    ),
                    link=Link(
                        id=link.id,
                        properties=link.properties,
                        service_ids=[s.id for s in link.services],
                    )
                    if link
                    else None,
                )
                if port is not None
                else None,
                # TODO add service_endpoint_type api
                type="1",
            )
        )
    return Service(
        name=result.name,
        id=result.id,
        externalId=result.externalId,
        customer=Customer(
            name=customer.name, id=customer.id, externalId=customer.externalId
        )
        if customer is not None
        else None,
        endpoints=endpoints,
        links=[
            Link(
                id=l.id, properties=l.properties, service_ids=[s.id for s in l.services]
            )
            for l in result.links
        ],
    )


def delete_service_type_with_services(
    client: SymphonyClient, service_type: ServiceType
) -> None:
    service_type_with_services = ServiceTypeServicesQuery.execute(
        client, id=service_type.id
    ).serviceType
    if not service_type_with_services:
        raise EntityNotFoundError(entity=Entity.ServiceType, entity_id=service_type.id)
    services = service_type_with_services.services
    for service in services:
        RemoveServiceMutation.execute(client, id=service.id)
    RemoveServiceTypeMutation.execute(client, id=service_type.id)
