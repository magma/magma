#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from pyinventory.api.customer import add_customer
from pyinventory.api.equipment import add_equipment
from pyinventory.api.equipment_type import add_equipment_type
from pyinventory.api.link import add_link
from pyinventory.api.location import add_location
from pyinventory.api.location_type import add_location_type
from pyinventory.api.port_type import add_equipment_port_type
from pyinventory.api.service import (
    add_service,
    add_service_endpoint,
    add_service_link,
    get_service,
    get_service_endpoints,
    get_service_links,
)
from pyinventory.api.service_type import add_service_type, edit_service_type
from pyinventory.common.cache import SERVICE_TYPES
from pyinventory.common.data_class import PropertyDefinition, ServiceEndpointDefinition
from pyinventory.graphql.enum.property_kind import PropertyKind
from pysymphony import SymphonyClient

from ..utils.base_test import BaseTest
from ..utils.grpc.rpc_pb2_grpc import TenantServiceStub


class TestService(BaseTest):
    def __init__(
        self, testName: str, client: SymphonyClient, stub: TenantServiceStub
    ) -> None:
        super().__init__(testName, client, stub)

    def setUp(self) -> None:
        super().setUp()
        self.service_type = add_service_type(
            client=self.client,
            name="Internet Access",
            has_customer=True,
            properties=[
                PropertyDefinition(
                    property_name="Service Package",
                    property_kind=PropertyKind.string,
                    default_raw_value="Public 5G",
                    is_fixed=False,
                ),
                PropertyDefinition(
                    property_name="Address Family",
                    property_kind=PropertyKind.string,
                    default_raw_value=None,
                    is_fixed=False,
                ),
            ],
            endpoint_definitions=[],
        )
        self.service = add_service(
            client=self.client,
            name="Room 201 Internet Access",
            external_id="S3232",
            service_type=self.service_type.name,
            customer=None,
            properties_dict={"Address Family": "v4"},
        )
        self.customer = add_customer(
            client=self.client, name="Donald", external_id="S322"
        )
        self.service_with_customer = add_service(
            client=self.client,
            name="Room 202 Internet Access",
            external_id="S32325",
            service_type=self.service_type.name,
            customer=self.customer,
            properties_dict={"Address Family": "v4"},
        )

    def test_service_created(self) -> None:
        fetched_service = get_service(client=self.client, id=self.service.id)
        self.assertEqual(fetched_service, self.service)

    def test_service_with_customer_created(self) -> None:
        fetched_service = get_service(
            client=self.client, id=self.service_with_customer.id
        )
        self.assertEqual(fetched_service, self.service_with_customer)
        fetched_customer = fetched_service.customer
        self.assertNotEqual(fetched_customer, None)
        self.assertEqual(fetched_customer, self.customer)

    def test_service_endpoint_added(self) -> None:
        add_equipment_port_type(
            self.client,
            name="port type 1",
            properties=[
                PropertyDefinition(
                    property_name="port property",
                    property_kind=PropertyKind.string,
                    default_raw_value="port property value",
                    is_fixed=False,
                )
            ],
            link_properties=[
                PropertyDefinition(
                    property_name="link property",
                    property_kind=PropertyKind.string,
                    default_raw_value="link property value",
                    is_fixed=False,
                )
            ],
        )
        add_location_type(
            client=self.client,
            name="Room",
            properties=[
                PropertyDefinition(
                    property_name="Contact",
                    property_kind=PropertyKind.email,
                    default_raw_value=None,
                    is_fixed=False,
                )
            ],
        )
        location = add_location(
            client=self.client,
            location_hirerchy=[("Room", "Room 201")],
            properties_dict={"Contact": "user@google.com"},
            lat=10,
            long=20,
        )
        equipment_type = add_equipment_type(
            client=self.client,
            name="Tp-Link T1600G",
            category="Router",
            properties=[
                PropertyDefinition(
                    property_name="IP",
                    property_kind=PropertyKind.string,
                    default_raw_value=None,
                    is_fixed=False,
                )
            ],
            ports_dict={"Port 1": "port type 1", "Port 2": "port type 1"},
            position_list=[],
        )
        router1 = add_equipment(
            client=self.client,
            name="TPLinkRouter1",
            equipment_type=equipment_type.name,
            location=location,
            properties_dict={"IP": "192.688.0.1"},
        )

        edit_service_type(
            client=self.client,
            service_type=self.service_type,
            new_endpoints=[
                ServiceEndpointDefinition(
                    id=None,
                    name="EndpointDefinition",
                    role="CPE",
                    endpoint_definition_index=0,
                    equipment_type_id=equipment_type.id,
                )
            ],
        )
        endpoint_definitions = SERVICE_TYPES[
            self.service_type.name
        ].endpoint_definitions
        endpoints = get_service_endpoints(
            client=self.client, service_id=self.service.id
        )
        self.assertFalse(endpoints)
        endpoint_definition_id = endpoint_definitions[0].id
        add_service_endpoint(
            client=self.client,
            service=self.service,
            equipment_id=router1.id,
            endpoint_definition_id=endpoint_definition_id,
        )
        endpoints = get_service_endpoints(
            client=self.client, service_id=self.service.id
        )
        self.assertEqual(len(endpoints), 1)

    def test_service_link_added(self) -> None:
        add_equipment_port_type(
            self.client,
            name="port type 1",
            properties=[
                PropertyDefinition(
                    property_name="port property",
                    property_kind=PropertyKind.string,
                    default_raw_value="port property value",
                    is_fixed=False,
                )
            ],
            link_properties=[
                PropertyDefinition(
                    property_name="link property",
                    property_kind=PropertyKind.string,
                    default_raw_value="link property value",
                    is_fixed=False,
                )
            ],
        )
        add_location_type(
            client=self.client,
            name="Room",
            properties=[
                PropertyDefinition(
                    property_name="Contact",
                    property_kind=PropertyKind.email,
                    default_raw_value=None,
                    is_fixed=False,
                )
            ],
        )
        location = add_location(
            client=self.client,
            location_hirerchy=[("Room", "Room 201")],
            properties_dict={"Contact": "user@google.com"},
            lat=10,
            long=20,
        )
        equipment_type = add_equipment_type(
            client=self.client,
            name="Tp-Link T1600G",
            category="Router",
            properties=[
                PropertyDefinition(
                    property_name="IP",
                    property_kind=PropertyKind.string,
                    default_raw_value=None,
                    is_fixed=False,
                )
            ],
            ports_dict={"Port 1": "port type 1", "Port 2": "port type 1"},
            position_list=[],
        )
        router1 = add_equipment(
            client=self.client,
            name="TPLinkRouter1",
            equipment_type=equipment_type.name,
            location=location,
            properties_dict={"IP": "192.688.0.1"},
        )
        router2 = add_equipment(
            client=self.client,
            name="TPLinkRouter2",
            equipment_type=equipment_type.name,
            location=location,
            properties_dict={"IP": "192.688.0.2"},
        )
        router3 = add_equipment(
            client=self.client,
            name="TPLinkRouter3",
            equipment_type=equipment_type.name,
            location=location,
            properties_dict={"IP": "192.688.0.3"},
        )
        link1 = add_link(
            client=self.client,
            equipment_a=router1,
            port_name_a="Port 1",
            equipment_b=router2,
            port_name_b="Port 1",
        )
        link2 = add_link(
            client=self.client,
            equipment_a=router2,
            port_name_a="Port 2",
            equipment_b=router3,
            port_name_b="Port 1",
        )
        links = get_service_links(client=self.client, service_id=self.service.id)
        self.assertFalse(links)
        for link in [link1, link2]:
            add_service_link(
                client=self.client, service_id=self.service.id, link_id=link.id
            )
        links = get_service_links(client=self.client, service_id=self.service.id)
        self.assertEqual(len(links), 2)
