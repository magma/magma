#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import argparse
import sys
from unittest import TestLoader, TextTestRunner

import pyinventory_tests
from pyinventory_tests.utils.constant import TESTS_PATTERN, XML_OUTPUT_DIRECTORY
from xmlrunner import XMLTestRunner


if __name__ == "__main__":

    parser = argparse.ArgumentParser(description="test")
    parser.add_argument("-p", "--pattern", default=TESTS_PATTERN)
    parser.add_argument("-o", "--output", default=XML_OUTPUT_DIRECTORY)
    parser.add_argument("-l", "--local", default=False, action="store_true")

    args = parser.parse_args()

    pyinventory_tests.utils.RUN_LOCALLY = args.local

    loader = TestLoader()
    loader.testNamePatterns = [args.pattern]
    suite = loader.loadTestsFromModule(pyinventory_tests)

    if args.output:
        runner = XMLTestRunner(output=args.output, verbosity=2)
    else:
        runner = TextTestRunner(verbosity=2)
    result = runner.run(suite)
    if len(result.errors) != 0 or len(result.failures) != 0:
        sys.exit(1)
    sys.exit(0)
