#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


from pyinventory.api.equipment import add_equipment
from pyinventory.api.equipment_type import add_equipment_type
from pyinventory.api.link import add_link, get_link_in_port_of_equipment
from pyinventory.api.location import add_location
from pyinventory.api.location_type import add_location_type
from pyinventory.api.port import edit_link_properties, get_port
from pyinventory.api.port_type import add_equipment_port_type
from pyinventory.common.data_class import PropertyDefinition
from pyinventory.exceptions import PortAlreadyOccupiedException
from pyinventory.graphql.enum.property_kind import PropertyKind
from pysymphony import SymphonyClient

from ..utils.base_test import BaseTest
from ..utils.grpc.rpc_pb2_grpc import TenantServiceStub


class TestLink(BaseTest):
    def __init__(
        self, testName: str, client: SymphonyClient, stub: TenantServiceStub
    ) -> None:
        super().__init__(testName, client, stub)

    def setUp(self) -> None:
        super().setUp()
        add_equipment_port_type(
            self.client,
            name="port type 1",
            properties=[
                PropertyDefinition(
                    property_name="port property",
                    property_kind=PropertyKind.string,
                    default_value="port property value",
                    is_fixed=False,
                )
            ],
            link_properties=[
                PropertyDefinition(
                    property_name="link property",
                    property_kind=PropertyKind.string,
                    default_value="link property value",
                    is_fixed=False,
                )
            ],
        )
        add_location_type(
            client=self.client,
            name="City",
            properties=[
                ("Mayor", "string", None, True),
                ("Contact", "email", None, True),
            ],
        )
        self.location = add_location(
            client=self.client,
            location_hirerchy=[("City", "Lima")],
            properties_dict={"Mayor": "Bernard King", "Contact": "limacity@peru.pe"},
            lat=10,
            long=20,
        )
        add_equipment_type(
            client=self.client,
            name="Tp-Link T1600G",
            category="Router",
            properties=[("IP", "string", None, True)],
            ports_dict={"Port 1": "port type 1", "Port 2": "port type 1"},
            position_list=[],
        )
        self.equipment1 = add_equipment(
            client=self.client,
            name="TPLinkRouter1",
            equipment_type="Tp-Link T1600G",
            location=self.location,
            properties_dict={"IP": "192.688.0.1"},
        )
        self.equipment2 = add_equipment(
            client=self.client,
            name="TPLinkRouter2",
            equipment_type="Tp-Link T1600G",
            location=self.location,
            properties_dict={"IP": "192.688.0.2"},
        )
        self.equipment3 = add_equipment(
            client=self.client,
            name="TPLinkRouter3",
            equipment_type="Tp-Link T1600G",
            location=self.location,
            properties_dict={"IP": "192.688.0.2"},
        )

    def test_add_link(self) -> None:
        link = add_link(
            client=self.client,
            equipment_a=self.equipment1,
            port_name_a="Port 1",
            equipment_b=self.equipment2,
            port_name_b="Port 1",
        )
        fetched_link1 = get_link_in_port_of_equipment(
            client=self.client, equipment=self.equipment1, port_name="Port 1"
        )
        fetched_link2 = get_link_in_port_of_equipment(
            client=self.client, equipment=self.equipment2, port_name="Port 1"
        )
        self.assertEqual(link, fetched_link1)
        self.assertEqual(link, fetched_link2)

    def test_cannot_create_link_if_port_occupied(self) -> None:
        link = add_link(
            client=self.client,
            equipment_a=self.equipment1,
            port_name_a="Port 1",
            equipment_b=self.equipment2,
            port_name_b="Port 1",
        )
        fetched_link1 = get_link_in_port_of_equipment(
            client=self.client, equipment=self.equipment1, port_name="Port 1"
        )
        fetched_link2 = get_link_in_port_of_equipment(
            client=self.client, equipment=self.equipment2, port_name="Port 1"
        )
        self.assertEqual(link, fetched_link1)
        self.assertEqual(link, fetched_link2)

        self.assertRaises(
            PortAlreadyOccupiedException,
            add_link,
            self.client,
            self.equipment2,
            "Port 1",
            self.equipment3,
            "Port 1",
        )

        link = add_link(
            client=self.client,
            equipment_a=self.equipment2,
            port_name_a="Port 2",
            equipment_b=self.equipment3,
            port_name_b="Port 1",
        )
        fetched_link1 = get_link_in_port_of_equipment(
            client=self.client, equipment=self.equipment2, port_name="Port 2"
        )
        fetched_link2 = get_link_in_port_of_equipment(
            client=self.client, equipment=self.equipment3, port_name="Port 1"
        )
        self.assertEqual(link, fetched_link1)
        self.assertEqual(link, fetched_link2)

    def test_edit_link_properties(self) -> None:
        add_link(
            client=self.client,
            equipment_a=self.equipment1,
            port_name_a="Port 1",
            equipment_b=self.equipment2,
            port_name_b="Port 1",
        )
        fetched_port = get_port(
            client=self.client, equipment=self.equipment1, port_name="Port 1"
        )
        edit_link_properties(
            client=self.client,
            equipment=self.equipment1,
            port_name="Port 1",
            new_link_properties={"link property": "test_link_property"},
        )
        fetched_port = get_port(
            client=self.client, equipment=self.equipment1, port_name="Port 1"
        )

        link = fetched_port.link
        link_properties = link.properties if link else []
        self.assertEqual(1, len(link_properties))
        link_property = link_properties[0]
        link_property_type = link_property.propertyType if link_property else None
        property_name = link_property_type.name if link_property_type else None
        value = link_property.stringValue if link_property else None

        self.assertEqual(property_name, "link property")
        self.assertEqual(value, "test_link_property")
