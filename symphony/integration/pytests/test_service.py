#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


from pyinventory.consts import ServiceEndpointRole
from utils.base_test import BaseTest


class TestService(BaseTest):
    def setUp(self):
        super().setUp()
        self.service_types_created = []
        self.service_types_created.append(
            self.client.add_service_type(
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
            self.client.add_location_type("Room", [("Contact", "email", None, True)])
        )
        self.equipment_types_created = []
        self.equipment_types_created.append(
            self.client.add_equipment_type(
                "Tp-Link T1600G",
                "Router",
                [("IP", "string", None, True)],
                {"Port 1": "Eth", "Port 2": "Eth"},
                [],
            )
        )

    def tearDown(self):
        for equipment_type in self.equipment_types_created:
            self.client.delete_equipment_type_with_equipments(equipment_type)
        for location_type in self.location_types_created:
            self.client.delete_location_type_with_locations(location_type)
        for service_type in self.service_types_created:
            self.client.delete_service_type_with_services(service_type)
        customers = self.client.get_all_customers()
        for customer in customers:
            self.client.delete_customer(customer)

    def test_service_created(self):
        service = self.client.add_service(
            "Room 201 Internet Access",
            "S3232",
            "Internet Access",
            None,
            {"Address Family": "v4"},
            [],
        )
        fetch_service = self.client.get_service(service.id)
        self.assertEqual(service, fetch_service)

    def test_service_with_topology_created(self):
        location = self.client.add_location(
            [("Room", "Room 201")], {"Contact": "user@google.com"}, 10, 20
        )
        router1 = self.client.add_equipment(
            "TPLinkRouter1", "Tp-Link T1600G", location, {"IP": "192.688.0.1"}
        )
        router2 = self.client.add_equipment(
            "TPLinkRouter2", "Tp-Link T1600G", location, {"IP": "192.688.0.2"}
        )
        router3 = self.client.add_equipment(
            "TPLinkRouter3", "Tp-Link T1600G", location, {"IP": "192.688.0.3"}
        )
        link1 = self.client.add_link(router1, "Port 1", router2, "Port 1")
        link2 = self.client.add_link(router2, "Port 2", router3, "Port 1")
        endpoint_port = self.client.get_port(router1, "Port 2")
        service = self.client.add_service(
            name="Room 201 Internet Access",
            external_id="S3232",
            service_type="Internet Access",
            customer=None,
            properties_dict={"Address Family": "v4"},
            links=[link1, link2],
        )
        self.client.add_service_endpoint(
            service, endpoint_port, ServiceEndpointRole.CONSUMER
        )
        service = self.client.get_service(service.id)

        self.assertEqual([endpoint_port], [e.port for e in service.endpoints])
        self.assertEqual([link1, link2], service.links)

    def test_service_with_customer_created(self):
        customer = self.client.add_customer("Donald", "S322")
        service = self.client.add_service(
            name="Room 201 Internet Access",
            external_id=None,
            service_type="Internet Access",
            customer=customer,
            properties_dict={},
            links=[],
        )
        self.assertEqual(customer, service.customer)
