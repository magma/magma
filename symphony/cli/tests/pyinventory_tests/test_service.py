#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


from pyinventory.api.customer import add_customer, delete_customer, get_all_customers
from pyinventory.api.equipment import add_equipment
from pyinventory.api.equipment_type import (
    add_equipment_type,
    delete_equipment_type_with_equipments,
)
from pyinventory.api.link import add_link, get_port
from pyinventory.api.location import add_location
from pyinventory.api.location_type import (
    add_location_type,
    delete_location_type_with_locations,
)
from pyinventory.api.port_type import (
    add_equipment_port_type,
    delete_equipment_port_type,
)
from pyinventory.api.service import (
    add_service,
    add_service_endpoint,
    add_service_type,
    delete_service_type_with_services,
    get_service,
)
from pyinventory.consts import PropertyDefinition
from pyinventory.graphql.property_kind_enum import PropertyKind
from pyinventory.graphql.service_endpoint_role_enum import ServiceEndpointRole

from .utils.base_test import BaseTest


class TestService(BaseTest):
    def setUp(self) -> None:
        self.port_type1 = add_equipment_port_type(
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
        self.service_types_created = []
        self.service_types_created.append(
            add_service_type(
                client=self.client,
                name="Internet Access",
                hasCustomer=True,
                properties=[
                    ("Service Package", "string", "Public 5G", True),
                    ("Address Family", "string", None, True),
                ],
            )
        )
        self.location_types_created = []
        self.location_types_created.append(
            add_location_type(
                client=self.client,
                name="Room",
                properties=[("Contact", "email", None, True)],
            )
        )
        self.equipment_types_created = []
        self.equipment_types_created.append(
            add_equipment_type(
                client=self.client,
                name="Tp-Link T1600G",
                category="Router",
                properties=[("IP", "string", None, True)],
                ports_dict={"Port 1": "port type 1", "Port 2": "port type 1"},
                position_list=[],
            )
        )

    def tearDown(self) -> None:
        for equipment_type in self.equipment_types_created:
            delete_equipment_type_with_equipments(
                client=self.client, equipment_type=equipment_type
            )
        for location_type in self.location_types_created:
            delete_location_type_with_locations(
                client=self.client, location_type=location_type
            )
        for service_type in self.service_types_created:
            delete_service_type_with_services(
                client=self.client, service_type=service_type
            )
        customers = get_all_customers(self.client)
        for customer in customers:
            delete_customer(self.client, customer)
        delete_equipment_port_type(
            client=self.client, equipment_port_type_id=self.port_type1.id
        )

    def test_service_created(self) -> None:
        service = add_service(
            client=self.client,
            name="Room 201 Internet Access",
            external_id="S3232",
            service_type="Internet Access",
            customer=None,
            properties_dict={"Address Family": "v4"},
            links=[],
        )
        fetch_service = get_service(client=self.client, id=service.id)
        self.assertEqual(service, fetch_service)

    def test_service_with_topology_created(self) -> None:
        location = add_location(
            client=self.client,
            location_hirerchy=[("Room", "Room 201")],
            properties_dict={"Contact": "user@google.com"},
            lat=10,
            long=20,
        )
        router1 = add_equipment(
            client=self.client,
            name="TPLinkRouter1",
            equipment_type="Tp-Link T1600G",
            location=location,
            properties_dict={"IP": "192.688.0.1"},
        )
        router2 = add_equipment(
            client=self.client,
            name="TPLinkRouter2",
            equipment_type="Tp-Link T1600G",
            location=location,
            properties_dict={"IP": "192.688.0.2"},
        )
        router3 = add_equipment(
            client=self.client,
            name="TPLinkRouter3",
            equipment_type="Tp-Link T1600G",
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
        endpoint_port = get_port(
            client=self.client, equipment=router1, port_name="Port 2"
        )
        service = add_service(
            self.client,
            name="Room 201 Internet Access",
            external_id="S3232",
            service_type="Internet Access",
            customer=None,
            properties_dict={"Address Family": "v4"},
            links=[link1, link2],
        )
        add_service_endpoint(
            client=self.client,
            service=service,
            port=endpoint_port,
            role=ServiceEndpointRole.CONSUMER,
        )
        service = get_service(client=self.client, id=service.id)

        self.assertEqual([endpoint_port.id], [e.port.id for e in service.endpoints])
        self.assertEqual([link1.id, link2.id], [s.id for s in service.links])

    def test_service_with_customer_created(self) -> None:
        customer = add_customer(client=self.client, name="Donald", external_id="S322")
        service = add_service(
            client=self.client,
            name="Room 201 Internet Access",
            external_id=None,
            service_type="Internet Access",
            customer=customer,
            properties_dict={},
            links=[],
        )
        self.assertEqual(customer, service.customer)
