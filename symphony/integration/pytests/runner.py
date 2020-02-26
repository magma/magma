#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from unittest import TestLoader, TestSuite, TextTestRunner

from test_equipment import TestEquipment
from test_link import TestLink
from test_location import TestLocation
from test_port_type import TestEquipmentPortType
from test_service import TestService
from test_site_survey import TestSiteSurvey
from utils.constant import XML_OUTPUT_DIRECTORY
from xmlrunner import XMLTestRunner


TEST_CASES = [
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
    runner.run(suite)
