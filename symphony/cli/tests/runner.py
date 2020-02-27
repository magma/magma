#!/usr/bin/env python3
# pyre-strict
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import sys
from typing import List
from unittest import TestCase, TestLoader, TestSuite, TextTestRunner

from pyinventory_tests.test_equipment import TestEquipment
from pyinventory_tests.test_link import TestLink
from pyinventory_tests.test_location import TestLocation
from pyinventory_tests.test_port_type import TestEquipmentPortType
from pyinventory_tests.test_service import TestService
from pyinventory_tests.test_site_survey import TestSiteSurvey
from pyinventory_tests.utils.constant import XML_OUTPUT_DIRECTORY
from xmlrunner import XMLTestRunner


TEST_CASES: List[TestCase] = [
    TestLocation,
    TestEquipment,
    TestLink,
    TestService,
    TestSiteSurvey,
    TestEquipmentPortType,
]

if __name__ == "__main__":
    suite = TestSuite()
    loader = TestLoader()
    for test_class in TEST_CASES:
        tests = loader.loadTestsFromTestCase(test_class)
        suite.addTests(tests)

    if XML_OUTPUT_DIRECTORY:
        runner = XMLTestRunner(output=XML_OUTPUT_DIRECTORY, verbosity=2)
    else:
        runner = TextTestRunner(verbosity=2)
    result = runner.run(suite)
    if len(result.errors) != 0 or len(result.failures) != 0:
        sys.exit(1)
    sys.exit(0)
