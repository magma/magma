#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import List, Mapping, Optional

from pysymphony import SymphonyClient

from .._utils import PropertyValue, get_graphql_property_inputs
from ..common.cache import SERVICE_TYPES
from ..common.data_class import Customer, Link, Service, ServiceEndpoint
from ..common.data_enum import Entity
from ..exceptions import EntityNotFoundError
from ..graphql.enum.service_status import ServiceStatus
from ..graphql.input.add_service_endpoint import AddServiceEndpointInput
from ..graphql.input.service_create_data import ServiceCreateData
from ..graphql.mutation.add_service import AddServiceMutation
from ..graphql.mutation.add_service_endpoint import AddServiceEndpointMutation
from ..graphql.mutation.add_service_link import AddServiceLinkMutation
from ..graphql.query.service_details import ServiceDetailsQuery
from ..graphql.query.service_endpoints import ServiceEndpointsQuery
from ..graphql.query.service_links import ServiceLinksQuery


def add_service(
    client: SymphonyClient,
    name: str,
    external_id: Optional[str],
    service_type: str,
    customer: Optional[Customer],
    properties_dict: Optional[Mapping[str, PropertyValue]],
) -> Service:
    """This function creates service.

        Args:
            name (str): service name
            external_id (Optional[str]): service external ID
            service_type (str): existing service type name
            customer (Optional[ `pyinventory.common.data_class.Customer` ]): existing customer object

            properties_dict (Optional[ Mapping[ str, PropertyValue ] ]): dictionary of property name to property value
            - str - property name
            - PropertyValue - new value of the same type for this property

        Returns:
            `pyinventory.common.data_class.Service`

        Raises:
            FailedOperationException: internal inventory error

        Example:
            ```
            service = client.add_service(
                name="Room 202 Internet Access",
                external_id="S32325",
                service_type=self.service_type.name,
                customer=None,
                properties_dict={"Address Family": "v4"},
            )
            ```
    """
    properties = []
    if properties_dict is not None:
        property_types = SERVICE_TYPES[service_type].property_types
        properties = get_graphql_property_inputs(property_types, properties_dict)
    service_create_data = ServiceCreateData(
        name=name,
        externalId=external_id,
        serviceTypeId=SERVICE_TYPES[service_type].id,
        status=ServiceStatus.PENDING,
        customerId=customer.id if customer is not None else None,
        properties=properties,
        upstreamServiceIds=[],
    )
    result = AddServiceMutation.execute(client, data=service_create_data)
    if customer is not None:
        customer = Customer(
            id=customer.id, name=customer.name, external_id=customer.external_id
        )
    return Service(
        id=result.id,
        name=result.name,
        external_id=result.externalId,
        service_type_name=result.serviceType.name,
        customer=customer,
        properties=result.properties,
    )


def get_service(client: SymphonyClient, id: str) -> Service:
    """This function returns service by ID.

        Args:
            id (str): existing service ID

        Returns:
            `pyinventory.common.data_class.Service`

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: service does not exist
            FailedOperationException: internal inventory error

        Example:
            ```
            service = client.get_service(id="12345")
            ```
    """
    result = ServiceDetailsQuery.execute(client, id=id)
    if result is None:
        raise EntityNotFoundError(entity=Entity.Service, entity_id=id)
    customer_result = result.customer if result.customer is not None else None
    customer: Optional[Customer] = None
    if customer_result is not None:
        customer = Customer(
            id=customer_result.id,
            name=customer_result.name,
            external_id=customer_result.externalId,
        )
    return Service(
        id=result.id,
        name=result.name,
        external_id=result.externalId if result.externalId else None,
        service_type_name=result.serviceType.name,
        customer=customer,
        properties=result.properties,
    )


def get_service_endpoints(
    client: SymphonyClient, service_id: str
) -> List[ServiceEndpoint]:
    """This function returns service endpoints list.

        Args:
            service_id (str): existing service ID

        Returns:
            List[ `pyinventory.common.data_class.ServiceEndpoint` ]

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: service does not exist
            FailedOperationException: internal inventory error

        Example:
            ```
            endpoints = client.get_service_endpoint_definitions(id="service_id")
            ```
    """
    service_data = ServiceEndpointsQuery.execute(client, id=service_id)

    if not service_data:
        raise EntityNotFoundError(entity=Entity.Service, entity_id=service_id)

    return [
        ServiceEndpoint(
            id=endpoint.id,
            equipment_id=endpoint.equipment.id,
            service_id=service_id,
            definition_id=endpoint.definition.id,
        )
        for endpoint in service_data.endpoints
    ]


def add_service_endpoint(
    client: SymphonyClient,
    service: Service,
    equipment_id: str,
    endpoint_definition_id: str,
) -> None:
    """This function adds existing endpoint to existing service.

        Args:
            service (str): existing service object
            equipment_id (str): existing equipment ID
            endpoint_definition_id (str): existing endpoint definition ID

        Raises:
            FailedOperationException: internal inventory error

        Example:
            ```
            service = client.get_service(id="service_id")
            location = client.get_location(location_hirerchy=[("Country", "LS_IND_Prod_Copy")])
            equipment = client.get_equipment(name="indProdCpy1_AIO", location=location)
            client.add_service_endpoint(
                service=service,
                equipment_id=equipment.id,
                endpoint_definition_id="endpoint_definition_id,
            )
            ```
    """
    endpoint_definition_ids = [
        ed.id for ed in SERVICE_TYPES[service.service_type_name].endpoint_definitions
    ]

    if endpoint_definition_id not in endpoint_definition_ids:
        raise EntityNotFoundError(
            entity=Entity.ServiceEndpointDefinition, entity_id=endpoint_definition_id
        )

    AddServiceEndpointMutation.execute(
        client,
        input=AddServiceEndpointInput(
            id=service.id, definition=endpoint_definition_id, equipmentID=equipment_id
        ),
    )


def get_service_links(client: SymphonyClient, service_id: str) -> List[Link]:
    """This function returns list of Links.

        Args:
            service_id (str): existing service ID

        Returns:
            List[ `pyinventory.common.data_class.Link` ]

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: service does not exist
            FailedOperationException: internal inventory error

        Example:
            ```
            links = client.get_service_links(id="service_id")
            ```
    """
    service_data = ServiceLinksQuery.execute(client, id=service_id)

    if not service_data:
        raise EntityNotFoundError(entity=Entity.Service, entity_id=service_id)

    return [
        Link(
            id=link.id,
            properties=link.properties,
            service_ids=[s.id for s in link.services],
        )
        for link in service_data.links
    ]


def add_service_link(client: SymphonyClient, service_id: str, link_id: str) -> None:
    """This function adds existing link to existing service.

        Args:
            service_id (str): existing service ID
            link_id (str): existing link ID

        Raises:
            FailedOperationException: internal inventory error

        Example:
            ```
            client.add_service_link(service_id=service.id, link_id=link.id)
            ```
    """
    AddServiceLinkMutation.execute(client, id=service_id, linkId=link_id)
