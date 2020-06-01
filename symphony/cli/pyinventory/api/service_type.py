#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import Dict, List, Optional

from pysymphony import SymphonyClient

from .._utils import format_property_definitions, get_graphql_property_type_inputs
from ..common.cache import SERVICE_TYPES
from ..common.data_class import (
    PropertyDefinition,
    PropertyValue,
    ServiceEndpointDefinition,
    ServiceType,
)
from ..common.data_enum import Entity
from ..exceptions import EntityNotFoundError
from ..graphql.input.service_type_create_data import ServiceTypeCreateData
from ..graphql.input.service_type_edit_data import (
    ServiceEndpointDefinitionInput,
    ServiceTypeEditData,
)
from ..graphql.mutation.add_service_type import AddServiceTypeMutation
from ..graphql.mutation.edit_service_type import EditServiceTypeMutation
from ..graphql.mutation.remove_service import RemoveServiceMutation
from ..graphql.mutation.remove_service_type import RemoveServiceTypeMutation
from ..graphql.query.service_type_services import ServiceTypeServicesQuery
from ..graphql.query.service_types import ServiceTypesQuery


def _populate_service_types(client: SymphonyClient) -> None:
    service_types = ServiceTypesQuery.execute(client)
    if not service_types:
        return
    edges = service_types.edges
    for edge in edges:
        node = edge.node
        if node is not None:
            definitions = []
            if node.endpointDefinitions:
                definitions = [
                    ServiceEndpointDefinition(
                        id=definition.id,
                        name=definition.name,
                        endpoint_definition_index=definition.index,
                        role=definition.role,
                        equipment_type_id=definition.equipmentType.id,
                    )
                    for definition in node.endpointDefinitions
                ]
            SERVICE_TYPES[node.name] = ServiceType(
                id=node.id,
                name=node.name,
                has_customer=node.hasCustomer,
                property_types=node.propertyTypes,
                endpoint_definitions=definitions,
            )


def add_service_type(
    client: SymphonyClient,
    name: str,
    has_customer: bool,
    properties: Optional[List[PropertyDefinition]],
    endpoint_definitions: Optional[List[ServiceEndpointDefinition]],
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
    definition_inputs = []
    if endpoint_definitions:
        for endpoint in endpoint_definitions:
            definition_inputs.append(
                ServiceEndpointDefinitionInput(
                    name=endpoint.name,
                    role=endpoint.role,
                    index=endpoint.endpoint_definition_index,
                    equipmentTypeID=endpoint.equipment_type_id,
                )
            )
    result = AddServiceTypeMutation.execute(
        client,
        data=ServiceTypeCreateData(
            name=name,
            hasCustomer=has_customer,
            properties=formated_property_types,
            endpoints=definition_inputs,
        ),
    )
    definitions = []
    if result.endpointDefinitions:
        definitions = [
            ServiceEndpointDefinition(
                id=definition.id,
                name=definition.name,
                endpoint_definition_index=definition.index,
                role=definition.role,
                equipment_type_id=definition.equipmentType.id,
            )
            for definition in result.endpointDefinitions
        ]
    service_type = ServiceType(
        id=result.id,
        name=result.name,
        has_customer=result.hasCustomer,
        property_types=result.propertyTypes,
        endpoint_definitions=definitions,
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


def edit_service_type(
    client: SymphonyClient,
    service_type: ServiceType,
    new_name: Optional[str] = None,
    new_has_customer: Optional[bool] = None,
    new_properties: Optional[Dict[str, PropertyValue]] = None,
    new_endpoints: Optional[List[ServiceEndpointDefinition]] = None,
) -> ServiceType:
    """Edit existing service type by ID.

        Args:
            service_type ( `pyinventory.common.data_class.ServiceType` ): existing service type object
            new_name (Optional[ str ]): new name
            new_has_customer (Optional[ bool ]): flag customer existance
            new_properties: (Optional[ Dict[ str, PropertyValue ] ]): dictionary
            - str - property type name
            - PropertyValue - new value of the same type for this property

            new_endpoints (Optional[ List[ `pyinventory.common.data_class.ServiceEndpointDefinition` ] ]): endpoint definitions list

        Returns:
            `pyinventory.common.data_class.ServiceType`

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: if service type with id=`service_type_id` does not found

        Example:
            ```
            service_type = client.edit_service_type(
                service_type=service_type,
                new_name="new service type name",
                new_properties={"existing property name": "new value"},
                new_endpoints=[
                    ServiceEndpointDefinition(
                        id="endpoint_def_id",
                        name="endpoint_def_name",
                        role="endpoint_def_role",
                        index=1,
                    ),
                ],
            )
            ```
    """
    new_name = service_type.name if new_name is None else new_name
    new_has_customer = (
        service_type.has_customer if new_has_customer is None else new_has_customer
    )

    new_property_type_inputs = []
    if new_properties:
        property_types = SERVICE_TYPES[service_type.name].property_types
        new_property_type_inputs = get_graphql_property_type_inputs(
            property_types, new_properties
        )

    new_endpoints_definition_inputs = []
    if new_endpoints:
        for endpoint in new_endpoints:
            new_endpoints_definition_inputs.append(
                ServiceEndpointDefinitionInput(
                    name=endpoint.name,
                    role=endpoint.role,
                    index=endpoint.endpoint_definition_index,
                    equipmentTypeID=endpoint.equipment_type_id,
                )
            )

    result = EditServiceTypeMutation.execute(
        client,
        ServiceTypeEditData(
            id=service_type.id,
            name=new_name,
            hasCustomer=new_has_customer,
            properties=new_property_type_inputs,
            endpoints=new_endpoints_definition_inputs,
        ),
    )
    definitions = []
    if result.endpointDefinitions is not None:
        definitions = [
            ServiceEndpointDefinition(
                id=definition.id,
                name=definition.name,
                endpoint_definition_index=definition.index,
                role=definition.role,
                equipment_type_id=definition.equipmentType.id,
            )
            for definition in result.endpointDefinitions
        ]
    service_type = ServiceType(
        id=result.id,
        name=result.name,
        has_customer=result.hasCustomer,
        property_types=result.propertyTypes,
        endpoint_definitions=definitions,
    )
    SERVICE_TYPES[service_type.name] = service_type
    return service_type


def delete_service_type(client: SymphonyClient, service_type: ServiceType) -> None:
    """This function deletes an service type.
        It can get only the requested service type ID

        Args:
            service_type ( `pyinventory.common.data_class.ServiceType` ): service type object

        Example:
            ```
            client.delete_service_type(service_type_id=service_type.id)
            ```
    """
    RemoveServiceTypeMutation.execute(client, id=service_type.id)
    del SERVICE_TYPES[service_type.name]


def delete_service_type_with_services(
    client: SymphonyClient, service_type: ServiceType
) -> None:
    """Delete service type with existing services.

        Args:
            service_type ( `pyinventory.common.data_class.ServiceType` ): service type object

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: if service_type does not exist

        Example:
            ```
            client.delete_service_type_with_services(service_type=service_type)
            ```
    """
    service_type_with_services = ServiceTypeServicesQuery.execute(
        client, id=service_type.id
    )
    if not service_type_with_services:
        raise EntityNotFoundError(entity=Entity.ServiceType, entity_id=service_type.id)
    services = service_type_with_services.services
    for service in services:
        RemoveServiceMutation.execute(client, id=service.id)
    RemoveServiceTypeMutation.execute(client, id=service_type.id)
    del SERVICE_TYPES[service_type.name]
