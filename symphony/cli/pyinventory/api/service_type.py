#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import List, Optional

from pysymphony import SymphonyClient

from .._utils import format_property_definitions
from ..common.cache import SERVICE_TYPES
from ..common.data_class import PropertyDefinition, ServiceType
from ..common.data_enum import Entity
from ..exceptions import EntityNotFoundError
from ..graphql.input.service_type_create_data import ServiceTypeCreateData
from ..graphql.mutation.add_service_type import AddServiceTypeMutation
from ..graphql.query.service_types import ServiceTypesQuery


def _populate_service_types(client: SymphonyClient) -> None:
    service_types = ServiceTypesQuery.execute(client)
    if not service_types:
        return
    edges = service_types.edges
    for edge in edges:
        node = edge.node
        if node is not None:
            SERVICE_TYPES[node.name] = ServiceType(
                id=node.id,
                name=node.name,
                has_customer=node.hasCustomer,
                property_types=node.propertyTypes,
            )


def add_service_type(
    client: SymphonyClient,
    name: str,
    has_customer: bool,
    properties: Optional[List[PropertyDefinition]],
) -> ServiceType:
    """This function creates new service type.

        Args:
            name (str): service type name
            has_customer (bool): flag for customenr existance
            properties: (Optional[ List[ `pyinventory.common.data_class.PropertyDefinition` ] ]): list of property definitions

        Returns:
            `pyinventory.common.data_class.ServiceType`

        Raises:
            FailedOperationException: internal inventory error

        Example:
            ```
            service_type = client.add_service_type(
                client=self.client,
                name="Internet Access",
                has_customer=True,
                properties=[
                    PropertyDefinition(
                        property_name="Service Package",
                        property_kind=PropertyKind.string,
                        default_value="Public 5G",
                        is_fixed=True,
                    ),
                    PropertyDefinition(
                        property_name="Address Family",
                        property_kind=PropertyKind.string,
                        default_value=None,
                        is_fixed=True,
                    ),
                )
            ```
    """

    formated_property_types = None
    if properties is not None:
        formated_property_types = format_property_definitions(properties=properties)
    result = AddServiceTypeMutation.execute(
        client,
        data=ServiceTypeCreateData(
            name=name, hasCustomer=has_customer, properties=formated_property_types
        ),
    )
    service_type = ServiceType(
        id=result.id,
        name=result.name,
        has_customer=result.hasCustomer,
        property_types=result.propertyTypes,
    )
    SERVICE_TYPES[name] = service_type
    return service_type


def get_service_type(client: SymphonyClient, service_type_id: str) -> ServiceType:
    """Get service type by ID.

        Args:
            service_type_id (str): service type ID

        Returns:
            `pyinventory.common.data_class.ServiceType`

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: if service type with id=`service_type_id` does not found

        Example:
            ```
            service_type = client.get_service_type(
                service_type_id="12345",
            )
            ```
    """
    for _, service_type in SERVICE_TYPES.items():
        if service_type.id == service_type_id:
            return service_type

    raise EntityNotFoundError(entity=Entity.ServiceType, entity_id=service_type_id)
