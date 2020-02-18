#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import List, Tuple

from dacite import Config, from_dict
from gql.gql.client import OperationException

from .._utils import PropertyValue, format_properties
from ..consts import EquipmentPortType
from ..graphql.add_equipment_port_type_mutation import (
    AddEquipmentPortTypeInput,
    AddEquipmentPortTypeMutation,
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

    formated_property_types = format_properties(properties)
    formated_link_property_types = format_properties(link_properties)

    add_equipment_port_type_input = {
        "name": name,
        "properties": formated_property_types,
        "linkProperties": formated_link_property_types,
    }

    try:
        result = AddEquipmentPortTypeMutation.execute(
            client,
            AddEquipmentPortTypeInput(
                name=name,
                properties=[
                    from_dict(
                        data_class=AddEquipmentPortTypeInput.PropertyTypeInput,
                        data=p,
                        config=Config(strict=True),
                    )
                    for p in formated_property_types
                ],
                linkProperties=[
                    from_dict(
                        data_class=AddEquipmentPortTypeInput.PropertyTypeInput,
                        data=p,
                        config=Config(strict=True),
                    )
                    for p in formated_link_property_types
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
