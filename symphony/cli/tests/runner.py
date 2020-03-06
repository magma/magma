#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import sys
from unittest import TestLoader, TextTestRunner

import pyinventory_tests
from pyinventory_tests.utils.constant import TESTS_PATTERN, XML_OUTPUT_DIRECTORY
from xmlrunner import XMLTestRunner


if __name__ == "__main__":
    loader = TestLoader()
    loader.testNamePatterns = [TESTS_PATTERN]
    suite = loader.loadTestsFromModule(pyinventory_tests)

    if XML_OUTPUT_DIRECTORY:
        runner = XMLTestRunner(output=XML_OUTPUT_DIRECTORY, verbosity=2)
    else:
        runner = TextTestRunner(verbosity=2)
    result = runner.run(suite)
    if len(result.errors) != 0 or len(result.failures) != 0:
        sys.exit(1)
    sys.exit(0)
