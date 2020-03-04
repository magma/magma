#!/usr/bin/env python3
# pyre-strict
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
from pyinventory.api.service import (
    add_service,
    add_service_endpoint,
    add_service_type,
    delete_service_type_with_services,
    get_service,
)
from pyinventory.graphql.service_endpoint_role_enum import ServiceEndpointRole

from .utils.base_test import BaseTest


class TestService(BaseTest):
    def setUp(self) -> None:
        super().setUp()
        self.service_types_created = []
        self.service_types_created.append(
            add_service_type(
                self.client,
                "Internet Access",
                True,
                [
                    ("Service Package", "string", "Public 5G", True),
                    ("Address Family", "string", None, True),
                ],
            )
        )
        self.location_types_created = []
        self.location_types_created.append(
            add_location_type(self.client, "Room", [("Contact", "email", None, True)])
        )
        self.equipment_types_created = []
        self.equipment_types_created.append(
            add_equipment_type(
                self.client,
                "Tp-Link T1600G",
                "Router",
                [("IP", "string", None, True)],
                {"Port 1": "Eth", "Port 2": "Eth"},
                [],
            )
        )

    def tearDown(self) -> None:
        for equipment_type in self.equipment_types_created:
            delete_equipment_type_with_equipments(self.client, equipment_type)
        for location_type in self.location_types_created:
            delete_location_type_with_locations(self.client, location_type)
        for service_type in self.service_types_created:
            delete_service_type_with_services(self.client, service_type)
        customers = get_all_customers(self.client)
        for customer in customers:
            delete_customer(self.client, customer)

    def test_service_created(self) -> None:
        service = add_service(
            self.client,
            "Room 201 Internet Access",
            "S3232",
            "Internet Access",
            None,
            {"Address Family": "v4"},
            [],
        )
        fetch_service = get_service(self.client, service.id)
        self.assertEqual(service, fetch_service)

    def test_service_with_topology_created(self) -> None:
        location = add_location(
            self.client, [("Room", "Room 201")], {"Contact": "user@google.com"}, 10, 20
        )
        router1 = add_equipment(
            self.client,
            "TPLinkRouter1",
            "Tp-Link T1600G",
            location,
            {"IP": "192.688.0.1"},
        )
        router2 = add_equipment(
            self.client,
            "TPLinkRouter2",
            "Tp-Link T1600G",
            location,
            {"IP": "192.688.0.2"},
        )
        router3 = add_equipment(
            self.client,
            "TPLinkRouter3",
            "Tp-Link T1600G",
            location,
            {"IP": "192.688.0.3"},
        )
        link1 = add_link(self.client, router1, "Port 1", router2, "Port 1")
        link2 = add_link(self.client, router2, "Port 2", router3, "Port 1")
        endpoint_port = get_port(self.client, router1, "Port 2")
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
            self.client, service, endpoint_port, ServiceEndpointRole.CONSUMER
        )
        service = get_service(self.client, service.id)

        self.assertEqual([endpoint_port.id], [e.port.id for e in service.endpoints])
        self.assertEqual([link1, link2], service.links)

    def test_service_with_customer_created(self) -> None:
        customer = add_customer(self.client, "Donald", "S322")
        service = add_service(
            self.client,
            name="Room 201 Internet Access",
            external_id=None,
            service_type="Internet Access",
            customer=customer,
            properties_dict={},
            links=[],
        )
        self.assertEqual(customer, service.customer)
