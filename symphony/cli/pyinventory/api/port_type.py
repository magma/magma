#!/usr/bin/env python3
# pyre-strict
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import List, Tuple

from dacite import Config, from_dict
from gql.gql.client import OperationException

from .._utils import PropertyValue, format_properties
from ..consts import EquipmentPortType
from ..exceptions import EntityNotFoundError
from ..graphql.add_equipment_port_type_input import AddEquipmentPortTypeInput
from ..graphql.add_equipment_port_type_mutation import AddEquipmentPortTypeMutation
from ..graphql.equipment_port_type_query import EquipmentPortTypeQuery
from ..graphql.property_type_input import PropertyTypeInput
from ..graphql.remove_equipment_port_type_mutation import (
    RemoveEquipmentPortTypeMutation,
)
from ..graphql_client import GraphqlClient
from ..reporter import FailedOperationException


ADD_EQUIPMENT_PORT_TYPE_MUTATION_NAME = "addEquipmentPortType"


def add_equipment_port_type(
    client: GraphqlClient,
    name: str,
    properties: List[Tuple[str, str, PropertyValue, bool]],
    link_properties: List[Tuple[str, str, PropertyValue, bool]],
) -> EquipmentPortType:
    """This function creates an equipment port type.

        Args:
            name (str): 
            properties: (list of tuple(str, str, PropertyValue, bool)): Optional, where
                str - port type name
                str - enum["string", "int", "bool", "float", "date", "enum", "range", 
                "email", "gps_location", "equipment", "location", "service", "datetime_local"]
                PropertyValue - default property value
                bool - fixed value
            link_properties: (list of tuple(str, str, PropertyValue, bool)): Optional, where
                str - port type name
                str - enum["string", "int", "bool", "float", "date", "enum", "range", 
                "email", "gps_location", "equipment", "location", "service", "datetime_local"]
                PropertyValue - default property value
                bool - fixed value

        Returns: 
            EquipmentPortType object

        Raises: 
            FailedOperationException

        Example:
            port_type1 = self.client.add_equipment_port_type(
                "port type 1",
                [("port property", "string", None, True)],
                [("link port property", "string", None, True)],
            )
    """

    new_property_types = format_properties(properties)
    new_link_property_types = format_properties(link_properties)
    add_equipment_port_type_input = {
        "name": name,
        "properties": new_property_types,
        "linkProperties": new_link_property_types,
    }

    try:
        result = AddEquipmentPortTypeMutation.execute(
            client,
            AddEquipmentPortTypeInput(
                name=name,
                properties=[
                    from_dict(
                        data_class=PropertyTypeInput, data=p, config=Config(strict=True)
                    )
                    for p in new_property_types
                ],
                linkProperties=[
                    from_dict(
                        data_class=PropertyTypeInput, data=p, config=Config(strict=True)
                    )
                    for p in new_link_property_types
                ],
            ),
        ).__dict__[ADD_EQUIPMENT_PORT_TYPE_MUTATION_NAME]
        client.reporter.log_successful_operation(
            ADD_EQUIPMENT_PORT_TYPE_MUTATION_NAME, add_equipment_port_type_input
        )
    except OperationException as e:
        raise FailedOperationException(
            client.reporter,
            e.err_msg,
            e.err_id,
            ADD_EQUIPMENT_PORT_TYPE_MUTATION_NAME,
            add_equipment_port_type_input,
        )
    return EquipmentPortType(id=result.id, name=result.name)


def get_equipment_port_type(
    client: GraphqlClient, equipment_port_type_id: str
) -> EquipmentPortType:
    """This function returns an equipment port type.
        It can get only the requested equipment port type ID

        Args:
            equipment_port_type_id (str): equipment port type ID

        Returns: 
            EquipmentPortType object

        Raises: 
            EntityNotFoundError for not found entity

        Example:
            port_type = self.client.get_equipment_port_type(self.port_type1.id)
    """
    result = EquipmentPortTypeQuery.execute(client, id=equipment_port_type_id).port_type
    if not result:
        raise EntityNotFoundError(
            entity="Equipment Port Type", entity_id=equipment_port_type_id
        )
    return EquipmentPortType(id=result.id, name=result.name)


def delete_equipment_port_type(
    client: GraphqlClient, equipment_port_type_id: str
) -> None:
    """This function deletes an equipment port type.
        It can get only the requested equipment port type ID

        Args:
            equipment_port_type_id (str): equipment port type ID

        Example:
            client.delete_equipment_port_type(self.port_type1.id)
    """
    RemoveEquipmentPortTypeMutation.execute(client, id=equipment_port_type_id)
