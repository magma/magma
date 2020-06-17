#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import Optional
from unittest import TestSuite
from unittest.loader import TestLoader

from grpc import insecure_channel

from .utils import get_grpc_server_address, init_client, wait_for_platform
from .utils.constant import TEST_USER_EMAIL


def load_tests(
    loader: TestLoader, tests: TestSuite, pattern: Optional[str]
) -> TestSuite:

    from .test_equipment import TestEquipment
    from .test_equipment_type import TestEquipmentType
    from .test_link import TestLink
    from .test_location import TestLocation
    from .test_port_type import TestEquipmentPortType
    from .test_service import TestService
    from .test_site_survey import TestSiteSurvey
    from .test_user import TestUser
    from .grpc.rpc_pb2_grpc import TenantServiceStub

    TESTS = [
        TestEquipment,
        TestEquipmentType,
        TestLink,
        TestLocation,
        TestEquipmentPortType,
        TestService,
        TestSiteSurvey,
        TestUser,
    ]

    print("Waiting for symphony to be ready")
    wait_for_platform()
    print("Initializing client")
    client = init_client(TEST_USER_EMAIL, TEST_USER_EMAIL)
    print("Initializing cleaner")
    address = get_grpc_server_address()
    channel = insecure_channel(address)
    stub = TenantServiceStub(channel)
    print("Packing tests")
    test_suite = TestSuite()
    for test_class in TESTS:
        testCaseNames = loader.getTestCaseNames(test_class)
        for test_case_name in testCaseNames:
            test_suite.addTest(test_class(test_case_name, client, stub))
    return test_suite
